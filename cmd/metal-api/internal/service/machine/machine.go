package machine

import (
	"time"

	mdm "github.com/metal-stack/masterdata-api/pkg/client"

	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"
	"github.com/metal-stack/metal-lib/bus"
)

const (
	waitForServerTimeout = 30 * time.Second
)

type machineResource struct {
	bus.Publisher
	ds     *datastore.RethinkStore
	ipamer ipam.IPAMer
	mdc    mdm.Client
}

// NewMachine returns a webservice for machine specific endpoints.
func NewMachine(
	ds *datastore.RethinkStore,
	pub bus.Publisher,
	ipamer ipam.IPAMer,
	mdc mdm.Client) *restful.WebService {
	r := machineResource{
		ds:        ds,
		Publisher: pub,
		ipamer:    ipamer,
		mdc:       mdc,
	}
	return r.webService()
}
