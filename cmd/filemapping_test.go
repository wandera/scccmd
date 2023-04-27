package cmd

import "testing"

func Test0FileMapping(t *testing.T) {
	mapping := FileMapping{}

	if mapping.source != "" || mapping.destination != "" {
		t.Error("Mapping does not initialize on empty values")
	}
}

func TestFileMapping(t *testing.T) {
	mapping := FileMapping{source: "src", destination: "dest"}

	if mapping.source != "src" {
		t.Error("Mapping does not contain source")
	}

	if mapping.destination != "dest" {
		t.Error("Mapping does not contain destination")
	}
}

func Test0FileMappings(t *testing.T) {
	mappings := FileMappings{}

	if len(mappings.Mappings()) != 0 {
		t.Error("Mappings does not initialize on empty array")
	}
}

func TestFileMappings(t *testing.T) {
	mappings := FileMappings{}
	mappings.Set("application.yaml:config.yaml")

	if len(mappings.Mappings()) != 1 {
		t.Error("Mappings expected to have 1 element")
	}

	if len(mappings.Sources()) != 1 || mappings.Sources()[0] != "application.yaml" {
		t.Error("Sources expected to have exactly 1 element 'application.yaml'")
	}

	if len(mappings.Destinations()) != 1 || mappings.Destinations()[0] != "config.yaml" {
		t.Error("Destinations expected to have exactly 1 element 'config.yaml'")
	}
}
