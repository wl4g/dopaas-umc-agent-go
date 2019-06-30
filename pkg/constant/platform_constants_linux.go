/**
 * Copyright 2017 ~ 2025 the original author or authors[983708408@qq.com].
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package constant

const (
	// Default network indicators commands.
	DefaultNetPortStateCmd = `
		ss -n sport == #{port} |
		awk 'BEGIN{up=0;down=0;n=0;c1=0;c2=0;c3=0;c4=0;c5=0;c6=0;}
		{up+=$3};
		{down+=$4};
		{n+=1}
		{if($0~"ESTAB") c1+=1};
		{if($0~"CLOSE-WAIT") c2+=1};
		{if($0~"TIME-WAIT") c3+=1};
		{if($0~"CLOSE"&&!$0~"CLOSE-WAIT") c4+=1};
		{if($0~"LISTEN") c5+=1};
		{if($0~"CLOSING") c6+=1};
		END {print up,down,n,c1,c2,c3,c4,c5,c6}'
	`
)
