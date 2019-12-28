package supervisor

import (
	"encoding/json"
	"fmt"
	"github.com/vntchain/go-vnt/accounts/keystore"
	"io/ioutil"
	"math/big"
	"path/filepath"
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
	// context := newcontext()
	// c := newSupervisorContext(context)
	// bizContract := BizContract{
	// 	Owner:     common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"),
	// 	Address:   common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"),
	// 	TimeStamp: big.NewInt(1531454152),
	// 	Status:    0,
	// 	BizType:   123,
	// 	Name:      "Hello",
	// 	Desc:      "Hello World",
	// }
	// if err := c.RegisterBizContract(bizContract); err != nil {
	// 	t.Errorf("addr: %s, error: %s", bizContract.Owner, err)
	// }
	// ret, _ := c.GetBizContract(common.HexToAddress("9ee97d274eb4c215f23238fee1f103d9ea10a234"))
	// if !reflect.DeepEqual(bizContract, ret) {
	// 	t.Errorf("not equal")
	// }
}

func Test_supervisorContext_RegBizMetaData(t *testing.T) {
	data := "0x62697a4e616d65203d20227472616e73616374696f6e220a62697a54797065203d202231220a64657363203d2022e5819ae4b8aae4baa4e69893220a76657273696f6e203d20310a0a5b5b64617461735d5d0a20206e616d65203d202270726f64756374220a20207265706f7274203d2022312c32220a0a20205b5b64617461732e6669656c64735d5d0a202020206e616d65203d2022706964220a2020202074797065203d2022737472696e67220a0a20205b5b64617461732e6669656c64735d5d0a202020206e616d65203d202264657363220a2020202074797065203d2022737472696e67220a0a5b5b64617461735d5d0a20206e616d65203d20226f72646572220a20207265706f7274203d202231220a0a20205b5b64617461732e6669656c64735d5d0a202020206e616d65203d202246726f6d220a2020202074797065203d2022737472696e67220a0a20205b5b64617461732e6669656c64735d5d0a202020206e616d65203d2022546f220a2020202074797065203d2022737472696e67220a0a20205b5b64617461732e6669656c64735d5d0a202020206e616d65203d2022706964220a0974797065203d2022737472696e67220a0a20205b5b64617461732e6669656c64735d5d0a202020206e616d65203d202276616c7565220a0974797065203d202275696e743332220a0a5b5b7461736b735d5d0a20206e616d65203d2022e58f91e5b883220a202064657363203d2022e58f91e5b883e59586e59381220a0a20205b5b7461736b732e616374696f6e735d5d0a2020202066756e634e616d65203d20227075626c697368220a202020206d757461626c65203d20747275650a2020202064617461526566203d202270726f64756374220a0a5b5b7461736b735d5d0a20206e616d65203d2022e4baa4e69893220a202064657363203d2022e4baa4e69893e59586e59381220a0a20205b5b7461736b732e616374696f6e735d5d0a2020202066756e634e616d65203d20227472616e73616374696f6e220a202020206d757461626c65203d20747275650a2020202064617461526566203d20226f72646572220a"
	context := newcontext()
	c := newSupervisorContext(context)
	c.RegBizMeta(BizMetaReq{data})
	meta := c.getBizMeta(123)
	fmt.Println(meta)
}

func Test_fromhex(t *testing.T) {
	data := "0x63be"
	bs := common.FromHex(data)
	fmt.Println(string(bs))
}

func Test_supervisorContext_UpdateConfig(t *testing.T) {
	// context := newcontext()
	// c := newSupervisorContext(context)
	workdir, err := ioutil.TempDir("", "")
	ks := keystore.NewKeyStore(filepath.Join(workdir, "keystore"), keystore.StandardScryptN, keystore.StandardScryptP)
	newAcc, err := ks.NewAccount("Creation password")
	if err != nil {
		t.Fatalf("Failed to create new account: %v", err)
	}
	fmt.Println(newAcc)
}

func Test_Config(t *testing.T) {
	cfg := Config{
		AccountBlackLists: []common.Address{
			common.HexToAddress("0x13f6c5ba187f165e64c651c0e1167091d067089c"),
			common.HexToAddress("0x23f6c5ba187f165e64c651c0e1167091d067089c")},
		MngNodes: map[string]MngNode{
			"123": {
				Id:     1,
				Name:   "监查院",
				Status: 1,
				Ip:     "127.0.0.1:9999",
				Pubkey: "12334",
				Addr:   common.HexToAddress("0xf3f6c5ba187f165e64c651c0e1167091d067089c"),
				Desc:   "陈萍萍",
			},
		},
	}
	b, _ := json.Marshal(cfg)
	n := &Config{}
	str := "{\"AccountBlackLists\":[\"0x13f6c5ba187f165e64c651c0e1167091d067089c\",\"0x23f6c5ba187f165e64c651c0e1167091d067089c\"],\"MngNodes\":{\"123\":{\"Id\":1,\"Name\":\"监查院\",\"Desc\":\"陈萍萍\",\"Ip\":\"127.0.0.1:9999\",\"Status\":1,\"Addr\":\"0xf3f6c5ba187f165e64c651c0e1167091d067089c\",\"Pubkey\":\"12334\"}}}"
	json.Unmarshal([]byte(str), n)
	fmt.Println(string(b))
	if !reflect.DeepEqual(*n, cfg) {
		t.Errorf("failed")
	}
}

func Test_ReportFiled(t *testing.T) {
	data := make(map[string]ReportField)
	data["0"] = ReportField{
		1,
		"123213",
	}
	data["1"] = ReportField{
		2,
		1,
	}
	data["2"] = ReportField{
		3,
		true,
	}
	bs, _ := json.Marshal(data)
	d := make(map[string]ReportField)
	if err := json.Unmarshal(bs, &d); err != nil {
		t.Errorf("failed")
	} else {
		t.Log("success")
	}

}
