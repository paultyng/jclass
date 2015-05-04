package jclass

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type FieldInfo struct {
	AccessFlags     FieldAccessFlags
	NameIndex       uint16
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []*AttributeInfo

	cp []*ConstantPoolInfo
}

func (i *FieldInfo) String() string {
	return fmt.Sprintf("%s %s %s",
		i.AccessFlagsString(),
		i.DescriptorString(),
		i.NameString())
}

func (i *FieldInfo) DescriptorString() string {
	return ((*ConstantUtf8Info)(i.cp[int(i.DescriptorIndex)])).Utf8()
}

func (i *FieldInfo) NameString() string {
	return ((*ConstantUtf8Info)(i.cp[int(i.NameIndex)])).Utf8()
}

func (i *FieldInfo) AccessFlagsString() string {
	s := bytes.NewBuffer(nil)

	switch {
	case i.AccessFlags&FIELD_ACC_PUBLIC != 0:
		s.WriteString("public")

	case i.AccessFlags&FIELD_ACC_PRIVATE != 0:
		s.WriteString("private")

	case i.AccessFlags&FIELD_ACC_PROTECTED != 0:
		s.WriteString("protected")

	default:
		s.WriteString("/* package */")
	}

	if i.AccessFlags&FIELD_ACC_STATIC != 0 {
		s.WriteString(" static")
	}

	switch {
	case i.AccessFlags&FIELD_ACC_FINAL != 0:
		s.WriteString(" final")

	case i.AccessFlags&FIELD_ACC_VOLATILE != 0:
		s.WriteString(" volatile")
	}

	if i.AccessFlags&FIELD_ACC_TRANSIENT != 0 {
		s.WriteString(" transient")
	}

	if i.AccessFlags&FIELD_ACC_SYNTHETIC != 0 {
		s.WriteString(" /* synthetic */")
	}

	if i.AccessFlags&FIELD_ACC_ENUM != 0 {
		s.WriteString(" /* enum */")
	}

	return s.String()
}

func NewFieldInfo(r io.Reader, buf []byte, cp []*ConstantPoolInfo) (*FieldInfo, []byte, error) {
	rs := FieldInfo{}
	byteOrder := binary.BigEndian

	_, err := io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, buf, err
	}
	rs.AccessFlags = FieldAccessFlags(byteOrder.Uint16(buf))

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, buf, err
	}
	rs.NameIndex = byteOrder.Uint16(buf)

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, buf, err
	}
	rs.DescriptorIndex = byteOrder.Uint16(buf)

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, buf, err
	}
	rs.AttributesCount = byteOrder.Uint16(buf)

	size := int(rs.AttributesCount)
	if size > 0 {
		rs.Attributes = make([]*AttributeInfo, size)
		var attr *AttributeInfo
		for i := 0; i < size; i++ {
			attr, buf, err = NewAttributeInfo(r, buf, cp)
			if err != nil {
				return nil, buf, err
			}
			rs.Attributes[i] = attr
		}
	}

	rs.cp = cp
	return &rs, buf, nil
}
