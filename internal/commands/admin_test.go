package commands

import (
	"testing"
)

func TestAdminCmd_Definition(t *testing.T) {
	if AdminCmd == nil {
		t.Fatal("AdminCmd should not be nil")
	}

	if AdminCmd.Use != "admin" {
		t.Errorf("AdminCmd.Use = %q, want %q", AdminCmd.Use, "admin")
	}

	if AdminCmd.Short == "" {
		t.Error("AdminCmd.Short should not be empty")
	}

	if AdminCmd.Run == nil {
		t.Error("AdminCmd.Run should not be nil")
	}
}

func TestReviewCmd_Definition(t *testing.T) {
	if ReviewCmd == nil {
		t.Fatal("ReviewCmd should not be nil")
	}

	if ReviewCmd.Use != "review" {
		t.Errorf("ReviewCmd.Use = %q, want %q", ReviewCmd.Use, "review")
	}

	if ReviewCmd.Short == "" {
		t.Error("ReviewCmd.Short should not be empty")
	}

	if ReviewCmd.Run == nil {
		t.Error("ReviewCmd.Run should not be nil")
	}
}
