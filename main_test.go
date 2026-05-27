// Podplane <https://podplane.dev>
// Copyright 2026 Nadrama Pty Ltd
// SPDX-License-Identifier: Apache-2.0

package main

import (
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
