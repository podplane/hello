// Podplane <https://podplane.dev>
// Copyright 2026 Nadrama Pty Ltd
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderIndexUsesConfiguredMessage(t *testing.T) {
	body, err := renderIndex(`Hello <@TestRenderIndexUsesConfiguredMessage>!`)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(body), "TestRenderIndexUsesConfiguredMessage") {
		t.Fatalf("body does not contain custom message:\n%s", body)
	}

	if !strings.Contains(string(body), "Hello &lt;@TestRenderIndexUsesConfiguredMessage&gt;!") {
		t.Fatalf("body does not contain escaped message:\n%s", body)
	}
}

func TestHelloMessageReadsAbsoluteFilePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "message.txt")
	if err := os.WriteFile(path, []byte("Secret <message>\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HELLO_MESSAGE", path)

	message := helloMessage()
	if message != "Secret <message>\n" {
		t.Fatalf("message = %q, want file contents", message)
	}

	body, err := renderIndex(message)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "Secret &lt;message&gt;") {
		t.Fatalf("body does not contain escaped file message:\n%s", body)
	}
}

func TestHelloMessageReturnsReadErrorForMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.txt")
	t.Setenv("HELLO_MESSAGE", path)

	message := helloMessage()
	if !strings.Contains(message, path) || !strings.Contains(message, "no such file") {
		t.Fatalf("message = %q, want missing file error", message)
	}
}

func TestHelloMessageRejectsOversizedFiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "message.txt")
	contents := strings.Repeat("a", maxMessageFileSize+1)
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HELLO_MESSAGE", path)

	message := helloMessage()
	if !strings.Contains(message, path) || !strings.Contains(message, "larger than 200 KiB") {
		t.Fatalf("message = %q, want size limit error", message)
	}
}

func TestHelloMessageRejectsBinaryData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "message.txt")
	if err := os.WriteFile(path, []byte{'h', 'i', 0x00}, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HELLO_MESSAGE", path)

	message := helloMessage()
	if !strings.Contains(message, path) || !strings.Contains(message, "binary or non-UTF-8 data") {
		t.Fatalf("message = %q, want binary data error", message)
	}
}
