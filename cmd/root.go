package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/apex/log"
	"github.com/spf13/viper"
)

var config string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "akv-entrypoint",
	Short: "Docker entrypoint that gets secrets from Azure KeyVault",
	Long: `This program fetches secrets from Azure KeyVault,
and exposes them as environment variable for the application.

Optionally, it can transform the env vars.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "TODO: Path to the config file (optional). Allows to set transformations")
	rootCmd.PersistentFlags().StringP("vault-name", "n", "", "KeyVault name")
	rootCmd.PersistentFlags().StringSliceP("secret-name", "s", []string{}, "Secret names to pull. If a name starts with json:, i.e. json:secret1, they are treated as JSON objects. If not set akv-entrypoint will List all secrets and read them all. Can be specified multiple times. If set via env var AKV_ENTRYPOINT_SECRET_NAMES, multiple values should be separated by a comma.")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Turn on debug logging")

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("vault-name", rootCmd.PersistentFlags().Lookup("vault-name"))
	viper.BindPFlag("secret-names", rootCmd.PersistentFlags().Lookup("secret-name"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if config != "" {
		viper.SetConfigFile(config)

		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			log.WithError(err).Fatal("Had some errors while parsing the config")
		}
	}

	viper.SetEnvPrefix("AKV_ENTRYPOINT")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	if viper.GetString("vault-name") == "" {
		fmt.Println("Please set the vault name")
		os.Exit(1)
	}
}
