package resource

import (
	"database/sql"
	"strings"
)

// Find returns a single row from the database, search by the primary key
func (b *Resource) Find(base *Base, id any) (map[string]any, error) {
	fields := b.GetFieldNames()

	sqlStatement := ConcatStr(`SELECT `, strings.Join(fields, ","), ` FROM `, b.Table, ` WHERE `, b.PrimaryKey, ` = ? LIMIT 1`)
	response := base.DB.QueryRow(sqlStatement, id)

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
