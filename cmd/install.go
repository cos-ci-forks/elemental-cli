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
	"errors"
	"os/exec"

	"github.com/rancher/elemental-cli/cmd/config"
	"github.com/rancher/elemental-cli/pkg/action"
	v1 "github.com/rancher/elemental-cli/pkg/types/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/mount-utils"
)

// NewInstallCmd returns a new instance of the install subcommand and appends it to
// the root command. requireRoot is to initiate it with or without the CheckRoot
// pre-run check. This method is mostly used for testing purposes.
func NewInstallCmd(root *cobra.Command, addCheckRoot bool) *cobra.Command {
	c := &cobra.Command{
		Use:   "install DEVICE",
		Short: "elemental installer",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if addCheckRoot {
				return CheckRoot()
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := exec.LookPath("mount")
			if err != nil {
				return err
			}
			mounter := mount.New(path)

			cfg, err := config.ReadConfigRun(viper.GetString("config-dir"), cmd.Flags(), mounter)
			if err != nil {
				cfg.Logger.Errorf("Error reading config: %s\n", err)
			}

			if err := validateInstallUpgradeFlags(cfg.Logger, cmd.Flags()); err != nil {
				return err
			}

			// Manage deprecated flags
			// Adapt 'docker-image' and 'directory'  deprecated flags to 'system' syntax
			adaptDockerImageAndDirectoryFlagsToSystem(cmd.Flags())

			//Adapt 'force-efi' and 'force-gpt' to 'firmware' and 'part-table'
			adaptEFIAndGPTFlags(cmd.Flags())

			cmd.SilenceUsage = true
			spec, err := config.ReadInstallSpec(cfg, cmd.Flags())
			if err != nil {
				cfg.Logger.Errorf("invalid install command setup %v", err)
				return err
			}

			if len(args) == 1 {
				spec.Target = args[0]
			}

			if spec.Target == "" {
				return errors.New("at least a target device must be supplied")
			}

			cfg.Logger.Infof("Install called")
			install := action.NewInstallAction(cfg, spec)
			err = install.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}
	firmType := newEnumFlag([]string{v1.EFI, v1.BIOS}, v1.EFI)
	pTableType := newEnumFlag([]string{v1.GPT, v1.MSDOS}, v1.GPT)

	root.AddCommand(c)
	c.Flags().StringP("cloud-init", "c", "", "Cloud-init config file")
	c.Flags().StringP("iso", "i", "", "Performs an installation from the ISO url")
	c.Flags().StringP("partition-layout", "p", "", "Partitioning layout file")
	_ = c.Flags().MarkDeprecated("partition-layout", "'partition-layout' is deprecated and ignored please use a config file instead")
	c.Flags().Bool("no-format", false, "Don’t format disks. It is implied that COS_STATE, COS_RECOVERY, COS_PERSISTENT, COS_OEM are already existing")

	c.Flags().Bool("force-efi", false, "Forces an EFI installation")
	_ = c.Flags().MarkDeprecated("force-efi", "'force-efi' is deprecated please use 'firmware' instead")
	c.Flags().Var(firmType, "firmware", "Firmware to install for ('esp' or 'bios')")

	c.Flags().Bool("force-gpt", false, "Forces a GPT partition table")
	_ = c.Flags().MarkDeprecated("force-gpt", "'force-gpt' is deprecated please use 'part-table' instead")
	c.Flags().Var(pTableType, "part-table", "Partition table type to use")

	c.Flags().String("tty", "", "Add named tty to grub")
	c.Flags().Bool("force", false, "Force install")
	c.Flags().Bool("eject-cd", false, "Try to eject the cd on reboot, only valid if booting from iso")
	addSharedInstallUpgradeFlags(c)
	addLocalImageFlag(c)
	return c
}

// register the subcommand into rootCmd
var _ = NewInstallCmd(rootCmd, true)
