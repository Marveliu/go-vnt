package supervisor

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestBizMetaRW(t *testing.T) {
	path := "tests/BizMeta.toml"
	task := BizMeta{
		BizName: "交易",
		BizType: "1",
		Desc:    "做个交易",
		Version: 1,
		Datas: []Data{
			{"Obj1", "1,2", []Field{
				{"From", "string"},
				{"To", "string"},
			}},
			{"Obj2", "1", []Field{{"To", "string"}}},
		},
		Tasks: []Task{
			{"发布", "发布商品", []Action{
				Action{"publish", true, "Obj1"},
			}},
			{"交易", "交易商品", []Action{
				Action{"transaction", false, "Obj2"},
			}},
		},
	}
	f, err := os.Create(path)
	check(err)
	defer f.Close()
	buffer, err := task.ToTOML()
	check(err)
	if _, err := f.Write(buffer.Bytes()); err != nil {
		panic(err)
	}
	dat, err := ioutil.ReadFile(path)
	check(err)
	n := BizMeta{}
	n.Decode(dat)
	fmt.Println(n.Valid())
	if !reflect.DeepEqual(n, task) {
		t.Errorf("读取不一致")
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}