package proxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// cssURLRe matches CSS url() references with root-relative paths.
// Groups: 1=quote char (or ""), 2=path after "/"
var cssURLRe = regexp.MustCompile(`(?m)url\(\s*(['"]?)/([^"')]+)`)

// rewriteCSS rewrites CSS url(/...) → url(/proxy/PORT/...)
func rewriteCSS(body, proxyPrefix []byte) []byte {
	return cssURLRe.ReplaceAllFunc(body, func(match []byte) []byte {
		subs := cssURLRe.FindSubmatch(match)
		if len(subs) < 3 {
			return match
		}
		path := strings.TrimPrefix(string(subs[2]), "/")
		return []byte("url(" + string(subs[1]) + string(proxyPrefix) + path)
	})
}

// stripProxyPrefix strips the leading "/proxy/PORT/" from a request path.
// Returns the port number and the remaining sub-path.
// If the path doesn't match the proxy pattern, returns empty port.
func stripProxyPrefix(requestPath string) (port, subPath string) {
	path := requestPath
	if !strings.HasPrefix(path, "/proxy/") {
		return "", path
	}
	rest := strings.TrimPrefix(path, "/proxy/")
	slash := strings.Index(rest, "/")
	if slash < 0 {
		return rest, "/"
	}
	port = rest[:slash]
	subPath = rest[slash:]
	return port, subPath
}

// rewriteLocation rewrites a relative Location header to include the proxy prefix.
// Nginx-style proxy_redirect: absolute URLs pass through, relative paths get prefixed.
func rewriteLocation(loc, prefix string) string {
	if loc == "" || strings.HasPrefix(loc, "http") {
		return loc
	}
	if strings.HasPrefix(loc, "/") {
		return strings.TrimRight(prefix, "/") + loc
	}
	return prefix + loc
}

// Handler returns the main proxy handler. Nginx-style reverse proxy with dual pattern support:
//
// Pattern 1 — /proxy/PORT/xxx → forward to localhost:PORT/xxx
//
//	This is the primary proxying mechanism (Nginx location /proxy/5000/).
//	Browser navigation and static resources flow through this path.
//
// Pattern 2 — /api/xxx → forward to the default port (typically 5000)
//
//	Browser XHR/fetch requests resolve absolute /api/xxx against the origin,
//	not the <base> tag. This catch-all handles those SPA API calls by forwarding
//	them to the registered default backend (Nginx location /api/).
//
// Go ServeMux routes longest-prefix first, so portal's own /api/health and
// /api/services will always be matched before this catch-all.
func Handler(defaultPort string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS preflight
		if r.Method == http.MethodOptions {
			origin := r.Header.Get("Origin")
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD, PATCH")
				w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
				w.Header().Set("Access-Control-Max-Age", "3600")
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		port, subPath := resolveTarget(r.URL.Path, r.Header.Get("Referer"), defaultPort)
		if port == "" {
			http.NotFound(w, r)
			return
		}

		proxyPrefix := "/proxy/" + port + "/"

		target, _ := url.Parse("http://localhost:" + port)

		director := func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = subPath
			req.URL.RawQuery = r.URL.RawQuery
			req.Host = target.Host
			req.Header.Set("X-Forwarded-For", r.RemoteAddr)
			if fwd := r.Header.Get("X-Real-IP"); fwd != "" {
				req.Header.Set("X-Real-IP", fwd)
			}
			req.Header.Set("X-Forwarded-Host", r.Host)
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			req.Header.Set("X-Forwarded-Proto", scheme)
			req.Header.Del("If-None-Match")
			req.Header.Del("If-Modified-Since")
		}

		errHandler := func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("proxy error (port %s): %v", port, err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, "<!DOCTYPE html><html><head><style>body{font-family:system-ui;background:#0d1117;color:#e6edf3;display:flex;align-items:center;justify-content:center;height:100vh}div{text-align:center}h1{font-size:20px}p{color:#8b949e}</style></head><body><div><h1>Service Unreachable</h1><p>Port %s is not responding</p></div></body></html>", port)
		}

		// Nginx-style pass-through for Next.js (port 20128):
		// Use raw http.Client + io.Copy to avoid any Go HTTP layer manipulation.
		// httputil.ReverseProxy, even with DisableCompression + no ModifyResponse,
		// may alter chunk boundaries → React hydration #418.
		if port == "20128" {
			upstreamURL := "http://localhost:" + port + subPath
			if r.URL.RawQuery != "" {
				upstreamURL += "?" + r.URL.RawQuery
			}
			upstreamReq, err := http.NewRequestWithContext(r.Context(), r.Method, upstreamURL, r.Body)
			if err != nil {
				http.Error(w, "Bad gateway", http.StatusBadGateway)
				return
			}
			copyHeaders(upstreamReq.Header, r.Header)
			upstreamReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
			upstreamReq.Header.Set("X-Forwarded-Host", r.Host)
			upstreamReq.Header.Del("If-None-Match")
			upstreamReq.Header.Del("If-Modified-Since")

			client := &http.Client{
				Transport: &http.Transport{
					DisableCompression: true,
				},
			}
			resp, err := client.Do(upstreamReq)
			if err != nil {
				log.Printf("proxy error (port %s): %v", port, err)
				http.Error(w, "Bad gateway", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			w.Header().Del("Content-Security-Policy")
			copyHeaders(w.Header(), resp.Header)
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
			return
		}

		proxy := &httputil.ReverseProxy{
			Director: director,
			ModifyResponse: func(resp *http.Response) error {
				origin := r.Header.Get("Origin")
				if origin != "" {
					resp.Header.Set("Access-Control-Allow-Origin", origin)
					resp.Header.Set("Access-Control-Allow-Credentials", "true")
				}

				if loc := resp.Header.Get("Location"); loc != "" {
					resp.Header.Set("Location", rewriteLocation(loc, proxyPrefix))
				}

				ct := resp.Header.Get("Content-Type")
				if strings.Contains(ct, "text/css") {
					orig := resp.Body
					body, err := io.ReadAll(orig)
					if err != nil {
						return nil
					}
					orig.Close()
					body = rewriteCSS(body, []byte(proxyPrefix))
					resp.Body = io.NopCloser(bytes.NewReader(body))
					resp.ContentLength = int64(len(body))
					resp.Header.Set("Content-Length", strconv.Itoa(len(body)))
				}
				if strings.Contains(ct, "text/html") {
					orig := resp.Body
					body, err := io.ReadAll(orig)
					if err != nil {
						return nil
					}
					orig.Close()
					body = rewriteHTML(body, []byte(proxyPrefix))
					resp.Body = io.NopCloser(bytes.NewReader(body))
					resp.ContentLength = int64(len(body))
					resp.Header.Set("Content-Length", strconv.Itoa(len(body)))
					resp.Header.Del("Etag")
					resp.Header.Del("X-Nextjs-Cache")
					resp.Header.Del("X-Nextjs-Stale-Time")
					resp.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
				}
				return nil
			},
			ErrorHandler:   errHandler,
			FlushInterval:  100 * 1000 * 1000,
		}
		w.Header().Del("Content-Security-Policy")
		proxy.ServeHTTP(w, r)
	}
}

// resolveTarget determines the target port and sub-path from a request URL.
// Nginx-style routing with Referer fallback (equivalent to Nginx map $http_referer).
//   - "/proxy/PORT/xxx" → port=PORT, subPath=/xxx
//   - other with Referer "/proxy/PORT/..." → port=PORT, subPath=as-is
//   - other → port=defaultPort, subPath=as-is
func resolveTarget(requestPath, referer, defaultPort string) (port, subPath string) {
	if port, sp := stripProxyPrefix(requestPath); port != "" {
		return port, sp
	}
	if port := extractPortFromReferer(referer); port != "" {
		return port, requestPath
	}
	return defaultPort, requestPath
}

// extractPortFromReferer extracts the port from /proxy/PORT/ pattern in the referer URL.
// Also detects /dashboard/ → 20128 (9Router SPA).
func extractPortFromReferer(referer string) string {
	idx := strings.Index(referer, "/proxy/")
	if idx >= 0 {
		rest := referer[idx+len("/proxy/"):]
		slash := strings.Index(rest, "/")
		if slash > 0 {
			return rest[:slash]
		}
	}
	u, err := url.Parse(referer)
	if err == nil {
		p := u.Path
		if p == "/dashboard" || strings.HasPrefix(p, "/dashboard/") {
			return "20128"
		}
	}
	return ""
}

// rewriteHTML rewrites absolute paths in HTML to include the proxy prefix.
// Only rewrites src/href attribute values — does NOT inject new elements (<base>/<script>)
// to avoid React hydration mismatch (#418).
func rewriteHTML(body, proxyPrefix []byte) []byte {
	prefix := string(proxyPrefix)

	// Step 1: Rewrite src="/..." and href="/..." paths. Skip non-absolute paths
	// and paths already carrying the proxy prefix (single-pass, no dedup needed).
	re := regexp.MustCompile(`(?i)\b(src|href)\s*=\s*(?:"([^"]*)"|'([^']*)'|(\S+))`)
	body = re.ReplaceAllFunc(body, func(match []byte) []byte {
		subs := re.FindSubmatch(match)
		if len(subs) < 5 {
			return match
		}
		var rawPath string
		switch {
		case len(subs[2]) > 0:
			rawPath = string(subs[2])
		case len(subs[3]) > 0:
			rawPath = string(subs[3])
		default:
			rawPath = strings.TrimSuffix(string(subs[4]), ">")
		}
		if !strings.HasPrefix(rawPath, "/") {
			return match
		}
		if strings.HasPrefix(rawPath, prefix) {
			return match
		}
		attr := string(subs[1])
		newPath := prefix + strings.TrimPrefix(rawPath, "/")
		if len(subs[2]) > 0 {
			return []byte(attr + `="` + newPath + `"`)
		}
		if len(subs[3]) > 0 {
			return []byte(attr + "='" + newPath + "'")
		}
		return []byte(attr + "=" + newPath + ">")
	})
	return body
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
