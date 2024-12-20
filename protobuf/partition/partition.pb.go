// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: partition.proto

package __

import (
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

type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Size       int64  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	ChunkSize  int64  `protobuf:"varint,4,opt,name=chunkSize,proto3" json:"chunkSize,omitempty"`
	NbReplicas int32  `protobuf:"varint,3,opt,name=nbReplicas,proto3" json:"nbReplicas,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	mi := &file_partition_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_partition_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_partition_proto_rawDescGZIP(), []int{0}
}

func (x *File) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *File) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *File) GetChunkSize() int64 {
	if x != nil {
		return x.ChunkSize
	}
	return 0
}

func (x *File) GetNbReplicas() int32 {
	if x != nil {
		return x.NbReplicas
	}
	return 0
}

type FilePartition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileUuid string   `protobuf:"bytes,1,opt,name=fileUuid,proto3" json:"fileUuid,omitempty"`
	Chunks   []*Chunk `protobuf:"bytes,2,rep,name=chunks,proto3" json:"chunks,omitempty"`
}

func (x *FilePartition) Reset() {
	*x = FilePartition{}
	mi := &file_partition_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FilePartition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilePartition) ProtoMessage() {}

func (x *FilePartition) ProtoReflect() protoreflect.Message {
	mi := &file_partition_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilePartition.ProtoReflect.Descriptor instead.
func (*FilePartition) Descriptor() ([]byte, []int) {
	return file_partition_proto_rawDescGZIP(), []int{1}
}

func (x *FilePartition) GetFileUuid() string {
	if x != nil {
		return x.FileUuid
	}
	return ""
}

func (x *FilePartition) GetChunks() []*Chunk {
	if x != nil {
		return x.Chunks
	}
	return nil
}

type Chunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Offset       int64    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Size         int64    `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	SendTo       string   `protobuf:"bytes,4,opt,name=sendTo,proto3" json:"sendTo,omitempty"`
	ReplicaChain []string `protobuf:"bytes,5,rep,name=replicaChain,proto3" json:"replicaChain,omitempty"`
}

func (x *Chunk) Reset() {
	*x = Chunk{}
	mi := &file_partition_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Chunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chunk) ProtoMessage() {}

func (x *Chunk) ProtoReflect() protoreflect.Message {
	mi := &file_partition_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chunk.ProtoReflect.Descriptor instead.
func (*Chunk) Descriptor() ([]byte, []int) {
	return file_partition_proto_rawDescGZIP(), []int{2}
}

func (x *Chunk) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Chunk) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *Chunk) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Chunk) GetSendTo() string {
	if x != nil {
		return x.SendTo
	}
	return ""
}

func (x *Chunk) GetReplicaChain() []string {
	if x != nil {
		return x.ReplicaChain
	}
	return nil
}

type FileDesc struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename string `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
}

func (x *FileDesc) Reset() {
	*x = FileDesc{}
	mi := &file_partition_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileDesc) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileDesc) ProtoMessage() {}

func (x *FileDesc) ProtoReflect() protoreflect.Message {
	mi := &file_partition_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileDesc.ProtoReflect.Descriptor instead.
func (*FileDesc) Descriptor() ([]byte, []int) {
	return file_partition_proto_rawDescGZIP(), []int{3}
}

func (x *FileDesc) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

var File_partition_proto protoreflect.FileDescriptor

var file_partition_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x6c, 0x0a, 0x04,
	0x46, 0x69, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x63, 0x68, 0x75, 0x6e, 0x6b, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x6e, 0x62,
	0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a,
	0x6e, 0x62, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x73, 0x22, 0x55, 0x0a, 0x0d, 0x46, 0x69,
	0x6c, 0x65, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x55, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x55, 0x75, 0x69, 0x64, 0x12, 0x28, 0x0a, 0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b,
	0x73, 0x22, 0x83, 0x01, 0x0a, 0x05, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x65, 0x6e, 0x64, 0x54, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x6e,
	0x64, 0x54, 0x6f, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x22, 0x26, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x44,
	0x65, 0x73, 0x63, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x32,
	0x81, 0x01, 0x0a, 0x09, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x34, 0x0a,
	0x05, 0x73, 0x70, 0x6c, 0x69, 0x74, 0x12, 0x0f, 0x2e, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x1a, 0x18, 0x2e, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x12, 0x13, 0x2e, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x44, 0x65, 0x73, 0x63, 0x1a, 0x18, 0x2e, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x00, 0x42, 0x03, 0x5a, 0x01, 0x2e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_partition_proto_rawDescOnce sync.Once
	file_partition_proto_rawDescData = file_partition_proto_rawDesc
)

func file_partition_proto_rawDescGZIP() []byte {
	file_partition_proto_rawDescOnce.Do(func() {
		file_partition_proto_rawDescData = protoimpl.X.CompressGZIP(file_partition_proto_rawDescData)
	})
	return file_partition_proto_rawDescData
}

var file_partition_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_partition_proto_goTypes = []any{
	(*File)(nil),          // 0: partition.File
	(*FilePartition)(nil), // 1: partition.FilePartition
	(*Chunk)(nil),         // 2: partition.Chunk
	(*FileDesc)(nil),      // 3: partition.FileDesc
}
var file_partition_proto_depIdxs = []int32{
	2, // 0: partition.FilePartition.chunks:type_name -> partition.Chunk
	0, // 1: partition.Partition.split:input_type -> partition.File
	3, // 2: partition.Partition.reconstruct:input_type -> partition.FileDesc
	1, // 3: partition.Partition.split:output_type -> partition.FilePartition
	1, // 4: partition.Partition.reconstruct:output_type -> partition.FilePartition
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_partition_proto_init() }
func file_partition_proto_init() {
	if File_partition_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_partition_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_partition_proto_goTypes,
		DependencyIndexes: file_partition_proto_depIdxs,
		MessageInfos:      file_partition_proto_msgTypes,
	}.Build()
	File_partition_proto = out.File
	file_partition_proto_rawDesc = nil
	file_partition_proto_goTypes = nil
	file_partition_proto_depIdxs = nil
}
