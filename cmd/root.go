/*
Copyright © 2022 SUSE LLC

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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elemental",
		Short: "elemental",
	}
	cmd.PersistentFlags().Bool("debug", false, "enable debug output")
	cmd.PersistentFlags().String("config-dir", "/etc/elemental", "set config dir (default is /etc/elemental)")
	cmd.PersistentFlags().String("logfile", "", "set logfile")
	cmd.PersistentFlags().Bool("quiet", false, "do not output to stdout")
	_ = viper.BindPFlag("debug", cmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("config-dir", cmd.PersistentFlags().Lookup("config-dir"))
	_ = viper.BindPFlag("logfile", cmd.PersistentFlags().Lookup("logfile"))
	_ = viper.BindPFlag("quiet", cmd.PersistentFlags().Lookup("quiet"))
	return cmd
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = NewRootCmd()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
