package local

import (
	"github.com/franciscoescher/gosimplerest/resource"
)

func (r Repository) Search(b *resource.Resource, query map[string][]string) ([]map[string]any, error) {
	results := make([]map[string]any, 0)

	for _, row := range r.data {
		match := false
		for field, value := range query {
			for _, v := range value {
				if row[field] == v {
					match = true
					break
				}
			}
		}
		if match {
			results = append(results, row)
		}
	}
	return results, nil
}
