package cmd

import (
	"fmt"
	"os"

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
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Turn on debug logging")

	rootCmd.MarkPersistentFlagRequired("vault-name")

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("vault-name", rootCmd.PersistentFlags().Lookup("vault-name"))
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

	viper.AutomaticEnv() // read in environment variables that match
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
}
