package mysql

import (
	"strings"

	"github.com/franciscoescher/gosimplerest/resource"
)

func (r Repository) Insert(b *resource.Resource, data map[string]any) (int64, error) {
	l := len(data)
	if b.AutoIncrementalPK {
		l--
	}
	in := strings.TrimSuffix(strings.Repeat("?,", l), ",")
	fields := make([]string, l)
	values := make([]any, l)
	i := 0
	for key, element := range data {
		if key == b.PrimaryKey && b.AutoIncrementalPK {
			continue
		}
		fields[i] = key
		values[i] = element
		i++
	}

	sql := concatStr(`INSERT INTO `, b.Table(), ` (`, strings.Join(fields, ","), `) VALUES (`, in, `)`)
	result, err := r.db.Exec(sql, values...)
	if err != nil {
		return 0, err
	}
	if b.AutoIncrementalPK {
		return result.LastInsertId()
	}
	return 0, nil
}
