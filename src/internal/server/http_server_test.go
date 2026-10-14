package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"country-iso-matcher/src/internal/version"
	"github.com/baditaflorin/go-common/graph"
)

func TestFleetGraphMiddlewareWritesAuthenticatedInboundEvent(t *testing.T) {
	keyPath := filepath.Join(t.TempDir(), "graph-writer.key")
	if err := os.WriteFile(keyPath, []byte("test-graph-writer-key"), 0o600); err != nil {
		t.Fatal(err)
	}

	batches := make(chan graph.Batch, 1)
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/events" {
			http.Error(w, "unexpected request", http.StatusNotFound)
			return
		}
		if got := r.Header.Get("X-API-Key"); got != "test-graph-writer-key" {
			http.Error(w, "writer key missing", http.StatusUnauthorized)
			return
		}
		var batch graph.Batch
		if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
			http.Error(w, "invalid graph batch", http.StatusBadRequest)
			return
		}
		batches <- batch
		w.WriteHeader(http.StatusAccepted)
	}))
	defer collector.Close()

	t.Setenv("GRAPH_ENABLED", "true")
	t.Setenv("GRAPH_SERVICE_ID", "country-iso-matcher")
	t.Setenv("GRAPH_COLLECTOR_URL", collector.URL)
	t.Setenv("GRAPH_API_KEY_FILE", keyPath)
	t.Setenv("GRAPH_FLUSH_INTERVAL", "1")
	t.Setenv("GRAPH_FLUSH_BATCH", "1")

	handler := fleetGraphMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/convert?country=Germany", nil))
	if response.Code != http.StatusNoContent {
		t.Fatalf("handler status = %d, want %d", response.Code, http.StatusNoContent)
	}

	select {
	case batch := <-batches:
		if batch.Service != "country-iso-matcher" || batch.Version != version.Version {
			t.Fatalf("batch identity = %q@%q, want country-iso-matcher@%q", batch.Service, batch.Version, version.Version)
		}
		if len(batch.Events) != 1 {
			t.Fatalf("batch event count = %d, want 1", len(batch.Events))
		}
		event := batch.Events[0]
		if event.Direction != "in" || event.Target != "country-iso-matcher" || event.Path != "/api/convert" || event.Status != http.StatusNoContent {
			t.Fatalf("inbound event = %+v", event)
		}
	case <-time.After(4 * time.Second):
		t.Fatal("timed out waiting for authenticated Fleet Graph event")
	}
}
