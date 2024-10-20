package sys

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Sys struct {
	Uname  string
	WhoAmI string
	Shell  string
	Pwd    string
}

var s *Sys = nil

func Get() *Sys {
	if s != nil {
		return s
	}

	s = &Sys{}

	s.Shell = os.Getenv("SHELL")

	ret, err := exec.Command("id", "-un").Output()
	if err == nil && len(ret) > 0 {
		s.WhoAmI = string(ret)
	}

	ret, err = exec.Command("uname", "-a").Output()
	if err == nil && len(ret) > 0 {
		s.Uname = string(ret)
	}

	s.Pwd, _ = os.Getwd()

	s.Shell = strings.TrimSpace(s.Shell)
	s.WhoAmI = strings.TrimSpace(s.WhoAmI)
	s.Uname = strings.TrimSpace(s.Uname)
	s.Pwd = strings.TrimSpace(s.Pwd)

	shellParts := strings.Split(s.Shell, "/")
	if len(shellParts) > 0 {
		s.Shell = shellParts[len(shellParts)-1]
	}

	return s
}

func (s *Sys) Text() string {
	ret := []string{}
	if s.Uname != "" {
		ret = append(ret, fmt.Sprintf("System Info: %s", s.Uname))
	}
	if s.WhoAmI != "" {
		ret = append(ret, fmt.Sprintf("Current User: %s", s.WhoAmI))
	}
	if s.Pwd != "" {
		ret = append(ret, fmt.Sprintf("Working Directory: %s", s.Pwd))
	}
	if s.Shell != "" {
		ret = append(ret, fmt.Sprintf("Shell: %s", s.Shell))
	}
	return strings.Join(ret, "\n")
}
