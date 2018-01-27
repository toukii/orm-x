package orm

import (
	"fmt"

	"github.com/didi/gendry/scanner"
)

type ExampleMask struct {
	Name     string `db:"u_name"`
	CardNo   string `db:"no"`
	BankName string `db:"b_name"`
}

func ExampleScan() {
	scanner.SetTagName("db")
	rows, err := Mysql().Query("SELECT u.name AS u_name, uc.no AS uc_no, b.name AS b_name FROM user u INNER JOIN user_card uc ON u.id=uc.user_id INNER JOIN bank b ON uc.bank_id=b.id")

	var mask []*ExampleMask
	err = scanner.Scan(rows, &mask)
	if err != nil {
		panic(err)
	}

	for _, m := range mask {
		fmt.Printf("%+v\n", m)
	}
}
