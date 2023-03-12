package local

import (
	"errors"

	"github.com/franciscoescher/gosimplerest/resource"
)

func (r Repository) Insert(b *resource.Resource, data map[string]any) (int64, error) {

	var pk any
	if b.AutoIncrementalPK {
		r.maxPK = r.maxPK + 1
		pk = r.maxPK
	} else {
		for key, element := range data {
			if key == b.PrimaryKey && !b.AutoIncrementalPK {
				pk = element
			}
		}
	}

	if pk == nil {
		return 0, errors.New("primary key not found")
	}

	// checks if pk already exists
	if _, ok := r.data[pk]; ok {
		return 0, errors.New("primary key already exists")
	}

	r.data[pk] = data

	if b.AutoIncrementalPK {
		return pk.(int64), nil
	}
	return 0, nil
}
