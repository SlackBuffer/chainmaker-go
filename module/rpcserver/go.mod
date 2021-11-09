module chainmaker.org/chainmaker-go/rpcserver

go 1.15

require (
	chainmaker.org/chainmaker-go/blockchain v0.0.0-00010101000000-000000000000
	chainmaker.org/chainmaker-go/subscriber v0.0.0
	chainmaker.org/chainmaker/common/v2 v2.1.0
	chainmaker.org/chainmaker/localconf/v2 v2.0.0-20211101111610-0d268248b5c8
	chainmaker.org/chainmaker/logger/v2 v2.1.0
	chainmaker.org/chainmaker/pb-go/v2 v2.1.0
	chainmaker.org/chainmaker/protocol/v2 v2.0.1-0.20211108075639-576c31f03396
	chainmaker.org/chainmaker/store/v2 v2.0.1-0.20211101060552-ab2318ad9c6e
	chainmaker.org/chainmaker/utils/v2 v2.0.0-20211108092352-2a3335a4ba15
	chainmaker.org/chainmaker/vm-native v0.0.0-20211029075817-046b5f1bbc31
	github.com/gogo/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/prometheus/client_golang v1.11.0
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	google.golang.org/grpc v1.41.0
)

replace (
	chainmaker.org/chainmaker-go/accesscontrol => ../accesscontrol
	chainmaker.org/chainmaker-go/blockchain => ../blockchain
	chainmaker.org/chainmaker-go/net => ../net
	chainmaker.org/chainmaker-go/snapshot => ../snapshot
	chainmaker.org/chainmaker-go/subscriber => ../subscriber
	chainmaker.org/chainmaker-go/sync => ../sync
	chainmaker.org/chainmaker-go/txpool => ../txpool
	chainmaker.org/chainmaker-go/vm => ../vm
	github.com/libp2p/go-libp2p-core => chainmaker.org/chainmaker/libp2p-core v0.0.2
	google.golang.org/grpc v1.40.0 => google.golang.org/grpc v1.26.0
)
