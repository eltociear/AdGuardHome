package home

import (
	"net"
	"net/netip"

	"github.com/AdguardTeam/AdGuardHome/internal/aghalg"
	"golang.org/x/exp/slices"
)

// macUID contains MAC and UID.
type macUID struct {
	mac net.HardwareAddr
	uid UID
}

// clientIndex stores all information about persistent clients.
type clientIndex struct {
	clientIDToUID map[string]UID

	ipToUID map[netip.Addr]UID

	subnetToUID aghalg.OrderedMap[netip.Prefix, UID]

	macUIDs []*macUID

	uidToClient map[UID]*persistentClient
}

// NewClientIndex initializes the new instance of client index.
func NewClientIndex() (ci *clientIndex) {
	return &clientIndex{
		clientIDToUID: map[string]UID{},
		ipToUID:       map[netip.Addr]UID{},
		subnetToUID:   aghalg.NewOrderedMap[netip.Prefix, UID](subnetCompare),
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
		ci.macUIDs = append(ci.macUIDs, &macUID{mac, c.UID})
	}

	ci.uidToClient[c.UID] = c
}

// contains returns true if the index already has information about persistent
// client.
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
		ci.subnetToUID.Range(func(p netip.Prefix, id UID) bool {
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
		ok = slices.ContainsFunc(ci.macUIDs, func(muid *macUID) bool {
			return slices.Compare(mac, muid.mac) == 0
		})

		if ok {
			return true
		}
	}

	return false
}

// find finds persistent client by string represenation of the client ID, IP
// address, or MAC.
func (ci *clientIndex) find(id string) (c *persistentClient, ok bool) {
	uid, found := ci.clientIDToUID[id]
	if found {
		return ci.uidToClient[uid], true
	}

	ip, err := netip.ParseAddr(id)
	if err == nil {
		return ci.findByIP(ip)
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

	ci.subnetToUID.Range(func(pref netip.Prefix, id UID) bool {
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
	var uid UID
	found = slices.ContainsFunc(ci.macUIDs, func(muid *macUID) bool {
		if slices.Compare(mac, muid.mac) == 0 {
			uid = muid.uid

			return true
		}

		return false
	})

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
		ci.macUIDs = append(ci.macUIDs, &macUID{mac, c.UID})
		slices.DeleteFunc(ci.macUIDs, func(muid *macUID) bool {
			return slices.Compare(mac, muid.mac) == 0
		})
	}

	delete(ci.uidToClient, c.UID)
}
