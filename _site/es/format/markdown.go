package format

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func ToHtml(input []byte) string {
	unsafe := blackfriday.MarkdownCommon(input)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return string(html)
}
