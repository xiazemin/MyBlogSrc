package format

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"

	gohtml "golang.org/x/net/html"
)

var (
	raw          *regexp.Regexp = regexp.MustCompile("\\{%\\s*raw\\s*%\\}")
	endraw       *regexp.Regexp = regexp.MustCompile("\\{%\\s*endraw\\s*%\\}")
	highlight    *regexp.Regexp = regexp.MustCompile("\\{%\\s*highlight\\w*\\s*%\\}")
	endhighlight *regexp.Regexp = regexp.MustCompile("\\{%\\s*endhighlight\\w*\\s*%\\}")
)

func ParseHtml(html, name string) (map[string]string, error) {
	header := make(map[string]string)
	//fmt.Println(raw.ReplaceAllString("12324{% raw %}", "<code>"))
	html = raw.ReplaceAllString(html, "<code>")
	html = endraw.ReplaceAllString(html, "</code>")
	html = highlight.ReplaceAllString(html, "<code>")
	html = endhighlight.ReplaceAllString(html, "</code>")

	utfBody, err := iconv.NewReader(strings.NewReader(html), "utf-8", "utf-8")
	if err != nil {
		return header, err
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		if _, err := gohtml.Parse(utfBody); err != nil {
			//ioutil.WriteFile("invalid.html", []byte(html), fs.ModeAppend)
			//panic(err)
			fmt.Println(err)
			name1 := name[:len(name)-9]
			nameParts := strings.Split(name1, "-")
			header["title"] = nameParts[len(nameParts)-1]
			header["category"] = nameParts[len(nameParts)-1]
			return header, nil
		}
		return header, err
	}

	// Find the review items
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			return
		}
		// For each item found, get the title
		//title := s.Find("br").Text()
		//fmt.Printf("Review %d: %s,%s\n", i, s.Text(), title)
		parts := strings.Split(s.Text(), "\n")
		// fmt.Println(parts)
		for _, part := range parts {
			if part == "" {
				continue
			}
			kv := strings.Split(strings.Trim(part, " "), ":")
			if len(kv) < 2 {
				fmt.Println(kv)
			}
			header[kv[0]] = kv[1]
		}
	})
	doc.Find("code").Each(func(i int, s *goquery.Selection) {
		//fmt.Println("code-------")
		//fmt.Println(s.Text())
	})
	return header, nil
}
