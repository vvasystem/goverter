package config

import (
	"fmt"
	"sort"

	"github.com/jmattheis/goverter/pkgload"
)

type RawLines struct {
	Location string
	Lines    []string
}

type RawConverter struct {
	Package       string
	InterfaceName string
	Converter     RawLines
	Methods       map[string]RawLines
	FileSource    string
}

type Raw struct {
	Converters []RawConverter
	Global     RawLines
}

func Parse(workDir string, raw *Raw) ([]*Converter, error) {
	loader, err := pkgload.New(workDir, getPackages(raw))
	if err != nil {
		return nil, err
	}

	global, err := parseGlobal(loader, raw.Global)
	if err != nil {
		return nil, err
	}

	converters := []*Converter{}
	for _, rawConverter := range raw.Converters {
		converter, err := parseConverter(loader, &rawConverter, *global)
		if err != nil {
			return nil, err
		}
		converters = append(converters, converter)
	}

	sort.Slice(converters, func(i, j int) bool {
		return converters[i].Name < converters[j].Name
	})

	return converters, nil
}

func formatLineError(lines RawLines, t, value string, err error) error {
	cmd, _ := parseCommand(value)
	msg := `error parsing 'goverter:%s' at
    %s
    %s

%s`
	return fmt.Errorf(msg, cmd, lines.Location, t, err)
}
