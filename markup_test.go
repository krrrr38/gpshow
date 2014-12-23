package main

import (
	"strings"
	"testing"
)

func TestMakeSlide(t *testing.T) {
	pages := MakeSlide([]byte("!SLIDE\nyuno"))
	if strings.Count(string(pages[:]), "container") != 1 {
		t.Error("should have one page")
	}
	if !strings.Contains(string(pages[:]), "<div class=\"container\">") {
		t.Error("should contain container div")
	}
	if !strings.Contains(string(pages[:]), "yuno") {
		t.Error("should contain content")
	}

	pages = MakeSlide([]byte("!SLIDE\nyuno\n!SLIDE\nmiyako"))
	if strings.Count(string(pages[:]), "container") != 2 {
		t.Error("should have two page")
	}

	pages = MakeSlide([]byte("!SLIDE\nyuno\n  !SLIDE\nmiyako"))
	if strings.Count(string(pages[:]), "container") != 1 {
		t.Error("should have one page")
	}
}
