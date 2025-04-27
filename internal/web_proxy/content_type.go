package webproxy

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

// processHtmlRespons takes the provided input stream and rewrites links and js import statements to add the
// provided base path to any path going from root
func processHtmlRespons(reader io.Reader, basePath string) (io.Reader, int, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, 0, err
	}

	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, 0, err
	}

	modifyLinkTags(doc, basePath)
	modifyJsTags(doc, basePath)

	head := findHeadNode(doc)
	if head != nil {
		modifyBaseTag(head, basePath)
	}

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, 0, err
	}

	return bytes.NewReader(buf.Bytes()), buf.Len(), nil
}

// processJavascriptResponse takes the provided input stream and rewrites all import statements going
// from root with the provided base path
func processJavascriptResponse(reader io.Reader, basePath string) (io.Reader, int, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, 0, err
	}

	processed := modifyJsImports(string(data), basePath)

	return bytes.NewReader([]byte(processed)), len(processed), nil
}
