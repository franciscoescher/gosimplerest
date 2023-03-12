package mysql

import (
	"database/sql"
	"strings"

	"github.com/franciscoescher/gosimplerest/resource"
)

func (r Repository) Find(b *resource.Resource, id any) (map[string]any, error) {
	fields := b.GetFieldNames()

	sqlStatement := concatStr(`SELECT `, strings.Join(fields, ","), ` FROM `, b.Table(), ` WHERE `, b.PrimaryKey, ` = ? LIMIT 1`)
	response := r.db.QueryRow(sqlStatement, id)

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

	return r.parseRow(b, values)
}
