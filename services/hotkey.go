package services

import (
	"fmt"
	"golang.design/x/hotkey"
	"strings"
)

type HotkeyService struct {
	hk      *hotkey.Hotkey
	handler func()
}

func NewHotkeyService(hotkeyStr string, handler func()) (*HotkeyService, error) {
	mods, key := parseHotkey(hotkeyStr)

	hk := hotkey.New(mods, key)
	if err := hk.Register(); err != nil {
		return nil, fmt.Errorf("failed to register hotkey: %w", err)
	}

	service := &HotkeyService{
		hk:      hk,
		handler: handler,
	}

	// Start listening for hotkey events
	go service.listen()

	return service, nil
}

func parseHotkey(hk string) ([]hotkey.Modifier, hotkey.Key) {
	parts := strings.Split(strings.ToLower(hk), "+")
	var mods []hotkey.Modifier

	for i := 0; i < len(parts)-1; i++ {
		switch parts[i] {
		case "ctrl":
			mods = append(mods, hotkey.ModCtrl)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "alt":
			mods = append(mods, hotkey.ModAlt)
		}
	}

	key := hotkey.Key(parts[len(parts)-1][0])
	return mods, key
}

func (h *HotkeyService) listen() {
	for range h.hk.Keydown() {
		fmt.Println("Detected hotkey press", h.hk.String())
		h.handler()
	}
}

func (h *HotkeyService) Stop() {
	if h.hk != nil {
		h.hk.Unregister()
	}
}
