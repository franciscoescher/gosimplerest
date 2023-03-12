package repository

import "github.com/franciscoescher/gosimplerest/resource"

type RepositoryInterface interface {
	// Delete deletes a row with the given primary key from the database
	Delete(b *resource.Resource, id any) error
	// Find returns a single row from the database, search by the primary key
	// return 0 rows if not found, but no error
	Find(b *resource.Resource, id any) (map[string]any, error)
	// Insert inserts a new row into the database
	// returns pk only if auto incremental
	Insert(b *resource.Resource, data map[string]any) (int64, error)
	// Search searches for rows in the database using the query parameters
	// returns 0 rows if not found, but no error
	// query is a map of field names and values
	// multiple values for the same field are ORed
	Search(b *resource.Resource, query map[string][]string) ([]map[string]any, error)
	// Update updates a row in the database
	// One of the fields must be the primary key or it will return an error
	// Returns true if the a row was updated, false if not found
	Update(b *resource.Resource, data map[string]any) (bool, error)
}
