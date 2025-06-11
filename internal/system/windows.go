//go:build windows

package system

import (
	"errors"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modKernel32 = windows.NewLazySystemDLL("kernel32.dll")
	modUser32   = windows.NewLazySystemDLL("user32.dll")
	modShell32  = windows.NewLazySystemDLL("shell32.dll")

	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
	procEnumWindows              = modUser32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = modUser32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible          = modUser32.NewProc("IsWindowVisible")
	procGetWindowTextLengthW     = modUser32.NewProc("GetWindowTextLengthW")
	procGetWindowTextW           = modUser32.NewProc("GetWindowTextW")
	procShellExecuteW            = modShell32.NewProc("ShellExecuteW")

	cb = syscall.NewCallback(enumWindowsCallback)
)

func enumWindowsCallback(hWnd windows.HWND, lParam uintptr) uintptr {
	visible, _, _ := procIsWindowVisible.Call(uintptr(hWnd))
	if visible == 0 {
		return 1
	}

	var pid uint32
	procGetWindowThreadProcessId.Call(uintptr(hWnd), uintptr(unsafe.Pointer(&pid)))

	searchData := (*WindowSearchData)(unsafe.Pointer(lParam))
	if pid == searchData.PID {
		searchData.Window = hWnd
		return 0
	}
	return 1
}

const (
	TH32CS_SNAPPROCESS = 0x00000002
	SW_SHOWNORMAL      = 1
)

type ProcessEntry32 struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	Threads           uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [260]uint16
}

type WindowSearchData struct {
	PID    uint32
	Window windows.HWND
}

type WindowsSystemController struct{}

func NewSystemController() (*WindowsSystemController, error) {
	return &WindowsSystemController{}
}

func (w *WindowsSystemController) getSpotifyWindow() (windows.HWND, error) {
	snapshot, _, _ := procCreateToolhelp32Snapshot.Call(TH32CS_SNAPPROCESS, 0)
	if snapshot == 0 || snapshot == uintptr(windows.InvalidHandle) {
		return 0, errors.New("failed to create snapshot")
	}
	defer syscall.CloseHandle(syscall.Handle(snapshot))

	var entry ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	r1, _, _ := procProcess32First.Call(snapshot, uintptr(unsafe.Pointer(&entry)))
	if r1 == 0 {
		return 0, errors.New("failed to get first process")
	}

	var spotifyPID uint32
	for {
		exe := syscall.UTF16ToString(entry.ExeFile[:])
		if strings.EqualFold(exe, "Spotify.exe") {
			spotifyPID = entry.ProcessID
			break
		}
		r, _, _ := procProcess32Next.Call(snapshot, uintptr(unsafe.Pointer(&entry)))
		if r == 0 {
			break
		}
	}

	if spotifyPID == 0 {
		return 0, errors.New("failed to find Spotify process")
	}

	searchData := WindowSearchData{
		PID:    spotifyPID,
		Window: 0,
	}

	procEnumWindows.Call(cb, uintptr(unsafe.Pointer(&searchData)))
	return searchData.Window, nil
}

func (w *WindowsSystemController) GetCurrentPlayingTrackTitle() (string, error) {
	hwnd, err := w.getSpotifyWindow()
	if hwnd == 0 {
		return "", err
	}

	length, _, _ := procGetWindowTextLengthW.Call(uintptr(hwnd))
	if length == 0 {
		return "", errors.New("failed to get window text length")
	}

	buf := make([]uint16, length+1)
	procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), length+1)
	windowTitle := syscall.UTF16ToString(buf)
	if windowTitle == "Spotify Premium" || windowTitle == "Spotify" {
		return "", errors.New("Spotify is running but not playing a track")
	}

	return windowTitle, nil
}

func (w *WindowsSystemController) OpenURLInBrowser(url string) error {
	urlPtr, _ := syscall.UTF16PtrFromString(url)
	openStr, _ := syscall.UTF16PtrFromString("open")
	r, _, _ := procShellExecuteW.Call(0, uintptr(unsafe.Pointer(openStr)),
		uintptr(unsafe.Pointer(urlPtr)), 0, 0, SW_SHOWNORMAL)
	if r <= 32 {
		return errors.New("failed to open URL in browser")
	}
	return nil
}
