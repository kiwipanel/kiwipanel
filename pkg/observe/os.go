package observe

import (
	"github.com/shirou/gopsutil/v3/mem"
)

type CurrentRam struct {
	Total     uint64
	Available uint64
	Free      uint64
}

func (cr CurrentRam) CheckRam() bool {
	if cr.Total < 10000000000 {
		return false
	} else {
		return true
	}
}

func CheckCurrentOs() CurrentRam {
	v, _ := mem.VirtualMemory()
	os := CurrentRam{Total: v.Total, Available: v.Available, Free: v.Free}
	return os
}
