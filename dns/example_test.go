package dns_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/thalassa-cloud/client-go/dns"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func ExampleClient_CreateZone() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusCreated, dns.DnsZone{
			Identity: "dnsz-abc123",
			Name:     "example.com",
		})
	}))
	defer server.Close()

	c, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	dnsClient, _ := dns.New(c)

	zone, _ := dnsClient.CreateZone(context.Background(), dns.CreateDnsZoneRequest{
		ZoneName: "example.com",
	})
	fmt.Println(zone.Name)
	// Output: example.com
}

func ExampleClient_CreateRecord() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusCreated, dns.DnsRecord{
			Identity: "dnsr-xyz789",
			Name:     "www",
			Type:     dns.DnsRecordTypeA,
			TTL:      300,
			Values:   []string{"192.0.2.1"},
		})
	}))
	defer server.Close()

	c, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	dnsClient, _ := dns.New(c)

	record, _ := dnsClient.CreateRecord(context.Background(), "dnsz-abc123", dns.CreateDnsRecordRequest{
		Name:   "www",
		Type:   dns.DnsRecordTypeA,
		TTL:    300,
		Values: []string{"192.0.2.1"},
	})
	fmt.Printf("%s %s\n", record.Name, record.Type)
	// Output: www A
}

func ExampleFormatMX() {
	fmt.Println(dns.FormatMX(10, "mail.example.com"))
	// Output: 10 mail.example.com
}

func ExampleRecordFQDN() {
	fmt.Println(dns.RecordFQDN("example.com", "@"))
	fmt.Println(dns.RecordFQDN("example.com", "www"))
	// Output:
	// example.com
	// www.example.com
}

func Example() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/dns/zones":
			writeJSON(w, http.StatusCreated, dns.DnsZone{
				Identity: "dnsz-abc123",
				Name:     "example.com",
			})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/dns/zones/dnsz-abc123/records":
			writeJSON(w, http.StatusCreated, dns.DnsRecord{
				Name:   "www",
				Type:   dns.DnsRecordTypeA,
				Values: []string{"192.0.2.1"},
			})
		}
	}))
	defer server.Close()

	tc, _ := thalassa.NewClient(
		client.WithBaseURL(server.URL),
		client.WithAuthCustom(),
	)
	ctx := context.Background()
	dnsClient := tc.DNS()

	zone, _ := dnsClient.CreateZone(ctx, dns.CreateDnsZoneRequest{ZoneName: "example.com"})
	record, _ := dnsClient.CreateRecord(ctx, zone.Identity, dns.CreateDnsRecordRequest{
		Name:   "www",
		Type:   dns.DnsRecordTypeA,
		Values: []string{"192.0.2.1"},
	})
	fmt.Printf("%s -> %s\n", dns.RecordFQDN(zone.Name, record.Name), record.Values[0])
	// Output: www.example.com -> 192.0.2.1
}

func ExampleClient_ListRecords() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("recordType") != "TXT" {
			http.Error(w, "unexpected filter", http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, []dns.DnsRecord{
			{Name: "_acme-challenge", Type: dns.DnsRecordTypeTXT, Values: []string{"token"}},
		})
	}))
	defer server.Close()

	c, _ := client.NewClient(client.WithBaseURL(server.URL), client.WithAuthCustom())
	dnsClient, _ := dns.New(c)

	records, _ := dnsClient.ListRecords(context.Background(), "dnsz-abc123", &dns.ListRecordsRequest{
		Filters: []dns.ListRecordsFilter{
			dns.ListRecordsFilterFromFilter(&filters.FilterKeyValue{
				Key:   filters.FilterKey("recordType"),
				Value: "TXT",
			}),
		},
	})
	fmt.Println(records[0].Name)
	// Output: _acme-challenge
}
