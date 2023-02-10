package gosimplerest

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	null "gopkg.in/guregu/null.v3"
)

// Resource represents a database table
type Resource struct {
	// Table is the name of the table
	Table string
	// PrimaryKey is the name of field that is the primary key
	PrimaryKey string
	// Fields is a list of fields in the table
	Fields map[string]Field
	// SoftDeleteField is the name of the field that is used for soft deletes
	// if null, no soft deletes are used
	SoftDeleteField null.String
	// CreatedAtField is the name of the field that is used as createion timestamp
	// if null, no creation timestamp is generated
	CreatedAtField null.String
	// UpdatedAtField is the name of the field that is used as update timestamp
	// if null, no update timestamp is generated
	UpdatedAtField null.String
	// BelongsToFields is a list of fields represent a belonging relation with another table,
	// usually also foreign keys to other tables
	BelongsToFields []BelongsTo
	// GeneratePrimaryKeyFunc is a function that generates a new primary key
	// if null, the defaultGeneratePrimaryKeyFunc is used
	GeneratePrimaryKeyFunc func() any
}

type Field struct {
	Validator    ValidatorFunc
	Unsearchable bool
}

type BelongsTo struct {
	Table string
	Field string
}

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

// Find returns a single row from the database, search by the primary key
func (b *Resource) Find(id any) (map[string]any, error) {
	fields := b.GetFieldNames()

	response := db.QueryRow(fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ? LIMIT 1`, strings.Join(fields, ","), b.Table, b.PrimaryKey), id)

	result := make(map[string]any, len(b.Fields))
	values := make([]any, len(b.Fields))
	scanArgs := make([]any, len(b.Fields))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	err := response.Scan(scanArgs...)
	if err != nil {
		return result, err
	}

	return b.parseRow(values)
}

// parseRow parses a row from the database, returning a map with
// the field names as keys and the values as values
func (b *Resource) parseRow(values []any) (map[string]any, error) {
	fields := b.GetFieldNames()
	result := make(map[string]any, len(b.Fields))
	for i, v := range values {
		// if nil, set to nil
		if v == nil {
			result[fields[i]] = nil
			continue
		}

		// if int
		n2, ok := v.(float64)
		if ok {
			result[fields[i]] = n2
			continue
		}

		// bool or string
		x, ok := v.([]byte)
		if ok {
			if p, ok := strconv.ParseBool(string(x)); ok == nil {
				result[fields[i]] = p
			} else if fmt.Sprintf("%T", string(x)) == "string" {
				result[fields[i]] = string(x)
			} else {
				return result, fmt.Errorf("failed on if for type %T of %v", x, x)
			}
			continue
		}

		// if time
		t, ok := v.(time.Time)
		if ok {
			result[fields[i]] = t
			continue
		}

		// if int64
		n, ok := v.(int64)
		if ok {
			result[fields[i]] = n
			continue
		}

		return result, fmt.Errorf("unmapped value (%b) field type for %s for model %s", v, fields[i], b.Table)
	}
	return result, nil
}

func (b *Resource) parseRows(response *sql.Rows) ([]map[string]any, error) {
	results := make([]map[string]any, 0)
	for response.Next() {
		values := make([]any, len(b.Fields))
		scanArgs := make([]any, len(b.Fields))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		response.Scan(scanArgs...)
		result, err := b.parseRow(values)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}
	return results, nil
}

// Delete deletes a row with the given primary key from the database
func (b *Resource) Delete(id string) error {
	var result sql.Result
	err := error(nil)
	if b.SoftDeleteField.Valid {
		result, err = db.Exec(fmt.Sprintf(`UPDATE %s SET %s = NOW() WHERE %s=?`, b.Table, b.SoftDeleteField.String, b.PrimaryKey), id)
		if err != nil {
			return err
		}
	} else {
		result, err = db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id=?`, b.Table), id)
		if err != nil {
			return err
		}
	}
	affect, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affect == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

// Insert inserts a new row into the database
func (b *Resource) Insert(data map[string]any) error {
	in := strings.TrimSuffix(strings.Repeat("?,", len(data)), ",")
	fields := make([]string, len(data))
	values := make([]any, len(data))
	i := 0
	for key, element := range data {
		fields[i] = key
		values[i] = element
		i++
	}

	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, b.Table, strings.Join(fields, ","), in)
	_, err := db.Exec(sql, values...)
	return err
}

// Update updates a row in the database
// data must contain the primary key
func (b *Resource) Update(data map[string]any) (int64, error) {
	fields := make([]string, len(data))
	values := make([]any, len(data))
	i := 0
	for key, element := range data {
		fields[i] = key
		values[i] = element
		i++
	}
	values = append(values, data[b.PrimaryKey])

	sql := fmt.Sprintf(`UPDATE %s SET %s WHERE %s=?`, b.Table, strings.Join(fields, "=?,")+"=?", b.PrimaryKey)
	result, err := db.Exec(sql, values...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}

// FindFromBelongsTo finds all rows of a model with the belongsTo relationship
func (b *Resource) FindFromBelongsTo(id any, belongsTo BelongsTo) ([]map[string]any, error) {
	fields := b.GetFieldNames()

	response, err := db.Query(fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?`, strings.Join(fields, ","), b.Table, belongsTo.Field), id)
	if err != nil {
		return nil, err
	}
	defer response.Close()
	return b.parseRows(response)
}

// Search searches for rows in the database with where clauses
func (b *Resource) Search(query map[string][]string) ([]map[string]any, error) {
	fields := b.GetFieldNames()

	// build query
	where := make([]string, len(query))
	values := make([]any, 0)
	i := 0
	for field, value := range query {
		if len(value) == 1 {
			where[i] = fmt.Sprintf("%s = ?", field)
			values = append(values, value[0])
		} else {
			where[i] = fmt.Sprintf("%s IN (%s)", field, strings.Repeat("?,", len(value)-1)+"?")
			for _, v := range value {
				values = append(values, v)
			}
		}
		i++
	}
	whereStr := ""
	if len(where) > 0 {
		whereStr = "WHERE " + strings.Join(where, " AND ")
	}
	response, err := db.Query(fmt.Sprintf(`SELECT %s FROM %s %s`, strings.Join(fields, ","), b.Table, whereStr), values...)
	if err != nil {
		return nil, err
	}
	defer response.Close()
	return b.parseRows(response)
}
