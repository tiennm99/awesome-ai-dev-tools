package main

import (
	"io"
	"os"
)

// atomicWriteFile writes through a temp file in the same directory and renames
// it over path, so a crash or a failed write never leaves a generated file
// truncated or half-written. The callback receives the temp file as an
// io.Writer, which suits streaming producers (templates, JSON encoders)
// without buffering the whole payload in memory.
//
// Errors from Sync and Close are reported too: on a write path they are the
// difference between "data reached the disk" and silent loss.
func atomicWriteFile(path string, write func(io.Writer) error) error {
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	writeErr := write(f)

	if syncErr := f.Sync(); syncErr != nil && writeErr == nil {
		writeErr = syncErr
	}
	if closeErr := f.Close(); closeErr != nil && writeErr == nil {
		writeErr = closeErr
	}

	if writeErr != nil {
		_ = os.Remove(tmp) // best-effort cleanup; the write error is what matters
		return writeErr
	}

	return os.Rename(tmp, path)
}
