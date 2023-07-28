package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	lib "github.com/springload/akv-entrypoint/lib"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the env vars fetched from the keyvault",
	Run: func(cmd *cobra.Command, args []string) {

		vaultName := viper.GetString("vault-name")
		if vaultName == "" {
			fmt.Println("Please set the vault name")
			os.Exit(1)
		}

		var secretNames []string
		if err := viper.UnmarshalKey("secret-names", &secretNames); err != nil {
			log.WithError(err).Fatal("can't unmarshal secret names")
		}
		if len(secretNames) > 0 {
			log.Debugf("getting secrets \"%s\" for %s", strings.Join(secretNames, ", "), vaultName)
		} else {
			log.Debugf("getting all secrets for %s", vaultName)
		}
		secrets, err := lib.GetSecrets(vaultName, secretNames...)
		if err != nil {
			log.WithError(err).Fatal("can't get secrets")
		}

		marshalled, err := json.MarshalIndent(secrets, "", "  ")
		if err != nil {
			log.WithError(err).Fatal("Can't marshal json")
		}
		fmt.Println(string(marshalled))
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
