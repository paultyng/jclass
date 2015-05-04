package jclass

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf16"
)

type ConstantPoolInfo struct {
	Tag  uint8
	Info []byte
}

func (i *ConstantPoolInfo) String() string {
	switch i.Tag {
	case 7:
		return fmt.Sprintf("%s", (*ConstantClassInfo)(i))

	case 9:
		return fmt.Sprintf("%s", (*ConstantFieldrefInfo)(i))

	case 10:
		return fmt.Sprintf("%s", (*ConstantMethodrefInfo)(i))

	case 11:
		return fmt.Sprintf("%s", (*ConstantInterfaceMethodrefInfo)(i))

	case 8:
		return fmt.Sprintf("%s", (*ConstantStringInfo)(i))

	case 3:
		return fmt.Sprintf("%s", (*ConstantIntegerInfo)(i))

	case 4:
		return fmt.Sprintf("%s", (*ConstantFloatInfo)(i))

	case 5:
		return fmt.Sprintf("%s", (*ConstantLongInfo)(i))

	case 6:
		return fmt.Sprintf("%s", (*ConstantDoubleInfo)(i))

	case 12:
		return fmt.Sprintf("%s", (*ConstantNameAndTypeInfo)(i))

	case 1:
		return fmt.Sprintf("%s", (*ConstantUtf8Info)(i))

	case 15:
		return fmt.Sprintf("%s", (*ConstantMethodHandleInfo)(i))

	case 16:
		return fmt.Sprintf("%s", (*ConstantMethodTypeInfo)(i))

	case 18:
		return fmt.Sprintf("%s", (*ConstantInvokeDynamicInfo)(i))

	default:
		panic(fmt.Errorf("invalid tag: %d", i.Tag))
	}
}

func readEnoughBytes(info *ConstantPoolInfo, r io.Reader, buf []byte, n int) error {
	_, err := io.ReadFull(r, buf[:n])
	if err != nil {
		return err
	}
	info.Info = make([]byte, n)
	copy(info.Info, buf)
	return nil
}

func NewConstantPoolInfo(r io.Reader, buf []byte) (*ConstantPoolInfo, []byte, error) {
	rs := ConstantPoolInfo{}

	_, err := io.ReadFull(r, buf[:1])
	if err != nil {
		return nil, buf, err
	}

	rs.Tag = buf[0]
	switch rs.Tag {
	case 7, 8, 16:
		err = readEnoughBytes(&rs, r, buf, 2)
		if err != nil {
			return nil, buf, err
		}

	case 3, 4, 9, 10, 11, 12, 18:
		err = readEnoughBytes(&rs, r, buf, 4)
		if err != nil {
			return nil, buf, err
		}

	case 5, 6:
		err = readEnoughBytes(&rs, r, buf, 8)
		if err != nil {
			return nil, buf, err
		}

	case 1:
		_, err = io.ReadFull(r, buf[:2])
		if err != nil {
			return nil, buf, err
		}

		length := binary.BigEndian.Uint16(buf)
		bufSize := 2 + length
		if len(buf) < int(bufSize) {
			newBuf := make([]byte, bufSize)
			copy(newBuf[:2], buf)
			buf = newBuf
		}

		_, err = io.ReadFull(r, buf[2:bufSize])
		if err != nil {
			return nil, buf, err
		}

		rs.Info = make([]byte, bufSize)
		copy(rs.Info, buf)

	case 15:
		err = readEnoughBytes(&rs, r, buf, 3)
		if err != nil {
			return nil, buf, err
		}

	default:
		panic(fmt.Errorf("invalid tag: %d", rs.Tag))
	}

	return &rs, buf, nil
}

// CONSTANT_Class 7
type ConstantClassInfo ConstantPoolInfo

func (i *ConstantClassInfo) NameIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info)
}

func (i *ConstantClassInfo) String() string {
	return fmt.Sprintf("ConstantClassInfo [NameIndex:%d]", i.NameIndex())
}

// CONSTANT_Fieldref 9
type ConstantFieldrefInfo ConstantPoolInfo

func (i *ConstantFieldrefInfo) ClassIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info)
}

func (i *ConstantFieldrefInfo) NameAndTypeIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info[2:])
}

func (i *ConstantFieldrefInfo) String() string {
	return fmt.Sprintf("ConstantFieldrefInfo [ClassIndex:%d, NameAndTypeIndex:%d]",
		i.ClassIndex(), i.NameAndTypeIndex())
}

// CONSTANT_Methodref 10
type ConstantMethodrefInfo ConstantPoolInfo

func (i *ConstantMethodrefInfo) ClassIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info)
}

func (i *ConstantMethodrefInfo) NameAndTypeIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info[2:])
}

func (i *ConstantMethodrefInfo) String() string {
	return fmt.Sprintf("ConstantMethodrefInfo [ClassIndex: %d, NameAndTypeIndex: %d]",
		i.ClassIndex(), i.NameAndTypeIndex())
}

// CONSTANT_InterfaceMethodref 11
type ConstantInterfaceMethodrefInfo ConstantPoolInfo

func (i *ConstantInterfaceMethodrefInfo) ClassIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info)
}

func (i *ConstantInterfaceMethodrefInfo) NameAndTypeIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info[2:])
}

func (i *ConstantInterfaceMethodrefInfo) String() string {
	return fmt.Sprintf("ConstantInterfaceMethodrefInfo [ClassIndex: %d, NameAndTypeIndex: %d]",
		i.ClassIndex(), i.NameAndTypeIndex())
}

// CONSTANT_String 8
type ConstantStringInfo ConstantPoolInfo

func (i *ConstantStringInfo) StringIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info)
}

func (i *ConstantStringInfo) String() string {
	return fmt.Sprintf("ConstantStringInfo [StringIndex: %d]", i.StringIndex())
}

// CONSTANT_Integer 3
type ConstantIntegerInfo ConstantPoolInfo

func (i *ConstantIntegerInfo) Integer() int32 {
	return uint32Toint32(binary.BigEndian.Uint32(i.Info))
}

func (i *ConstantIntegerInfo) String() string {
	return fmt.Sprintf("ConstantIntegerInfo [value: %d]", i.Integer())
}

// CONSTANT_Float 4
type ConstantFloatInfo ConstantPoolInfo

func (i *ConstantFloatInfo) Float() float32 {
	return uint32ToFloat32(binary.BigEndian.Uint32(i.Info))
}

func (i *ConstantFloatInfo) String() string {
	return fmt.Sprintf("ConstantFloatInfo [value: %.7f]", i.Float())
}

// CONSTANT_Long 5
type ConstantLongInfo ConstantPoolInfo

func (i *ConstantLongInfo) Long() int64 {
	return uint64ToInt64(binary.BigEndian.Uint64(i.Info))
}

func (i *ConstantLongInfo) String() string {
	return fmt.Sprintf("ConstantLongInfo [value: %d]", i.Long())
}

// CONSTANT_Double 6
type ConstantDoubleInfo ConstantPoolInfo

func (i *ConstantDoubleInfo) Double() float64 {
	return uint64ToFloat64(binary.BigEndian.Uint64(i.Info))
}

func (i *ConstantDoubleInfo) String() string {
	return fmt.Sprintf("ConstantDoubleInfo [value: %.16f]", i.Double())
}

// CONSTANT_NameAndType 12
type ConstantNameAndTypeInfo ConstantPoolInfo

func (i *ConstantNameAndTypeInfo) NameIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info)
}

func (i *ConstantNameAndTypeInfo) DescriptorIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info[2:])
}

func (i *ConstantNameAndTypeInfo) String() string {
	return fmt.Sprintf("ConstantNameAndTypeInfo [NameIndex: %d, DescriptorIndex: %d]",
		i.NameIndex(), i.DescriptorIndex())
}

// CONSTANT_Utf8 1
type ConstantUtf8Info ConstantPoolInfo

func (i *ConstantUtf8Info) Length() uint16 {
	return binary.BigEndian.Uint16(i.Info)
}

func (i *ConstantUtf8Info) Bytes() []byte {
	return i.Info[2:]
}

func (info *ConstantUtf8Info) Utf8() string {
	ch := make([]uint16, 0, 512)
	var c, cc uint16
	buf := info.Info[2:]
	var st int

	for i, length := 0, len(buf); i < length; i++ {
		c = uint16(buf[i])
		switch st {
		case 0:
			if c < 0x80 {
				ch = append(ch, c)
			} else if c < 0xE0 && c > 0xBF {
				cc = c & 0x1F
				st = 1
			} else {
				cc = c & 0x0F
				st = 2
			}

		case 1:
			ch = append(ch, (cc<<6)|(c&0x3F))
			st = 0

		case 2:
			cc = (cc << 6) | (c & 0x3F)
			st = 1
		}
	}

	return string(utf16.Decode(ch))
}

func (i *ConstantUtf8Info) String() string {
	return fmt.Sprintf("ConstantUtf8Info [len: %d, utf8: %s]",
		i.Length(), i.Utf8())
}

// CONSTANT_MethodHandle 15
type ConstantMethodHandleInfo ConstantPoolInfo

func (i *ConstantMethodHandleInfo) ReferenceKind() uint8 {
	return i.Info[0]
}

func (i *ConstantMethodHandleInfo) ReferenceIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info[1:])
}

func (i *ConstantMethodHandleInfo) String() string {
	return fmt.Sprintf("ConstantMethodHandleInfo [ReferenceKind: %d, ReferenceIndex: %d]",
		i.ReferenceKind(), i.ReferenceIndex())
}

// CONSTANT_MethodType 16
type ConstantMethodTypeInfo ConstantPoolInfo

func (i *ConstantMethodTypeInfo) DescriptorIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info)
}

func (i *ConstantMethodTypeInfo) String() string {
	return fmt.Sprintf("ConstantMethodTypeInfo [DescriptorIndex: %d]", i.DescriptorIndex())
}

// CONSTANT_InvokeDynamic 18
type ConstantInvokeDynamicInfo ConstantPoolInfo

func (i *ConstantInvokeDynamicInfo) BootstrapMethodAttrIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info)
}

func (i *ConstantInvokeDynamicInfo) NameAndTypeIndex() uint16 {
	return binary.BigEndian.Uint16(i.Info[2:])
}

func (i *ConstantInvokeDynamicInfo) String() string {
	return fmt.Sprintf("ConstantInvokeDynamicInfo [BootstrapMethodAttrIndex: %d, NameAndTypeIndex: %d]",
		i.BootstrapMethodAttrIndex(), i.NameAndTypeIndex())
}
