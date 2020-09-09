package mediainfo

// #cgo CFLAGS: -DUNICODE
// #cgo LDFLAGS: -lz -lzen -lpthread -lstdc++ -lmediainfo -ldl
// #include "go_mediainfo.h"
// #include <MediaInfoDLL/MediaInfoDLL.h>
import "C"

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strconv"
	"unsafe"
)

var (
	ErrorOpenFile = errors.New("Failed open file")
	ErrorMemmoryOpen = errors.New("Failed get mediainfo from memmory")
)

// MediaInfo - represents MediaInfo class, all interaction with libmediainfo through it
type MediaInfo struct {
	handle unsafe.Pointer
}

const (
	StreamGeneral int = iota
	StreamVideo
	StreamAudio
	StreamText
	StreamOther
	StreamImage
	StreamMenu
)

func init() {
	C.setlocale(C.LC_CTYPE, C.CString(""))
	C.MediaInfoDLL_Load()

	if C.MediaInfoDLL_IsLoaded() == 0 {
		panic("Cannot load mediainfo")
	}
}

// NewMediaInfo - constructs new MediaInfo
func NewMediaInfo() *MediaInfo {
	result := &MediaInfo{handle: C.GoMediaInfo_New()}
	runtime.SetFinalizer(result, func(h *MediaInfo) {
		C.GoMediaInfo_Close(h.handle)
		C.GoMediaInfo_Delete(h.handle)
	})
	return result
}

// Open - get MediaInfo for given file path
func Open(path string) (mi *MediaInfo, err error) {
	handle := C.GoMediaInfo_New()
	p :=  C.CString(path)
	defer C.free(unsafe.Pointer(p))
	if C.GoMediaInfo_OpenFile(handle, p) != 1 {
		return nil, fmt.Errorf("%w: %s", ErrorOpenFile, path)
	}

	mi =&MediaInfo{handle: handle}
	runtime.SetFinalizer(mi, func(h *MediaInfo) {
		C.GoMediaInfo_Close(h.handle)
		C.GoMediaInfo_Delete(h.handle)
		// TODO: check for memmory leak
	})
	return
}

// Read - use generic io.Reader as data source. Dangerouse beacous of
// usage inmemory buffer
func Read(r io.Reader) (mi *MediaInfo, err error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	handle := C.GoMediaInfo_New()
	rc := C.GoMediaInfo_OpenMemory(handle,
		(*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
	if rc != 0 {
		return nil, ErrorMemmoryOpen
	}
	mi =&MediaInfo{handle: handle}
	runtime.SetFinalizer(mi, func(h *MediaInfo) {
		C.GoMediaInfo_Close(h.handle)
		C.GoMediaInfo_Delete(h.handle)
		// TODO: check for memmory leak
	})
	return
}

// OpenFile - opens file
func (mi *MediaInfo) OpenFile(path string) error {
	p := C.CString(path)
	s := C.GoMediaInfo_OpenFile(mi.handle, p)
	if s == 0 {
		return fmt.Errorf("MediaInfo can't open file: %s", path)
	}
	C.free(unsafe.Pointer(p))
	return nil
}

// OpenMemory - opens memory buffer
func (mi *MediaInfo) OpenMemory(bytes []byte) error {
	if len(bytes) == 0 {
		return fmt.Errorf("Buffer is empty")
	}
	s := C.GoMediaInfo_OpenMemory(mi.handle, (*C.char)(unsafe.Pointer(&bytes[0])), C.size_t(len(bytes)))
	if s == 0 {
		return fmt.Errorf("MediaInfo can't open memory buffer")
	}
	return nil
}

// Close - closes file
func (mi *MediaInfo) Close() {
	C.GoMediaInfo_Close(mi.handle)
}

// Get - allow to read info from file
func (mi *MediaInfo) Get(param string) (result string) {
	p := C.CString(param)
	r := C.GoMediaInfoGet(mi.handle, p)
	result = C.GoString(r)
	C.free(unsafe.Pointer(p))
	C.free(unsafe.Pointer(r))
	return
}

// GetStream - get single stream information (only for the first one)
func (mi *MediaInfo) GetStream(stream int, param string) (result string) {
	p := C.CString(param)
	defer C.free(unsafe.Pointer(p))
	resp := C.GoMediaInfoStreamGet(mi.handle, C.MediaInfo_stream_C(stream), p)
	result = C.GoString(resp)
	C.free(unsafe.Pointer(resp))
	return
}

// Inform returns string with summary file information, like mediainfo util
func (mi *MediaInfo) Inform() (result string) {
	r := C.GoMediaInfoInform(mi.handle)
	result = C.GoString(r)
	C.free(unsafe.Pointer(r))
	return
}

// Option configure or get information about MediaInfoLib
func (mi *MediaInfo) Option(option string, value string) (result string) {
	o := C.CString(option)
	v := C.CString(value)
	r := C.GoMediaInfoOption(mi.handle, o, v)
	C.free(unsafe.Pointer(o))
	C.free(unsafe.Pointer(v))
	result = C.GoString(r)
	C.free(unsafe.Pointer(r))
	return
}

// AvailableParameters returns string with all available Get params and it's descriptions
func (mi *MediaInfo) AvailableParameters() string {
	return mi.Option("Info_Parameters", "")
}

// Duration returns file duration
func (mi *MediaInfo) Duration() int {
	duration, _ := strconv.Atoi(mi.Get("Duration"))
	return duration
}

// Codec returns file codec
func (mi *MediaInfo) Codec() string {
	return mi.Get("Format")
}

// Format returns file codec
func (mi *MediaInfo) Format() string {
	return mi.Get("Format")
}
