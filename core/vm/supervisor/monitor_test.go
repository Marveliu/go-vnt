package supervisor

// var (
// 	from = common.BytesToAddress([]byte("0"))
// 	to   = common.BytesToAddress([]byte("1"))
// 	addr = common.BytesToAddress([]byte("2"))
// )
//
// type context struct {
// 	Log   Monitor
// 	Depth int
// }
//
// func (c *context) recursive() {
// 	if err := c.Log.SampleStart(from, to, BizType_CALL, 123, big.NewInt(int64(rand.Intn(100)))); err != nil {
// 		panic(err)
// 	}
// 	defer func() {
// 		var (
// 			err error
// 		)
// 		if r := recover(); r != nil {
// 			err = errors.New("faied")
// 		}
// 		if err := c.Log.SampleEnd(uint64(100), false, err); err != nil {
// 			panic(err)
// 		}
// 	}()
// 	if c.Depth >= 10 {
// 		return
// 	}
// 	c.Depth++
// 	c.recursive()
// }
//
// func TestLogger(t *testing.T) {
// 	vm := wavm.WAVM{}
// 	vm.Wavm = &wavm.Wavm{}
// 	vm.Wavm.SetFuncName("1234")
// 	log := NewMonitor(vm)
// 	ctx := &context{
// 		Log: log,
// 	}
// 	ctx.recursive()
// 	ctx.Log.Export()
// }
