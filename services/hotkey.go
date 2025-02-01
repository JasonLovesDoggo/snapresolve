package services

import (
	"context"
	"fmt"
	"strings"

	"golang.design/x/hotkey"
)

type HotkeyService struct {
	ctx      context.Context
	hk       *hotkey.Hotkey
	callback func()
}

func NewHotkeyService() *HotkeyService {
	return &HotkeyService{}
}

func (h *HotkeyService) Init(ctx context.Context, callback func()) {
	h.ctx = ctx
	h.callback = callback
}

func (h *HotkeyService) Register(hotkeyStr string) error {
	if h.hk != nil {
		h.hk.Unregister()
	}

	mods, key := h.parseHotkey(hotkeyStr)
	h.hk = hotkey.New(mods, key)

	if err := h.hk.Register(); err != nil {
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	go func() {
		for range h.hk.Keydown() {
			h.callback()
		}
	}()

	return nil
}

func (h *HotkeyService) Cleanup() {
	if h.hk != nil {
		h.hk.Unregister()
	}
}

func (h *HotkeyService) parseHotkey(hk string) ([]hotkey.Modifier, hotkey.Key) {
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
