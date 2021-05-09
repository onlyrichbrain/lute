// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"github.com/88250/lute/ast"
	"testing"

	"github.com/88250/lute"
)

var md2BlockDOMTests = []parseTest{

	{"5", "{{{col\n\n{{{\nfoo\n}}}\n\n}}}\n\nbar\n\n", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeSuperBlock\" class=\"sb\" data-sb-layout=\"col\"><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeSuperBlock\" class=\"sb\" data-sb-layout=\"row\"><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
	{"4", "<u>foo</u>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><u>foo</u></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
	{"3", "* [x] foo", "<div data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li protyle-task--done\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconCheck\"></use></svg></div><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
	{"2", "{{name:foo}}", "<div data-content=\"name:foo\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeBlockQueryEmbed\" class=\"render-node\"></div>"},
	{"1", "<kbd>foo</kbd>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><kbd>foo</kbd></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
	{"0", "<audio src=\"assets/foo\"></audio>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeAudio\" class=\"iframe\"><span class=\"protyle-action\"><svg class=\"svg\"><use xlink:href=\"#iconMore\"></use></svg></span><audio src=\"/siyuan/0/测试笔记/assets/foo\" data-src=\"assets/foo\"></audio>\u200b<div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestMd2BlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetToC(true)
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.ParseOptions.KramdownBlockIAL = true
	luteEngine.RenderOptions.KramdownBlockIAL = true
	luteEngine.ParseOptions.Tag = true
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.SetLinkBase("/siyuan/0/测试笔记/")
	luteEngine.ParseOptions.GitConflict = true
	luteEngine.ParseOptions.LinkRef = false

	ast.Testing = true
	for _, test := range md2BlockDOMTests {
		result := luteEngine.Md2BlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
	ast.Testing = false
}
