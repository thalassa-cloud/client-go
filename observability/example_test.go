package observability_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/thalassa-cloud/client-go/observability"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func ExampleClient_CreateObservabilityWorkspace() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusCreated, observability.ObservabilityWorkspace{
			Identity:          "obsw-abc",
			Name:              "prod",
			Status:            observability.ObservabilityWorkspaceStatusReady,
			PrometheusEnabled: true,
			LokiEnabled:       true,
			RemoteWriteURL:    "https://prometheus.example/workspace/obsw-abc/api/v1/push",
			PushURL:           "https://loki.example/workspace/obsw-abc/loki/api/v1/push",
		})
	}))
	defer server.Close()

	c, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	obs, _ := observability.New(c)

	ws, _ := obs.CreateObservabilityWorkspace(context.Background(), observability.CreateObservabilityWorkspaceRequest{
		Name:           "prod",
		RegionIdentity: "nl-01",
	})
	fmt.Println(ws.Identity, ws.PrometheusEnabled, ws.LokiEnabled)
	// Output: obsw-abc true true
}
