package main

import (
	"github.com/hashicorp/mdns"
	"github.com/miekg/dns"
	"log"
	"os"
	"regexp"
	"time"
)

type Proxy struct{}

var (
	pattern     string = ".*\\.local"
	PROTO       string = "tcp"
	LOG         *log.Logger
	LOCAL_REGEX *regexp.Regexp
)

func (p Proxy) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	c := new(dns.Client)

	c.Net = PROTO
	for _, q := range r.Question {
		if LOCAL_REGEX.MatchString(q.Name) && q.Qtype == dns.TypeA && q.Qclass == dns.ClassINET {
			LOG.Printf("Local regex match: %s", q.Name)
			p.mDNSLookup(w, r, q.Name)
			break
		}
	}
	dns.HandleFailed(w, r)
}

func (p Proxy) mDNSLookup(w dns.ResponseWriter, r *dns.Msg, name string) {
	client, err := mdns.NewClient()
	if err != nil {
		LOG.Fatalf("Failed to create mdns client")
	}
	defer client.Close()

	// Start listening for response packets
	msgCh := make(chan *dns.Msg, 32)
	defer close(msgCh)

	go client.Recv(client.Ipv4UnicastConn, msgCh)
	go client.Recv(client.Ipv6UnicastConn, msgCh)
	go client.Recv(client.Ipv4MulticastConn, msgCh)
	go client.Recv(client.Ipv6MulticastConn, msgCh)

	finish := time.After(time.Duration(3) * time.Second)

	// Start the lookup
	if err := client.SendQuery(r); err != nil {
		LOG.Fatalf("Failed to send msg via client")
	}
	for {
		select {
		case resp := <-msgCh:
			for _, answer := range append(resp.Answer, resp.Extra...) {
				LOG.Printf("ans:", answer)
				switch rr := answer.(type) {
				case *dns.A:
					if rr.Hdr.Name == name {
						w.WriteMsg(resp)
						return
					}
				}
			}
		case <-finish:
			return
		}
	}
}

func main() {
	LOG = log.New(os.Stderr, "[DNS PROXY] ", log.LstdFlags)

	if re, err := regexp.Compile(pattern); err != nil {
		LOG.Fatalf("Compiling pattern [%s] was %s", pattern, err)
	} else {
		LOCAL_REGEX = re
	}
	proxyer := Proxy{}

	if err := dns.ListenAndServe("0.0.0.0:53", "udp", proxyer); err != nil {
		LOG.Fatal(err)
	}
}
