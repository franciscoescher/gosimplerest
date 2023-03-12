package local

import (
	"github.com/franciscoescher/gosimplerest/resource"
)

func (r Repository) Find(b *resource.Resource, id any) (map[string]any, error) {
	row, ok := r.data[id]
	if !ok {
		return make(map[string]any, 0), nil
	}
	return row, nil
}
