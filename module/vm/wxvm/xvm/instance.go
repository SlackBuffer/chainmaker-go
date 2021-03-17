package xvm

import (
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	"chainmaker.org/chainmaker-go/wxvm/xvm/exec"
	"chainmaker.org/chainmaker-go/wxvm/xvm/runtime/emscripten"
	"errors"
)

func CreateInstance(contextId int64, code exec.Code, method string, contractId *commonPb.ContractId, gasUsed uint64, gasLimit int64) (*wxvmInstance, error) {
	execCtx, err := code.NewContext(&exec.ContextConfig{
		GasLimit: gasLimit,
	})
	if err != nil {
		return nil, err
	}

	if err = emscripten.Init(execCtx); err != nil {
		return nil, err
	}

	execCtx.SetGasUsed(gasUsed)
	execCtx.SetUserData(contextIDKey, contextId)
	instance := &wxvmInstance{
		method:  method,
		execCtx: execCtx,
	}
	return instance, nil
}

type wxvmInstance struct {
	method  string
	execCtx exec.Context
}

func (x *wxvmInstance) Exec() error {
	mem := x.execCtx.Memory()
	if mem == nil {
		return errors.New("bad contract, no memory")
	}

	function := "_" + x.method
	_, err := x.execCtx.Exec(function, []int64{})
	return err
}

func (x *wxvmInstance) ResourceUsed() Limits {
	limits := Limits{
		Cpu: x.execCtx.GasUsed(),
	}
	return limits
}

func (x *wxvmInstance) Release() {
	x.execCtx.Release()
}

func (x *wxvmInstance) Abort(msg string) {
	exec.Throw(exec.NewTrap(msg))
}