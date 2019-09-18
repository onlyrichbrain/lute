// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

// +build !js

package lute

import (
	"bytes"

	"github.com/alecthomas/chroma"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	chromalexers "github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// languagesNoHighlight 中定义的语言不要进行代码语法高亮。这些代码块会在前端进行渲染，比如各种图表。
var languagesNoHighlight = []string{"mermaid", "echarts", "abc"}

// renderCodeBlockHTML 进行代码块 HTML 渲染，实现语法高亮。
func (r *HTMLRenderer) renderCodeBlockHTML(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		tokens := node.tokens
		if nil != node.codeBlockInfo {
			infoWords := bytes.Split(node.codeBlockInfo, items(" "))
			language := string(infoWords[0])
			rendered := false
			if r.option.CodeSyntaxHighlight && !noHighlight(language) {
				rendered = highlightChroma(tokens, language, r)
			}

			if !rendered {
				r.writeString("<pre><code class=\"language-")
				r.writeString(language)
				r.writeString("\">")
				tokens = escapeHTML(tokens)
				r.write(tokens)
			}
		} else {
			rendered := false
			if r.option.CodeSyntaxHighlight {
				rendered = highlightChroma(tokens, "", r)
				if !rendered {
					tokens = escapeHTML(tokens)
					r.write(tokens)
				}
			} else {
				r.writeString("<pre><code>")
				tokens = escapeHTML(tokens)
				r.write(tokens)
			}
		}
		return WalkSkipChildren, nil
	}
	r.writeString("</code></pre>")
	r.newline()
	return WalkContinue, nil
}

func highlightChroma(tokens items, language string, r *HTMLRenderer) (rendered bool) {
	codeBlock := fromItems(tokens)
	var lexer chroma.Lexer
	if "" != language {
		lexer = chromalexers.Get(language)
	} else {
		lexer = chromalexers.Analyse(codeBlock)
	}
	if nil == lexer {
		lexer = chromalexers.Fallback
		language = lexer.Config().Name
	}
	lexer = chroma.Coalesce(lexer)
	iterator, err := lexer.Tokenise(nil, codeBlock)
	if nil == err {
		chromahtmlOpts := []chromahtml.Option{
			chromahtml.PreventSurroundingPre(),
			chromahtml.ClassPrefix("highlight-"),
		}
		if !r.option.CodeSyntaxHighlightInlineStyle {
			chromahtmlOpts = append(chromahtmlOpts, chromahtml.WithClasses())
		}
		if r.option.CodeSyntaxHighlightLineNum {
			chromahtmlOpts = append(chromahtmlOpts, chromahtml.WithLineNumbers())
		}
		formatter := chromahtml.New(chromahtmlOpts...)
		style := styles.Get(r.option.CodeSyntaxHighlightStyleName)
		var b bytes.Buffer
		if err = formatter.Format(&b, style, iterator); nil == err {
			if !r.option.CodeSyntaxHighlightInlineStyle {
				r.writeString("<pre>")
			} else {
				r.writeString("<pre style=\"" + chromahtml.StyleEntryToCSS(style.Get(chroma.Background)) + "\">")
			}
			r.writeString("<code class=\"language-" + language)
			if !r.option.CodeSyntaxHighlightInlineStyle {
				r.writeString(" highlight-chroma")
			}
			r.writeString("\">")
			r.write(b.Bytes())
			rendered = true
		}
	}
	return
}

func noHighlight(language string) bool {
	for _, langNoHighlight := range languagesNoHighlight {
		if language == langNoHighlight {
			return true
		}
	}
	return false
}
