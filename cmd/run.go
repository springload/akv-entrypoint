package cmd

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	lib "github.com/springload/akv-entrypoint/lib"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the specified program with the env vars fetched from the keyvault",
	Run: func(cmd *cobra.Command, args []string) {

		vaultName := viper.GetString("vault-name")
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

		for key, value := range secrets {
			os.Setenv(key, value)
		}

		// sanity check
		command, err := exec.LookPath(args[0])
		ctx := log.WithFields(log.Fields{"command": command})
		if err != nil {
			ctx.WithError(err).Fatal("Cant find the command")
		}
		if err := syscall.Exec(command, args, os.Environ()); err != nil {
			ctx.WithError(err).Fatal("Can't run the command")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
