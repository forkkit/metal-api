package datastore

import (
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	v1 "github.com/metal-stack/metal-api/cmd/metal-api/internal/service/proto/v1"
)

// FindNetworkByID returns an network of a given id.
func (rs *RethinkStore) FindNetworkByID(id string) (*metal.Network, error) {
	var nw metal.Network
	err := rs.findEntityByID(rs.NetworkTable(), &nw, id)
	if err != nil {
		return nil, err
	}
	return &nw, nil
}

// FindNetwork returns a machine by the given query, fails if there is no record or multiple records found.
func (rs *RethinkStore) FindNetwork(q *v1.NetworkSearchQuery, n *metal.Network) error {
	return rs.findEntity(q.GenerateTerm(rs), &n)
}

// SearchNetworks returns the networks that match the given properties
func (rs *RethinkStore) SearchNetworks(q *v1.NetworkSearchQuery, ns *metal.Networks) error {
	return rs.searchEntities(q.GenerateTerm(rs), ns)
}

// ListNetworks returns all networks.
func (rs *RethinkStore) ListNetworks() (metal.Networks, error) {
	nws := make(metal.Networks, 0)
	err := rs.listEntities(rs.NetworkTable(), &nws)
	return nws, err
}

// CreateNetwork creates a new network.
func (rs *RethinkStore) CreateNetwork(nw *metal.Network) error {
	return rs.createEntity(rs.NetworkTable(), nw)
}

// DeleteNetwork deletes an network.
func (rs *RethinkStore) DeleteNetwork(nw *metal.Network) error {
	return rs.deleteEntity(rs.NetworkTable(), nw)
}

// UpdateNetwork updates an network.
func (rs *RethinkStore) UpdateNetwork(oldNetwork *metal.Network, newNetwork *metal.Network) error {
	return rs.updateEntity(rs.NetworkTable(), newNetwork, oldNetwork)
}
