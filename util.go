package jclass

import (
	"io"
	"unsafe"
)

func uint32Toint32(in uint32) int32 {
	return *((*int32)(unsafe.Pointer(&in)))
}

func uint32ToFloat32(in uint32) float32 {
	return *((*float32)(unsafe.Pointer(&in)))
}

func uint64ToInt64(in uint64) int64 {
	return *((*int64)(unsafe.Pointer(&in)))
}

func uint64ToFloat64(in uint64) float64 {
	return *((*float64)(unsafe.Pointer(&in)))
}

func readEnoughBytes(r io.Reader, buf []byte, n int) ([]byte, error) {
	_, err := io.ReadFull(r, buf[:n])
	if err != nil {
		return nil, err
	}
	info := make([]byte, n)
	copy(info, buf)
	return info, nil
}
