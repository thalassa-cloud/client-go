package iaas

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLoadbalancer_instanceSizeJSON(t *testing.T) {
	t.Run("set size is encoded", func(t *testing.T) {
		b, err := json.Marshal(CreateLoadbalancer{
			Name:         "x",
			Subnet:       "sub",
			InstanceSize: "lb-gp-small",
		})
		require.NoError(t, err)
		var m map[string]any
		require.NoError(t, json.Unmarshal(b, &m))
		assert.Equal(t, "lb-gp-small", m["instanceSize"])
	})
	t.Run("empty size is omitted", func(t *testing.T) {
		b, err := json.Marshal(CreateLoadbalancer{Name: "x", Subnet: "sub"})
		require.NoError(t, err)
		assert.NotContains(t, string(b), "instanceSize")
	})
}

func TestUpdateLoadbalancer_instanceSizePointerSemantics(t *testing.T) {
	t.Run("nil omits field", func(t *testing.T) {
		b, err := json.Marshal(UpdateLoadbalancer{Name: "n"})
		require.NoError(t, err)
		assert.NotContains(t, string(b), "instanceSize")
	})
	t.Run("non-empty is encoded", func(t *testing.T) {
		size := "lb-gp-large"
		b, err := json.Marshal(UpdateLoadbalancer{Name: "n", InstanceSize: &size})
		require.NoError(t, err)
		var m map[string]any
		require.NoError(t, json.Unmarshal(b, &m))
		assert.Equal(t, "lb-gp-large", m["instanceSize"])
	})
}

func TestCreateVpcNatGateway_instanceSizeJSON(t *testing.T) {
	b, err := json.Marshal(CreateVpcNatGateway{
		Name:           "egress",
		SubnetIdentity: "s-1",
		InstanceSize:   "ng-gp-small",
	})
	require.NoError(t, err)
	var m map[string]any
	require.NoError(t, json.Unmarshal(b, &m))
	assert.Equal(t, "ng-gp-small", m["instanceSize"])
}

func TestUpdateVpcNatGateway_instanceSizePointerSemantics(t *testing.T) {
	t.Run("nil omits field", func(t *testing.T) {
		b, err := json.Marshal(UpdateVpcNatGateway{Name: "n"})
		require.NoError(t, err)
		assert.NotContains(t, string(b), "instanceSize")
	})
	t.Run("non-empty is encoded", func(t *testing.T) {
		size := "ng-gp-medium"
		b, err := json.Marshal(UpdateVpcNatGateway{Name: "n", InstanceSize: &size})
		require.NoError(t, err)
		var m map[string]any
		require.NoError(t, json.Unmarshal(b, &m))
		assert.Equal(t, "ng-gp-medium", m["instanceSize"])
	})
}
