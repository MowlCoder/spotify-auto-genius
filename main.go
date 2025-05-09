package main

import (
	"errors"
	"log"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
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

func findSpotifyWindow() windows.HWND {
	snapshot, _, _ := procCreateToolhelp32Snapshot.Call(TH32CS_SNAPPROCESS, 0)
	if snapshot == 0 || snapshot == uintptr(windows.InvalidHandle) {
		return 0
	}
	defer syscall.CloseHandle(syscall.Handle(snapshot))

	var entry ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	r1, _, _ := procProcess32First.Call(snapshot, uintptr(unsafe.Pointer(&entry)))
	if r1 == 0 {
		return 0
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
		return 0
	}

	searchData := WindowSearchData{
		PID:    spotifyPID,
		Window: 0,
	}

	procEnumWindows.Call(cb, uintptr(unsafe.Pointer(&searchData)))
	return searchData.Window
}

func getSpotifyWindowTitle() string {
	hwnd := findSpotifyWindow()
	if hwnd == 0 {
		return ""
	}

	length, _, _ := procGetWindowTextLengthW.Call(uintptr(hwnd))
	if length == 0 {
		return ""
	}

	buf := make([]uint16, length+1)
	procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), length+1)
	return syscall.UTF16ToString(buf)
}

func openURLInBrowser(url string) bool {
	urlPtr, _ := syscall.UTF16PtrFromString(url)
	openStr, _ := syscall.UTF16PtrFromString("open")
	r, _, _ := procShellExecuteW.Call(0, uintptr(unsafe.Pointer(openStr)),
		uintptr(unsafe.Pointer(urlPtr)), 0, 0, SW_SHOWNORMAL)
	return r > 32
}

func openGeniusPage(title string) error {
	url := "https://genius.com/search?q=" + title

	browser := rod.New().Timeout(10 * time.Second)
	if err := browser.Connect(); err != nil {
		return errors.New("failed to connect to browser: " + err.Error())
	}
	defer browser.Close()

	page, err := browser.Page(proto.TargetCreateTarget{URL: url})
	if err != nil {
		return errors.New("failed to create page: " + err.Error())
	}

	if err := page.WaitLoad(); err != nil {
		return errors.New("failed to load page: " + err.Error())
	}

	html, err := page.HTML()
	if err != nil {
		return errors.New("failed to get HTML: " + err.Error())
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("Failed to parse HTML: %v", err)
		log.Println("Opening search page", url)
		openURLInBrowser(url)
		return nil
	}

	href, exists := doc.Find("a.mini_card").First().Attr("href")
	if !exists {
		log.Println("No exact match found, opening search page:", url)
		openURLInBrowser(url)
	} else {
		log.Println("Found exact match:", href)
		openURLInBrowser(href)
	}

	return nil
}

func main() {
	prevTitle := ""
	log.Println("Starting scanning Spotify...")

	for {
		title := getSpotifyWindowTitle()
		if title == "" {
			log.Println("Spotify is not running...Waiting 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		if title != prevTitle && title != "Spotify Premium" && title != "Spotify" {
			log.Println("New track:", title)
			prevTitle = title

			if err := openGeniusPage(title); err != nil {
				log.Printf("Failed to open Genius page: %v", err)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
