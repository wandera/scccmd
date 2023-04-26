package cmd

import (
	"fmt"
	"strings"
)

const (
	mappingSeparator    = ","
	sourceDestSeparator = ":"
)

// FileMapping single file mapping source:dest
type FileMapping struct {
	source      string
	destination string
}

// FileMappings file mappings source:dest,source:dest...
type FileMappings struct {
	mappings []FileMapping
}

func (m *FileMappings) String() string {
	return ""
}

// Set parse mappings from string
func (m *FileMappings) Set(value string) error {
	mappings := strings.Split(value, mappingSeparator)

	m.mappings = make([]FileMapping, len(mappings))
	for i, rawMapping := range mappings {
		mapping := strings.Split(rawMapping, sourceDestSeparator)
		if len(mapping) != 2 {
			return fmt.Errorf("unable to create mapping from: %s", value)
		}

		m.mappings[i].source = mapping[0]
		m.mappings[i].destination = mapping[1]
	}

	return nil
}

// Type type name (for cobra)
func (m *FileMappings) Type() string {
	return "FileMappings"
}

// Mappings all mappings
func (m *FileMappings) Mappings() []FileMapping {
	return m.mappings
}

// Sources all sources
func (m *FileMappings) Sources() []string {
	sources := make([]string, len(m.mappings))
	for i, mapping := range m.mappings {
		sources[i] = mapping.source
	}

	return sources
}

// Destinations all destinations
func (m *FileMappings) Destinations() []string {
	destinations := make([]string, len(m.mappings))
	for i, mapping := range m.mappings {
		destinations[i] = mapping.destination
	}

	return destinations
}
