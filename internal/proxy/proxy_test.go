package proxy

import (
	"strings"
	"testing"
)

func TestRewriteLocation(t *testing.T) {
	tests := []struct {
		loc    string
		prefix string
		want   string
	}{
		{"/dashboard", "/proxy/5000/", "/proxy/5000/dashboard"},
		{"/funds", "/proxy/5000/", "/proxy/5000/funds"},
		{"/", "/proxy/5000/", "/proxy/5000/"},
		{"http://localhost:5000/dashboard", "/proxy/5000/", "http://localhost:5000/dashboard"},
		{"https://other.com/path", "/proxy/5000/", "https://other.com/path"},
		{"", "/proxy/5000/", ""},
		{"/api/data", "/proxy/8747/", "/proxy/8747/api/data"},
	}
	for _, tt := range tests {
		got := rewriteLocation(tt.loc, tt.prefix)
		if got != tt.want {
			t.Errorf("rewriteLocation(%q, %q) = %q, want %q", tt.loc, tt.prefix, got, tt.want)
		}
	}
}

func TestRewriteHTML_BasicRewrites(t *testing.T) {
	prefix := "/proxy/5000/"
	html := `<html><head><title>Fund</title></head><body><img src="/assets/logo.png"><a href="/dashboard">Dash</a></body></html>`

	result := string(rewriteHTML([]byte(html), []byte(prefix)))

	if !strings.Contains(result, `src="/proxy/5000/assets/logo.png"`) {
		t.Errorf("asset src not rewritten: %s", result)
	}
	if !strings.Contains(result, `href="/proxy/5000/dashboard"`) {
		t.Errorf("href not rewritten: %s", result)
	}
	if !strings.Contains(result, `<base href="/proxy/5000/">`) {
		t.Errorf("base tag not injected: %s", result)
	}
	// Verify <title> preserved
	if !strings.Contains(result, "<title>Fund</title>") {
		t.Errorf("title tag was destroyed: %s", result)
	}
}

func TestRewriteHTML_BaseInjection(t *testing.T) {
	prefix := "/proxy/5000/"

	t.Run("no existing base", func(t *testing.T) {
		html := `<html><head></head><body></body></html>`
		result := string(rewriteHTML([]byte(html), []byte(prefix)))
		if !strings.Contains(result, `<base href="/proxy/5000/">`) {
			t.Errorf("base tag not injected: %s", result)
		}
	})

	t.Run("existing base updated", func(t *testing.T) {
		html := `<html><head><base href="/"></head><body></body></html>`
		result := string(rewriteHTML([]byte(html), []byte(prefix)))
		if strings.Contains(result, `<base href="/">`) {
			t.Errorf("original base not replaced: %s", result)
		}
		if !strings.Contains(result, `<base href="/proxy/5000/">`) {
			t.Errorf("base not updated to proxy prefix: %s", result)
		}
	})

	t.Run("single quoted base", func(t *testing.T) {
		html := `<html><head><base href='/app/'><title>T</title></head><body></body></html>`
		result := string(rewriteHTML([]byte(html), []byte(prefix)))
		if !strings.Contains(result, `href="/proxy/5000/app/"`) {
			t.Errorf("single-quoted base not updated: %s", result)
		}
	})
}

func TestRewriteHTML_DoublePrefixProtection(t *testing.T) {
	prefix := "/proxy/5000/"
	// Simulate a response that already contains the proxy prefix
	html := `<html><head></head><body><img src="/proxy/5000/assets/logo.png"></body></html>`

	result := string(rewriteHTML([]byte(html), []byte(prefix)))

	if strings.Contains(result, `/proxy/5000/proxy/5000/`) {
		t.Errorf("double prefix detected: %s", result)
	}
	if !strings.Contains(result, `src="/proxy/5000/assets/logo.png"`) {
		t.Errorf("already-prefixed path was mangled: %s", result)
	}
}

func TestRewriteHTML_MultipleRewrites(t *testing.T) {
	prefix := "/proxy/5000/"
	html := `<html><head></head><body>
		<link rel="stylesheet" href="/assets/main.css">
		<script src="/assets/app.js"></script>
		<img src="/favicon.ico">
		<a href="/settings">Settings</a>
	</body></html>`

	result := string(rewriteHTML([]byte(html), []byte(prefix)))

	checks := []string{
		`href="/proxy/5000/assets/main.css"`,
		`src="/proxy/5000/assets/app.js"`,
		`src="/proxy/5000/favicon.ico"`,
		`href="/proxy/5000/settings"`,
	}
	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("missing expected rewrite %q in: %s", check, result)
		}
	}
}

func TestRewriteHTML_NonPathAttributes(t *testing.T) {
	prefix := "/proxy/5000/"
	// charset, content, id etc. should NOT be rewritten
	html := `<html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width"><div id="/test"></div></head><body></body></html>`

	result := string(rewriteHTML([]byte(html), []byte(prefix)))

	if strings.Contains(result, `charset="/proxy/5000/`) {
		t.Errorf("charset attribute was incorrectly rewritten: %s", result)
	}
	if strings.Contains(result, `content="/proxy/5000/`) {
		t.Errorf("content attribute was incorrectly rewritten: %s", result)
	}
}

func TestRewriteCSS(t *testing.T) {
	prefix := "/proxy/5000/"
	css := `.icon { background: url("/assets/icon.png"); } .icon2 { background: url('/assets/icon2.png'); }`

	result := string(rewriteCSS([]byte(css), []byte(prefix)))

	if !strings.Contains(result, `url("/proxy/5000/assets/icon.png")`) {
		t.Errorf("quoted url not rewritten: %s", result)
	}
	if !strings.Contains(result, `url('/proxy/5000/assets/icon2.png')`) {
		t.Errorf("single-quoted url not rewritten: %s", result)
	}
}

func TestStripProxyPrefix(t *testing.T) {
	tests := []struct {
		path     string
		wantPort string
		wantSub  string
	}{
		{"/proxy/5000/dashboard", "5000", "/dashboard"},
		{"/proxy/5000/", "5000", "/"},
		{"/proxy/5000", "5000", "/"},
		{"/proxy/8747/api/data", "8747", "/api/data"},
		{"/health", "", "/health"},
		{"/", "", "/"},
	}
	for _, tt := range tests {
		port, sub := stripProxyPrefix(tt.path)
		if port != tt.wantPort {
			t.Errorf("stripProxyPrefix(%q) port = %q, want %q", tt.path, port, tt.wantPort)
		}
		if sub != tt.wantSub {
			t.Errorf("stripProxyPrefix(%q) subPath = %q, want %q", tt.path, sub, tt.wantSub)
		}
	}
}
