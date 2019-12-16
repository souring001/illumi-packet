package main

import (
    "fmt"
    "flag"
    "strings"
    "time"
    "log"
    "os"
    "net"
    "errors"
    "github.com/jgarff/rpi_ws281x/golang/ws2811"
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
)

type layerMeta struct{
    color uint32
    show bool
}

const (
    ipv4Len    = 4
    pin        = 18     // GPIO
    series     = 6     // length of trail
    count      = 60    // number of LEDs
    brightness = 50     // max 255
)

var (
    speed int     // speed of flowing packet
    device string //eth0
    debug        = flag.Bool("debug", true, "print packet details")
    showip       = flag.Bool("ipaddr", false, "display ip address")
    reset        = flag.Bool("reset", false, "reset LEDs")
    narp         = flag.Bool("narp", false, "disable arp")
    ntcp         = flag.Bool("ntcp", false, "disable tcp")
    nudp         = flag.Bool("nudp", false, "disable udp")
    snapshotLen  = int32(1024)
    promiscuous  = false
    timeout      = 50 * time.Millisecond
    colors = []uint32{
        // GRB color
        0xFFFFFF, //0 White others
        0x880000, //1 Green
        0x00FF00, //2 Red Anomaly
        0x0000FF, //3 Blue TCP
        0x0066cc, //4 Purple ARP
        0x33FF99, //5 Pink ICMP
        0xFFFF00, //6 Yellow UDP
        0x88FF00, //7 Orange IGMP
        0xFF00FF, //8 Cyan DHCP
        0xFF0000, //9 Lime DNS
        0x888888, //10 GRAY
    }
    layerMap = map[string]layerMeta{
        "ARP":      layerMeta{ color: colors[7], show: true },
        "ICMP":     layerMeta{ color: colors[5], show: true },
        "TCP":      layerMeta{ color: colors[3], show: true },
        "UDP":      layerMeta{ color: colors[6], show: true },
        "IGMP":     layerMeta{ color: colors[4], show: true },
        "DNS":      layerMeta{ color: colors[9], show: true },
        "DHCP":     layerMeta{ color: colors[8], show: true },
        "Anomaly":  layerMeta{ color: colors[2], show: true },
        "Others":   layerMeta{ color: colors[0], show: true },
    }
)

func main() {
    // Option flag
    flag.IntVar(&speed, "speed", 1, "set speed of flowing packet")
    flag.StringVar(&device, "device", "eth0", "set network interface")
    flag.Parse()

    meta := layerMap["ARP"]
    meta.show = !*narp
    layerMap["ARP"] = meta

    meta = layerMap["TCP"]
    meta.show = !*ntcp
    layerMap["TCP"] = meta

    meta = layerMap["UDP"]
    meta.show = !*nudp
    layerMap["UDP"] = meta

    // Open device
    handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
    if err != nil { log.Fatal(err) }
    defer handle.Close()

    // Initialize LED strip
    errl := ws2811.Init(pin, count, brightness)
    if errl != nil { log.Fatal(errl) }
    defer ws2811.Fini()

    led := make([]uint32, count)

    if *reset {
        resetLeds(led)
        os.Exit(0)
    }

    // Set IP Address
    ipv4Addr, ipv6Addr, err := externalIP()
    if err != nil { log.Fatal(err) }
    if *debug {
        fmt.Println("IPv4 address:", ipv4Addr)
        fmt.Println("IPv6 address:", ipv6Addr)
    }

    if *showip {
        showIPAddress(led, ipv4Addr)
        os.Exit(0)
    }

    // Use the handle as a packet source to process all packets
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    if *debug {
        fmt.Println("Press Ctr-C to quit.")
        fmt.Println("Start capturing...")
    }

    for packet := range packetSource.Packets() {
        if *debug { fmt.Println("----------------") }

        // Direction of the packet
        reverse := true
        if net := packet.NetworkLayer(); net != nil {
            src, _ := net.NetworkFlow().Endpoints()
            if strings.Contains(src.String(), ipv4Addr) || strings.Contains(src.String(), ipv6Addr) {
                reverse = false
            }
        }

        packetName := categorizePacket(packet)
        layerMeta := layerMap[packetName]

        if *debug {
            fmt.Println(packetName)
            fmt.Println(packet)
        }

        packetTime := packet.Metadata().Timestamp
        nowTime := time.Now()
        diffTime := nowTime.Sub(packetTime)
        if *debug { fmt.Println("delay:", diffTime) }
        if diffTime > 5 * time.Second {
            if *debug { fmt.Println("skip\n") }
            continue
        }

        if layerMeta.show {
            castPacket(led, series, layerMeta.color, reverse)
        }
    }
}

func reverseLeds(led []uint32) {
    for i, j := 0, len(led)-1; i < j; i, j = i+1, j-1 {
        led[i], led[j] = led[j], led[i]
    }
}

func initLeds(led []uint32) {
    for i, _ := range led {
        led[i] = 0
    }
}

func castPacket(led []uint32, k int, color uint32,reverse bool) {
    for i := -(k-1); i < len(led)+series+speed; i += speed {
        initLeds(led)

        for j := 0; j < k; j++ {
            if t := i + j; 0 <= t && t < len(led) {
                // packet color gradiation
                g := (((color & 0xFF0000) >> 16) * uint32(j+1) / uint32(k)) << 16
                r := (((color & 0x00FF00) >> 8)* uint32(j+1) / uint32(k)) << 8
                b := (color & 0x0000FF)* uint32(j+1) / uint32(k)
                led[t] = g|r|b
            }
        }

        if reverse {
            reverseLeds(led)
        }

        setLeds(led)
        err := ws2811.Render()
        if err != nil {
            ws2811.Clear()
            fmt.Println("Error during wipe " + err.Error())
            os.Exit(-1)
        }
    }
}

func setLeds(led []uint32) {
	for i := 0; i < count; i++ {
		ws2811.SetLed(i, led[i])
	}
}

func isAnomaly(packet gopacket.Packet) bool {
    anml := false
    if tcp := packet.Layer(layers.LayerTypeTCP); tcp != nil {
        tcpl, _ := tcp.(*layers.TCP)
        // Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
        if tcpl.FIN && tcpl.URG && tcpl.PSH {
            anml = true
        }
    }
    return anml
}

func externalIP() (string, string, error) {
    ifaces, err := net.Interfaces()
	if err != nil {
		return "", "", err
	}
    var ipv4Addr net.IP
    var ipv6Addr net.IP
	for _, iface := range ifaces {
        if iface.Name != device {
            continue // select device "eth0" or something
        }
		addrs, err := iface.Addrs()
		if err != nil {
			return "", "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ipv4 := ip.To4(); ipv4 != nil {
				ipv4Addr = ipv4
			}else if ipv6 := ip.To16(); ipv6 != nil {
				ipv6Addr = ipv6
			}
		}
    }
    if ipv4Addr == nil{
        return "", "", errors.New("are you connected to the network?")
    }else{
        return ipv4Addr.String(), ipv6Addr.String(), nil
    }
}

func showIPAddress(led []uint32, ipaddr string) {
	ip := net.ParseIP(ipaddr)
	ipv4 := ip.To4()

    initLeds(led)
	for i := 0; i < ipv4Len; i++ {

		// number
		for j := 0; j < 8; j++ {
			t := i * 9 + j
	        if (ipv4[i]>>uint(7-j))&1 == 1 { // nth bit of Y = (X>>n)&1;
                led[t] = colors[10]
            }
	    }

		// period
		led[(i+1) * 9 - 1] = colors[1]
	}
    setLeds(led)
    err := ws2811.Render()
    if err != nil {
        ws2811.Clear()
        fmt.Println("Error during wipe " + err.Error())
        os.Exit(-1)
    }
}

func resetLeds(led []uint32) {
    initLeds(led)
    setLeds(led)
    err := ws2811.Render()
    if err != nil {
        ws2811.Clear()
        fmt.Println("Error during wipe " + err.Error())
        os.Exit(-1)
    }
}

func categorizePacket(packet gopacket.Packet) string {
    packetName := "Others";
    if isAnomaly(packet) {
        packetName = "Anomaly"
    }else if lldp := packet.Layer(layers.LayerTypeLinkLayerDiscovery); lldp != nil {
        packetName = "LLDP"
    }else if dns := packet.Layer(layers.LayerTypeDNS); dns != nil {
        packetName = "DNS"
    }else if icmpv4 := packet.Layer(layers.LayerTypeICMPv4); icmpv4 != nil {
        packetName = "ICMP"
    }else if icmpv6 := packet.Layer(layers.LayerTypeICMPv6); icmpv6 != nil {
        packetName = "ICMP"
    }else if dhcpv4 := packet.Layer(layers.LayerTypeDHCPv4); dhcpv4 != nil {
        packetName = "DHCP"
    }else if arp := packet.Layer(layers.LayerTypeARP); arp != nil {
        packetName = "ARP"
    }else if igmp := packet.Layer(layers.LayerTypeIGMP); igmp != nil {
        packetName = "IGMP"
    }else if udp := packet.Layer(layers.LayerTypeUDP); udp != nil {
        packetName = "UDP"
    }else if tcp := packet.Layer(layers.LayerTypeTCP); tcp != nil {
        packetName = "TCP"
    }
    return packetName
}
