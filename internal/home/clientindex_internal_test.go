package home

import (
	"testing"

	"github.com/AdguardTeam/AdGuardHome/internal/filtering"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientIndex(t *testing.T) {
	const (
		cliIPNone = "1.2.3.4"
		cliIP1    = "1.1.1.1"
		cliIP2    = "2.2.2.2"

		cliIPv6 = "1:2:3::4"

		cliSubnet   = "2.2.2.0/24"
		cliSubnetIP = "2.2.2.222"

		cliID  = "client-id"
		cliMAC = "11:11:11:11:11:11"
	)

	objs := []clientObject{{
		Name:            "client1",
		IDs:             []string{cliIP1, cliIPv6},
		BlockedServices: &filtering.BlockedServices{},
	}, {
		Name:            "client2",
		IDs:             []string{cliIP2, cliSubnet},
		BlockedServices: &filtering.BlockedServices{},
	}, {
		Name:            "client_with_mac",
		IDs:             []string{cliMAC},
		BlockedServices: &filtering.BlockedServices{},
	}, {
		Name:            "client_with_id",
		IDs:             []string{cliID},
		BlockedServices: &filtering.BlockedServices{},
	}}

	clients := []*persistentClient{}
	for _, o := range objs {
		cli, err := o.toPersistent(&filtering.Config{}, nil)
		require.NoError(t, err)

		clients = append(clients, cli)
	}

	client1 := clients[0]
	client2 := clients[1]
	clientWithMAC := clients[2]
	clientWithID := clients[3]

	ci := NewClientIndex()

	t.Run("add_find", func(t *testing.T) {
		ci.add(client1)
		ci.add(client2)
		ci.add(clientWithMAC)
		ci.add(clientWithID)

		c, ok := ci.find(cliIP1)
		require.True(t, ok)

		assert.Equal(t, client1.Name, c.Name)

		c, ok = ci.find(cliIPv6)
		require.True(t, ok)

		assert.Equal(t, client1.Name, c.Name)

		c, ok = ci.find(cliIP2)
		require.True(t, ok)

		assert.Equal(t, client2.Name, c.Name)

		c, ok = ci.find(cliSubnetIP)
		require.True(t, ok)

		assert.Equal(t, client2.Name, c.Name)

		c, ok = ci.find(cliMAC)
		require.True(t, ok)

		assert.Equal(t, clientWithMAC.Name, c.Name)

		c, ok = ci.find(cliID)
		require.True(t, ok)

		assert.Equal(t, clientWithID.Name, c.Name)

		_, ok = ci.find(cliIPNone)
		assert.False(t, ok)
	})

	t.Run("contains_delete", func(t *testing.T) {
		ok := ci.contains(client1)
		require.True(t, ok)

		ci.del(client1)
		ok = ci.contains(client1)
		require.False(t, ok)
	})
}
