module chainmaker.org/chainmaker-go

go 1.15

require (
	chainmaker.org/chainmaker-go/accesscontrol v0.0.0
	chainmaker.org/chainmaker-go/blockchain v0.0.0
	chainmaker.org/chainmaker-go/net v0.0.0
	chainmaker.org/chainmaker-go/rpcserver v0.0.0
	chainmaker.org/chainmaker-go/txpool v0.0.0
	chainmaker.org/chainmaker-go/vm v0.0.0
	chainmaker.org/chainmaker/common/v2 v2.1.1-0.20211125080705-8c4f85e6ca19
	chainmaker.org/chainmaker/localconf/v2 v2.1.0
	chainmaker.org/chainmaker/logger/v2 v2.1.1-0.20211109074349-f79af5e1892d
	chainmaker.org/chainmaker/pb-go/v2 v2.1.1-0.20211130072802-e9d8fae57083
	chainmaker.org/chainmaker/protocol/v2 v2.1.1-0.20211130023410-5df7eb63bfb8
	chainmaker.org/chainmaker/sdk-go/v2 v2.0.1-0.20211126034817-c7a1903e2f83
	chainmaker.org/chainmaker/txpool-batch/v2 v2.1.0
	chainmaker.org/chainmaker/txpool-single/v2 v2.1.0
	chainmaker.org/chainmaker/utils/v2 v2.1.1-0.20211109074701-81d58330e787
	chainmaker.org/chainmaker/vm-docker-go v0.0.0-20211122064347-23cc0a21ba97
	chainmaker.org/chainmaker/vm-evm/v2 v2.0.0-20211130025150-83cc8c18e71b
	chainmaker.org/chainmaker/vm-gasm/v2 v2.1.1-0.20211117071021-e940cf22a033
	chainmaker.org/chainmaker/vm-wasmer/v2 v2.1.1-0.20211117071057-eee2c7f08b39
	chainmaker.org/chainmaker/vm-wxvm/v2 v2.1.1-0.20211117041254-a6cd6aecd988
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	github.com/common-nighthawk/go-figure v0.0.0-20200609044655-c4b36f998cf2
	github.com/ethereum/go-ethereum v1.10.4
	github.com/gogo/protobuf v1.3.2
	github.com/mr-tron/base58 v1.2.0
	github.com/prometheus/client_golang v1.11.0
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tjfoc/gmtls v1.2.1 // indirect
	google.golang.org/grpc v1.41.0
)

replace (
	chainmaker.org/chainmaker-go/accesscontrol => ./module/accesscontrol
	chainmaker.org/chainmaker-go/blockchain => ./module/blockchain
	chainmaker.org/chainmaker-go/consensus => ./module/consensus
	chainmaker.org/chainmaker-go/core => ./module/core
	chainmaker.org/chainmaker-go/net => ./module/net
	chainmaker.org/chainmaker-go/rpcserver => ./module/rpcserver
	chainmaker.org/chainmaker-go/snapshot => ./module/snapshot
	chainmaker.org/chainmaker-go/subscriber => ./module/subscriber
	chainmaker.org/chainmaker-go/sync => ./module/sync
	chainmaker.org/chainmaker-go/txpool => ./module/txpool
	chainmaker.org/chainmaker-go/vm => ./module/vm
	github.com/libp2p/go-libp2p-core => chainmaker.org/chainmaker/libp2p-core v1.0.0
	github.com/spf13/afero => github.com/spf13/afero v1.5.1 //for go1.15 build
	github.com/spf13/viper => github.com/spf13/viper v1.7.1 //for go1.15 build
	google.golang.org/grpc v1.40.0 => google.golang.org/grpc v1.26.0
)
