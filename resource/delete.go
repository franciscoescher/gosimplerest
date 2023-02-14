package resource

import (
	"database/sql"
	"fmt"
)

// Delete deletes a row with the given primary key from the database
func (b *Resource) Delete(base *Base, id string) error {
	var result sql.Result
	err := error(nil)
	if b.SoftDeleteField().Valid {
		result, err = base.DB.Exec(fmt.Sprintf(`UPDATE %s SET %s = NOW() WHERE %s=?`, b.GetName(), b.SoftDeleteField().String, b.PrimaryKey()), id)
		if err != nil {
			return err
		}
	} else {
		result, err = base.DB.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id=?`, b.GetName()), id)
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
