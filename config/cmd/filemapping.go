package cmd

import (
	"errors"
	"fmt"
	"strings"
)

const mappingSeparator = ","
const sourceDestSeparator = ":"

type FileMapping struct {
	source      string
	destination string
}

type FileMappings struct {
	mappings []FileMapping
}

func (m *FileMappings) String() string {
	return ""
}

func (m *FileMappings) Set(value string) error {
	mappings := strings.Split(value, mappingSeparator)

	m.mappings = make([]FileMapping, len(mappings))
	for i, rawMapping := range mappings {
		mapping := strings.Split(rawMapping, sourceDestSeparator)
		if len(mapping) != 2 {
			return errors.New(fmt.Sprintf("Unable to create mapping from: %s", value))
		}

		m.mappings[i].source = mapping[0]
		m.mappings[i].destination = mapping[1]
	}

	return nil
}

func (m *FileMappings) Type() string {
	return "fileMappings"
}

func (m *FileMappings) Mappings() []FileMapping {
	return m.mappings
}

func (m *FileMappings) Sources() []string {
	sources := make([]string, len(m.mappings))
	for i, mapping := range m.mappings {
		sources[i] = mapping.source
	}

	return sources
}

func (m *FileMappings) Destinations() []string {
	destinations := make([]string, len(m.mappings))
	for i, mapping := range m.mappings {
		destinations[i] = mapping.destination
	}

	return destinations
}
