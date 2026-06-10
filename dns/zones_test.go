package dns

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

func TestRecordFQDN(t *testing.T) {
	tests := []struct {
		zone   string
		record string
		want   string
	}{
		{"example.com", "@", "example.com"},
		{"example.com", "*", "*.example.com"},
		{"example.com", "www", "www.example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.record, func(t *testing.T) {
			assert.Equal(t, tt.want, RecordFQDN(tt.zone, tt.record))
		})
	}
}

func TestFormatRecordValues(t *testing.T) {
	assert.Equal(t, "10 mail.example.com", FormatMX(10, "mail.example.com"))
	assert.Equal(t, "0 issue letsencrypt.org", FormatCAA(0, "issue", "letsencrypt.org"))
	assert.Equal(t, "10 5 5060 sip.example.com", FormatSRV(10, 5, 5060, "sip.example.com"))
}

func TestListZones(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/dns/zones", r.URL.Path)
		assert.Equal(t, "prod", r.URL.Query().Get("matchLabels[env]"))

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode([]DnsZone{
			{Identity: "dnsz-1", Name: "example.com"},
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	dnsClient, err := New(c)
	require.NoError(t, err)

	zones, err := dnsClient.ListZones(context.Background(), &ListZonesRequest{
		Filters: []ListZonesFilter{
			ListZonesFilterFromFilter(&filters.LabelFilter{
				MatchLabels: map[string]string{"env": "prod"},
			}),
		},
	})
	require.NoError(t, err)
	require.Len(t, zones, 1)
	assert.Equal(t, "example.com", zones[0].Name)
}

func TestCreateZone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/dns/zones", r.URL.Path)

		var body CreateDnsZoneRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "example.com", body.ZoneName)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		require.NoError(t, json.NewEncoder(w).Encode(DnsZone{
			Identity: "dnsz-abc123",
			Name:     body.ZoneName,
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	dnsClient, err := New(c)
	require.NoError(t, err)

	zone, err := dnsClient.CreateZone(context.Background(), CreateDnsZoneRequest{
		ZoneName: "example.com",
	})
	require.NoError(t, err)
	assert.Equal(t, "dnsz-abc123", zone.Identity)
	assert.Equal(t, "example.com", zone.Name)
}

func TestCreateRecordConflict(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		require.NoError(t, json.NewEncoder(w).Encode(map[string]string{
			"message": "record already exists",
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	dnsClient, err := New(c)
	require.NoError(t, err)

	_, err = dnsClient.CreateRecord(context.Background(), "dnsz-abc123", CreateDnsRecordRequest{
		Name:   "www",
		Type:   DnsRecordTypeA,
		Values: []string{"192.0.2.1"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "409")
}

func TestUpdateRecord(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/v1/dns/zones/dnsz-abc123/records/dnsr-xyz789", r.URL.Path)

		var body UpdateDnsRecordRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, 600, body.TTL)
		assert.Equal(t, []string{"192.0.2.10"}, body.Values)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(DnsRecord{
			Identity: "dnsr-xyz789",
			Name:     "www",
			Type:     DnsRecordTypeA,
			TTL:      body.TTL,
			Values:   body.Values,
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	dnsClient, err := New(c)
	require.NoError(t, err)

	record, err := dnsClient.UpdateRecord(context.Background(), "dnsz-abc123", "dnsr-xyz789", UpdateDnsRecordRequest{
		TTL:    600,
		Values: []string{"192.0.2.10"},
	})
	require.NoError(t, err)
	assert.Equal(t, 600, record.TTL)
}

func TestExportImportZoneFile(t *testing.T) {
	bindText := "$ORIGIN example.com.\nwww 300 IN A 192.0.2.1\n"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/dns/zones/dnsz-abc123/export":
			require.NoError(t, json.NewEncoder(w).Encode(ExportDnsZoneFileResponse{
				ZoneName: "example.com",
				ZoneFile: bindText,
			}))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/dns/zones/dnsz-abc123/import":
			var body ImportDnsZoneFileRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, bindText, body.ZoneFile)

			require.NoError(t, json.NewEncoder(w).Encode(ImportDnsZoneFileResponse{
				Created: 1,
				Records: []DnsRecord{{Name: "www", Type: DnsRecordTypeA, Values: []string{"192.0.2.1"}}},
			}))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	dnsClient, err := New(c)
	require.NoError(t, err)
	ctx := context.Background()

	exported, err := dnsClient.ExportZoneFile(ctx, "dnsz-abc123")
	require.NoError(t, err)
	assert.Equal(t, bindText, exported.ZoneFile)

	imported, err := dnsClient.ImportZoneFile(ctx, "dnsz-abc123", ImportDnsZoneFileRequest{
		ZoneFile: exported.ZoneFile,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, imported.Created)
}

func TestDeleteZone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/dns/zones/dnsz-abc123", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	dnsClient, err := New(c)
	require.NoError(t, err)

	err = dnsClient.DeleteZone(context.Background(), "dnsz-abc123")
	require.NoError(t, err)
}

func TestGetDnssec(t *testing.T) {
	now := time.Now()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/dns/zones/dnsz-abc123/dnssec", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(DnsZoneDnssecStatus{
			Enabled:      true,
			DsDelegated:  true,
			LastSignedAt: &now,
			Region:       "nl-01",
		}))
	}))
	defer server.Close()

	c, err := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	require.NoError(t, err)
	dnsClient, err := New(c)
	require.NoError(t, err)

	status, err := dnsClient.GetDnssec(context.Background(), "dnsz-abc123")
	require.NoError(t, err)
	assert.True(t, status.Enabled)
	assert.True(t, status.DsDelegated)
}

func TestZonePath(t *testing.T) {
	assert.Equal(t, "/v1/dns/zones/dnsz-abc123/records", zonePath("dnsz-abc123", "records"))
}
