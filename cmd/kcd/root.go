/*
Copyright © 2019 Zedge, Inc.

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
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/zedge/kubecd/pkg/model"
	"github.com/zedge/kubecd/pkg/updates"
)

var cfgFile string

var environmentsFile string
var verbosity int

type NotYetImplementedError string

func (e NotYetImplementedError) Error() string {
	return "not yet implemented: " + string(e)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kcd",
	Short: "kcd is the command line interface for KubeCD",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubecd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	defaultEnvFile := os.Getenv("KUBECD_ENVIRONMENTS")
	if defaultEnvFile == "" {
		defaultEnvFile = "environments.yaml"
	}
	rootCmd.PersistentFlags().StringVarP(&environmentsFile, "environments-file", "f", defaultEnvFile, `KubeCD config file file (default $KUBECD_ENVIRONMENTS or "environments.yaml")`)
	rootCmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "Increase verbosity level")
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

		// Search config in home directory with name ".kubecd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kubecd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func runCommand(dryRun, disableColors bool, argv []string) error {
	printCmd := strings.Join(argv, " ")

	if !disableColors {
		_, _ = colorstring.Fprintf(os.Stderr, "[yellow]%s\n", printCmd)
	}

	if !dryRun {
		cmd := exec.Command(argv[0], argv[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("command %q failed: %w", printCmd, err)
		}
	}
	return nil
}

func matchAllArgs(checks ...cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		for _, check := range checks {
			if err := check(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}

func environmentsFromArgs(kcdConfig *model.KubeCDConfig, cluster string, args []string) ([]*model.Environment, error) {
	if len(args) > 0 {
		for _, envName := range args {
			env := kcdConfig.GetEnvironment(envName)
			if env == nil {
				return nil, fmt.Errorf(`unknown environment: %q`, envName)
			}
			return []*model.Environment{env}, nil
		}
	}
	if cluster == "" {
		return nil, fmt.Errorf("specify --cluster flag or ENV arg")
	}
	if !kcdConfig.HasCluster(cluster) {
		return nil, fmt.Errorf("unknown cluster: %q", cluster)
	}
	return kcdConfig.GetEnvironmentsInCluster(cluster), nil
}

func clusterFlagOrEnvArg(clusterFlag *string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if *clusterFlag == "" && len(args) != 1 {
			return fmt.Errorf("specify --cluster flag or ENV arg")
		}
		return nil
	}
}

func makeReleaseFilters(args []string, cluster string, releases []string, image string) []updates.ReleaseFilterFunc {
	filters := make([]updates.ReleaseFilterFunc, 0)
	if cluster != "" {
		filters = append(filters, updates.ClusterReleaseFilter(pollCluster))
	} else if len(args) == 1 {
		filters = append(filters, updates.EnvironmentReleaseFilter(args[0]))
	}
	if len(releases) > 0 {
		filters = append(filters, updates.ReleaseFilter(pollReleases))
	}
	if image != "" {
		filters = append(filters, updates.ImageReleaseFilter(pollImage))
	}
	return filters
}
