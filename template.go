package gosimplerest

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	null "gopkg.in/guregu/null.v3"
)

// Resource represents a database table
type Resource struct {
	// Table is the name of the table
	Table string
	// PrimaryKey is the name of field that is the primary key
	PrimaryKey string
	// Fields is a list of fields in the table
	Fields []Field
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
	GeneratePrimaryKeyFunc func() interface{}
}

type Field struct {
	Name string
}

type BelongsTo struct {
	Table string
	Field string
}

// HasField returns true if the model has the given field
func (b *Resource) HasField(field string) bool {
	for _, f := range b.Fields {
		if f.Name == field {
			return true
		}
	}
	return false
}

// GeneratePrimaryKey generates a new primary key
func (b *Resource) GeneratePrimaryKey() interface{} {
	if b.GeneratePrimaryKeyFunc != nil {
		return b.GeneratePrimaryKeyFunc()
	}
	return b.defaultGeneratePrimaryKeyFunc()
}

func (b *Resource) defaultGeneratePrimaryKeyFunc() string {
	return b.PrimaryKey
}

// Find returns a single row from the database, search by the primary key
func (b *Resource) Find(id interface{}) (map[string]interface{}, error) {
	fields := make([]string, len(b.Fields))
	for i, field := range b.Fields {
		fields[i] = field.Name
	}

	response := db.QueryRow(fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ? LIMIT 1`, strings.Join(fields, ","), b.Table, b.PrimaryKey), id)

	result := make(map[string]interface{}, len(b.Fields))
	values := make([]interface{}, len(b.Fields))
	scanArgs := make([]interface{}, len(b.Fields))
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
func (b *Resource) parseRow(values []interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(b.Fields))
	for i, v := range values {
		// if nil, set to nil
		if v == nil {
			result[b.Fields[i].Name] = nil
			continue
		}

		// if number, bool or string
		x, ok := v.([]byte)
		if ok {
			if p, ok := strconv.ParseFloat(string(x), 64); ok == nil {
				result[b.Fields[i].Name] = p
			} else if p, ok := strconv.ParseBool(string(x)); ok == nil {
				result[b.Fields[i].Name] = p
			} else if fmt.Sprintf("%T", string(x)) == "string" {
				result[b.Fields[i].Name] = string(x)
			} else {
				return result, fmt.Errorf("failed on if for type %T of %v", x, x)
			}
			continue
		}

		// if time
		t, ok := v.(time.Time)
		if ok {
			result[b.Fields[i].Name] = t
			continue
		}

		// if int
		n, ok := v.(int64)
		if ok {
			result[b.Fields[i].Name] = n
			continue
		}

		return result, fmt.Errorf("unmapped type for model %s", b.Table)
	}
	return result, nil
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
func (b *Resource) Insert(data map[string]interface{}) error {
	in := strings.TrimSuffix(strings.Repeat("?,", len(data)), ",")
	fields := make([]string, len(data))
	values := make([]interface{}, len(data))
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
func (b *Resource) Update(data map[string]interface{}) (int64, error) {
	fields := make([]string, len(data))
	values := make([]interface{}, len(data))
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
func (b *Resource) FindFromBelongsTo(id interface{}, belongsTo BelongsTo) ([]map[string]interface{}, error) {
	fields := make([]string, len(b.Fields))
	for i, field := range b.Fields {
		fields[i] = field.Name
	}

	response, err := db.Query(fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?`, strings.Join(fields, ","), b.Table, belongsTo.Field), id)
	if err != nil {
		return nil, err
	}
	defer response.Close()

	results := make([]map[string]interface{}, 0)
	for response.Next() {
		values := make([]interface{}, len(b.Fields))
		scanArgs := make([]interface{}, len(b.Fields))
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
