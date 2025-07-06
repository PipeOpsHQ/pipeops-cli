/*
Copyright © 2024 9trocode

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// Version is set at build time
var Version = "dev"

type Config struct {
	Version VersionInfo
}

type VersionInfo struct {
	Version string
}

var Conf Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pipeops",
	Short:   "🚀 PipeOps CLI - Manage your cloud-native environment",
	Long:    `🚀 PipeOps CLI is a command-line interface for managing your cloud-native environment and deployments.`,
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set global JSON output flag
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			// Set a global flag that other commands can check
			cmd.Root().SetContext(context.WithValue(cmd.Root().Context(), "json", true))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("version").Changed {
			fmt.Println("🚀 PipeOps CLI Version:", Version)
			return
		}

		// Show help by default
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	SaveConfig()
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().Bool("quiet", false, "Suppress non-essential output")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pipeops.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "Prints out the current version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pipeops-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pipeops")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func GetConfig() Config {
	var filename string

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filename = fmt.Sprintf("%s/%s", home, ".pipeops.json")

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("Config file does not exist")
		os.Exit(1)
	}

	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(dataBytes, &Conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return Conf
}

func SaveConfig() {
	var filename string

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filename = fmt.Sprintf("%s/%s", home, ".pipeops.json")

	Conf.Version.Version = Version

	dataBytes, err := json.Marshal(Conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.WriteFile(filename, dataBytes, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.Chmod(filename, 0600); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
