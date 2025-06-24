package procf

import (
	"bytes"
	"fmt"
	"net"
)

func NetCounters() (NetCounter, error) {

	iCounters, err := procInterfaceCounters()
	if err != nil {
		return NetCounter{}, err
	}

	sCounters, err := procSnmpCounters()
	if err != nil {
		return NetCounter{}, err
	}

	data := NetCounter{
		Interface: iCounters,
		Snmp:      sCounters,
	}

	return data, nil
}

func parseInterfaceCounters(data []byte) (InterfaceCounters, error) {

	fileLines := bytes.Split(data, []byte("\n"))

	if len(fileLines) <= 2 {
		return nil, fmt.Errorf("net counter fetch error: invalid file")
	}

	ifaces := make(InterfaceCounters, len(fileLines)-2)

	for _, byteLine := range fileLines[2 : len(fileLines)-1] {

		var (
			fields = bytes.Fields(byteLine)
			iface  = fields[0]
		)

		if len(fields) < 17 {
			return InterfaceCounters{}, fmt.Errorf("net counter fetch error: invalid file")
		}

		ifaces[string(iface[:len(iface)-1])] = InterfaceCounter{
			Receive:  parseNetRx(fields[1:9]),
			Transmit: parseNetTx(fields[9:17]),
		}
	}

	return ifaces, nil
}

func procInterfaceCounters() (InterfaceCounters, error) {

	data, err := procNetDev.Data()
	if err != nil {
		return nil, fmt.Errorf("net counter fetch error: %v", err)
	}

	return parseInterfaceCounters(data)
}

func parseNetRx(data [][]byte) NetRX {

	if len(data) < 8 {
		return NetRX{}
	}

	return NetRX{
		Bytes:      bytesToUint64(data[0]),
		Packets:    bytesToUint64(data[1]),
		Errs:       bytesToUint64(data[2]),
		Drop:       bytesToUint64(data[3]),
		Fifo:       bytesToUint64(data[4]),
		Frame:      bytesToUint64(data[5]),
		Compressed: bytesToUint64(data[6]),
		Multicast:  bytesToUint64(data[7]),
	}
}

func parseNetTx(data [][]byte) NetTX {

	if len(data) < 8 {
		return NetTX{}
	}

	return NetTX{
		Bytes:      bytesToUint64(data[0]),
		Packets:    bytesToUint64(data[1]),
		Errs:       bytesToUint64(data[2]),
		Drop:       bytesToUint64(data[3]),
		Fifo:       bytesToUint64(data[4]),
		Colls:      bytesToUint64(data[5]),
		Carrier:    bytesToUint64(data[6]),
		Compressed: bytesToUint64(data[7]),
	}
}

func procSnmpCounters() (SnmpCounter, error) {
	data, err := procNetSnmp.Data()
	if err != nil {
		return SnmpCounter{}, err
	}

	fileLines := bytes.Split(data, []byte("\n"))

	if len(fileLines) < 10 {
		return SnmpCounter{}, fmt.Errorf("net counter fetch error: invalid file")
	}

	cntr := SnmpCounter{
		Ip:      parseIpSnmp(bytes.Fields(fileLines[1])),
		Icmp:    parseIcmpSnmp(bytes.Fields(fileLines[3])),
		Tcp:     parseTcpSnmp(bytes.Fields(fileLines[5])),
		Udp:     parseUdpSnmp(bytes.Fields(fileLines[7])),
		UdpLite: parseUdpLiteSnmp(bytes.Fields(fileLines[9])),
	}

	return cntr, nil
}

func parseIpSnmp(data [][]byte) IpSnmp {

	if len(data) < 21 {
		return IpSnmp{}
	}

	return IpSnmp{
		Forwarding:      bytesToInt64(data[1]),
		DefaultTTL:      bytesToInt64(data[2]),
		InReceives:      bytesToInt64(data[3]),
		InHdrErrors:     bytesToInt64(data[4]),
		InAddrErrors:    bytesToInt64(data[5]),
		ForwDatagrams:   bytesToInt64(data[6]),
		InUnknownProtos: bytesToInt64(data[7]),
		InDiscards:      bytesToInt64(data[8]),
		InDelivers:      bytesToInt64(data[9]),
		OutRequests:     bytesToInt64(data[10]),
		OutDiscards:     bytesToInt64(data[11]),
		OutNoRoutes:     bytesToInt64(data[12]),
		ReasmTimeout:    bytesToInt64(data[13]),
		ReasmReqds:      bytesToInt64(data[14]),
		ReasmOKs:        bytesToInt64(data[15]),
		ReasmFails:      bytesToInt64(data[16]),
		FragOKs:         bytesToInt64(data[17]),
		FragFails:       bytesToInt64(data[18]),
		FragCreates:     bytesToInt64(data[19]),
		OutTransmits:    bytesToInt64(data[20]),
	}
}

func parseIcmpSnmp(data [][]byte) IcmpSnmp {

	if len(data) < 30 {
		return IcmpSnmp{}
	}

	return IcmpSnmp{
		InMsgs:             bytesToInt64(data[1]),
		InErrors:           bytesToInt64(data[2]),
		InCsumErrors:       bytesToInt64(data[3]),
		InDestUnreachs:     bytesToInt64(data[4]),
		InTimeExcds:        bytesToInt64(data[5]),
		InParmProbs:        bytesToInt64(data[6]),
		InSrcQuenchs:       bytesToInt64(data[7]),
		InRedirects:        bytesToInt64(data[8]),
		InEchos:            bytesToInt64(data[9]),
		InEchoReps:         bytesToInt64(data[10]),
		InTimestamps:       bytesToInt64(data[11]),
		InTimestampReps:    bytesToInt64(data[12]),
		InAddrMasks:        bytesToInt64(data[13]),
		InAddrMaskReps:     bytesToInt64(data[14]),
		OutMsgs:            bytesToInt64(data[15]),
		OutErrors:          bytesToInt64(data[16]),
		OutRateLimitGlobal: bytesToInt64(data[17]),
		OutRateLimitHost:   bytesToInt64(data[18]),
		OutDestUnreachs:    bytesToInt64(data[19]),
		OutTimeExcds:       bytesToInt64(data[20]),
		OutParmProbs:       bytesToInt64(data[21]),
		OutSrcQuenchs:      bytesToInt64(data[22]),
		OutRedirects:       bytesToInt64(data[23]),
		OutEchos:           bytesToInt64(data[24]),
		OutEchoReps:        bytesToInt64(data[25]),
		OutTimestamps:      bytesToInt64(data[26]),
		OutTimestampReps:   bytesToInt64(data[27]),
		OutAddrMasks:       bytesToInt64(data[28]),
		OutAddrMaskReps:    bytesToInt64(data[29]),
	}
}

func parseTcpSnmp(data [][]byte) TcpSnmp {

	if len(data) < 16 {
		return TcpSnmp{}
	}

	return TcpSnmp{
		RtoAlgorithm: bytesToInt64(data[1]),
		RtoMin:       bytesToInt64(data[2]),
		RtoMax:       bytesToInt64(data[3]),
		MaxConn:      bytesToInt64(data[4]),
		ActiveOpens:  bytesToInt64(data[5]),
		PassiveOpens: bytesToInt64(data[6]),
		AttemptFails: bytesToInt64(data[7]),
		EstabResets:  bytesToInt64(data[8]),
		CurrEstab:    bytesToInt64(data[9]),
		InSegs:       bytesToInt64(data[10]),
		OutSegs:      bytesToInt64(data[11]),
		RetransSegs:  bytesToInt64(data[12]),
		InErrs:       bytesToInt64(data[13]),
		OutRsts:      bytesToInt64(data[14]),
		InCsumErrors: bytesToInt64(data[15]),
	}
}

func parseUdpSnmp(data [][]byte) UdpSnmp {

	if len(data) < 10 {
		return UdpSnmp{}
	}

	return UdpSnmp{
		InDatagrams:  bytesToInt64(data[1]),
		NoPorts:      bytesToInt64(data[2]),
		InErrors:     bytesToInt64(data[3]),
		OutDatagrams: bytesToInt64(data[4]),
		RcvbufErrors: bytesToInt64(data[5]),
		SndbufErrors: bytesToInt64(data[6]),
		InCsumErrors: bytesToInt64(data[7]),
		IgnoredMulti: bytesToInt64(data[8]),
		MemErrors:    bytesToInt64(data[9]),
	}
}

func parseUdpLiteSnmp(data [][]byte) UdpSnmp {

	if len(data) < 10 {
		return UdpSnmp{}
	}

	return UdpSnmp{
		InDatagrams:  bytesToInt64(data[1]),
		NoPorts:      bytesToInt64(data[2]),
		InErrors:     bytesToInt64(data[3]),
		OutDatagrams: bytesToInt64(data[4]),
		RcvbufErrors: bytesToInt64(data[5]),
		SndbufErrors: bytesToInt64(data[6]),
		InCsumErrors: bytesToInt64(data[7]),
		IgnoredMulti: bytesToInt64(data[8]),
		MemErrors:    bytesToInt64(data[9]),
	}
}

type NetInterfaces struct {
}

func FetchNetInterfaces() {
	net.Interfaces()
}
