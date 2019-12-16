package supervisor

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
)

const (
	PATH           = "templates/"
	BIZCONTRACT_TP = "bizContract"
)

func GenFile(src string, out string) {

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

func Gen(cfg string) string {

	// load templates
	t1, err := template.ParseFiles(PATH + BIZCONTRACT_TP)
	if err != nil {
		panic(err)
	}

	// parse bizMeta
	bizMeta := BizMeta{}
	if _, err := bizMeta.Decode([]byte(cfg)); err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	// render
	if t1.Execute(buf, bizMeta) != nil {
		panic("Gen failed !")
	}
	return buf.String()
}
