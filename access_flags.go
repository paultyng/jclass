package jclass

type ClassAccessFlags uint16

const (
	CLASS_ACC_PUBLIC     ClassAccessFlags = 0x0001
	CLASS_ACC_FINAL      ClassAccessFlags = 0x0010
	CLASS_ACC_SUPER      ClassAccessFlags = 0x0020
	CLASS_ACC_INTERFACE  ClassAccessFlags = 0x0200
	CLASS_ACC_ABSTRACT   ClassAccessFlags = 0x0400
	CLASS_ACC_SYNTHETIC  ClassAccessFlags = 0x1000
	CLASS_ACC_ANNOTATION ClassAccessFlags = 0x2000
	CLASS_ACC_ENUM       ClassAccessFlags = 0x4000
)

type FieldAccessFlags uint16

const (
	FIELD_ACC_PUBLIC    FieldAccessFlags = 0x0001
	FIELD_ACC_PRIVATE   FieldAccessFlags = 0x0002
	FIELD_ACC_PROTECTED FieldAccessFlags = 0x0004
	FIELD_ACC_STATIC    FieldAccessFlags = 0x0008
	FIELD_ACC_FINAL     FieldAccessFlags = 0x0010
	FIELD_ACC_VOLATILE  FieldAccessFlags = 0x0040
	FIELD_ACC_TRANSIENT FieldAccessFlags = 0x0080
	FIELD_ACC_SYNTHETIC FieldAccessFlags = 0x1000
	FIELD_ACC_ENUM      FieldAccessFlags = 0x4000
)

type MethodAccessFlags uint16

const (
	METHOD_ACC_PUBLIC       MethodAccessFlags = 0x0001
	METHOD_ACC_PRIVATE      MethodAccessFlags = 0x0002
	METHOD_ACC_PROTECTED    MethodAccessFlags = 0x0004
	METHOD_ACC_STATIC       MethodAccessFlags = 0x0008
	METHOD_ACC_FINAL        MethodAccessFlags = 0x0010
	METHOD_ACC_SYNCHRONIZED MethodAccessFlags = 0x0020
	METHOD_ACC_BRIDGE       MethodAccessFlags = 0x0040
	METHOD_ACC_VARARGS      MethodAccessFlags = 0x0080
	METHOD_ACC_NATIVE       MethodAccessFlags = 0x0100
	METHOD_ACC_ABSTRACT     MethodAccessFlags = 0x0400
	METHOD_ACC_STRICT       MethodAccessFlags = 0x0800
	METHOD_ACC_SYNTHETIC    MethodAccessFlags = 0x1000
)
