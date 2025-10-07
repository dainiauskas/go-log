package log

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEnrichHTTPMeta_HeadersAndRequestInfo(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/foo?x=1", nil)
	req.Header.Set("X-Request-Id", "req-123")
	req.Header.Set("UserName", "tester")
	req.Header.Set("Application-Version", "v1.2.3")
	req.Header.Set("User-Agent", "unittest-agent")

	meta := EnrichHTTPMeta(400, req, nil, 1)

	if meta["method"] != "GET" {
		t.Fatalf("expected method GET, got %v", meta["method"])
	}
	if meta["path"] != "/foo" {
		t.Fatalf("expected path /foo, got %v", meta["path"])
	}
	if meta["query"] != "x=1" {
		t.Fatalf("expected query x=1, got %v", meta["query"])
	}
	if meta["X-Request-Id"] != "req-123" {
		t.Fatalf("expected X-Request-Id req-123, got %v", meta["X-Request-Id"])
	}
	if meta["UserName"] != "tester" {
		t.Fatalf("expected UserName tester, got %v", meta["UserName"])
	}
	if meta["Application-Version"] != "v1.2.3" {
		t.Fatalf("expected Application-Version v1.2.3, got %v", meta["Application-Version"])
	}
	if meta["User-Agent"] != "unittest-agent" {
		t.Fatalf("expected User-Agent unittest-agent, got %v", meta["User-Agent"])
	}
}

func TestEnrichHTTPMeta_TimestampAndHostnamePresent(t *testing.T) {
	meta := EnrichHTTPMeta(400, nil, nil, 1)
	if meta["timestamp"] == nil {
		t.Fatalf("timestamp not set")
	}
	if meta["hostname"] == nil {
		t.Fatalf("hostname not set")
	}
}

func TestEnrichHTTPMeta_OriginPreservedAndPopulated(t *testing.T) {
	// Provide a prefilled origin and ensure it's preserved
	initial := map[string]interface{}{"origin": "custom"}
	meta := EnrichHTTPMeta(400, nil, initial, 1)
	if meta["origin"] != "custom" {
		t.Fatalf("expected origin preserved as custom, got %v", meta["origin"])
	}

	// Without prefill, origin should be populated
	meta2 := EnrichHTTPMeta(400, nil, nil, 1)
	if meta2["origin"] == nil {
		t.Fatalf("expected origin to be populated, was nil")
	}
}

func TestEnrichHTTPMeta_StackIncludedFor5xx(t *testing.T) {
	meta := EnrichHTTPMeta(500, nil, nil, 1)
	st, ok := meta["stack"].(string)
	if !ok {
		t.Fatalf("expected stack string present for 500")
	}
	if !strings.Contains(st, "TestEnrichHTTPMeta_StackIncludedFor5xx") && !strings.Contains(st, "enrich_http_meta_test") {
		// We don't require exact frame names across platforms, but expect the test name or file to appear.
		t.Fatalf("stack does not appear to contain test frame: %v", st)
	}
}
