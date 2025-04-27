package webproxy

import (
	"path"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// findHeadNode traverses the html AST to extract the <head> node
func findHeadNode(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "head" {
		return n
	}

	// Recursively check child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if head := findHeadNode(c); head != nil {
			return head
		}
	}

	return nil
}

// modifyBaseTag searches for the <base> tag and prefixes any root path found there with the provided base path
func modifyBaseTag(node *html.Node, basePath string) bool {
	if node.Type == html.ElementNode && node.Data == "base" {
		for i, attr := range node.Attr {
			if attr.Key == "href" {
				if !strings.HasPrefix(attr.Val, basePath) {
					node.Attr[i].Val = path.Join(basePath, attr.Val)
				}
				return true
			}
		}

		node.Attr = append(node.Attr, html.Attribute{Key: "href", Val: basePath + "/"})
		return true
	}

	// Recursively check child nodes
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if modifyBaseTag(child, basePath) {
			return true
		}
	}

	if node.Type == html.ElementNode && node.Data == "head" {
		baseTag := html.Node{
			Type: html.ElementNode,
			Data: "base",
			Attr: []html.Attribute{
				{Key: "href", Val: basePath + "/"},
			},
		}
		if node.FirstChild != nil {
			node.InsertBefore(&baseTag, node.FirstChild)
		} else {
			node.AppendChild(&baseTag)
		}
	}

	return false
}

// modifyLinkTags locates all tags with an href/src attribute and prefixes any that are root links
// (start with /) with the provided basePath
func modifyLinkTags(node *html.Node, basePath string) {
	if node.Type == html.ElementNode {
		for i, attr := range node.Attr {
			hasRootScopedHrefOrSrc := (attr.Key == "href" || attr.Key == "src") &&
				attr.Val[0] == '/' &&
				!strings.HasPrefix(attr.Val, basePath)

			if hasRootScopedHrefOrSrc {
				node.Attr[i].Val = path.Join(basePath, attr.Val)
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		modifyLinkTags(child, basePath)
	}
}

// modifyJsTags locates all <script> tags that have a body containing the import keyword and attempts
// to prefix and imports from root with the provided basePath
func modifyJsTags(node *html.Node, basePath string) {
	isScriptNodeWithImportStatements := node.Type == html.ElementNode &&
		node.Data == "script" &&
		node.FirstChild != nil &&
		node.FirstChild.Type == html.TextNode &&
		strings.Contains(node.FirstChild.Data, "import")

	if isScriptNodeWithImportStatements {
		node.FirstChild.Data = modifyJsImports(node.FirstChild.Data, basePath)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		modifyJsTags(child, basePath)
	}
}

// importJsRe regex used to identify and extract js import statements
var importJsRe = regexp.MustCompile(`import(?:["'\s]*([\w*{}\n, ]+)from\s*)?["'\s]*([@\w/_\-.]+)["'\s].*`)

// modifyJsImports extracts all import statments and prefixes any imports from root with the provided
// basePath
func modifyJsImports(src, basePath string) string {
	matches := importJsRe.FindAllStringSubmatch(src, -1)

	if len(matches) == 0 {
		return src
	}

	for _, parts := range matches {
		if !strings.HasPrefix(parts[2], "/") {
			continue
		}

		prefixed := ""
		if parts[1] == "" {
			prefixed = importJsRe.ReplaceAllString(parts[0], "import '"+basePath+"${2}';")
		} else {
			prefixed = importJsRe.ReplaceAllString(parts[0], "import ${1} from '"+basePath+"${2}';")
		}

		src = strings.ReplaceAll(src, parts[0], prefixed)
	}

	return src
}
