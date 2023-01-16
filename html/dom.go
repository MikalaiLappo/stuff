package dom

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
)


func ExtractAttr(node *html.Node, attr string) *string {
	for _, a := range node.Attr {
		if a.Key != attr {
			continue
		}
		s := a.Val
		if s == "" {
			s = "true"
		}
		return &s
	}
	return nil
}

func ExtractClassList(node *html.Node) *[]string {
	classAttr := ExtractAttr(node, "class")
	if classAttr == nil {
		return nil
	}
	classList := strings.Fields(*classAttr)
	return &classList
}


func InnerText(entryNode *html.Node) string {
	text := ""

	CbAllByCond(entryNode, func(node *html.Node) bool {
		return node.Type == html.TextNode
	}, func(node *html.Node) {
		text += node.Data
	})
	
	return text
}

func CbInnerText(entryNode *html.Node, Cb func(string)) {
	Cb(InnerText(entryNode))
}

func GetChildren(entryNode *html.Node) []*html.Node {
	var result []*html.Node
	result = nil

	for child := entryNode.FirstChild; child != nil; child = child.NextSibling {
		result = append(result, child)
	}

	return result
}

func CbChildren(entryNode *html.Node, Cb func(*html.Node)) {
	var children []*html.Node
	children = GetChildren(entryNode)

	if children == nil {
		return
	}
	
	for _, child := range children {
		Cb(child)
	}
}

func CbDrill(entryNode *html.Node, Cb func(*html.Node)) {
	var rec func(*html.Node)

	rec = func(node *html.Node) {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			Cb(child)
			rec(child)
		}
	}

	rec(entryNode)
}

func GetByCond(entryNode *html.Node, cond func(*html.Node) bool) *html.Node {
	var rec func(*html.Node) *html.Node

	rec = func(node *html.Node) *html.Node {
		if cond(node) {
			return node
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			recChild := rec(child)
			if recChild != nil {
				return recChild
			}
		}
		return nil
	}
	return rec(entryNode)
}

func CbByCond(entryNode *html.Node, cond func(*html.Node) bool, Cb func(*html.Node)) {
	Cb(GetByCond(entryNode, cond))
}

func GetAllByCond(entryNode *html.Node, cond func(*html.Node) bool) []*html.Node {
	result := []*html.Node{}
	var rec func(*html.Node)

	rec = func(node *html.Node) {
		if cond(node) {
			result = append(result, node)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			rec(child)
		}
	}
	rec(entryNode)
	return result
}

func CbAllByCond(entryNode *html.Node, cond func(*html.Node) bool, Cb func(*html.Node)) {
	for _, node := range GetAllByCond(entryNode, cond) {
		Cb(node)
	}
}

func tagCond(tagName string) func(*html.Node) bool {
	return func(node *html.Node) bool {
		return node.Type == html.ElementNode && node.Data == tagName
	}
}

func GetByTag(tagName string, entryNode *html.Node) *html.Node {
	return GetByCond(entryNode, tagCond(tagName))
}

func GetAllByTag(tagName string, entryNode *html.Node) []*html.Node {
	return GetAllByCond(entryNode, tagCond(tagName))
}

func CbAllByTag(tagName string, entryNode *html.Node, Cb func(*html.Node)) {
	CbAllByCond(entryNode, tagCond(tagName), Cb)
}

func classNameCond(className string) func(*html.Node) bool {
	return func(node *html.Node) bool {
		if node.Type != html.ElementNode {
			return false
		}
		cla := ExtractClassList(node)
		if cla == nil {
			return false
		}
		for _, c := range *cla {
			if c == className {
				return true
			}
		}
		return false
	}
}

func GetByClass(tagName string, entryNode *html.Node) *html.Node {
	return GetByCond(entryNode, classNameCond(tagName))
}

func GetAllByClass(tagName string, entryNode *html.Node) []*html.Node {
	return GetAllByCond(entryNode, classNameCond(tagName))
}

func CbAllByClass(tagName string, entryNode *html.Node, Cb func(*html.Node)) {
	CbAllByCond(entryNode, classNameCond(tagName), Cb)
}

func Ntos(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

