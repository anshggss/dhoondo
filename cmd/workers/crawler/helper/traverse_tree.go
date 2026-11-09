package workerhelper

import (
	"fmt"
	"strings"

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
	} else {
		fmt.Println("condition didn't match")
	}
	for temp := root.FirstChild; temp != nil; temp = temp.NextSibling {
		links = append(links, ExtractLinks(temp, depth+1)...)
	}
	return links
}

func ExtractWords(root *html.Node, depth int) []string {
	if root == nil {
		return nil
	}
	words := []string{}
	if root.Type == html.TextNode {
		words = append(words, strings.Fields(root.Data)...)
	}
	for temp := root.FirstChild; temp != nil; temp = temp.NextSibling {
		words = append(words, ExtractWords(temp, depth+1)...)
	}
	return words
}
