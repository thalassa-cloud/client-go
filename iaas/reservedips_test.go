package iaas

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssociateReservedIpRequest_JSON(t *testing.T) {
	lb := "lb-abc"
	nat := "nat-xyz"

	tests := []struct {
		name string
		req  AssociateReservedIpRequest
		want map[string]any
	}{
		{
			name: "load balancer only",
			req:  AssociateReservedIpRequest{LoadbalancerIdentity: &lb},
			want: map[string]any{"loadbalancerIdentity": "lb-abc"},
		},
		{
			name: "nat gateway only",
			req:  AssociateReservedIpRequest{NatGatewayIdentity: &nat},
			want: map[string]any{"natGatewayIdentity": "nat-xyz"},
		},
		{
			name: "empty struct omits both",
			req:  AssociateReservedIpRequest{},
			want: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.req)
			require.NoError(t, err)
			var got map[string]any
			require.NoError(t, json.Unmarshal(b, &got))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateLoadbalancer_reservedIpIdJSON(t *testing.T) {
	fip := "fip-123"
	b, err := json.Marshal(CreateLoadbalancer{
		Name:                 "x",
		Subnet:               "sub",
		InternalLoadbalancer: false,
		DeleteProtection:     false,
		Listeners:            nil,
		ReservedIpID:         &fip,
	})
	require.NoError(t, err)
	var m map[string]any
	require.NoError(t, json.Unmarshal(b, &m))
	assert.Equal(t, "fip-123", m["reservedIpId"])
}

func TestUpdateLoadbalancer_reservedIpIdPointerSemantics(t *testing.T) {
	t.Run("empty string detaches", func(t *testing.T) {
		empty := ""
		b, err := json.Marshal(UpdateLoadbalancer{
			Name:             "n",
			Description:      "",
			DeleteProtection: false,
			ReservedIpID:     &empty,
		})
		require.NoError(t, err)
		assert.Contains(t, string(b), `"reservedIpId":""`)
	})
	t.Run("nil omits field", func(t *testing.T) {
		b, err := json.Marshal(UpdateLoadbalancer{
			Name:             "n",
			Description:      "",
			DeleteProtection: false,
			ReservedIpID:     nil,
		})
		require.NoError(t, err)
		assert.NotContains(t, string(b), "reservedIpId")
	})
}

func TestReservedIP_unmarshal(t *testing.T) {
	raw := `{
		"identity": "fip-test",
		"name": "pub",
		"slug": "pub",
		"description": "",
		"createdAt": "2025-01-01T12:00:00Z",
		"status": "available",
		"ipv4Address": "203.0.113.1"
	}`
	var f ReservedIP
	require.NoError(t, json.Unmarshal([]byte(raw), &f))
	assert.Equal(t, "fip-test", f.Identity)
	assert.Equal(t, ReservedIpStatusAvailable, f.Status)
	assert.Equal(t, "203.0.113.1", f.IPv4Address)
}
