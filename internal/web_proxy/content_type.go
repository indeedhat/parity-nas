package webproxy

import (
	"bytes"
	"errors"
	"io"

	"github.com/indeedhat/parity-nas/internal/config"
	"golang.org/x/net/html"
)

var errProcessSkipped = errors.New("process skipped")

// processHtmlRespons takes the provided input stream and rewrites links and js import statements to add the
// provided base path to any path going from root
func processHtmlRespons(cfg config.WebProxyMutators, reader io.Reader, basePath string) (io.Reader, int, error) {
	if !cfg.HtmlLinks && !cfg.JsImports && !cfg.BaseTag {
		return reader, 0, errProcessSkipped
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, 0, err
	}

	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, 0, err
	}

	if cfg.HtmlLinks {
		modifyLinkTags(doc, basePath)
	}
	if cfg.JsImports {
		modifyJsTags(doc, basePath)
	}
	if cfg.BaseTag {
		head := findHeadNode(doc)
		if head != nil {
			modifyBaseTag(head, basePath)
		}
	}

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, 0, err
	}

	return bytes.NewReader(buf.Bytes()), buf.Len(), nil
}

// processJavascriptResponse takes the provided input stream and rewrites all import statements going
// from root with the provided base path
func processJavascriptResponse(cfg config.WebProxyMutators, reader io.Reader, basePath string) (io.Reader, int, error) {
	if !cfg.JsImports {
		return reader, 0, errProcessSkipped
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, 0, err
	}

	processed := modifyJsImports(string(data), basePath)

	return bytes.NewReader([]byte(processed)), len(processed), nil
}
