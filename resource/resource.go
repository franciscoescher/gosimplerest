package resource

import (
	"database/sql"
	"fmt"
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

// gosimplerest.Resource represents a database table
type Resource struct {
	// OverrideName is the name of the resource, if null, the name is generated from the struct name
	OverrideName string `json:"override_name"`
	// GeneratePrimaryKeyFunc is a function that generates a new primary key
	// if null, the defaultGeneratePrimaryKeyFunc is used
	GeneratePrimaryKeyFunc GeneratePrimaryKeyFunc `json:"-"`
	// Ommmit<Route Type>Route are flags that omit the generation of the specific route from the router
	OmitCreateRoute     bool `json:"omit_create_route"`
	OmitRetrieveRoute   bool `json:"omit_retrieve_route"`
	OmitUpdateRoute     bool `json:"omit_update_route"`
	OmitDeleteRoute     bool `json:"omit_delete_route"`
	OmitSearchRoute     bool `json:"omit_search_route"`
	OmitBelongsToRoutes bool `json:"omit_belongs_to_routes"`
	// AutoIncrementalPK is a flag that indicates if the primary key is autoincremental
	// and will not use it when inserting a new row
	AutoIncrementalPK bool `json:"incremental_pk"`
	//
	Data any `json:"data"`
}

func (b *Resource) PrimaryKey() string {
	f := b.findTaggedFieldName("primary_key")
	if f.Valid {
		return f.String
	}
	return "id"
}

func (b *Resource) BelongsToFields() []BelongsTo {
	out := make([]BelongsTo, 0)
	t := reflect.TypeOf(b.Data)
	if t == nil {
		return out
	}
	// iterate over fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// get the tag value
		tag := field.Tag.Get("belongs_to")
		// if the tag is not empty, use it as the field name
		if tag != "" {
			out = append(out, BelongsTo{Field: field.Name, Table: tag})
		}
	}
	return out
}

func (b *Resource) SoftDeleteField() null.String {
	return b.findTaggedFieldName("soft_delete")
}

func (b *Resource) CreatedAtField() null.String {
	return b.findTaggedFieldName("created_at")
}

func (b *Resource) UpdatedAtField() null.String {
	return b.findTaggedFieldName("updated_at")
}

func (b *Resource) findTaggedFieldName(tag string) null.String {
	t := reflect.TypeOf(b.Data)
	if t == nil {
		return null.NewString("", false)
	}
	// iterate over fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// get the tag value
		tag := field.Tag.Get(tag)
		// if the tag is not empty, use it as the field name
		if tag == "true" {
			return null.NewString(field.Name, true)
		}
	}
	return null.NewString("", false)
}

// GetName returns the name of the resource using the reflection package
func (b *Resource) GetName() string {
	t := reflect.TypeOf(b.Data)
	if t != nil {
		if b.OverrideName != "" {
			return b.OverrideName
		}
		return strcase.KebabCase(t.Name())
	}
	return ""
}

func (b *Resource) GetFields() []Field {
	t := reflect.TypeOf(b.Data)
	if t == nil {
		return make([]Field, 0)
	}
	// iterate over fields
	fields := make([]Field, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// get the tag value
		tag := field.Tag.Get("json")
		// if the tag is not empty, use it as the field name
		if tag != "" {
			fields[i].Validator = field.Tag.Get("validate")
			fields[i].Unsearchable = field.Tag.Get("unsearchable") == "true"
		}
	}
	return fields
}

func (b *Resource) GetFieldNames() []string {
	t := reflect.TypeOf(b.Data)
	if t == nil {
		return make([]string, 0)
	}
	// iterate over fields
	fields := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	sort.Strings(fields)
	return fields
}

type GeneratePrimaryKeyFunc func() any

type Field struct {
	Validator    string `json:"validator"`
	Unsearchable bool   `json:"unsearchable"`
}

type BelongsTo struct {
	Table string `json:"table"`
	Field string `json:"field"`
}

// HasField returns true if the model has the given field
func (b *Resource) HasField(field string) bool {
	t := reflect.TypeOf(b.Data)
	if t == nil {
		return false
	}
	_, ok := t.FieldByName(field)
	return ok
}

// HasField returns true if the model has the given field
func (b *Resource) IsSearchable(field string) bool {
	t := reflect.TypeOf(b.Data)
	if t == nil {
		return false
	}
	f, ok := t.FieldByName(field)
	return ok && f.Tag.Get("unsearchable") != "true"
}

func (b *Resource) ValidateFields(v *validator.Validate, data map[string]interface{}) (map[string]interface{}, error) {
	t := reflect.TypeOf(b.Data)
	if t == nil {
		return nil, fmt.Errorf("model has no fields")
	}
	rules := make(map[string]interface{}, len(data))
	for k := range data {
		f, ok := t.FieldByName(k)
		if !ok {
			return nil, fmt.Errorf("field %s not found", k)
		}
		rules[k] = f.Tag.Get("validate")
	}
	return v.ValidateMap(data, rules), nil
}

func (b *Resource) ValidateField(v *validator.Validate, field string, value any) error {
	t := reflect.TypeOf(b.Data)
	if t == nil {
		return fmt.Errorf("model has no fields")
	}
	f, ok := t.FieldByName(field)
	if !ok {
		return fmt.Errorf("field %s not found", field)
	}
	vf := f.Tag.Get("validate")
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

// parseRow parses a row from the database, returning a map with
// the field names as keys and the values as values
func (b *Resource) parseRow(values []any) (map[string]any, error) {
	fields := b.GetFieldNames()
	result := make(map[string]any, len(b.GetFieldNames()))
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
		values := make([]any, len(b.GetFieldNames()))
		scanArgs := make([]any, len(b.GetFieldNames()))
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
