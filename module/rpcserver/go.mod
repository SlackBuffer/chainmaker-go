module chainmaker.org/chainmaker-go/rpcserver

go 1.15

require (
	chainmaker.org/chainmaker-go/blockchain v0.0.0-00010101000000-000000000000
	chainmaker.org/chainmaker-go/subscriber v0.0.0
	chainmaker.org/chainmaker/common/v2 v2.0.1-0.20211025091405-b52d2356a39e
	chainmaker.org/chainmaker/localconf/v2 v2.0.0-20211014134424-9431ffcc5bbc
	chainmaker.org/chainmaker/logger/v2 v2.0.1-0.20211015125919-8e5199930ac9
	chainmaker.org/chainmaker/pb-go/v2 v2.0.1-0.20211021024710-9329804d1c21
	chainmaker.org/chainmaker/protocol/v2 v2.0.1-0.20211014144951-97323532a236
	chainmaker.org/chainmaker/store/v2 v2.0.1-0.20211027072818-556aa7111997
	chainmaker.org/chainmaker/utils/v2 v2.0.0-20211027124954-09b710bd9ce8
	chainmaker.org/chainmaker/vm-native v0.0.0-20211027120103-3b4f11ba02ba
	github.com/gogo/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/prometheus/client_golang v1.11.0
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	google.golang.org/grpc v1.41.0
)

replace (
	chainmaker.org/chainmaker-go/accesscontrol => ../accesscontrol
	chainmaker.org/chainmaker-go/blockchain => ../blockchain
	chainmaker.org/chainmaker-go/consensus => ../consensus
	chainmaker.org/chainmaker-go/core => ../core
	chainmaker.org/chainmaker-go/net => ../net
	chainmaker.org/chainmaker-go/snapshot => ../snapshot
	chainmaker.org/chainmaker-go/subscriber => ../subscriber
	chainmaker.org/chainmaker-go/sync => ../sync
	chainmaker.org/chainmaker-go/txpool => ../txpool
	chainmaker.org/chainmaker-go/vm => ../vm
	github.com/libp2p/go-libp2p-core => chainmaker.org/chainmaker/libp2p-core v0.0.2
	google.golang.org/grpc v1.40.0 => google.golang.org/grpc v1.26.0
)
