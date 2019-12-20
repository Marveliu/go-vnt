package supervisor

// var (
// 	ErrCreateVM = "must create with wavm"
// )
//
// type SupervisorVM struct {
// 	wavm.WAVM
// 	Supervised bool
// 	Monitor    Monitor
// }
//
// func NewSupervisorVM(vm vm.VM) *SupervisorVM {
// 	wavm, ok := vm.(*wavm.WAVM)
// 	if !ok {
// 		panic(ErrCreateVM)
// 	}
// 	svm := &SupervisorVM{*wavm, true, Monitor{}}
// 	return svm
// }
//
// func (s *SupervisorVM) Cancel() {
// 	s.WAVM.Cancel()
// }
//
// func (s *SupervisorVM) Create(caller vm.ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
// 	s.Monitor.SampleStart(caller.Address(), contractAddr, BizType_CREATE, gas, value)
// 	defer func() {
// 		s.Monitor.SampleEnd(s.getFuncNmae(), leftOverGas, false, err)
// 	}()
// 	return s.WAVM.Create(caller, code, gas, value)
// }
//
// func (s *SupervisorVM) Call(caller vm.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
// 	s.Monitor.SampleStart(caller.Address(), addr, BizType_CALL, gas, value)
// 	defer func() {
// 		s.Monitor.SampleEnd(s.getFuncNmae(), leftOverGas, false, err)
// 	}()
// 	return s.WAVM.Call(caller, addr, input, gas, value)
// }
//
// func (s *SupervisorVM) CallCode(caller vm.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
// 	return s.WAVM.CallCode(caller, addr, input, gas, value)
// }
//
// func (s *SupervisorVM) DelegateCall(caller vm.ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
// 	return s.WAVM.DelegateCall(caller, addr, input, gas)
// }
//
// func (s *SupervisorVM) StaticCall(caller vm.ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
// 	return s.WAVM.StaticCall(caller, addr, input, gas)
// }
//
// func (s *SupervisorVM) GetStateDb() inter.StateDB {
// 	return s.WAVM.GetStateDb()
// }
//
// func (s *SupervisorVM) ChainConfig() *params.ChainConfig {
// 	return s.WAVM.ChainConfig()
// }
//
// func (s *SupervisorVM) GetContext() vm.Context {
// 	return s.WAVM.GetContext()
// }
//
// func (s *SupervisorVM) getFuncNmae() string {
// 	if s.Wavm == nil {
// 		log.Error("Supervisor Not found wavm")
// 		return ""
// 	}
// 	return s.Wavm.GetFuncName()
// }
