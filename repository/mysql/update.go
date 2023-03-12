package mysql

import (
	"strings"

	"github.com/franciscoescher/gosimplerest/resource"
)

func (r Repository) Update(b *resource.Resource, data map[string]any) (bool, error) {
	fields := make([]string, len(data))
	values := make([]any, len(data))
	i := 0
	for key, element := range data {
		fields[i] = key
		values[i] = element
		i++
	}
	values = append(values, data[b.PrimaryKey])

	sql := concatStr(`UPDATE `, b.Table(), ` SET `, strings.Join(fields, "=?,")+"=?", ` WHERE `, b.PrimaryKey, `=?`)
	result, err := r.db.Exec(sql, values...)
	if err != nil {
		return false, err
	}
	affect, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affect > 0, nil
}
