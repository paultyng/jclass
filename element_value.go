package jclass

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ElementValue struct {
	Tag  string
	Info []byte

	cpi uint16
	cp  []*ConstantPoolInfo
}

func (ev *ElementValue) ConstantPoolInfo(i uint16) *ConstantPoolInfo {
	return ev.cp[int(i)]
}

func NewElementValue(r io.Reader, buf []byte, cp []*ConstantPoolInfo) (*ElementValue, []byte, error) {
	rs := ElementValue{cp: cp}
	byteOrder := binary.BigEndian

	_, err := io.ReadFull(r, buf[:1])
	if err != nil {
		return nil, buf, err
	}

	rs.Tag = string(buf[0])
	switch rs.Tag {
	case "B":
		fallthrough
	case "C":
		fallthrough
	case "D":
		fallthrough
	case "F":
		fallthrough
	case "I":
		fallthrough
	case "J":
		fallthrough
	case "S":
		fallthrough
	case "Z":
		fallthrough
	case "s":
		_, err := io.ReadFull(r, buf[:2])
		if err != nil {
			return nil, buf, err
		}
		rs.cpi = byteOrder.Uint16(buf)
	default:
		panic(fmt.Errorf("invalid element value tag: %s", rs.Tag))
	}

	return &rs, buf, nil
}

type ElementValuePair struct {
	ElementNameIndex uint16
	Value            *ElementValue

	cp []*ConstantPoolInfo
}

func (evp *ElementValuePair) ElementNameString() string {
	return ((*ConstantUtf8Info)(evp.cp[evp.ElementNameIndex])).Utf8()
}

func (evp *ElementValuePair) ConstantPoolInfo(i uint16) *ConstantPoolInfo {
	return evp.cp[int(i)]
}

func NewElementValuePair(r io.Reader, buf []byte, cp []*ConstantPoolInfo) (*ElementValuePair, []byte, error) {
	rs := ElementValuePair{}
	byteOrder := binary.BigEndian

	_, err := io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, buf, err
	}
	rs.ElementNameIndex = byteOrder.Uint16(buf)

	v, buf, err := NewElementValue(r, buf, cp)
	if err != nil {
		return nil, buf, err
	}

	rs.Value = v
	return &rs, buf, nil
}
