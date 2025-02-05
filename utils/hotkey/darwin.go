//go:build darwin

package hotkey

import "golang.design/x/hotkey"

func GetModKey() hotkey.Modifier {
	return hotkey.ModCmd
}
