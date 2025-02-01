package backend

import (
	"context"
	"log"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
)

type HotkeyManager struct {
	ctx    context.Context
	hk     *hotkey.Hotkey
	active bool
}

func NewHotkeyManager() *HotkeyManager {
	return &HotkeyManager{}
}

func (h *HotkeyManager) Startup(ctx context.Context) {
	h.ctx = ctx
}

func (h *HotkeyManager) Register(hotkeyStr string) {
	if h.active {
		h.hk.Unregister()
	}

	mods, key := parseHotkey(hotkeyStr)
	h.hk = hotkey.New(mods, key)

	go func() {
		if err := h.hk.Register(); err != nil {
			log.Printf("Failed to register hotkey: %v", err)
			return
		}
		h.active = true
		log.Println("Hotkey registered:", hotkeyStr)

		for range h.hk.Keydown() {
			runtime.EventsEmit(h.ctx, "hotkey-pressed")
		}
	}()
}

func parseHotkey(hk string) ([]hotkey.Modifier, hotkey.Key) {
	parts := strings.Split(hk, "+")
	mods := make([]hotkey.Modifier, 0, len(parts)-1)

	var key hotkey.Key
	for _, part := range parts {
		switch strings.ToLower(part) {
		case "ctrl":
			mods = append(mods, hotkey.ModCtrl)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "alt":
			mods = append(mods, hotkey.ModAlt)
		default:
			key = hotkey.Key(part[0])
		}
	}
	return mods, key
}
