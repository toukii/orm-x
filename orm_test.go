package orm

import (
	"fmt"
	"testing"

	"github.com/didi/gendry/scanner"
	"github.com/toukii/assert"
)

// = - = - = - = - = - = - = - = -D=E-F=I-N=E- = - = - = - = - = - = - = - =
type UserCardMask struct {
	Name     string `db:"u_name"`
	CardNo   string `db:"no"`
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
	mt, err := MultiParse("test_yaml/model.yaml")
	if err != nil {
		t.Error(err)
	}

	u := mt["User"]
	uc := mt["UserCard"]
	b := mt["Bank"]

	m := NewMask("Mask", u, uc, b)
	u.Select(m).InnerJoin(uc, "u.id=uc.user_id").InnerJoin(b, "uc.bank_id=b.id")
	fmt.Println(u.Sql())

	rows, err := Mysql().Query(u.Sql())
	equal.Equal(nil, err, nil)

	var mask []*UserCardMask
	err = scanner.Scan(rows, &mask)
	equal.Equal(nil, err, nil)
	fmt.Printf("%+v\n", mask)
	for _, m := range mask {
		fmt.Printf("%+v\n", m)
	}
}

func TestExample(t *testing.T) {
	ExampleScan()
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

	t1.Select(mask).InnerJoin(t2, "u.id=uc.user_id")
	equal.Equal(nil, t1.Sql(), "SELECT u.id AS u_id, u.name AS u_name, uc.user_id AS uc_user_id, uc.no AS uc_no FROM user u INNER JOIN user_card uc ON u.id=uc.user_id")
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
