package supervisor

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/vm/interface"
	"github.com/vntchain/go-vnt/log"
	"math/big"
)

const (
	ContractAddr = "0x0000000000000000000000000000000000000008"
	BizMetaKey   = "BizMetas"
)

var (
	ErrDescInvalid    = errors.New("the length of contract's description should between [10, 200]")
	ErrNotExisted     = errors.New("the target is not existed")
	ErrBizMetaExisted = errors.New("the target is  existed")
)

var (
	contractAddr = common.HexToAddress(ContractAddr)
	bizMetaKey   = common.BytesToHash([]byte(BizMetaKey))
	emptyAddress = common.Address{}
)

type Supervisor struct{}

func (s *Supervisor) RequiredGas(input []byte) uint64 {
	return 0
}

func (s *Supervisor) Run(ctx inter.ChainContext, input []byte, value *big.Int) ([]byte, error) {

	nonce := ctx.GetStateDb().GetNonce(contractAddr)
	ctx.GetStateDb().SetNonce(contractAddr, nonce+1)
	supervisorAbi, err := GetSuervisorABI()
	if err != nil {
		return nil, err
	}

	if len(input) < 4 {
		return nil, nil
	}

	// input的组成见abi.Pack函数
	methodId := input[:4]
	methodArgs := input[4:]

	methodName := "None"
	isMethod := func(name string) bool {
		if bytes.Equal(methodId, supervisorAbi.Methods[name].Id()) {
			methodName = name
			return true
		}
		return false
	}

	c := newSupervisorContext(ctx)
	sender := ctx.GetOrigin()

	var (
		ret []byte
	)

	ret = nil
	switch {
	case isMethod("RegisterBizContract"):
		var bizContract BizContract
		if err = supervisorAbi.UnpackInput(&bizContract, methodName, methodArgs); err == nil {
			bizContract.Address = sender
			c.RegisterBizContract(bizContract)
		}
	// case isMethod("updateConfig"):
	case isMethod("Report"):
		var data ReportReq
		if err = supervisorAbi.UnpackInput(&data, methodName, methodArgs); err == nil {
			log.Info(string(data.Msg))
		}
	case isMethod("RegBizMeta"):
		var data BizMetaReq
		if err = supervisorAbi.UnpackInput(&data, methodName, methodArgs); err == nil {
			err = c.RegBizMeta(data)
			// topics := make([]common.Hash, 0)
			// ctx.GetStateDb().AddLog(&types.Log{
			// 	Address:     contractAddr,
			// 	Topics:      topics,
			// 	Data:        abi.U256(big.NewInt(int64(n))),
			// 	BlockNumber: ctx.GetBlockNum().Uint64(),
			// })
		}
	case isMethod("GetBizMetaTemplate"):
		var data BizContractTpReq
		if err = supervisorAbi.UnpackInput(&data, methodName, methodArgs); err == nil {
			meta := c.getBizMeta(data.BizType)
			if meta != nil {
				ret, err = PackOutPut(supervisorAbi, methodName, GenFromBizMeta(*meta))
			}
		}
	}

	if err != nil {
		log.Error("call supervisor contract err:", "method", methodName, "err", err)
	} else {
		log.Debug("Supervisor call", "method", methodName)
	}

	return ret, err
}

type supervisorContext struct {
	context inter.ChainContext
}

func newSupervisorContext(ctx inter.ChainContext) supervisorContext {
	return supervisorContext{
		context: ctx,
	}
}

func (sc supervisorContext) RegisterBizContract(contract BizContract) error {
	return sc.setObject(PREFIX_BIZCONTRACT, contract.Owner, contract)
}

func (sc supervisorContext) GetBizContract(addr common.Address) (BizContract, error) {
	ret := BizContract{}
	if sc.getObject(PREFIX_BIZCONTRACT, addr, &ret) != nil {
		return ret, ErrNotExisted
	}
	return ret, nil
}

func (sc supervisorContext) RegBizMeta(data BizMetaReq) error {
	var meta BizMeta
	if _, err := meta.Decode(common.FromHex(data.Cfg)); err != nil {
		return err
	}
	// generate No and store
	if sc.getBizMeta(meta.No) != nil {
		return ErrBizMetaExisted
	}
	return sc.setObject(PREFIX_BIZMETA, common.BigToAddress(big.NewInt(int64(meta.No))), meta)
}

func (sc supervisorContext) GetBizMeta(n uint32) *BizMeta {
	// todo check
	return sc.getBizMeta(n)
}

func (sc supervisorContext) UpdateConfig(str string) error {
	config := &Config{}
	if err := json.Unmarshal([]byte(str), &config); err != nil {
		log.Error("Parse supervisor config error ", str)
		return err
	}
	return nil
}

// func (sc supervisorContext) GetConfig() (Config, error) {
//
// }

type BizContract struct {
	Address   common.Address // 合约地址
	Owner     common.Address // 所有者地址
	Name      []byte         // 合约名称
	Desc      []byte         // 描述
	BizType   *big.Int       // 合约类型
	Status    *big.Int       // 状态
	TimeStamp *big.Int       // 创建时间
}

type MngNode struct {
	Id     *big.Int // 编号
	Name   []byte   // 节点名称
	Desc   []byte   // 节点描述
	Ip     []byte   // 地址
	Status *big.Int // 状态
	Pubkey []byte   // 公钥
}

type ReportConfig struct {
	Id      *big.Int // 编号
	BizType *big.Int // 类型
	Status  *big.Int // 状态
	Config  []byte   // 配置
	Version *big.Int // 版本
}

type Config struct {
	AccountBlackLists []byte                  // 账户黑名单
	MngNodes          map[string]MngNode      // 管理员账户
	ReportConfig      map[string]ReportConfig // 配置文件格式
}

type ReportReq struct {
	Msg []byte
}

type BizMetaReq struct {
	Cfg string
}

type BizContractTpReq struct {
	BizType uint32
}
