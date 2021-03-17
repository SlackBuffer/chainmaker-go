// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sync/sync.proto

package sync

import (
	common "chainmaker.org/chainmaker-go/pb/protogo/common"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// specific syncblockmessage types
type SyncMsg_MsgType int32

const (
	SyncMsg_NODE_STATUS_REQ  SyncMsg_MsgType = 0
	SyncMsg_NODE_STATUS_RESP SyncMsg_MsgType = 1
	SyncMsg_BLOCK_SYNC_REQ   SyncMsg_MsgType = 2
	SyncMsg_BLOCK_SYNC_RESP  SyncMsg_MsgType = 3
)

var SyncMsg_MsgType_name = map[int32]string{
	0: "NODE_STATUS_REQ",
	1: "NODE_STATUS_RESP",
	2: "BLOCK_SYNC_REQ",
	3: "BLOCK_SYNC_RESP",
}

var SyncMsg_MsgType_value = map[string]int32{
	"NODE_STATUS_REQ":  0,
	"NODE_STATUS_RESP": 1,
	"BLOCK_SYNC_REQ":   2,
	"BLOCK_SYNC_RESP":  3,
}

func (x SyncMsg_MsgType) String() string {
	return proto.EnumName(SyncMsg_MsgType_name, int32(x))
}

func (SyncMsg_MsgType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_221a5c59bc60326f, []int{0, 0}
}

// network message of synchronization module
type SyncMsg struct {
	// sync message type
	Type SyncMsg_MsgType `protobuf:"varint,1,opt,name=type,proto3,enum=sync.SyncMsg_MsgType" json:"type,omitempty"`
	// payload for the message
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *SyncMsg) Reset()         { *m = SyncMsg{} }
func (m *SyncMsg) String() string { return proto.CompactTextString(m) }
func (*SyncMsg) ProtoMessage()    {}
func (*SyncMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_221a5c59bc60326f, []int{0}
}
func (m *SyncMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SyncMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SyncMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SyncMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncMsg.Merge(m, src)
}
func (m *SyncMsg) XXX_Size() int {
	return m.Size()
}
func (m *SyncMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncMsg.DiscardUnknown(m)
}

var xxx_messageInfo_SyncMsg proto.InternalMessageInfo

func (m *SyncMsg) GetType() SyncMsg_MsgType {
	if m != nil {
		return m.Type
	}
	return SyncMsg_NODE_STATUS_REQ
}

func (m *SyncMsg) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

// response message for node status
type BlockHeightBCM struct {
	BlockHeight int64 `protobuf:"varint,1,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
}

func (m *BlockHeightBCM) Reset()         { *m = BlockHeightBCM{} }
func (m *BlockHeightBCM) String() string { return proto.CompactTextString(m) }
func (*BlockHeightBCM) ProtoMessage()    {}
func (*BlockHeightBCM) Descriptor() ([]byte, []int) {
	return fileDescriptor_221a5c59bc60326f, []int{1}
}
func (m *BlockHeightBCM) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BlockHeightBCM) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BlockHeightBCM.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BlockHeightBCM) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockHeightBCM.Merge(m, src)
}
func (m *BlockHeightBCM) XXX_Size() int {
	return m.Size()
}
func (m *BlockHeightBCM) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockHeightBCM.DiscardUnknown(m)
}

var xxx_messageInfo_BlockHeightBCM proto.InternalMessageInfo

func (m *BlockHeightBCM) GetBlockHeight() int64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

// block request message
type BlockSyncReq struct {
	BlockHeight int64 `protobuf:"varint,1,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	BatchSize   int64 `protobuf:"varint,2,opt,name=batchSize,proto3" json:"batchSize,omitempty"`
	WithRwset   bool  `protobuf:"varint,3,opt,name=with_rwset,json=withRwset,proto3" json:"with_rwset,omitempty"`
}

func (m *BlockSyncReq) Reset()         { *m = BlockSyncReq{} }
func (m *BlockSyncReq) String() string { return proto.CompactTextString(m) }
func (*BlockSyncReq) ProtoMessage()    {}
func (*BlockSyncReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_221a5c59bc60326f, []int{2}
}
func (m *BlockSyncReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BlockSyncReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BlockSyncReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BlockSyncReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockSyncReq.Merge(m, src)
}
func (m *BlockSyncReq) XXX_Size() int {
	return m.Size()
}
func (m *BlockSyncReq) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockSyncReq.DiscardUnknown(m)
}

var xxx_messageInfo_BlockSyncReq proto.InternalMessageInfo

func (m *BlockSyncReq) GetBlockHeight() int64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *BlockSyncReq) GetBatchSize() int64 {
	if m != nil {
		return m.BatchSize
	}
	return 0
}

func (m *BlockSyncReq) GetWithRwset() bool {
	if m != nil {
		return m.WithRwset
	}
	return false
}

// batch blocks
type BlockBatch struct {
	Batchs []*common.Block `protobuf:"bytes,1,rep,name=batchs,proto3" json:"batchs,omitempty"`
}

func (m *BlockBatch) Reset()         { *m = BlockBatch{} }
func (m *BlockBatch) String() string { return proto.CompactTextString(m) }
func (*BlockBatch) ProtoMessage()    {}
func (*BlockBatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_221a5c59bc60326f, []int{3}
}
func (m *BlockBatch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BlockBatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BlockBatch.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BlockBatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockBatch.Merge(m, src)
}
func (m *BlockBatch) XXX_Size() int {
	return m.Size()
}
func (m *BlockBatch) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockBatch.DiscardUnknown(m)
}

var xxx_messageInfo_BlockBatch proto.InternalMessageInfo

func (m *BlockBatch) GetBatchs() []*common.Block {
	if m != nil {
		return m.Batchs
	}
	return nil
}

// information of batch blocks
type BlockInfoBatch struct {
	Batch []*common.BlockInfo `protobuf:"bytes,1,rep,name=batch,proto3" json:"batch,omitempty"`
}

func (m *BlockInfoBatch) Reset()         { *m = BlockInfoBatch{} }
func (m *BlockInfoBatch) String() string { return proto.CompactTextString(m) }
func (*BlockInfoBatch) ProtoMessage()    {}
func (*BlockInfoBatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_221a5c59bc60326f, []int{4}
}
func (m *BlockInfoBatch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BlockInfoBatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BlockInfoBatch.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BlockInfoBatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockInfoBatch.Merge(m, src)
}
func (m *BlockInfoBatch) XXX_Size() int {
	return m.Size()
}
func (m *BlockInfoBatch) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockInfoBatch.DiscardUnknown(m)
}

var xxx_messageInfo_BlockInfoBatch proto.InternalMessageInfo

func (m *BlockInfoBatch) GetBatch() []*common.BlockInfo {
	if m != nil {
		return m.Batch
	}
	return nil
}

// block response message
type SyncBlockBatch struct {
	// Types that are valid to be assigned to Data:
	//	*SyncBlockBatch_BlockBatch
	//	*SyncBlockBatch_BlockinfoBatch
	Data isSyncBlockBatch_Data `protobuf_oneof:"Data"`
}

func (m *SyncBlockBatch) Reset()         { *m = SyncBlockBatch{} }
func (m *SyncBlockBatch) String() string { return proto.CompactTextString(m) }
func (*SyncBlockBatch) ProtoMessage()    {}
func (*SyncBlockBatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_221a5c59bc60326f, []int{5}
}
func (m *SyncBlockBatch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SyncBlockBatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SyncBlockBatch.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SyncBlockBatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncBlockBatch.Merge(m, src)
}
func (m *SyncBlockBatch) XXX_Size() int {
	return m.Size()
}
func (m *SyncBlockBatch) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncBlockBatch.DiscardUnknown(m)
}

var xxx_messageInfo_SyncBlockBatch proto.InternalMessageInfo

type isSyncBlockBatch_Data interface {
	isSyncBlockBatch_Data()
	MarshalTo([]byte) (int, error)
	Size() int
}

type SyncBlockBatch_BlockBatch struct {
	BlockBatch *BlockBatch `protobuf:"bytes,1,opt,name=block_batch,json=blockBatch,proto3,oneof" json:"block_batch,omitempty"`
}
type SyncBlockBatch_BlockinfoBatch struct {
	BlockinfoBatch *BlockInfoBatch `protobuf:"bytes,2,opt,name=blockinfo_batch,json=blockinfoBatch,proto3,oneof" json:"blockinfo_batch,omitempty"`
}

func (*SyncBlockBatch_BlockBatch) isSyncBlockBatch_Data()     {}
func (*SyncBlockBatch_BlockinfoBatch) isSyncBlockBatch_Data() {}

func (m *SyncBlockBatch) GetData() isSyncBlockBatch_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SyncBlockBatch) GetBlockBatch() *BlockBatch {
	if x, ok := m.GetData().(*SyncBlockBatch_BlockBatch); ok {
		return x.BlockBatch
	}
	return nil
}

func (m *SyncBlockBatch) GetBlockinfoBatch() *BlockInfoBatch {
	if x, ok := m.GetData().(*SyncBlockBatch_BlockinfoBatch); ok {
		return x.BlockinfoBatch
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*SyncBlockBatch) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*SyncBlockBatch_BlockBatch)(nil),
		(*SyncBlockBatch_BlockinfoBatch)(nil),
	}
}

func init() {
	proto.RegisterEnum("sync.SyncMsg_MsgType", SyncMsg_MsgType_name, SyncMsg_MsgType_value)
	proto.RegisterType((*SyncMsg)(nil), "sync.SyncMsg")
	proto.RegisterType((*BlockHeightBCM)(nil), "sync.BlockHeightBCM")
	proto.RegisterType((*BlockSyncReq)(nil), "sync.BlockSyncReq")
	proto.RegisterType((*BlockBatch)(nil), "sync.BlockBatch")
	proto.RegisterType((*BlockInfoBatch)(nil), "sync.BlockInfoBatch")
	proto.RegisterType((*SyncBlockBatch)(nil), "sync.SyncBlockBatch")
}

func init() { proto.RegisterFile("sync/sync.proto", fileDescriptor_221a5c59bc60326f) }

var fileDescriptor_221a5c59bc60326f = []byte{
	// 453 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xc1, 0x8e, 0x93, 0x40,
	0x1c, 0xc6, 0x99, 0x52, 0x5b, 0xf7, 0xdf, 0x4a, 0x71, 0x5c, 0x13, 0x62, 0x94, 0x54, 0x12, 0x23,
	0x26, 0x0a, 0x49, 0x7b, 0xf2, 0x64, 0xa4, 0xbb, 0xa6, 0x46, 0xbb, 0xbb, 0x0e, 0xf5, 0xa0, 0x89,
	0x21, 0x80, 0x2c, 0x90, 0xdd, 0x32, 0x58, 0x48, 0x36, 0xf8, 0x10, 0xc6, 0x17, 0xf1, 0x3d, 0x3c,
	0xee, 0xd1, 0xa3, 0x69, 0x5f, 0xc4, 0xcc, 0x1f, 0x76, 0xd9, 0xde, 0xbc, 0x10, 0xfe, 0xdf, 0xf7,
	0xfd, 0x66, 0xe6, 0xcb, 0x0c, 0x8c, 0x8a, 0x2a, 0x0b, 0x6d, 0xf1, 0xb1, 0xf2, 0x35, 0x2f, 0x39,
	0xed, 0x8a, 0xff, 0x07, 0x34, 0xe4, 0xab, 0x15, 0xcf, 0xec, 0xe0, 0x9c, 0x87, 0x67, 0xb5, 0x63,
	0xfc, 0x22, 0xd0, 0x77, 0xab, 0x2c, 0x5c, 0x14, 0x31, 0x7d, 0x06, 0xdd, 0xb2, 0xca, 0x23, 0x8d,
	0x8c, 0x89, 0xa9, 0x4c, 0xee, 0x5b, 0xb8, 0x40, 0x63, 0x5a, 0x8b, 0x22, 0x5e, 0x56, 0x79, 0xc4,
	0x30, 0x42, 0x35, 0xe8, 0xe7, 0x7e, 0x75, 0xce, 0xfd, 0xaf, 0x5a, 0x67, 0x4c, 0xcc, 0x21, 0xbb,
	0x1a, 0x8d, 0x2f, 0xd0, 0x6f, 0xa2, 0xf4, 0x1e, 0x8c, 0x8e, 0x8e, 0x0f, 0x0e, 0x3d, 0x77, 0xf9,
	0x7a, 0xf9, 0xd1, 0xf5, 0xd8, 0xe1, 0x07, 0x55, 0xa2, 0xfb, 0xa0, 0xee, 0x8a, 0xee, 0x89, 0x4a,
	0x28, 0x05, 0xc5, 0x79, 0x7f, 0x3c, 0x7b, 0xe7, 0xb9, 0x9f, 0x8e, 0x66, 0x98, 0xec, 0x08, 0x7c,
	0x47, 0x73, 0x4f, 0x54, 0xd9, 0x98, 0x82, 0xe2, 0x88, 0xe3, 0xcf, 0xa3, 0x34, 0x4e, 0x4a, 0x67,
	0xb6, 0xa0, 0x8f, 0x61, 0x88, 0x85, 0xbc, 0x04, 0x25, 0x3c, 0xbd, 0xcc, 0x06, 0x41, 0x9b, 0x32,
	0x32, 0x18, 0x22, 0x24, 0xba, 0xb0, 0xe8, 0xdb, 0x7f, 0x20, 0xf4, 0x21, 0xec, 0x05, 0x7e, 0x19,
	0x26, 0x6e, 0xfa, 0x3d, 0xc2, 0x8a, 0x32, 0x6b, 0x05, 0xfa, 0x08, 0xe0, 0x22, 0x2d, 0x13, 0x6f,
	0x7d, 0x51, 0x44, 0xa5, 0x26, 0x8f, 0x89, 0x79, 0x9b, 0xed, 0x09, 0x85, 0x09, 0xc1, 0x98, 0x02,
	0xe0, 0x7e, 0x8e, 0x00, 0xe8, 0x13, 0xe8, 0x21, 0x59, 0x68, 0x64, 0x2c, 0x9b, 0x83, 0xc9, 0x1d,
	0xab, 0xbe, 0x07, 0x0b, 0x33, 0xac, 0x31, 0x8d, 0x97, 0x4d, 0xb3, 0xb7, 0xd9, 0x29, 0xaf, 0xc1,
	0xa7, 0x70, 0x0b, 0xbd, 0x86, 0xbb, 0xbb, 0xc3, 0x89, 0x18, 0xab, 0x7d, 0xe3, 0x07, 0x01, 0x45,
	0x74, 0xbb, 0xb1, 0xe9, 0x14, 0xea, 0x3a, 0xde, 0xd5, 0x0a, 0xc4, 0x1c, 0x4c, 0xd4, 0xfa, 0x4a,
	0xdb, 0xd8, 0x5c, 0x62, 0x10, 0xb4, 0xd0, 0x2b, 0x18, 0xe1, 0x94, 0x66, 0xa7, 0xbc, 0x01, 0x3b,
	0x08, 0xee, 0xdf, 0x00, 0xaf, 0xcf, 0x37, 0x97, 0x98, 0x72, 0x1d, 0x47, 0xc5, 0xe9, 0x41, 0xf7,
	0xc0, 0x2f, 0x7d, 0xe7, 0xcd, 0xef, 0x8d, 0x4e, 0x2e, 0x37, 0x3a, 0xf9, 0xbb, 0xd1, 0xc9, 0xcf,
	0xad, 0x2e, 0x5d, 0x6e, 0x75, 0xe9, 0xcf, 0x56, 0x97, 0x3e, 0x3f, 0x0f, 0x13, 0x3f, 0xcd, 0x56,
	0xfe, 0x59, 0xb4, 0xb6, 0xf8, 0x3a, 0xb6, 0xdb, 0xf1, 0x45, 0xcc, 0xed, 0x3c, 0xb0, 0xf1, 0x59,
	0xc6, 0x1c, 0x5f, 0x6f, 0xd0, 0xc3, 0x69, 0xfa, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x14, 0x37, 0xbe,
	0x90, 0xd1, 0x02, 0x00, 0x00,
}

func (m *SyncMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SyncMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SyncMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Payload) > 0 {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintSync(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0x12
	}
	if m.Type != 0 {
		i = encodeVarintSync(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *BlockHeightBCM) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BlockHeightBCM) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BlockHeightBCM) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.BlockHeight != 0 {
		i = encodeVarintSync(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *BlockSyncReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BlockSyncReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BlockSyncReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.WithRwset {
		i--
		if m.WithRwset {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.BatchSize != 0 {
		i = encodeVarintSync(dAtA, i, uint64(m.BatchSize))
		i--
		dAtA[i] = 0x10
	}
	if m.BlockHeight != 0 {
		i = encodeVarintSync(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *BlockBatch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BlockBatch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BlockBatch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Batchs) > 0 {
		for iNdEx := len(m.Batchs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Batchs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSync(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *BlockInfoBatch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BlockInfoBatch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BlockInfoBatch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Batch) > 0 {
		for iNdEx := len(m.Batch) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Batch[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSync(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *SyncBlockBatch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SyncBlockBatch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SyncBlockBatch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Data != nil {
		{
			size := m.Data.Size()
			i -= size
			if _, err := m.Data.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *SyncBlockBatch_BlockBatch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SyncBlockBatch_BlockBatch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.BlockBatch != nil {
		{
			size, err := m.BlockBatch.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSync(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *SyncBlockBatch_BlockinfoBatch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SyncBlockBatch_BlockinfoBatch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.BlockinfoBatch != nil {
		{
			size, err := m.BlockinfoBatch.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSync(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func encodeVarintSync(dAtA []byte, offset int, v uint64) int {
	offset -= sovSync(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SyncMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovSync(uint64(m.Type))
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovSync(uint64(l))
	}
	return n
}

func (m *BlockHeightBCM) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BlockHeight != 0 {
		n += 1 + sovSync(uint64(m.BlockHeight))
	}
	return n
}

func (m *BlockSyncReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BlockHeight != 0 {
		n += 1 + sovSync(uint64(m.BlockHeight))
	}
	if m.BatchSize != 0 {
		n += 1 + sovSync(uint64(m.BatchSize))
	}
	if m.WithRwset {
		n += 2
	}
	return n
}

func (m *BlockBatch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Batchs) > 0 {
		for _, e := range m.Batchs {
			l = e.Size()
			n += 1 + l + sovSync(uint64(l))
		}
	}
	return n
}

func (m *BlockInfoBatch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Batch) > 0 {
		for _, e := range m.Batch {
			l = e.Size()
			n += 1 + l + sovSync(uint64(l))
		}
	}
	return n
}

func (m *SyncBlockBatch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Data != nil {
		n += m.Data.Size()
	}
	return n
}

func (m *SyncBlockBatch_BlockBatch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BlockBatch != nil {
		l = m.BlockBatch.Size()
		n += 1 + l + sovSync(uint64(l))
	}
	return n
}
func (m *SyncBlockBatch_BlockinfoBatch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BlockinfoBatch != nil {
		l = m.BlockinfoBatch.Size()
		n += 1 + l + sovSync(uint64(l))
	}
	return n
}

func sovSync(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSync(x uint64) (n int) {
	return sovSync(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SyncMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSync
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SyncMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SyncMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= SyncMsg_MsgType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSync
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSync(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSync
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BlockHeightBCM) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSync
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BlockHeightBCM: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BlockHeightBCM: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSync(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSync
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BlockSyncReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSync
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BlockSyncReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BlockSyncReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BatchSize", wireType)
			}
			m.BatchSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BatchSize |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field WithRwset", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.WithRwset = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipSync(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSync
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BlockBatch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSync
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BlockBatch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BlockBatch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Batchs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSync
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Batchs = append(m.Batchs, &common.Block{})
			if err := m.Batchs[len(m.Batchs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSync(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSync
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BlockInfoBatch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSync
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BlockInfoBatch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BlockInfoBatch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Batch", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSync
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Batch = append(m.Batch, &common.BlockInfo{})
			if err := m.Batch[len(m.Batch)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSync(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSync
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SyncBlockBatch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSync
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SyncBlockBatch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SyncBlockBatch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockBatch", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSync
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &BlockBatch{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Data = &SyncBlockBatch_BlockBatch{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockinfoBatch", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSync
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSync
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSync
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &BlockInfoBatch{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Data = &SyncBlockBatch_BlockinfoBatch{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSync(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSync
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSync(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSync
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSync
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSync
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthSync
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSync
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSync
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSync        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSync          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSync = fmt.Errorf("proto: unexpected end of group")
)