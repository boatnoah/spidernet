package crawler

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

func extractLinks(body []byte) ([]string, error) {
	doc, err := html.Parse(bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	links := make([]string, 0)

	seen := make(map[string]struct{}, 0)

	var walk func(n *html.Node)

	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					url := attr.Val

					_, ok := seen[url]

					if ok {
						continue
					}

					hasProto := strings.Index(url, "http") == 0
					if hasProto {
						links = append(links, url)
						seen[url] = struct{}{}
					}

				}

			}

		}

		for c := range n.ChildNodes() {
			if c != nil {
				walk(c)
			}
		}
	}

	walk(doc)

	return links, nil
}
