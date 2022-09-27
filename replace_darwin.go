package monkey

import (
	"reflect"
	"syscall"
	"unsafe"
)

func copyToLocation(target uintptr, bytes []byte) {
	modifyBinary(target, bytes)
}

func modifyBinary(target uintptr, bytes []byte) {
	function := entryAddress(target, len(bytes))
	mprotectCrossPage(target, len(bytes), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC)
	copy(function, bytes)
	mprotectCrossPage(target, len(bytes), syscall.PROT_READ|syscall.PROT_EXEC)
}

// func mprotectCrossPage(addr uintptr, length int, prot int) error {
func mprotectCrossPage(addr uintptr, length int, prot int) {
	pageSize := syscall.Getpagesize()
	for p := pageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := entryAddress(p, pageSize)
		if err := syscall.Mprotect(page, prot); err != nil {
			panic(err)
			// return err
		}
	}
	// return nil
}

func entryAddress(p uintptr, l int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: p, Len: l, Cap: l}))
}
