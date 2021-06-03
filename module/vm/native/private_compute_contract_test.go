package native

import (
	"chainmaker.org/chainmaker-go/logger"
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var certCaFilename = "testdata/remote_attestation/enclave_cacert.crt"
var certFilename = "testdata/remote_attestation/in_teecert.pem"
var proofFilename = "testdata/remote_attestation/proof.hex"
var reportFilename = "testdata/remote_attestation/report.dat"

func readFileData(filename string, t *testing.T) []byte {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("open file '%s' error: %v", certCaFilename, err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("read file '%v' error: %v", certCaFilename, err)
	}

	return data
}

func TestSaveEnclaveCaCert(t *testing.T) {
	ds := map[string][]byte{}
	mockCtx := newTxContextMock(ds)

	privateComputeRuntime := PrivateComputeRuntime{
		log: logger.GetLogger("test-logger"),
	}

	// 读取CA证书
	certCaPem := readFileData(certCaFilename, t)

	params := map[string]string{}
	params["ca_cert"] = string(certCaPem)
	result, err := privateComputeRuntime.SaveEnclaveCACert(mockCtx, params)
	if err != nil {
		t.Fatalf("Save enclave ca cert error: %v", err)
	}

	fmt.Printf("result = %v \n", string(result))
}

func TestSaveEnclaveReport(t *testing.T) {
	ds := map[string][]byte{}
	mockCtx := newTxContextMock(ds)

	privateComputeRuntime := PrivateComputeRuntime{
		log: logger.GetLogger("test-logger"),
	}

	// 读取report
	report := readFileData(reportFilename, t)

	params := map[string]string{}
	params["enclave_id"] = "global_enclave_id"
	params["report"] = string(report)
	_, err := privateComputeRuntime.SaveEnclaveReport(mockCtx, params)
	if err != nil {
		t.Fatalf("Save enclave report error: %v", err)
	}

	for key, val := range ds {
		fmt.Printf("%s ==>\n%s \n", key, val)
	}
}

func TestSaveRemoteAttestation(t *testing.T) {

	caCertPem := readFileData(certCaFilename, t)
	report := readFileData("testdata/remote_attestation/report.dat", t)

	ds := map[string][]byte{}
	mockCtx := newTxContextMock(ds)
	reportKey := commonPb.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String() + "global_enclave_id::report"
	ds[reportKey] = report
	caCertKey := commonPb.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String() + "::ca_cert"
	ds[caCertKey] = caCertPem

	proofFile, err := os.Open(proofFilename)
	if err != nil {
		t.Fatalf("open file '%s' error: %v", proofFilename, err)
	}

	proofHex, err := ioutil.ReadAll(proofFile)
	if err != nil {
		t.Fatalf("read file '%v' error: %v", proofFile, err)
	}

	proof, err := hex.DecodeString(string(proofHex))
	if err != nil {
		t.Fatalf("decode hex string error: %v", err)
	}

	privateComputeRuntime := PrivateComputeRuntime{
		log: logger.GetLogger("test-logger"),
	}
	params := map[string]string{}
	params["proof"] = string(proof)
	result, err := privateComputeRuntime.SaveRemoteAttestation(mockCtx, params)
	if err != nil {
		t.Fatalf("Save remote attestation error: %v", err)
	}

	fmt.Printf("result = %v \n", string(result));
	for key, val := range ds {
		fmt.Printf("key = %v, val = %x \n", key, val)
	}
}

func TestInTeecertPemFile(t *testing.T) {
	file, err := os.Open("/Users/caizhihong/证书测试/in_teecert.pem")
	if err != nil {
		t.Fatalf("open file error: %v", err)
	}

	crtPEM, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("read file error: %v", err)
	}

	signingPubkey, cryptoPubkey, err := getPubkeyPairFromCert(crtPEM)
	if err != nil {
		t.Fatalf("get pubkey pair error: %v", err)
	}
	fmt.Printf("signing pub key ==> %v \n", signingPubkey.ToStandardKey())
	fmt.Printf("crypto pub key ==> %v \n", cryptoPubkey.ToStandardKey())
}