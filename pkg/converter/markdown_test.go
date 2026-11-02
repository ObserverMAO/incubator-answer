/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package converter

import (
	"regexp"
	"strings"
	"testing"
)

// matches dangerous constructs as real markup, not as percent-encoded text inside a URL
var dangerous = regexp.MustCompile(`(?i)<script|\son[a-z]+\s*=|javascript:|evil\.example`)

func TestMarkdown2HTML_BlocksEnclaveInjection(t *testing.T) {
	cases := map[string]string{
		// payloads seen in production on 2026-11-01
		"youtube id breakout": "![](https://www.youtube.com/watch?v=x%22%3E%3Cscript%3Ealert(document.domain)%3C/script%3E)",
		"image size breakout": "![](https://example.com/a.png?w=100px%22%3E%3Cimg%20src=x%20onerror=alert(1)%3E&h=100px \"cap\")",
		"raw script":          "<script>alert(1)</script>",
		"raw iframe":          `<iframe src="https://evil.example/"></iframe>`,
		"js link":             "[x](javascript:alert(1))",
	}
	for name, md := range cases {
		out := Markdown2HTML(md)
		if m := dangerous.FindString(out); m != "" {
			t.Errorf("%s: output contains %q:\n%s", name, m, out)
		}
	}
}

func TestMarkdown2HTML_KeepsLegitEmbeds(t *testing.T) {
	cases := map[string]string{
		"![](https://www.youtube.com/watch?v=KN6Cj2-eX_Q)":            `src="https://www.youtube.com/embed/KN6Cj2-eX_Q"`,
		"![](https://www.bilibili.com/video/BV113xieaEHQ)":            `src="//player.bilibili.com/player.html?bvid=BV113xieaEHQ&amp;page=1&amp;autoplay=0&amp;poster=1"`,
		"| a |\n|:--|\n| b |":                                         `style="text-align: left"`,
		"![](https://discuss.mixin.one/uploads/post/5pL4152MgYu.mov)": `src="https://discuss.mixin.one/uploads/post/5pL4152MgYu.mov"`,
	}
	for md, want := range cases {
		out := Markdown2HTML(md)
		if !strings.Contains(out, want) {
			t.Errorf("input %q: want %s in output:\n%s", md, want, out)
		}
	}
}
