package local

import (
	"github.com/franciscoescher/gosimplerest/resource"
)

func (r Repository) Update(b *resource.Resource, data map[string]any) (int64, error) {

	if _, ok := data[b.PrimaryKey]; !ok {
		return 0, nil
	}

	inPlaceData := r.data[data[b.PrimaryKey]]
	for key, element := range data {
		inPlaceData[key] = element
	}

	r.data[data[b.PrimaryKey]] = inPlaceData

	return 1, nil
}
