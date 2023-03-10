package local

import (
	"github.com/franciscoescher/gosimplerest/repository"
)

// Repository is the implementation of the RepositoryInterface for local in memory database.
// It is a simple map of maps where the key of external map is the primary key
// and the key of the internal map is the column name.
// Only use it for testing purposes.
type Repository struct {
	// data is the local database
	data map[any]map[string]any
	// pk counter in case auto incremental pk
	maxPK int64
}

// NewRepository returns a new local Repository
func NewRepository() Repository {
	return Repository{data: make(map[any]map[string]any, 0)}
}

// NewRepository returns a new local Repository
func NewRepositoryWithData(data map[any]map[string]any) Repository {
	return Repository{
		data: data,
	}
}

// Compile-time check that Repository implements the Repository interface
var _ repository.RepositoryInterface = (*Repository)(nil)
