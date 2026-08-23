package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ToolState struct {
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installed_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type State struct {
	mu       sync.Mutex
	Tools    map[string]ToolState `json:"tools"`
	Theme    string               `json:"theme,omitempty"`
	Version  string               `json:"cps_version,omitempty"`
	LastInit time.Time            `json:"last_init"`
	path     string
}

func (s *State) CPSVersion() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Version
}

func (s *State) SetCPSVersion(v string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Version = v
}

func (s *State) IsFresh() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.Tools) == 0
}

func (s *State) CurrentTheme() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Theme
}

func (s *State) SetTheme(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Theme = name
}

func Load(path string) (*State, error) {
	s := &State{
		Tools: make(map[string]ToolState),
		path:  path,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	s.path = path
	if s.Tools == nil {
		s.Tools = make(map[string]ToolState)
	}
	return s, nil
}

func (s *State) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *State) ToolVersion(name string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if t, ok := s.Tools[name]; ok {
		return t.Version
	}
	return ""
}

func (s *State) SetToolVersion(name, version string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	if t, ok := s.Tools[name]; ok {
		t.Version = version
		t.UpdatedAt = now
		s.Tools[name] = t
	} else {
		s.Tools[name] = ToolState{
			Version:     version,
			InstalledAt: now,
			UpdatedAt:   now,
		}
	}
}

func (s *State) Remove(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Tools, name)
}
