package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/greensea/so/common"
	"github.com/greensea/so/sys"
)

func TestExtractCodeBlocks1(t *testing.T) {
	raw := "```\nfoo\n```\n```\nbar\n```"
	cmds, err := extractCodeBlocks([]byte(raw))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if cmds[0] != "foo" {
		t.Fatalf("Expected foo, got %s", cmds[0])
	}
	if cmds[1] != "bar" {
		t.Fatalf("Expected bar, got %s", cmds[1])
	}
}

func TestExtractCodeBlocks2(t *testing.T) {
	raw := "```bash\nfoo\n```\n```bash\nbar\n```"
	cmds, err := extractCodeBlocks([]byte(raw))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if cmds[0] != "foo" {
		t.Fatalf("Expected foo, got %s", cmds[0])
	}
	if cmds[1] != "bar" {
		t.Fatalf("Expected bar, got %s", cmds[1])
	}
}

func TestExtractCodeBlocks3(t *testing.T) {
	raw := "```bash\nfoo\nbar\n```\n```bash\nfoobar\n```"
	cmds, err := extractCodeBlocks([]byte(raw))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if cmds[0] != "foo\nbar" {
		t.Fatalf("Expected foo\nbar, got %s", cmds[0])
	}

}

func TestLang(t *testing.T) {
	os.Setenv("LANG", "zh")
	if common.Lang() != "zh" {
		t.Fatalf("Expected zh, got %s", common.Lang())
	}
}

func TestExecCmd(t *testing.T) {
	s := sys.Get()
	if s.Shell == "" {
		s.Shell = "sh"
	}

	input := "for i in {1..2}; do echo $i; done"
	buf := new(bytes.Buffer)
	err := ExecCmd(input, buf, nil)
	if err != nil {
		t.Fatalf("Command failed: %v", err)
	}

	if buf.String() != "1\n2\n" {
		t.Fatalf("Expected 1\n2\n, got %s", buf.String())
	}
}

func TestExecCmd2(t *testing.T) {
	s := sys.Get()
	if s.Shell == "" {
		s.Shell = "sh"
	}

	input := "for i in {1..2}; do echo $i; done | wc -l"
	buf := new(bytes.Buffer)
	err := ExecCmd(input, buf, nil)
	if err != nil {
		t.Fatalf("Command failed: %v", err)
	}

	if buf.String() != "2\n" {
		t.Fatalf("Expected `2\n', got `%s'", buf.String())
	}
}

func TestExecCmd3(t *testing.T) {
	s := sys.Get()
	if s.Shell == "" {
		s.Shell = "sh"
	}

	input := `find /not_exist -type f -exec ls {} \; | wc -l`
	// input := `find . -type f -mtime +10 -print | wc -l`
	buf := new(bytes.Buffer)
	err := ExecCmd(input, buf, nil)
	if err != nil {
		t.Fatalf("Command failed: %v", err)
	}

	if buf.String() != "0\n" {
		t.Fatalf("Expected `0\n', got `%s'", buf.String())
	}
}
