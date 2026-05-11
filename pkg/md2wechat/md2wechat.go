package md2wechat

import (
	"bytes"
	"net/url"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
)

var (
	unsupportedTagRe = regexp.MustCompile(`</?(iframe|script|style|object|embed|form|input|select|textarea|button|canvas|svg|video|audio|source)\b[^>]*/?>`)
	classAttrRe      = regexp.MustCompile(`\s+class="[^"]*"`)
	styleAttrRe      = regexp.MustCompile(`\s+style="[^"]*"`)
	srcAttrRe        = regexp.MustCompile(`src="([^"]*)"`)
)

func Convert(markdown string, cdnURL string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}

	html := buf.String()

	html = unsupportedTagRe.ReplaceAllString(html, "")

	html = cleanHTML(html)

	if cdnURL != "" {
		html = resolveRelativeURLs(html, strings.TrimRight(cdnURL, "/"))
	}

	html = strings.TrimSpace(html)
	return html, nil
}

func cleanHTML(html string) string {
	html = classAttrRe.ReplaceAllString(html, "")
	html = styleAttrRe.ReplaceAllString(html, "")
	return html
}

func resolveRelativeURLs(html, cdnURL string) string {
	return srcAttrRe.ReplaceAllStringFunc(html, func(match string) string {
		urlMatch := srcAttrRe.FindStringSubmatch(match)
		if len(urlMatch) < 2 {
			return match
		}
		imgSrc := urlMatch[1]

		if strings.HasPrefix(imgSrc, "http://") || strings.HasPrefix(imgSrc, "https://") || strings.HasPrefix(imgSrc, "data:") {
			return match
		}

		parsed, err := url.Parse(imgSrc)
		if err != nil {
			return match
		}

		absURL := cdnURL + "/" + strings.TrimPrefix(parsed.Path, "/")
		return `src="` + absURL + `"`
	})
}
