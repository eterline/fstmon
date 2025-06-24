package procf

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	psdisk "github.com/shirou/gopsutil/v4/disk"
)

/*
PartitionUsage - Represents disk partition usage statistics, including space and inode allocation metrics

	┌──────────────────────┬─────────────────────────────────────────────────────────┐
	│ Field                │ Description                                             │
	├──────────────────────┼─────────────────────────────────────────────────────────┤
	│ Total                │ Total size of the partition (bytes)                     │
	│ Free                 │ Free space available (bytes)                            │
	│ Used                 │ Used space (bytes)                                      │
	│ UsedPercent          │ Percentage of used space (0-100)                        │
	│ InodesTotal          │ Total number of inodes (file/directory metadata slots)  │
	│ InodesUsed           │ Number of inodes currently used                         │
	│ InodesFree           │ Number of free inodes remaining                         │
	│ InodesUsedPercent    │ Percentage of inodes used (0-100)                       │
	└──────────────────────┴─────────────────────────────────────────────────────────┘
*/
type PartitionUsage struct {
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"used_percent"`
	InodesTotal       uint64  `json:"inodes_total"`
	InodesUsed        uint64  `json:"inodes_used"`
	InodesFree        uint64  `json:"inodes_free"`
	InodesUsedPercent float64 `json:"inodes_used_percent"`
}

/*
Partition - Represents a filesystem partition's metadata and usage statistics

	┌──────────────┬─────────────────────────────────────────────────────────────────┐
	│ Field        │ Description                                                     │
	├──────────────┼─────────────────────────────────────────────────────────────────┤
	│ Device       │ Block device path (e.g., /dev/sda1)                             │
	│ MountPoint   │ Directory where the partition is mounted (e.g., /)              │
	│ FsType       │ Filesystem type (e.g., ext4, xfs, ntfs)                         │
	│ Opts         │ Mount options (e.g., ["rw", "noatime", "relatime"])             │
	│ Usage        │ Disk/inode usage stats (PartitionUsage), omitted if not mounted │
	└──────────────┴─────────────────────────────────────────────────────────────────┘
*/
type Partition struct {
	Device     string          `json:"device"`
	MountPoint string          `json:"mount_point"`
	FsType     string          `json:"fs_type"`
	Opts       []string        `json:"opts"`
	Usage      *PartitionUsage `json:"usage,omitempty"`
}

/*
Partitions - types of partitions in system

	┌──────────────┬───────────────────────────────────────────────────────────────┐
	│ Field        │ Description                                                   │
	├──────────────┼───────────────────────────────────────────────────────────────┤
	│ Main         │ Primary physical partitions (e.g., /dev/sda1, /dev/nvme0n1p2) │
	│ LoopBack     │ Loop device partitions (e.g., /dev/loop0 for disk images)     │
	│ Docker       │ Docker-specific mounts (e.g., overlayfs for containers)       │
	│ Mount        │ Manually mounted filesystems (/mnt/***)                       │
	│ Run          │ Ephemeral tmpfs mounts (e.g., /run, /dev/shm)                 │
	└──────────────┴───────────────────────────────────────────────────────────────┘
*/
type Partitions struct {
	Main     []Partition `json:"main_partitions"`
	LoopBack []Partition `json:"loop_partitions"`
	Docker   []Partition `json:"docker_partitions"`
	Mount    []Partition `json:"mount_partitions"`
	Run      []Partition `json:"run_partitions"`
}

/*
SmartAttribute - Represents a single SMART (Self-Monitoring, Analysis and Reporting Technology)
attribute from a storage device (HDD/SSD), as reported by smartctl

	┌───────────────┬───────────────────────────────────────────────────────────────────┐
	│ Field         │ Description                                                       │
	├───────────────┼───────────────────────────────────────────────────────────────────┤
	│ ID            │ SMART attribute ID (e.g., 5 for Reallocated Sectors Count)        │
	│ AttributeName │ Human-readable name (e.g., "Power_On_Hours")                      │
	│ Flag          │ Bitmask flags (vendor-specific, indicates attribute properties)   │
	│ Value         │ Normalized current value (1–253, higher is usually better)        │
	│ Worst         │ Lowest normalized value recorded                                  │
	│ Thresh        │ Threshold value (if Value ≤ Thresh, failure is imminent)          │
	│ Type          │ Attribute type (e.g., "Pre-fail" or "Old_age")                    │
	│ Updated       │ When the value updates (e.g., "Always", "Offline")                │
	│ WhenFailed    │ Failure status (e.g., "FAILING_NOW", "Past", or empty if healthy) │
	│ RawValue      │ Raw data (vendor-specific, e.g., "42" for Power_On_Hours in hex)  │
	└───────────────┴───────────────────────────────────────────────────────────────────┘
*/
type SmartAttribute struct {
	ID            uint16 `json:"id"`
	AttributeName string `json:"attribute_name"`
	Flag          uint16 `json:"flag"`
	Value         uint8  `json:"value"`
	Worst         uint8  `json:"worst"`
	Thresh        uint8  `json:"thresh"`
	Type          string `json:"type"`
	Updated       string `json:"updated"`
	WhenFailed    string `json:"when_failed"`
	RawValue      string `json:"raw_value"`
}

type DiskStatus struct {
	Serial     string                    `json:"serial"`
	MountPoint string                    `json:"mount_point"`
	SMART      map[string]SmartAttribute `json:"smart"`
	Partitions []Partition               `json:"partitions"`
}

func FetchPartitions() (Partitions, error) {

	psList, err := psdisk.Partitions(true)
	if err != nil {
		return Partitions{}, fmt.Errorf("partitions fetch error: %v", err)
	}

	var wg sync.WaitGroup

	data := Partitions{
		Main:     []Partition{},
		LoopBack: []Partition{},
		Docker:   []Partition{},
		Mount:    []Partition{},
		Run:      []Partition{},
	}

	for _, stat := range psList {

		p := readPartitionStat(stat)

		switch {
		case belongsWithPrefix(dockerPrefixes, p.MountPoint):
			data.Docker = append(data.Docker, p)
		case belongsWithPrefix(mntPrefixes, p.MountPoint):
			data.Mount = append(data.Mount, p)
		case belongsWithPrefix(loopPrefixes, p.MountPoint):
			data.LoopBack = append(data.LoopBack, p)
		case belongsWithPrefix(runPrefixes, p.MountPoint):
			data.Run = append(data.Run, p)
		default:
			data.Main = append(data.Main, p)
		}
	}

	wg.Wait()
	return data, nil
}

func ReadDiskStatus(dev string) (DiskStatus, error) {

	smart, err := readSmartctl(dev)
	if err != nil {
		return DiskStatus{}, err
	}

	serial, err := psdisk.SerialNumber(dev)
	if err != nil {
		return DiskStatus{}, err
	}

	parts := devPrefixPartitions(dev)

	status := DiskStatus{
		Serial:     serial,
		MountPoint: dev,
		SMART:      smart,
		Partitions: parts,
	}

	return status, nil
}

func readSmartctl(p string) (smart map[string]SmartAttribute, err error) {

	data, err := exec.Command("smartctl", "-A", p).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'smartctl' command: %v", err)
	}

	var (
		SMART     = map[string]SmartAttribute{}
		dataLines = bytes.Split(data, []byte("\n"))
	)

	for _, smartLine := range dataLines {

		if len(smartLine) < 87 {
			continue
		}

		attrID := bytesToUint32(bytes.TrimSpace(smartLine[:3]))
		if attrID == errUint32 {
			continue
		}

		attrName := strings.TrimSpace(string(smartLine[4:27]))

		SMART[attrName] = SmartAttribute{
			ID:            uint16(attrID),
			AttributeName: attrName,
			Flag:          uint16(bytesToUint32(bytes.TrimSpace(smartLine[31:36]))),
			Value:         uint8(bytesToUint32(bytes.TrimSpace(smartLine[37:42]))),
			Worst:         uint8(bytesToUint32(bytes.TrimSpace(smartLine[43:48]))),
			Thresh:        uint8(bytesToUint32(bytes.TrimSpace(smartLine[49:55]))),
			Type:          strings.TrimSpace(string(smartLine[56:65])),
			Updated:       strings.TrimSpace(string(smartLine[66:74])),
			WhenFailed:    strings.TrimSpace(string(smartLine[75:86])),
			RawValue:      string(smartLine[87:]),
		}
	}

	if len(SMART) == 0 {
		return nil, fmt.Errorf("smartctl did not report smart status")
	}

	return SMART, nil
}

func devPrefixPartitions(p string) []Partition {

	parts := make([]Partition, 0)

	allPartitions, err := psdisk.Partitions(true)
	if err != nil {
		return parts
	}

	for _, stat := range allPartitions {

		if !strings.HasPrefix(stat.Device, p) {
			continue
		}

		part := readPartitionStat(stat)
		parts = append(parts, part)
	}

	return parts
}

/*
/proc/diskstats — Displays I/O statistics for block devices.

Each line contains the following fields:

	┌────┬────────────────────────────────────────────────────────┐
	│ #  │ Description                                            │
	├────┼────────────────────────────────────────────────────────┤
	│  1 │ Major number                                           │
	│  2 │ Minor number                                           │
	│  3 │ Device name (e.g., sda)                                │
	│  4 │ Reads completed successfully                           │
	│  5 │ Reads merged                                           │
	│  6 │ Sectors read                                           │
	│  7 │ Time spent reading (ms)                                │
	│  8 │ Writes completed                                       │
	│  9 │ Writes merged                                          │
	│ 10 │ Sectors written                                        │
	│ 11 │ Time spent writing (ms)                                │
	│ 12 │ I/Os currently in progress                             │
	│ 13 │ Time spent doing I/Os (ms)                             │
	│ 14 │ Weighted time spent doing I/Os (ms)                    │
	└────┴────────────────────────────────────────────────────────┘

Linux Kernel 4.18+ appends 4 discard-related fields:

	┌────┬────────────────────────────────────────────────────────┐
	│ 15 │ Discards completed successfully                        │
	│ 16 │ Discards merged                                        │
	│ 17 │ Sectors discarded                                      │
	│ 18 │ Time spent discarding (ms)                             │
	└────┴────────────────────────────────────────────────────────┘

Linux Kernel 5.5+ appends 2 flush-related fields:

	┌────┬────────────────────────────────────────────────────────┐
	│ 19 │ Flush requests completed successfully                  │
	│ 20 │ Time spent flushing (ms)                               │
	└────┴────────────────────────────────────────────────────────┘
*/
type ProcDiskStats struct {
	Major             uint32
	Minor             uint32
	Device            string
	ReadsComplete     uint64
	ReadsMerged       uint64
	SectorsRead       uint64
	TimeReading       time.Duration
	WritesCompleted   uint64
	WritesMerged      uint64
	SectorsWritten    uint64
	TimeWriting       time.Duration
	CurrentIO         uint64
	Time              time.Duration
	WeightTime        time.Duration
	DiscardsCompleted uint64
	DiscardsMerged    uint64
	SectorsDiscarded  uint64
	TimeDiscarding    time.Duration
	FlushCompleted    uint64
	TimeFlushing      time.Duration
}

func ReadProcDisksStats() ([]ProcDiskStats, error) {
	data, err := procDiskStats.Data()
	if err != nil {
		return nil, err
	}

	var (
		dataLines = bytes.Split(data, []byte{'\n'})
		stats     []ProcDiskStats
	)

	for _, line := range dataLines {
		fList := bytes.Fields(line)
		if len(fList) < 20 {
			continue
		}

		stat := ProcDiskStats{
			Major:             bytesToUint32(fList[0]),
			Minor:             bytesToUint32(fList[1]),
			Device:            string(fList[2]),
			ReadsComplete:     bytesToUint64(fList[3]),
			ReadsMerged:       bytesToUint64(fList[4]),
			SectorsRead:       bytesToUint64(fList[5]),
			TimeReading:       time.Duration(bytesToUint64(fList[6])) * time.Millisecond,
			WritesCompleted:   bytesToUint64(fList[7]),
			WritesMerged:      bytesToUint64(fList[8]),
			SectorsWritten:    bytesToUint64(fList[9]),
			TimeWriting:       time.Duration(bytesToUint64(fList[10])) * time.Millisecond,
			CurrentIO:         bytesToUint64(fList[11]),
			Time:              time.Duration(bytesToUint64(fList[12])) * time.Millisecond,
			WeightTime:        time.Duration(bytesToUint64(fList[13])) * time.Millisecond,
			DiscardsCompleted: bytesToUint64(fList[14]),
			DiscardsMerged:    bytesToUint64(fList[15]),
			SectorsDiscarded:  bytesToUint64(fList[16]),
			TimeDiscarding:    time.Duration(bytesToUint64(fList[17])) * time.Millisecond,
			FlushCompleted:    bytesToUint64(fList[18]),
			TimeFlushing:      time.Duration(bytesToUint64(fList[19])) * time.Millisecond,
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func readPartitionStat(stat psdisk.PartitionStat) Partition {
	partUse := usageOfPartition(stat.Mountpoint)

	p := Partition{
		Device:     stat.Device,
		MountPoint: stat.Mountpoint,
		FsType:     stat.Fstype,
		Opts:       stat.Opts,
		Usage:      partUse,
	}

	return p
}

func usageOfPartition(p string) *PartitionUsage {

	stat, err := psdisk.Usage(p)
	if err != nil {
		return nil
	}

	var (
		inodeSum = (stat.InodesTotal + stat.InodesFree + stat.InodesUsed)
		useSum   = (stat.Total + stat.Free + stat.Used)
	)

	if inodeSum > 0 || useSum > 0 {
		return &PartitionUsage{
			Total:             stat.Total,
			Free:              stat.Free,
			Used:              stat.Used,
			UsedPercent:       stat.UsedPercent,
			InodesTotal:       stat.InodesTotal,
			InodesFree:        stat.InodesFree,
			InodesUsed:        stat.InodesUsed,
			InodesUsedPercent: stat.InodesUsedPercent,
		}
	}

	return nil
}
