/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/lesnuages/thm/pkg/utils"
	"github.com/lesnuages/thm/pkg/vmware"
	"github.com/spf13/cobra"
)

var config *vmware.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "thm",
	Short: "The Manager",
	Long:  `A simple CLI tool to manage a VSphere cluster.`,
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
	var err error

	configPath := path.Join(os.Getenv("HOME"), ".config", "thm")
	config, err = vmware.LoadConfig(configPath)
	if err != nil {
		err = os.MkdirAll(configPath, 0755)
		if err != nil {
			utils.PrintError("%v", err)
			return
		}
		err = ioutil.WriteFile(path.Join(configPath, "config.env"), []byte(vmware.GetDefaultConfig()), 0644)
		if err != nil {
			utils.PrintError("%v", err)
			return
		}
		config, err = vmware.LoadConfig(configPath)
		if err != nil || config == nil {
			utils.PrintError("%v", err)
			return
		}
	}
}
