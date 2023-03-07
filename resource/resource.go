package resource

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/stoewer/go-strcase"
	null "gopkg.in/guregu/null.v3"
)

// gosimplerest.Resource represents a data resource, such as
// a database table, a file in a storage system, etc.
type Resource struct {
	// Name of the resource
	Name string `json:"name"`
	// OverwriteTableName is the name of the table in case it is not the same as the resource name
	OverwriteTableName null.String `json:"overwrite_table_name"`
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
	// Ommmit<Route Type>Route are flags that omit the generation of the specific route from the router
	OmitCreateRoute        bool `json:"omit_create_route"`
	OmitRetrieveRoute      bool `json:"omit_retrieve_route"`
	OmitUpdateRoute        bool `json:"omit_update_route"`
	OmitPartialUpdateRoute bool `json:"omit_partial_update_route"`
	OmitDeleteRoute        bool `json:"omit_delete_route"`
	OmitSearchRoute        bool `json:"omit_search_route"`
	OmitHeadRoutes         bool `json:"omit_head_routes"`
}

type GeneratePrimaryKeyFunc func() any

type Field struct {
	// Validator is the validation rules for the field, using the
	// package github.com/go-playground/validator/v10
	Validator string `json:"validator"`
	// Unsearchable is a flag that indicates that a field can not be used
	// as query parameter in the search route
	Unsearchable bool `json:"unsearchable"`
	// Immutable is a flag that indicates that a field can not be updated
	Immutable bool `json:"immutable"`
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

// Table returns the table name
func (b *Resource) Table() string {
	if b.OverwriteTableName.Valid {
		return b.OverwriteTableName.String
	}
	return b.Name
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
  - validate: used to get the validation rules
  - unsearchable: used to get the unsearchable fields
  - pk: used to get the primary key

The omit route flags, OverwriteTableName and GeneratePrimaryKeyFunc are not populated by this function
*/
func (b *Resource) FromStruct(s any) error {
	// Table name
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("not struct type: %s", t.Kind())
	}
	b.Name = strcase.SnakeCase(t.Name())
	b.OverwriteTableName = null.NewString("", false)

	// Fields
	// iterate over fields
	fields := make(map[string]Field, t.NumField())
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
		presentOrTrue := func(tag string) bool {
			val, ok := field.Tag.Lookup(tag)
			return ok && (val == "" || val == "true")
		}
		// get the field struct
		fields[name] = Field{
			Validator:    field.Tag.Get("validate"),
			Immutable:    presentOrTrue("immutable"),
			Unsearchable: presentOrTrue("unsearchable"),
		}
		// get the primary key
		if presentOrTrue("pk") {
			b.PrimaryKey = name
		} else if field.Tag.Get("pk") == "autoincremental" {
			b.PrimaryKey = name
			b.AutoIncrementalPK = true
		}
		// get the soft delete field
		if presentOrTrue("soft_delete") {
			b.SoftDeleteField = null.StringFrom(name)
		}
		// get the created at field
		if presentOrTrue("created_at") {
			b.CreatedAtField = null.StringFrom(name)
		}
		// get the updated at field
		if presentOrTrue("updated_at") {
			b.UpdatedAtField = null.StringFrom(name)
		}
	}
	b.Fields = fields

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

// ValidateAllFields validates all fields of the model against the given data
func (b *Resource) ValidateAllFields(v *validator.Validate, data map[string]interface{}) map[string]interface{} {
	values := make(map[string]interface{}, len(data))
	rules := make(map[string]interface{}, len(data))
	for field := range b.Fields {
		if _, ok := data[field]; ok {
			values[field] = data[field]
		} else {
			values[field] = nil
		}
		rules[field] = b.Fields[field].Validator
	}
	return v.ValidateMap(values, rules)
}

// ValidateInputFields validates the given fields of the model against the given data
func (b *Resource) ValidateInputFields(v *validator.Validate, data map[string]interface{}) map[string]interface{} {
	rules := make(map[string]interface{}, len(data))
	for k := range data {
		rules[k] = b.Fields[k].Validator
	}
	return v.ValidateMap(data, rules)
}

// ValidateField validates the given field of the model against the given data
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

// defaultGeneratePrimaryKeyFunc generates a new primary key using the uuid package
func (b *Resource) defaultGeneratePrimaryKeyFunc() string {
	id, _ := uuid.NewV4()
	return id.String()
}

// GetFieldNames returns a list of strings with the field names of the model
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
