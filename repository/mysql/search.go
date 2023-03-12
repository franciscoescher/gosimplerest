package mysql

import (
	"strings"

	"github.com/franciscoescher/gosimplerest/resource"
)

func (r Repository) Search(b *resource.Resource, query map[string][]string) ([]map[string]any, error) {
	fields := b.GetFieldNames()

	// build query
	where := make([]string, len(query))
	values := make([]any, 0)
	i := 0
	for field, value := range query {
		if len(value) == 1 {
			where[i] = concatStr(field, " = ?")
			values = append(values, value[0])
		} else {
			where[i] = concatStr(field, " IN (", strings.Repeat("?,", len(value)-1)+"?", ")")
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
	sqlStr := concatStr(`SELECT `, strings.Join(fields, ","), ` FROM `, b.Table(), ` `, whereStr, ` ORDER BY `, b.PrimaryKey)
	response, err := r.db.Query(sqlStr, values...)
	if err != nil {
		return nil, err
	}
	defer response.Close()
	return r.parseRows(b, response)
}
