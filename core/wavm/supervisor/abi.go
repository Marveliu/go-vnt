package supervisor

import (
	"github.com/vntchain/go-vnt/accounts/abi"
	"strings"
)

const SupervisorAbiJSON = `[
{"name":"RegisterBizContract","constant":false,"inputs":[{"name":"addr","type":"address","indexed":false},{"name":"name","type":"string","indexed":false},{"name":"desc","type":"string","indexed":false},{"name":"type","type":"int32","indexed":false},{"name":"owner","type":"string","indexed":false}],"outputs":[{"name":"output","type":"bool","indexed":false}],"type":"function"},
{"name":"ReportData","constant":false,"inputs":[{"name":"msg","type":"string","indexed":false}],"outputs":[{"name":"output","type":"bool","indexed":false}],"type":"function"}
]`

func GetSuervisorABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(SupervisorAbiJSON))
}

func PackInput(abiobj abi.ABI, name string, args ...interface{}) ([]byte, error) {
	abires := abiobj
	var res []byte
	var err error
	if len(args) == 0 {
		res, err = abires.Pack(name)
	} else {
		res, err = abires.Pack(name, args)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}
