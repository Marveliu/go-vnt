package supervisor

import (
	"html/template"
	"io/ioutil"
	"os"
)

const (
	PATH           = "templates/"
	BIZCONTRACT_TP = "bizContract"
)

func Gen(src string, out string) {

	// load templates
	t1, err := template.ParseFiles(PATH + BIZCONTRACT_TP)
	if err != nil {
		panic(err)
	}

	// parse bizMeta
	dat, err := ioutil.ReadFile(src)
	bizMeta := BizMeta{}
	if _, err := bizMeta.Decode(dat); err != nil {
		panic(err)
	}

	// render
	f, err := os.Create(out)
	if t1.Execute(f, bizMeta) != nil {
		panic("Gen failed !")
	}
}

func Reg_Biz(meta BizMeta) string {
	// reg to supervisor contract

	return ""
}
