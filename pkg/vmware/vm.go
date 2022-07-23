package vmware

import (
	"context"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
)

var powerStates = map[string]string{
	"poweredOff": "off",
	"poweredOn":  "on",
	"suspended":  "suspended",
}

type VirtualMachine struct {
	Name     string
	vm       mo.VirtualMachine
	OS       string
	IP       string
	CPU      int32
	Mem      int32
	Hostname string
}

func (v *VirtualMachine) IsTemplate() bool {
	return v.vm.Summary.Config.Template
}

func (v *VirtualMachine) GetPowerState() string {
	return powerStates[string(v.vm.Summary.Runtime.PowerState)]
}

func ListVMs(conf *Config) []*VirtualMachine {
	var out []*VirtualMachine
	Run(conf, func(ctx context.Context, c *vim25.Client) error {
		// Create view of VirtualMachine objects
		m := view.NewManager(c)

		v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
		if err != nil {
			return err
		}

		defer v.Destroy(ctx)

		// Retrieve summary property for all machines
		// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
		var vms []mo.VirtualMachine
		err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"guest", "summary"}, &vms)
		if err != nil {
			return err
		}

		// Print summary per vm (see also: govc/vm/info.go)
		for _, vm := range vms {
			hostname := ""
			if vm.Guest != nil {
				hostname = vm.Guest.HostName
			}
			out = append(out, &VirtualMachine{
				vm:       vm,
				Name:     vm.Summary.Config.Name,
				OS:       vm.Guest.GuestFullName,
				IP:       vm.Summary.Guest.IpAddress,
				CPU:      vm.Summary.Config.NumCpu,
				Mem:      vm.Summary.Config.MemorySizeMB,
				Hostname: hostname,
			})
		}
		return nil
	})
	return out
}

// GetVM - Retrieve a VirtualMachine from the given name
func GetVM(conf *Config, name string) (*VirtualMachine, error) {
	var virtMachine *VirtualMachine
	err := Run(conf, func(ctx context.Context, c *vim25.Client) error {
		// Create view of VirtualMachine objects
		m := view.NewManager(c)

		v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
		if err != nil {
			return err
		}

		defer v.Destroy(ctx)

		// Retrieve summary property for all machines
		// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
		var vms []mo.VirtualMachine
		err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"guest", "summary"}, &vms)
		if err != nil {
			return err
		}

		// Print summary per vm (see also: govc/vm/info.go)
		for _, vm := range vms {
			if vm.Summary.Config.Name == name {
				virtMachine = &VirtualMachine{
					vm:       vm,
					Name:     vm.Summary.Config.Name,
					OS:       vm.Guest.GuestFullName,
					IP:       vm.Summary.Guest.IpAddress,
					CPU:      vm.Summary.Config.NumCpu,
					Mem:      vm.Summary.Config.MemorySizeMB,
					Hostname: vm.Guest.HostName,
				}
				return nil
			}
		}
		return nil
	})
	return virtMachine, err
}

func PowerChange(conf *Config, v *VirtualMachine, state string) error {
	return Run(conf, func(ctx context.Context, c *vim25.Client) error {
		finder := find.NewFinder(c)
		dc, err := finder.Datacenter(ctx, conf.GovcDataCenter)
		if err != nil {
			return err
		}
		finder.SetDatacenter(dc)
		vm, err := finder.VirtualMachine(ctx, v.Name)
		if err != nil {
			return err
		}
		var task *object.Task
		switch state {
		case "on":
			task, err = vm.PowerOn(ctx)
		case "off":
			task, err = vm.PowerOff(ctx)
		case "suspended":
			task, err = vm.Suspend(ctx)
		}
		if err != nil {
			return err
		}
		err = task.Wait(ctx)
		return err
	})
}
