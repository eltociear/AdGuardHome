package home

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientIndex(t *testing.T) {
	var (
		cliNone = "1.2.3.4"
		cli1    = "1.1.1.1"
		cli2    = "2.2.2.2"

		cli1IP = netip.MustParseAddr(cli1)
		cli2IP = netip.MustParseAddr(cli2)

		cliIPv6 = netip.MustParseAddr("1:2:3::4")
	)

	ci := NewClientIndex()

	uid, err := NewUID()
	require.NoError(t, err)

	client1 := &persistentClient{
		Name: "client1",
		IPs:  []netip.Addr{cli1IP, cliIPv6},
		UID:  uid,
	}

	uid, err = NewUID()
	require.NoError(t, err)

	client2 := &persistentClient{
		Name: "client2",
		IPs:  []netip.Addr{cli2IP},
		UID:  uid,
	}

	t.Run("add_find", func(t *testing.T) {
		ci.add(client1)
		ci.add(client2)

		c, ok := ci.find(cli1)
		require.True(t, ok)

		assert.Equal(t, "client1", c.Name)

		c, ok = ci.find("1:2:3::4")
		require.True(t, ok)

		assert.Equal(t, "client1", c.Name)

		c, ok = ci.find(cli2)
		require.True(t, ok)

		assert.Equal(t, "client2", c.Name)

		_, ok = ci.find(cliNone)
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
