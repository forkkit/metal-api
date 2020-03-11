package ip

import (
	"fmt"
	"github.com/emicklei/go-restful"
	mdm "github.com/metal-stack/masterdata-api/pkg/client"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
)

type ipResource struct {
	ds     *datastore.RethinkStore
	ipamer ipam.IPAMer
	mdc    mdm.Client
}

// NewIP returns a webservice for ip specific endpoints.
func NewIP(ds *datastore.RethinkStore, ipamer ipam.IPAMer, mdc mdm.Client) *restful.WebService {
	r := ipResource{
		ds:     ds,
		ipamer: ipamer,
		mdc:    mdc,
	}
	return r.webService()
}

func AllocateIP(parent *metal.Network, specificIP string, ipamer ipam.IPAMer) (string, string, error) {
	var errors []error
	var err error
	var ipAddress string
	var parentPrefixCidr string
	for _, prefix := range parent.Prefixes {
		if specificIP == "" {
			ipAddress, err = ipamer.AllocateIP(prefix)
		} else {
			ipAddress, err = ipamer.AllocateSpecificIP(prefix, specificIP)
		}
		if err != nil {
			errors = append(errors, err)
			continue
		}
		if ipAddress != "" {
			parentPrefixCidr = prefix.String()
			break
		}
	}
	if ipAddress == "" {
		if len(errors) > 0 {
			return "", "", fmt.Errorf("cannot allocate free ip in ipam: %v", errors)
		}
		return "", "", fmt.Errorf("cannot allocate free ip in ipam")
	}

	return ipAddress, parentPrefixCidr, nil
}
