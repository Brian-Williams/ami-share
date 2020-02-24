/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/Brian-Williams/ami-share/common"
	"github.com/Brian-Williams/ami-share/core"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	cliName    = "ami-share"
	cliExample = "AWS_SDK_LOAD_CONFIG=true AWS_PROFILE=staging-ami ./ami-share -v -c example.yaml -p plan.yaml"
)

var (
	cfgFile      string
	amiCfgFile   string
	buildVersion = "unknown"
	buildHash    = "unknown"
	buildDate    = "unknown"
	buildInfo    = fmt.Sprintf("Version=%s, Build=%s, Date=%s", buildVersion, buildHash, buildDate)
	params       common.ShareParams
	verbose      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   cliName,
	Short: "AWS ami share utility",
	Long: fmt.Sprintf(`AWS AMI Share is a utility for sharing AMIs across accounts.
%s`, buildInfo),
	Example: cliExample,
	Version: buildVersion,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.WithFields(log.Fields{
			"context":   "share-command",
			"operation": "validation",
		})

		if config, err := common.LoadConfig(amiCfgFile); err != nil {
			logger.Errorf("failed to parse config file: %v", err)
			return err
		} else {
			logger.Info("Validating config")
			if err := config.Validate(); err != nil {
				return err
			}
			params.Config = config
		}

		logger.Info("Initializing")
		shareAMI, err := core.NewAWSShareAMI(&params)
		if err != nil {
			return err
		}

		logger.Info("Validating accounts")
		if err := shareAMI.ValidateAccounts(); err != nil {
			return err
		}
		return shareAMI.Run()
	},
	PreRun: logLevel,
}

func logLevel(_ *cobra.Command, args []string) {
	log.SetLevel(log.InfoLevel)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ami-share.yaml)")
	rootCmd.PersistentFlags().StringVarP(&amiCfgFile, "ami-config", "c", "",
		"Path to the config file")
	rootCmd.PersistentFlags().StringVarP(&params.PlanFile, "plan", "p", "",
		"(required) Path to output file for plan.")

	// Non-required flags aren't persistent for shim usage
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enables debug output.")
	rootCmd.Flags().BoolVar(&params.NoDryRun, "no-dry-run", false,
		"If specified, it shares AMIs. Otherwise it just list target candidates in plan file.")
	rootCmd.Flags().BoolVar(&params.ShareSnapshots, "share-snapshots", false,
		"Whether to share snapshots attached to AMIs.")

	for _, required := range []string{"ami-config", "plan"} {
		rootCmd.MarkPersistentFlagRequired(required)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ami-share" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ami-share")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
