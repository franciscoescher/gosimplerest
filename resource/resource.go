package resource

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	null "gopkg.in/guregu/null.v3"
)

type Base struct {
	Logger   *logrus.Logger
	DB       *sql.DB
	Resource *Resource
	Validate *validator.Validate
}

// gosimplerest.Resource represents a data resource, such as
// a database table, a file in a storage system, etc.
type Resource struct {
	// Table is the name of the table
	Table string `json:"table"`
	// Fields is a list of fields in the table
	Fields map[string]Field `json:"fields"`
	// PrimaryKey is the name of field that is the primary key
	PrimaryKey string `json:"primary_key"`
	// AutoIncrementalPK is a flag that indicates if the primary key is autoincremental
	// and will not use it when inserting a new row
	AutoIncrementalPK bool `json:"incremental_pk"`
	// GeneratePrimaryKeyFunc is a function that generates a new primary key
	// if null, the defaultGeneratePrimaryKeyFunc is used
	GeneratePrimaryKeyFunc GeneratePrimaryKeyFunc `json:"-"`
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
	// Ommmit<Route Type>Route are flags that omit the generation of the specific route from the router
	OmitCreateRoute     bool `json:"omit_create_route"`
	OmitRetrieveRoute   bool `json:"omit_retrieve_route"`
	OmitUpdateRoute     bool `json:"omit_update_route"`
	OmitDeleteRoute     bool `json:"omit_delete_route"`
	OmitSearchRoute     bool `json:"omit_search_route"`
	OmitBelongsToRoutes bool `json:"omit_belongs_to_routes"`
}

type GeneratePrimaryKeyFunc func() any

type Field struct {
	// Validator is the validation rules for the field, using the
	// package github.com/go-playground/validator/v10
	Validator string `json:"validator"`
	// Unsearchable is a flag that indicates that a field can not be used
	// as query parameter in the search route
	Unsearchable bool `json:"unsearchable"`
}

type BelongsTo struct {
	// Table of the other resource that this resource belongs to
	Table string `json:"table"`
	// Field of the current resource that is the foreign key to the table
	Field string `json:"field"`
}

// FromJSON reads a JSON file and populates the model
func (b *Resource) FromJSON(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(file), &b)
	if err != nil {
		return err
	}
	return nil
}

/*
FromStruct reads a struct and populates the resource

The struct should have the following tags:
  - db: used to get the field name
  - json: used to get the field name (if db is not present)
  - pk: used to get the primary key field name:
    set to autoincremental if the primary key is autoincremental
    set to true if the primary key is not autoincremental
  - soft_delete: used to get the soft delete field
  - created_at: used to get the created at field
  - updated_at: used to get the updated at field
  - belongs_to: used to get the belongs to fields
  - validate: used to get the validation rules
  - unsearchable: used to get the unsearchable fields
  - pk: used to get the primary key

The omit flags and GeneratePrimaryKeyFunc are not populated by this function
*/
func (b *Resource) FromStruct(s any) error {
	// Table name
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("not struct type: %s", t.Kind())
	}
	b.Table = strcase.SnakeCase(t.Name())

	// Fields
	// iterate over fields
	fields := make(map[string]Field, t.NumField())
	belongs := make([]BelongsTo, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// get the field name
		name := field.Tag.Get("json")
		db := field.Tag.Get("db")
		if db != "-" && db != "" {
			name = db
		} else if db == "" {
			json := field.Tag.Get("json")
			if json != "-" && json != "" {
				name = json
			} else if json == "" {
				name = strcase.SnakeCase(field.Name)
			}
		}
		// get the field struct
		fields[name] = Field{
			Validator:    field.Tag.Get("validate"),
			Unsearchable: field.Tag.Get("unsearchable") == "true"}
		// get the primary key
		if field.Tag.Get("pk") == "true" {
			b.PrimaryKey = name
		} else if field.Tag.Get("pk") == "autoincremental" {
			b.PrimaryKey = name
			b.AutoIncrementalPK = true
		}
		// get the soft delete field
		if field.Tag.Get("soft_delete") == "true" {
			b.SoftDeleteField = null.StringFrom(name)
		}
		// get the created at field
		if field.Tag.Get("created_at") == "true" {
			b.CreatedAtField = null.StringFrom(name)
		}
		// get the updated at field
		if field.Tag.Get("updated_at") == "true" {
			b.UpdatedAtField = null.StringFrom(name)
		}
		// get the belongs to fields
		if field.Tag.Get("belongs_to") != "" {
			belongs = append(b.BelongsToFields, BelongsTo{
				Field: name,
				Table: field.Tag.Get("belongs_to"),
			})
		}
	}
	b.Fields = fields
	b.BelongsToFields = belongs

	return nil
}

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

func (b *Resource) ValidateFields(v *validator.Validate, data map[string]interface{}) map[string]interface{} {
	rules := make(map[string]interface{}, len(data))
	for k := range data {
		rules[k] = b.Fields[k].Validator
	}
	return v.ValidateMap(data, rules)
}

func (b *Resource) ValidateField(v *validator.Validate, field string, value any) error {
	vf := b.Fields[field].Validator
	if vf == "" {
		return nil
	}
	err := v.Var(value, vf)
	if err != nil {
		return fmt.Errorf("field %s is invalid for validation rule: %s", field, vf)
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
