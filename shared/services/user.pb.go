// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user.proto

package services

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SessionToken struct {
	Token                string   `protobuf:"bytes,1,opt,name=Token,proto3" json:"Token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionToken) Reset()         { *m = SessionToken{} }
func (m *SessionToken) String() string { return proto.CompactTextString(m) }
func (*SessionToken) ProtoMessage()    {}
func (*SessionToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{0}
}

func (m *SessionToken) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionToken.Unmarshal(m, b)
}
func (m *SessionToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionToken.Marshal(b, m, deterministic)
}
func (m *SessionToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionToken.Merge(m, src)
}
func (m *SessionToken) XXX_Size() int {
	return xxx_messageInfo_SessionToken.Size(m)
}
func (m *SessionToken) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionToken.DiscardUnknown(m)
}

var xxx_messageInfo_SessionToken proto.InternalMessageInfo

func (m *SessionToken) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type UserSignature struct {
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserSignature) Reset()         { *m = UserSignature{} }
func (m *UserSignature) String() string { return proto.CompactTextString(m) }
func (*UserSignature) ProtoMessage()    {}
func (*UserSignature) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{1}
}

func (m *UserSignature) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserSignature.Unmarshal(m, b)
}
func (m *UserSignature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserSignature.Marshal(b, m, deterministic)
}
func (m *UserSignature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserSignature.Merge(m, src)
}
func (m *UserSignature) XXX_Size() int {
	return xxx_messageInfo_UserSignature.Size(m)
}
func (m *UserSignature) XXX_DiscardUnknown() {
	xxx_messageInfo_UserSignature.DiscardUnknown(m)
}

var xxx_messageInfo_UserSignature proto.InternalMessageInfo

func (m *UserSignature) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UserSignature) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type User struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Email                string   `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Nickname             string   `protobuf:"bytes,3,opt,name=nickname,proto3" json:"nickname,omitempty"`
	AvatarPath           string   `protobuf:"bytes,4,opt,name=avatarPath,proto3" json:"avatarPath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{2}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetNickname() string {
	if m != nil {
		return m.Nickname
	}
	return ""
}

func (m *User) GetAvatarPath() string {
	if m != nil {
		return m.AvatarPath
	}
	return ""
}

type NewUserData struct {
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Nickname             string   `protobuf:"bytes,3,opt,name=nickname,proto3" json:"nickname,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NewUserData) Reset()         { *m = NewUserData{} }
func (m *NewUserData) String() string { return proto.CompactTextString(m) }
func (*NewUserData) ProtoMessage()    {}
func (*NewUserData) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{3}
}

func (m *NewUserData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewUserData.Unmarshal(m, b)
}
func (m *NewUserData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewUserData.Marshal(b, m, deterministic)
}
func (m *NewUserData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewUserData.Merge(m, src)
}
func (m *NewUserData) XXX_Size() int {
	return xxx_messageInfo_NewUserData.Size(m)
}
func (m *NewUserData) XXX_DiscardUnknown() {
	xxx_messageInfo_NewUserData.DiscardUnknown(m)
}

var xxx_messageInfo_NewUserData proto.InternalMessageInfo

func (m *NewUserData) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *NewUserData) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *NewUserData) GetNickname() string {
	if m != nil {
		return m.Nickname
	}
	return ""
}

type PageData struct {
	Since                uint64   `protobuf:"varint,1,opt,name=since,proto3" json:"since,omitempty"`
	Limit                uint64   `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PageData) Reset()         { *m = PageData{} }
func (m *PageData) String() string { return proto.CompactTextString(m) }
func (*PageData) ProtoMessage()    {}
func (*PageData) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{4}
}

func (m *PageData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PageData.Unmarshal(m, b)
}
func (m *PageData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PageData.Marshal(b, m, deterministic)
}
func (m *PageData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PageData.Merge(m, src)
}
func (m *PageData) XXX_Size() int {
	return xxx_messageInfo_PageData.Size(m)
}
func (m *PageData) XXX_DiscardUnknown() {
	xxx_messageInfo_PageData.DiscardUnknown(m)
}

var xxx_messageInfo_PageData proto.InternalMessageInfo

func (m *PageData) GetSince() uint64 {
	if m != nil {
		return m.Since
	}
	return 0
}

func (m *PageData) GetLimit() uint64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type LeaderBoardPage struct {
	PagesCount           uint64     `protobuf:"varint,1,opt,name=pagesCount,proto3" json:"pagesCount,omitempty"`
	Leaders              []*Profile `protobuf:"bytes,2,rep,name=leaders,proto3" json:"leaders,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *LeaderBoardPage) Reset()         { *m = LeaderBoardPage{} }
func (m *LeaderBoardPage) String() string { return proto.CompactTextString(m) }
func (*LeaderBoardPage) ProtoMessage()    {}
func (*LeaderBoardPage) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{5}
}

func (m *LeaderBoardPage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LeaderBoardPage.Unmarshal(m, b)
}
func (m *LeaderBoardPage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LeaderBoardPage.Marshal(b, m, deterministic)
}
func (m *LeaderBoardPage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaderBoardPage.Merge(m, src)
}
func (m *LeaderBoardPage) XXX_Size() int {
	return xxx_messageInfo_LeaderBoardPage.Size(m)
}
func (m *LeaderBoardPage) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaderBoardPage.DiscardUnknown(m)
}

var xxx_messageInfo_LeaderBoardPage proto.InternalMessageInfo

func (m *LeaderBoardPage) GetPagesCount() uint64 {
	if m != nil {
		return m.PagesCount
	}
	return 0
}

func (m *LeaderBoardPage) GetLeaders() []*Profile {
	if m != nil {
		return m.Leaders
	}
	return nil
}

type Profile struct {
	User                 *User    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	WinRate              float32  `protobuf:"fixed32,2,opt,name=winRate,proto3" json:"winRate,omitempty"`
	Rating               uint32   `protobuf:"varint,3,opt,name=rating,proto3" json:"rating,omitempty"`
	Matches              uint32   `protobuf:"varint,4,opt,name=matches,proto3" json:"matches,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Profile) Reset()         { *m = Profile{} }
func (m *Profile) String() string { return proto.CompactTextString(m) }
func (*Profile) ProtoMessage()    {}
func (*Profile) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{6}
}

func (m *Profile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Profile.Unmarshal(m, b)
}
func (m *Profile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Profile.Marshal(b, m, deterministic)
}
func (m *Profile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Profile.Merge(m, src)
}
func (m *Profile) XXX_Size() int {
	return xxx_messageInfo_Profile.Size(m)
}
func (m *Profile) XXX_DiscardUnknown() {
	xxx_messageInfo_Profile.DiscardUnknown(m)
}

var xxx_messageInfo_Profile proto.InternalMessageInfo

func (m *Profile) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *Profile) GetWinRate() float32 {
	if m != nil {
		return m.WinRate
	}
	return 0
}

func (m *Profile) GetRating() uint32 {
	if m != nil {
		return m.Rating
	}
	return 0
}

func (m *Profile) GetMatches() uint32 {
	if m != nil {
		return m.Matches
	}
	return 0
}

type Nothing struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Nothing) Reset()         { *m = Nothing{} }
func (m *Nothing) String() string { return proto.CompactTextString(m) }
func (*Nothing) ProtoMessage()    {}
func (*Nothing) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{7}
}

func (m *Nothing) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Nothing.Unmarshal(m, b)
}
func (m *Nothing) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Nothing.Marshal(b, m, deterministic)
}
func (m *Nothing) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nothing.Merge(m, src)
}
func (m *Nothing) XXX_Size() int {
	return xxx_messageInfo_Nothing.Size(m)
}
func (m *Nothing) XXX_DiscardUnknown() {
	xxx_messageInfo_Nothing.DiscardUnknown(m)
}

var xxx_messageInfo_Nothing proto.InternalMessageInfo

type UserId struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserId) Reset()         { *m = UserId{} }
func (m *UserId) String() string { return proto.CompactTextString(m) }
func (*UserId) ProtoMessage()    {}
func (*UserId) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{8}
}

func (m *UserId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserId.Unmarshal(m, b)
}
func (m *UserId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserId.Marshal(b, m, deterministic)
}
func (m *UserId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserId.Merge(m, src)
}
func (m *UserId) XXX_Size() int {
	return xxx_messageInfo_UserId.Size(m)
}
func (m *UserId) XXX_DiscardUnknown() {
	xxx_messageInfo_UserId.DiscardUnknown(m)
}

var xxx_messageInfo_UserId proto.InternalMessageInfo

func (m *UserId) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type RoomSettings struct {
	MaxConnections       uint64   `protobuf:"varint,1,opt,name=maxConnections,proto3" json:"maxConnections,omitempty"`
	Talkers              []uint64 `protobuf:"varint,2,rep,packed,name=Talkers,proto3" json:"Talkers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomSettings) Reset()         { *m = RoomSettings{} }
func (m *RoomSettings) String() string { return proto.CompactTextString(m) }
func (*RoomSettings) ProtoMessage()    {}
func (*RoomSettings) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{9}
}

func (m *RoomSettings) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomSettings.Unmarshal(m, b)
}
func (m *RoomSettings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomSettings.Marshal(b, m, deterministic)
}
func (m *RoomSettings) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomSettings.Merge(m, src)
}
func (m *RoomSettings) XXX_Size() int {
	return xxx_messageInfo_RoomSettings.Size(m)
}
func (m *RoomSettings) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomSettings.DiscardUnknown(m)
}

var xxx_messageInfo_RoomSettings proto.InternalMessageInfo

func (m *RoomSettings) GetMaxConnections() uint64 {
	if m != nil {
		return m.MaxConnections
	}
	return 0
}

func (m *RoomSettings) GetTalkers() []uint64 {
	if m != nil {
		return m.Talkers
	}
	return nil
}

type RoomId struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoomId) Reset()         { *m = RoomId{} }
func (m *RoomId) String() string { return proto.CompactTextString(m) }
func (*RoomId) ProtoMessage()    {}
func (*RoomId) Descriptor() ([]byte, []int) {
	return fileDescriptor_116e343673f7ffaf, []int{10}
}

func (m *RoomId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomId.Unmarshal(m, b)
}
func (m *RoomId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomId.Marshal(b, m, deterministic)
}
func (m *RoomId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomId.Merge(m, src)
}
func (m *RoomId) XXX_Size() int {
	return xxx_messageInfo_RoomId.Size(m)
}
func (m *RoomId) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomId.DiscardUnknown(m)
}

var xxx_messageInfo_RoomId proto.InternalMessageInfo

func (m *RoomId) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func init() {
	proto.RegisterType((*SessionToken)(nil), "services.SessionToken")
	proto.RegisterType((*UserSignature)(nil), "services.UserSignature")
	proto.RegisterType((*User)(nil), "services.User")
	proto.RegisterType((*NewUserData)(nil), "services.NewUserData")
	proto.RegisterType((*PageData)(nil), "services.PageData")
	proto.RegisterType((*LeaderBoardPage)(nil), "services.LeaderBoardPage")
	proto.RegisterType((*Profile)(nil), "services.Profile")
	proto.RegisterType((*Nothing)(nil), "services.Nothing")
	proto.RegisterType((*UserId)(nil), "services.UserId")
	proto.RegisterType((*RoomSettings)(nil), "services.RoomSettings")
	proto.RegisterType((*RoomId)(nil), "services.RoomId")
}

func init() { proto.RegisterFile("user.proto", fileDescriptor_116e343673f7ffaf) }

var fileDescriptor_116e343673f7ffaf = []byte{
	// 584 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xd1, 0x4b, 0x1b, 0x4f,
	0x10, 0xd6, 0x98, 0x5f, 0xa2, 0x13, 0xe3, 0xcf, 0x2e, 0xad, 0xbd, 0xde, 0x43, 0x91, 0xa5, 0x14,
	0xa1, 0x20, 0x34, 0x52, 0xe9, 0x53, 0xa1, 0xc6, 0x22, 0x81, 0x2a, 0xe1, 0xa2, 0x4f, 0x85, 0xc2,
	0x34, 0x37, 0xcd, 0x2d, 0xb9, 0xdb, 0x0d, 0xbb, 0xab, 0xa9, 0xff, 0x62, 0xff, 0xaa, 0xb2, 0x7b,
	0x77, 0xe6, 0x72, 0xa7, 0x85, 0xbe, 0xe5, 0x9b, 0x99, 0xef, 0xfb, 0x66, 0x67, 0x26, 0x07, 0x70,
	0x6b, 0x48, 0x1f, 0x2f, 0xb4, 0xb2, 0x8a, 0x6d, 0x1b, 0xd2, 0x77, 0x62, 0x4a, 0x86, 0xbf, 0x81,
	0xdd, 0x09, 0x19, 0x23, 0x94, 0xbc, 0x56, 0x73, 0x92, 0xec, 0x39, 0xfc, 0xe7, 0x7f, 0x04, 0x9b,
	0x87, 0x9b, 0x47, 0x3b, 0x51, 0x0e, 0xf8, 0x67, 0xe8, 0xdf, 0x18, 0xd2, 0x13, 0x31, 0x93, 0x68,
	0x6f, 0x35, 0xb9, 0x32, 0xca, 0x50, 0xa4, 0x65, 0x99, 0x07, 0x2c, 0x84, 0xed, 0x05, 0x1a, 0xb3,
	0x54, 0x3a, 0x0e, 0x5a, 0x3e, 0xf1, 0x80, 0x79, 0x02, 0x6d, 0x27, 0xc1, 0xf6, 0xa0, 0x25, 0x62,
	0x4f, 0x6b, 0x47, 0x2d, 0x11, 0xaf, 0x94, 0x5a, 0x35, 0x25, 0x29, 0xa6, 0x73, 0x89, 0x19, 0x05,
	0x5b, 0xb9, 0x52, 0x89, 0xd9, 0x6b, 0x00, 0xbc, 0x43, 0x8b, 0x7a, 0x8c, 0x36, 0x09, 0xda, 0x3e,
	0x5b, 0x89, 0xf0, 0x6f, 0xd0, 0xbb, 0xa2, 0xa5, 0x33, 0x3b, 0x47, 0x8b, 0xff, 0xde, 0xea, 0xdf,
	0xcc, 0xf9, 0x29, 0x6c, 0x8f, 0x71, 0x46, 0xa5, 0xb2, 0x11, 0x72, 0x4a, 0xc5, 0x6b, 0x72, 0xe0,
	0xa2, 0xa9, 0xc8, 0x84, 0xf5, 0xb2, 0xed, 0x28, 0x07, 0xfc, 0x3b, 0xfc, 0xff, 0x95, 0x30, 0x26,
	0x7d, 0xa6, 0x50, 0xc7, 0x4e, 0xc2, 0xbd, 0x63, 0x81, 0x33, 0x32, 0x43, 0x75, 0x2b, 0x6d, 0xa1,
	0x51, 0x89, 0xb0, 0x77, 0xd0, 0x4d, 0x3d, 0xc5, 0x04, 0xad, 0xc3, 0xad, 0xa3, 0xde, 0xe0, 0xd9,
	0x71, 0xb9, 0xb6, 0xe3, 0xb1, 0x56, 0x3f, 0x45, 0x4a, 0x51, 0x59, 0xc1, 0xef, 0xa1, 0x5b, 0xc4,
	0x18, 0x87, 0xb6, 0x5b, 0xb5, 0x57, 0xec, 0x0d, 0xf6, 0x56, 0x24, 0x37, 0x92, 0xc8, 0xe7, 0x58,
	0x00, 0xdd, 0xa5, 0x90, 0x11, 0x5a, 0xf2, 0x6d, 0xb6, 0xa2, 0x12, 0xb2, 0x03, 0xe8, 0x68, 0xb4,
	0x42, 0xce, 0xfc, 0xd3, 0xfb, 0x51, 0x81, 0x1c, 0x23, 0x43, 0x3b, 0x4d, 0xc8, 0xf8, 0x91, 0xf7,
	0xa3, 0x12, 0xf2, 0x1d, 0xe8, 0x5e, 0x29, 0x9b, 0x08, 0x39, 0xe3, 0x01, 0x74, 0x9c, 0xc9, 0x28,
	0xae, 0xaf, 0x99, 0x8f, 0x61, 0x37, 0x52, 0x2a, 0x9b, 0x90, 0x75, 0x6a, 0x86, 0xbd, 0x85, 0xbd,
	0x0c, 0x7f, 0x0d, 0x95, 0x94, 0x34, 0xb5, 0x42, 0x49, 0x53, 0xd4, 0xd6, 0xa2, 0xce, 0xf6, 0x1a,
	0xd3, 0x79, 0x39, 0x84, 0x76, 0x54, 0x42, 0xe7, 0xe5, 0x14, 0x9b, 0x5e, 0x83, 0xdf, 0x5b, 0x79,
	0x1b, 0x97, 0x13, 0x76, 0x0a, 0x30, 0x4c, 0x68, 0x3a, 0xcf, 0x8f, 0xfb, 0x60, 0x35, 0x8b, 0xea,
	0xd1, 0x87, 0xb5, 0x19, 0xf1, 0x0d, 0xf6, 0x09, 0x76, 0x2e, 0x71, 0x4e, 0x39, 0xed, 0xe5, 0x7a,
	0xfa, 0xe1, 0x5f, 0x10, 0x3e, 0xa1, 0xc7, 0x37, 0xd8, 0x07, 0x80, 0xa1, 0x26, 0xb4, 0xe4, 0x6f,
	0xfe, 0xc5, 0xaa, 0xae, 0x72, 0x99, 0x8f, 0xd8, 0xbe, 0x07, 0xb8, 0x20, 0x5b, 0x2e, 0xb2, 0x96,
	0x0f, 0x9b, 0xfb, 0xf7, 0x94, 0xfe, 0xcd, 0x22, 0x46, 0x4b, 0x4f, 0xb1, 0x9a, 0x2e, 0x5f, 0x80,
	0x5d, 0x90, 0xad, 0x9f, 0x23, 0xab, 0xa8, 0x17, 0x17, 0x1e, 0xbe, 0x5a, 0xc5, 0x6a, 0xe5, 0xde,
	0xb9, 0x3b, 0x4c, 0x09, 0xf5, 0xf9, 0x19, 0xab, 0x74, 0x56, 0x9c, 0x42, 0xd8, 0x0c, 0x79, 0x4a,
	0xef, 0x82, 0xac, 0x6b, 0xe3, 0xec, 0x7e, 0x14, 0xb3, 0xfd, 0xf5, 0xd6, 0x46, 0x71, 0xb3, 0xd9,
	0xc1, 0x12, 0x3a, 0xc3, 0x04, 0xed, 0xe5, 0x84, 0x7d, 0x2c, 0x67, 0xea, 0xd6, 0x5e, 0xdd, 0x65,
	0xf5, 0xb0, 0xc2, 0xfd, 0xf5, 0xf8, 0x28, 0xe6, 0x1b, 0xec, 0x04, 0xe0, 0x9c, 0x52, 0x2a, 0x98,
	0x8d, 0x8a, 0x47, 0x7b, 0xfd, 0xd1, 0xf1, 0x9f, 0xca, 0x93, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff,
	0xe5, 0xbd, 0x87, 0x3d, 0x38, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UserMSClient is the client API for UserMS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserMSClient interface {
	CheckToken(ctx context.Context, in *SessionToken, opts ...grpc.CallOption) (*User, error)
	MakeToken(ctx context.Context, in *UserSignature, opts ...grpc.CallOption) (*SessionToken, error)
	CreateUser(ctx context.Context, in *NewUserData, opts ...grpc.CallOption) (*User, error)
	GetProfile(ctx context.Context, in *User, opts ...grpc.CallOption) (*Profile, error)
	UpdateProfile(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
	GetLeaderBoardPage(ctx context.Context, in *PageData, opts ...grpc.CallOption) (*LeaderBoardPage, error)
	ClearDB(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
	GetUserById(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*User, error)
}

type userMSClient struct {
	cc *grpc.ClientConn
}

func NewUserMSClient(cc *grpc.ClientConn) UserMSClient {
	return &userMSClient{cc}
}

func (c *userMSClient) CheckToken(ctx context.Context, in *SessionToken, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/services.UserMS/CheckToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userMSClient) MakeToken(ctx context.Context, in *UserSignature, opts ...grpc.CallOption) (*SessionToken, error) {
	out := new(SessionToken)
	err := c.cc.Invoke(ctx, "/services.UserMS/MakeToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userMSClient) CreateUser(ctx context.Context, in *NewUserData, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/services.UserMS/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userMSClient) GetProfile(ctx context.Context, in *User, opts ...grpc.CallOption) (*Profile, error) {
	out := new(Profile)
	err := c.cc.Invoke(ctx, "/services.UserMS/GetProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userMSClient) UpdateProfile(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/services.UserMS/UpdateProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userMSClient) GetLeaderBoardPage(ctx context.Context, in *PageData, opts ...grpc.CallOption) (*LeaderBoardPage, error) {
	out := new(LeaderBoardPage)
	err := c.cc.Invoke(ctx, "/services.UserMS/GetLeaderBoardPage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userMSClient) ClearDB(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/services.UserMS/ClearDB", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userMSClient) GetUserById(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/services.UserMS/GetUserById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserMSServer is the server API for UserMS service.
type UserMSServer interface {
	CheckToken(context.Context, *SessionToken) (*User, error)
	MakeToken(context.Context, *UserSignature) (*SessionToken, error)
	CreateUser(context.Context, *NewUserData) (*User, error)
	GetProfile(context.Context, *User) (*Profile, error)
	UpdateProfile(context.Context, *User) (*User, error)
	GetLeaderBoardPage(context.Context, *PageData) (*LeaderBoardPage, error)
	ClearDB(context.Context, *Nothing) (*Nothing, error)
	GetUserById(context.Context, *UserId) (*User, error)
}

func RegisterUserMSServer(s *grpc.Server, srv UserMSServer) {
	s.RegisterService(&_UserMS_serviceDesc, srv)
}

func _UserMS_CheckToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserMSServer).CheckToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.UserMS/CheckToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserMSServer).CheckToken(ctx, req.(*SessionToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserMS_MakeToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserSignature)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserMSServer).MakeToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.UserMS/MakeToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserMSServer).MakeToken(ctx, req.(*UserSignature))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserMS_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewUserData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserMSServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.UserMS/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserMSServer).CreateUser(ctx, req.(*NewUserData))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserMS_GetProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserMSServer).GetProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.UserMS/GetProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserMSServer).GetProfile(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserMS_UpdateProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserMSServer).UpdateProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.UserMS/UpdateProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserMSServer).UpdateProfile(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserMS_GetLeaderBoardPage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PageData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserMSServer).GetLeaderBoardPage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.UserMS/GetLeaderBoardPage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserMSServer).GetLeaderBoardPage(ctx, req.(*PageData))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserMS_ClearDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nothing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserMSServer).ClearDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.UserMS/ClearDB",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserMSServer).ClearDB(ctx, req.(*Nothing))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserMS_GetUserById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserMSServer).GetUserById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.UserMS/GetUserById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserMSServer).GetUserById(ctx, req.(*UserId))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserMS_serviceDesc = grpc.ServiceDesc{
	ServiceName: "services.UserMS",
	HandlerType: (*UserMSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckToken",
			Handler:    _UserMS_CheckToken_Handler,
		},
		{
			MethodName: "MakeToken",
			Handler:    _UserMS_MakeToken_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _UserMS_CreateUser_Handler,
		},
		{
			MethodName: "GetProfile",
			Handler:    _UserMS_GetProfile_Handler,
		},
		{
			MethodName: "UpdateProfile",
			Handler:    _UserMS_UpdateProfile_Handler,
		},
		{
			MethodName: "GetLeaderBoardPage",
			Handler:    _UserMS_GetLeaderBoardPage_Handler,
		},
		{
			MethodName: "ClearDB",
			Handler:    _UserMS_ClearDB_Handler,
		},
		{
			MethodName: "GetUserById",
			Handler:    _UserMS_GetUserById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}

// ChatMSClient is the client API for ChatMS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChatMSClient interface {
	CreateRoom(ctx context.Context, in *RoomSettings, opts ...grpc.CallOption) (*RoomId, error)
	DeleteRoom(ctx context.Context, in *RoomId, opts ...grpc.CallOption) (*Nothing, error)
}

type chatMSClient struct {
	cc *grpc.ClientConn
}

func NewChatMSClient(cc *grpc.ClientConn) ChatMSClient {
	return &chatMSClient{cc}
}

func (c *chatMSClient) CreateRoom(ctx context.Context, in *RoomSettings, opts ...grpc.CallOption) (*RoomId, error) {
	out := new(RoomId)
	err := c.cc.Invoke(ctx, "/services.ChatMS/CreateRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatMSClient) DeleteRoom(ctx context.Context, in *RoomId, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/services.ChatMS/DeleteRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatMSServer is the server API for ChatMS service.
type ChatMSServer interface {
	CreateRoom(context.Context, *RoomSettings) (*RoomId, error)
	DeleteRoom(context.Context, *RoomId) (*Nothing, error)
}

func RegisterChatMSServer(s *grpc.Server, srv ChatMSServer) {
	s.RegisterService(&_ChatMS_serviceDesc, srv)
}

func _ChatMS_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomSettings)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatMSServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.ChatMS/CreateRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatMSServer).CreateRoom(ctx, req.(*RoomSettings))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatMS_DeleteRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatMSServer).DeleteRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.ChatMS/DeleteRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatMSServer).DeleteRoom(ctx, req.(*RoomId))
	}
	return interceptor(ctx, in, info, handler)
}

var _ChatMS_serviceDesc = grpc.ServiceDesc{
	ServiceName: "services.ChatMS",
	HandlerType: (*ChatMSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRoom",
			Handler:    _ChatMS_CreateRoom_Handler,
		},
		{
			MethodName: "DeleteRoom",
			Handler:    _ChatMS_DeleteRoom_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
