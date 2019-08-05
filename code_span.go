// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

// CodeSpan 描述了代码节点结构。
type CodeSpan struct {
	*BaseNode
}

func (t *Tree) parseCodeSpan(tokens items) (ret Node) {
	startPos := t.context.pos
	length := len(tokens)
	n := 0
	for ; startPos+n < length; n++ {
		if itemBacktick != tokens[startPos+n] {
			break
		}
	}

	backticks := tokens[startPos:startPos+n].rawText()
	if length <= startPos+n {
		t.context.pos += n
		ret = &Text{&BaseNode{typ: NodeText, rawText: backticks, value: backticks}}
		return
	}

	endPos := t.matchCodeSpanEnd(tokens[startPos+n:], n)
	if 1 > endPos {
		t.context.pos += n
		ret = &Text{&BaseNode{typ: NodeText, rawText: backticks, value: backticks}}
		return
	}
	endPos = startPos + endPos + n

	var textTokens = items{}
	for i := startPos + n; i < len(tokens) && i < endPos; i++ {
		token := tokens[i]
		if itemNewline == token {
			textTokens = append(textTokens, itemSpace)
		} else {
			textTokens = append(textTokens, token)
		}
	}

	if 2 < len(textTokens) && itemSpace == textTokens[0] && itemSpace == textTokens[len(textTokens)-1] && !textTokens.isBlankLine() {
		textTokens = textTokens[1:]
		textTokens = textTokens[:len(textTokens)-1]
	}

	baseNode := &BaseNode{typ: NodeCodeSpan, tokens: textTokens, value: textTokens.rawText()}
	ret = &CodeSpan{baseNode}
	t.context.pos = endPos + n

	return
}

func (t *Tree) matchCodeSpanEnd(tokens items, num int) (pos int) {
	length := len(tokens)
	for ; pos < length; {
		len := tokens[pos:].accept(itemBacktick)
		if num == len {
			next := pos + len
			if length-1 > next && itemBacktick == tokens[next] {
				continue
			}

			return pos
		}
		if 0 < len {
			pos += len
		} else {
			pos++
		}
	}

	return -1
}
