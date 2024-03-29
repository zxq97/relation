// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: app/relationship/pkg/bizdata/relationship_biz.proto

package bizdata

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type FollowItem struct {
	ToUid                int64    `protobuf:"varint,1,opt,name=to_uid,json=toUid,proto3" json:"to_uid,omitempty"`
	CreateTime           int64    `protobuf:"varint,2,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FollowItem) Reset()         { *m = FollowItem{} }
func (m *FollowItem) String() string { return proto.CompactTextString(m) }
func (*FollowItem) ProtoMessage()    {}
func (*FollowItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec8f742d0bf97fa6, []int{0}
}
func (m *FollowItem) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FollowItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FollowItem.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FollowItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FollowItem.Merge(m, src)
}
func (m *FollowItem) XXX_Size() int {
	return m.Size()
}
func (m *FollowItem) XXX_DiscardUnknown() {
	xxx_messageInfo_FollowItem.DiscardUnknown(m)
}

var xxx_messageInfo_FollowItem proto.InternalMessageInfo

func (m *FollowItem) GetToUid() int64 {
	if m != nil {
		return m.ToUid
	}
	return 0
}

func (m *FollowItem) GetCreateTime() int64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

type FollowList struct {
	List                 []*FollowItem `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *FollowList) Reset()         { *m = FollowList{} }
func (m *FollowList) String() string { return proto.CompactTextString(m) }
func (*FollowList) ProtoMessage()    {}
func (*FollowList) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec8f742d0bf97fa6, []int{1}
}
func (m *FollowList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FollowList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FollowList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FollowList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FollowList.Merge(m, src)
}
func (m *FollowList) XXX_Size() int {
	return m.Size()
}
func (m *FollowList) XXX_DiscardUnknown() {
	xxx_messageInfo_FollowList.DiscardUnknown(m)
}

var xxx_messageInfo_FollowList proto.InternalMessageInfo

func (m *FollowList) GetList() []*FollowItem {
	if m != nil {
		return m.List
	}
	return nil
}

type RelationCount struct {
	Uid                  int64    `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	FollowCount          int32    `protobuf:"varint,2,opt,name=follow_count,json=followCount,proto3" json:"follow_count,omitempty"`
	FollowerCount        int32    `protobuf:"varint,3,opt,name=follower_count,json=followerCount,proto3" json:"follower_count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RelationCount) Reset()         { *m = RelationCount{} }
func (m *RelationCount) String() string { return proto.CompactTextString(m) }
func (*RelationCount) ProtoMessage()    {}
func (*RelationCount) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec8f742d0bf97fa6, []int{2}
}
func (m *RelationCount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RelationCount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RelationCount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RelationCount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelationCount.Merge(m, src)
}
func (m *RelationCount) XXX_Size() int {
	return m.Size()
}
func (m *RelationCount) XXX_DiscardUnknown() {
	xxx_messageInfo_RelationCount.DiscardUnknown(m)
}

var xxx_messageInfo_RelationCount proto.InternalMessageInfo

func (m *RelationCount) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *RelationCount) GetFollowCount() int32 {
	if m != nil {
		return m.FollowCount
	}
	return 0
}

func (m *RelationCount) GetFollowerCount() int32 {
	if m != nil {
		return m.FollowerCount
	}
	return 0
}

type UserRelation struct {
	Relation             int32    `protobuf:"varint,1,opt,name=relation,proto3" json:"relation,omitempty"`
	FollowTime           int64    `protobuf:"varint,2,opt,name=follow_time,json=followTime,proto3" json:"follow_time,omitempty"`
	FollowedTime         int64    `protobuf:"varint,3,opt,name=followed_time,json=followedTime,proto3" json:"followed_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserRelation) Reset()         { *m = UserRelation{} }
func (m *UserRelation) String() string { return proto.CompactTextString(m) }
func (*UserRelation) ProtoMessage()    {}
func (*UserRelation) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec8f742d0bf97fa6, []int{3}
}
func (m *UserRelation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UserRelation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UserRelation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UserRelation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserRelation.Merge(m, src)
}
func (m *UserRelation) XXX_Size() int {
	return m.Size()
}
func (m *UserRelation) XXX_DiscardUnknown() {
	xxx_messageInfo_UserRelation.DiscardUnknown(m)
}

var xxx_messageInfo_UserRelation proto.InternalMessageInfo

func (m *UserRelation) GetRelation() int32 {
	if m != nil {
		return m.Relation
	}
	return 0
}

func (m *UserRelation) GetFollowTime() int64 {
	if m != nil {
		return m.FollowTime
	}
	return 0
}

func (m *UserRelation) GetFollowedTime() int64 {
	if m != nil {
		return m.FollowedTime
	}
	return 0
}

func init() {
	proto.RegisterType((*FollowItem)(nil), "relationship.pkg.bizdata.FollowItem")
	proto.RegisterType((*FollowList)(nil), "relationship.pkg.bizdata.FollowList")
	proto.RegisterType((*RelationCount)(nil), "relationship.pkg.bizdata.RelationCount")
	proto.RegisterType((*UserRelation)(nil), "relationship.pkg.bizdata.UserRelation")
}

func init() {
	proto.RegisterFile("app/relationship/pkg/bizdata/relationship_biz.proto", fileDescriptor_ec8f742d0bf97fa6)
}

var fileDescriptor_ec8f742d0bf97fa6 = []byte{
	// 309 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x5d, 0x63, 0x8b, 0x4c, 0x52, 0x91, 0x05, 0x21, 0x78, 0x88, 0x35, 0x56, 0xe8, 0x29,
	0x01, 0x0b, 0x22, 0x78, 0x53, 0x29, 0x08, 0x9e, 0x82, 0xbd, 0x78, 0x09, 0x49, 0xb3, 0xd6, 0x25,
	0x49, 0x77, 0x49, 0xa6, 0x08, 0x7d, 0x12, 0x1f, 0xc9, 0xa3, 0x8f, 0x20, 0xf1, 0x45, 0x24, 0xbb,
	0x59, 0x5b, 0x0f, 0x7a, 0x4a, 0xf6, 0x9f, 0xef, 0x9f, 0x7f, 0x98, 0x81, 0x49, 0x22, 0x65, 0x58,
	0xb1, 0x22, 0x41, 0x2e, 0x96, 0xf5, 0x0b, 0x97, 0xa1, 0xcc, 0x17, 0x61, 0xca, 0xd7, 0x59, 0x82,
	0xc9, 0xaf, 0x42, 0x9c, 0xf2, 0x75, 0x20, 0x2b, 0x81, 0x82, 0xba, 0xdb, 0x7a, 0x20, 0xf3, 0x45,
	0xd0, 0x19, 0xfc, 0x3b, 0x80, 0xa9, 0x28, 0x0a, 0xf1, 0x7a, 0x8f, 0xac, 0xa4, 0x47, 0xd0, 0x47,
	0x11, 0xaf, 0x78, 0xe6, 0x92, 0x21, 0x19, 0x5b, 0x51, 0x0f, 0xc5, 0x8c, 0x67, 0xf4, 0x04, 0xec,
	0x79, 0xc5, 0x12, 0x64, 0x31, 0xf2, 0x92, 0xb9, 0xbb, 0xaa, 0x06, 0x5a, 0x7a, 0xe4, 0x25, 0xf3,
	0xa7, 0xa6, 0xcb, 0x03, 0xaf, 0x91, 0x5e, 0xc1, 0x5e, 0xc1, 0x6b, 0x74, 0xc9, 0xd0, 0x1a, 0xdb,
	0x17, 0xa3, 0xe0, 0xaf, 0xf0, 0x60, 0x93, 0x1c, 0x29, 0x87, 0x9f, 0xc3, 0x20, 0xea, 0xe0, 0x5b,
	0xb1, 0x5a, 0x22, 0x3d, 0x04, 0x6b, 0x33, 0x4d, 0xfb, 0x4b, 0x4f, 0xc1, 0x79, 0x56, 0xb6, 0x78,
	0xde, 0x12, 0x6a, 0x98, 0x5e, 0x64, 0x6b, 0x4d, 0x9b, 0xce, 0xe1, 0x40, 0x3f, 0x59, 0xd5, 0x41,
	0x96, 0x82, 0x06, 0x46, 0x55, 0x98, 0x2f, 0xc1, 0x99, 0xd5, 0xac, 0x32, 0x81, 0xf4, 0x18, 0xf6,
	0xcd, 0xa4, 0x2a, 0xb0, 0x17, 0xfd, 0xbc, 0xdb, 0x0d, 0x74, 0xa9, 0xdb, 0x1b, 0xd0, 0x52, 0xbb,
	0x01, 0x7a, 0x06, 0xa6, 0x7b, 0xa6, 0x11, 0x4b, 0x21, 0x8e, 0x11, 0x5b, 0xe8, 0xe6, 0xf2, 0xbd,
	0xf1, 0xc8, 0x47, 0xe3, 0x91, 0xcf, 0xc6, 0x23, 0x6f, 0x5f, 0xde, 0xce, 0xd3, 0xe8, 0xbf, 0x6b,
	0x5e, 0x77, 0xdf, 0xb4, 0xaf, 0xae, 0x38, 0xf9, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x06, 0x73, 0xba,
	0xd9, 0xfc, 0x01, 0x00, 0x00,
}

func (m *FollowItem) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FollowItem) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FollowItem) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.CreateTime != 0 {
		i = encodeVarintRelationshipBiz(dAtA, i, uint64(m.CreateTime))
		i--
		dAtA[i] = 0x10
	}
	if m.ToUid != 0 {
		i = encodeVarintRelationshipBiz(dAtA, i, uint64(m.ToUid))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *FollowList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FollowList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FollowList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.List) > 0 {
		for iNdEx := len(m.List) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.List[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintRelationshipBiz(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *RelationCount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RelationCount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RelationCount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.FollowerCount != 0 {
		i = encodeVarintRelationshipBiz(dAtA, i, uint64(m.FollowerCount))
		i--
		dAtA[i] = 0x18
	}
	if m.FollowCount != 0 {
		i = encodeVarintRelationshipBiz(dAtA, i, uint64(m.FollowCount))
		i--
		dAtA[i] = 0x10
	}
	if m.Uid != 0 {
		i = encodeVarintRelationshipBiz(dAtA, i, uint64(m.Uid))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *UserRelation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UserRelation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UserRelation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.FollowedTime != 0 {
		i = encodeVarintRelationshipBiz(dAtA, i, uint64(m.FollowedTime))
		i--
		dAtA[i] = 0x18
	}
	if m.FollowTime != 0 {
		i = encodeVarintRelationshipBiz(dAtA, i, uint64(m.FollowTime))
		i--
		dAtA[i] = 0x10
	}
	if m.Relation != 0 {
		i = encodeVarintRelationshipBiz(dAtA, i, uint64(m.Relation))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintRelationshipBiz(dAtA []byte, offset int, v uint64) int {
	offset -= sovRelationshipBiz(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *FollowItem) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ToUid != 0 {
		n += 1 + sovRelationshipBiz(uint64(m.ToUid))
	}
	if m.CreateTime != 0 {
		n += 1 + sovRelationshipBiz(uint64(m.CreateTime))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *FollowList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.List) > 0 {
		for _, e := range m.List {
			l = e.Size()
			n += 1 + l + sovRelationshipBiz(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *RelationCount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Uid != 0 {
		n += 1 + sovRelationshipBiz(uint64(m.Uid))
	}
	if m.FollowCount != 0 {
		n += 1 + sovRelationshipBiz(uint64(m.FollowCount))
	}
	if m.FollowerCount != 0 {
		n += 1 + sovRelationshipBiz(uint64(m.FollowerCount))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *UserRelation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Relation != 0 {
		n += 1 + sovRelationshipBiz(uint64(m.Relation))
	}
	if m.FollowTime != 0 {
		n += 1 + sovRelationshipBiz(uint64(m.FollowTime))
	}
	if m.FollowedTime != 0 {
		n += 1 + sovRelationshipBiz(uint64(m.FollowedTime))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovRelationshipBiz(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRelationshipBiz(x uint64) (n int) {
	return sovRelationshipBiz(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *FollowItem) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRelationshipBiz
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
			return fmt.Errorf("proto: FollowItem: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FollowItem: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ToUid", wireType)
			}
			m.ToUid = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelationshipBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ToUid |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreateTime", wireType)
			}
			m.CreateTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelationshipBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreateTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRelationshipBiz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRelationshipBiz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *FollowList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRelationshipBiz
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
			return fmt.Errorf("proto: FollowList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FollowList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field List", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelationshipBiz
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
				return ErrInvalidLengthRelationshipBiz
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRelationshipBiz
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.List = append(m.List, &FollowItem{})
			if err := m.List[len(m.List)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRelationshipBiz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRelationshipBiz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RelationCount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRelationshipBiz
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
			return fmt.Errorf("proto: RelationCount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RelationCount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uid", wireType)
			}
			m.Uid = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelationshipBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Uid |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FollowCount", wireType)
			}
			m.FollowCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelationshipBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FollowCount |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FollowerCount", wireType)
			}
			m.FollowerCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelationshipBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FollowerCount |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRelationshipBiz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRelationshipBiz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UserRelation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRelationshipBiz
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
			return fmt.Errorf("proto: UserRelation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UserRelation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Relation", wireType)
			}
			m.Relation = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelationshipBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Relation |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FollowTime", wireType)
			}
			m.FollowTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelationshipBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FollowTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FollowedTime", wireType)
			}
			m.FollowedTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelationshipBiz
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FollowedTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRelationshipBiz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRelationshipBiz
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipRelationshipBiz(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRelationshipBiz
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
					return 0, ErrIntOverflowRelationshipBiz
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
					return 0, ErrIntOverflowRelationshipBiz
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
				return 0, ErrInvalidLengthRelationshipBiz
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRelationshipBiz
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRelationshipBiz
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRelationshipBiz        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRelationshipBiz          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRelationshipBiz = fmt.Errorf("proto: unexpected end of group")
)
