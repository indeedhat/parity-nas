package webproxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/servermux"
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

			var (
				err             error
				newData         io.Reader
				contentLength   int
				contentType     = r.Header.Get("Content-Type")
				contentEncoding = r.Header.Get("Content-Encoding")
			)

			bodyData, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			r.Body.Close()

			var bodyReader io.Reader = bytes.NewReader(bodyData)
			log.Print("content-type: ", contentType)

			switch {
			case strings.Contains(contentType, "html") || contentType == "":
				newData, contentLength, err = processHtmlRespons(bodyReader, basePath)
			case strings.Contains(contentType, "javascript"):
				bodyReader, err = decompress(bodyReader, contentEncoding)
				newData, contentLength, err = processJavascriptResponse(bodyReader, basePath)
				if err == nil {
					newData, err = compress(newData, contentEncoding)
				}
			default:
				log.Print("unhandled content type: ", contentType)
			}

			if newData != nil {
				r.Body = io.NopCloser(newData)
				r.ContentLength = int64(contentLength)
				r.Header.Set("Content-Length", strconv.Itoa(contentLength))
			} else {
				log.Print("error: ", err)
				r.Body = io.NopCloser(bytes.NewReader(bodyData))
			}

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
