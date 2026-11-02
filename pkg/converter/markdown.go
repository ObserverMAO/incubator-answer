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
	"bytes"
	"regexp"

	"github.com/asaskevich/govalidator"
	enclave "github.com/lixvyang/goldmark-enclave"
	"github.com/microcosm-cc/bluemonday"
	"github.com/segmentfault/pacman/log"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	goldmarkHTML "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// Markdown2HTML convert markdown to html
func Markdown2HTML(source string) string {
	mdConverter := goldmark.New(
		goldmark.WithExtensions(&DangerousHTMLFilterExtension{}, extension.GFM, enclave.New(&enclave.Config{})),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkHTML.WithHardWraps(),
		),
	)
	var buf bytes.Buffer
	if err := mdConverter.Convert([]byte(source), &buf); err != nil {
		log.Error(err)
		return source
	}
	// The enclave extension writes raw HTML with unescaped user input
	// (YouTube IDs, image params, error messages), so its output must never be
	// trusted: always run the final document through the sanitizer.
	return outputPolicy.Sanitize(buf.String())
}

var (
	// Embeds we render on purpose. Anything else in an iframe src is dropped.
	embedIframeSrc = regexp.MustCompile(
		`^(?:https:)?//(?:www\.youtube\.com/embed/[A-Za-z0-9_-]{1,64}|player\.bilibili\.com/player\.html\?[A-Za-z0-9=&_-]{1,256})$`)
	iframeAllow  = regexp.MustCompile(`^[a-z\-; ]*$`)
	iframeLength = regexp.MustCompile(`^[0-9]{1,4}%?$`)
	iframeFlag   = regexp.MustCompile(`^(?:true|false|0|1|no|yes|)$`)
	videoSrc     = regexp.MustCompile(`^https://[^\s"'<>]+$`)

	outputPolicy = newOutputPolicy()
)

func newOutputPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowStyling()
	p.RequireNoFollowOnLinks(false)
	p.RequireParseableURLs(false)
	p.RequireNoFollowOnFullyQualifiedLinks(false)
	p.AllowElements("kbd")
	p.AllowAttrs("title").Matching(regexp.MustCompile(`^[\p{L}\p{N}\s\-_',\[\]!\./\\\(\)]*$|^@embed?$`)).Globally()
	p.AllowAttrs("start").OnElements("ol")

	// GFM table alignment
	p.AllowStyles("text-align").OnElements("th", "td")
	// enclave image wrapper sizing (values are validated by bluemonday's CSS handlers)
	p.AllowStyles("width", "height", "margin", "display").OnElements("span", "img")

	// YouTube / Bilibili embeds
	p.AllowAttrs("src").Matching(embedIframeSrc).OnElements("iframe")
	p.AllowAttrs("width", "height").Matching(iframeLength).OnElements("iframe")
	p.AllowAttrs("allow").Matching(iframeAllow).OnElements("iframe")
	p.AllowAttrs("allowfullscreen", "frameborder", "border", "framespacing", "scrolling").
		Matching(iframeFlag).OnElements("iframe")

	// uploaded videos
	p.AllowAttrs("src").Matching(videoSrc).OnElements("video")
	p.AllowAttrs("width", "height").Matching(iframeLength).OnElements("video")
	p.AllowAttrs("controls").Matching(iframeFlag).OnElements("video")
	return p
}

// Markdown2BasicHTML convert markdown to html ,Only basic syntax can be used
func Markdown2BasicHTML(source string) string {
	content := Markdown2HTML(source)
	filter := bluemonday.NewPolicy()
	filter.AllowElements("p", "b", "br", "strong", "em")
	filter.AllowAttrs("src").OnElements("img")
	filter.AddSpaceWhenStrippingTag(true)
	content = filter.Sanitize(content)
	return content
}

type DangerousHTMLFilterExtension struct {
}

func (e *DangerousHTMLFilterExtension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&DangerousHTMLRenderer{
			Config: goldmarkHTML.NewConfig(),
			Filter: bluemonday.UGCPolicy(),
		}, 1),
	))
}

type DangerousHTMLRenderer struct {
	goldmarkHTML.Config
	Filter *bluemonday.Policy
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *DangerousHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
}

func (r *DangerousHTMLRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}
	n := node.(*ast.RawHTML)
	l := n.Segments.Len()
	for i := 0; i < l; i++ {
		segment := n.Segments.At(i)
		if string(source[segment.Start:segment.Stop]) == "<kbd>" || string(source[segment.Start:segment.Stop]) == "</kbd>" {
			_, _ = w.Write(segment.Value(source))
		} else {
			_, _ = w.Write(r.Filter.SanitizeBytes(segment.Value(source)))
		}
	}
	return ast.WalkSkipChildren, nil
}

func (r *DangerousHTMLRenderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.HTMLBlock)
	if entering {
		l := n.Lines().Len()
		for i := 0; i < l; i++ {
			line := n.Lines().At(i)
			r.Writer.SecureWrite(w, r.Filter.SanitizeBytes(line.Value(source)))
		}
	} else {
		if n.HasClosure() {
			closure := n.ClosureLine
			r.Writer.SecureWrite(w, closure.Value(source))
		}
	}
	return ast.WalkContinue, nil
}

func (r *DangerousHTMLRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering && r.renderLinkIsUrl(string(n.Destination)) {
		_, _ = w.WriteString("<a href=\"")
		// _, _ = w.WriteString("<a test=\"1\" rel=\"nofollow\" href=\"")
		if r.Unsafe || !html.IsDangerousURL(n.Destination) {
			_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		}
		_ = w.WriteByte('"')
		if n.Title != nil {
			_, _ = w.WriteString(` title="`)
			r.Writer.Write(w, n.Title)
			_ = w.WriteByte('"')
		}
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, html.LinkAttributeFilter)
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

func (r *DangerousHTMLRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.AutoLink)

	if !entering || !r.renderLinkIsUrl(string(n.URL(source))) {
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString(`<a href="`)
	url := n.URL(source)
	label := n.Label(source)
	if n.AutoLinkType == ast.AutoLinkEmail && !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
		_, _ = w.WriteString("mailto:")
	}
	_, _ = w.Write(util.EscapeHTML(util.URLEscape(url, false)))
	if n.Attributes() != nil {
		_ = w.WriteByte('"')
		html.RenderAttributes(w, n, html.LinkAttributeFilter)
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString(`">`)
	}
	_, _ = w.Write(util.EscapeHTML(label))
	_, _ = w.WriteString(`</a>`)
	return ast.WalkContinue, nil
}

func (r *DangerousHTMLRenderer) renderLinkIsUrl(verifyUrl string) bool {
	isURL := govalidator.IsURL(verifyUrl)
	isPath, _ := regexp.MatchString(`^/`, verifyUrl)
	return isURL || isPath
}
