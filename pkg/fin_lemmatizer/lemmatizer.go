package fin_lemmatizer

/*
#cgo CFLAGS: -I/usr/include/libvoikko
#cgo LDFLAGS: -lvoikko
#include <voikko.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type VoikkoHandle *C.struct_VoikkoHandle

func InitVoikko() (VoikkoHandle, error) {
	var handle VoikkoHandle
	var errorStr *C.char
	cLang := C.CString("fi") // Replace "fi" with appropriate language code
	defer C.free(unsafe.Pointer(cLang))

	handlePtr := C.voikkoInit((**C.char)(&errorStr), cLang, nil)
	if handlePtr == nil {
		if errorStr != nil {
			return handle, errors.New(C.GoString(errorStr))
		}
		return handle, errors.New("unknown error during initialization")
	}
	handle = VoikkoHandle(handlePtr)
	return handle, nil
}

// Batch lemmatizes all words in input slice. Returns a slice of string and an error.
func Batch(handle VoikkoHandle, words []string) ([]string, error) {
	result := make([]string, len(words))
	for i, word := range words {
		lemma, err := Single(handle, word)
		if err != nil {
			return nil, fmt.Errorf("error lemmatizing word \"%s\": %v", word, err)
		}
		result[i] = lemma
	}
	return result, nil
}

// Single lemmatizes the word and returns it and an error.
func Single(handle VoikkoHandle, word string) (string, error) {
	cWord := toWideCharString(word)
	defer C.free(unsafe.Pointer(cWord))

	// Perform morphology analysis
	cAnalysis := C.voikkoAnalyzeWordUcs4((*C.struct_VoikkoHandle)(handle), cWord)
	if cAnalysis == nil {
		return "", errors.New("failed to analyze word")
	}

	// If word doesn't have Voikko analysis information, function returns the word as it is.
	if *cAnalysis == nil {
		return word, nil
	}

	cKey := C.voikko_mor_analysis_keys(*cAnalysis)
	if cKey == nil {
		return "", errors.New("failed to analyze word")
	}

	cWordLemma := C.voikko_mor_analysis_value_ucs4(*cAnalysis, *cKey)
	if cWordLemma == nil {
		return "", errors.New("failed to analyze word")
	}
	wordLemma := wideCharToGoString(cWordLemma)

	// Free memory allocated by C function
	C.voikko_free_mor_analysis(cAnalysis)

	return wordLemma, nil
}

func Free(handle VoikkoHandle) {
	C.voikkoTerminate((*C.struct_VoikkoHandle)(handle))
}

// Convert Go string to wide character string
func toWideCharString(s string) *C.wchar_t {
	runes := []rune(s)
	ptr := C.malloc(C.size_t(len(runes)+1) * C.size_t(unsafe.Sizeof(C.wchar_t(0))))
	wStr := (*[1 << 30]C.wchar_t)(ptr)
	for i, r := range runes {
		wStr[i] = C.wchar_t(r)
	}
	wStr[len(runes)] = 0
	return (*C.wchar_t)(ptr)
}

// Convert wide character string to Go string
func wideCharToGoString(wc *C.wchar_t) string {
	length := 0
	for p := wc; *p != 0; p = (*C.wchar_t)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(*wc))) {
		length++
	}

	wStr := (*[1 << 30]C.wchar_t)(unsafe.Pointer(wc))
	runes := make([]rune, length)
	for i := 0; i < length; i++ {
		runes[i] = rune(wStr[i])
	}

	return string(runes)
}
