package wavm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/common/math"
	"github.com/vntchain/go-vnt/core/vm"
	"github.com/vntchain/go-vnt/core/vm/interface"
)

const (
	BizType_CREATE       = "create"
	BizType_CALL         = "call"
	BizType_CallCode     = "callCode"
	BizType_DelegateCall = "delegateCall"
	BizType_StaticCall   = "staticCall"
)

var (
	ErrUnpairInvoke = errors.New("Capture start&end must come in pairs ")
)

type WasmLogger struct {
	cfg       vm.LogConfig
	logs      []StructLog
	debugLogs []DebugLog

	// Dig
	digests []*Digest
	depth   int

	output []byte
	err    error
}

type Digest struct {
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Contract common.Address `json:"contract"`
	Method   []byte         `json:"method"`
	Value    *big.Int       `json:"value"`
	Noce     uint64         `json:"bizNoce"`
	BizType  string         `json:"bizType"`
	GasPrice *big.Int       `json:"gas_price"`
	Gas      uint64         `json:"gas"`
	GasCost  uint64         `json:"gasCost"`
	Result   bool           `json:"result"`
	StartAt  time.Time      `json:"startAt"`
	Duration time.Duration  `json:"duration"`
	Err      string         `json:"err"`
	Extra    string         `json:"extra"`
}

type StructLog struct {
	Pc         uint64                      `json:"pc"`
	Op         vm.OPCode                   `json:"op"`
	Gas        uint64                      `json:"gas"`
	GasCost    uint64                      `json:"gasCost"`
	Memory     []byte                      `json:"-"`
	MemorySize int                         `json:"-"`
	Stack      []*big.Int                  `json:"-"`
	Storage    map[common.Hash]common.Hash `json:"-"`
	Depth      int                         `json:"depth"`
	Err        error                       `json:"error"`
}

type DebugLog struct {
	PrintMsg string `json:"printMsg"`
}

// overrides for gencodec
type structLogMarshaling struct {
	Gas         math.HexOrDecimal64
	GasCost     math.HexOrDecimal64
	OpName      string `json:"opName"` // adds call to OpName() in MarshalJSON
	ErrorString string `json:"error"`  // adds call to ErrorString() in MarshalJSON
}

func (s *StructLog) OpName() string {
	return s.Op.String()
}

func (s *StructLog) ErrorString() string {
	if s.Err != nil {
		return s.Err.Error()
	}
	return ""
}

func NewWasmLogger(cfg *vm.LogConfig) *WasmLogger {
	logger := &WasmLogger{}
	if cfg != nil {
		logger.cfg = *cfg
	}
	return logger
}

func (l *WasmLogger) CaptureStart(from common.Address, to common.Address, call bool, input []byte, gas uint64, value *big.Int) error {
	return nil
}
func (l *WasmLogger) CaptureState(env vm.VM, pc uint64, op vm.OPCode, gas, cost uint64, contract inter.Contract, depth int, err error) error {
	// check if already accumulated the specified number of logs
	if l.cfg.Limit != 0 && l.cfg.Limit <= len(l.logs) {
		return vm.ErrTraceLimitReached
	}

	// create a new snaptshot of the VM.
	log := StructLog{pc, op, gas, cost, nil, 0, nil, nil, depth, err}

	l.logs = append(l.logs, log)
	return nil
}
func (l *WasmLogger) CaptureLog(env vm.VM, msg string) error {
	log := DebugLog{msg}
	l.debugLogs = append(l.debugLogs, log)
	return nil
}
func (l *WasmLogger) CaptureFault(env vm.VM, pc uint64, op vm.OPCode, gas, cost uint64, contract inter.Contract, depth int, err error) error {
	return nil
}
func (l *WasmLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) error {
	return nil
}

// Error returns the VM error captured by the trace.
func (l *WasmLogger) Error() error { return l.err }

// Output returns the VM return value captured by the trace.
func (l *WasmLogger) Output() []byte { return l.output }

// WasmLogger returns the captured log entries.
func (l *WasmLogger) StructLogs() []StructLog { return l.logs }

// DebugLogs returns the captured debug log entries.
func (l *WasmLogger) DebugLogs() []DebugLog { return l.debugLogs }

// WriteTrace writes a formatted trace to the given writer
func WriteTrace(writer io.Writer, logs []StructLog) {

}

func (l *WasmLogger) SampleStart(from common.Address, to common.Address, bizType string, input []byte, gas uint64, value *big.Int) error {
	dig := &Digest{
		From:    from,
		To:      to,
		BizType: bizType,
		Gas:     gas,
		Value:   value,
		StartAt: time.Now(),
	}
	if input != nil {
		dig.Method = input[:4]
	}
	l.digests = append(l.digests, dig)
	l.depth++
	return nil
}

func (l *WasmLogger) SampleEnd(contract common.Address, usedGas uint64, failed bool, err error) error {
	if l.depth > len(l.digests) || l.depth < 0 {
		return ErrUnpairInvoke
	}
	l.depth--
	dig := l.digests[l.depth]
	dig.Result = true
	if err != nil {
		dig.Err = err.Error()
	}
	dig.GasCost = usedGas
	dig.Contract = contract
	dig.Duration = time.Since(dig.StartAt)
	if l.depth == 0 {
		// 落入到区块再导出
		l.Export()
	}
	return nil
}

func (l *WasmLogger) Export() error {
	if l.depth != 0 {
		return ErrUnpairInvoke
	}
	for i, entry := range l.digests {
		bytes, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		fmt.Printf("index: %d, content: %s\n", i, string(bytes))
	}
	return nil
}
