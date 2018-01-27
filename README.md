# orm-x

orm for `join on` table, sql like this:

## output

```
SELECT u.name AS u_name, uc.no AS uc_no, b.name AS b_name FROM user u INNER JOIN user_card uc ON u.id=uc.user_id INNER JOIN bank b ON uc.bank_id=b.id WHERE u.name='toukii'
```

Because refer to multi tables, columns should have an alias(use the default alias: [table alias]\_[column name]).

## usage

```
	m := NewMask("Mask", user, userCard, bank)
	user.Select(m).InnerJoin(userCard, "u.id=uc.user_id").InnerJoin(bank, "uc.bank_id=b.id").Where("u.name='toukii'")
	fmt.Println(user.Sql())
```

`user`,`userCard`,`bank` are tableModel:

 - 1 use yaml parse

```
	tb, err := MultiParse("test_yaml/model.yaml")

	user := tb["User"]
	userCard := tb["UserCard"]
	bank := tb["Bank"]
```

 - 2 new

```
	user := NewT("user", "u", []string{"id", "name"})
	user.Mask["Mask"] = []string{"name"}
```


You may want this usefull lib (sq.Rows scanner): __[github.com/didi/gendry/scanner](https://github.com/didi/gendry/tree/master/scanner)__


_scanner example code:_

```
package orm

import (
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
```

result:

```
&{Name:toukii CardNo:6214xxxxxxxx8890 BankName:中国招商银行}
&{Name:toukii CardNo:6214xxxxxxxx5760 BankName:上海浦发银行}
```