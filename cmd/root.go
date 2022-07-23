/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/lesnuages/thm/pkg/vmware"
	"github.com/spf13/cobra"
)

var config *vmware.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "thm",
	Short: "The Manager",
	Long:  `A simple CLI tool to manage a VSphere cluster.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	configPath := path.Join(os.Getenv("HOME"), ".config", "thm")
	config, err = vmware.LoadConfig(configPath)
	if err != nil {
		err = os.MkdirAll(configPath, 0755)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(path.Join(configPath, "config.env"), []byte(vmware.GetDefaultConfig()), 0644)
		if err != nil {
			panic(err)
		}
		config, err = vmware.LoadConfig(configPath)
		if err != nil || config == nil {
			panic(err)
		}
	}
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.thm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
