// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: wire/messages.proto

/*
	Package wire is a generated protocol buffer package.

	It is generated from these files:
		wire/messages.proto

	It has these top-level messages:
		Handshake
		GetAddr
		Addr
		Address
		Ping
		Pong
		Reject
*/
package wire

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import bytes "bytes"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Handshake represents the first message between peers
type Handshake struct {
	SubVersion string `protobuf:"bytes,1,opt,name=subVersion,proto3" json:"subVersion,omitempty"`
}

func (m *Handshake) Reset()                    { *m = Handshake{} }
func (*Handshake) ProtoMessage()               {}
func (*Handshake) Descriptor() ([]byte, []int) { return fileDescriptorMessages, []int{0} }

func (m *Handshake) GetSubVersion() string {
	if m != nil {
		return m.SubVersion
	}
	return ""
}

// GetAddr is used to request for addresses from other peers
type GetAddr struct {
}

func (m *GetAddr) Reset()                    { *m = GetAddr{} }
func (*GetAddr) ProtoMessage()               {}
func (*GetAddr) Descriptor() ([]byte, []int) { return fileDescriptorMessages, []int{1} }

// GetAddrResponse is used to send addresses in response to a GetAddr
type Addr struct {
	Addresses []*Address `protobuf:"bytes,1,rep,name=addresses" json:"addresses,omitempty"`
}

func (m *Addr) Reset()                    { *m = Addr{} }
func (*Addr) ProtoMessage()               {}
func (*Addr) Descriptor() ([]byte, []int) { return fileDescriptorMessages, []int{2} }

func (m *Addr) GetAddresses() []*Address {
	if m != nil {
		return m.Addresses
	}
	return nil
}

type Address struct {
	Address   string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Timestamp int64  `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (m *Address) Reset()                    { *m = Address{} }
func (*Address) ProtoMessage()               {}
func (*Address) Descriptor() ([]byte, []int) { return fileDescriptorMessages, []int{3} }

func (m *Address) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Address) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

// Ping represents a ping message
type Ping struct {
}

func (m *Ping) Reset()                    { *m = Ping{} }
func (*Ping) ProtoMessage()               {}
func (*Ping) Descriptor() ([]byte, []int) { return fileDescriptorMessages, []int{4} }

// Pong represents a pong message
type Pong struct {
}

func (m *Pong) Reset()                    { *m = Pong{} }
func (*Pong) ProtoMessage()               {}
func (*Pong) Descriptor() ([]byte, []int) { return fileDescriptorMessages, []int{5} }

// Reject is used to inform a node that its message was rejected
type Reject struct {
	Message   string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Code      int32  `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	Reason    string `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	ExtraData []byte `protobuf:"bytes,4,opt,name=extraData,proto3" json:"extraData,omitempty"`
}

func (m *Reject) Reset()                    { *m = Reject{} }
func (*Reject) ProtoMessage()               {}
func (*Reject) Descriptor() ([]byte, []int) { return fileDescriptorMessages, []int{6} }

func (m *Reject) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Reject) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Reject) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *Reject) GetExtraData() []byte {
	if m != nil {
		return m.ExtraData
	}
	return nil
}

func init() {
	proto.RegisterType((*Handshake)(nil), "wire.Handshake")
	proto.RegisterType((*GetAddr)(nil), "wire.GetAddr")
	proto.RegisterType((*Addr)(nil), "wire.Addr")
	proto.RegisterType((*Address)(nil), "wire.Address")
	proto.RegisterType((*Ping)(nil), "wire.Ping")
	proto.RegisterType((*Pong)(nil), "wire.Pong")
	proto.RegisterType((*Reject)(nil), "wire.Reject")
}
func (this *Handshake) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Handshake)
	if !ok {
		that2, ok := that.(Handshake)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.SubVersion != that1.SubVersion {
		return false
	}
	return true
}
func (this *GetAddr) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GetAddr)
	if !ok {
		that2, ok := that.(GetAddr)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	return true
}
func (this *Addr) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Addr)
	if !ok {
		that2, ok := that.(Addr)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Addresses) != len(that1.Addresses) {
		return false
	}
	for i := range this.Addresses {
		if !this.Addresses[i].Equal(that1.Addresses[i]) {
			return false
		}
	}
	return true
}
func (this *Address) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Address)
	if !ok {
		that2, ok := that.(Address)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Address != that1.Address {
		return false
	}
	if this.Timestamp != that1.Timestamp {
		return false
	}
	return true
}
func (this *Ping) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Ping)
	if !ok {
		that2, ok := that.(Ping)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	return true
}
func (this *Pong) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Pong)
	if !ok {
		that2, ok := that.(Pong)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	return true
}
func (this *Reject) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Reject)
	if !ok {
		that2, ok := that.(Reject)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Message != that1.Message {
		return false
	}
	if this.Code != that1.Code {
		return false
	}
	if this.Reason != that1.Reason {
		return false
	}
	if !bytes.Equal(this.ExtraData, that1.ExtraData) {
		return false
	}
	return true
}
func (this *Handshake) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&wire.Handshake{")
	s = append(s, "SubVersion: "+fmt.Sprintf("%#v", this.SubVersion)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *GetAddr) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 4)
	s = append(s, "&wire.GetAddr{")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Addr) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&wire.Addr{")
	if this.Addresses != nil {
		s = append(s, "Addresses: "+fmt.Sprintf("%#v", this.Addresses)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Address) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&wire.Address{")
	s = append(s, "Address: "+fmt.Sprintf("%#v", this.Address)+",\n")
	s = append(s, "Timestamp: "+fmt.Sprintf("%#v", this.Timestamp)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Ping) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 4)
	s = append(s, "&wire.Ping{")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Pong) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 4)
	s = append(s, "&wire.Pong{")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Reject) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 8)
	s = append(s, "&wire.Reject{")
	s = append(s, "Message: "+fmt.Sprintf("%#v", this.Message)+",\n")
	s = append(s, "Code: "+fmt.Sprintf("%#v", this.Code)+",\n")
	s = append(s, "Reason: "+fmt.Sprintf("%#v", this.Reason)+",\n")
	s = append(s, "ExtraData: "+fmt.Sprintf("%#v", this.ExtraData)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringMessages(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Handshake) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Handshake) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.SubVersion) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintMessages(dAtA, i, uint64(len(m.SubVersion)))
		i += copy(dAtA[i:], m.SubVersion)
	}
	return i, nil
}

func (m *GetAddr) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetAddr) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *Addr) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Addr) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Addresses) > 0 {
		for _, msg := range m.Addresses {
			dAtA[i] = 0xa
			i++
			i = encodeVarintMessages(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Address) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Address) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintMessages(dAtA, i, uint64(len(m.Address)))
		i += copy(dAtA[i:], m.Address)
	}
	if m.Timestamp != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintMessages(dAtA, i, uint64(m.Timestamp))
	}
	return i, nil
}

func (m *Ping) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Ping) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *Pong) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Pong) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *Reject) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Reject) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintMessages(dAtA, i, uint64(len(m.Message)))
		i += copy(dAtA[i:], m.Message)
	}
	if m.Code != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintMessages(dAtA, i, uint64(m.Code))
	}
	if len(m.Reason) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintMessages(dAtA, i, uint64(len(m.Reason)))
		i += copy(dAtA[i:], m.Reason)
	}
	if len(m.ExtraData) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintMessages(dAtA, i, uint64(len(m.ExtraData)))
		i += copy(dAtA[i:], m.ExtraData)
	}
	return i, nil
}

func encodeVarintMessages(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Handshake) Size() (n int) {
	var l int
	_ = l
	l = len(m.SubVersion)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	return n
}

func (m *GetAddr) Size() (n int) {
	var l int
	_ = l
	return n
}

func (m *Addr) Size() (n int) {
	var l int
	_ = l
	if len(m.Addresses) > 0 {
		for _, e := range m.Addresses {
			l = e.Size()
			n += 1 + l + sovMessages(uint64(l))
		}
	}
	return n
}

func (m *Address) Size() (n int) {
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.Timestamp != 0 {
		n += 1 + sovMessages(uint64(m.Timestamp))
	}
	return n
}

func (m *Ping) Size() (n int) {
	var l int
	_ = l
	return n
}

func (m *Pong) Size() (n int) {
	var l int
	_ = l
	return n
}

func (m *Reject) Size() (n int) {
	var l int
	_ = l
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.Code != 0 {
		n += 1 + sovMessages(uint64(m.Code))
	}
	l = len(m.Reason)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	l = len(m.ExtraData)
	if l > 0 {
		n += 1 + l + sovMessages(uint64(l))
	}
	return n
}

func sovMessages(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozMessages(x uint64) (n int) {
	return sovMessages(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Handshake) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Handshake{`,
		`SubVersion:` + fmt.Sprintf("%v", this.SubVersion) + `,`,
		`}`,
	}, "")
	return s
}
func (this *GetAddr) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&GetAddr{`,
		`}`,
	}, "")
	return s
}
func (this *Addr) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Addr{`,
		`Addresses:` + strings.Replace(fmt.Sprintf("%v", this.Addresses), "Address", "Address", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Address) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Address{`,
		`Address:` + fmt.Sprintf("%v", this.Address) + `,`,
		`Timestamp:` + fmt.Sprintf("%v", this.Timestamp) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Ping) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Ping{`,
		`}`,
	}, "")
	return s
}
func (this *Pong) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Pong{`,
		`}`,
	}, "")
	return s
}
func (this *Reject) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Reject{`,
		`Message:` + fmt.Sprintf("%v", this.Message) + `,`,
		`Code:` + fmt.Sprintf("%v", this.Code) + `,`,
		`Reason:` + fmt.Sprintf("%v", this.Reason) + `,`,
		`ExtraData:` + fmt.Sprintf("%v", this.ExtraData) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringMessages(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Handshake) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Handshake: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Handshake: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubVersion", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SubVersion = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
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
func (m *GetAddr) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GetAddr: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetAddr: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
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
func (m *Addr) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Addr: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Addr: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addresses", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addresses = append(m.Addresses, &Address{})
			if err := m.Addresses[len(m.Addresses)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
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
func (m *Address) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Address: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Address: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
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
func (m *Ping) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Ping: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Ping: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
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
func (m *Pong) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Pong: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Pong: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
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
func (m *Reject) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Reject: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Reject: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Reason", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Reason = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtraData", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExtraData = append(m.ExtraData[:0], dAtA[iNdEx:postIndex]...)
			if m.ExtraData == nil {
				m.ExtraData = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
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
func skipMessages(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMessages
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
					return 0, ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessages
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthMessages
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowMessages
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipMessages(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthMessages = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMessages   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("wire/messages.proto", fileDescriptorMessages) }

var fileDescriptorMessages = []byte{
	// 293 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xb1, 0x4e, 0xf3, 0x30,
	0x10, 0xc7, 0x73, 0x5f, 0xf3, 0xa5, 0xca, 0x01, 0x8b, 0x91, 0x50, 0x06, 0x74, 0x8a, 0x32, 0x45,
	0x2a, 0x2a, 0x12, 0x7d, 0x82, 0x22, 0x24, 0x18, 0x51, 0x06, 0x76, 0xb7, 0x3e, 0x95, 0x00, 0x8d,
	0x2b, 0xdb, 0x08, 0x46, 0x1e, 0x81, 0xc7, 0xe0, 0x51, 0x18, 0x3b, 0x32, 0x52, 0xb3, 0x30, 0xf6,
	0x11, 0x50, 0x82, 0xab, 0x30, 0x9d, 0xff, 0x3f, 0x5b, 0xbf, 0xbf, 0x7c, 0x78, 0xf8, 0x54, 0x1b,
	0x3e, 0x5d, 0xb2, 0xb5, 0x72, 0xc1, 0x76, 0xbc, 0x32, 0xda, 0x69, 0x11, 0xb7, 0xb0, 0x18, 0x61,
	0x7a, 0x25, 0x1b, 0x65, 0x6f, 0xe5, 0x3d, 0x0b, 0x42, 0xb4, 0x8f, 0xb3, 0x1b, 0x36, 0xb6, 0xd6,
	0x4d, 0x06, 0x39, 0x94, 0x69, 0xf5, 0x87, 0x14, 0x29, 0x0e, 0x2f, 0xd9, 0x4d, 0x95, 0x32, 0xc5,
	0x04, 0xe3, 0x76, 0x8a, 0x11, 0xa6, 0x52, 0x29, 0xc3, 0xd6, 0xb2, 0xcd, 0x20, 0x1f, 0x94, 0x7b,
	0x67, 0x07, 0xe3, 0xd6, 0x3c, 0x9e, 0xfe, 0xe2, 0xaa, 0xbf, 0x2f, 0xa6, 0x38, 0x0c, 0x54, 0x64,
	0x38, 0x0c, 0x3c, 0xf4, 0xec, 0xa2, 0x38, 0xc6, 0xd4, 0xd5, 0x4b, 0xb6, 0x4e, 0x2e, 0x57, 0xd9,
	0xbf, 0x1c, 0xca, 0x41, 0xd5, 0x83, 0x22, 0xc1, 0xf8, 0xba, 0x6e, 0x16, 0xdd, 0xd4, 0xcd, 0xa2,
	0x78, 0xc0, 0xa4, 0xe2, 0x3b, 0x9e, 0xbb, 0xd6, 0x18, 0x7e, 0xb8, 0x33, 0x86, 0x28, 0x04, 0xc6,
	0x73, 0xad, 0xb8, 0x93, 0xfd, 0xaf, 0xba, 0xb3, 0x38, 0xc2, 0xc4, 0xb0, 0xb4, 0xba, 0xc9, 0x06,
	0xdd, 0xe3, 0x90, 0xda, 0x76, 0x7e, 0x76, 0x46, 0x5e, 0x48, 0x27, 0xb3, 0x38, 0x87, 0x72, 0xbf,
	0xea, 0xc1, 0xf9, 0xc9, 0x7a, 0x43, 0xd1, 0xc7, 0x86, 0xa2, 0xed, 0x86, 0xe0, 0xc5, 0x13, 0xbc,
	0x79, 0x82, 0x77, 0x4f, 0xb0, 0xf6, 0x04, 0x9f, 0x9e, 0xe0, 0xdb, 0x53, 0xb4, 0xf5, 0x04, 0xaf,
	0x5f, 0x14, 0xcd, 0x92, 0x6e, 0xd1, 0x93, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x77, 0xe9, 0x89,
	0xc4, 0x7f, 0x01, 0x00, 0x00,
}
