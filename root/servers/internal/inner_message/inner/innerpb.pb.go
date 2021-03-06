// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        (unknown)
// source: inner/innerpb.proto

package inner

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

type INNER_MSG int32

const (
	//------------------ 服务器内部通信使用 id段 10000-100000 --------------------------------------
	INNER_MSG_ERROR                      INNER_MSG = 0
	INNER_MSG_INNER_SEGMENT_BEGIN        INNER_MSG = 10000  // 消息段begin
	INNER_MSG_L2D_ALL_USER_REQ           INNER_MSG = 10001  // L2dAllUserReq   login 向db请求所有用户角色数据
	INNER_MSG_D2L_ALL_USER_RES           INNER_MSG = 10002  // D2lAllUserRes
	INNER_MSG_L2D_USER_SAVE              INNER_MSG = 10003  // L2dUserSave
	INNER_MSG_G2D_MODEL_SAVE             INNER_MSG = 10006  // G2DModelSave  存储数据
	INNER_MSG_G2D_ROLE_REQ               INNER_MSG = 10007  // G2DRoleReq  game请求role数据
	INNER_MSG_D2G_ROLE_RES               INNER_MSG = 10008  // D2GRoleRes
	INNER_MSG_GT2G_SESSION_CLOSED        INNER_MSG = 10010  // GT2GSessionClosed      gatesession通知用户网络断开
	INNER_MSG_L2GT_SESSION_ASSIGN_GAME   INNER_MSG = 10011  // L2GSessionAssignGame    通知gate,gatesession分配gameActor
	INNER_MSG_L2GT_USER_SESSION_DISABLED INNER_MSG = 10012  // L2GTUserSessionDisabled login通知game,玩家顶号，旧session失效
	INNER_MSG_G2L_ROLE_OFFLINE           INNER_MSG = 10013  // G2LRoleOffline          game通知login，玩家离线
	INNER_MSG_G2D_GAME_STOP              INNER_MSG = 10014  // G2DGameStop             game通知gamedb 停服
	INNER_MSG_INNER_SEGMENT_END          INNER_MSG = 100000 // 消息段end
)

// Enum value maps for INNER_MSG.
var (
	INNER_MSG_name = map[int32]string{
		0:      "ERROR",
		10000:  "INNER_SEGMENT_BEGIN",
		10001:  "L2D_ALL_USER_REQ",
		10002:  "D2L_ALL_USER_RES",
		10003:  "L2D_USER_SAVE",
		10006:  "G2D_MODEL_SAVE",
		10007:  "G2D_ROLE_REQ",
		10008:  "D2G_ROLE_RES",
		10010:  "GT2G_SESSION_CLOSED",
		10011:  "L2GT_SESSION_ASSIGN_GAME",
		10012:  "L2GT_USER_SESSION_DISABLED",
		10013:  "G2L_ROLE_OFFLINE",
		10014:  "G2D_GAME_STOP",
		100000: "INNER_SEGMENT_END",
	}
	INNER_MSG_value = map[string]int32{
		"ERROR":                      0,
		"INNER_SEGMENT_BEGIN":        10000,
		"L2D_ALL_USER_REQ":           10001,
		"D2L_ALL_USER_RES":           10002,
		"L2D_USER_SAVE":              10003,
		"G2D_MODEL_SAVE":             10006,
		"G2D_ROLE_REQ":               10007,
		"D2G_ROLE_RES":               10008,
		"GT2G_SESSION_CLOSED":        10010,
		"L2GT_SESSION_ASSIGN_GAME":   10011,
		"L2GT_USER_SESSION_DISABLED": 10012,
		"G2L_ROLE_OFFLINE":           10013,
		"G2D_GAME_STOP":              10014,
		"INNER_SEGMENT_END":          100000,
	}
)

func (x INNER_MSG) Enum() *INNER_MSG {
	p := new(INNER_MSG)
	*p = x
	return p
}

func (x INNER_MSG) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (INNER_MSG) Descriptor() protoreflect.EnumDescriptor {
	return file_inner_innerpb_proto_enumTypes[0].Descriptor()
}

func (INNER_MSG) Type() protoreflect.EnumType {
	return &file_inner_innerpb_proto_enumTypes[0]
}

func (x INNER_MSG) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use INNER_MSG.Descriptor instead.
func (INNER_MSG) EnumDescriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{0}
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{0}
}

// 存储信息
type G2DModelSave struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data      []byte `protobuf:"bytes,1,opt,name=Data,json=data,proto3" json:"Data"`
	ModelName string `protobuf:"bytes,2,opt,name=ModelName,json=modelName,proto3" json:"ModelName"`
	RID       int64  `protobuf:"varint,3,opt,name=RID,json=rID,proto3" json:"RID"`
}

func (x *G2DModelSave) Reset() {
	*x = G2DModelSave{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *G2DModelSave) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*G2DModelSave) ProtoMessage() {}

func (x *G2DModelSave) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use G2DModelSave.ProtoReflect.Descriptor instead.
func (*G2DModelSave) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{1}
}

func (x *G2DModelSave) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *G2DModelSave) GetModelName() string {
	if x != nil {
		return x.ModelName
	}
	return ""
}

func (x *G2DModelSave) GetRID() int64 {
	if x != nil {
		return x.RID
	}
	return 0
}

// login 请求所有角色数据
type L2DAllUserReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *L2DAllUserReq) Reset() {
	*x = L2DAllUserReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *L2DAllUserReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*L2DAllUserReq) ProtoMessage() {}

func (x *L2DAllUserReq) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use L2DAllUserReq.ProtoReflect.Descriptor instead.
func (*L2DAllUserReq) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{2}
}

type D2LAllUserRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=Data,json=data,proto3" json:"Data"` // models.GameUser
}

func (x *D2LAllUserRes) Reset() {
	*x = D2LAllUserRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *D2LAllUserRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*D2LAllUserRes) ProtoMessage() {}

func (x *D2LAllUserRes) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use D2LAllUserRes.ProtoReflect.Descriptor instead.
func (*D2LAllUserRes) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{3}
}

func (x *D2LAllUserRes) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// login 请求回存新玩家
type L2DUserSave struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=Data,json=data,proto3" json:"Data"` // models.GameUser
}

func (x *L2DUserSave) Reset() {
	*x = L2DUserSave{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *L2DUserSave) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*L2DUserSave) ProtoMessage() {}

func (x *L2DUserSave) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use L2DUserSave.ProtoReflect.Descriptor instead.
func (*L2DUserSave) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{4}
}

func (x *L2DUserSave) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// game请求role数据
type G2DRoleReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CallbackId int64 `protobuf:"varint,1,opt,name=CallbackId,json=callbackId,proto3" json:"CallbackId"` // 回调Id
	RID        int64 `protobuf:"varint,2,opt,name=RID,json=rID,proto3" json:"RID"`
}

func (x *G2DRoleReq) Reset() {
	*x = G2DRoleReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *G2DRoleReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*G2DRoleReq) ProtoMessage() {}

func (x *G2DRoleReq) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use G2DRoleReq.ProtoReflect.Descriptor instead.
func (*G2DRoleReq) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{5}
}

func (x *G2DRoleReq) GetCallbackId() int64 {
	if x != nil {
		return x.CallbackId
	}
	return 0
}

func (x *G2DRoleReq) GetRID() int64 {
	if x != nil {
		return x.RID
	}
	return 0
}

type D2GRoleRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CallbackId int64             `protobuf:"varint,1,opt,name=CallbackId,json=callbackId,proto3" json:"CallbackId"` // 回调Id
	Valid      int32             `protobuf:"varint,2,opt,name=valid,proto3" json:"valid"`                           // 数据是否找到 0.成功 1.未找到角色
	RID        int64             `protobuf:"varint,3,opt,name=RID,json=rID,proto3" json:"RID"`
	WholeInfo  map[string][]byte `protobuf:"bytes,4,rep,name=WholeInfo,json=wholeInfo,proto3" json:"WholeInfo" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // 角色各个模块的数据
}

func (x *D2GRoleRes) Reset() {
	*x = D2GRoleRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *D2GRoleRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*D2GRoleRes) ProtoMessage() {}

func (x *D2GRoleRes) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use D2GRoleRes.ProtoReflect.Descriptor instead.
func (*D2GRoleRes) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{6}
}

func (x *D2GRoleRes) GetCallbackId() int64 {
	if x != nil {
		return x.CallbackId
	}
	return 0
}

func (x *D2GRoleRes) GetValid() int32 {
	if x != nil {
		return x.Valid
	}
	return 0
}

func (x *D2GRoleRes) GetRID() int64 {
	if x != nil {
		return x.RID
	}
	return 0
}

func (x *D2GRoleRes) GetWholeInfo() map[string][]byte {
	if x != nil {
		return x.WholeInfo
	}
	return nil
}

// gatesession断开,gate通知game
type GT2GSessionClosed struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GateSession string `protobuf:"bytes,1,opt,name=GateSession,json=gateSession,proto3" json:"GateSession"` // 断开的session
}

func (x *GT2GSessionClosed) Reset() {
	*x = GT2GSessionClosed{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GT2GSessionClosed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GT2GSessionClosed) ProtoMessage() {}

func (x *GT2GSessionClosed) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GT2GSessionClosed.ProtoReflect.Descriptor instead.
func (*GT2GSessionClosed) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{7}
}

func (x *GT2GSessionClosed) GetGateSession() string {
	if x != nil {
		return x.GateSession
	}
	return ""
}

// login通知gate 为session分配gameActor
type L2GTSessionAssignGame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GateSession string `protobuf:"bytes,1,opt,name=GateSession,json=gateSession,proto3" json:"GateSession"` //
	GameActorId string `protobuf:"bytes,2,opt,name=GameActorId,json=gameActorId,proto3" json:"GameActorId"` // 分配的游戏Actor
}

func (x *L2GTSessionAssignGame) Reset() {
	*x = L2GTSessionAssignGame{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *L2GTSessionAssignGame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*L2GTSessionAssignGame) ProtoMessage() {}

func (x *L2GTSessionAssignGame) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use L2GTSessionAssignGame.ProtoReflect.Descriptor instead.
func (*L2GTSessionAssignGame) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{8}
}

func (x *L2GTSessionAssignGame) GetGateSession() string {
	if x != nil {
		return x.GateSession
	}
	return ""
}

func (x *L2GTSessionAssignGame) GetGameActorId() string {
	if x != nil {
		return x.GameActorId
	}
	return ""
}

// game通知login，玩家离线
type G2LRoleOffline struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UID int64 `protobuf:"varint,1,opt,name=UID,json=uID,proto3" json:"UID"`
	RID int64 `protobuf:"varint,2,opt,name=RID,json=rID,proto3" json:"RID"`
}

func (x *G2LRoleOffline) Reset() {
	*x = G2LRoleOffline{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *G2LRoleOffline) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*G2LRoleOffline) ProtoMessage() {}

func (x *G2LRoleOffline) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use G2LRoleOffline.ProtoReflect.Descriptor instead.
func (*G2LRoleOffline) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{9}
}

func (x *G2LRoleOffline) GetUID() int64 {
	if x != nil {
		return x.UID
	}
	return 0
}

func (x *G2LRoleOffline) GetRID() int64 {
	if x != nil {
		return x.RID
	}
	return 0
}

// game通知login，玩家离线
type G2DGameStop struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *G2DGameStop) Reset() {
	*x = G2DGameStop{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *G2DGameStop) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*G2DGameStop) ProtoMessage() {}

func (x *G2DGameStop) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use G2DGameStop.ProtoReflect.Descriptor instead.
func (*G2DGameStop) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{10}
}

// login通知game,玩家顶号，旧session失效
type L2GTUserSessionDisabled struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GateSession string `protobuf:"bytes,1,opt,name=GateSession,json=gateSession,proto3" json:"GateSession"`
	UID         int64  `protobuf:"varint,2,opt,name=UID,json=uID,proto3" json:"UID"`
}

func (x *L2GTUserSessionDisabled) Reset() {
	*x = L2GTUserSessionDisabled{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inner_innerpb_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *L2GTUserSessionDisabled) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*L2GTUserSessionDisabled) ProtoMessage() {}

func (x *L2GTUserSessionDisabled) ProtoReflect() protoreflect.Message {
	mi := &file_inner_innerpb_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use L2GTUserSessionDisabled.ProtoReflect.Descriptor instead.
func (*L2GTUserSessionDisabled) Descriptor() ([]byte, []int) {
	return file_inner_innerpb_proto_rawDescGZIP(), []int{11}
}

func (x *L2GTUserSessionDisabled) GetGateSession() string {
	if x != nil {
		return x.GateSession
	}
	return ""
}

func (x *L2GTUserSessionDisabled) GetUID() int64 {
	if x != nil {
		return x.UID
	}
	return 0
}

var File_inner_innerpb_proto protoreflect.FileDescriptor

var file_inner_innerpb_proto_rawDesc = []byte{
	0x0a, 0x13, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x70, 0x62, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x52, 0x0a,
	0x0c, 0x47, 0x32, 0x44, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x53, 0x61, 0x76, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x1c, 0x0a, 0x09, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x52, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x72, 0x49,
	0x44, 0x22, 0x0f, 0x0a, 0x0d, 0x4c, 0x32, 0x44, 0x41, 0x6c, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x22, 0x23, 0x0a, 0x0d, 0x44, 0x32, 0x4c, 0x41, 0x6c, 0x6c, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x21, 0x0a, 0x0b, 0x4c, 0x32, 0x44, 0x55, 0x73,
	0x65, 0x72, 0x53, 0x61, 0x76, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3e, 0x0a, 0x0a, 0x47, 0x32,
	0x44, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x61, 0x6c, 0x6c,
	0x62, 0x61, 0x63, 0x6b, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x61,
	0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x52, 0x49, 0x44, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x72, 0x49, 0x44, 0x22, 0xda, 0x01, 0x0a, 0x0a, 0x44,
	0x32, 0x47, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x61, 0x6c,
	0x6c, 0x62, 0x61, 0x63, 0x6b, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63,
	0x61, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x52, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x72, 0x49,
	0x44, 0x12, 0x46, 0x0a, 0x09, 0x57, 0x68, 0x6f, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x44, 0x32, 0x47, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x2e,
	0x57, 0x68, 0x6f, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09,
	0x77, 0x68, 0x6f, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x3c, 0x0a, 0x0e, 0x57, 0x68, 0x6f,
	0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x35, 0x0a, 0x11, 0x47, 0x54, 0x32, 0x47, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x12, 0x20, 0x0a, 0x0b,
	0x47, 0x61, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x67, 0x61, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x5b,
	0x0a, 0x15, 0x4c, 0x32, 0x47, 0x54, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x41, 0x73, 0x73,
	0x69, 0x67, 0x6e, 0x47, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x47, 0x61, 0x74, 0x65, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x67, 0x61,
	0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x47, 0x61, 0x6d,
	0x65, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x67, 0x61, 0x6d, 0x65, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x22, 0x34, 0x0a, 0x0e, 0x47,
	0x32, 0x4c, 0x52, 0x6f, 0x6c, 0x65, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x49, 0x44, 0x12,
	0x10, 0x0a, 0x03, 0x52, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x72, 0x49,
	0x44, 0x22, 0x0d, 0x0a, 0x0b, 0x47, 0x32, 0x44, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x6f, 0x70,
	0x22, 0x4d, 0x0a, 0x17, 0x4c, 0x32, 0x47, 0x54, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x47,
	0x61, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x67, 0x61, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a,
	0x03, 0x55, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x49, 0x44, 0x2a,
	0xcb, 0x02, 0x0a, 0x09, 0x49, 0x4e, 0x4e, 0x45, 0x52, 0x5f, 0x4d, 0x53, 0x47, 0x12, 0x09, 0x0a,
	0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x13, 0x49, 0x4e, 0x4e, 0x45,
	0x52, 0x5f, 0x53, 0x45, 0x47, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x42, 0x45, 0x47, 0x49, 0x4e, 0x10,
	0x90, 0x4e, 0x12, 0x15, 0x0a, 0x10, 0x4c, 0x32, 0x44, 0x5f, 0x41, 0x4c, 0x4c, 0x5f, 0x55, 0x53,
	0x45, 0x52, 0x5f, 0x52, 0x45, 0x51, 0x10, 0x91, 0x4e, 0x12, 0x15, 0x0a, 0x10, 0x44, 0x32, 0x4c,
	0x5f, 0x41, 0x4c, 0x4c, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x52, 0x45, 0x53, 0x10, 0x92, 0x4e,
	0x12, 0x12, 0x0a, 0x0d, 0x4c, 0x32, 0x44, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x41, 0x56,
	0x45, 0x10, 0x93, 0x4e, 0x12, 0x13, 0x0a, 0x0e, 0x47, 0x32, 0x44, 0x5f, 0x4d, 0x4f, 0x44, 0x45,
	0x4c, 0x5f, 0x53, 0x41, 0x56, 0x45, 0x10, 0x96, 0x4e, 0x12, 0x11, 0x0a, 0x0c, 0x47, 0x32, 0x44,
	0x5f, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x52, 0x45, 0x51, 0x10, 0x97, 0x4e, 0x12, 0x11, 0x0a, 0x0c,
	0x44, 0x32, 0x47, 0x5f, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x52, 0x45, 0x53, 0x10, 0x98, 0x4e, 0x12,
	0x18, 0x0a, 0x13, 0x47, 0x54, 0x32, 0x47, 0x5f, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f,
	0x43, 0x4c, 0x4f, 0x53, 0x45, 0x44, 0x10, 0x9a, 0x4e, 0x12, 0x1d, 0x0a, 0x18, 0x4c, 0x32, 0x47,
	0x54, 0x5f, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x41, 0x53, 0x53, 0x49, 0x47, 0x4e,
	0x5f, 0x47, 0x41, 0x4d, 0x45, 0x10, 0x9b, 0x4e, 0x12, 0x1f, 0x0a, 0x1a, 0x4c, 0x32, 0x47, 0x54,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x49,
	0x53, 0x41, 0x42, 0x4c, 0x45, 0x44, 0x10, 0x9c, 0x4e, 0x12, 0x15, 0x0a, 0x10, 0x47, 0x32, 0x4c,
	0x5f, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x4f, 0x46, 0x46, 0x4c, 0x49, 0x4e, 0x45, 0x10, 0x9d, 0x4e,
	0x12, 0x12, 0x0a, 0x0d, 0x47, 0x32, 0x44, 0x5f, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x53, 0x54, 0x4f,
	0x50, 0x10, 0x9e, 0x4e, 0x12, 0x17, 0x0a, 0x11, 0x49, 0x4e, 0x4e, 0x45, 0x52, 0x5f, 0x53, 0x45,
	0x47, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x45, 0x4e, 0x44, 0x10, 0xa0, 0x8d, 0x06, 0x42, 0x08, 0x5a,
	0x06, 0x2f, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_inner_innerpb_proto_rawDescOnce sync.Once
	file_inner_innerpb_proto_rawDescData = file_inner_innerpb_proto_rawDesc
)

func file_inner_innerpb_proto_rawDescGZIP() []byte {
	file_inner_innerpb_proto_rawDescOnce.Do(func() {
		file_inner_innerpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_inner_innerpb_proto_rawDescData)
	})
	return file_inner_innerpb_proto_rawDescData
}

var file_inner_innerpb_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_inner_innerpb_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_inner_innerpb_proto_goTypes = []interface{}{
	(INNER_MSG)(0),                  // 0: inner_message.INNER_MSG
	(*Error)(nil),                   // 1: inner_message.Error
	(*G2DModelSave)(nil),            // 2: inner_message.G2DModelSave
	(*L2DAllUserReq)(nil),           // 3: inner_message.L2DAllUserReq
	(*D2LAllUserRes)(nil),           // 4: inner_message.D2LAllUserRes
	(*L2DUserSave)(nil),             // 5: inner_message.L2DUserSave
	(*G2DRoleReq)(nil),              // 6: inner_message.G2DRoleReq
	(*D2GRoleRes)(nil),              // 7: inner_message.D2GRoleRes
	(*GT2GSessionClosed)(nil),       // 8: inner_message.GT2GSessionClosed
	(*L2GTSessionAssignGame)(nil),   // 9: inner_message.L2GTSessionAssignGame
	(*G2LRoleOffline)(nil),          // 10: inner_message.G2LRoleOffline
	(*G2DGameStop)(nil),             // 11: inner_message.G2DGameStop
	(*L2GTUserSessionDisabled)(nil), // 12: inner_message.L2GTUserSessionDisabled
	nil,                             // 13: inner_message.D2GRoleRes.WholeInfoEntry
}
var file_inner_innerpb_proto_depIdxs = []int32{
	13, // 0: inner_message.D2GRoleRes.WholeInfo:type_name -> inner_message.D2GRoleRes.WholeInfoEntry
	1,  // [1:1] is the sub-list for method output_type
	1,  // [1:1] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_inner_innerpb_proto_init() }
func file_inner_innerpb_proto_init() {
	if File_inner_innerpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_inner_innerpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
		file_inner_innerpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*G2DModelSave); i {
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
		file_inner_innerpb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*L2DAllUserReq); i {
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
		file_inner_innerpb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*D2LAllUserRes); i {
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
		file_inner_innerpb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*L2DUserSave); i {
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
		file_inner_innerpb_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*G2DRoleReq); i {
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
		file_inner_innerpb_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*D2GRoleRes); i {
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
		file_inner_innerpb_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GT2GSessionClosed); i {
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
		file_inner_innerpb_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*L2GTSessionAssignGame); i {
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
		file_inner_innerpb_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*G2LRoleOffline); i {
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
		file_inner_innerpb_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*G2DGameStop); i {
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
		file_inner_innerpb_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*L2GTUserSessionDisabled); i {
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
			RawDescriptor: file_inner_innerpb_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_inner_innerpb_proto_goTypes,
		DependencyIndexes: file_inner_innerpb_proto_depIdxs,
		EnumInfos:         file_inner_innerpb_proto_enumTypes,
		MessageInfos:      file_inner_innerpb_proto_msgTypes,
	}.Build()
	File_inner_innerpb_proto = out.File
	file_inner_innerpb_proto_rawDesc = nil
	file_inner_innerpb_proto_goTypes = nil
	file_inner_innerpb_proto_depIdxs = nil
}
