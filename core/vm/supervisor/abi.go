package supervisor

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/common/math"
	"math/big"
	"reflect"
	"strings"
)

const SupervisorAbiJSON = `[
{"name":"Supervisor","constant":false,"inputs":[],"outputs":[],"type":"constructor"},
{"name":"RegisterBizContract","constant":false,"inputs":[{"name":"addr","type":"address","indexed":false},{"name":"owner","type":"address","indexed":false},{"name":"name","type":"string","indexed":false},{"name":"bizType","type":"uint32","indexed":false},{"name":"info","type":"string","indexed":false}],"outputs":[{"name":"output","type":"bool","indexed":false}],"type":"function"},
{"name":"Report","constant":true,"inputs":[{"name":"msg","type":"string","indexed":false}],"outputs":[],"type":"function"},
{"name":"RegBizMeta","constant":false,"inputs":[{"name":"cfg","type":"string","indexed":false}],"outputs":[{"name":"output","type":"uint32","indexed":false}],"type":"function"},
{"name":"GetBizMetaTemplate","constant":true,"inputs":[{"name":"bizType","type":"uint32","indexed":false}],"outputs":[{"name":"output","type":"string","indexed":false}],"type":"function"}
]`

var endianess = binary.LittleEndian

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

func PackOutPut(abiobj abi.ABI, name string, ret []byte) ([]byte, error) {
	var res []byte
	var err error
	if err != nil {
		return nil, err
	}
	if val, ok := abiobj.Methods[name]; ok {
		outputs := val.Outputs
		if len(outputs) != 0 {
			output := outputs[0].Type.T
			switch output {
			case abi.StringTy:
				l, err := packNum(reflect.ValueOf(32))
				if err != nil {
					return nil, err
				}
				s, err := packBytesSlice(ret, len(ret))
				if err != nil {
					return nil, err
				}
				return append(l, s...), nil
			case abi.UintTy, abi.IntTy:
				bigint := getU256(ret)
				return abi.U256(bigint), nil
			case abi.BoolTy:
				if ret != nil {
					return math.PaddedBigBytes(common.Big1, 32), nil
				}
				return math.PaddedBigBytes(common.Big0, 32), nil
			case abi.AddressTy:
				return common.LeftPadBytes(ret, 32), nil
			default:
				return nil, errors.New("未知返回类型")
			}
		} else { // 无返回类型
			return I32ToBytes(0), nil
		}
	}
	return res, nil
}

func packNum(value reflect.Value) ([]byte, error) {
	switch kind := value.Kind(); kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return abi.U256(new(big.Int).SetUint64(value.Uint())), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return abi.U256(big.NewInt(value.Int())), nil
	case reflect.Ptr:
		return abi.U256(value.Interface().(*big.Int)), nil
	default:
		return nil, fmt.Errorf("abi: fatal error")
	}
}

func packBytesSlice(bytes []byte, l int) ([]byte, error) {
	len, err := packNum(reflect.ValueOf(l))
	return append(len, common.RightPadBytes(bytes, (l+31)/32*32)...), err
}

func getU256(mem []byte) *big.Int {
	bigint := new(big.Int)
	var toStr string
	if len(mem) == 0 {
		toStr = "0"
	} else {
		toStr = string(mem)
	}
	_, success := bigint.SetString(toStr, 10)
	if success == false {
		panic("Illegal uint256 input " + toStr)
	}
	return math.U256(bigint)
}

func I32ToBytes(i32 uint32) []byte {
	bytes := make([]byte, 4)
	endianess.PutUint32(bytes, i32)
	return bytes
}
