package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chijiajian/zstack-webhook/config"
	"github.com/chijiajian/zstack-webhook/handler"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the webhook server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(viper.GetString("config"))
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}

		//log.Printf("Loaded config: %+v\n", cfg)

		outputFormat := viper.GetString("output")

		http.HandleFunc("/webhook", handler.WebhookHandler(cfg, outputFormat))

		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		fmt.Printf("Webhook service is starting, listening on port %d...\n", cfg.Server.Port)
		if cfg.Server.HTTPS {
			fmt.Println("Starting service with HTTPS...")
			log.Fatal(http.ListenAndServeTLS(addr, cfg.Server.CertFile, cfg.Server.KeyFile, nil))
		} else {
			fmt.Println("Starting service with HTTP...")
			log.Fatal(http.ListenAndServe(addr, nil))
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("config", "c", "config.yaml", "Path to the configuration file")
	viper.BindPFlag("config", serveCmd.Flags().Lookup("config"))

	serveCmd.Flags().StringP("output", "o", "text", "Output format: 'text' or 'json'")
	viper.BindPFlag("output", serveCmd.Flags().Lookup("output"))
}
