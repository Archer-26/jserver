// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        (unknown)
// source: island.proto

package message

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

//岛屿数据
//[Inject(IslandMeta)]
//[Sync(SyncIsland)]
type Island struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           int32           `protobuf:"varint,1,opt,name=Id,json=id,proto3" json:"Id"`                                                                                                                     //岛屿Id
	CheckpointId int32           `protobuf:"varint,2,opt,name=CheckpointId,json=checkpointId,proto3" json:"CheckpointId"`                                                                                       //关卡Id
	BuildingMap  map[int32]int32 `protobuf:"bytes,3,rep,name=BuildingMap,json=buildingMap,proto3" json:"BuildingMap" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"` //岛上的建筑<槽位，建筑ID>
	JsonInfo     string          `protobuf:"bytes,4,opt,name=JsonInfo,json=jsonInfo,proto3" json:"JsonInfo"`                                                                                                    //
}

func (x *Island) Reset() {
	*x = Island{}
	if protoimpl.UnsafeEnabled {
		mi := &file_island_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Island) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Island) ProtoMessage() {}

func (x *Island) ProtoReflect() protoreflect.Message {
	mi := &file_island_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Island.ProtoReflect.Descriptor instead.
func (*Island) Descriptor() ([]byte, []int) {
	return file_island_proto_rawDescGZIP(), []int{0}
}

func (x *Island) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Island) GetCheckpointId() int32 {
	if x != nil {
		return x.CheckpointId
	}
	return 0
}

func (x *Island) GetBuildingMap() map[int32]int32 {
	if x != nil {
		return x.BuildingMap
	}
	return nil
}

func (x *Island) GetJsonInfo() string {
	if x != nil {
		return x.JsonInfo
	}
	return ""
}

// 同步岛屿信息
type NotifyIslandInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//[Sync(Root.Role.Islands, Update, v => v.Id)]
	Islands []*Island `protobuf:"bytes,1,rep,name=Islands,json=islands,proto3" json:"Islands"`
}

func (x *NotifyIslandInfo) Reset() {
	*x = NotifyIslandInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_island_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyIslandInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyIslandInfo) ProtoMessage() {}

func (x *NotifyIslandInfo) ProtoReflect() protoreflect.Message {
	mi := &file_island_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyIslandInfo.ProtoReflect.Descriptor instead.
func (*NotifyIslandInfo) Descriptor() ([]byte, []int) {
	return file_island_proto_rawDescGZIP(), []int{1}
}

func (x *NotifyIslandInfo) GetIslands() []*Island {
	if x != nil {
		return x.Islands
	}
	return nil
}

//放置建筑--需要推送岛屿信息
type PlaceBuildingReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Land *Island `protobuf:"bytes,1,opt,name=land,proto3" json:"land"`
}

func (x *PlaceBuildingReq) Reset() {
	*x = PlaceBuildingReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_island_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlaceBuildingReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlaceBuildingReq) ProtoMessage() {}

func (x *PlaceBuildingReq) ProtoReflect() protoreflect.Message {
	mi := &file_island_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlaceBuildingReq.ProtoReflect.Descriptor instead.
func (*PlaceBuildingReq) Descriptor() ([]byte, []int) {
	return file_island_proto_rawDescGZIP(), []int{2}
}

func (x *PlaceBuildingReq) GetLand() *Island {
	if x != nil {
		return x.Land
	}
	return nil
}

type PlaceBuildingRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result MSG_RESULT `protobuf:"varint,1,opt,name=Result,json=result,proto3,enum=message.MSG_RESULT" json:"Result"` //成功码
}

func (x *PlaceBuildingRes) Reset() {
	*x = PlaceBuildingRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_island_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlaceBuildingRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlaceBuildingRes) ProtoMessage() {}

func (x *PlaceBuildingRes) ProtoReflect() protoreflect.Message {
	mi := &file_island_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlaceBuildingRes.ProtoReflect.Descriptor instead.
func (*PlaceBuildingRes) Descriptor() ([]byte, []int) {
	return file_island_proto_rawDescGZIP(), []int{3}
}

func (x *PlaceBuildingRes) GetResult() MSG_RESULT {
	if x != nil {
		return x.Result
	}
	return MSG_RESULT_SUCCESS
}

//关卡完成
type CheckpointFinishReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Gold         int64 `protobuf:"varint,1,opt,name=Gold,json=gold,proto3" json:"Gold"`                         // 关卡获取金币
	IslandId     int64 `protobuf:"varint,2,opt,name=IslandId,json=islandId,proto3" json:"IslandId"`             // 岛屿Id
	CheckpointId int64 `protobuf:"varint,3,opt,name=CheckpointId,json=checkpointId,proto3" json:"CheckpointId"` // 关卡Id
	Success      int32 `protobuf:"varint,4,opt,name=Success,json=success,proto3" json:"Success"`                // 0 通关. 1 为通关
}

func (x *CheckpointFinishReq) Reset() {
	*x = CheckpointFinishReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_island_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckpointFinishReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckpointFinishReq) ProtoMessage() {}

func (x *CheckpointFinishReq) ProtoReflect() protoreflect.Message {
	mi := &file_island_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckpointFinishReq.ProtoReflect.Descriptor instead.
func (*CheckpointFinishReq) Descriptor() ([]byte, []int) {
	return file_island_proto_rawDescGZIP(), []int{4}
}

func (x *CheckpointFinishReq) GetGold() int64 {
	if x != nil {
		return x.Gold
	}
	return 0
}

func (x *CheckpointFinishReq) GetIslandId() int64 {
	if x != nil {
		return x.IslandId
	}
	return 0
}

func (x *CheckpointFinishReq) GetCheckpointId() int64 {
	if x != nil {
		return x.CheckpointId
	}
	return 0
}

func (x *CheckpointFinishReq) GetSuccess() int32 {
	if x != nil {
		return x.Success
	}
	return 0
}

type CheckpointFinishRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result MSG_RESULT `protobuf:"varint,1,opt,name=Result,json=result,proto3,enum=message.MSG_RESULT" json:"Result"` //成功码
}

func (x *CheckpointFinishRes) Reset() {
	*x = CheckpointFinishRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_island_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckpointFinishRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckpointFinishRes) ProtoMessage() {}

func (x *CheckpointFinishRes) ProtoReflect() protoreflect.Message {
	mi := &file_island_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckpointFinishRes.ProtoReflect.Descriptor instead.
func (*CheckpointFinishRes) Descriptor() ([]byte, []int) {
	return file_island_proto_rawDescGZIP(), []int{5}
}

func (x *CheckpointFinishRes) GetResult() MSG_RESULT {
	if x != nil {
		return x.Result
	}
	return MSG_RESULT_SUCCESS
}

// 保存岛屿信息
type SaveIslandInfoReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IslandId int64  `protobuf:"varint,1,opt,name=IslandId,json=islandId,proto3" json:"IslandId"`
	JsonInfo string `protobuf:"bytes,2,opt,name=JsonInfo,json=jsonInfo,proto3" json:"JsonInfo"`
}

func (x *SaveIslandInfoReq) Reset() {
	*x = SaveIslandInfoReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_island_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveIslandInfoReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveIslandInfoReq) ProtoMessage() {}

func (x *SaveIslandInfoReq) ProtoReflect() protoreflect.Message {
	mi := &file_island_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveIslandInfoReq.ProtoReflect.Descriptor instead.
func (*SaveIslandInfoReq) Descriptor() ([]byte, []int) {
	return file_island_proto_rawDescGZIP(), []int{6}
}

func (x *SaveIslandInfoReq) GetIslandId() int64 {
	if x != nil {
		return x.IslandId
	}
	return 0
}

func (x *SaveIslandInfoReq) GetJsonInfo() string {
	if x != nil {
		return x.JsonInfo
	}
	return ""
}

type SaveIslandInfoRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result MSG_RESULT `protobuf:"varint,1,opt,name=Result,json=result,proto3,enum=message.MSG_RESULT" json:"Result"` //成功码
}

func (x *SaveIslandInfoRes) Reset() {
	*x = SaveIslandInfoRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_island_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveIslandInfoRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveIslandInfoRes) ProtoMessage() {}

func (x *SaveIslandInfoRes) ProtoReflect() protoreflect.Message {
	mi := &file_island_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveIslandInfoRes.ProtoReflect.Descriptor instead.
func (*SaveIslandInfoRes) Descriptor() ([]byte, []int) {
	return file_island_proto_rawDescGZIP(), []int{7}
}

func (x *SaveIslandInfoRes) GetResult() MSG_RESULT {
	if x != nil {
		return x.Result
	}
	return MSG_RESULT_SUCCESS
}

var File_island_proto protoreflect.FileDescriptor

var file_island_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x69, 0x73, 0x6c, 0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0d, 0x6d, 0x73, 0x67, 0x74, 0x79, 0x70, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdc, 0x01, 0x0a, 0x06, 0x49, 0x73, 0x6c, 0x61, 0x6e,
	0x64, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x22, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x42, 0x0a, 0x0b, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x69, 0x6e,
	0x67, 0x4d, 0x61, 0x70, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x49, 0x73, 0x6c, 0x61, 0x6e, 0x64, 0x2e, 0x42, 0x75, 0x69, 0x6c,
	0x64, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x62, 0x75,
	0x69, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x4a, 0x73, 0x6f,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6a, 0x73, 0x6f,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x3e, 0x0a, 0x10, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x69, 0x6e,
	0x67, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x3d, 0x0a, 0x10, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x49,
	0x73, 0x6c, 0x61, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x29, 0x0a, 0x07, 0x49, 0x73, 0x6c,
	0x61, 0x6e, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x49, 0x73, 0x6c, 0x61, 0x6e, 0x64, 0x52, 0x07, 0x69, 0x73, 0x6c,
	0x61, 0x6e, 0x64, 0x73, 0x22, 0x37, 0x0a, 0x10, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x42, 0x75, 0x69,
	0x6c, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x12, 0x23, 0x0a, 0x04, 0x6c, 0x61, 0x6e, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x49, 0x73, 0x6c, 0x61, 0x6e, 0x64, 0x52, 0x04, 0x6c, 0x61, 0x6e, 0x64, 0x22, 0x3f, 0x0a,
	0x10, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x73, 0x12, 0x2b, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x13, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x53, 0x47, 0x5f,
	0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x83,
	0x01, 0x0a, 0x13, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x46, 0x69, 0x6e,
	0x69, 0x73, 0x68, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x47, 0x6f, 0x6c, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x67, 0x6f, 0x6c, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x73,
	0x6c, 0x61, 0x6e, 0x64, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x69, 0x73,
	0x6c, 0x61, 0x6e, 0x64, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x63, 0x68,
	0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x22, 0x42, 0x0a, 0x13, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x52, 0x65, 0x73, 0x12, 0x2b, 0x0a, 0x06, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x53, 0x47, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54,
	0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x4b, 0x0a, 0x11, 0x53, 0x61, 0x76, 0x65,
	0x49, 0x73, 0x6c, 0x61, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a,
	0x08, 0x49, 0x73, 0x6c, 0x61, 0x6e, 0x64, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x08, 0x69, 0x73, 0x6c, 0x61, 0x6e, 0x64, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x4a, 0x73, 0x6f,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6a, 0x73, 0x6f,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x40, 0x0a, 0x11, 0x53, 0x61, 0x76, 0x65, 0x49, 0x73, 0x6c,
	0x61, 0x6e, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x12, 0x2b, 0x0a, 0x06, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x53, 0x47, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x52,
	0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x0a, 0x5a, 0x08, 0x2f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_island_proto_rawDescOnce sync.Once
	file_island_proto_rawDescData = file_island_proto_rawDesc
)

func file_island_proto_rawDescGZIP() []byte {
	file_island_proto_rawDescOnce.Do(func() {
		file_island_proto_rawDescData = protoimpl.X.CompressGZIP(file_island_proto_rawDescData)
	})
	return file_island_proto_rawDescData
}

var file_island_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_island_proto_goTypes = []interface{}{
	(*Island)(nil),              // 0: message.Island
	(*NotifyIslandInfo)(nil),    // 1: message.NotifyIslandInfo
	(*PlaceBuildingReq)(nil),    // 2: message.PlaceBuildingReq
	(*PlaceBuildingRes)(nil),    // 3: message.PlaceBuildingRes
	(*CheckpointFinishReq)(nil), // 4: message.CheckpointFinishReq
	(*CheckpointFinishRes)(nil), // 5: message.CheckpointFinishRes
	(*SaveIslandInfoReq)(nil),   // 6: message.SaveIslandInfoReq
	(*SaveIslandInfoRes)(nil),   // 7: message.SaveIslandInfoRes
	nil,                         // 8: message.Island.BuildingMapEntry
	(MSG_RESULT)(0),             // 9: message.MSG_RESULT
}
var file_island_proto_depIdxs = []int32{
	8, // 0: message.Island.BuildingMap:type_name -> message.Island.BuildingMapEntry
	0, // 1: message.NotifyIslandInfo.Islands:type_name -> message.Island
	0, // 2: message.PlaceBuildingReq.land:type_name -> message.Island
	9, // 3: message.PlaceBuildingRes.Result:type_name -> message.MSG_RESULT
	9, // 4: message.CheckpointFinishRes.Result:type_name -> message.MSG_RESULT
	9, // 5: message.SaveIslandInfoRes.Result:type_name -> message.MSG_RESULT
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_island_proto_init() }
func file_island_proto_init() {
	if File_island_proto != nil {
		return
	}
	file_msgtype_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_island_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Island); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_island_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyIslandInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_island_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlaceBuildingReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_island_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlaceBuildingRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_island_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckpointFinishReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_island_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckpointFinishRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_island_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveIslandInfoReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_island_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveIslandInfoRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_island_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_island_proto_goTypes,
		DependencyIndexes: file_island_proto_depIdxs,
		MessageInfos:      file_island_proto_msgTypes,
	}.Build()
	File_island_proto = out.File
	file_island_proto_rawDesc = nil
	file_island_proto_goTypes = nil
	file_island_proto_depIdxs = nil
}
