/*
Copyright (C) BABEC. All rights reserved.

SPDX-License-Identifier: Apache-2.0

Wacsi WebAssembly chainmaker system interface
*/
package wasi

import (
	"fmt"
	"regexp"
	"sync/atomic"

	"chainmaker.org/chainmaker-go/common/serialize"
	"chainmaker.org/chainmaker-go/logger"
	"chainmaker.org/chainmaker-go/pb/protogo/common"
	"chainmaker.org/chainmaker-go/protocol"
	"chainmaker.org/chainmaker-go/store/statedb/statesqldb"
	"chainmaker.org/chainmaker-go/store/types"
	"chainmaker.org/chainmaker-go/utils"
)

var ErrorNotManageContract = fmt.Errorf("method not init_contract or upgrade")

type Bool int32

const boolTrue Bool = 1
const boolFalse Bool = 0

// Wacsi WebAssembly chainmaker system interface
type Wacsi interface {
	// state operation
	PutState(requestBody []byte, contractName string, txSimContext protocol.TxSimContext) error
	GetState(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte, data []byte, isLen bool) ([]byte, error)
	DeleteState(requestBody []byte, contractName string, txSimContext protocol.TxSimContext) error
	// call other contract
	CallContract(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte, data []byte, gasUsed uint64, isLen bool) ([]byte, error, uint64)
	// result record
	SuccessResult(contractResult *common.ContractResult, data []byte) int32
	ErrorResult(contractResult *common.ContractResult, data []byte) int32
	// emit event
	EmitEvent(requestBody []byte, txSimContext protocol.TxSimContext, contractId *common.ContractId, log *logger.CMLogger) (*common.ContractEvent, error)
	// kv iterator
	KvIterator(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte) error
	KvIteratorHasNext(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte) error
	KvIteratorNext(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte, data []byte, contractName string, isLen bool) ([]byte, error)
	KvIteratorClose(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte) error

	// sql operation
	ExecuteQuery(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte, chainId string) error
	ExecuteQueryOne(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte, data []byte, chainId string, isLen bool) ([]byte, error)
	ExecuteUpdate(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte, chainId string) error
	ExecuteDDL(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte, method string) error
	RSHasNext(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte) error
	RSNext(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte, data []byte, isLen bool) ([]byte, error)
	RSClose(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte) error
}

type WacsiImpl struct {
	verifySql *types.StandardSqlVerify
	rowIndex  int32
}

func NewWacsi() Wacsi {
	return &WacsiImpl{
		verifySql: &types.StandardSqlVerify{},
		rowIndex:  0,
	}
}

func (*WacsiImpl) PutState(requestBody []byte, contractName string, txSimContext protocol.TxSimContext) error {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	key, _ := ec.GetString("key")
	field, _ := ec.GetString("field")
	value, _ := ec.GetBytes("value")
	if err := protocol.CheckKeyFieldStr(key, field); err != nil {
		return err
	}
	err := txSimContext.Put(contractName, protocol.GetKeyStr(key, field), value)
	return err
}

func (*WacsiImpl) GetState(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte, data []byte, isLen bool) ([]byte, error) {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	key, _ := ec.GetString("key")
	field, _ := ec.GetString("field")
	valuePtr, _ := ec.GetInt32("value_ptr")
	if err := protocol.CheckKeyFieldStr(key, field); err != nil {
		return nil, err
	}

	if !isLen {
		copy(memory[valuePtr:valuePtr+int32(len(data))], data)
		return nil, nil
	}
	value, err := txSimContext.Get(contractName, protocol.GetKeyStr(key, field))
	if err != nil {
		msg := fmt.Errorf("method getStateCore get fail. key=%s, field=%s, error:%s", key, field, err.Error())
		return nil, msg
	}
	copy(memory[valuePtr:valuePtr+4], utils.IntToBytes(int32(len(value))))
	if len(value) == 0 {
		return nil, nil
	}
	return value, nil
}

func (*WacsiImpl) DeleteState(requestBody []byte, contractName string, txSimContext protocol.TxSimContext) error {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	key, _ := ec.GetString("key")
	field, _ := ec.GetString("field")
	if err := protocol.CheckKeyFieldStr(key, field); err != nil {
		return err
	}

	err := txSimContext.Del(contractName, protocol.GetKeyStr(key, field))
	if err != nil {
		return err
	}
	return nil
}
func (*WacsiImpl) CallContract(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte, data []byte, gasUsed uint64, isLen bool) ([]byte, error, uint64) {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	valuePtr, _ := ec.GetInt32("value_ptr")
	contractName, _ := ec.GetString("contract_name")
	method, _ := ec.GetString("method")
	param, _ := ec.GetBytes("param")

	ecData := serialize.NewEasyCodecWithBytes(param)
	paramItem := ecData.GetItems()

	if !isLen { // get value from cache
		result := txSimContext.GetCurrentResult()
		copy(memory[valuePtr:valuePtr+int32(len(result))], result)
		return nil, nil, gasUsed
	}

	// check param
	if len(contractName) == 0 {
		return nil, fmt.Errorf("CallContract contractName is null"), gasUsed
	}
	if len(method) == 0 {
		return nil, fmt.Errorf("CallContract method is null"), gasUsed
	}
	if len(paramItem) > protocol.ParametersKeyMaxCount {
		return nil, fmt.Errorf("expect less than %d parameters, but get %d", protocol.ParametersKeyMaxCount, len(paramItem)), gasUsed
	}
	for _, item := range paramItem {
		if len(item.Key) > protocol.DefaultStateLen {
			return nil, fmt.Errorf("CallContract param expect Key length less than %d, but get %d", protocol.DefaultStateLen, len(item.Key)), gasUsed
		}
		match, err := regexp.MatchString(protocol.DefaultStateRegex, item.Key)
		if err != nil || !match {
			return nil, fmt.Errorf("CallContract param expect Key no special characters, but get %s. letter, number, dot and underline are allowed", item.Key), gasUsed
		}
		if len(item.Value.(string)) > protocol.ParametersValueMaxLength {
			return nil, fmt.Errorf("expect Value length less than %d, but get %d", protocol.ParametersValueMaxLength, len(item.Value.(string))), gasUsed
		}
	}
	if err := protocol.CheckKeyFieldStr(contractName, method); err != nil {
		return nil, err, gasUsed
	}

	// call contract
	gasUsed += protocol.CallContractGasOnce
	paramMap := ecData.ToMap()
	result, code := txSimContext.CallContract(&common.ContractId{ContractName: contractName}, method, nil, paramMap, gasUsed, common.TxType_INVOKE_USER_CONTRACT)
	gasUsed += uint64(result.GasUsed)
	if code != common.TxStatusCode_SUCCESS {
		return nil, fmt.Errorf("CallContract %s, , msg: %s", code.String(), result.Message), gasUsed
	}
	// set value length to memory
	l := utils.IntToBytes(int32(len(result.Result)))
	copy(memory[valuePtr:valuePtr+4], l)
	if len(result.Result) == 0 {
		return nil, nil, gasUsed
	}
	return result.Result, nil, gasUsed
}

func (*WacsiImpl) SuccessResult(contractResult *common.ContractResult, data []byte) int32 {
	if contractResult.Code == common.ContractResultCode_FAIL {
		return protocol.ContractSdkSignalResultFail
	}
	contractResult.Code = common.ContractResultCode_OK
	contractResult.Result = data
	return protocol.ContractSdkSignalResultSuccess
}

func (*WacsiImpl) ErrorResult(contractResult *common.ContractResult, data []byte) int32 {
	contractResult.Code = common.ContractResultCode_FAIL
	if len(contractResult.Message) > 0 {
		contractResult.Message += ". contract message:" + string(data)
	} else {
		contractResult.Message = "contract message:" + string(data)
	}
	return protocol.ContractSdkSignalResultSuccess
}

// EmitEvent emit event to chain
func (w *WacsiImpl) EmitEvent(requestBody []byte, txSimContext protocol.TxSimContext, contractId *common.ContractId, log *logger.CMLogger) (*common.ContractEvent, error) {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	topic, err := ec.GetString("topic")
	if err != nil {
		return nil, fmt.Errorf("emit event : get topic err")
	}
	if err := protocol.CheckTopicStr(topic); err != nil {
		return nil, err
	}

	req := ec.GetItems()
	var eventData []string
	for i := 1; i < len(req); i++ {
		data := req[i].Value.(string)
		eventData = append(eventData, data)
		log.Debugf("EmitEvent EventData :%v", data)
	}

	if err := protocol.CheckEventData(eventData); err != nil {
		return nil, err
	}

	contractEvent := &common.ContractEvent{
		ContractName:    contractId.ContractName,
		ContractVersion: contractId.ContractVersion,
		Topic:           topic,
		TxId:            txSimContext.GetTx().Header.TxId,
		EventData:       eventData,
	}
	ddl := utils.GenerateSaveContractEventDdl(contractEvent, "chainId", 1, 1)
	count := utils.GetSqlStatementCount(ddl)
	if count != 1 {
		return nil, fmt.Errorf("contract event parameter error,exist sql injection")
	}

	return contractEvent, nil
}

//author:whang1234
func (w *WacsiImpl) KvIterator(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte) error {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	startKey, _ := ec.GetString("start_key")
	startField, _ := ec.GetString("start_field")
	limitKey, _ := ec.GetString("limit_key")
	limitField, _ := ec.GetString("limit_field")
	valuePtr, _ := ec.GetInt32("value_ptr")
	if err := protocol.CheckKeyFieldStr(startKey, startField); err != nil { //加判断
		return err
	}
	if err := protocol.CheckKeyFieldStr(limitKey, limitField); err != nil { //加判断
		return err
	}

	key := protocol.GetKeyStr(startKey, startField)
	limit := protocol.GetKeyStr(limitKey, limitField)
	iter, err := txSimContext.Select(contractName, key, limit)
	if err != nil {
		return fmt.Errorf("ctx query error, %s", err.Error())
	}

	index := atomic.AddInt32(&w.rowIndex, 1)
	txSimContext.SetStateKvHandle(index, iter)
	copy(memory[valuePtr:valuePtr+4], utils.IntToBytes(index))
	return nil
}

func (*WacsiImpl) KvIteratorHasNext(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte) error {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	kvIndex, _ := ec.GetInt32("rs_index")
	valuePtr, _ := ec.GetInt32("value_ptr")

	// get
	kvRows, ok := txSimContext.GetStateKvHandle(kvIndex)
	if !ok {
		return fmt.Errorf("KvHasNext:ctx can not found rs_index[%d]", kvIndex)
	}

	index := boolFalse
	if kvRows.Next() {
		index = boolTrue
	}
	copy(memory[valuePtr:valuePtr+4], utils.IntToBytes(int32(index)))
	return nil
}

func (*WacsiImpl) KvIteratorNext(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte, data []byte, contractname string, isLen bool) ([]byte, error) {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	kvIndex, _ := ec.GetInt32("rs_index")
	ptr, _ := ec.GetInt32("value_ptr")

	// get handle
	kvRows, ok := txSimContext.GetStateKvHandle(kvIndex)
	if !ok {
		return nil, fmt.Errorf("KvGetNextState:ctx can not found rs_index[%d]", kvIndex)
	}
	// get data
	if !isLen {
		copy(memory[ptr:ptr+int32(len(data))], data)
		return nil, nil
	}
	// get len
	ec = serialize.NewEasyCodec()
	if kvRows != nil {
		kvRow, err := kvRows.Value()
		if err != nil {
			return nil, fmt.Errorf("ctx iterator next data error, %s", err.Error())
		}
		key := parseStateKey(kvRow.Key, contractname)
		value := kvRow.Value
		ec.AddString("key", string(key))
		ec.AddBytes("value", value)
	}
	kvBytes := ec.Marshal()
	copy(memory[ptr:ptr+4], utils.IntToBytes(int32(len(kvBytes))))
	return kvBytes, nil
}

// parseStateKey corresponding to the constructStateKey(),  delete contract name from leveldb key
func parseStateKey(key []byte, contractName string) []byte {
	return key[len(contractName)+1:]
}

func (w *WacsiImpl) KvIteratorClose(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte) error {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	kvIndex, _ := ec.GetInt32("rs_index")
	valuePtr, _ := ec.GetInt32("value_ptr")
	// get
	kvRows, ok := txSimContext.GetStateKvHandle(kvIndex)
	if !ok {
		return fmt.Errorf("kv close:ctx can not found rs_index[%d]", kvIndex)
	}
	kvRows.Release()
	copy(memory[valuePtr:valuePtr+4], utils.IntToBytes(1))
	return nil
}

func (w *WacsiImpl) ExecuteQuery(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte, chainId string) error {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	sql, _ := ec.GetString("sql")
	ptr, _ := ec.GetInt32("value_ptr")

	// verify
	if err := w.verifySql.VerifyDQLSql(sql); err != nil {
		return fmt.Errorf("verify query sql error, %s", err.Error())
	}

	// execute query
	var rows protocol.SqlRows
	var err error
	if txSimContext.GetTx().GetHeader().TxType == common.TxType_QUERY_USER_CONTRACT {
		rows, err = txSimContext.GetBlockchainStore().QueryMulti(contractName, sql)
		if err != nil {
			return fmt.Errorf("ctx query error, %s", err.Error())
		}
	} else {
		txKey := common.GetTxKeyWith(txSimContext.GetBlockProposer(), txSimContext.GetBlockHeight())
		transaction, err := txSimContext.GetBlockchainStore().GetDbTransaction(txKey)
		if err != nil {
			return fmt.Errorf("ctx get db transaction error, [%s]", err.Error())
		}
		changeCurrentDB(chainId, contractName, transaction)
		rows, err = transaction.QueryMulti(sql)
	}

	index := atomic.AddInt32(&w.rowIndex, 1)
	txSimContext.SetStateSqlHandle(index, rows)
	copy(memory[ptr:ptr+4], utils.IntToBytes(index))
	return nil
}

func (w *WacsiImpl) ExecuteQueryOne(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte, data []byte, chainId string, isLen bool) ([]byte, error) {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	sql, _ := ec.GetString("sql")
	ptr, _ := ec.GetInt32("value_ptr")

	// verify
	if err := w.verifySql.VerifyDQLSql(sql); err != nil {
		return nil, fmt.Errorf("verify query one sql error, %s", err.Error())
	}

	// get len
	if isLen {
		// execute
		var row protocol.SqlRow
		var err error
		if txSimContext.GetTx().GetHeader().TxType == common.TxType_QUERY_USER_CONTRACT {
			row, err = txSimContext.GetBlockchainStore().QuerySingle(contractName, sql)
			if err != nil {
				return nil, fmt.Errorf("ctx query one error, %s", err.Error())
			}
		} else {
			txKey := common.GetTxKeyWith(txSimContext.GetBlockProposer(), txSimContext.GetBlockHeight())
			transaction, err := txSimContext.GetBlockchainStore().GetDbTransaction(txKey)
			if err != nil {
				return nil, fmt.Errorf("ctx get db transaction error, [%s]", err.Error())
			}
			changeCurrentDB(chainId, contractName, transaction)
			row, err = transaction.QuerySingle(sql)
		}
		var dataRow map[string]string
		if row.IsEmpty() {
			dataRow = make(map[string]string, 0)
		} else {
			dataRow, err = row.Data()
			if err != nil {
				return nil, fmt.Errorf("ctx query get data to map error, %s", err.Error())
			}
		}
		ec := serialize.NewEasyCodecWithMap(dataRow)
		rsBytes := ec.Marshal()
		copy(memory[ptr:ptr+4], utils.IntToBytes(int32(len(rsBytes))))
		if len(rsBytes) == 0 {
			return nil, nil
		}
		return rsBytes, nil
	} else { // get data
		if data != nil && len(data) > 0 {
			copy(memory[ptr:ptr+int32(len(data))], data)
		}
		return nil, nil
	}
}

func (*WacsiImpl) RSHasNext(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte) error {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	rsIndex, _ := ec.GetInt32("rs_index")
	valuePtr, _ := ec.GetInt32("value_ptr")

	// get
	rows, ok := txSimContext.GetStateSqlHandle(rsIndex)
	if !ok {
		return fmt.Errorf("ctx can not found rs_index[%d]", rsIndex)
	}
	index := boolFalse
	if rows.Next() {
		index = boolTrue
	}
	copy(memory[valuePtr:valuePtr+4], utils.IntToBytes(int32(index)))
	return nil
}

func (*WacsiImpl) RSNext(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte, data []byte, isLen bool) ([]byte, error) {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	rsIndex, _ := ec.GetInt32("rs_index")
	ptr, _ := ec.GetInt32("value_ptr")

	// get handle
	rows, ok := txSimContext.GetStateSqlHandle(rsIndex)
	if !ok {
		return nil, fmt.Errorf("ctx can not found rs_index[%d]", rsIndex)
	}

	// get len
	if isLen {
		var dataRow map[string]string
		var err error
		if rows == nil {
			dataRow = make(map[string]string, 0)
		} else {
			dataRow, err = rows.Data()
			if err != nil {
				return nil, fmt.Errorf("ctx query next data error, %s", err.Error())
			}
		}
		ec := serialize.NewEasyCodecWithMap(dataRow)
		rsBytes := ec.Marshal()
		copy(memory[ptr:ptr+4], utils.IntToBytes(int32(len(rsBytes))))
		if len(rsBytes) == 0 {
			return nil, nil
		}
		return rsBytes, nil
	} else { // get data
		if len(data) > 0 {
			copy(memory[ptr:ptr+int32(len(data))], data)
		}
		return nil, nil
	}
}

func (*WacsiImpl) RSClose(requestBody []byte, txSimContext protocol.TxSimContext, memory []byte) error {
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	rsIndex, _ := ec.GetInt32("rs_index")
	valuePtr, _ := ec.GetInt32("value_ptr")

	// get
	rows, ok := txSimContext.GetStateSqlHandle(rsIndex)
	if !ok {
		return fmt.Errorf("ctx can not found rs_index[%d]", rsIndex)
	}
	var index int32 = 1
	if err := rows.Close(); err != nil {
		return fmt.Errorf("ctx close rows error, [%s]", err.Error())
	}
	copy(memory[valuePtr:valuePtr+4], utils.IntToBytes(index))
	return nil
}

func (w *WacsiImpl) ExecuteUpdate(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte, chainId string) error {
	if txSimContext.GetTx().GetHeader().TxType == common.TxType_QUERY_USER_CONTRACT {
		return fmt.Errorf(" Query transaction cannot be updated")
	}
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	sql, _ := ec.GetString("sql")
	ptr, _ := ec.GetInt32("value_ptr")

	// verify
	if err := w.verifySql.VerifyDMLSql(sql); err != nil {
		return fmt.Errorf("verify update sql error, [%s]", err.Error())
	}

	txKey := common.GetTxKeyWith(txSimContext.GetBlockProposer(), txSimContext.GetBlockHeight())
	transaction, err := txSimContext.GetBlockchainStore().GetDbTransaction(txKey)
	if err != nil {
		return fmt.Errorf("ctx get db transaction error, [%s]", err.Error())
	}

	// execute
	changeCurrentDB(chainId, contractName, transaction)
	affectedCount, err := transaction.ExecSql(sql)
	if err != nil {
		return fmt.Errorf("ctx execute update sql error, [%s], sql[%s]", err.Error(), sql)
	}
	txSimContext.PutRecord(contractName, []byte(sql))
	copy(memory[ptr:ptr+4], utils.IntToBytes(int32(affectedCount)))
	return nil
}

func (w *WacsiImpl) ExecuteDDL(requestBody []byte, contractName string, txSimContext protocol.TxSimContext, memory []byte, method string) error {
	if !w.isManageContract(method) {
		return ErrorNotManageContract
	}
	ec := serialize.NewEasyCodecWithBytes(requestBody)
	sql, _ := ec.GetString("sql")
	ptr, _ := ec.GetInt32("value_ptr")

	// verify
	if err := w.verifySql.VerifyDDLSql(sql); err != nil {
		return fmt.Errorf("verify ddl sql error,  [%s], sql[%s]", err.Error(), sql)
	}

	// execute
	if err := txSimContext.GetBlockchainStore().ExecDdlSql(contractName, sql); err != nil {
		return fmt.Errorf("ctx ExecDdlSql error, %s, sql[%s]", err.Error(), sql)
	}
	txSimContext.PutRecord(contractName, []byte(sql))
	copy(memory[ptr:ptr+4], utils.IntToBytes(0))
	return nil
}

func (w *WacsiImpl) isManageContract(method string) bool {
	return method == protocol.ContractInitMethod || method == protocol.ContractUpgradeMethod
}

func changeCurrentDB(chainId string, contractName string, transaction protocol.SqlDBTransaction) {
	dbName := statesqldb.GetContractDbName(chainId, contractName)
	//currentDbName := getCurrentDb(chainId)
	//if contractName != "" && dbName != currentDbName {
	transaction.ChangeContextDb(dbName)
	//setCurrentDb(chainId, dbName)
	//}
}
