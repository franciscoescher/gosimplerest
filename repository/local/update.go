package local

import (
	"errors"

	"github.com/franciscoescher/gosimplerest/resource"
)

func (r Repository) Update(b *resource.Resource, data map[string]any) (bool, error) {

	// if the primary key is not in the data, return false
	if _, ok := data[b.PrimaryKey]; !ok {
		return false, errors.New("primary key not in data")
	}

	inPlaceData, ok := r.data[data[b.PrimaryKey]]

	// if the row does not exist, return false
	if !ok {
		return false, nil
	}

	for key, element := range data {
		inPlaceData[key] = element
	}

	r.data[data[b.PrimaryKey]] = inPlaceData

	return true, nil
}
