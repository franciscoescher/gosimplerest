package local

import (
	"github.com/franciscoescher/gosimplerest/resource"
)

func (r Repository) Update(b *resource.Resource, data map[string]any) (bool, error) {

	if _, ok := data[b.PrimaryKey]; !ok {
		return false, nil
	}

	inPlaceData := r.data[data[b.PrimaryKey]]
	for key, element := range data {
		inPlaceData[key] = element
	}

	r.data[data[b.PrimaryKey]] = inPlaceData

	return true, nil
}
