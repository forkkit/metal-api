package datastore

import (
	"strconv"

	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/utils"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
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
func (rs *RethinkStore) FindNetwork(q *NetworkSearchQuery, n *metal.Network) error {
	return rs.findEntity(q.generateTerm(rs), &n)
}

// SearchNetworks returns the networks that match the given properties
func (rs *RethinkStore) SearchNetworks(q *NetworkSearchQuery, ns *metal.Networks) error {
	return rs.searchEntities(q.generateTerm(rs), ns)
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
