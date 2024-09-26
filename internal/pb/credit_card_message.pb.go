// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: credit_card_message.proto

package pb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BillingAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Country  string `protobuf:"bytes,1,opt,name=country,proto3" json:"country,omitempty"`
	Street   string `protobuf:"bytes,2,opt,name=street,proto3" json:"street,omitempty"`
	City     string `protobuf:"bytes,3,opt,name=city,proto3" json:"city,omitempty"`
	Postcode string `protobuf:"bytes,4,opt,name=postcode,proto3" json:"postcode,omitempty"`
}

func (x *BillingAddress) Reset() {
	*x = BillingAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_credit_card_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BillingAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BillingAddress) ProtoMessage() {}

func (x *BillingAddress) ProtoReflect() protoreflect.Message {
	mi := &file_credit_card_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BillingAddress.ProtoReflect.Descriptor instead.
func (*BillingAddress) Descriptor() ([]byte, []int) {
	return file_credit_card_message_proto_rawDescGZIP(), []int{0}
}

func (x *BillingAddress) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

func (x *BillingAddress) GetStreet() string {
	if x != nil {
		return x.Street
	}
	return ""
}

func (x *BillingAddress) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *BillingAddress) GetPostcode() string {
	if x != nil {
		return x.Postcode
	}
	return ""
}

type CreditCard struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata       string          `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Name           string          `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Month          string          `protobuf:"bytes,3,opt,name=month,proto3" json:"month,omitempty"`
	Year           string          `protobuf:"bytes,4,opt,name=year,proto3" json:"year,omitempty"`
	Number         string          `protobuf:"bytes,5,opt,name=number,proto3" json:"number,omitempty"`
	BillingAddress *BillingAddress `protobuf:"bytes,6,opt,name=billingAddress,proto3" json:"billingAddress,omitempty"`
}

func (x *CreditCard) Reset() {
	*x = CreditCard{}
	if protoimpl.UnsafeEnabled {
		mi := &file_credit_card_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreditCard) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreditCard) ProtoMessage() {}

func (x *CreditCard) ProtoReflect() protoreflect.Message {
	mi := &file_credit_card_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreditCard.ProtoReflect.Descriptor instead.
func (*CreditCard) Descriptor() ([]byte, []int) {
	return file_credit_card_message_proto_rawDescGZIP(), []int{1}
}

func (x *CreditCard) GetMetadata() string {
	if x != nil {
		return x.Metadata
	}
	return ""
}

func (x *CreditCard) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreditCard) GetMonth() string {
	if x != nil {
		return x.Month
	}
	return ""
}

func (x *CreditCard) GetYear() string {
	if x != nil {
		return x.Year
	}
	return ""
}

func (x *CreditCard) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *CreditCard) GetBillingAddress() *BillingAddress {
	if x != nil {
		return x.BillingAddress
	}
	return nil
}

var File_credit_card_message_proto protoreflect.FileDescriptor

var file_credit_card_message_proto_rawDesc = []byte{
	0x0a, 0x19, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x67, 0x6f, 0x70,
	0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x22, 0x72, 0x0a, 0x0e, 0x42, 0x69, 0x6c, 0x6c, 0x69,
	0x6e, 0x67, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63,
	0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x69, 0x74, 0x79, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x63, 0x6f, 0x64, 0x65, 0x22, 0xc2, 0x01, 0x0a, 0x0a,
	0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x43, 0x61, 0x72, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f,
	0x6e, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x6f, 0x6e, 0x74, 0x68,
	0x12, 0x12, 0x0a, 0x04, 0x79, 0x65, 0x61, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x79, 0x65, 0x61, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x42, 0x0a, 0x0e,
	0x62, 0x69, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65,
	0x72, 0x2e, 0x42, 0x69, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x52, 0x0e, 0x62, 0x69, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b,
	0x61, 0x61, 0x2d, 0x69, 0x74, 0x2f, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_credit_card_message_proto_rawDescOnce sync.Once
	file_credit_card_message_proto_rawDescData = file_credit_card_message_proto_rawDesc
)

func file_credit_card_message_proto_rawDescGZIP() []byte {
	file_credit_card_message_proto_rawDescOnce.Do(func() {
		file_credit_card_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_credit_card_message_proto_rawDescData)
	})
	return file_credit_card_message_proto_rawDescData
}

var file_credit_card_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_credit_card_message_proto_goTypes = []any{
	(*BillingAddress)(nil), // 0: gophkeeper.BillingAddress
	(*CreditCard)(nil),     // 1: gophkeeper.CreditCard
}
var file_credit_card_message_proto_depIdxs = []int32{
	0, // 0: gophkeeper.CreditCard.billingAddress:type_name -> gophkeeper.BillingAddress
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_credit_card_message_proto_init() }
func file_credit_card_message_proto_init() {
	if File_credit_card_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_credit_card_message_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*BillingAddress); i {
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
		file_credit_card_message_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreditCard); i {
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
			RawDescriptor: file_credit_card_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_credit_card_message_proto_goTypes,
		DependencyIndexes: file_credit_card_message_proto_depIdxs,
		MessageInfos:      file_credit_card_message_proto_msgTypes,
	}.Build()
	File_credit_card_message_proto = out.File
	file_credit_card_message_proto_rawDesc = nil
	file_credit_card_message_proto_goTypes = nil
	file_credit_card_message_proto_depIdxs = nil
}
