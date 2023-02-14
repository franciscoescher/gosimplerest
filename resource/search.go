package resource

import (
	"fmt"
	"strings"
)

// Search searches for rows in the database with where clauses
func (b *Resource) Search(base *Base, query map[string][]string) ([]map[string]any, error) {
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
	response, err := base.DB.Query(fmt.Sprintf(`SELECT %s FROM %s %s ORDER BY %s`, strings.Join(fields, ","),
		b.GetName(), whereStr, b.PrimaryKey()), values...)
	if err != nil {
		return nil, err
	}
	defer response.Close()
	return b.parseRows(response)
}
