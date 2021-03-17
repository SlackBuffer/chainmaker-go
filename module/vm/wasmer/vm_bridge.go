/*
Copyright (C) BABEC. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package wasmer

import (
	"chainmaker.org/chainmaker-go/common/json"
	"chainmaker.org/chainmaker-go/logger"
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	"chainmaker.org/chainmaker-go/protocol"
	wasm "chainmaker.org/chainmaker-go/wasmer/wasmer-go"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"unsafe"
)

// #include <stdlib.h>

// extern int sysCall(void *context, int requestHeaderPtr, int requestHeaderLen, int requestBodyPtr, int requestBodyLen);
// extern void logMessage(void *context, int pointer, int length);
// extern int getStateLen(void *context, int ctxPtr, int keyPtr, int keyLen,  int fieldPtr, int fieldLen, int valueLenPtr);
// extern int getState(void *context, int ctxPtr, int keyPtr, int keyLen,  int fieldPtr, int fieldLen, int valuePtr, int valueLen);
// extern int putState(void *context, int ctxPtr, int keyPtr, int keyLen,  int fieldPtr, int fieldLen, int valuePtr, int valueLen);
// extern int deleteState(void *context,  int ctxPtr, int keyPtr, int keyLen,  int fieldPtr, int fieldLen);
// extern int successResult(void *context, int ctxPtr, int dataPtr, int dataLen);
// extern int errorResult(void *context, int ctxPtr, int dataPtr, int dataLen);
// extern int callContract(void *context, int ctxParameterPtr, int ctxParameterLen, int parameterPtr, int parameterLen);
// extern int callContractLen(void *context, int ctxParameterPtr, int ctxParameterLen, int parameterPtr, int parameterLen);
// extern int fdWrite(void *contextfd,int iovs,int iovsPtr ,int iovsLen,int nwrittenPtr);
// extern int fdRead(void *contextfd,int iovs,int iovsPtr ,int iovsLen,int nwrittenPtr);
// extern int fdClose(void *contextfd,int iovs,int iovsPtr ,int iovsLen,int nwrittenPtr);
// extern int fdSeek(void *contextfd,int iovs,int iovsPtr ,int iovsLen,int nwrittenPtr);
// extern void procExit(void *contextfd,int exitCode);
import "C"

var log = logger.GetLogger(logger.MODULE_VM)

// sdkRequestCtx record wasmer vm request parameter
type sdkRequestCtx struct {
	Sc            *SimContext
	RequestHeader *RequestHeader // sdk request common json param
	RequestBody   []byte         // sdk request param
	Memory        []byte         // cache call method GetStateLen value result
}

// LogMessage print log to file
func (s *sdkRequestCtx) LogMessage() int32 {
	s.Sc.Log.Debugf("waci log>> [%s] %s", s.Sc.TxSimContext.GetTx().Header.TxId, string(s.RequestBody))
	return protocol.ContractSdkSignalResultSuccess
}

// logMessage print log to file
//export logMessage
func logMessage(context unsafe.Pointer, pointer int32, length int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	gotText := string(memory[pointer : pointer+length])
	log.Debugf("wasm log>> " + gotText)
}

// sysCall wasmer vm call chain entry
//export sysCall
func sysCall(context unsafe.Pointer, requestHeaderPtr int32, requestHeaderLen int32, requestBodyPtr int32, requestBodyLen int32) int32 {
	if requestHeaderLen == 0 {
		log.Error("wasm log>>requestHeader is null.")
		return protocol.ContractSdkSignalResultFail
	}
	// get param from memory
	instanceContext := wasm.IntoInstanceContext(context)
	memory := instanceContext.Memory().Data()

	requestHeaderByte := make([]byte, requestHeaderLen)
	copy(requestHeaderByte, memory[requestHeaderPtr:requestHeaderPtr+requestHeaderLen])
	requestBody := make([]byte, requestBodyLen)
	copy(requestBody, memory[requestBodyPtr:requestBodyPtr+requestBodyLen])

	requestHeader := &RequestHeader{}
	if err := json.Unmarshal(requestHeaderByte, requestHeader); err != nil {
		log.Errorf("unmarshal requestHeader param fail: %s requestHeader=%s  requestBody=%s", err.Error(), string(requestHeaderByte), string(requestBody))
		return protocol.ContractSdkSignalResultFail
	}

	vbm := GetVmBridgeManager()
	sc := vbm.get(requestHeader.CtxPtr)

	s := &sdkRequestCtx{
		Sc:            sc,
		RequestHeader: requestHeader,
		RequestBody:   requestBody,
		Memory:        memory,
	}

	method := requestHeader.Method
	switch method {
	case protocol.ContractMethodLogMessage:
		return s.LogMessage()
	case protocol.ContractMethodGetStateLen:
		return s.GetStateLen()
	case protocol.ContractMethodGetState:
		return s.GetState()
	case protocol.ContractMethodPutState:
		return s.PutState()
	case protocol.ContractMethodDeleteState:
		return s.DeleteState()
	case protocol.ContractMethodSuccessResult:
		return s.SuccessResult()
	case protocol.ContractMethodErrorResult:
		return s.ErrorResult()
	case protocol.ContractMethodCallContract:
		return s.CallContract()
	case protocol.ContractMethodCallContractLen:
		return s.CallContractLen()
	default:
		log.Errorf("method is %s not match.", method)
	}
	return protocol.ContractSdkSignalResultFail
}

// GetStateLen get state length from chain
func (s *sdkRequestCtx) GetStateLen() int32 {
	return s.getStateCore(true)
}

// GetStateLen get state from chain
func (s *sdkRequestCtx) GetState() int32 {
	return s.getStateCore(false)
}

func (s *sdkRequestCtx) getStateCore(isGetLen bool) int32 {
	req := &GetStateRequest{}
	if err := json.Unmarshal(s.RequestBody, req); err != nil {
		msg := fmt.Sprintf("unmarshal method getStateCore request param fail. param=%s error: %s", string(s.RequestBody), err.Error())
		return s.recordMsg(msg)
	}
	if err := protocol.CheckKeyFieldStr(req.Key, req.Field); err != nil {
		return s.recordMsg(err.Error())
	}

	valuePtr := req.ValuePtr
	if isGetLen {
		contractName := s.Sc.ContractId.ContractName
		value, err := s.Sc.TxSimContext.Get(contractName, protocol.GetKeyStr(req.Key, req.Field))
		if err != nil {
			msg := fmt.Sprintf("method getStateCore get fail. key=%s, field=%s, error:%s", req.Key, req.Field, err.Error())
			return s.recordMsg(msg)
		}
		copy(s.Memory[valuePtr:valuePtr+4], IntToBytes(int32(len(value))))
		s.Sc.GetStateCache = value
	} else {
		len := int32(len(s.Sc.GetStateCache))
		if len != 0 {
			copy(s.Memory[valuePtr:valuePtr+len], s.Sc.GetStateCache)
			s.Sc.GetStateCache = nil
		}
	}
	return protocol.ContractSdkSignalResultSuccess
}

// PutState put state to chain
func (s *sdkRequestCtx) PutState() int32 {

	req := &PutStateRequest{}
	if err := json.Unmarshal(s.RequestBody, req); err != nil {
		msg := fmt.Sprintf("unmarshal method PutState request param fail. param=%s error: %s", string(s.RequestBody), err.Error())
		return s.recordMsg(msg)
	}
	if err := protocol.CheckKeyFieldStr(req.Key, req.Field); err != nil {
		return s.recordMsg(err.Error())
	}

	contractName := s.Sc.ContractId.ContractName
	err := s.Sc.TxSimContext.Put(contractName, protocol.GetKeyStr(req.Key, req.Field), []byte(req.Value))
	if err != nil {
		return s.recordMsg("method PutState put fail. " + err.Error())
	}
	return protocol.ContractSdkSignalResultSuccess
}

// DeleteState delete state from chain
func (s *sdkRequestCtx) DeleteState() int32 {

	req := &DeleteStateRequest{}
	if err := json.Unmarshal(s.RequestBody, req); err != nil {
		return s.recordMsg("unmarshal method DeleteState request param fail: " + err.Error())
	}
	if err := protocol.CheckKeyFieldStr(req.Key, req.Field); err != nil {
		return s.recordMsg(err.Error())
	}

	contractName := s.Sc.ContractId.ContractName
	err := s.Sc.TxSimContext.Del(contractName, protocol.GetKeyStr(req.Key, req.Field))
	if err != nil {
		return s.recordMsg(err.Error())
	}

	return protocol.ContractSdkSignalResultSuccess
}

// SuccessResult record the results of contract execution success
func (s *sdkRequestCtx) SuccessResult() int32 {
	s.Sc.ContractResult.Code = commonPb.ContractResultCode_OK
	s.Sc.ContractResult.Result = s.RequestBody
	return protocol.ContractSdkSignalResultSuccess
}

// ErrorResult record the results of contract execution error
func (s *sdkRequestCtx) ErrorResult() int32 {
	s.Sc.ContractResult.Code = commonPb.ContractResultCode_FAIL
	s.Sc.ContractResult.Message = string(s.RequestBody)
	return protocol.ContractSdkSignalResultSuccess
}

//  CallContractLen invoke cross contract calls, save result to cache and putout result length
func (s *sdkRequestCtx) CallContractLen() int32 {
	return s.callContractCore(true)
}

//  CallContractLen get cross contract call result from cache
func (s *sdkRequestCtx) CallContract() int32 {
	return s.callContractCore(false)
}

func (s *sdkRequestCtx) callContractCore(isGetLen bool) int32 {

	req := &CallContractRequest{}
	if err := json.Unmarshal(s.RequestBody, req); err != nil {
		return s.recordMsg("unmarshal method CallContract request param fail: " + err.Error())
	}

	param := req.Param
	valuePtr := req.ValuePtr
	contractName := req.ContractName
	method := req.Method

	if !isGetLen { // get value from cache
		result := s.Sc.TxSimContext.GetCurrentResult()
		copy(s.Memory[valuePtr:valuePtr+int32(len(result))], result)
		return protocol.ContractSdkSignalResultSuccess
	}

	// check param
	if len(contractName) == 0 {
		return s.recordMsg("CallContract contractName is null")
	}
	if len(method) == 0 {
		return s.recordMsg("CallContract method is null")
	}
	if len(param) > protocol.ParametersKeyMaxCount {
		return s.recordMsg("expect less than 20 parameters, but get " + strconv.Itoa(len(param)))
	}
	for key, value := range param {
		if len(key) > protocol.DefaultStateLen {
			msg := fmt.Sprintf("CallContract param expect Key length less than %d, but get %d", protocol.DefaultStateLen, len(key))
			return s.recordMsg(msg)
		}
		match, err := regexp.MatchString(protocol.DefaultStateRegex, key)
		if err != nil || !match {
			msg := fmt.Sprintf("CallContract param expect Key no special characters, but get %s. letter, number, dot and underline are allowed", key)
			return s.recordMsg(msg)
		}
		if len(value) > protocol.ParametersValueMaxLength {
			msg := fmt.Sprintf("expect Value length less than %d, but get %d", protocol.ParametersValueMaxLength, len(value))
			return s.recordMsg(msg)
		}
	}
	if err := protocol.CheckKeyFieldStr(contractName, method); err != nil {
		return s.recordMsg(err.Error())
	}

	// call contract
	usedGas := s.Sc.Instance.GetGasUsed() + protocol.CallContractGasOnce
	s.Sc.Instance.SetGasUsed(usedGas)
	result, code := s.Sc.TxSimContext.CallContract(&commonPb.ContractId{ContractName: contractName}, method, nil, param, usedGas, commonPb.TxType_INVOKE_USER_CONTRACT)
	usedGas = s.Sc.Instance.GetGasUsed() + uint64(result.GasUsed)
	s.Sc.Instance.SetGasUsed(usedGas)
	if code != commonPb.TxStatusCode_SUCCESS {
		return s.recordMsg("CallContract " + code.String() + ", msg:" + result.Message)
	}
	// set value length to memory
	l := IntToBytes(int32(len(result.Result)))
	copy(s.Memory[valuePtr:valuePtr+4], l)
	return protocol.ContractSdkSignalResultSuccess
}

// wasi
//export fdWrite
func fdWrite(context unsafe.Pointer, fd int32, iovsPtr int32, iovsLen int32, nwrittenPtr int32) (err int32) {
	return protocol.ContractSdkSignalResultSuccess
}

//export fdRead
func fdRead(context unsafe.Pointer, fd int32, iovsPtr int32, iovsLen int32, nwrittenPtr int32) (err int32) {
	return protocol.ContractSdkSignalResultSuccess
}

//export fdClose
func fdClose(context unsafe.Pointer, fd int32, iovsPtr int32, iovsLen int32, nwrittenPtr int32) (err int32) {
	return protocol.ContractSdkSignalResultSuccess
}

//export fdSeek
func fdSeek(context unsafe.Pointer, fd int32, iovsPtr int32, iovsLen int32, nwrittenPtr int32) (err int32) {
	return protocol.ContractSdkSignalResultSuccess
}

//export procExit
func procExit(context unsafe.Pointer, exitCode int32) {
	panic("exit called by contract, code:" + strconv.Itoa(int(exitCode)))
}

func (s *sdkRequestCtx) recordMsg(msg string) int32 {
	s.Sc.ContractResult.Message += msg
	s.Sc.ContractResult.Code = commonPb.ContractResultCode_FAIL
	msg = s.Sc.ContractId.ContractName + ":" + msg
	s.Sc.Log.Errorf("wasm log>> " + msg)
	return protocol.ContractSdkSignalResultFail
}

var (
	vmBridgeManagerMutex = &sync.Mutex{}
	bridgeSingleton      *vmBridgeManager
)

type vmBridgeManager struct {
	//wasmImports *wasm.Imports
	pointerLock     sync.Mutex
	simContextCache map[int32]*SimContext
}

// GetVmBridgeManager get singleton vmBridgeManager struct
func GetVmBridgeManager() *vmBridgeManager {
	if bridgeSingleton == nil {
		vmBridgeManagerMutex.Lock()
		defer vmBridgeManagerMutex.Unlock()
		if bridgeSingleton == nil {
			log.Debugf("init vmBridgeManager")
			bridgeSingleton = &vmBridgeManager{}
			bridgeSingleton.simContextCache = make(map[int32]*SimContext)
			//bridgeSingleton.wasmImports = bridgeSingleton.GetImports()
		}
	}
	return bridgeSingleton
}

// put the context
func (b *vmBridgeManager) put(k int32, v *SimContext) {
	b.pointerLock.Lock()
	defer b.pointerLock.Unlock()
	b.simContextCache[k] = v
}

// get the context
func (b *vmBridgeManager) get(k int32) *SimContext {
	b.pointerLock.Lock()
	defer b.pointerLock.Unlock()
	return b.simContextCache[k]
}

// remove the context
func (b *vmBridgeManager) remove(k int32) {
	b.pointerLock.Lock()
	defer b.pointerLock.Unlock()
	delete(b.simContextCache, k)
}

// NewWasmInstance new wasm instance. Apply for new memory.
func (b *vmBridgeManager) NewWasmInstance(byteCode []byte) (wasm.Instance, error) {
	return wasm.NewInstanceWithImports(byteCode, b.GetImports())
}

// GetImports return export interface to cgo
func (b *vmBridgeManager) GetImports() *wasm.Imports {
	imports := wasm.NewImports().Namespace("env")
	// parameter explain:  1、["log_message"]: rust extern "C" method name 2、[logMessage] go method ptr 3、[C.logMessage] cgo function pointer.
	imports.Append("sys_call", sysCall, C.sysCall)
	imports.Append("log_message", logMessage, C.logMessage)
	// for waci empty interface
	imports.Namespace("wasi_unstable")
	imports.Append("fd_write", fdWrite, C.fdWrite)
	imports.Append("fd_read", fdRead, C.fdRead)
	imports.Append("fd_close", fdClose, C.fdClose)
	imports.Append("fd_seek", fdSeek, C.fdSeek)

	imports.Namespace("wasi_snapshot_preview1")
	imports.Append("proc_exit", procExit, C.procExit)
	//imports.Append("fd_write", fdWrite2, C.fdWrite2)
	//imports.Append("environ_sizes_get", fdWrite, C.fdWrite)
	//imports.Append("proc_exit", fdWrite, C.fdWrite)
	//imports.Append("environ_get", fdWrite, C.fdWrite)

	return imports
}