package datastore

import (
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
)

// FindIPByID returns an ip of a given id.
func (rs *RethinkStore) FindIPByID(id string) (*metal.IP, error) {
	var ip metal.IP
	err := rs.findEntityByID(rs.ipTable(), &ip, id)
	if err != nil {
		return nil, err
	}
	return &ip, nil
}

// FindIPs returns an IP by the given query, fails if there is no record or multiple records found.
func (rs *RethinkStore) FindIPs(q *v1.IPFindRequest, ip *metal.IP) error {
	return rs.findEntity(q.GenerateTerm(*rs.ipTable()), &ip)
}

// SearchIPs returns the result of the ips search request query.
func (rs *RethinkStore) SearchIPs(q *v1.IPFindRequest, ips *metal.IPs) error {
	return rs.searchEntities(q.GenerateTerm(*rs.ipTable()), ips)
}

// ListIPs returns all ips.
func (rs *RethinkStore) ListIPs() (metal.IPs, error) {
	ips := make([]metal.IP, 0)
	err := rs.listEntities(rs.ipTable(), &ips)
	return ips, err
}

// CreateIP creates a new ip.
func (rs *RethinkStore) CreateIP(ip *metal.IP) error {
	return rs.createEntity(rs.ipTable(), ip)
}

// DeleteIP deletes an ip.
func (rs *RethinkStore) DeleteIP(ip *metal.IP) error {
	return rs.deleteEntity(rs.ipTable(), ip)
}

// UpdateIP updates an ip.
func (rs *RethinkStore) UpdateIP(oldIP *metal.IP, newIP *metal.IP) error {
	return rs.updateEntity(rs.ipTable(), newIP, oldIP)
}
