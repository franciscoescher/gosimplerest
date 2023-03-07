package repository

import "github.com/franciscoescher/gosimplerest/resource"

type RepositoryInterface interface {
	Delete(b *resource.Resource, id string) error
	Find(b *resource.Resource, id any) (map[string]any, error)
	Insert(b *resource.Resource, data map[string]any) (int64, error)
	Search(b *resource.Resource, query map[string][]string) ([]map[string]any, error)
	Update(b *resource.Resource, data map[string]any) (int64, error)
}
