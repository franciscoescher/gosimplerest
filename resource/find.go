package resource

import (
	"database/sql"
	"fmt"
	"strings"
)

// Find returns a single row from the database, search by the primary key
func (b *Resource) Find(base *Base, id any) (map[string]any, error) {
	fields := b.GetFieldNames()

	response := base.DB.QueryRow(fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ? LIMIT 1`, strings.Join(fields, ","), b.Table, b.PrimaryKey), id)

	values := make([]any, len(b.Fields))
	scanArgs := make([]any, len(b.Fields))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	err := response.Scan(scanArgs...)
	if err != nil {
		if err == sql.ErrNoRows {
			return make(map[string]any, 0), nil
		}
		return make(map[string]any, 0), err
	}

	return b.parseRow(values)
}
