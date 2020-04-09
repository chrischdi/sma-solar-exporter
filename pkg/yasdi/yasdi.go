package yasdi

// #cgo CFLAGS: -g -I../../yasdi/libs -I../../yasdi/include -I../../yasdi/smalib
// #cgo LDFLAGS: -L../../yasdi/projects/generic-cmake/build-gcc -lyasdimaster
// #include <stdlib.h>
// #include "libyasdimaster.h"
import "C"
import (
	"fmt"
	"time"
	"unsafe"

	"github.com/pkg/errors"
)

// TChanType are the predefined channel types in yasdi/libs/libyasdimaster.h:120
type TChanType int

const (
	// MaxDeviceName is the maximum number of chars in a device name
	MaxDeviceName = 50
	// MaxChannelName is the maximum number of chars in a channel name
	MaxChannelName = 50
	// MaxChannelUnit is the maximum number of chars in a channel unit
	MaxChannelUnit = 17

	// MaxDeviceChannelHandles is the default value taken from yasdi/shell/CommonShellUIMain.c:130
	MaxDeviceChannelHandles = 300
)
const (
	// TChanType are the channel types available
	SPOTCHANNELS TChanType = iota
	PARAMCHANNELS
	TESTCHANNELS
	ALLCHANNELS
)

var (
	// ErrInvalidHandle is the error YE_UNKNOWN_HANDLE
	ErrInvalidHandle = errors.New("function called with invalid handle (dev or chan)")
	// ErrNotAllDevicesFound is the error YE_NOT_ALL_DEVS_FOUND
	ErrNotAllDevicesFound = errors.New("device detection, not all devices found...")
	// ErrYasdiIsShutdown is the error YE_SHUTDOWN
	ErrYasdiIsShutdown = errors.New("YASDI is in shutdwon mode. Function can't be called")
	// ErrTimeout is the error YE_TIMEOUT
	ErrTimeout = errors.New("an timeout while getting or setting channel value had occurred")
	// ErrNoValueRange is the error YE_NO_RANGE
	ErrNoValueRange = errors.New("GetChannelValRange: Channel has no value range...")
	// ErrInvalidChannelValue is the error YE_VALUE_NOT_VALID
	ErrInvalidChannelValue = errors.New("channel value is not valid..")
	// ErrPermissionDenied is the error YE_NO_ACCESS_RIGHTS
	ErrPermissionDenied = errors.New("you does not have needed right's to get the value...")
	// ErrOperationNotPossible is the error YE_CHAN_TYPE_MISMATCH
	ErrOperationNotPossible = errors.New("an operation is not possible on the channel")
	// ErrInvalidArgument is the error YE_INVAL_ARGUMENT
	ErrInvalidArgument = errors.New("function was called with an invalid argument or pointer")
	// ErrUnsupported is the error YE_NOT_SUPPORTED
	ErrUnsupported = errors.New("function not supported anymore...")
	// ErrDetectionAlreadyRunning is the error YE_DEV_DETECT_IN_PROGRESS
	ErrDetectionAlreadyRunning = errors.New("Device detection is already in progress")
	// ErrTooManyRequests is the error YE_TOO_MANY_REQUESTS
	ErrTooManyRequests = errors.New("Sync functions: too many requests by user API...")
	ErrUnknown         = errors.New("unknown error")

	errs = map[int]error{
		-1: ErrInvalidHandle,
		// -1: ErrNotAllDevicesFound,
		-2: ErrYasdiIsShutdown,
		-3: ErrTimeout,
		// -3: ErrNoValueRange,
		-4:  ErrInvalidChannelValue,
		-5:  ErrPermissionDenied,
		-6:  ErrOperationNotPossible,
		-7:  ErrInvalidArgument,
		-8:  ErrUnsupported,
		-9:  ErrDetectionAlreadyRunning,
		-20: ErrTooManyRequests,
	}
)

// YasdiMasterInitialize initializes YASDI by reading information from the INI-
// file and returns the number of found drivers.
// This function must be called before all others.
func YasdiMasterInitialize(iniFilePath string) (int, error) {
	cPath := C.CString(iniFilePath)
	defer C.free(unsafe.Pointer(cPath))
	cDriverNum := C.DWORD(0)
	ret := C.yasdiMasterInitialize(cPath, &cDriverNum)

	if int(ret) != 0 {
		if int(ret) == -1 {
			return 0, fmt.Errorf("file %q not found or not readable", iniFilePath)
		}
		return 0, ErrUnknown
	}
	return int(cDriverNum), nil
}

// YasdiMasterShutdown shutsdown YASDI.
// This function must be called after using YASDI.
func YasdiMasterShutdown() {
	C.yasdiMasterShutdown()
}

// YasdiReset reset YASDI to the state after loading the library.
func YasdiReset() {
	C.yasdiReset()
}

func YasdiMasterSetDriverOnline(driverID int) error {
	cDriverID := C.DWORD(driverID)
	ret := C.yasdiMasterSetDriverOnline(cDriverID)
	if int(ret) == 0 {
		return fmt.Errorf("error setting driver %d online", driverID)
	}
	return nil
}

// DoStartDeviceDetection tries to find the number of given devices.
func DoStartDeviceDetection(devices int, wait bool) error {
	ret := C.DoStartDeviceDetection(C.int(devices), C.int(boolToInt(wait)))
	err := errs[int(ret)]
	// in this case -1 is not InvalidHandle
	if errors.Cause(err) == ErrInvalidHandle {
		err = ErrNotAllDevicesFound
	}
	return err
}

// GetDeviceHandles returns all available device handles
func GetDeviceHandles(handleCount int) ([]int, error) {
	cHandles := make([]C.DWORD, handleCount)
	cHandleCount := C.DWORD(handleCount)
	ret := C.GetDeviceHandles(&cHandles[0], cHandleCount)
	if err := errs[int(ret)]; err != nil {
		return nil, err
	}

	deviceHandles := []int{}
	for _, handle := range cHandles {
		deviceHandles = append(deviceHandles, int(handle))
	}
	return deviceHandles, nil
}

// GetDeviceName returns the name of a device
func GetDeviceName(deviceHandle int) (string, error) {
	cDeviceHandle := C.DWORD(deviceHandle)
	buf := make([]C.char, MaxDeviceName+1)

	ret := C.GetDeviceName(cDeviceHandle, &buf[0], MaxDeviceName)
	if err := errs[int(ret)]; err != nil {
		return "", err
	}

	return C.GoString(&buf[0]), nil
}

// GetDeviceSN returns the serial number of a device
func GetDeviceSN(deviceHandle int) (uint, error) {
	cDeviceHandle := C.DWORD(deviceHandle)
	snBuffer := C.DWORD(0)
	ret := C.GetDeviceSN(cDeviceHandle, &snBuffer)
	if err := errs[int(ret)]; err != nil {
		return 0, err
	}

	return uint(snBuffer), nil
}

// FindChannelName searches for a channel handle for a device by a channel name
func FindChannelName(deviceHandle int, channelName string) (int, error) {
	cDeviceHandle := C.DWORD(deviceHandle)
	cChannelName := C.CString(channelName)
	ret := C.FindChannelName(cDeviceHandle, cChannelName)
	C.free(unsafe.Pointer(cChannelName))

	if err := errs[int(ret)]; err != nil {
		return 0, err
	}

	return int(ret), nil
}

// GetChannelHandlesEx returns all channel handles for a given device
func GetChannelHandlesEx(deviceHandle, maxChannelHandles int, t TChanType) ([]int, error) {
	cDeviceHandle := C.DWORD(deviceHandle)
	cMaxChannels := C.DWORD(maxChannelHandles)
	cChanType := C.TChanType(t)
	cChannelHandles := make([]C.DWORD, maxChannelHandles)

	ret := C.GetChannelHandlesEx(cDeviceHandle, &cChannelHandles[0], cMaxChannels, cChanType)
	if int(ret) == 0 {
		return nil, fmt.Errorf("no channel handles found for deviceHandle %d", deviceHandle)
	}
	if err := errs[int(ret)]; err != nil {
		return nil, err
	}

	channelHandles := []int{}
	for _, handle := range cChannelHandles {
		if int(handle) != 0 {
			channelHandles = append(channelHandles, int(handle))
		}
	}

	return channelHandles, nil
}

// GetChannelName returns the name of a channel
func GetChannelName(channelHandle int) (string, error) {
	cChannelHandle := C.DWORD(channelHandle)
	buf := make([]C.char, MaxChannelName+1)

	ret := C.GetChannelName(cChannelHandle, &buf[0], MaxChannelName)
	if err := errs[int(ret)]; err != nil {
		return "", err
	}

	return C.GoString(&buf[0]), nil
}

// GetChannelUnit returns the unit of a channel
func GetChannelUnit(channelHandle int) (string, error) {
	cChannelHandle := C.DWORD(channelHandle)
	buf := make([]C.char, MaxChannelUnit+1)

	ret := C.GetChannelUnit(cChannelHandle, &buf[0], MaxChannelUnit)
	if err := errs[int(ret)]; err != nil {
		return "", err
	}

	return C.GoString(&buf[0]), nil
}

// GetChannelValRange returns the value range for a channel
func GetChannelValRange(channelHandle int) (float64, float64, error) {
	cChannelHandle := C.DWORD(channelHandle)
	cMin := C.double(0)
	cMax := C.double(0)

	ret := C.GetChannelValRange(cChannelHandle, &cMin, &cMax)
	if err := errs[int(ret)]; err != nil {
		// in this case -3 is not Timeout
		if errors.Cause(err) == ErrTimeout {
			err = ErrNoValueRange
		}
		return 0, 0, err
	}

	return float64(cMin), float64(cMax), nil
}

// GetChannelValue returns the value of a channel
func GetChannelValue(channelHandle, deviceHandle, maxAge int) (float64, string, error) {
	cChannelHandle := C.DWORD(channelHandle)
	cDeviceHandle := C.DWORD(deviceHandle)
	cValue := C.double(0)
	buf := make([]C.char, 17)

	ret := C.GetChannelValue(cChannelHandle, cDeviceHandle, &cValue, &buf[0], C.DWORD(16), C.DWORD(maxAge))
	if err := errs[int(ret)]; err != nil {
		return 0, "", err
	}

	return float64(cValue), C.GoString(&buf[0]), nil
}

// GetChannelValueTimeStamp returns the time of a channel value
func GetChannelValueTimeStamp(channelHandle, deviceHandle int) (time.Time, error) {
	cChannelHandle := C.DWORD(channelHandle)
	cDeviceHandle := C.DWORD(deviceHandle)

	ret := C.GetChannelValueTimeStamp(cChannelHandle, cDeviceHandle)
	if int64(ret) == 0 {
		return time.Time{}, ErrInvalidHandle
	}

	return time.Unix(int64(ret), 0), nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
