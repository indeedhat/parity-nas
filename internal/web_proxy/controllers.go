package webproxy

import (
	"bytes"
	"errors"
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
	"github.com/indeedhat/parity-nas/pkg/server_mux"
)

// WebProxyController sets up a dynamic web proxy to forward traffic from the /./* path to other services
// running in parity-nas
func WebProxyController(rw http.ResponseWriter, r *http.Request) {
	cfg, err := config.WebProxy()
	if err != nil {
		servermux.InternalError(rw, "Failed to load proxy config")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")[1:]
	if len(pathParts) < 2 {
		servermux.WriteError(rw, http.StatusUnprocessableEntity, "Invalid proxy url")
		return
	}

	basePath := "/" + path.Join(pathParts[0], pathParts[1])

	handlerCfg := cfg.FIndHandler(pathParts[1])
	if handlerCfg == nil {
		servermux.WriteError(
			rw,
			http.StatusNotFound,
			fmt.Sprintf("Cannot find proxy config for prefix /%s/%s", cfg.Prefix, pathParts[1]),
		)
		return
	}

	outUrl, err := r.URL.Parse(fmt.Sprintf("%s://%s:%d/%s",
		handlerCfg.Scheme,
		handlerCfg.Host,
		handlerCfg.Port,
		path.Join(pathParts[2:]...),
	))
	if err != nil {
		servermux.InternalError(rw, "Failed to build up destination url")
		return
	}

	proxy := httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			log.Printf("%s -> %s", r.URL.String(), outUrl.String())

			pr.SetURL(outUrl)
			pr.Out.URL.Path = outUrl.Path
			pr.Out.URL.RawPath = outUrl.RawPath
		},

		ModifyResponse: func(pr *http.Response) error {
			location := pr.Header.Get("location")
			if location != "" {
				pr.Header.Set("location", path.Join(basePath, location))
			}

			var (
				err             error
				newData         io.Reader
				contentLength   int
				contentType     = pr.Header.Get("Content-Type")
				contentEncoding = pr.Header.Get("Content-Encoding")
			)

			bodyData, err := io.ReadAll(pr.Body)
			if err != nil {
				return err
			}
			pr.Body.Close()

			bodyReader, err := decompress(bytes.NewReader(bodyData), contentEncoding)

			switch {
			case strings.Contains(contentType, "html") || contentType == "":
				newData, contentLength, err = processHtmlRespons(handlerCfg.Mutators, bodyReader, basePath)
			case strings.Contains(contentType, "javascript"):
				newData, contentLength, err = processJavascriptResponse(handlerCfg.Mutators, bodyReader, basePath)
			default:
				log.Print("unhandled content type: ", contentType)
			}

			if errors.Is(err, errProcessSkipped) {
				contentLength = int(pr.ContentLength)
			} else if err != nil {
				newData, err = compress(newData, contentEncoding)
			}

			if newData != nil {
				pr.Body = io.NopCloser(newData)
				pr.ContentLength = int64(contentLength)
				pr.Header.Set("Content-Length", strconv.Itoa(contentLength))
			} else {
				log.Print("error: ", err)
				pr.Body = io.NopCloser(bytes.NewReader(bodyData))
			}

			return nil
		},
	}

	proxy.Transport = http.DefaultTransport
	proxy.Transport.(*http.Transport).MaxIdleConns = 3000
	proxy.Transport.(*http.Transport).MaxIdleConnsPerHost = 3000
	proxy.Transport.(*http.Transport).IdleConnTimeout = 10 * time.Second
	proxy.Transport.(*http.Transport).MaxConnsPerHost = 0

	proxy.ServeHTTP(rw, r)
}
