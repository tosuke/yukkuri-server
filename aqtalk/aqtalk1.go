package aqtalk

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>

unsigned char* call_Synthe(void *p, const char *koe, int speed, int *size) {
	unsigned char* (*synthe)(const char*, int, int*) = p;
	return synthe(koe, speed, size);
}

void call_FreeWave(void *p, unsigned char *wav) {
	void (*freeWave)(unsigned char*) = p;
	return freeWave(wav);
}

*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Synthesizer struct {
	handle unsafe.Pointer
	pSynthe unsafe.Pointer
	pFreeWave unsafe.Pointer
}

func NewAqTalk1Synthesizer(typ string) (*Synthesizer, error) {
	switch typ {
	case "dvd", "f1", "f2", "imd1", "jgr", "m1", "m2", "r1": break
	default:
		return nil, ErrInvalidKoeType
	}

	dllName := fmt.Sprintf("./aqtk1-lnx/lib64/%s/libAquesTalk.so", typ)

	h := C.dlopen(C.CString(dllName), C.RTLD_LAZY)
	if h == nil {
		mes := C.GoString(C.dlerror())
		if(mes == "") {
			return nil, fmt.Errorf("failed to open %s: unknown error", dllName)
		}
		return nil, fmt.Errorf("failed to open: %s", mes)
	}

	pSynthe, err := getSymbol(h, "AquesTalk_Synthe_Utf8")
	if err != nil {
		return nil, err
	}
	pFreeWave, err := getSymbol(h, "AquesTalk_FreeWave")
	if err != nil {
		return nil, err
	}

	return &Synthesizer{
		handle: h,
		pSynthe: pSynthe,
		pFreeWave: pFreeWave,
	}, nil
}

func (s *Synthesizer) Close() {
	C.dlclose(s.handle)
}

func (s *Synthesizer) Synthe(koe string, speed uint32) ([]byte, error) {
	var size C.int
	r := C.call_Synthe(s.pSynthe, C.CString(koe), C.int(speed), &size)
	if r == nil {
		return nil, &AqTalk1Error{Code: int(size)}
	}
	defer C.call_FreeWave(s.pFreeWave, r)

	pWav := C.GoBytes(unsafe.Pointer(r), size)
	buf := make([]byte, size)
	copy(buf, pWav)

	return buf, nil
}

func getSymbol(mod unsafe.Pointer, name string) (unsafe.Pointer, error) {
	sym := C.dlsym(mod, C.CString(name))
	if sym == nil {
		return nil, fmt.Errorf("failed to get symbol %s", name)
	}
	return sym, nil
}
