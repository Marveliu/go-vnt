package supervisor

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/core/vm/interface"
	"github.com/vntchain/go-vnt/log"
	"math/big"
)

const (
	ContractAddr = "0x0000000000000000000000000000000000000008"
	BizMetaKey   = "BizMetas"
	ConfigKey    = "Config"
)

var (
	ErrDescInvalid        = errors.New("the length of contract's description should between [10, 200]")
	ErrNotExisted         = errors.New("the target is not existed")
	ErrBizMetaExisted     = errors.New("the target is  existed")
	ErrLackField          = errors.New("ReportData invalid")
	ErrReportDataNotFound = errors.New("ReportData Not Found")
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

	var (
		ret []byte
	)

	ret = nil
	switch {
	case isMethod("RegisterBizContract"):
		var req RegContractReq
		if err = supervisorAbi.UnpackInput(&req, methodName, methodArgs); err == nil {
			err = c.RegisterBizContract(req)
		}
	case isMethod("GetBizContract"):
		var req GetBizContractReq
		if err = supervisorAbi.UnpackInput(&req, methodName, methodArgs); err == nil {
			bizContract := c.GetBizContract(req.Addr)
			if bizContract != nil {
				bs, _ := json.Marshal(bizContract)
				ret, err = PackOutPut(supervisorAbi, methodName, bs)
			}
		}
	case isMethod("UpdateConfig"):
		var req UpdateConfigReq
		if err = supervisorAbi.UnpackInput(&req, methodName, methodArgs); err == nil {
			err = c.UpdateConfig(req.Cfg)
		}
	case isMethod("GetConfig"):
		cfg := c.GetConfig()
		if cfg != nil {
			bs, _ := json.Marshal(cfg)
			ret, err = PackOutPut(supervisorAbi, methodName, bs)
		}
	case isMethod("Report"):
		var data ReportReq
		if err = supervisorAbi.UnpackInput(&data, methodName, methodArgs); err == nil {
			// TODO 校验
			// 定位 Contract -> Meta + Data "reportDataName"
			err = c.RecordData(data.Addr, data.DataName, data.Msg)
		}
	case isMethod("RegBizMeta"):
		var data BizMetaReq
		if err = supervisorAbi.UnpackInput(&data, methodName, methodArgs); err == nil {
			err = c.RegBizMeta(data)
		}
	case isMethod("GetBizMetaTemplate"):
		var data BizContractTpReq
		if err = supervisorAbi.UnpackInput(&data, methodName, methodArgs); err == nil {
			meta := c.getBizMeta(data.BizType)
			if meta != nil {
				ret, err = PackOutPut(supervisorAbi, methodName, GenFromBizMeta(*meta))
			}
		}
	case isMethod("ValidateBizContract"):
		var req GetBizContractReq
		if err = supervisorAbi.UnpackInput(&req, methodName, methodArgs); err == nil {
			// TODO 优化布尔类型值的转化
			if c.ValidateBizContract(req.Addr) {
				ret, err = PackOutPut(supervisorAbi, methodName, []byte("1"))
			} else {
				ret, err = PackOutPut(supervisorAbi, methodName, nil)
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

func (sc supervisorContext) RegisterBizContract(req RegContractReq) error {
	// TODO req 检查
	meta := sc.getBizMeta(req.BizType)
	if meta == nil {
		return ErrNotExisted
	}
	contract := BizContract{
		Address:   req.Addr,
		Owner:     req.Owner,
		BizType:   req.BizType,
		Desc:      meta.Desc,
		Name:      meta.BizName,
		Status:    0,
		TimeStamp: sc.context.GetTime(),
	}

	return sc.setObject(PREFIX_BIZCONTRACT, contract.Address, contract)
}

func (sc supervisorContext) GetBizContract(addr common.Address) *BizContract {
	ret := &BizContract{}
	if sc.getObject(PREFIX_BIZCONTRACT, addr, ret) != nil {
		return nil
	}
	return ret
}

func (sc supervisorContext) ValidateBizContract(address common.Address) bool {
	contract := sc.GetBizContract(address)
	return contract != nil
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
	return sc.setObject(PREFIX_CONFIG, contractAddr, config)
}

func (sc supervisorContext) GetConfig() *Config {
	config := &Config{}
	if sc.getObject(PREFIX_CONFIG, contractAddr, config) != nil {
		return nil
	}
	return config
}

func (sc supervisorContext) RecordData(addr common.Address, dataName string, msg string) error {
	bizContract := sc.GetBizContract(addr)
	if bizContract == nil {
		return ErrNotExisted
	}
	meta := sc.getBizMeta(bizContract.BizType)
	if meta == nil {
		return ErrBizMetaExisted
	}
	data := make(map[string]ReportField)
	if err := json.Unmarshal([]byte(msg), &data); err != nil {
		return err
	}

	checkType := func(t string, b byte) error {
		switch b {
		case abi.AddressTy:
			if t != "address" {
				return errors.New("expect 'address' to present Address")
			}
		case abi.StringTy:
			if t != "string" {
				return errors.New("expect 'string' to present string")
			}
		case abi.UintTy, abi.IntTy:
			if t != "int" && t != "uint32" && t != "uint64" && t != "int32" && t != "int64" {
				return errors.New("expect 'int' to present int")
			}
		case abi.BoolTy:
			if t != "bool" {
				return errors.New("expect 'bool' to present bool")
			}
		}
		return nil
	}

	// valid msg
	existed := false
	for _, d := range meta.Datas {
		if dataName == d.Name {
			existed = true
			for _, f := range d.Fields {
				if v, ok := data[f.Name]; ok {
					if err := checkType(f.Type, v.FieldType); err != nil {
						return err
					}
				} else {
					return ErrLackField
				}
			}
		}
	}

	if !existed {
		return ErrReportDataNotFound
	}

	topics := make([]common.Hash, 0)
	topics = append(topics, addr.Hash())
	vs := make([]ReportField, 0, len(data))
	for _, v := range data {
		vs = append(vs, v)
	}
	d := StructReport{
		BizType: meta.BizType,
		Datas:   vs,
	}
	bs, _ := json.Marshal(d)
	sc.context.GetStateDb().AddLog(&types.Log{
		Address:     contractAddr,
		Topics:      topics,
		Data:        bs,
		BlockNumber: sc.context.GetBlockNum().Uint64(),
	})
	return nil
}

type BizContract struct {
	Address   common.Address // 合约地址
	Owner     common.Address // 所有者地址
	Name      string         // 合约名称
	Desc      string         // 描述
	BizType   uint32         // 合约类型
	Status    uint32         // 状态
	TimeStamp *big.Int       // 创建时间
}

type MngNode struct {
	Id     uint32         // 编号
	Name   string         // 节点名称
	Desc   string         // 节点描述
	Ip     string         // 地址
	Status uint32         // 状态
	Addr   common.Address // 监管账户
	Pubkey string         // 公钥
}

type Config struct {
	AccountBlackLists []common.Address   // 账户黑名单
	MngNodes          map[string]MngNode // 管理员账户
}

type UpdateConfigReq struct {
	Cfg string
}

type RegContractReq struct {
	Addr    common.Address // 合约地址
	Owner   common.Address // 所有者地址
	BizType uint32         // 合约类型
}

type GetBizContractReq struct {
	Addr common.Address // 合约地址
}

type ReportReq struct {
	Addr     common.Address
	DataName string
	Msg      string
}

type BizMetaReq struct {
	Cfg string
}

type BizContractTpReq struct {
	BizType uint32
}

type ReportField struct {
	FieldType byte
	Value     interface{}
}

type StructReport struct {
	BizType string
	Datas   []ReportField
}
