package supervisor

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/vntchain/go-vnt/common"
)

var (
	bizContract = BizContract{
		Owner:     common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"),
		Address:   common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"),
		TimeStamp: big.NewInt(1531454152),
		Status:    big.NewInt(0),
		BizType:   big.NewInt(0),
		Name:      []byte("Hello"),
		Desc:      []byte("Hello World"),
	}
)

func TestObjDataLayer(t *testing.T) {
	context := newcontext()
	c := newSupervisorContext(context)
	if c.setObject(PREFIX_BIZCONTRACT, bizContract.Owner, bizContract) != nil {
		t.Errorf("test: failed, setVal: %s", bizContract)
	}
	nb := BizContract{}
	if c.getObject(PREFIX_BIZCONTRACT, bizContract.Owner, &nb) != nil {
		t.Errorf("test: failed, getVal: %s", bizContract)
	}
	if !reflect.DeepEqual(bizContract, nb) {
		t.Errorf("test: failed, want: %s, get: %s", bizContract, nb)
	}
}

func TestStringDataLayer(t *testing.T) {
	tests := []string{
		"",
		"1",
		"9ee97d274eb4c215f23238fee1f103d9ea10a2341",
		"9ee97d274eb4c215f23238fee1f103d9ea10a234",
		"{\"book\":[{\"id\":\"01\",\"language\":\"Java\",\"edition\":\"third\",\"author\":\"Herbert Schildt\"}]}",
	}
	context := newcontext()
	c := newSupervisorContext(context)
	for i, tt := range tests {
		key := common.BigToHash(big.NewInt(int64(i)))
		if c.setStringToDB(key, tt) != nil {
			t.Errorf("test: %d failed, setVal: %s", i, tt)
		}
		if val, err := c.getStringFromDB(key); err != nil || strings.Compare(tt, val) != 0 {
			t.Errorf("test: %d failed, want: %s, get: %s", i, tt, val)
		}
	}
}

func Test_supervisorContext_getObjKey(t *testing.T) {
	tests := []string{
		"",
		"1",
		"9ee97d274eb4c215f23238fee1f103d9ea10a2341",
		"9ee97d274eb4c215f23238fee1f103d9ea10a234",
	}
	context := newcontext()
	c := newSupervisorContext(context)
	for _, tt := range tests {
		src := common.HexToAddress(tt)
		fmt.Println(src)
		str := c.getObjKey(PREFIX_CONFIG, src)
		fmt.Println(str)
	}
}
