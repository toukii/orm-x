package orm

import (
	"bytes"
	"fmt"
	"strings"
)

func NewT(name, alias string, cols []string) *T {
	t := &T{
		Name:    name,
		Alias:   alias,
		Columns: cols,
		Mask:    make(map[string][]string, 1),
	}
	t.init()
	return t
}

type T struct {
	Name    string              `yaml:"name"`
	Alias   string              `yaml:"alias"`
	Columns []string            `yaml:"columns"`
	Mask    map[string][]string `yaml:"mask"`

	_sql *bytes.Buffer
}

func (t *T) Sql() string {
	return t._sql.String()
}

func (t *T) init() {
	t.Mask[""] = t.Columns
	t._sql = bytes.NewBuffer(nil)
}

func (t *T) aliasTable() string {
	return t.Name + " " + t.Alias
}

func (t *T) aliasColumns() []string {
	return prefixAlias(t.Alias, t.Columns)
}

func (t *T) InnerJoin(t2 *T, on string) *T {
	ret := fmt.Sprintf(" INNER JOIN %s ON %s", t2.aliasTable(), on)
	t._sql.WriteString(ret)
	return t
}

func (t *T) LeftJoin(t2 *T, on string) *T {
	ret := fmt.Sprintf(" LEFT JOIN %s ON %s", t2.aliasTable(), on)
	t._sql.WriteString(ret)
	return t
}

func (t *T) Select(mask *Mask) *T {
	scol := ""
	if mask == nil {
		scol = strings.Join(t.aliasColumns(), ", ")
	} else {
		scol = strings.Join(mask.cols, ", ")
	}
	ret := fmt.Sprintf("SELECT %s FROM %s", scol, t.aliasTable())
	t._sql.Reset()
	t._sql.WriteString(ret)
	return t
}

func (t *T) mask(mask string) []string {
	return prefixAlias(t.Alias, t.Mask[mask])
}

func prefixAlias(alias string, cols []string) []string {
	ret := make([]string, len(cols))
	for i, it := range cols {
		ret[i] = fmt.Sprintf("%s.%s AS %s_%s", alias, it, alias, it)
	}
	return ret
}

// ====================
type Mask struct {
	ts   []*T
	cols []string
}

func (m *Mask) init(mask string) {
	cols := make([]string, 0, 10)
	for _, t := range m.ts {
		cols = append(cols, t.mask(mask)...)
	}
	m.cols = cols
}

func NewMask(mask string, ts ...*T) *Mask {
	m := &Mask{
		ts: ts,
	}
	m.init(mask)
	return m
}
