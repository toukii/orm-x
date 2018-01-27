package orm

import (
	"github.com/toukii/goutils"
	yaml "gopkg.in/yaml.v2"
)

func Parse(filename string) (*T, error) {
	bs := goutils.ReadFile(filename)
	t := new(T)
	err := yaml.Unmarshal(bs, t)
	t.init()
	return t, err
}

func MultiParse(filename string) (map[string]*T, error) {
	bs := goutils.ReadFile(filename)
	mt := make(map[string]*T)
	err := yaml.Unmarshal(bs, &mt)
	for _, t := range mt {
		t.init()
	}
	return mt, err
}
