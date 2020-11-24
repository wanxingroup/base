package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const FlagServiceId = "serviceId"

var rootCmd = &cobra.Command{
	Use:   "launcher",
	Short: "Launch Service",
	Long:  `To help start service quickly`,
}

var cfgFile string

type CommandOption func(command *cobra.Command)

func SetCommandUsage(use string) CommandOption {
	return func(command *cobra.Command) {

		command.Use = use
	}
}

func SetCommandShortDescription(short string) CommandOption {
	return func(command *cobra.Command) {

		command.Short = short
	}
}

func SetCommandLongDescription(long string) CommandOption {
	return func(command *cobra.Command) {

		command.Long = long
	}
}

func AddSubCommand(subCommand *cobra.Command) CommandOption {
	return func(command *cobra.Command) {

		command.AddCommand(subCommand)
	}
}

func InitRootCommand(options ...CommandOption) {

	for _, option := range options {

		option(rootCmd)
	}
}

func GetRootCommand() *cobra.Command {

	return rootCmd
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.service.yaml)")
	rootCmd.PersistentFlags().String(FlagServiceId, "0", "service id")
}

func printError(msg interface{}) {

	fmt.Println("Error:", msg)
	os.Exit(1)
}

// To read configuration, please use viper
// https://github.com/spf13/viper#getting-values-from-viper
// Or Unmarshal to your configuration variable
// https://github.com/spf13/viper#unmarshaling
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			printError(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".service")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
