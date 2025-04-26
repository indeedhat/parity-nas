package webproxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"path"
	"strings"
	"time"

	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/servermux"
	"golang.org/x/net/html"
)

// WebProxyController sets up a dynamic web proxy to forward traffic from the /./* path to other services
// running in parity-nas
func WebProxyController(ctx servermux.Context) error {
	cfg, err := config.WebProxy()
	if err != nil {
		return ctx.InternalError("Failed to load proxy config")
	}

	url := ctx.Request().URL
	pathParts := strings.Split(url.Path, "/")[1:]
	if len(pathParts) < 2 {
		return ctx.Error(http.StatusUnprocessableEntity, "Invalid proxy url")
	}

	basePath := "/" + path.Join(pathParts[0], pathParts[1])

	handlerCfg := cfg.FIndHandler(pathParts[1])
	if handlerCfg == nil {
		return ctx.Error(
			http.StatusNotFound,
			fmt.Sprintf("Cannot find proxy config for prefix /%s/%s", cfg.Prefix, pathParts[1]),
		)
	}

	outUrl, err := url.Parse(fmt.Sprintf("%s://%s:%d/%s",
		handlerCfg.Scheme,
		handlerCfg.Host,
		handlerCfg.Port,
		path.Join(pathParts[2:]...),
	))
	if err != nil {
		return ctx.InternalError("Failed to build up destination url")
	}

	proxy := httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			log.Printf("%s -> %s", url.String(), outUrl.String())

			r.SetURL(outUrl)
			r.Out.URL.Path = outUrl.Path
			r.Out.URL.RawPath = outUrl.RawPath
		},

		ModifyResponse: func(r *http.Response) error {
			location := r.Header.Get("location")
			if location != "" {
				r.Header.Set("location", path.Join(basePath, location))
			}

			if !strings.Contains(r.Header.Get("Content-Type"), "html") {
				return nil
			}

			doc, err := html.Parse(r.Body)
			if err != nil {
				return err
			}

			modifyLinks(doc, basePath)

			head := findHeadNode(doc)
			if head != nil {
				modifyBaseTag(head, basePath)
			}

			var buf bytes.Buffer
			if err := html.Render(&buf, doc); err != nil {
				return err
			}

			log.Print(string(buf.Bytes()))
			r.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))

			return nil
		},
	}

	proxy.Transport = http.DefaultTransport
	proxy.Transport.(*http.Transport).MaxIdleConns = 3000
	proxy.Transport.(*http.Transport).MaxIdleConnsPerHost = 3000
	proxy.Transport.(*http.Transport).IdleConnTimeout = 10 * time.Second
	proxy.Transport.(*http.Transport).MaxConnsPerHost = 0

	proxy.ServeHTTP(ctx.Writer(), ctx.Request())

	return nil
}

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

func modifyLinks(node *html.Node, basePath string) {
	if node.Type == html.ElementNode {
		for i, attr := range node.Attr {
			if (attr.Key == "href" || attr.Key == "src") &&
				attr.Val[0] == '/' &&
				!strings.HasPrefix(attr.Val, basePath) {

				node.Attr[i].Val = path.Join(basePath, attr.Val)
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		modifyLinks(child, basePath)
	}
}
