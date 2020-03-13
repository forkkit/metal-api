package size

import (
	"github.com/emicklei/go-restful"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/datastore"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
	v1 "github.com/metal-stack/metal-api/pkg/proto/v1"
)

type sizeResource struct {
	ds *datastore.RethinkStore
}

// NewSizeService returns a webservice for size specific endpoints.
func NewSizeService(ds *datastore.RethinkStore) *restful.WebService {
	r := sizeResource{
		ds: ds,
	}
	return r.webService()
}

func mapSizeConstraintType(constraint v1.SizeConstraint_Type) metal.ConstraintType {
	t := metal.CoreConstraint
	switch constraint {
	case v1.SizeConstraint_MEMORY:
		t = metal.MemoryConstraint
	case v1.SizeConstraint_STORAGE:
		t = metal.StorageConstraint
	}
	return t
}
