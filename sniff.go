package main

import (
	"fmt"
	"syscall"
        "strings"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
        "time"
        "strconv"
        "os"
)

type Kv struct{
	IP string
	Port int
	Data string
        Ts int64
}

func main() {
    proto := (syscall.ETH_P_ALL<<8)&0xff00 | syscall.ETH_P_ALL>>8 // change to Big-Endian order
    fd, _ := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, proto)
    port, _ := strconv.Atoi(os.Args[1])
    p := make(map[string]*Kv)
    buf := make([]byte, 65536)
    replacer := strings.NewReplacer("\r", " ", "\n", " ")
    for {
	n, _, _ := syscall.Recvfrom(fd, buf, 0)
	packet := gopacket.NewPacket(buf[:n], layers.LayerTypeEthernet, gopacket.Default)
        if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer == nil { continue }
        ts := time.Now().UnixNano() / 1e3
        data := packet.TransportLayer().LayerPayload()
        if len(data) < 1 {continue}
        ipLayer := packet.Layer(layers.LayerTypeIPv4)
        ip, _ := ipLayer.(*layers.IPv4)
        tcpLayer := packet.Layer(layers.LayerTypeTCP)
        tcp, _ := tcpLayer.(*layers.TCP)
        if int(tcp.DstPort) == port {
            key := fmt.Sprintf("%s:%s", ip.SrcIP,tcp.SrcPort)
            p[key] = new(Kv)
            p[key].IP = fmt.Sprintf("%s",ip.SrcIP)
            p[key].Port = int(tcp.SrcPort)
            p[key].Data = replacer.Replace(string(data))
            p[key].Ts = ts
        }
        if int(tcp.SrcPort) == port {
            key := fmt.Sprintf("%s:%s", ip.DstIP,tcp.DstPort)
            value, ok := p[key]
            if ok == false {continue}
            if p[key].Ts == -1 {continue}
            now := time.Now()
            hour, min, sec := now.Clock()
            fmt.Printf("%02d:%02d:%02d %d %s:%d,%s",hour, min, sec, ts-value.Ts, value.IP, value.Port, value.Data)
            fmt.Printf("<===> %s\n",replacer.Replace(string(data)))
            p[key].Ts = -1
        }
     }
}
