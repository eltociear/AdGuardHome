package home

import (
	"net"
	"net/netip"
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
		err := ci.clashes(client1)
		require.NoError(t, err)

		dup := &persistentClient{
			Name: "client_with_the_same_ip_as_client1",
			IPs:  []netip.Addr{netip.MustParseAddr(cliIP1)},
			UID:  MustNewUID(),
		}
		err = ci.clashes(dup)
		require.Error(t, err)

		ci.del(client1)
		err = ci.clashes(dup)
		require.NoError(t, err)
	})
}

func TestMACToKey(t *testing.T) {
	macs := []string{
		"00:00:5e:00:53:01",
		"02:00:5e:10:00:00:00:01",
		"00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01",
		"00-00-5e-00-53-01",
		"02-00-5e-10-00-00-00-01",
		"00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01",
		"0000.5e00.5301",
		"0200.5e10.0000.0001",
		"0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001",
	}

	for _, m := range macs {
		mac, err := net.ParseMAC(m)
		require.NoError(t, err)

		key := macToKey(mac)
		assert.Len(t, key, len(mac))
	}

	assert.Panics(t, func() {
		mac := net.HardwareAddr([]byte{1, 2, 3})
		_ = macToKey(mac)
	})
}
