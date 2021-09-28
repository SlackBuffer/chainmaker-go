module chainmaker.org/chainmaker-go/blockchain

go 1.15

require (
	chainmaker.org/chainmaker-go/accesscontrol v0.0.0
	chainmaker.org/chainmaker-go/consensus v0.0.0
	chainmaker.org/chainmaker-go/core v0.0.0
	chainmaker.org/chainmaker-go/net v0.0.0
	chainmaker.org/chainmaker-go/snapshot v0.0.0
	chainmaker.org/chainmaker-go/subscriber v0.0.0
	chainmaker.org/chainmaker-go/sync v0.0.0
	chainmaker.org/chainmaker-go/txpool v0.0.0
	chainmaker.org/chainmaker/chainconf/v2 v2.0.0-20210928073120-a1aa7224dd8e
	chainmaker.org/chainmaker/common/v2 v2.0.1-0.20210927025216-3d740cb6258e
	chainmaker.org/chainmaker/localconf/v2 v2.0.0-20210928020228-3ab2986d5ecd
	chainmaker.org/chainmaker/logger/v2 v2.0.0-20210907134457-53647922a89d
	chainmaker.org/chainmaker/pb-go/v2 v2.0.1-0.20210926113446-c38a67e6150e
	chainmaker.org/chainmaker/protocol/v2 v2.0.1-0.20210928034542-33bb9e319825
	chainmaker.org/chainmaker/store/v2 v2.0.0-20210927063334-95fec89a7435
	chainmaker.org/chainmaker/utils/v2 v2.0.0-20210916084713-abd13154c26b
	chainmaker.org/chainmaker/vm v0.0.0-20210918104424-239140ec3366
	chainmaker.org/chainmaker/vm-evm v0.0.0-20210916091920-b915815eb88b
	chainmaker.org/chainmaker/vm-gasm v0.0.0-20210918095814-3f0ddfe29968
	chainmaker.org/chainmaker/vm-wasmer v0.0.0-20210918173526-7cd16f1a1d3a
	chainmaker.org/chainmaker/vm-wxvm v0.0.0-20210918101823-dce1c76fb189
	github.com/mitchellh/mapstructure v1.4.1
)

replace (
	chainmaker.org/chainmaker-go/accesscontrol => ../accesscontrol
	chainmaker.org/chainmaker-go/consensus => ../consensus
	chainmaker.org/chainmaker-go/core => ../core
	chainmaker.org/chainmaker-go/monitor => ../monitor
	chainmaker.org/chainmaker-go/net => ../net
	chainmaker.org/chainmaker-go/snapshot => ../snapshot
	chainmaker.org/chainmaker-go/subscriber => ../subscriber
	chainmaker.org/chainmaker-go/sync => ../sync
	chainmaker.org/chainmaker-go/txpool => ../txpool
	github.com/libp2p/go-libp2p-core => chainmaker.org/chainmaker/libp2p-core v0.0.2
	//chainmaker.org/chainmaker-go/txpool/batchtxpool => ./../txpool/batch
	google.golang.org/grpc v1.40.0 => google.golang.org/grpc v1.26.0
)
