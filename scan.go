package orm

import (
	"database/sql"
	"fmt"

	"github.com/didi/gendry/scanner"
)

type ExampleMask struct {
	Name     string `db:"u_name"`
	CardNo   string `db:"uc_no"`
	BankName string `db:"b_name"`
}

func ExampleScan() {
	scanner.SetTagName("db")
	rows, err := Mysql().Query("SELECT u.name AS u_name, uc.no AS uc_no, b.name AS b_name FROM user u INNER JOIN user_card uc ON u.id=uc.user_id INNER JOIN bank b ON uc.bank_id=b.id WHERE u.name='toukii'")

	var mask []*ExampleMask
	err = scanner.Scan(rows, &mask)
	if err != nil {
		panic(err)
	}

	for _, m := range mask {
		fmt.Printf("%+v\n", m)
	}
}

type ScanLoopFunc func(*sql.Rows)
type ScanEachFunc func(rows *sql.Rows) error

// if objOrFunc is func: either ScanLoopFunc or ScanEachFunc; or data-struct
func (db *DBStore) Scan(t *T, objOrFunc interface{}, args ...interface{}) error {
	rows, err := db.Query(t.Sql(), args...)
	if err != nil {
		return err
	}
	scanLoopFunc, ok := objOrFunc.(ScanLoopFunc)
	if ok {
		scanLoopFunc(rows)
		return nil
	}
	scanEachFunc, ok := objOrFunc.(ScanEachFunc)
	if ok {
		for rows.Next() {
			err := scanEachFunc(rows)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return scanner.Scan(rows, objOrFunc)
}
