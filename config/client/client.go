package client

import (
	"fmt"
	"github.com/go-resty/resty"
	"strings"
)

//Extension format of downloaded config
type Extension string

//ParseExtension parse string into Extension type
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
		return unknown, fmt.Errorf("failed to parse extension: '%s'", str)
	}
}

const (
	configPathFmt     = "/%s/%s-%s.%s"
	configFilePathFmt = "/%s/%s/%s/%s"
	encryptPath       = "/encrypt"
	decryptPath       = "/decrypt"
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

	//Encrypt encrypts the value server side and returns result
	Encrypt(value string) (string, error)

	//Decrypt decrypts the value server side and returns result
	Decrypt(value string) (string, error)
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
func NewClient(c Config) Client {
	client := &client{
		config: &c,
	}

	resty.SetHostURL(c.URI)
	resty.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		if response.StatusCode() >= 300 || response.StatusCode() < 200 {
			return fmt.Errorf("unexpected response code '%s'", response.Status())
		}
		return nil
	})

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
	resp, err := resty.R().Get(c.formatValuesURI(extension))
	return resp.String(), err
}

//Encrypt encrypts the value server side and returns result
func (c *client) Encrypt(value string) (string, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "text/plain").
		SetBody(value).
		Post(encryptPath)
	return resp.String(), err
}

//Decrypt decrypts the value server side and returns result
func (c *client) Decrypt(value string) (string, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "text/plain").
		SetBody(value).
		Post(decryptPath)
	return resp.String(), err
}

func (c *client) formatValuesURI(extension Extension) string {
	return fmt.Sprintf(configPathFmt, c.config.Label, c.config.Application, c.config.Profile, extension)
}

func (c *client) formatFileURI(source string) string {
	return fmt.Sprintf(configFilePathFmt, c.config.Application, c.config.Profile, c.config.Label, source)
}
