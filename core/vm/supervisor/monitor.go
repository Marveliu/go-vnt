package supervisor

// var (
// 	ErrUnpairInvoke = errors.New("Capture start&end must come in pairs ")
// )
//
// const (
// 	BizType_CREATE       = "create"
// 	BizType_CALL         = "call"
// 	BizType_CallCode     = "callCode"
// 	BizType_DelegateCall = "delegateCall"
// 	BizType_StaticCall   = "staticCall"
// )
//
// type Digest struct {
// 	From     common.Address `json:"from"`
// 	To       common.Address `json:"to"`
// 	Contract common.Address `json:"contract"`
// 	Method   string         `json:"method"`
// 	Value    *big.Int       `json:"value"`
// 	Noce     uint64         `json:"bizNoce"`
// 	BizType  string         `json:"bizType"`
// 	GasPrice *big.Int       `json:"gas_price"`
// 	Gas      uint64         `json:"gas"`
// 	GasCost  uint64         `json:"gasCost"`
// 	Result   bool           `json:"result"`
// 	StartAt  time.Time      `json:"startAt"`
// 	Duration time.Duration  `json:"duration"`
// 	Err      error          `json:"error"`
// 	Extra    string         `json:"extra"`
// }
//
// type Monitor struct {
// 	Depth int
// 	Digs  []*Digest
// 	Err   error
// }
//
// func NewMonitor(vm wavm.WAVM) Monitor {
// 	m := Monitor{}
// 	return m
// }
//
// func (l *Monitor) SampleStart(from common.Address, to common.Address, bizType string, gas uint64, value *big.Int) error {
// 	dig := &Digest{
// 		From:    from,
// 		To:      to,
// 		BizType: bizType,
// 		Gas:     gas,
// 		Value:   value,
// 		StartAt: time.Now(),
// 	}
// 	l.Digs = append(l.Digs, dig)
// 	l.Depth++
// 	return nil
// }
//
// func (l *Monitor) SampleEnd(funcName string, usedGas uint64, failed bool, err error) error {
// 	if l.Depth > len(l.Digs) || l.Depth < 0 {
// 		return ErrUnpairInvoke
// 	}
// 	l.Depth--
// 	dig := l.Digs[l.Depth]
// 	dig.Result = !failed
// 	dig.Err = err
// 	dig.GasCost = usedGas
// 	dig.Duration = time.Since(dig.StartAt)
// 	dig.Method = funcName
// 	if l.Depth == 0 {
// 		l.Export()
// 	}
// 	return nil
// }
//
// func (l *Monitor) Export() error {
// 	if l.Depth != 0 {
// 		return ErrUnpairInvoke
// 	}
// 	for i, entry := range l.Digs {
// 		bytes, err := json.Marshal(entry)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Printf("index: %d, content: %s\n", i, string(bytes))
// 	}
// 	return nil
// }
