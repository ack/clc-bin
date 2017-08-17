package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	clc "github.com/CenturyLinkCloud/clc-sdk"
	"github.com/CenturyLinkCloud/clc-sdk/api"
	"github.com/CenturyLinkCloud/clc-sdk/server"
	"github.com/CenturyLinkCloud/clc-sdk/status"
)

// Version of binary
const Version = "0.1"

func supportedProtocol(proto string) bool {
	proto = strings.ToLower(proto)
	switch proto {
	case
		"tcp",
		"udp",
		"icmp":
		return true
	}
	return false
}

func main() {
	un := flag.String("username", "", "clc username")
	pw := flag.String("password", "", "clc password")
	sid := flag.String("server", "", "server id")
	spec := flag.String("ports", "", "ports to open - eg. '80/tcp 22/tcp'")
	intl := flag.String("internal", "", "(optional) internal ip")
	output := flag.Bool("output", false, "echo provisioned IP")
	verbose:= flag.Bool("verbose", false, "debug logging")
	if *verbose {
		os.Setenv("DEBUG", "on") // turns on wire tracing
	}
	flag.Parse()
	var client *clc.Client

	if *spec == "" {
		log.Panic("missing flag -ports")
	}
	// when not passed, use local hostname
	if *sid == "" {
		*sid, _ = os.Hostname()
	}
	if *intl != "" && *verbose{
		log.Printf("Allocating on internal IP: %v", *intl)
	}


	token:= os.Getenv("CLC_V2_API_TOKEN")
	alias := os.Getenv("CLC_ACCT_ALIAS")
	if token != "" && alias != "" {
		client = clc.NewFromAliasToken(alias, token)
	} else if (*un != "" && *pw != "") { 
		config, _ := api.NewConfig(*un, *pw)
		config.UserAgent = fmt.Sprintf("natip: %s", Version)
		client = clc.New(config)
		if err := client.Authenticate(); err != nil {
			log.Panicf("Failed to auth: %v", err)
		}

	} else {
		log.Panic("token/alias not available and user/pass not provided")
	}

	pubip := server.PublicIP{}
	var ports []server.Port
	for _, s := range strings.Split(*spec, " ") {
		x := strings.Split(s, "/")
		portrange, proto := x[0], x[1]
		if !supportedProtocol(proto) {
			log.Panicf("Unsupported protocol: %v", proto)
		}
		fromto := strings.Split(portrange, "-")
		var i, j int
		if len(fromto) > 1 {
			i, _ = strconv.Atoi(fromto[0])
			j, _ = strconv.Atoi(fromto[1])
		} else {
			i, _ = strconv.Atoi(fromto[0])
			j = -1
		}
		port := server.Port{
			Port:     i,
			Protocol: proto,
		}
		if j != -1 {
			port.PortTo = j
		}
		ports = append(ports, port)
	}
	pubip.Ports = ports

	var addr string
	var st *status.Status
	var svr *server.Response
	svr, err := client.Server.Get(*sid)
	if err != nil {
		log.Panicf("Failed fetching server: %v - %v", *sid, err)
	}

	for _, ip := range svr.Details.IPaddresses {
		addr = ip.Public
		pubip.InternalIP = ip.Internal
		if *intl == ip.Internal {
			// specific NIC requested
			break
		}
	}

	if addr != "" {
		if *verbose {
			log.Printf("updating existing public ip %v on server %v", addr, *sid)
		}
		st, err = client.Server.UpdatePublicIP(*sid, addr, pubip)
	} else {
		if *verbose {
			log.Printf("creating public ip on %v. internal: %v", *sid, pubip.InternalIP)
		}
		st, err = client.Server.AddPublicIP(*sid, pubip)
	}
	if err != nil {
		log.Panicf("error sending public ip: %v", err)
	}
	if *verbose {
		log.Printf("polling status on %v", st.ID)
	}
	poll := make(chan *status.Response, 1)
	_ = client.Status.Poll(st.ID, poll)
	status := <-poll
	if *verbose {
		log.Printf("done. status: %v", status)
	}

	// fetch/print public ips
	svr, _ = client.Server.Get(*sid)
	pub := ""
	for _, ip := range svr.Details.IPaddresses {
		if ip.Public != "" {
			pub = ip.Public
		}
		if *verbose {
			log.Printf("IP: %v \t => %v", ip.Internal, ip.Public)
		}
	}
	if *output {
		fmt.Printf(pub)
	}
}
