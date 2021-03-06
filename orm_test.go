package orm

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/didi/gendry/scanner"
	"github.com/toukii/assert"
)

// = - = - = - = - = - = - = - = -D=E-F=I-N=E- = - = - = - = - = - = - = - =
type UserCardMask struct {
	Name     string `db:"u_name"`
	CardNo   string `db:"uc_no"`
	BankName string `db:"b_name"`
}

// = - = - = - = - = - = - = - = -D=B- = - = - = - = - = - = - = - = - = - =
func init() {
	err := SetUpMysql(&Config{
		User:     "root",
		Passwd:   "ezorm_pass",
		Host:     "localhost",
		Port:     3306,
		Database: "jikang",
	})
	equal.Equal(nil, err, nil)

	scanner.SetTagName("db")
}

func TestQuery(t *testing.T) {
	tb, err := MultiParse("test_yaml/model.yaml")
	if err != nil {
		t.Error(err)
	}

	user := tb["User"]
	userCard := tb["UserCard"]
	bank := tb["Bank"]

	m := NewMask("Mask", user, userCard, bank)
	user.Select(m).InnerJoin(userCard, "u.id=uc.user_id").InnerJoin(bank, "uc.bank_id=b.id").Where("u.name='toukii'")
	fmt.Println(user.Sql())

	rows, err := Mysql().Query(user.Sql())
	equal.Equal(nil, err, nil)

	var mask []*UserCardMask
	err = scanner.Scan(rows, &mask)
	equal.Equal(nil, err, nil)
	for _, m := range mask {
		fmt.Printf("%+v\n", m)
	}

	user.Select(m).InnerJoin(userCard, "u.id=uc.user_id").InnerJoin(bank, "uc.bank_id=b.id").Where("u.name=?")
	rows, err = Mysql().Query(user.Sql(), "toukii")
	equal.Equal(nil, err, nil)

	var mask2 []*UserCardMask
	err = scanner.Scan(rows, &mask2)
	equal.Equal(nil, err, nil)
	for _, m := range mask2 {
		fmt.Printf("%+v\n", m)
	}
}

func TestExample(t *testing.T) {
	ExampleScan()
}

func TestScan(t *testing.T) {
	tb, err := MultiParse("test_yaml/model.yaml")
	if err != nil {
		t.Error(err)
	}

	user := tb["User"]
	userCard := tb["UserCard"]
	bank := tb["Bank"]

	m := NewMask("Mask", user, userCard, bank)

	var mask []*UserCardMask
	err = Mysql().Scan(user.Select(m).InnerJoin(userCard, "u.id=uc.user_id").InnerJoin(bank, "uc.bank_id=b.id").Where("u.name='toukii'"), &mask)
	equal.Equal(nil, err, nil)
	for _, m := range mask {
		fmt.Printf("%+v\n", m)
	}

	var mask2 []*UserCardMask
	err = Mysql().Scan(user.Select(m).InnerJoin(userCard, "u.id=uc.user_id").InnerJoin(bank, "uc.bank_id=b.id").Where("u.name=?"), &mask2, "toukii")
	equal.Equal(nil, err, nil)
	for _, m := range mask2 {
		fmt.Printf("%+v\n", m)
	}

	var mask3 []*UserCardMask
	scanLoopFun := ScanLoopFunc(func(rows *sql.Rows) {
		for rows.Next() {
			var mk UserCardMask
			err := rows.Scan(&mk.Name, &mk.CardNo, &mk.BankName)
			if err != nil {
				fmt.Println(err)
				continue
			}
			mask3 = append(mask3, &mk)
		}
	})

	err = Mysql().Scan(user.Select(m).InnerJoin(userCard, "u.id=uc.user_id").InnerJoin(bank, "uc.bank_id=b.id").Where("u.name=?"), scanLoopFun, "toukii")
	equal.Equal(nil, err, nil)
	for _, m := range mask3 {
		fmt.Printf("%+v\n", m)
	}

	var mask4 []*UserCardMask
	scanEachFun := ScanEachFunc(func(rows *sql.Rows) error {
		var mk UserCardMask
		err := rows.Scan(&mk.Name, &mk.CardNo, &mk.BankName)
		if err != nil {
			return err
		}
		mask4 = append(mask4, &mk)
		return nil
	})

	err = Mysql().Scan(user.Select(m).InnerJoin(userCard, "u.id=uc.user_id").InnerJoin(bank, "uc.bank_id=b.id").Where("u.name=?"), scanEachFun, "toukii")
	equal.Equal(nil, err, nil)
	for _, m := range mask4 {
		fmt.Printf("%+v\n", m)
	}
}

// = - = - = - = - = - = - = - = -S=Q-L= - = - = - = - = - = - = - = - = - =
func TestSelect(t *testing.T) {
	t1 := NewT("user", "u", []string{"id", "name"})

	t1.Select(nil)
	equal.Equal(nil, t1.Sql(), "SELECT u.id AS u_id, u.name AS u_name FROM user u")
}

func TestOnlyJoin(t *testing.T) {
	t1 := NewT("user", "u", []string{"id", "name"})
	t2 := NewT("user_card", "uc", []string{"user_id", "no"})

	t1.InnerJoin(t2, "u.id=uc.user_id")
	equal.Equal(nil, t1.Sql(), " INNER JOIN user_card uc ON u.id=uc.user_id")
}

func TestInnerJoin(t *testing.T) {
	t1 := NewT("user", "u", []string{"id", "name"})
	t2 := NewT("user_card", "uc", []string{"user_id", "no"})
	mask := NewMask("", t1, t2)

	t1.Select(mask).InnerJoin(t2, "u.id=uc.user_id").Where("u.name='toukii'")
	equal.Equal(nil, t1.Sql(), "SELECT u.id AS u_id, u.name AS u_name, uc.user_id AS uc_user_id, uc.no AS uc_no FROM user u INNER JOIN user_card uc ON u.id=uc.user_id WHERE u.name='toukii'")
}

func TestParse(t *testing.T) {
	u, err := Parse("test_yaml/user.yaml")
	if err != nil {
		t.Error(err)
	}

	uc, err := Parse("test_yaml/user_card.yaml")
	if err != nil {
		t.Error(err)
	}

	u.Select(nil)
	equal.Equal(nil, u.Sql(), "SELECT u.id AS u_id, u.name AS u_name FROM user u")

	u = NewT("user", "u", []string{"id", "name"})
	u.Mask["Mask"] = []string{"name"}

	u.Select(nil)
	equal.Equal(nil, u.Sql(), "SELECT u.id AS u_id, u.name AS u_name FROM user u")

	m := NewMask("Mask", u, uc)
	u.Select(m).InnerJoin(uc, "u.id=uc.user_id")
	equal.Equal(nil, u.Sql(), "SELECT u.name AS u_name, uc.no AS uc_no FROM user u INNER JOIN user_card uc ON u.id=uc.user_id")
}

func TestMultiParse(t *testing.T) {
	mt, err := MultiParse("test_yaml/model.yaml")
	if err != nil {
		t.Error(err)
	}

	u := mt["User"]
	uc := mt["UserCard"]
	b := mt["Bank"]

	m := NewMask("Mask", u, uc, b)
	u.Select(m).InnerJoin(uc, "u.id=uc.user_id").InnerJoin(b, "uc.bank_id=b.id")
	equal.Equal(nil, u.Sql(), "SELECT u.name AS u_name, uc.no AS uc_no, b.name AS b_name FROM user u INNER JOIN user_card uc ON u.id=uc.user_id INNER JOIN bank b ON uc.bank_id=b.id")
}
