package supervisor

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/state"
	"github.com/vntchain/go-vnt/core/vm/interface"
	"github.com/vntchain/go-vnt/vntdb"
)

var url = []byte("/ip4/127.0.0.1/tcp/30303/ipfs/1kHGq5zZFRW5FBJ9YMbbvSiW4AzGg5CKMCtDeg6FNnjCbGS")

type testContext struct {
	Origin      common.Address
	Time        *big.Int
	StateDB     inter.StateDB
	BlockNumber *big.Int
}

func (tc *testContext) GetOrigin() common.Address {
	return tc.Origin
}

func (tc *testContext) GetStateDb() inter.StateDB {
	return tc.StateDB
}

func (tc *testContext) GetTime() *big.Int {
	return tc.Time
}

func (tc *testContext) SetTime(t *big.Int) {
	tc.Time = t
}

func (tc *testContext) GetBlockNum() *big.Int {
	return tc.BlockNumber
}

func newcontext() inter.ChainContext {
	db := vntdb.NewMemDatabase()
	stateDB, _ := state.New(common.Hash{}, state.NewDatabase(db))
	c := testContext{
		Origin:      common.BytesToAddress([]byte{111}),
		Time:        big.NewInt(1531328510),
		StateDB:     stateDB,
		BlockNumber: big.NewInt(1),
	}
	return &c
}

func TestSetSupervisor(t *testing.T) {
	context := newcontext()
	c := newSupervisorContext(context)
	bizContract := BizContract{
		Owner:     common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"),
		Address:   common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"),
		TimeStamp: big.NewInt(1531454152),
		Status:    big.NewInt(0),
		BizType:   big.NewInt(0),
		Name:      []byte("Hello"),
		Desc:      []byte("Hello World"),
	}
	if err := c.RegisterBizContract(bizContract); err != nil {
		t.Errorf("addr: %s, error: %s", bizContract.Owner, err)
	}

	ret, _ := c.GetBizContract(common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"))
	if !reflect.DeepEqual(bizContract, ret) {
		t.Errorf("not equal")
	}
}

func Test_supervisorContext_RegBizMetaData(t *testing.T) {

}
