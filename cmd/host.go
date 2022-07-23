/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/lesnuages/thm/pkg/utils"
	"github.com/lesnuages/thm/pkg/vmware"
	"github.com/spf13/cobra"
)

// datacenterCmd represents the datacenter command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Manage an ESXi host",
	Run: func(cmd *cobra.Command, args []string) {
		host(cmd, args)
	},
}

func host(cmd *cobra.Command, args []string) {
	hostsInfo, err := vmware.GethostInfos(config)
	if err != nil {
		utils.PrintError("%v", err)
		return
	}
	for _, h := range hostsInfo {
		cmd.Println()
		printHostInfo(cmd, h)
	}

}

func printHostInfo(cmd *cobra.Command, h *vmware.HostSystem) {
	hostBuf := bytes.NewBufferString("")
	hostTable := tabwriter.NewWriter(hostBuf, 0, 2, 2, ' ', 0)
	fmt.Fprintf(hostTable, "%s\t%s\t\n", utils.Bold("Hostname"), h.Name)
	fmt.Fprintf(hostTable, "%s\t%.2f MHz\t\n", utils.Bold("CPU Usage"), h.UsedCPU)
	fmt.Fprintf(hostTable, "%s\t%.2f GHz\t\n", utils.Bold("Total CPU"), h.TotalCPU)
	fmt.Fprintf(hostTable, "%s\t%s\t\n", utils.Bold("Memory Usage"), h.MemUsage)
	fmt.Fprintf(hostTable, "%s\t%s\t\n", utils.Bold("Free Memory"), h.FreeMemory)
	hostTable.Flush()
	cmd.Println(hostBuf.String())
}

func init() {
	rootCmd.AddCommand(hostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hostCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
