package main

import (
	"errors"
	"io"
	"os"
	"testing"
)

func TestAtomicWriteFile_WritesAndLeavesNoTemp(t *testing.T) {
	path := t.TempDir() + "/out.txt"

	if err := atomicWriteFile(path, func(w io.Writer) error {
		_, err := io.WriteString(w, "hello")
		return err
	}); err != nil {
		t.Fatalf("atomicWriteFile: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("content = %q, want %q", got, "hello")
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("temp file still present after success")
	}
}

func TestAtomicWriteFile_FailedWriteKeepsPreviousContent(t *testing.T) {
	// The point of the temp-file dance: a generated file (README.md,
	// history.jsonl) must never be truncated by a write that fails midway.
	path := t.TempDir() + "/out.txt"
	if err := os.WriteFile(path, []byte("previous"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	boom := errors.New("boom")
	err := atomicWriteFile(path, func(w io.Writer) error {
		_, _ = io.WriteString(w, "partial")
		return boom
	})
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want %v", err, boom)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "previous" {
		t.Errorf("content = %q, want the file left untouched", got)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("temp file left behind after a failed write")
	}
}

func TestAtomicWriteFile_UncreatableTempReportsError(t *testing.T) {
	if err := atomicWriteFile(t.TempDir()+"/missing-dir/out.txt", func(w io.Writer) error {
		return nil
	}); err == nil {
		t.Fatal("expected an error when the temp file cannot be created")
	}
}
