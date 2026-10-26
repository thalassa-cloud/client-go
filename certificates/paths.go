package certificates

import "strings"

const Endpoint = "/v1/certificates"

func regionPath(region string, segments ...string) string {
	parts := append([]string{Endpoint, region}, segments...)
	return strings.Join(parts, "/")
}
