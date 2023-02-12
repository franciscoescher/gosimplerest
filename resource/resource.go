package resource

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

type Base struct {
	Logger   *logrus.Logger
	DB       *sql.DB
	Resource *Resource
}

// gosimplerest.Resource represents a database table
type Resource struct {
	// Table is the name of the table
	Table string `json:"table"`
	// PrimaryKey is the name of field that is the primary key
	PrimaryKey string `json:"primary_key"`
	// Fields is a list of fields in the table
	Fields map[string]Field `json:"fields"`
	// SoftDeleteField is the name of the field that is used for soft deletes
	// if null, no soft deletes are used
	SoftDeleteField null.String `json:"soft_delete_field"`
	// CreatedAtField is the name of the field that is used as createion timestamp
	// if null, no creation timestamp is generated
	CreatedAtField null.String `json:"created_at_field"`
	// UpdatedAtField is the name of the field that is used as update timestamp
	// if null, no update timestamp is generated
	UpdatedAtField null.String `json:"updated_at_field"`
	// BelongsToFields is a list of fields represent a belonging relation with another table,
	// usually also foreign keys to other tables
	BelongsToFields []BelongsTo `json:"belongs_to_fields"`
	// GeneratePrimaryKeyFunc is a function that generates a new primary key
	// if null, the defaultGeneratePrimaryKeyFunc is used
	GeneratePrimaryKeyFunc func() any `json:"-"`
	// Ommmit<Route Type>Route are flags that omit the generation of the specific route from the router
	OmitCreateRoute     bool `json:"omit_create_route"`
	OmitRetrieveRoute   bool `json:"omit_retrieve_route"`
	OmitUpdateRoute     bool `json:"omit_update_route"`
	OmitDeleteRoute     bool `json:"omit_delete_route"`
	OmitSearchRoute     bool `json:"omit_search_route"`
	OmitBelongsToRoutes bool `json:"omit_belongs_to_routes"`
}

type Field struct {
	Validator    ValidatorFunc `json:"-"`
	Unsearchable bool          `json:"unsearchable"`
}

type BelongsTo struct {
	Table string `json:"table"`
	Field string `json:"field"`
}

// ValidatorFunc is a function that validates a field
// This function should expect the value to be either string (for the query routes) or
// the format that the database driver expects (for the insert/update routes)
type ValidatorFunc func(field string, val any) error

// HasField returns true if the model has the given field
func (b *Resource) HasField(field string) bool {
	_, ok := b.Fields[field]
	return ok
}

// HasField returns true if the model has the given field
func (b *Resource) IsSearchable(field string) bool {
	val, ok := b.Fields[field]
	if !ok {
		return false
	}
	return !val.Unsearchable
}

func (b *Resource) ValidateField(field string, value any) error {
	vf := b.Fields[field].Validator
	if vf != nil {
		err := vf(field, value)
		return err
	}
	return nil
}

// GeneratePrimaryKey generates a new primary key
func (b *Resource) GeneratePrimaryKey() any {
	if b.GeneratePrimaryKeyFunc != nil {
		return b.GeneratePrimaryKeyFunc()
	}
	return b.defaultGeneratePrimaryKeyFunc()
}

func (b *Resource) defaultGeneratePrimaryKeyFunc() string {
	id, _ := uuid.NewV4()
	return id.String()
}

func (b *Resource) GetFieldNames() []string {
	fields := make([]string, len(b.Fields))
	i := 0
	for field := range b.Fields {
		fields[i] = field
		i++
	}
	sort.Strings(fields)
	return fields
}

// parseRow parses a row from the database, returning a map with
// the field names as keys and the values as values
func (b *Resource) parseRow(values []any) (map[string]any, error) {
	fields := b.GetFieldNames()
	result := make(map[string]any, len(b.Fields))
	for i, v := range values {
		casted, err := castVal(v)
		if err != nil {
			return result, fmt.Errorf("failed on if for type %T of %v", v, v)
		}
		result[fields[i]] = casted
	}
	return result, nil
}

func castVal(v any) (any, error) {
	// if nil, set to nil
	if v == nil {
		return nil, nil
	}

	n, ok := v.(int)
	if ok {
		logrus.Info("AQUIIII1")
		logrus.Info(n)
		return n, nil
	}

	n3, ok := v.(int64)
	if ok {
		logrus.Info("AQUIIII2")
		logrus.Info(n3)
		return n3, nil
	}

	n2, ok := v.(float64)
	if ok {
		logrus.Info("AQUIIII3")
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

func (b *Resource) parseRows(rows *sql.Rows) ([]map[string]any, error) {
	results := make([]map[string]any, 0)
	for rows.Next() {
		values := make([]any, len(b.Fields))
		scanArgs := make([]any, len(b.Fields))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		rows.Scan(scanArgs...)
		result, err := b.parseRow(values)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (b *Resource) FromJSON(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(file), &b)
}
