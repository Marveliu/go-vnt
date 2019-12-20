package supervisor

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
)

const (
	PATH           = "/Users/mac/gopath/src/github.com/vntchain/go-vnt/core/vm/supervisor/templates/"
	BIZCONTRACT_TP = "bizContract"
)

func GenFile(src string, out string) {

	// parse bizMeta
	dat, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}
	ret := Gen(dat)

	// render
	f, err := os.Create(out)
	if _, err := f.Write(ret); err != nil {
		panic("Gen failed !")
	}
}

func Gen(cfg []byte) []byte {

	// load templates
	t1, err := template.ParseFiles(PATH + BIZCONTRACT_TP)
	if err != nil {
		panic(err)
	}

	// parse bizMeta
	bizMeta := BizMeta{}
	if _, err := bizMeta.Decode(cfg); err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	// render
	if t1.Execute(buf, bizMeta) != nil {
		panic("Gen failed !")
	}
	return buf.Bytes()
}

func GenFromBizMeta(meta BizMeta) []byte {
	t1, err := template.ParseFiles(PATH + BIZCONTRACT_TP)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	// render
	if t1.Execute(buf, meta) != nil {
		panic("Gen failed !")
	}
	return buf.Bytes()
}
