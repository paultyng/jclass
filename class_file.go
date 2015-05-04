package jclass

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	MAGIC uint32 = 0xCAFEBABE
	DEBUG bool   = false
)

var (
	ERR_NOT_CLASS_FILE = errors.New("not class file")
)

type ClassFile struct {
	Magic uint32

	MinorVersion uint16
	MajorVersion uint16

	ConstantPoolCount uint16
	ConstantPool      []*ConstantPoolInfo

	AccessFlags ClassAccessFlags
	ThisClass   uint16
	SuperClass  uint16

	InterfaceCount uint16
	Interfaces     []uint16

	FieldsCount uint16
	Fields      []*FieldInfo

	MethodsCount uint16
	Methods      []*MethodInfo

	AttributesCount uint16
	Attributes      []*AttributeInfo
}

func (cf *ClassFile) String() string {
	s := &bytes.Buffer{}

	fmt.Fprintf(s, "// version: %d.%d\n", cf.MajorVersion, cf.MinorVersion)
	fmt.Fprintf(s, "// constant pool count: %d\n", cf.ConstantPoolCount)
	fmt.Fprintf(s, "// field count: %d\n", cf.FieldsCount)
	fmt.Fprintf(s, "// method count: %d\n", cf.MethodsCount)

	s.WriteString(cf.AccessFlagsString())
	s.WriteString(" ")
	s.WriteString(cf.ThisClassString())

	if cf.AccessFlags&CLASS_ACC_INTERFACE != 0 {
		if cf.HasInterfaces() {
			s.WriteString("\n\t\textends ")
			header := true
			for _, name := range cf.InterfaceStrings() {
				if header {
					header = false
				} else {
					s.WriteString(",\n\t\t        ")
				}
				s.WriteString(name)
			}
		}
	} else {
		superClass := cf.SuperClassString()
		if superClass != "" {
			s.WriteString("\n\t\textends ")
			s.WriteString(superClass)
		}

		if cf.HasInterfaces() {
			s.WriteString("\n\t\timplements ")
			header := true
			for _, name := range cf.InterfaceStrings() {
				if header {
					header = false
				} else {
					s.WriteString(",\n\t\t           ")
				}
				s.WriteString(name)
			}
		}
	}

	s.WriteString(" {")

	if cf.HasField() {
		s.WriteString("\n\t//---------- Fields ----------")
	}

	for _, field := range cf.Fields {
		s.WriteString("\n\t")
		s.WriteString(field.String())
	}

	if cf.HasMethod() {
		if cf.HasField() {
			s.WriteString("\n")
		}
		s.WriteString("\n\t//---------- Methods ----------")
	}

	for _, method := range cf.Methods {
		s.WriteString("\n\t")
		s.WriteString(method.String())
	}

	s.WriteString("\n}")

	return s.String()
}

// 是 interface 但不是 @interface
func (cf *ClassFile) IsInterface() bool {
	return cf.AccessFlags&CLASS_ACC_INTERFACE != 0 && cf.AccessFlags&CLASS_ACC_ANNOTATION == 0
}

// 是 @interface
func (cf *ClassFile) IsAnnotation() bool {
	return cf.AccessFlags&CLASS_ACC_ANNOTATION != 0
}

// 是 class 但不是 enum
func (cf *ClassFile) IsClass() bool {
	return cf.AccessFlags&CLASS_ACC_INTERFACE == 0 && cf.AccessFlags&CLASS_ACC_ENUM == 0
}

// 是 enum
func (cf *ClassFile) IsEnum() bool {
	return cf.AccessFlags&CLASS_ACC_ENUM != 0
}

// 是抽象类或者接口
func (cf *ClassFile) IsAbstract() bool {
	return cf.AccessFlags&CLASS_ACC_ABSTRACT != 0
}

// 如果为 true 则表示是 public 的类、interface、@interface 或 enum；否则为 package 的
func (cf *ClassFile) IsPublic() bool {
	return cf.AccessFlags&CLASS_ACC_PUBLIC != 0
}

// 如果为 true 则表明一定为类或 enum，不是接口、@interface 或者抽象类
func (cf *ClassFile) IsFinal() bool {
	return cf.AccessFlags&CLASS_ACC_FINAL != 0
}

// 是编译器生成的类或者接口
func (cf *ClassFile) IsSynthetic() bool {
	return cf.AccessFlags&CLASS_ACC_SYNTHETIC != 0
}

func (cf *ClassFile) AccessFlagsString() string {
	s := bytes.NewBuffer(nil)

	if cf.AccessFlags&CLASS_ACC_PUBLIC == 0 {
		s.WriteString("/* package */")
	} else {
		s.WriteString("public")
	}

	if cf.AccessFlags&CLASS_ACC_SYNTHETIC != 0 {
		s.WriteString(" /* synthetic */")
	}

	if cf.AccessFlags&CLASS_ACC_INTERFACE != 0 {
		if cf.AccessFlags&CLASS_ACC_ANNOTATION != 0 {
			s.WriteString(" @interface")
		} else {
			s.WriteString(" interface")
		}
	} else {
		if cf.AccessFlags&CLASS_ACC_FINAL != 0 {
			s.WriteString(" final")
		} else if cf.AccessFlags&CLASS_ACC_ABSTRACT != 0 {
			s.WriteString(" abstract")
		}

		if cf.AccessFlags&CLASS_ACC_ENUM != 0 {
			s.WriteString(" enum")
		} else {
			s.WriteString(" class")
		}
	}

	return s.String()
}

func (cf *ClassFile) getClassName(index uint16) string {
	classInfo := (*ConstantClassInfo)(cf.ConstantPool[index])
	nameIndex := classInfo.NameIndex()
	utf8 := ((*ConstantUtf8Info)(cf.ConstantPool[nameIndex])).Utf8()
	var debug string
	if DEBUG {
		debug = fmt.Sprintf("/* %d -> %d */", index, nameIndex)
	}
	return fmt.Sprintf("%s %s", utf8, debug)
}

func (cf *ClassFile) ThisClassString() string {
	return cf.getClassName(cf.ThisClass)
}

func (cf *ClassFile) SuperClassString() string {
	if cf.SuperClass != 0 {
		return cf.getClassName(cf.SuperClass)
	} else {
		return ""
	}
}

func (cf *ClassFile) HasInterfaces() bool {
	return cf.InterfaceCount > 0
}

func (cf *ClassFile) InterfaceStrings() []string {
	if cf.InterfaceCount == 0 {
		return nil
	}

	rs := make([]string, cf.InterfaceCount)
	for i, nameIndex := range cf.Interfaces {
		rs[i] = cf.getClassName(nameIndex)
	}

	return rs
}

func (cf *ClassFile) HasField() bool {
	return cf.FieldsCount > 0
}

func (cf *ClassFile) HasMethod() bool {
	return cf.MethodsCount > 0
}

func NewClassFileFromPath(path string) (*ClassFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return NewClassFile(file)
}

func NewClassFile(r io.Reader) (*ClassFile, error) {
	rs := ClassFile{}
	byteOrder := binary.BigEndian
	buf := make([]byte, 512)

	_, err := io.ReadFull(r, buf[:4])
	if err != nil {
		return nil, err
	}
	rs.Magic = byteOrder.Uint32(buf)
	if MAGIC != rs.Magic {
		return nil, ERR_NOT_CLASS_FILE
	}

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, err
	}
	rs.MinorVersion = byteOrder.Uint16(buf)

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, err
	}
	rs.MajorVersion = byteOrder.Uint16(buf)

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, err
	}
	rs.ConstantPoolCount = byteOrder.Uint16(buf)

	rs.ConstantPool = make([]*ConstantPoolInfo, rs.ConstantPoolCount)
	var info *ConstantPoolInfo
	for i := 1; i < int(rs.ConstantPoolCount); i++ {
		info, buf, err = NewConstantPoolInfo(r, buf)
		if err != nil {
			log.Fatalln(err)
		}

		rs.ConstantPool[i] = info
		if info.Tag == 5 || info.Tag == 6 {
			// All 8-byte constants take up two entries in the constant_pool table
			i++
		}
	}

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, err
	}
	rs.AccessFlags = ClassAccessFlags(byteOrder.Uint16(buf))

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, err
	}
	rs.ThisClass = byteOrder.Uint16(buf)

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, err
	}
	rs.SuperClass = byteOrder.Uint16(buf)

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, err
	}
	rs.InterfaceCount = byteOrder.Uint16(buf)

	rs.Interfaces = make([]uint16, rs.InterfaceCount)
	for i := 0; i < int(rs.InterfaceCount); i++ {
		_, err = io.ReadFull(r, buf[:2])
		if err != nil {
			return nil, err
		}
		rs.Interfaces[i] = byteOrder.Uint16(buf)
	}

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, err
	}
	rs.FieldsCount = byteOrder.Uint16(buf)

	rs.Fields = make([]*FieldInfo, rs.FieldsCount)
	var field *FieldInfo
	for i := 0; i < int(rs.FieldsCount); i++ {
		field, buf, err = NewFieldInfo(r, buf, rs.ConstantPool)
		if err != nil {
			return nil, err
		}
		rs.Fields[i] = field
	}

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, err
	}
	rs.MethodsCount = byteOrder.Uint16(buf)

	rs.Methods = make([]*MethodInfo, rs.MethodsCount)
	var method *MethodInfo
	for i := 0; i < int(rs.MethodsCount); i++ {
		method, buf, err = NewMethodInfo(r, buf, rs.ConstantPool)
		if err != nil {
			return nil, err
		}
		rs.Methods[i] = method
	}

	_, err = io.ReadFull(r, buf[:2])
	if err != nil {
		return nil, err
	}
	rs.AttributesCount = byteOrder.Uint16(buf)

	rs.Attributes = make([]*AttributeInfo, rs.AttributesCount)
	var attr *AttributeInfo
	for i := 0; i < int(rs.AttributesCount); i++ {
		attr, buf, err = NewAttributeInfo(r, buf, rs.ConstantPool)
		if err != nil {
			return nil, err
		}
		rs.Attributes[i] = attr
	}

	return &rs, nil
}
