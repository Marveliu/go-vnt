package supervisor

import (
	"bytes"
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/vntchain/go-vnt/log"
)

var (
	ErrDataNotExisted     = errors.New("data definition must be declared before action")
	ErrDuplicatedDeclared = errors.New("definition duplication")
)

type Action struct {
	FuncName string `toml:"funcName"`
	Mutable  bool   `toml:"mutable"`
	DataRef  string `toml:"dataRef"`
}

type Task struct {
	Name    string   `toml:"name"`
	Desc    string   `toml:"desc"`
	Actions []Action `toml:"actions"`
}

type Field struct {
	Name string `toml:"name"`
	Type string `toml:"type"`
}

type Data struct {
	Name   string  `toml:"name"`
	Report string  `toml:"report"`
	Fields []Field `toml:"fields"`
}

type BizMeta struct {
	no      int    // 编号
	BizName string `toml:"bizName"`
	BizType string `toml:"bizType"`
	Desc    string `toml:"desc"`
	Version int    `toml:"version"`
	Datas   []Data `toml:"datas"`
	Tasks   []Task `toml:"tasks"`
}

func (t *BizMeta) ToTOML() (*bytes.Buffer, error) {
	b := &bytes.Buffer{}
	encoder := toml.NewEncoder(b)
	if err := encoder.Encode(t); err != nil {
		return nil, err
	}
	return b, nil
}

func (t *BizMeta) Decode(data []byte) (toml.MetaData, error) {
	return toml.Decode(string(data), t)
}

func (t *BizMeta) Valid() error {

	dataTable := map[string]Data{}
	for _, data := range t.Datas {
		if _, ok := dataTable[data.Name]; ok {
			return ErrDuplicatedDeclared
		}
		dataTable[data.Name] = data
	}
	for _, task := range t.Tasks {
		for _, action := range task.Actions {
			if _, ok := dataTable[action.DataRef]; !ok {
				log.Error("Parase BizConfig failed: cannot find data definition", action.DataRef)
				return ErrDataNotExisted
			}
		}
	}

	// TODO 参数配置格式解析
	return nil
}
