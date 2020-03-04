package service

import (
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
)

// Some predefined users
var (
	BasePath = "/"
)

type WebResource struct {
	DS *datastore.RethinkStore
}
