package jclass

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type MethodInfo struct {
	AccessFlags     MethodAccessFlags
	NameIndex       uint16
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []*AttributeInfo

	cp []*ConstantPoolInfo
}

func (i *MethodInfo) String() string {
	return fmt.Sprintf("%s %s %s",
		i.AccessFlagsString(),
		i.NameString(),
		i.DescriptorString())
}

func (i *MethodInfo) DescriptorString() string {
	return ((*ConstantUtf8Info)(i.cp[int(i.DescriptorIndex)])).Utf8()
}

func (i *MethodInfo) NameString() string {
	return ((*ConstantUtf8Info)(i.cp[int(i.NameIndex)])).Utf8()
}

func (i *MethodInfo) AccessFlagsString() string {
	s := bytes.NewBuffer(nil)

	switch {
	case i.AccessFlags&METHOD_ACC_PUBLIC != 0:
		s.WriteString("public")

	case i.AccessFlags&METHOD_ACC_PRIVATE != 0:
		s.WriteString("private")

	case i.AccessFlags&METHOD_ACC_PROTECTED != 0:
		s.WriteString("protected")

	default:
		s.WriteString("/* package */")
	}

	if i.AccessFlags&METHOD_ACC_ABSTRACT != 0 {
		s.WriteString(" abstract")
	} else {
		if i.AccessFlags&METHOD_ACC_STATIC != 0 {
			s.WriteString(" static")
		}

		if i.AccessFlags&METHOD_ACC_FINAL != 0 {
			s.WriteString(" final")
		}

		if i.AccessFlags&METHOD_ACC_SYNCHRONIZED != 0 {
			s.WriteString(" synchronized")
		}

		if i.AccessFlags&METHOD_ACC_NATIVE != 0 {
			s.WriteString(" native")
		}

		if i.AccessFlags&METHOD_ACC_STRICT != 0 {
			s.WriteString(" strict")
		}
	}

	if i.AccessFlags&METHOD_ACC_BRIDGE != 0 {
		s.WriteString(" /* bridge */")
	}

	if i.AccessFlags&METHOD_ACC_VARARGS != 0 {
		s.WriteString(" /* varargs */")
	}

	if i.AccessFlags&METHOD_ACC_SYNTHETIC != 0 {
		s.WriteString(" /* synthetic */")
	}

	return s.String()
}

func NewMethodInfo(r io.Reader, buf []byte, cp []*ConstantPoolInfo) (*MethodInfo, []byte, error) {
	rs := MethodInfo{}
	byteOrder := binary.BigEndian

	_, err := io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, buf, err
	}
	rs.AccessFlags = MethodAccessFlags(byteOrder.Uint16(buf))

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
