package bridge

import (
	"context"
	"errors"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var errNativeNotWired = errors.New("native adapter is not wired yet")

type SystemBridge struct { ctx context.Context }

func NewSystemBridge() *SystemBridge { return &SystemBridge{} }
func (b *SystemBridge) SetContext(ctx context.Context) { b.ctx = ctx }

func (b *SystemBridge) ChooseDirectory(title string) (string, error) {
	return runtime.OpenDirectoryDialog(b.ctx, runtime.OpenDialogOptions{Title: title})
}

func (b *SystemBridge) OpenExternal(url string) { runtime.BrowserOpenURL(b.ctx, url) }
func (b *SystemBridge) RevealPath(path string) error { return errNativeNotWired }
func (b *SystemBridge) SecureGet(key string) (string, error) { return "", errNativeNotWired }
func (b *SystemBridge) SecureSet(key, value string) error { return errNativeNotWired }
func (b *SystemBridge) SecureDelete(key string) error { return errNativeNotWired }
