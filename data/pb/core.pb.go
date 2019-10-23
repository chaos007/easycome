// Code generated by protoc-gen-go.
// source: core.proto
// DO NOT EDIT!

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// TestToAgent TestToAgent
type TestToAgent struct {
	FieldA string `protobuf:"bytes,1,opt,name=FieldA" json:"FieldA,omitempty"`
}

func (m *TestToAgent) Reset()                    { *m = TestToAgent{} }
func (m *TestToAgent) String() string            { return proto.CompactTextString(m) }
func (*TestToAgent) ProtoMessage()               {}
func (*TestToAgent) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

// TestToGame TestToGame
type TestToGame struct {
	FieldA string `protobuf:"bytes,1,opt,name=FieldA" json:"FieldA,omitempty"`
	FieldB string `protobuf:"bytes,2,opt,name=FieldB" json:"FieldB,omitempty"`
}

func (m *TestToGame) Reset()                    { *m = TestToGame{} }
func (m *TestToGame) String() string            { return proto.CompactTextString(m) }
func (*TestToGame) ProtoMessage()               {}
func (*TestToGame) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

// UpPing UpPing
type UpPing struct {
}

func (m *UpPing) Reset()                    { *m = UpPing{} }
func (m *UpPing) String() string            { return proto.CompactTextString(m) }
func (*UpPing) ProtoMessage()               {}
func (*UpPing) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

// DownPing DownPing
type DownPing struct {
}

func (m *DownPing) Reset()                    { *m = DownPing{} }
func (m *DownPing) String() string            { return proto.CompactTextString(m) }
func (*DownPing) ProtoMessage()               {}
func (*DownPing) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

// web请求
// register?data={{json}}&syn={{syn}}
// WebUpRegister 上发游客登录，使用http上发游客登录
type WebUpRegister struct {
	MachineCode string `protobuf:"bytes,1,opt,name=MachineCode" json:"MachineCode,omitempty"`
}

func (m *WebUpRegister) Reset()                    { *m = WebUpRegister{} }
func (m *WebUpRegister) String() string            { return proto.CompactTextString(m) }
func (*WebUpRegister) ProtoMessage()               {}
func (*WebUpRegister) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

// WebDownRegister 下发游客登录
type WebDownRegister struct {
	UserID      string      `protobuf:"bytes,1,opt,name=UserID" json:"UserID,omitempty"`
	Password    string      `protobuf:"bytes,2,opt,name=Password" json:"Password,omitempty"`
	Info        *ServerInfo `protobuf:"bytes,3,opt,name=Info" json:"Info,omitempty"`
	Cookie      string      `protobuf:"bytes,4,opt,name=Cookie" json:"Cookie,omitempty"`
	MachineCode string      `protobuf:"bytes,5,opt,name=MachineCode" json:"MachineCode,omitempty"`
}

func (m *WebDownRegister) Reset()                    { *m = WebDownRegister{} }
func (m *WebDownRegister) String() string            { return proto.CompactTextString(m) }
func (*WebDownRegister) ProtoMessage()               {}
func (*WebDownRegister) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *WebDownRegister) GetInfo() *ServerInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

// web请求
// login?data={{json}}&syn={{syn}}
// WebUpGetServerInfo 使用http登录
type WebUpGetServerInfo struct {
	UserName string `protobuf:"bytes,1,opt,name=UserName" json:"UserName,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=Password" json:"Password,omitempty"`
}

func (m *WebUpGetServerInfo) Reset()                    { *m = WebUpGetServerInfo{} }
func (m *WebUpGetServerInfo) String() string            { return proto.CompactTextString(m) }
func (*WebUpGetServerInfo) ProtoMessage()               {}
func (*WebUpGetServerInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

// WebDownGetServerInfo 下发http登录
type WebDownGetServerInfo struct {
	Info   *ServerInfo `protobuf:"bytes,1,opt,name=Info" json:"Info,omitempty"`
	UserID string      `protobuf:"bytes,2,opt,name=UserID" json:"UserID,omitempty"`
	Cookie string      `protobuf:"bytes,3,opt,name=Cookie" json:"Cookie,omitempty"`
}

func (m *WebDownGetServerInfo) Reset()                    { *m = WebDownGetServerInfo{} }
func (m *WebDownGetServerInfo) String() string            { return proto.CompactTextString(m) }
func (*WebDownGetServerInfo) ProtoMessage()               {}
func (*WebDownGetServerInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

func (m *WebDownGetServerInfo) GetInfo() *ServerInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

// ServerInfo ServerInfo
type ServerInfo struct {
	AgentIP   string `protobuf:"bytes,1,opt,name=AgentIP" json:"AgentIP,omitempty"`
	AgentPort string `protobuf:"bytes,2,opt,name=AgentPort" json:"AgentPort,omitempty"`
}

func (m *ServerInfo) Reset()                    { *m = ServerInfo{} }
func (m *ServerInfo) String() string            { return proto.CompactTextString(m) }
func (*ServerInfo) ProtoMessage()               {}
func (*ServerInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{8} }

// UpGetConfigToCenter 向中心服请求配置表
type UpGetConfigToCenter struct {
}

func (m *UpGetConfigToCenter) Reset()                    { *m = UpGetConfigToCenter{} }
func (m *UpGetConfigToCenter) String() string            { return proto.CompactTextString(m) }
func (*UpGetConfigToCenter) ProtoMessage()               {}
func (*UpGetConfigToCenter) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{9} }

// SyncConfigData 同步数据
type SyncConfigDataToAll struct {
	Data string `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *SyncConfigDataToAll) Reset()                    { *m = SyncConfigDataToAll{} }
func (m *SyncConfigDataToAll) String() string            { return proto.CompactTextString(m) }
func (*SyncConfigDataToAll) ProtoMessage()               {}
func (*SyncConfigDataToAll) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{10} }

// LoginKickPlayerToAgent 登录服踢出玩家
type LoginKickPlayerToAgent struct {
	UserID int64 `protobuf:"varint,1,opt,name=UserID" json:"UserID,omitempty"`
}

func (m *LoginKickPlayerToAgent) Reset()                    { *m = LoginKickPlayerToAgent{} }
func (m *LoginKickPlayerToAgent) String() string            { return proto.CompactTextString(m) }
func (*LoginKickPlayerToAgent) ProtoMessage()               {}
func (*LoginKickPlayerToAgent) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{11} }

// AgentCheckUserToLogin  向登陆服请求玩家数据正确
type AgentCheckUserToLogin struct {
	UserID int64  `protobuf:"varint,1,opt,name=UserID" json:"UserID,omitempty"`
	Cookie string `protobuf:"bytes,2,opt,name=Cookie" json:"Cookie,omitempty"`
}

func (m *AgentCheckUserToLogin) Reset()                    { *m = AgentCheckUserToLogin{} }
func (m *AgentCheckUserToLogin) String() string            { return proto.CompactTextString(m) }
func (*AgentCheckUserToLogin) ProtoMessage()               {}
func (*AgentCheckUserToLogin) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{12} }

// LoginCheckPlayerToAgent 检查登录的玩家是否合法
type LoginCheckPlayerToAgent struct {
	UserID  int64 `protobuf:"varint,1,opt,name=UserID" json:"UserID,omitempty"`
	IsLegal bool  `protobuf:"varint,2,opt,name=IsLegal" json:"IsLegal,omitempty"`
}

func (m *LoginCheckPlayerToAgent) Reset()                    { *m = LoginCheckPlayerToAgent{} }
func (m *LoginCheckPlayerToAgent) String() string            { return proto.CompactTextString(m) }
func (*LoginCheckPlayerToAgent) ProtoMessage()               {}
func (*LoginCheckPlayerToAgent) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{13} }

// PlayerLogin 玩家请求登陆
type PlayerLogin struct {
	UserID int64  `protobuf:"varint,1,opt,name=UserID" json:"UserID,omitempty"`
	Cookie string `protobuf:"bytes,2,opt,name=Cookie" json:"Cookie,omitempty"`
}

func (m *PlayerLogin) Reset()                    { *m = PlayerLogin{} }
func (m *PlayerLogin) String() string            { return proto.CompactTextString(m) }
func (*PlayerLogin) ProtoMessage()               {}
func (*PlayerLogin) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{14} }

// PlayerLoginToGame 玩家向游戏服请求数据
type PlayerLoginToGame struct {
	UserID int64 `protobuf:"varint,1,opt,name=UserID" json:"UserID,omitempty"`
}

func (m *PlayerLoginToGame) Reset()                    { *m = PlayerLoginToGame{} }
func (m *PlayerLoginToGame) String() string            { return proto.CompactTextString(m) }
func (*PlayerLoginToGame) ProtoMessage()               {}
func (*PlayerLoginToGame) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{15} }

// DownPlayerLogin 玩家登陆成功
type DownPlayerLogin struct {
	RoleID    int32 `protobuf:"varint,1,opt,name=RoleID" json:"RoleID,omitempty"`
	IsSucceed bool  `protobuf:"varint,2,opt,name=IsSucceed" json:"IsSucceed,omitempty"`
}

func (m *DownPlayerLogin) Reset()                    { *m = DownPlayerLogin{} }
func (m *DownPlayerLogin) String() string            { return proto.CompactTextString(m) }
func (*DownPlayerLogin) ProtoMessage()               {}
func (*DownPlayerLogin) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{16} }

// WebUpRegisterPassword?data={{json}}&syn={{syn}}
// WebUpRegisterPassword 上发注册，使用http上发注册 加上密码
type WebUpRegisterPassword struct {
	UserID   string `protobuf:"bytes,1,opt,name=UserID" json:"UserID,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=Password" json:"Password,omitempty"`
	Code     string `protobuf:"bytes,3,opt,name=Code" json:"Code,omitempty"`
}

func (m *WebUpRegisterPassword) Reset()                    { *m = WebUpRegisterPassword{} }
func (m *WebUpRegisterPassword) String() string            { return proto.CompactTextString(m) }
func (*WebUpRegisterPassword) ProtoMessage()               {}
func (*WebUpRegisterPassword) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{17} }

// WebDownRegisterPassword 下发注册
type WebDownRegisterPassword struct {
	UserID string      `protobuf:"bytes,1,opt,name=UserID" json:"UserID,omitempty"`
	Info   *ServerInfo `protobuf:"bytes,2,opt,name=Info" json:"Info,omitempty"`
	Cookie string      `protobuf:"bytes,3,opt,name=Cookie" json:"Cookie,omitempty"`
}

func (m *WebDownRegisterPassword) Reset()                    { *m = WebDownRegisterPassword{} }
func (m *WebDownRegisterPassword) String() string            { return proto.CompactTextString(m) }
func (*WebDownRegisterPassword) ProtoMessage()               {}
func (*WebDownRegisterPassword) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{18} }

func (m *WebDownRegisterPassword) GetInfo() *ServerInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func init() {
	proto.RegisterType((*TestToAgent)(nil), "pb.TestToAgent")
	proto.RegisterType((*TestToGame)(nil), "pb.TestToGame")
	proto.RegisterType((*UpPing)(nil), "pb.UpPing")
	proto.RegisterType((*DownPing)(nil), "pb.DownPing")
	proto.RegisterType((*WebUpRegister)(nil), "pb.WebUpRegister")
	proto.RegisterType((*WebDownRegister)(nil), "pb.WebDownRegister")
	proto.RegisterType((*WebUpGetServerInfo)(nil), "pb.WebUpGetServerInfo")
	proto.RegisterType((*WebDownGetServerInfo)(nil), "pb.WebDownGetServerInfo")
	proto.RegisterType((*ServerInfo)(nil), "pb.ServerInfo")
	proto.RegisterType((*UpGetConfigToCenter)(nil), "pb.UpGetConfigToCenter")
	proto.RegisterType((*SyncConfigDataToAll)(nil), "pb.SyncConfigDataToAll")
	proto.RegisterType((*LoginKickPlayerToAgent)(nil), "pb.LoginKickPlayerToAgent")
	proto.RegisterType((*AgentCheckUserToLogin)(nil), "pb.AgentCheckUserToLogin")
	proto.RegisterType((*LoginCheckPlayerToAgent)(nil), "pb.LoginCheckPlayerToAgent")
	proto.RegisterType((*PlayerLogin)(nil), "pb.PlayerLogin")
	proto.RegisterType((*PlayerLoginToGame)(nil), "pb.PlayerLoginToGame")
	proto.RegisterType((*DownPlayerLogin)(nil), "pb.DownPlayerLogin")
	proto.RegisterType((*WebUpRegisterPassword)(nil), "pb.WebUpRegisterPassword")
	proto.RegisterType((*WebDownRegisterPassword)(nil), "pb.WebDownRegisterPassword")
}

func init() { proto.RegisterFile("core.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 504 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x55, 0x9c, 0x34, 0xa4, 0x13, 0x41, 0xc5, 0x96, 0xb4, 0x16, 0xe2, 0x50, 0xad, 0x84, 0x04,
	0x42, 0x8a, 0xf8, 0xb8, 0xc2, 0xa1, 0x75, 0x44, 0x64, 0x35, 0x20, 0xcb, 0x49, 0xc4, 0x11, 0x39,
	0xf6, 0xd4, 0x35, 0x71, 0xbd, 0xd6, 0xda, 0x50, 0xf5, 0xf7, 0xf0, 0x47, 0x59, 0x8f, 0xd7, 0x5f,
	0x45, 0x69, 0xab, 0xde, 0xf6, 0x3d, 0xef, 0xbc, 0x79, 0x33, 0x3b, 0x63, 0x00, 0x5f, 0x48, 0x9c,
	0xa6, 0x52, 0xe4, 0x82, 0x19, 0xe9, 0x86, 0xbf, 0x86, 0xf1, 0x0a, 0xb3, 0x7c, 0x25, 0x4e, 0x43,
	0x4c, 0x72, 0x76, 0x04, 0xc3, 0xaf, 0x11, 0xc6, 0xc1, 0xa9, 0xd9, 0x3b, 0xe9, 0xbd, 0xd9, 0x77,
	0x35, 0xe2, 0x9f, 0x01, 0xca, 0x6b, 0x73, 0xef, 0x0a, 0x77, 0xdd, 0xaa, 0xf9, 0x33, 0xd3, 0x68,
	0xf1, 0x67, 0x7c, 0x04, 0xc3, 0x75, 0xea, 0x44, 0x49, 0xc8, 0x01, 0x46, 0x33, 0x71, 0x9d, 0xd0,
	0xf9, 0x03, 0x3c, 0xfd, 0x81, 0x9b, 0x75, 0xea, 0x62, 0x18, 0x65, 0x39, 0x4a, 0x76, 0x02, 0xe3,
	0x6f, 0x9e, 0x7f, 0x19, 0x25, 0x68, 0x89, 0x00, 0xb5, 0x76, 0x9b, 0xe2, 0x7f, 0x7b, 0x70, 0xa0,
	0x62, 0x0a, 0x89, 0x3a, 0x4a, 0x25, 0x5d, 0x67, 0x28, 0xed, 0x59, 0x65, 0xa6, 0x44, 0xec, 0x25,
	0x8c, 0x1c, 0x2f, 0xcb, 0xae, 0x85, 0x0c, 0xb4, 0x9d, 0x1a, 0x33, 0x0e, 0x03, 0x3b, 0xb9, 0x10,
	0x66, 0x5f, 0xf1, 0xe3, 0x8f, 0xcf, 0xa6, 0xe9, 0x66, 0xba, 0x44, 0xf9, 0x47, 0xc5, 0x29, 0xd6,
	0xa5, 0x6f, 0x85, 0xae, 0x25, 0xc4, 0x36, 0x42, 0x73, 0x50, 0xea, 0x96, 0xe8, 0xb6, 0xcb, 0xbd,
	0xff, 0x5d, 0x2e, 0x80, 0x51, 0x61, 0x73, 0xcc, 0x1b, 0xd5, 0xc2, 0x4f, 0xe1, 0xec, 0xbb, 0x6a,
	0xa0, 0x76, 0x5a, 0xe3, 0xbb, 0xbc, 0xf2, 0x5f, 0xf0, 0x42, 0x97, 0xdc, 0xd5, 0xab, 0x6a, 0xe8,
	0xdd, 0x5d, 0x83, 0xee, 0x8d, 0xd1, 0xe9, 0x4d, 0x53, 0x5b, 0xbf, 0x5d, 0x1b, 0x9f, 0x01, 0xb4,
	0x32, 0x98, 0xf0, 0x84, 0xa6, 0xc2, 0x76, 0xb4, 0xe1, 0x0a, 0xb2, 0x57, 0xb0, 0x4f, 0x47, 0x47,
	0xc8, 0x5c, 0x4b, 0x37, 0x04, 0x9f, 0xc0, 0x21, 0x15, 0x6f, 0x89, 0xe4, 0x22, 0x0a, 0x57, 0xc2,
	0x52, 0x1f, 0x50, 0xf2, 0xb7, 0x70, 0xb8, 0xbc, 0x49, 0xfc, 0x92, 0x9d, 0x79, 0xb9, 0xa7, 0x86,
	0x2e, 0x8e, 0x19, 0x83, 0x41, 0xa0, 0x80, 0x4e, 0x41, 0x67, 0xfe, 0x1e, 0x8e, 0x16, 0x22, 0x8c,
	0x92, 0xf3, 0xc8, 0xdf, 0x3a, 0xb1, 0x77, 0x83, 0xb2, 0x35, 0xa0, 0xad, 0xd7, 0xee, 0x57, 0x15,
	0xf1, 0x39, 0x4c, 0xe8, 0x82, 0x75, 0x89, 0xfe, 0xb6, 0xe0, 0x56, 0x82, 0x14, 0x76, 0x05, 0xb4,
	0x5a, 0x60, 0x74, 0x5a, 0x70, 0x0e, 0xc7, 0x14, 0x48, 0x42, 0x0f, 0xca, 0x5d, 0xf4, 0xc9, 0xce,
	0x16, 0x18, 0x7a, 0x31, 0x69, 0x8d, 0xdc, 0x0a, 0xf2, 0x2f, 0x30, 0x2e, 0x25, 0x1e, 0xe7, 0xe5,
	0x1d, 0x3c, 0x6f, 0x85, 0x37, 0xcb, 0xb7, 0xa3, 0x03, 0x07, 0xb4, 0x5a, 0xdd, 0x7c, 0xae, 0x88,
	0x51, 0x5f, 0xdd, 0x73, 0x35, 0x2a, 0x9e, 0xcf, 0xce, 0x96, 0xbf, 0x7d, 0x1f, 0x31, 0xd0, 0x96,
	0x1b, 0x82, 0xff, 0x84, 0x49, 0x67, 0x2f, 0xeb, 0xad, 0x79, 0xcc, 0xa6, 0xa9, 0xd7, 0xa5, 0x35,
	0x29, 0xe7, 0x8c, 0xce, 0xfc, 0x0a, 0x8e, 0x6f, 0x2d, 0xf1, 0xbd, 0x29, 0xaa, 0x61, 0x37, 0x1e,
	0xb4, 0xb0, 0x9d, 0xa1, 0xde, 0x0c, 0xe9, 0x6f, 0xf7, 0xe9, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x00, 0xf0, 0x3b, 0x49, 0xfb, 0x04, 0x00, 0x00,
}