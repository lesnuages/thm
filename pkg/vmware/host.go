package vmware

import (
	"context"

	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
)

type HostSystem struct {
	Name       string
	UsedCPU    float64
	TotalCPU   float64
	FreeCPU    float64
	MemUsage   units.ByteSize
	FreeMemory units.ByteSize
	Storage    int32
	hs         mo.HostSystem
}

func GethostInfos(conf *Config) ([]*HostSystem, error) {
	var out []*HostSystem
	err := Run(conf, func(ctx context.Context, c *vim25.Client) error {
		m := view.NewManager(c)

		v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
		if err != nil {
			return err
		}

		defer v.Destroy(ctx)

		// Retrieve summary property for all hosts
		// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.HostSystem.html
		var hss []mo.HostSystem
		err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
		if err != nil {
			return err
		}

		for _, hs := range hss {
			totalCPU := int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores)
			freeCPU := int64(totalCPU) - int64(hs.Summary.QuickStats.OverallCpuUsage)
			freeMemory := int64(hs.Summary.Hardware.MemorySize) - (int64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
			out = append(out, &HostSystem{
				Name:       hs.Summary.Config.Name,
				UsedCPU:    float64(hs.Summary.QuickStats.OverallCpuUsage),
				TotalCPU:   float64(totalCPU / 1024),
				FreeCPU:    float64(freeCPU / 1024),
				MemUsage:   units.ByteSize(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024,
				FreeMemory: units.ByteSize(freeMemory),
				hs:         hs,
			})
		}

		return nil
	})
	return out, err
}
