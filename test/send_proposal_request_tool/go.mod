module chainmaker.org/chainmaker-go/test/send_proposal_request_tool

go 1.15

require (
	chainmaker.org/chainmaker-go/accesscontrol v0.0.0
	chainmaker.org/chainmaker/common/v2 v2.0.1-0.20211111105523-441fed4a8603
	chainmaker.org/chainmaker/logger/v2 v2.0.1-0.20211014131951-892d098049bc
	chainmaker.org/chainmaker/pb-go/v2 v2.0.1-0.20211021024710-9329804d1c21
	chainmaker.org/chainmaker/protocol/v2 v2.0.1-0.20211014144951-97323532a236
	chainmaker.org/chainmaker/utils/v2 v2.0.0-20211027124954-09b710bd9ce8
	github.com/Rican7/retry v0.1.0
	github.com/ethereum/go-ethereum v1.10.2
	github.com/gogo/protobuf v1.3.2
	github.com/mr-tron/base58 v1.2.0
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/pretty v1.2.0
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98 // indirect
	google.golang.org/grpc v1.37.0
)

replace (
	chainmaker.org/chainmaker-go/accesscontrol => ../../module/accesscontrol

	chainmaker.org/chainmaker-go/localconf => ../../module/conf/localconf
	chainmaker.org/chainmaker-go/logger => ../../module/logger

)
