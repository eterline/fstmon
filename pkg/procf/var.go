package procf

import (
	"fmt"
	"os"
)

type ProcFile string

func (pf ProcFile) Data() (data []byte, err error) {
	data, err = os.ReadFile(string(pf))
	if err != nil {
		return nil, fmt.Errorf("failed to read proc file '%s': %v", pf, err)
	}
	return data, nil
}

const (
	procNetDev  ProcFile = "/proc/net/dev"  // network devices
	procNetSnmp ProcFile = "/proc/net/snmp" // network snmp counters

	procSelfStatus ProcFile = "/proc/self/status" // information of current go program

	procCpuInfo   ProcFile = "/proc/cpuinfo"   // cpu info
	procMemInfo   ProcFile = "/proc/meminfo"   //
	procVmStat    ProcFile = "/proc/vmstat"    // information of current go program
	procDiskStats ProcFile = "/proc/diskstats" //
	procCrypto    ProcFile = "/proc/crypto"    //
	procLoadAvg   ProcFile = "/proc/loadavg"   //
	procUptime    ProcFile = "/proc/uptime"    //
)

const (
	psCmdFieldsArg = "pid:10,ppid:10,pcpu:5,pmem:5,user:15,comm:20,class:6,stat:5,args"
)

var (
	dockerPrefixes = []string{
		"/var/lib/docker-volumes",
		"/var/lib/docker",
		"/var/run/docker",
	}

	mntPrefixes = []string{
		"/mnt",
		"/volumes",
	}

	loopPrefixes = []string{
		"/dev/loop",
	}

	runPrefixes = []string{
		"/run",
	}
)
