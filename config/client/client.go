package client

import (
	"fmt"
	"github.com/go-resty/resty"
	"errors"
	"strings"
)

type Extension string

func ParseExtension(str string) (Extension, error) {
	switch value := strings.TrimRight(str, "/n"); value {
	case "json":
		return json, nil
	case "properties":
		return properties, nil
	case "yaml":
		return yaml, nil
	case "yml":
		return yaml, nil
	default:
		return unknown, errors.New(fmt.Sprintf("failed to parse extension: '%s'", str))
	}
}

const (
	configPathFmt     = "%s/%s/%s-%s.%s"
	configFilePathFmt = "%s/%s/%s/%s/%s"
)

const (
	json       Extension = "json"
	properties Extension = "properties"
	yaml       Extension = "yml"
	unknown    Extension = "_"
)

//Client Spring Cloud Config Client
type Client interface {
	//Config of the client
	Config() *Config

	//FetchFile queries the remote configuration service and returns the resulting file
	FetchFile(source string) ([]byte, error)

	//FetchAs queries the remote configuration service and returns the result in specified format
	FetchAs(extension Extension) (string, error)

	//FetchAsJSON queries the remote configuration service and returns the result as a JSON string
	FetchAsJSON() (string, error)

	//FetchAsYAML queries the remote configuration service and returns the result as a YAML string
	FetchAsYAML() (string, error)

	//FetchAsProperties queries the remote configuration service and returns the result as a Properties string
	FetchAsProperties() (string, error)
}

//Config needed to fetch a remote configuration
type Config struct {
	URI         string
	Profile     string
	Application string
	Label       string
}

type client struct {
	config *Config
}

//NewClient creates instance of the Client
func NewClient(c Config) (Client) {
	client := &client{
		config: &c,
	}
	return client
}

//Config of the client
func (c *client) Config() *Config {
	return c.config
}

//FetchFile queries the remote configuration service and returns the resulting file
func (c *client) FetchFile(source string) ([]byte, error) {
	resp, err := resty.R().Get(c.formatFileURI(source))

	return resp.Body(), err
}

//FetchAsProperties queries the remote configuration service and returns the result as a Properties string
func (c *client) FetchAsProperties() (string, error) {
	return c.FetchAs(properties)
}

//FetchAsJSON queries the remote configuration service and returns the result as a JSON string
func (c *client) FetchAsJSON() (string, error) {
	return c.FetchAs(json)
}

//FetchAsYAML queries the remote configuration service and returns the result as a YAML string
func (c *client) FetchAsYAML() (string, error) {
	return c.FetchAs(yaml)
}

//FetchAs queries the remote configuration service and returns the result in specified format
func (c *client) FetchAs(extension Extension) (string, error) {
	resp, err := resty.R().Get(c.formatURI(extension))
	return resp.String(), err
}

func (c *client) formatURI(extension Extension) string {
	return fmt.Sprintf(configPathFmt, c.config.URI, c.config.Label, c.config.Application, c.config.Profile, extension)
}

func (c *client) formatFileURI(source string) string {
	return fmt.Sprintf(configFilePathFmt, c.config.URI, c.config.Application, c.config.Profile, c.config.Label, source)
}
