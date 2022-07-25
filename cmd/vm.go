/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lesnuages/thm/pkg/utils"
	"github.com/lesnuages/thm/pkg/vmware"
	"github.com/spf13/cobra"
)

// vmCmd represents the vm command
var (
	vmName string
	vmCmd  = &cobra.Command{
		Use:   "vm",
		Short: "Interact with virtual machines",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	lsCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List virtual machines",
		Run: func(cmd *cobra.Command, args []string) {
			listVMs(cmd, args)
		},
	}

	detailsCmd = &cobra.Command{
		Use:     "details",
		Aliases: []string{"detail", "info"},
		Short:   "Show details of a virtual machine",
		Run: func(cmd *cobra.Command, args []string) {
			detailsVM(cmd, args)
		},
	}

	powerCmd = &cobra.Command{
		Use:       "power [on|off|suspend]",
		Aliases:   []string{"pwr"},
		Short:     "Manage virtual machine states (on, off, suspend)",
		ValidArgs: []string{"on", "off", "suspend"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}
			state := args[0]
			switch state {
			case "on":
				state = "on"
			case "off":
				state = "off"
			case "suspend":
				state = "suspended"
			default:
				utils.PrintError("Invalid state: %s\n", utils.Bold(state))
				return
			}
			powerVM(cmd, state)
		},
	}
)

func listVMs(cmd *cobra.Command, args []string) {
	vms := vmware.ListVMs(config)
	if vms == nil {
		return
	}
	tw := table.NewWriter()
	tw.SetStyle(table.StyleLight)
	tw.SortBy([]table.SortBy{
		{Name: "Power State", Mode: table.Dsc},
		{Name: "Name", Mode: table.Asc},
	})
	tw.AppendHeader(table.Row{"Name", "IP Address", "Hostname", "Operating System", "CPU", "Mem", "Power State"})
	rows := []table.Row{}
	for _, vm := range vms {
		powerState := vm.GetPowerState()
		switch powerState {
		case "on":
			powerState = utils.MsgInfo(powerState)
		case "off":
			powerState = utils.MsgError(powerState)
		case "suspended":
			powerState = utils.MsgWarning(powerState)
		}
		if !vm.IsTemplate() {
			rows = append(rows, table.Row{vm.Name, vm.IP, vm.Hostname, vm.OS, fmt.Sprintf("%d ", vm.CPU), fmt.Sprintf("%d MB", vm.Mem), powerState})
		}
	}
	tw.AppendRows(rows)
	cmd.Printf("%s\n", tw.Render())
}

func detailsVM(cmd *cobra.Command, args []string) {
	var (
		vm  *vmware.VirtualMachine
		err error
	)
	if len(args) == 0 && vmName == "" {
		vm, err = selectVM()
	} else {
		if len(args) != 0 {
			vmName = args[0]
		}
		vm, err = vmware.GetVM(config, vmName)
	}
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	if vm == nil {
		return
	}
	printVM(cmd, vm)
}

func printVM(cmd *cobra.Command, vm *vmware.VirtualMachine) {
	vmBuf := bytes.NewBufferString("")
	vmTable := tabwriter.NewWriter(vmBuf, 0, 2, 2, ' ', 0)
	cmd.Println()
	fmt.Fprintf(vmTable, "%s\t%s\t\n", utils.Bold("Name"), vm.Name)
	fmt.Fprintf(vmTable, "%s\t%s\t\n", utils.Bold("IP Address"), vm.IP)
	fmt.Fprintf(vmTable, "%s\t%s\t\n", utils.Bold("Hostname"), vm.Hostname)
	fmt.Fprintf(vmTable, "%s\t%s\t\n", utils.Bold("Operating System"), vm.OS)
	fmt.Fprintf(vmTable, "%s\t%d\t\n", utils.Bold("CPU"), vm.CPU)
	fmt.Fprintf(vmTable, "%s\t%d MB\t\n", utils.Bold("Memory"), vm.Mem)
	fmt.Fprintf(vmTable, "%s\t%s\t\n", utils.Bold("Power State"), vm.GetPowerState())
	vmTable.Flush()
	cmd.Println(vmBuf.String())
}

func powerVM(cmd *cobra.Command, state string) {
	var (
		virtMachine *vmware.VirtualMachine
		err         error
	)
	if vmName == "" {
		virtMachine, err = selectVM()
	} else {
		virtMachine, err = vmware.GetVM(config, vmName)
	}
	if err != nil {
		utils.PrintError("%v", err)
		return
	}
	if virtMachine.GetPowerState() == state {
		utils.PrintError("%s is already %s\n", virtMachine.Name, state)
		return
	}
	err = vmware.PowerChange(config, virtMachine, state)
	if err != nil {
		utils.PrintError("Error powering VM: %s", err)
	}
	utils.PrintInfo("Changed %s state to %s\n", virtMachine.Name, state)
	printVM(cmd, virtMachine)
}

func selectVM() (*vmware.VirtualMachine, error) {
	vms := vmware.ListVMs(config)
	if vms == nil {
		return nil, fmt.Errorf("no virtual machines found")
	}
	buf := bytes.NewBufferString("")
	table := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	displayedVMs := []*vmware.VirtualMachine{}
	for _, v := range vms {
		displayedVMs = append(displayedVMs, v)
		fmt.Fprintf(table, "%s\t%s\t\n", v.Name, v.GetPowerState())
	}
	table.Flush()
	options := strings.Split(buf.String(), "\n")
	options = options[:len(options)-1]
	prompt := survey.Select{
		Message: "Select a virtual machine:",
		Options: options,
	}
	var selected string
	err := survey.AskOne(&prompt, &selected)
	if err != nil {
		return nil, err
	}
	for index, value := range options {
		if value == selected {
			return displayedVMs[index], nil
		}
	}
	return nil, fmt.Errorf("invalid selection")
}

func init() {
	vmCmd.Flags().StringVarP(&vmName, "name", "n", "", "Name of the virtual machine")
	vmCmd.AddCommand(detailsCmd)
	vmCmd.AddCommand(lsCmd)
	vmCmd.AddCommand(powerCmd)
	rootCmd.AddCommand(vmCmd)
}
