package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dedelala/sysexits"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "container-builder",
	Short: "A tool that builds opinionated containers",
	Long: `A tool that builds opinionated containers.
  
More information can be found as to what constitutes a "well 
constructed" container at the following URL:

  https://l.littleman.co/2AY1PX3
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(sysexits.Usage)
	},
}

// Execute starts the initial root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(sysexits.Usage)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLog)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/.container-builder.yaml)")
	rootCmd.PersistentFlags().StringP("log-level", "v", "warn", "log level for the application (default is warn)")

	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".container-builder" (without extension).
		viper.AddConfigPath("/etc")
		viper.SetConfigName("container-builder")
	}

	viper.SetEnvPrefix("CONTAINER_BUILDER")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Unable to read in configuration: %s\n", err.Error())
	}
}

func initLog() {
	// Set up logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	level, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		log.SetLevel(log.WarnLevel)
		log.WithFields(log.Fields{
			"log-level": viper.GetString("log-level"),
		}).Warn("Unable to parse log level. Defaulting to warn")
	} else {
		log.SetLevel(level)
	}
}
