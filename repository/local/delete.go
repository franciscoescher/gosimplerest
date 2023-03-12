package local

import (
	"fmt"

	"github.com/franciscoescher/gosimplerest/resource"
)

func (r Repository) Delete(b *resource.Resource, id any) error {
	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("no rows affected")
	}
	delete(r.data, id)
	return nil
}
