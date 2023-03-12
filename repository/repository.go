package repository

import "github.com/franciscoescher/gosimplerest/resource"

type RepositoryInterface interface {
	// Delete deletes a row with the given primary key from the database
	Delete(b *resource.Resource, id any) error
	// Find returns a single row from the database, search by the primary key
	Find(b *resource.Resource, id any) (map[string]any, error)
	// Insert inserts a new row into the database
	Insert(b *resource.Resource, data map[string]any) (int64, error)
	// Search searches for rows in the database with where clauses
	Search(b *resource.Resource, query map[string][]string) ([]map[string]any, error)
	// Update updates a row in the database, one of the fields must be the primary key
	Update(b *resource.Resource, data map[string]any) (int64, error)
}
