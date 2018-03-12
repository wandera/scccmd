package client

import (
	"fmt"
	"github.com/go-resty/resty"
)

const (
	configPathFmt     = "%s/%s/%s-%s.%s"
	configFilePathFmt = "%s/%s/%s/%s/%s"
	json              = "json"
	properties        = "properties"
	yaml              = "yml"
)

//Client Spring Cloud Config Client
type Client interface {
	//Config of the client
	Config() *Config

	//FetchFile queries the remote configuration service and returns the resulting file
	FetchFile(source string) ([]byte, error)

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

//FetchFile
func (c *client) FetchFile(source string) ([]byte, error) {
	resp, err := resty.R().Get(c.formatFileURI(source))

	return resp.Body(), err
}

//FetchAsProperties
func (c *client) FetchAsProperties() (string, error) {
	return c.fetchAsString(properties)
}

//FetchAsJSON
func (c *client) FetchAsJSON() (string, error) {
	return c.fetchAsString(json)
}

//FetchAsYAML
func (c *client) FetchAsYAML() (string, error) {
	return c.fetchAsString(yaml)
}

func (c *client) fetchAsString(extension string) (string, error) {
	resp, err := resty.R().Get(c.formatURI(extension))
	return resp.String(), err
}

func (c *client) formatURI(extension string) string {
	return fmt.Sprintf(configPathFmt, c.config.URI, c.config.Label, c.config.Application, c.config.Profile, extension)
}

func (c *client) formatFileURI(source string) string {
	return fmt.Sprintf(configFilePathFmt, c.config.URI, c.config.Application, c.config.Profile, c.config.Label, source)
}
