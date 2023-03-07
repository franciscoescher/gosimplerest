package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/franciscoescher/gosimplerest/repository"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/sirupsen/logrus"
)

// Repository is the implementation of the RepositoryInterface for MySQL database.
type Repository struct {
	db *sql.DB
}

// NewRepository returns a new MySQL Repository
func NewRepository(db *sql.DB) Repository {
	return Repository{db: db}
}

// Compile-time check that Repository implements the Repository interface
var _ repository.RepositoryInterface = (*Repository)(nil)

// ConcatStr concatenates a list of strings
func concatStr(strs ...string) string {
	var sb strings.Builder
	for _, s := range strs {
		sb.WriteString(s)
	}
	return sb.String()
}

// parseRow parses a row from the database, returning a map with
// the field names as keys and the values as values
func (r Repository) parseRow(b *resource.Resource, values []any) (map[string]any, error) {
	fields := b.GetFieldNames()
	result := make(map[string]any, len(b.Fields))
	for i, v := range values {
		casted, err := r.castVal(v)
		if err != nil {
			return result, fmt.Errorf("failed on if for type %T of %v", v, v)
		}
		result[fields[i]] = casted
	}
	return result, nil
}

// castVal casts the value incomming from the database to a valid type
func (r Repository) castVal(v any) (any, error) {
	// if nil, set to nil
	if v == nil {
		return nil, nil
	}

	n3, ok := v.(int64)
	if ok {
		logrus.Info(n3)
		return n3, nil
	}

	n2, ok := v.(float64)
	if ok {
		logrus.Info(n2)
		return n2, nil
	}

	// bool or string
	x, ok := v.([]byte)
	if ok {
		if p, ok := strconv.ParseBool(string(x)); ok == nil {
			return p, nil
		} else {
			return string(x), nil
		}
	}

	t, ok := v.(time.Time)
	if ok {
		return t, nil
	}

	return nil, fmt.Errorf("failed on if for type %T of %v", v, v)
}

// parseRows parses a row from the database, returning a map with the field names as keys and the values as values
func (r Repository) parseRows(b *resource.Resource, rows *sql.Rows) ([]map[string]any, error) {
	results := make([]map[string]any, 0)
	for rows.Next() {
		values := make([]any, len(b.Fields))
		scanArgs := make([]any, len(b.Fields))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err := rows.Scan(scanArgs...)
		if err != nil {
			return make([]map[string]any, 0), err
		}
		result, err := r.parseRow(b, values)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}
	return results, nil
}
