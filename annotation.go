package jclass

import (
	"encoding/binary"
	"io"
)

type Annotation struct {
	TypeIndex            uint16
	NumElementValuePairs uint16
	ElementValuePairs    []*ElementValuePair

	cp []*ConstantPoolInfo
}

func (a *Annotation) TypeString() string {
	return ((*ConstantUtf8Info)(a.cp[a.TypeIndex])).Utf8()
}

func (a *Annotation) ConstantPoolInfo(i uint16) *ConstantPoolInfo {
	return a.cp[int(i)]
}

func NewAnnotation(r io.Reader, buf []byte, cp []*ConstantPoolInfo) (*Annotation, []byte, error) {
	rs := Annotation{}
	byteOrder := binary.BigEndian

	_, err := io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, buf, err
	}
	rs.TypeIndex = byteOrder.Uint16(buf)

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, buf, err
	}
	rs.NumElementValuePairs = byteOrder.Uint16(buf)

	num := int(rs.NumElementValuePairs)
	rs.ElementValuePairs = make([]*ElementValuePair, num)
	for i := 0; i < num; i++ {
		evp, buf, err := NewElementValuePair(r, buf, cp)
		if err != nil {
			return nil, buf, err
		}
		rs.ElementValuePairs = append(rs.ElementValuePairs, evp)
	}

	rs.cp = cp
	return &rs, buf, err
}
