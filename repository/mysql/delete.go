package mysql

import (
	"database/sql"
	"fmt"

	"github.com/franciscoescher/gosimplerest/resource"
)

// Delete deletes a row with the given primary key from the database
func (r Repository) Delete(b *resource.Resource, id string) error {
	var result sql.Result
	err := error(nil)
	if b.SoftDeleteField.Valid {
		sqlStr := concatStr(`UPDATE `, b.Table(), ` SET `, b.SoftDeleteField.String, ` = NOW() WHERE `, b.PrimaryKey, `=?`)
		result, err = r.db.Exec(sqlStr, id)
		if err != nil {
			return err
		}
	} else {
		result, err = r.db.Exec(concatStr(`DELETE FROM `, b.Table(), ` WHERE id=?`), id)
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
