package resource

import (
	"fmt"
	"strings"
)

// Update updates a row in the database
// data must contain the primary key
func (b *Resource) Update(base *Base, data map[string]any) (int64, error) {
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
	result, err := base.DB.Exec(sql, values...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affect, nil
}
