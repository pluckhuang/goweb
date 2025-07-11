// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: user/v1/user.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Gender int32

const (
	Gender_GENDER_UNKNOWN Gender = 0
	Gender_GENDER_MALE    Gender = 1
	Gender_GENDER_FEMALE  Gender = 2
)

// Enum value maps for Gender.
var (
	Gender_name = map[int32]string{
		0: "GENDER_UNKNOWN",
		1: "GENDER_MALE",
		2: "GENDER_FEMALE",
	}
	Gender_value = map[string]int32{
		"GENDER_UNKNOWN": 0,
		"GENDER_MALE":    1,
		"GENDER_FEMALE":  2,
	}
)

func (x Gender) Enum() *Gender {
	p := new(Gender)
	*p = x
	return p
}

func (x Gender) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Gender) Descriptor() protoreflect.EnumDescriptor {
	return file_user_v1_user_proto_enumTypes[0].Descriptor()
}

func (Gender) Type() protoreflect.EnumType {
	return &file_user_v1_user_proto_enumTypes[0]
}

func (x Gender) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Gender.Descriptor instead.
func (Gender) EnumDescriptor() ([]byte, []int) {
	return file_user_v1_user_proto_rawDescGZIP(), []int{0}
}

type User struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Avatar        string                 `protobuf:"bytes,4,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Attributes    map[string]string      `protobuf:"bytes,5,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Nicknames     []string               `protobuf:"bytes,6,rep,name=nicknames,proto3" json:"nicknames,omitempty"`
	Address       *Address               `protobuf:"bytes,7,opt,name=address,proto3" json:"address,omitempty"`
	Gender        Gender                 `protobuf:"varint,8,opt,name=gender,proto3,enum=Gender" json:"gender,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_user_v1_user_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_user_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_user_v1_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *User) GetAttributes() map[string]string {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *User) GetNicknames() []string {
	if x != nil {
		return x.Nicknames
	}
	return nil
}

func (x *User) GetAddress() *Address {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *User) GetGender() Gender {
	if x != nil {
		return x.Gender
	}
	return Gender_GENDER_UNKNOWN
}

type Address struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Province      string                 `protobuf:"bytes,1,opt,name=province,proto3" json:"province,omitempty"`
	City          string                 `protobuf:"bytes,3,opt,name=city,proto3" json:"city,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Address) Reset() {
	*x = Address{}
	mi := &file_user_v1_user_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_user_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_user_v1_user_proto_rawDescGZIP(), []int{1}
}

func (x *Address) GetProvince() string {
	if x != nil {
		return x.Province
	}
	return ""
}

func (x *Address) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

type GetByIDRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetByIDRequest) Reset() {
	*x = GetByIDRequest{}
	mi := &file_user_v1_user_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIDRequest) ProtoMessage() {}

func (x *GetByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_user_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIDRequest.ProtoReflect.Descriptor instead.
func (*GetByIDRequest) Descriptor() ([]byte, []int) {
	return file_user_v1_user_proto_rawDescGZIP(), []int{2}
}

func (x *GetByIDRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetByIDResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetByIDResponse) Reset() {
	*x = GetByIDResponse{}
	mi := &file_user_v1_user_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIDResponse) ProtoMessage() {}

func (x *GetByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_user_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIDResponse.ProtoReflect.Descriptor instead.
func (*GetByIDResponse) Descriptor() ([]byte, []int) {
	return file_user_v1_user_proto_rawDescGZIP(), []int{3}
}

func (x *GetByIDResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

var File_user_v1_user_proto protoreflect.FileDescriptor

const file_user_v1_user_proto_rawDesc = "" +
	"\n" +
	"\x12user/v1/user.proto\"\x9b\x02\n" +
	"\x04User\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x16\n" +
	"\x06avatar\x18\x04 \x01(\tR\x06avatar\x125\n" +
	"\n" +
	"attributes\x18\x05 \x03(\v2\x15.User.AttributesEntryR\n" +
	"attributes\x12\x1c\n" +
	"\tnicknames\x18\x06 \x03(\tR\tnicknames\x12\"\n" +
	"\aaddress\x18\a \x01(\v2\b.AddressR\aaddress\x12\x1f\n" +
	"\x06gender\x18\b \x01(\x0e2\a.GenderR\x06gender\x1a=\n" +
	"\x0fAttributesEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\"9\n" +
	"\aAddress\x12\x1a\n" +
	"\bprovince\x18\x01 \x01(\tR\bprovince\x12\x12\n" +
	"\x04city\x18\x03 \x01(\tR\x04city\" \n" +
	"\x0eGetByIDRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\",\n" +
	"\x0fGetByIDResponse\x12\x19\n" +
	"\x04user\x18\x01 \x01(\v2\x05.UserR\x04user*@\n" +
	"\x06Gender\x12\x12\n" +
	"\x0eGENDER_UNKNOWN\x10\x00\x12\x0f\n" +
	"\vGENDER_MALE\x10\x01\x12\x11\n" +
	"\rGENDER_FEMALE\x10\x022;\n" +
	"\vUserService\x12,\n" +
	"\aGetByID\x12\x0f.GetByIDRequest\x1a\x10.GetByIDResponseBEB\tUserProtoP\x01Z6github.com/pluckhuang/goweb/aweb/api/proto/gen/user/v1b\x06proto3"

var (
	file_user_v1_user_proto_rawDescOnce sync.Once
	file_user_v1_user_proto_rawDescData []byte
)

func file_user_v1_user_proto_rawDescGZIP() []byte {
	file_user_v1_user_proto_rawDescOnce.Do(func() {
		file_user_v1_user_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_user_v1_user_proto_rawDesc), len(file_user_v1_user_proto_rawDesc)))
	})
	return file_user_v1_user_proto_rawDescData
}

var file_user_v1_user_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_user_v1_user_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_user_v1_user_proto_goTypes = []any{
	(Gender)(0),             // 0: Gender
	(*User)(nil),            // 1: User
	(*Address)(nil),         // 2: Address
	(*GetByIDRequest)(nil),  // 3: GetByIDRequest
	(*GetByIDResponse)(nil), // 4: GetByIDResponse
	nil,                     // 5: User.AttributesEntry
}
var file_user_v1_user_proto_depIdxs = []int32{
	5, // 0: User.attributes:type_name -> User.AttributesEntry
	2, // 1: User.address:type_name -> Address
	0, // 2: User.gender:type_name -> Gender
	1, // 3: GetByIDResponse.user:type_name -> User
	3, // 4: UserService.GetByID:input_type -> GetByIDRequest
	4, // 5: UserService.GetByID:output_type -> GetByIDResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_user_v1_user_proto_init() }
func file_user_v1_user_proto_init() {
	if File_user_v1_user_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_user_v1_user_proto_rawDesc), len(file_user_v1_user_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_v1_user_proto_goTypes,
		DependencyIndexes: file_user_v1_user_proto_depIdxs,
		EnumInfos:         file_user_v1_user_proto_enumTypes,
		MessageInfos:      file_user_v1_user_proto_msgTypes,
	}.Build()
	File_user_v1_user_proto = out.File
	file_user_v1_user_proto_goTypes = nil
	file_user_v1_user_proto_depIdxs = nil
}
