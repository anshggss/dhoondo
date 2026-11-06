package crawler

import (
	"golang.org/x/net/html"
)

func ExtractLinks(root *html.Node, depth int) []string {
	if root == nil {
		return nil
	}
	links := []string{}
	if root.Type == html.ElementNode && root.Data == "a" {
		for _, attr := range root.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}
	for temp := root.FirstChild; temp != nil; temp = temp.NextSibling {
		links = append(links, ExtractLinks(temp, depth+1)...)
	}
	return links
}
