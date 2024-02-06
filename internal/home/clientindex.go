package home

import (
	"net"
	"net/netip"

	"github.com/AdguardTeam/AdGuardHome/internal/aghalg"
)

// macKey contains MAC as byte array of 6, 8, or 20 bytes.
type macKey any

func macToKey(mac net.HardwareAddr) (key macKey) {
	switch len(mac) {
	case 6:
		arr := [6]byte{}
		copy(arr[:], mac[:])

		return arr
	case 8:
		arr := [8]byte{}
		copy(arr[:], mac[:])

		return arr
	default:
		arr := [20]byte{}
		copy(arr[:], mac[:])

		return arr
	}
}

// clientIndex stores all information about persistent clients.
type clientIndex struct {
	clientIDToUID map[string]UID

	ipToUID map[netip.Addr]UID

	subnetToUID aghalg.SortedMap[netip.Prefix, UID]

	macToUID map[macKey]UID

	uidToClient map[UID]*persistentClient
}

// NewClientIndex initializes the new instance of client index.
func NewClientIndex() (ci *clientIndex) {
	return &clientIndex{
		clientIDToUID: map[string]UID{},
		ipToUID:       map[netip.Addr]UID{},
		subnetToUID:   aghalg.NewSortedMap[netip.Prefix, UID](subnetCompare),
		macToUID:      map[macKey]UID{},
		uidToClient:   map[UID]*persistentClient{},
	}
}

// add stores information about a persistent client in the index.
func (ci *clientIndex) add(c *persistentClient) {
	for _, id := range c.ClientIDs {
		ci.clientIDToUID[id] = c.UID
	}

	for _, ip := range c.IPs {
		ci.ipToUID[ip] = c.UID
	}

	for _, pref := range c.Subnets {
		ci.subnetToUID.Set(pref, c.UID)
	}

	for _, mac := range c.MACs {
		k := macToKey(mac)
		ci.macToUID[k] = c.UID
	}

	ci.uidToClient[c.UID] = c
}

// contains returns true if the index contains a persistent client with at least
// a single identifier contained by c.
func (ci *clientIndex) contains(c *persistentClient) (ok bool) {
	for _, id := range c.ClientIDs {
		_, ok = ci.clientIDToUID[id]
		if ok {
			return true
		}
	}

	for _, ip := range c.IPs {
		_, ok = ci.ipToUID[ip]
		if ok {
			return true
		}
	}

	for _, pref := range c.Subnets {
		ci.subnetToUID.Range(func(p netip.Prefix, _ UID) (cont bool) {
			if pref == p {
				ok = true

				return false
			}

			return true
		})

		if ok {
			return true
		}
	}

	for _, mac := range c.MACs {
		k := macToKey(mac)
		_, ok = ci.macToUID[k]
		if ok {
			return true
		}
	}

	return false
}

// find finds persistent client by string representation of the client ID, IP
// address, or MAC.
func (ci *clientIndex) find(id string) (c *persistentClient, ok bool) {
	uid, found := ci.clientIDToUID[id]
	if found {
		return ci.uidToClient[uid], true
	}

	ip, err := netip.ParseAddr(id)
	if err == nil {
		// MAC addresses can be successfully parsed as IP addresses.
		c, found = ci.findByIP(ip)
		if found {
			return c, true
		}
	}

	mac, err := net.ParseMAC(id)
	if err == nil {
		return ci.findByMAC(mac)
	}

	return nil, false
}

// find finds persistent client by IP address.
func (ci *clientIndex) findByIP(ip netip.Addr) (c *persistentClient, found bool) {
	uid, found := ci.ipToUID[ip]
	if found {
		return ci.uidToClient[uid], true
	}

	ci.subnetToUID.Range(func(pref netip.Prefix, id UID) (cont bool) {
		if pref.Contains(ip) {
			uid, found = id, true

			return false
		}

		return true
	})

	if found {
		return ci.uidToClient[uid], true
	}

	return nil, false
}

// find finds persistent client by MAC.
func (ci *clientIndex) findByMAC(mac net.HardwareAddr) (c *persistentClient, found bool) {
	k := macToKey(mac)
	uid, found := ci.macToUID[k]
	if found {
		return ci.uidToClient[uid], true
	}

	return nil, false
}

// del removes information about persistent client from the index.
func (ci *clientIndex) del(c *persistentClient) {
	for _, id := range c.ClientIDs {
		delete(ci.clientIDToUID, id)
	}

	for _, ip := range c.IPs {
		delete(ci.ipToUID, ip)
	}

	for _, pref := range c.Subnets {
		ci.subnetToUID.Del(pref)
	}

	for _, mac := range c.MACs {
		k := macToKey(mac)
		delete(ci.macToUID, k)
	}

	delete(ci.uidToClient, c.UID)
}
