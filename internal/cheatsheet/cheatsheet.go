package cheatsheet

import "fmt"

type Sheet struct {
	Name        string
	Aliases     []string
	Description string
	Content     string
}

var sheets = []Sheet{
	cpsSheet,
	goSheet,
	uvSheet,
	fnmSheet,
	rustSheet,
	tmuxSheet,
	nvimSheet,
	fzfSheet,
	regexSheet,
}

func List() []Sheet {
	return sheets
}

func Get(name string) *Sheet {
	for i := range sheets {
		if sheets[i].Name == name {
			return &sheets[i]
		}
		for _, alias := range sheets[i].Aliases {
			if alias == name {
				return &sheets[i]
			}
		}
	}
	return nil
}

func AllNames() []string {
	var names []string
	for _, s := range sheets {
		names = append(names, s.Name)
		names = append(names, s.Aliases...)
	}
	return names
}

func AvailableList() string {
	s := ""
	for _, sheet := range sheets {
		if s != "" {
			s += ", "
		}
		s += sheet.Name
	}
	return s
}

func Print(name string) error {
	sheet := Get(name)
	if sheet == nil {
		return fmt.Errorf("unknown cheat sheet: %s (available: %s)", name, AvailableList())
	}
	fmt.Print(sheet.Content)
	return nil
}
