package supervisor

import (
	"testing"
)

func Test(t *testing.T) {
	src := "tests/BizMeta.toml"
	out := "/Users/mac/gopath/src/github.com/vntchain/bottle/dev/contracts/biz.c"
	Gen(src, out)
}
