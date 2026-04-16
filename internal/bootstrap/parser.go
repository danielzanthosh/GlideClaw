package bootstrap

import (
	"bufio"
	"os"
	"strings"
)

type Profile struct {
	Identity                 string
	SecurityMode             string
	AllowedAutonomousActions []string
	ConfirmationRequired     []string
	BlockedActions           []string
	MemoryHints              []string
	RawSections              map[string][]string
}

func Load(path string) (Profile, error) {
	p := Profile{SecurityMode: "strict", RawSections: map[string][]string{}}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return p, nil
		}
		return p, err
	}
	defer f.Close()

	section := ""
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "## ") {
			section = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			continue
		}
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "- ") {
			v := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			p.RawSections[section] = append(p.RawSections[section], v)
		}
	}

	p.Identity = first(p.RawSections["Identity"])
	if m := first(p.RawSections["Security mode"]); m != "" {
		p.SecurityMode = m
	}
	p.AllowedAutonomousActions = p.RawSections["Allowed autonomous actions"]
	p.ConfirmationRequired = p.RawSections["Confirmation-required actions"]
	p.BlockedActions = p.RawSections["Blocked actions"]
	p.MemoryHints = p.RawSections["Memory hints"]
	return p, s.Err()
}

func first(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[0]
}
