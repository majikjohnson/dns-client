package main

import (
	"fmt"
	"net"
	"os"

	"github.com/miekg/dns"
)

func do53(q string, qt uint16) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(q), qt)
	m.RecursionDesired = true

	r, _, _ := c.Exchange(m, "1.1.1.1:53")

	if r == nil {
		fmt.Println("r is nil")
	}

	fmt.Println(r)
}

func dnsRawUDP(args []string) {
	fmt.Println("Making a DNS request using raw UDP")
	if len(args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: client.exe dns_host:port domain_to_resolve\n")
		os.Exit(1)
	}

	service := args[1]
	domain := args[2]
	qtype := args[3]

	//qtype, err := strconv.ParseUint(args[3], 10, 16)

	//checkError(err)

	//Resolve the hostname of the DNS server.
	udpAddr, err := net.ResolveUDPAddr("udp", service)

	checkError(err)

	conn, err := net.DialUDP("udp", nil, udpAddr)

	query, err := toWire(domain, qtype)
	checkError(err)

	_, err = conn.Write(query)
	checkError(err)

	var buf [1024]byte
	n, err := conn.Read(buf[0:])
	checkError(err)

	r := fromWire(buf[0:n])

	fmt.Println(r)

	os.Exit(0)

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func toWire(q string, qt string) ([]byte, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("test.com"), dns.TypeA)
	m.RecursionDesired = true
	return m.Pack()
}

func fromWire(buf []byte) string {
	m := new(dns.Msg)
	if err := m.Unpack(buf); err != nil {
		checkError(err)
	}
	return m.String()
}

// func fromWire(resp []byte) (string, error) {

// }

func main() {
	//do53("testingninja.co.uk", dns.TypeA)
	dnsRawUDP(os.Args)
}
