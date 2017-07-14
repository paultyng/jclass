package jclass

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type AttributeInfo struct {
	NameIndex uint16
	Length    uint32
	Info      []byte

	Annotations []*Annotation

	cp []*ConstantPoolInfo
}

func (i *AttributeInfo) String() string {
	return fmt.Sprintf("AttributeInfo: {%s}", i.NameString())
}

func (i *AttributeInfo) NameString() string {
	return ((*ConstantUtf8Info)(i.cp[i.NameIndex])).Utf8()
}

func (a *AttributeInfo) ConstantPoolInfo(i uint16) *ConstantPoolInfo {
	return a.cp[int(i)]
}

func NewAttributeInfo(r io.Reader, buf []byte, cp []*ConstantPoolInfo) (*AttributeInfo, []byte, error) {
	rs := AttributeInfo{}
	byteOrder := binary.BigEndian

	_, err := io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, buf, err
	}
	rs.NameIndex = byteOrder.Uint16(buf)

	_, err = io.ReadFull(r, buf[:4])
	if err != nil {
		return nil, buf, err
	}
	rs.Length = byteOrder.Uint32(buf)

	size := int(rs.Length)
	if cap(buf) < size {
		buf = make([]byte, size)
	}
	_, err = io.ReadFull(r, buf[:size])
	if err != nil {
		return nil, buf, err
	}
	rs.Info = make([]byte, size)
	copy(rs.Info, buf)

	rs.cp = cp

	switch rs.NameString() {
	case "RuntimeVisibleAnnotations":
		r := bytes.NewReader(rs.Info)
		_, err = io.ReadFull(r, buf[:2])
		if err != nil {
			return nil, buf, err
		}
		num := int(byteOrder.Uint16(buf))
		rs.Annotations = make([]*Annotation, num)
		for i := 0; i < num; i++ {
			ann, buf, err := NewAnnotation(r, buf, cp)
			if err != nil {
				return nil, buf, err
			}
			rs.Annotations = append(rs.Annotations, ann)
		}
	}

	return &rs, buf, nil
}
