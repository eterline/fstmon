// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package wol

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

var (
	ErrInvalidMac = errors.New("invalid mac string")
)

var (
	wolMagicBytes = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
)

const (
	WoLPort1 = 7
	WoLPort2 = 9
)

/*
NewWoLSender - wol sending type

	send WoL packets to target host and MAC
*/
type WoLSender struct {
	addr *net.UDPAddr
	mac  [6]byte
}

/*
NewWoLSender - create new sender instance

	Signature example:
	target, err := NewWoLSender("60:be:b4:1a:22:68", "192.168.1.16", 7)

	err returns in invalid MAC string or IP parse fail
*/
func NewWoLSender(mac string, addr string, port int) (*WoLSender, error) {

	macAddr, ok := parseMAC(mac)
	if !ok {
		return nil, ErrInvalidMac
	}

	resAddr := fmt.Sprintf("%s:%d", addr, port)
	udpAddr, err := net.ResolveUDPAddr("udp", resAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid udp addr: %v", err)
	}

	snd := &WoLSender{
		addr: udpAddr,
		mac:  macAddr,
	}

	return snd, nil
}

/*
Mac - return interface mac addr as a string

	Example: target.Mac() -> "60:be:b4:1a:22:68"
*/
func (wol WoLSender) Mac() string {
	return mac(wol.mac[:])
}

/*
Addr - return host target full addr

	Example: target.Addr() -> "192.168.1.16:7"
*/
func (wol WoLSender) Addr() string {
	return wol.addr.String()
}

/*
Addr - return host target IP

	Example: target.IP() -> 192.168.1.16 (net.IP)
*/
func (wol WoLSender) IP() net.IP {
	return wol.addr.IP
}

/*
Addr - return host target wol port

	Example: target.Port() -> 7
*/
func (wol WoLSender) Port() int {
	return wol.addr.Port
}

/*
Wake - send WoL packets to target

	localAddr - source address from client host.
	 If field will be "", the packets will be Broadcasted from all IPs
	times - times of packet repeat before cycle stop

	Example: Wake("", 10) -> broadcast 10 times
	Example: Wake("192.168.1.12:8001", 10) -> send from 192.168.1.12:8001 port 10 times
*/
func (wol WoLSender) Wake(localAddr string, times int) error {
	return wol.WakeWithContext(context.Background(), localAddr, times)
}

/*
Wake - send WoL packets to target

	Can be cancelled by context

	localAddr - source address from client host.
	 If field will be "", the packets will be Broadcasted from all IPs
	times - times of packet repeat before cycle stop

	Example: Wake("", 10) -> broadcast 10 times
	Example: Wake("192.168.1.12:8001", 10) -> send from 192.168.1.12:8001 port 10 times
*/
func (wol WoLSender) WakeWithContext(ctx context.Context, localAddr string, times int) error {
	udpAddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil && localAddr != "" {
		return fmt.Errorf("invalid udp addr: %v", err)
	}
	return wol.wake(ctx, udpAddr, times)
}

// wake - internal realization
func (wol WoLSender) wake(ctx context.Context, localAddr *net.UDPAddr, times int) error {

	// connect to host
	conn, err := net.DialUDP("udp", localAddr, wol.addr)
	if err != nil {
		return fmt.Errorf("wol packet send failed: %v", err)
	}
	defer conn.Close()

	if times < 1 {
		times = 1
	}

	// generate packet
	paket := wol.magicPacket()

	for i := 0; i < times; i++ {

		// cancelling loop if done
		if ctx.Err() != nil {
			return fmt.Errorf("wol packet send error: %v", err)
		}

		_, err := conn.Write(paket)
		if err != nil {
			return fmt.Errorf("wol packet send error: %v", err)
		}
	}

	return nil
}

func (wol WoLSender) Ping() []byte {

	out, err := exec.Command("ping", "-c", "3", wol.IP().String()).CombinedOutput()
	if err != nil {
		fmt.Println("Ping failed:", err)
	}

	fmt.Println(string(out))

	return nil
}

// magicPacket - generate magic packet from bytes and mac
func (wol WoLSender) magicPacket() []byte {
	buf := &bytes.Buffer{}

	buf.Write(wolMagicBytes)

	for i := 0; i < 6; i++ {
		buf.Write(wol.mac[:])
	}

	return buf.Bytes()
}

// parseMAC - parse mac from string
func parseMAC(s string) ([6]byte, bool) {

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	parts := strings.Split(s, ":")

	if len(parts) != 6 {
		return [6]byte{}, false
	}

	mac := [6]byte{}

	for i := 0; i < 6; i++ {
		part, err := hex.DecodeString(parts[i])
		if err != nil {
			return [6]byte{}, false
		}
		mac[i] = part[0]
	}

	return mac, true
}

// mac - convert hex slice to MAC
func mac(p []byte) string {
	if len(p) < 6 {
		return "NULL"
	}
	return fmt.Sprintf(
		"%02x:%02x:%02x:%02x:%02x:%02x",
		p[0], p[1], p[2],
		p[3], p[4], p[5],
	)
}
