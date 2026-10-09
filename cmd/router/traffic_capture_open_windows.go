//go:build windows

package main

import (
	"fmt"
	"os"
)

// openCaptureFile opens the JSONL capture file. Windows exposes no O_NOFOLLOW
// equivalent in syscall, and CreateFile follows reparse points, so a symlinked
// capture path is rejected before the open instead of being redirected into the
// link target. The check is a pre-open inspection, so it closes the accidental
// redirection case rather than a determined local race.
func openCaptureFile(path string) (*os.File, error) {
	if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("HTTP capture path %q is a symbolic link", path)
	}
	return os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		trafficCaptureFilePermissions,
	)
}
