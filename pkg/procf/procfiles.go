// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package procf

import (
	"time"
)

/*
ProcLoadAvg holds average load

	┌───────────────┬──────────────────────────────────────────────────────────────────────────┐
	│ Field         │ Description                                                              │
	├───────────────┼──────────────────────────────────────────────────────────────────────────┤
	│ Load1         │ The system load average over the last 1 minute.                          │
	│ Load5         │ The system load average over the last 5 minutes.                         │
	│ Load15        │ The system load average over the last 15 minutes.                        │
	│ RunningProcs  │ Number of currently running processes and total number of processes, e.g.│
	│               │ "2/150" means 2 running out of 150 total processes.                      │
	│ LastPID       │ The process ID of the last process created.                              │
	└───────────────┴──────────────────────────────────────────────────────────────────────────┘
*/
type ProcLoadAvg struct {
	Load1        float64 `json:"load_1"`
	Load5        float64 `json:"load_5"`
	Load15       float64 `json:"load_15"`
	RunningProcs string  `json:"running_procs"`
	LastPID      int32   `json:"last_pid"`
}

/*
ProcUptime holds host running time

	┌─────────────┬───────────────────────────────────────────────────────────────────────┐
	│ Field       │ Description                                                           │
	├─────────────┼───────────────────────────────────────────────────────────────────────┤
	│ Uptime      │ The total number of seconds the system has been running since boot.   │
	│ IdleTime    │ The total number of seconds all CPUs have spent idle since boot.      │
	└─────────────┴───────────────────────────────────────────────────────────────────────┘
*/
type ProcUptime struct {
	Uptime   time.Duration `json:"uptime"`
	IdleTime time.Duration `json:"idle_time"`
}

/*
ProcCrypto holds kernel crypto modules

	┌─────────────┬──────────────────────────────────────────────────────────────────────────────┐
	│ Field       │ Description                                                                  │
	├─────────────┼──────────────────────────────────────────────────────────────────────────────┤
	│ Name        │ The name of the cryptographic module, identifying it in the system.          │
	│ Driver      │ The driver name implementing the cryptographic module.                       │
	│ Module      │ The kernel module or component loaded to provide this cryptographic support. │
	│ Priority    │ The priority of the module, influencing selection order in crypto operations.│
	│ Refcnt      │ Reference count indicating how many users are currently using the module.    │
	│ SelfTest    │ Result of the module self-test: usually passed, failed, or not supported.    │
	│ Internal    │ Indicates whether the module is internal to the system or for internal use.  │
	│ Type        │ The type of cryptographic module, e.g., cipher, hash, rng.                   │
	└─────────────┴──────────────────────────────────────────────────────────────────────────────┘
*/
type ProcCrypto []ProcCryptoModule

type ProcCryptoModule struct {
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	Module   string `json:"module"`
	Priority int32  `json:"priority"`
	Refcnt   int32  `json:"ref_cnt"`
	SelfTest string `json:"self_test"`
	Internal string `json:"internal"`
	Type     string `json:"type"`
}

/*
ProcCpuInfo holds kernel crypto modules

	┌────────────────────┬──────────────────────────────────────────────────────────────┐
	│ Field              │ Description                                                  │
	├────────────────────┼──────────────────────────────────────────────────────────────┤
	│ Processor          │ Logical processor number                                     │
	│ VendorID           │ CPU manufacturer (e.g., "GenuineIntel", "AuthenticAMD")      │
	│ CpuFamily          │ CPU family as defined by the vendor                          │
	│ Model              │ Model number within the CPU family                           │
	│ ModelName          │ Full model name of the processor                             │
	│ Stepping           │ CPU stepping/revision number                                 │
	│ Microcode          │ Microcode version currently loaded                           │
	│ CpuMhz             │ Current CPU frequency in MHz                                 │
	│ CacheSize          │ L2 (or L3) cache size in KB                                  │
	│ PhysicalID         │ Identifier for the physical processor package (socket)       │
	│ Siblings           │ Total number of logical processors on the same physical CPU  │
	│ CoreID             │ Identifier for the core within a physical processor          │
	│ CpuCores           │ Number of physical cores in the processor                    │
	│ Apicid             │ APIC (Advanced Programmable Interrupt Controller) ID         │
	│ InitialApicid      │ Initial APIC ID assigned at boot                             │
	│ Fpu                │ Indicates if a Floating Point Unit is present                │
	│ FpuException       │ Indicates if FPU exceptions are supported                    │
	│ CpuIDLevel         │ Maximum CPUID function supported                             │
	│ WP                 │ Write Protect flag status (memory page protection)           │
	│ Flags              │ Supported CPU instruction flags (e.g., sse, avx, etc.)       │
	│ VmxFlags           │ Virtualization-specific CPU flags                            │
	│ Bugs               │ Known CPU bugs or errata                                     │
	│ BogoMips           │ Bogus MIPS value (used for kernel timing)                    │
	│ TLB                │ Sizes of different TLBs (Translation Lookaside Buffers)      │
	│ ClflushSize        │ Cache line size for CLFLUSH instruction                      │
	│ CacheAlignment     │ Memory cache alignment in bytes                              │
	│ AddressSizes       │ Width of physical and virtual address buses                  │
	│ PowerManagement    │ Supported power management features                          │
	└────────────────────┴──────────────────────────────────────────────────────────────┘
*/
type ProcCpuInfo struct {
	Cores []CoreInfo `json:"cores"`
}

// ProcCpuInfo detail info in ProcCpuInfo description
type CoreInfo struct {
	Processor       int64        `json:"processor"`
	VendorID        string       `json:"vendor_id"`
	CpuFamily       int64        `json:"cpu_family"`
	Model           int64        `json:"model"`
	ModelName       string       `json:"model_name"`
	Stepping        int64        `json:"stepping"`
	Microcode       string       `json:"microcode"`
	CpuMhz          float64      `json:"cpu_mhz"`
	CacheSize       int64        `json:"cache_size"`
	PhysicalID      int64        `json:"physical_id"`
	Siblings        int64        `json:"siblings"`
	CoreID          int64        `json:"core_id"`
	CpuCores        int64        `json:"cpu_cores"`
	Apicid          int64        `json:"apicid"`
	InitialApicid   int64        `json:"initial_apicid"`
	Fpu             bool         `json:"fpu"`
	FpuException    bool         `json:"fpu_exception"`
	CpuIDLevel      int64        `json:"cpuid_level"`
	WP              bool         `json:"wp"`
	Flags           Instructions `json:"flags"`
	VmxFlags        []string     `json:"vmx_flags"`
	Bugs            []string     `json:"bugs"`
	BogoMips        float64      `json:"bogomips"`
	TLB             SizesTLB     `json:"tlb_sizes"`
	ClflushSize     int64        `json:"clflush_size"`
	CacheAlignment  int64        `json:"cache_alignment"`
	AddressSizes    AddrBits     `json:"address_sizes"`
	PowerManagement []string     `json:"power_management"`
}

// InterfaceCounters - map of interface counters stats (RX and TX)
type InterfaceCounters map[string]InterfaceCounter

/*
NetRX – Receive segment

	┌────────────┬────────────────────────────────────────────────────────────┐
	│ Field      │ Description                                                │
	├────────────┼────────────────────────────────────────────────────────────┤
	│ Bytes      │ Total number of bytes received                             │
	│ Packets    │ Total number of packets received                           │
	│ Errs       │ Number of receive errors                                   │
	│ Drop       │ Number of received packets dropped                         │
	│ Fifo       │ Number of FIFO buffer errors on receive                    │
	│ Frame      │ Number of frame alignment errors                           │
	│ Compressed │ Number of compressed packets received (if any)             │
	│ Multicast  │ Number of multicast packets received                       │
	└────────────┴────────────────────────────────────────────────────────────┘
*/
type NetRX struct {
	Bytes      uint64 `json:"bytes"`
	Packets    uint64 `json:"packets"`
	Errs       uint64 `json:"errs"`
	Drop       uint64 `json:"drop"`
	Fifo       uint64 `json:"fifo"`
	Frame      uint64 `json:"frame"`
	Compressed uint64 `json:"compressed"`
	Multicast  uint64 `json:"multicast"`
}

/*
NetTX – Transmit segment

	┌────────────┬────────────────────────────────────────────────────────────┐
	│ Field      │ Description                                                │
	├────────────┼────────────────────────────────────────────────────────────┤
	│ Bytes      │ Total number of bytes transmitted                          │
	│ Packets    │ Total number of packets transmitted                        │
	│ Errs       │ Number of transmission errors                              │
	│ Drop       │ Number of transmitted packets dropped                      │
	│ Fifo       │ Number of FIFO buffer errors on transmit                   │
	│ Colls      │ Number of collisions detected during transmission          │
	│ Carrier    │ Number of carrier sense errors                             │
	│ Compressed │ Number of compressed packets transmitted (if any)          │
	└────────────┴────────────────────────────────────────────────────────────┘
*/
type NetTX struct {
	Bytes      uint64 `json:"bytes"`
	Packets    uint64 `json:"packets"`
	Errs       uint64 `json:"errs"`
	Drop       uint64 `json:"drop"`
	Fifo       uint64 `json:"fifo"`
	Colls      uint64 `json:"colls"`
	Carrier    uint64 `json:"carrier"`
	Compressed uint64 `json:"compressed"`
}

/*
IpSnmp - IP layer statistics, including packet forwarding, errors, fragmentation, and reassembly

	┌───────────────────┬─────────────────────────────────────────────────────────────┐
	│ Field             │ Description                                                 │
	├───────────────────┼─────────────────────────────────────────────────────────────┤
	│ Forwarding        │ IP packet forwarding enabled (1) or disabled (0)            │
	│ DefaultTTL        │ Default Time-To-Live (TTL) value for outgoing packets       │
	│ InReceives        │ Total IP packets received                                   │
	│ InHdrErrors       │ Packets dropped due to header errors (e.g., checksum)       │
	│ InAddrErrors      │ Packets dropped due to invalid destination addresses        │
	│ ForwDatagrams     │ Packets forwarded to another interface                      │
	│ InUnknownProtos   │ Packets with unknown protocol discarded                     │
	│ InDiscards        │ Incoming packets discarded (e.g., buffer full)              │
	│ InDelivers        │ Packets successfully delivered to higher layers             │
	│ OutRequests       │ Outgoing IP packets requested to be sent                    │
	│ OutDiscards       │ Outgoing packets discarded (e.g., buffer full)              │
	│ OutNoRoutes       │ Packets dropped due to missing route                        │
	│ ReasmTimeout      │ Timeout for IP fragment reassembly                          │
	│ ReasmReqds        │ Fragments needing reassembly                                │
	│ ReasmOKs          │ Successfully reassembled packets                            │
	│ ReasmFails        │ Failed reassembly attempts                                  │
	│ FragOKs           │ Successfully fragmented packets                             │
	│ FragFails         │ Failed fragmentation attempts                               │
	│ FragCreates       │ Fragments generated                                         │
	│ OutTransmits      │ Outgoing IP packets transmitted                             │
	└───────────────────┴─────────────────────────────────────────────────────────────┘
*/
type IpSnmp struct {
	Forwarding      int64 `json:"forwarding"`
	DefaultTTL      int64 `json:"default_ttl"`
	InReceives      int64 `json:"in_receives"`
	InHdrErrors     int64 `json:"in_hdr_errors"`
	InAddrErrors    int64 `json:"in_addr_errors"`
	ForwDatagrams   int64 `json:"forw_datagrams"`
	InUnknownProtos int64 `json:"in_unknown_protos"`
	InDiscards      int64 `json:"in_discards"`
	InDelivers      int64 `json:"in_delivers"`
	OutRequests     int64 `json:"out_requests"`
	OutDiscards     int64 `json:"out_discards"`
	OutNoRoutes     int64 `json:"out_no_routes"`
	ReasmTimeout    int64 `json:"reasm_timeout"`
	ReasmReqds      int64 `json:"reasm_reqds"`
	ReasmOKs        int64 `json:"reasm_oks"`
	ReasmFails      int64 `json:"reasm_fails"`
	FragOKs         int64 `json:"frag_oks"`
	FragFails       int64 `json:"frag_fails"`
	FragCreates     int64 `json:"frag_creates"`
	OutTransmits    int64 `json:"out_transmits"`
}

/*
IcmpSnmp - TCP protocol statistics, including connection states, retransmissions, and errors

	┌──────────────────────┬───────────────────────────────────────────────────┐
	│ Field                │ Description                                       │
	├──────────────────────┼───────────────────────────────────────────────────┤
	│ InMsgs               │ Total ICMP messages received                      │
	│ InErrors             │ Corrupted/invalid ICMP messages received          │
	│ InCsumErrors         │ ICMP checksum errors                              │
	│ InDestUnreachs       │ "Destination Unreachable" messages received       │
	│ InTimeExcds          │ "Time Exceeded" messages received                 │
	│ InParmProbs          │ "Parameter Problem" messages received             │
	│ InSrcQuenchs         │ "Source Quench" messages received (legacy)        │
	│ InRedirects          │ Redirect messages received                        │
	│ InEchos              │ ICMP Echo (ping) requests received                │
	│ InEchoReps           │ ICMP Echo replies received                        │
	│ InTimestamps         │ Timestamp requests received                       │
	│ InTimestampReps      │ Timestamp replies received                        │
	│ InAddrMasks          │ Address Mask requests received                    │
	│ InAddrMaskReps       │ Address Mask replies received                     │
	│ OutMsgs              │ Total ICMP messages sent                          │
	│ OutErrors            │ Failed ICMP message transmissions                 │
	│ OutRateLimitGlobal   │ Global rate-limited outgoing ICMP messages        │
	│ OutRateLimitHost     │ Host-specific rate-limited outgoing ICMP messages │
	│ OutDestUnreachs      │ "Destination Unreachable" messages sent           │
	│ OutTimeExcds         │ "Time Exceeded" messages sent                     │
	│ OutParmProbs         │ "Parameter Problem" messages sent                 │
	│ OutSrcQuenchs        │ "Source Quench" messages sent (legacy)            │
	│ OutRedirects         │ Redirect messages sent                            │
	│ OutEchos             │ ICMP Echo (ping) requests sent                    │
	│ OutEchoReps          │ ICMP Echo replies sent                            │
	│ OutTimestamps        │ Timestamp requests sent                           │
	│ OutTimestampReps     │ Timestamp replies sent                            │
	│ OutAddrMasks         │ Address Mask requests sent                        │
	│ OutAddrMaskReps      │ Address Mask replies sent                         │
	└──────────────────────┴───────────────────────────────────────────────────┘
*/
type IcmpSnmp struct {
	InMsgs             int64 `json:"in_msgs"`
	InErrors           int64 `json:"in_errors"`
	InCsumErrors       int64 `json:"in_csum_errors"`
	InDestUnreachs     int64 `json:"in_dest_unreachs"`
	InTimeExcds        int64 `json:"in_time_excds"`
	InParmProbs        int64 `json:"in_parm_probs"`
	InSrcQuenchs       int64 `json:"in_src_quenchs"`
	InRedirects        int64 `json:"in_redirects"`
	InEchos            int64 `json:"in_echos"`
	InEchoReps         int64 `json:"in_echo_reps"`
	InTimestamps       int64 `json:"in_timestamps"`
	InTimestampReps    int64 `json:"in_timestamp_reps"`
	InAddrMasks        int64 `json:"in_addr_masks"`
	InAddrMaskReps     int64 `json:"in_addr_mask_reps"`
	OutMsgs            int64 `json:"out_msgs"`
	OutErrors          int64 `json:"out_errors"`
	OutRateLimitGlobal int64 `json:"out_rate_limit_global"`
	OutRateLimitHost   int64 `json:"out_rate_limit_host"`
	OutDestUnreachs    int64 `json:"out_dest_unreachs"`
	OutTimeExcds       int64 `json:"out_time_excds"`
	OutParmProbs       int64 `json:"out_parm_probs"`
	OutSrcQuenchs      int64 `json:"out_src_quenchs"`
	OutRedirects       int64 `json:"out_redirects"`
	OutEchos           int64 `json:"out_echos"`
	OutEchoReps        int64 `json:"out_echo_reps"`
	OutTimestamps      int64 `json:"out_timestamps"`
	OutTimestampReps   int64 `json:"out_timestamp_reps"`
	OutAddrMasks       int64 `json:"out_addr_masks"`
	OutAddrMaskReps    int64 `json:"out_addr_mask_reps"`
}

/*
TcpSnmp - TCP protocol statistics, including connection states, retransmissions, and errors

	┌──────────────────┬───────────────────────────────────────────────────┐
	│ Field            │ Description                                       │
	├──────────────────┼───────────────────────────────────────────────────┤
	│ RtoAlgorithm     │ Retransmission timeout algorithm (e.g., RFC 6298) │
	│ RtoMin           │ Minimum retransmission timeout (ms)               │
	│ RtoMax           │ Maximum retransmission timeout (ms)               │
	│ MaxConn          │ Maximum allowed TCP connections (deprecated)      │
	│ ActiveOpens      │ Connections initiated                             │
	│ PassiveOpens     │ Connections accepted                              │
	│ AttemptFails     │ Failed connection attempts                        │
	│ EstabResets      │ Connections reset (RST packets sent)              │
	│ CurrEstab        │ Currently established connections                 │
	│ InSegs           │ Incoming TCP segments                             │
	│ OutSegs          │ Outgoing TCP segments                             │
	│ RetransSegs      │ Retransmitted segments                            │
	│ InErrs           │ Corrupted/invalid segments received               │
	│ OutRsts          │ Outgoing RST packets (resets)                     │
	│ InCsumErrors     │ TCP checksum errors                               │
	└──────────────────┴───────────────────────────────────────────────────┘
*/
type TcpSnmp struct {
	RtoAlgorithm int64 `json:"rto_algorithm"`
	RtoMin       int64 `json:"rto_min"`
	RtoMax       int64 `json:"rto_max"`
	MaxConn      int64 `json:"max_conn"`
	ActiveOpens  int64 `json:"active_opens"`
	PassiveOpens int64 `json:"passive_opens"`
	AttemptFails int64 `json:"attempt_fails"`
	EstabResets  int64 `json:"estab_resets"`
	CurrEstab    int64 `json:"curr_estab"`
	InSegs       int64 `json:"in_segs"`
	OutSegs      int64 `json:"out_segs"`
	RetransSegs  int64 `json:"retrans_segs"`
	InErrs       int64 `json:"in_errs"`
	OutRsts      int64 `json:"out_rsts"`
	InCsumErrors int64 `json:"in_csum_errors"`
}

/*
UdpSnmp - UDP and UDP-Lite statistics, covering datagrams, errors, and buffer issues

	┌──────────────────┬────────────────────────────────────────┐
	│ Field            │ Description                            │
	├──────────────────┼────────────────────────────────────────┤
	│ InDatagrams      │ Datagrams received                     │
	│ NoPorts          │ Datagrams dropped (no listening port)  │
	│ InErrors         │ Corrupted/invalid datagrams received   │
	│ OutDatagrams     │ Datagrams sent                         │
	│ RcvbufErrors     │ Receive buffer errors (overflow)       │
	│ SndbufErrors     │ Send buffer errors (overflow)          │
	│ InCsumErrors     │ Checksum errors (UDP-Lite only)        │
	│ IgnoredMulti     │ Multicast datagrams ignored            │
	│ MemErrors        │ Memory allocation errors               │
	└──────────────────┴────────────────────────────────────────┘
*/
type UdpSnmp struct {
	InDatagrams  int64 `json:"in_datagrams"`
	NoPorts      int64 `json:"no_ports"`
	InErrors     int64 `json:"in_errors"`
	OutDatagrams int64 `json:"out_datagrams"`
	RcvbufErrors int64 `json:"rcvbuf_errors"`
	SndbufErrors int64 `json:"sndbuf_errors"`
	InCsumErrors int64 `json:"in_csum_errors"`
	IgnoredMulti int64 `json:"ignored_multi"`
	MemErrors    int64 `json:"mem_errors"`
}

/*
SnmpCounter - Aggregates all SNMP counters for IP, ICMP, TCP, UDP, and UDP-Lite

	┌──────────┬─────────────────────────────────────────────┐
	│ Field    │ Description                                 │
	├──────────┼─────────────────────────────────────────────┤
	│ Ip       │ IP layer statistics (IpSnmp)                │
	│ Icmp     │ ICMP message statistics (IcmpSnmp)          │
	│ Tcp      │ TCP connection/segment statistics (TcpSnmp) │
	│ Udp      │ UDP datagram statistics (UdpSnmp)           │
	│ UdpLite  │ UDP-Lite datagram statistics (UdpSnmp)      │
	└──────────┴─────────────────────────────────────────────┘
*/
type SnmpCounter struct {
	Ip      IpSnmp   `json:"ip_snmp"`
	Icmp    IcmpSnmp `json:"icmp_snmp"`
	Tcp     TcpSnmp  `json:"tcp_snmp"`
	Udp     UdpSnmp  `json:"udp_snmp"`
	UdpLite UdpSnmp  `json:"udp_lite_snmp"`
}

/*
NetCounter - Aggregates network statistics, including per-interface counter and protocol-level SNMP metrics

	┌───────────┬─────────────────────────────────────────────────────────────┐
	│ Field     │ Description                                                 │
	├───────────┼─────────────────────────────────────────────────────────────┤
	│ Interface │ Per-interface traffic stats                                 │
	│           │ (RX/TX bytes, packets, errors)                              │
	│ Snmp      │ Protocol-level stats (IP, ICMP, TCP, UDP, UDP-Lite)         │
	└───────────┴─────────────────────────────────────────────────────────────┘
*/
type NetCounter struct {
	Interface InterfaceCounters `json:"interface_counters"`
	Snmp      SnmpCounter       `json:"snmp_counter"`
}

/*
InterfaceCounter - Tracks per-interface network statistics, split into receive (RX) and transmit (TX) metrics

	┌──────────┬────────────────────────────────┐
	│ Field    │ Description                    │
	├──────────┼────────────────────────────────┤
	│ Receive  │ Incoming traffic stats (NetRX) │
	│ Transmit │ Outgoing traffic stats (NetTX) │
	└──────────┴────────────────────────────────┘
*/
type InterfaceCounter struct {
	Receive  NetRX `json:"receive"`
	Transmit NetTX `json:"transmit"`
}
