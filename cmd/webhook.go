package cmd

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wandera/scccmd/pkg/inject"
)

var wp = struct {
	configFile string
	port       int
	certFile   string
	keyFile    string
}{}

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Runs K8s webhook for injecting config from Cloud Config Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeWebhook()
	},
}

func executeWebhook() error {
	wh, err := inject.NewWebhook(inject.WebhookParameters{
		Port:       wp.port,
		ConfigFile: wp.configFile,
		CertFile:   wp.certFile,
		KeyFile:    wp.keyFile,
	})
	if err != nil {
		return err
	}

	stop := make(chan struct{})
	go wh.Run(stop)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Config webhook started")
	<-signalChan
	log.Info("Shutdown signal received, exiting...")
	close(stop)

	return nil
}

func init() {
	webhookCmd.Flags().StringVarP(&wp.configFile, "config-file", "f", "config/config.yaml", "the configuration file")
	webhookCmd.Flags().StringVarP(&wp.certFile, "cert-file", "c", "keys/publickey.cer", "location of public part of SSL certificate")
	webhookCmd.Flags().StringVarP(&wp.keyFile, "key-file", "k", "keys/private.key", "location of private key of SSL certificate")
	webhookCmd.Flags().IntVarP(&wp.port, "port", "p", 443, "webhook port")
}
