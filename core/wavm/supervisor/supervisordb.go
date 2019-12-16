package supervisor

import (
	"encoding/binary"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/rlp"
	"math/big"
	"reflect"
)

const (
	PREFIX_CONFIG      = byte(0)
	PREFIX_BIZCONTRACT = byte(1)
	PREFIX_BIZMETA     = byte(2)

	PREFIXLENGTH = 4 // key的结构为，4位表前缀，20位address，8位的value在struct中的位置
)

var KeyNotExistErr = errors.New("the key do not exist")

func (sc supervisorContext) getBizMeta(n int) BizMeta {
	meta := BizMeta{}
	sc.getObject(PREFIX_BIZMETA, common.BigToAddress(big.NewInt(int64(n))), meta)
	return meta
}

func (sc supervisorContext) getObject(prefix byte, key common.Address, v interface{}) error {
	str, err := sc.getStringFromDB(sc.getObjKey(prefix, key))
	if err != nil {
		log.Error("Key not existed.", key)
		return nil
	}
	return json.Unmarshal([]byte(str), v)
}

func (sc supervisorContext) setObject(prefix byte, key common.Address, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Error("Marshal supervisor fail.", v)
		return nil
	}
	if err := sc.setStringToDB(sc.getObjKey(prefix, key), string(bytes)); err != nil {
		log.Error("setObject error", "err", err, "type", reflect.ValueOf(v).Kind())
	}
	return nil
}

func (sc supervisorContext) getFromDB(key common.Hash) common.Hash {
	return sc.context.GetStateDb().GetState(contractAddr, key)
}

func (sc supervisorContext) setToDB(key common.Hash, value common.Hash) {
	sc.context.GetStateDb().SetState(contractAddr, key, value)
}

func (sc supervisorContext) getStringFromDB(key common.Hash) (string, error) {
	valByte := sc.getFromDB(key)
	if valByte == (common.Hash{}) {
		return "", KeyNotExistErr
	}

	// 部分byte数组过长，是拆分了之后存储的
	var val []byte
	err := rlp.DecodeBytes(valByte.Big().Bytes(), &val)
	if err == nil {
		return string(val), nil
	} else {
		val = valByte.Big().Bytes()
		var tmp []byte
		for j := 1; ; j++ {
			binary.BigEndian.PutUint32(key[PREFIXLENGTH+common.AddressLength:], uint32(j))
			arrayByte := sc.getFromDB(key)
			if arrayByte.Big().Sign() == 0 {
				return "", KeyNotExistErr
			}
			val = append(val, arrayByte.Bytes()...)
			if err = rlp.DecodeBytes(val, &tmp); err == nil {
				return string(tmp), nil
			}
		}
	}
}

func (sc supervisorContext) setStringToDB(key common.Hash, value string) error {
	elem, err := rlp.EncodeToBytes(value)
	if err != nil {
		return err
	}
	// 如果要存储的字节过长，就拆分了存
	// 0号位置存储切分的长度，后面按右对齐方式存储，若需要补空位，补在第一个元素处
	valLen := len(elem)/32 + 1
	var j int
	for j = valLen - 1; j >= 0; j-- {
		var subKey common.Hash
		copy(subKey[:], key[:])
		binary.BigEndian.PutUint32(subKey[PREFIXLENGTH+common.AddressLength:], uint32(j))
		cutPos := len(elem) - 32
		if cutPos < 0 {
			sc.setToDB(subKey, common.BytesToHash(elem))
			break
		}
		tmpElem := elem[cutPos:]
		elem = elem[:cutPos]
		sc.setToDB(subKey, common.BytesToHash(tmpElem))
	}
	return nil
}

func (sc supervisorContext) getObjKey(bizType byte, key common.Address) common.Hash {
	res := common.Hash{}
	res[0] = bizType
	copy(res[PREFIXLENGTH:], key.Bytes())
	return res
}
