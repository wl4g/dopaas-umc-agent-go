package proc

import (
	"bytes"
	"debug/dwarf"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-delve/delve/pkg/dwarf/frame"
	"github.com/go-delve/delve/pkg/dwarf/godwarf"
	"github.com/go-delve/delve/pkg/dwarf/line"
	"github.com/go-delve/delve/pkg/dwarf/op"
	"github.com/go-delve/delve/pkg/dwarf/reader"
	"github.com/go-delve/delve/pkg/goversion"
)

// BinaryInfo holds information on the binaries being executed (this
// includes both the executable and also any loaded libraries).
type BinaryInfo struct {
	// Architecture of this binary.
	Arch Arch

	// GOOS operating system this binary is executing on.
	GOOS string

	debugInfoDirectories []string

	// Functions is a list of all DW_TAG_subprogram entries in debug_info, sorted by entry point
	Functions []Function
	// Sources is a list of all source files found in debug_line.
	Sources []string
	// LookupFunc maps function names to a description of the function.
	LookupFunc map[string]*Function

	// Images is a list of loaded shared libraries (also known as
	// shared objects on linux or DLLs on windows).
	Images []*Image

	ElfDynamicSection ElfDynamicSection

	lastModified time.Time // Time the executable of this process was last modified

	closer         io.Closer
	sepDebugCloser io.Closer

	staticBase uint64

	// Maps package names to package paths, needed to lookup types inside DWARF info
	packageMap map[string]string

	frameEntries frame.FrameDescriptionEntries
	compileUnits []*compileUnit
	types        map[string]dwarfRef
	packageVars  []packageVar // packageVars is a list of all global/package variables in debug_info, sorted by address

	gStructOffset uint64

	// nameOfRuntimeType maps an address of a runtime._type struct to its
	// decoded name. Used with versions of Go <= 1.10 to figure out the DIE of
	// the concrete type of interfaces.
	nameOfRuntimeType map[uintptr]nameOfRuntimeTypeEntry

	// consts[off] lists all the constants with the type defined at offset off.
	consts constantsMap

	loadErrMu sync.Mutex
	loadErr   error

	initialized bool
}

// ErrUnsupportedLinuxArch is returned when attempting to debug a binary compiled for an unsupported architecture.
var ErrUnsupportedLinuxArch = errors.New("unsupported architecture - only linux/amd64 is supported")

// ErrUnsupportedWindowsArch is returned when attempting to debug a binary compiled for an unsupported architecture.
var ErrUnsupportedWindowsArch = errors.New("unsupported architecture of windows/386 - only windows/amd64 is supported")

// ErrUnsupportedDarwinArch is returned when attempting to debug a binary compiled for an unsupported architecture.
var ErrUnsupportedDarwinArch = errors.New("unsupported architecture - only darwin/amd64 is supported")

// ErrCouldNotDetermineRelocation is an error returned when Delve could not determine the base address of a
// position independant executable.
var ErrCouldNotDetermineRelocation = errors.New("could not determine the base address of a PIE")

// ErrNoDebugInfoFound is returned when Delve cannot open the debug_info
// section or find an external debug info file.
var ErrNoDebugInfoFound = errors.New("could not open debug info")

const dwarfGoLanguage = 22 // DW_LANG_Go (from DWARF v5, section 7.12, page 231)

type compileUnit struct {
	name   string // univocal name for non-go compile units
	lowPC  uint64
	ranges [][2]uint64

	entry              *dwarf.Entry        // debug_info entry describing this compile unit
	isgo               bool                // true if this is the go compile unit
	lineInfo           *line.DebugLineInfo // debug_line segment associated with this compile unit
	concreteInlinedFns []inlinedFn         // list of concrete inlined functions within this compile unit
	optimized          bool                // this compile unit is optimized
	producer           string              // producer attribute

	startOffset, endOffset dwarf.Offset // interval of offsets contained in this compile unit

	image *Image // parent image of this compilation unit.
}

// dwarfRef is a reference to a Debug Info Entry inside a shared object.
type dwarfRef struct {
	imageIndex int
	offset     dwarf.Offset
}

// inlinedFn represents a concrete inlined function, e.g.
// an entry for the generated code of an inlined function.
type inlinedFn struct {
	Name          string    // Name of the function that was inlined
	LowPC, HighPC uint64    // Address range of the generated inlined instructions
	CallFile      string    // File of the call site of the inlined function
	CallLine      int64     // Line of the call site of the inlined function
	Parent        *Function // The function that contains this inlined function
}

// Function describes a function in the target program.
type Function struct {
	Name       string
	Entry, End uint64 // same as DW_AT_lowpc and DW_AT_highpc
	offset     dwarf.Offset
	cu         *compileUnit
}

// PackageName returns the package part of the symbol name,
// or the empty string if there is none.
// Borrowed from $GOROOT/debug/gosym/symtab.go
func (fn *Function) PackageName() string {
	return packageName(fn.Name)
}

func packageName(name string) string {
	pathend := strings.LastIndex(name, "/")
	if pathend < 0 {
		pathend = 0
	}

	if i := strings.Index(name[pathend:], "."); i != -1 {
		return name[:pathend+i]
	}
	return ""
}

// ReceiverName returns the receiver type name of this symbol,
// or the empty string if there is none.
// Borrowed from $GOROOT/debug/gosym/symtab.go
func (fn *Function) ReceiverName() string {
	pathend := strings.LastIndex(fn.Name, "/")
	if pathend < 0 {
		pathend = 0
	}
	l := strings.Index(fn.Name[pathend:], ".")
	r := strings.LastIndex(fn.Name[pathend:], ".")
	if l == -1 || r == -1 || l == r {
		return ""
	}
	return fn.Name[pathend+l+1 : pathend+r]
}

// BaseName returns the symbol name without the package or receiver name.
// Borrowed from $GOROOT/debug/gosym/symtab.go
func (fn *Function) BaseName() string {
	if i := strings.LastIndex(fn.Name, "."); i != -1 {
		return fn.Name[i+1:]
	}
	return fn.Name
}

// Optimized returns true if the function was optimized by the compiler.
func (fn *Function) Optimized() bool {
	return fn.cu.optimized
}

type constantsMap map[dwarfRef]*constantType

type constantType struct {
	initialized bool
	values      []constantValue
}

type constantValue struct {
	name      string
	fullName  string
	value     int64
	singleBit bool
}

// packageVar represents a package-level variable (or a C global variable).
// If a global variable does not have an address (for example it's stored in
// a register, or non-contiguously) addr will be 0.
type packageVar struct {
	name   string
	cu     *compileUnit
	offset dwarf.Offset
	addr   uint64
}

type loclistReader struct {
	data  []byte
	cur   int
	ptrSz int
}

func (rdr *loclistReader) Seek(off int) {
	rdr.cur = off
}

func (rdr *loclistReader) read(sz int) []byte {
	r := rdr.data[rdr.cur : rdr.cur+sz]
	rdr.cur += sz
	return r
}

func (rdr *loclistReader) oneAddr() uint64 {
	switch rdr.ptrSz {
	case 4:
		addr := binary.LittleEndian.Uint32(rdr.read(rdr.ptrSz))
		if addr == ^uint32(0) {
			return ^uint64(0)
		}
		return uint64(addr)
	case 8:
		addr := uint64(binary.LittleEndian.Uint64(rdr.read(rdr.ptrSz)))
		return addr
	default:
		panic("bad address size")
	}
}

func (rdr *loclistReader) Next(e *loclistEntry) bool {
	e.lowpc = rdr.oneAddr()
	e.highpc = rdr.oneAddr()

	if e.lowpc == 0 && e.highpc == 0 {
		return false
	}

	if e.BaseAddressSelection() {
		e.instr = nil
		return true
	}

	instrlen := binary.LittleEndian.Uint16(rdr.read(2))
	e.instr = rdr.read(int(instrlen))
	return true
}

type loclistEntry struct {
	lowpc, highpc uint64
	instr         []byte
}

type runtimeTypeDIE struct {
	offset dwarf.Offset
	kind   int64
}

func (e *loclistEntry) BaseAddressSelection() bool {
	return e.lowpc == ^uint64(0)
}

type buildIDHeader struct {
	Namesz uint32
	Descsz uint32
	Type   uint32
}

// ElfDynamicSection describes the .dynamic section of an ELF executable.
type ElfDynamicSection struct {
	Addr uint64 // relocated address of where the .dynamic section is mapped in memory
	Size uint64 // size of the .dynamic section of the executable
}

// NewBinaryInfo returns an initialized but unloaded BinaryInfo struct.
func NewBinaryInfo(goos, goarch string) *BinaryInfo {
	r := &BinaryInfo{GOOS: goos, nameOfRuntimeType: make(map[uintptr]nameOfRuntimeTypeEntry)}

	// TODO: find better way to determine proc arch (perhaps use executable file info).
	switch goarch {
	case "amd64":
		r.Arch = AMD64Arch(goos)
	}

	return r
}

// LoadBinaryInfo will load and store the information from the binary at 'path'.
func (bi *BinaryInfo) LoadBinaryInfo(path string, entryPoint uint64, debugInfoDirs []string) error {
	fi, err := os.Stat(path)
	if err == nil {
		bi.lastModified = fi.ModTime()
	}

	bi.debugInfoDirectories = debugInfoDirs

	return bi.AddImage(path, entryPoint)
}

func loadBinaryInfo(bi *BinaryInfo, image *Image, path string, entryPoint uint64) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	switch bi.GOOS {
	case "linux":
		return loadBinaryInfoElf(bi, image, path, entryPoint, &wg)
	case "windows":
		return loadBinaryInfoPE(bi, image, path, entryPoint, &wg)
	case "darwin":
		return loadBinaryInfoMacho(bi, image, path, entryPoint, &wg)
	}
	return errors.New("unsupported operating system")
}

// GStructOffset returns the offset of the G
// struct in thread local storage.
func (bi *BinaryInfo) GStructOffset() uint64 {
	return bi.gStructOffset
}

// LastModified returns the last modified time of the binary.
func (bi *BinaryInfo) LastModified() time.Time {
	return bi.lastModified
}

// DwarfReader returns a reader for the dwarf data
func (so *Image) DwarfReader() *reader.Reader {
	return reader.New(so.dwarf)
}

// Types returns list of types present in the debugged program.
func (bi *BinaryInfo) Types() ([]string, error) {
	types := make([]string, 0, len(bi.types))
	for k := range bi.types {
		types = append(types, k)
	}
	return types, nil
}

// PCToLine converts an instruction address to a file/line/function.
func (bi *BinaryInfo) PCToLine(pc uint64) (string, int, *Function) {
	fn := bi.PCToFunc(pc)
	if fn == nil {
		return "", 0, nil
	}
	f, ln := fn.cu.lineInfo.PCToLine(fn.Entry, pc)
	return f, ln, fn
}

// LineToPC converts a file:line into a memory address.
func (bi *BinaryInfo) LineToPC(filename string, lineno int) (pc uint64, fn *Function, err error) {
	fileFound := false
	for _, cu := range bi.compileUnits {
		if cu.lineInfo.Lookup[filename] != nil {
			fileFound = true
			pc := cu.lineInfo.LineToPC(filename, lineno)
			if pc == 0 {
				// Check to see if this file:line belongs to the call site
				// of an inlined function.
				for _, ifn := range cu.concreteInlinedFns {
					if strings.Contains(ifn.CallFile, filename) && ifn.CallLine == int64(lineno) {
						return ifn.LowPC, ifn.Parent, nil
					}
				}
			}
			if fn := bi.PCToFunc(pc); fn != nil {
				return pc, fn, nil
			}
		}
	}
	if fileFound {
		return 0, nil, fmt.Errorf("could not find statement at %s:%d, please use a line with a statement", filename, lineno)
	} else {
		return 0, nil, fmt.Errorf("could not find file %s", filename)
	}
}

// AllPCsForFileLine returns all PC addresses for the given filename:lineno.
func (bi *BinaryInfo) AllPCsForFileLine(filename string, lineno int) []uint64 {
	r := make([]uint64, 0, 1)
	for _, cu := range bi.compileUnits {
		if cu.lineInfo.Lookup[filename] != nil {
			r = append(r, cu.lineInfo.AllPCsForFileLine(filename, lineno)...)
		}
	}
	return r
}

// PCToFunc returns the function containing the given PC address
func (bi *BinaryInfo) PCToFunc(pc uint64) *Function {
	i := sort.Search(len(bi.Functions), func(i int) bool {
		fn := bi.Functions[i]
		return pc <= fn.Entry || (fn.Entry <= pc && pc < fn.End)
	})
	if i != len(bi.Functions) {
		fn := &bi.Functions[i]
		if fn.Entry <= pc && pc < fn.End {
			return fn
		}
	}
	return nil
}

// pcToImage returns the image containing the given PC address.
func (bi *BinaryInfo) pcToImage(pc uint64) *Image {
	fn := bi.PCToFunc(pc)
	return bi.funcToImage(fn)
}

// Image represents a loaded library file (shared object on linux, DLL on windows).
type Image struct {
	Path       string
	StaticBase uint64
	addr       uint64

	index int // index of this object in BinaryInfo.SharedObjects

	closer         io.Closer
	sepDebugCloser io.Closer

	dwarf       *dwarf.Data
	dwarfReader *dwarf.Reader
	loclist     loclistReader

	typeCache map[dwarf.Offset]godwarf.Type

	// runtimeTypeToDIE maps between the offset of a runtime._type in
	// runtime.moduledata.types and the offset of the DIE in debug_info. This
	// map is filled by using the extended attribute godwarf.AttrGoRuntimeType
	// which was added in go 1.11.
	runtimeTypeToDIE map[uint64]runtimeTypeDIE

	loadErrMu sync.Mutex
	loadErr   error
}

// AddImage adds the specified image to bi, loading data asynchronously.
// Addr is the relocated entry point for the executable and staticBase (i.e.
// the relocation offset) for all other images.
// The first image added must be the executable file.
func (bi *BinaryInfo) AddImage(path string, addr uint64) error {
	// Check if the image is already present.
	if len(bi.Images) > 0 && !strings.HasPrefix(path, "/") {
		return nil
	}
	for _, image := range bi.Images {
		if image.Path == path && image.addr == addr {
			return nil
		}
	}

	// Actually add the image.
	image := &Image{Path: path, addr: addr, typeCache: make(map[dwarf.Offset]godwarf.Type)}
	// add Image regardless of error so that we don't attempt to re-add it every time we stop
	image.index = len(bi.Images)
	bi.Images = append(bi.Images, image)
	err := loadBinaryInfo(bi, image, path, addr)
	if err != nil {
		bi.Images[len(bi.Images)-1].loadErr = err
	}
	return err
}

// moduleDataToImage finds the image corresponding to the given module data object.
func (bi *BinaryInfo) moduleDataToImage(md *moduleData) *Image {
	return bi.funcToImage(bi.PCToFunc(uint64(md.text)))
}

// imageToModuleData finds the module data in mds corresponding to the given image.
func (bi *BinaryInfo) imageToModuleData(image *Image, mds []moduleData) *moduleData {
	for _, md := range mds {
		im2 := bi.moduleDataToImage(&md)
		if im2.index == image.index {
			return &md
		}
	}
	return nil
}

// typeToImage returns the image containing the give type.
func (bi *BinaryInfo) typeToImage(typ godwarf.Type) *Image {
	return bi.Images[typ.Common().Index]
}

var errBinaryInfoClose = errors.New("multiple errors closing executable files")

// Close closes all internal readers.
func (bi *BinaryInfo) Close() error {
	var errs []error
	for _, image := range bi.Images {
		if err := image.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errBinaryInfoClose
	}
}

func (image *Image) Close() error {
	var err1, err2 error
	if image.sepDebugCloser != nil {
		err := image.sepDebugCloser.Close()
		if err != nil {
			err1 = fmt.Errorf("closing shared object %q (split dwarf): %v", image.Path, err)
		}
	}
	if image.closer != nil {
		err := image.closer.Close()
		if err != nil {
			err2 = fmt.Errorf("closing shared object %q: %v", image.Path, err)
		}
	}
	if err1 != nil && err2 != nil {
		return errBinaryInfoClose
	}
	if err1 != nil {
		return err1
	}
	return err2
}

func (image *Image) setLoadError(fmtstr string, args ...interface{}) {
	image.loadErrMu.Lock()
	image.loadErr = fmt.Errorf(fmtstr, args...)
	image.loadErrMu.Unlock()
}

// LoadError returns any error incurred while loading this image.
func (image *Image) LoadError() error {
	return image.loadErr
}

type nilCloser struct{}

func (c *nilCloser) Close() error { return nil }

// LoadImageFromData creates a new Image, using the specified data, and adds it to bi.
// This is used for debugging BinaryInfo, you should use LoadBinary instead.
func (bi *BinaryInfo) LoadImageFromData(dwdata *dwarf.Data, debugFrameBytes, debugLineBytes, debugLocBytes []byte) {
	image := &Image{}
	image.closer = (*nilCloser)(nil)
	image.sepDebugCloser = (*nilCloser)(nil)
	image.dwarf = dwdata
	image.typeCache = make(map[dwarf.Offset]godwarf.Type)

	if debugFrameBytes != nil {
		bi.frameEntries = frame.Parse(debugFrameBytes, frame.DwarfEndian(debugFrameBytes), 0)
	}

	image.loclistInit(debugLocBytes, bi.Arch.PtrSize())

	bi.loadDebugInfoMaps(image, debugLineBytes, nil, nil)

	bi.Images = append(bi.Images, image)
}

func (image *Image) loclistInit(data []byte, ptrSz int) {
	image.loclist.data = data
	image.loclist.ptrSz = ptrSz
}

func (bi *BinaryInfo) locationExpr(entry reader.Entry, attr dwarf.Attr, pc uint64) ([]byte, string, error) {
	a := entry.Val(attr)
	if a == nil {
		return nil, "", fmt.Errorf("no location attribute %s", attr)
	}
	if instr, ok := a.([]byte); ok {
		var descr bytes.Buffer
		fmt.Fprintf(&descr, "[block] ")
		op.PrettyPrint(&descr, instr)
		return instr, descr.String(), nil
	}
	off, ok := a.(int64)
	if !ok {
		return nil, "", fmt.Errorf("could not interpret location attribute %s", attr)
	}
	instr := bi.loclistEntry(off, pc)
	if instr == nil {
		return nil, "", fmt.Errorf("could not find loclist entry at %#x for address %#x", off, pc)
	}
	var descr bytes.Buffer
	fmt.Fprintf(&descr, "[%#x:%#x] ", off, pc)
	op.PrettyPrint(&descr, instr)
	return instr, descr.String(), nil
}

// Location returns the location described by attribute attr of entry.
// This will either be an int64 address or a slice of Pieces for locations
// that don't correspond to a single memory address (registers, composite
// locations).
func (bi *BinaryInfo) Location(entry reader.Entry, attr dwarf.Attr, pc uint64, regs op.DwarfRegisters) (int64, []op.Piece, string, error) {
	instr, descr, err := bi.locationExpr(entry, attr, pc)
	if err != nil {
		return 0, nil, "", err
	}
	addr, pieces, err := op.ExecuteStackProgram(regs, instr)
	return addr, pieces, descr, err
}

// loclistEntry returns the loclist entry in the loclist starting at off,
// for address pc.
func (bi *BinaryInfo) loclistEntry(off int64, pc uint64) []byte {
	var base uint64
	image := bi.Images[0]
	if cu := bi.findCompileUnit(pc); cu != nil {
		base = cu.lowPC
		image = cu.image
	}
	if image == nil || image.loclist.data == nil {
		return nil
	}

	image.loclist.Seek(int(off))
	var e loclistEntry
	for image.loclist.Next(&e) {
		if e.BaseAddressSelection() {
			base = e.highpc
			continue
		}
		if pc >= e.lowpc+base && pc < e.highpc+base {
			return e.instr
		}
	}

	return nil
}

// findCompileUnit returns the compile unit containing address pc.
func (bi *BinaryInfo) findCompileUnit(pc uint64) *compileUnit {
	for _, cu := range bi.compileUnits {
		for _, rng := range cu.ranges {
			if pc >= rng[0] && pc < rng[1] {
				return cu
			}
		}
	}
	return nil
}

func (bi *BinaryInfo) findCompileUnitForOffset(off dwarf.Offset) *compileUnit {
	for _, cu := range bi.compileUnits {
		if off >= cu.startOffset && off < cu.endOffset {
			return cu
		}
	}
	return nil
}

// Producer returns the value of DW_AT_producer.
func (bi *BinaryInfo) Producer() string {
	for _, cu := range bi.compileUnits {
		if cu.isgo && cu.producer != "" {
			return cu.producer
		}
	}
	return ""
}

// Type returns the Dwarf type entry at `offset`.
func (image *Image) Type(offset dwarf.Offset) (godwarf.Type, error) {
	return godwarf.ReadType(image.dwarf, image.index, offset, image.typeCache)
}

// funcToImage returns the Image containing function fn, or the
// executable file as a fallback.
func (bi *BinaryInfo) funcToImage(fn *Function) *Image {
	if fn == nil {
		return bi.Images[0]
	}
	return fn.cu.image
}

// ELF ///////////////////////////////////////////////////////////////

// ErrNoBuildIDNote is used in openSeparateDebugInfo to signal there's no
// build-id note on the binary, so LoadBinaryInfoElf will return
// the error message coming from elfFile.DWARF() instead.
type ErrNoBuildIDNote struct{}

func (e *ErrNoBuildIDNote) Error() string {
	return "can't find build-id note on binary"
}

// openSeparateDebugInfo searches for a file containing the separate
// debug info for the binary using the "build ID" method as described
// in GDB's documentation [1], and if found returns two handles, one
// for the bare file, and another for its corresponding elf.File.
// [1] https://sourceware.org/gdb/onlinedocs/gdb/Separate-Debug-Files.html
//
// Alternatively, if the debug file cannot be found be the build-id, Delve
// will look in directories specified by the debug-info-directories config value.
func (bi *BinaryInfo) openSeparateDebugInfo(image *Image, exe *elf.File, debugInfoDirectories []string) (*os.File, *elf.File, error) {
	var debugFilePath string
	for _, dir := range debugInfoDirectories {
		var potentialDebugFilePath string
		if strings.Contains(dir, "build-id") {
			desc1, desc2, err := parseBuildID(exe)
			if err != nil {
				continue
			}
			potentialDebugFilePath = fmt.Sprintf("%s/%s/%s.debug", dir, desc1, desc2)
		} else {
			potentialDebugFilePath = fmt.Sprintf("%s/%s.debug", dir, filepath.Base(image.Path))
		}
		_, err := os.Stat(potentialDebugFilePath)
		if err == nil {
			debugFilePath = potentialDebugFilePath
			break
		}
	}
	if debugFilePath == "" {
		return nil, nil, ErrNoDebugInfoFound
	}
	sepFile, err := os.OpenFile(debugFilePath, 0, os.ModePerm)
	if err != nil {
		return nil, nil, errors.New("can't open separate debug file: " + err.Error())
	}

	elfFile, err := elf.NewFile(sepFile)
	if err != nil {
		sepFile.Close()
		return nil, nil, fmt.Errorf("can't open separate debug file %q: %v", debugFilePath, err.Error())
	}

	if elfFile.Machine != elf.EM_X86_64 {
		sepFile.Close()
		return nil, nil, fmt.Errorf("can't open separate debug file %q: %v", debugFilePath, ErrUnsupportedLinuxArch.Error())
	}

	return sepFile, elfFile, nil
}

func parseBuildID(exe *elf.File) (string, string, error) {
	buildid := exe.Section(".note.gnu.build-id")
	if buildid == nil {
		return "", "", &ErrNoBuildIDNote{}
	}

	br := buildid.Open()
	bh := new(buildIDHeader)
	if err := binary.Read(br, binary.LittleEndian, bh); err != nil {
		return "", "", errors.New("can't read build-id header: " + err.Error())
	}

	name := make([]byte, bh.Namesz)
	if err := binary.Read(br, binary.LittleEndian, name); err != nil {
		return "", "", errors.New("can't read build-id name: " + err.Error())
	}

	if strings.TrimSpace(string(name)) != "GNU\x00" {
		return "", "", errors.New("invalid build-id signature")
	}

	descBinary := make([]byte, bh.Descsz)
	if err := binary.Read(br, binary.LittleEndian, descBinary); err != nil {
		return "", "", errors.New("can't read build-id desc: " + err.Error())
	}
	desc := hex.EncodeToString(descBinary)
	return desc[:2], desc[2:], nil
}

// loadBinaryInfoElf specifically loads information from an ELF binary.
func loadBinaryInfoElf(bi *BinaryInfo, image *Image, path string, addr uint64, wg *sync.WaitGroup) error {
	exe, err := os.OpenFile(path, 0, os.ModePerm)
	if err != nil {
		return err
	}
	image.closer = exe
	elfFile, err := elf.NewFile(exe)
	if err != nil {
		return err
	}
	if elfFile.Machine != elf.EM_X86_64 {
		return ErrUnsupportedLinuxArch
	}

	if image.index == 0 {
		// adding executable file:
		// - addr is entryPoint therefore staticBase needs to be calculated by
		//   subtracting the entry point specified in the executable file from addr.
		// - memory address of the .dynamic section needs to be recorded in
		//   BinaryInfo so that we can find loaded libraries.
		if addr != 0 {
			image.StaticBase = addr - elfFile.Entry
		} else if elfFile.Type == elf.ET_DYN {
			return ErrCouldNotDetermineRelocation
		}
		if dynsec := elfFile.Section(".dynamic"); dynsec != nil {
			bi.ElfDynamicSection.Addr = dynsec.Addr + image.StaticBase
			bi.ElfDynamicSection.Size = dynsec.Size
		}
	} else {
		image.StaticBase = addr
	}

	dwarfFile := elfFile

	image.dwarf, err = elfFile.DWARF()
	if err != nil {
		var sepFile *os.File
		var serr error
		sepFile, dwarfFile, serr = bi.openSeparateDebugInfo(image, elfFile, bi.debugInfoDirectories)
		if serr != nil {
			return serr
		}
		image.sepDebugCloser = sepFile
		image.dwarf, err = dwarfFile.DWARF()
		if err != nil {
			return err
		}
	}

	image.dwarfReader = image.dwarf.Reader()

	debugLineBytes, err := godwarf.GetDebugSectionElf(dwarfFile, "line")
	if err != nil {
		return err
	}
	debugLocBytes, _ := godwarf.GetDebugSectionElf(dwarfFile, "loc")
	image.loclistInit(debugLocBytes, bi.Arch.PtrSize())

	wg.Add(2)
	go bi.parseDebugFrameElf(image, dwarfFile, wg)
	go bi.loadDebugInfoMaps(image, debugLineBytes, wg, nil)
	if image.index == 0 {
		// determine g struct offset only when loading the executable file
		wg.Add(1)
		go bi.setGStructOffsetElf(image, dwarfFile, wg)
	}
	return nil
}

func (bi *BinaryInfo) parseDebugFrameElf(image *Image, exe *elf.File, wg *sync.WaitGroup) {
	defer wg.Done()

	debugFrameData, err := godwarf.GetDebugSectionElf(exe, "frame")
	if err != nil {
		image.setLoadError("could not get .debug_frame section: %v", err)
		return
	}
	debugInfoData, err := godwarf.GetDebugSectionElf(exe, "info")
	if err != nil {
		image.setLoadError("could not get .debug_info section: %v", err)
		return
	}

	bi.frameEntries = bi.frameEntries.Append(frame.Parse(debugFrameData, frame.DwarfEndian(debugInfoData), image.StaticBase))
}

func (bi *BinaryInfo) setGStructOffsetElf(image *Image, exe *elf.File, wg *sync.WaitGroup) {
	defer wg.Done()

	// This is a bit arcane. Essentially:
	// - If the program is pure Go, it can do whatever it wants, and puts the G
	//   pointer at %fs-8.
	// - Otherwise, Go asks the external linker to place the G pointer by
	//   emitting runtime.tlsg, a TLS symbol, which is relocated to the chosen
	//   offset in libc's TLS block.
	symbols, err := exe.Symbols()
	if err != nil {
		image.setLoadError("could not parse ELF symbols: %v", err)
		return
	}
	var tlsg *elf.Symbol
	for _, symbol := range symbols {
		if symbol.Name == "runtime.tlsg" {
			s := symbol
			tlsg = &s
			break
		}
	}
	if tlsg == nil {
		bi.gStructOffset = ^uint64(8) + 1 // -8
		return
	}
	var tls *elf.Prog
	for _, prog := range exe.Progs {
		if prog.Type == elf.PT_TLS {
			tls = prog
			break
		}
	}
	if tls == nil {
		bi.gStructOffset = ^uint64(8) + 1 // -8
		return
	}
	memsz := tls.Memsz

	memsz = (memsz + uint64(bi.Arch.PtrSize()) - 1) & ^uint64(bi.Arch.PtrSize()-1) // align to pointer-sized-boundary
	// The TLS register points to the end of the TLS block, which is
	// tls.Memsz long. runtime.tlsg is an offset from the beginning of that block.
	bi.gStructOffset = ^(memsz) + 1 + tlsg.Value // -tls.Memsz + tlsg.Value
}

// PE ////////////////////////////////////////////////////////////////

const _IMAGE_DLLCHARACTERISTICS_DYNAMIC_BASE = 0x0040

// loadBinaryInfoPE specifically loads information from a PE binary.
func loadBinaryInfoPE(bi *BinaryInfo, image *Image, path string, entryPoint uint64, wg *sync.WaitGroup) error {
	peFile, closer, err := openExecutablePathPE(path)
	if err != nil {
		return err
	}
	image.closer = closer
	if peFile.Machine != pe.IMAGE_FILE_MACHINE_AMD64 {
		return ErrUnsupportedWindowsArch
	}
	image.dwarf, err = peFile.DWARF()
	if err != nil {
		return err
	}

	//TODO(aarzilli): actually test this when Go supports PIE buildmode on Windows.
	opth := peFile.OptionalHeader.(*pe.OptionalHeader64)
	if entryPoint != 0 {
		image.StaticBase = entryPoint - opth.ImageBase
	} else {
		if opth.DllCharacteristics&_IMAGE_DLLCHARACTERISTICS_DYNAMIC_BASE != 0 {
			return ErrCouldNotDetermineRelocation
		}
	}

	image.dwarfReader = image.dwarf.Reader()

	debugLineBytes, err := godwarf.GetDebugSectionPE(peFile, "line")
	if err != nil {
		return err
	}
	debugLocBytes, _ := godwarf.GetDebugSectionPE(peFile, "loc")
	image.loclistInit(debugLocBytes, bi.Arch.PtrSize())

	wg.Add(2)
	go bi.parseDebugFramePE(image, peFile, wg)
	go bi.loadDebugInfoMaps(image, debugLineBytes, wg, nil)

	// Use ArbitraryUserPointer (0x28) as pointer to pointer
	// to G struct per:
	// https://golang.org/src/runtime/cgo/gcc_windows_amd64.c

	bi.gStructOffset = 0x28
	return nil
}

func openExecutablePathPE(path string) (*pe.File, io.Closer, error) {
	f, err := os.OpenFile(path, 0, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}
	peFile, err := pe.NewFile(f)
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	return peFile, f, nil
}

func (bi *BinaryInfo) parseDebugFramePE(image *Image, exe *pe.File, wg *sync.WaitGroup) {
	defer wg.Done()

	debugFrameBytes, err := godwarf.GetDebugSectionPE(exe, "frame")
	if err != nil {
		image.setLoadError("could not get .debug_frame section: %v", err)
		return
	}
	debugInfoBytes, err := godwarf.GetDebugSectionPE(exe, "info")
	if err != nil {
		image.setLoadError("could not get .debug_info section: %v", err)
		return
	}

	bi.frameEntries = bi.frameEntries.Append(frame.Parse(debugFrameBytes, frame.DwarfEndian(debugInfoBytes), image.StaticBase))
}

// Borrowed from https://golang.org/src/cmd/internal/objfile/pe.go
func findPESymbol(f *pe.File, name string) (*pe.Symbol, error) {
	for _, s := range f.Symbols {
		if s.Name != name {
			continue
		}
		if s.SectionNumber <= 0 {
			return nil, fmt.Errorf("symbol %s: invalid section number %d", name, s.SectionNumber)
		}
		if len(f.Sections) < int(s.SectionNumber) {
			return nil, fmt.Errorf("symbol %s: section number %d is larger than max %d", name, s.SectionNumber, len(f.Sections))
		}
		return s, nil
	}
	return nil, fmt.Errorf("no %s symbol found", name)
}

// MACH-O ////////////////////////////////////////////////////////////

// loadBinaryInfoMacho specifically loads information from a Mach-O binary.
func loadBinaryInfoMacho(bi *BinaryInfo, image *Image, path string, entryPoint uint64, wg *sync.WaitGroup) error {
	exe, err := macho.Open(path)
	if err != nil {
		return err
	}
	image.closer = exe
	if exe.Cpu != macho.CpuAmd64 {
		return ErrUnsupportedDarwinArch
	}
	image.dwarf, err = exe.DWARF()
	if err != nil {
		return err
	}

	image.dwarfReader = image.dwarf.Reader()

	debugLineBytes, err := godwarf.GetDebugSectionMacho(exe, "line")
	if err != nil {
		return err
	}
	debugLocBytes, _ := godwarf.GetDebugSectionMacho(exe, "loc")
	image.loclistInit(debugLocBytes, bi.Arch.PtrSize())

	wg.Add(2)
	go bi.parseDebugFrameMacho(image, exe, wg)
	go bi.loadDebugInfoMaps(image, debugLineBytes, wg, bi.setGStructOffsetMacho)
	return nil
}

func (bi *BinaryInfo) setGStructOffsetMacho() {
	// In go1.11 it's 0x30, before 0x8a0, see:
	// https://github.com/golang/go/issues/23617
	// and go commit b3a854c733257c5249c3435ffcee194f8439676a
	producer := bi.Producer()
	if producer != "" && goversion.ProducerAfterOrEqual(producer, 1, 11) {
		bi.gStructOffset = 0x30
		return
	}
	bi.gStructOffset = 0x8a0
}

func (bi *BinaryInfo) parseDebugFrameMacho(image *Image, exe *macho.File, wg *sync.WaitGroup) {
	defer wg.Done()

	debugFrameBytes, err := godwarf.GetDebugSectionMacho(exe, "frame")
	if err != nil {
		image.setLoadError("could not get __debug_frame section: %v", err)
		return
	}
	debugInfoBytes, err := godwarf.GetDebugSectionMacho(exe, "info")
	if err != nil {
		image.setLoadError("could not get .debug_info section: %v", err)
		return
	}

	bi.frameEntries = bi.frameEntries.Append(frame.Parse(debugFrameBytes, frame.DwarfEndian(debugInfoBytes), image.StaticBase))
}
