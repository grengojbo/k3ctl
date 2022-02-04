package util

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/cloverstd/tcping/ping"
)

type PingArgs struct {
	Host string
	Port int
	Protocol string
	ShowResponse bool
	Timeout int
	Counter int
	Retry int
}

func PingRetry(args *PingArgs) (err error) {
	if args.Timeout == 0 {
		args.Timeout = 1
	}
	if args.Counter == 0 {
		args.Counter = 1
	}
	if args.Retry == 0 {
		args.Retry = 1
	}

	for i := 0; i < args.Retry; i++ {
		err = runPing(args)
		if err == nil {
			return nil
		}
		args.Timeout = args.Timeout * 2
	}
	return err
}

func runPing(args *PingArgs) (err error) {
	var pinger ping.Pinger
	var protocol ping.Protocol
	var timeoutDuration time.Duration
	var intervalDuration time.Duration
	
	httpMethod := "HEAD"
	switch args.Protocol {
	case "tcp":
		pinger = ping.NewTCPing()
		protocol = ping.TCP
	case "http":
		pinger = ping.NewHTTPing(httpMethod)
		protocol = ping.HTTP
	case "https":
		pinger = ping.NewHTTPing(httpMethod)
		protocol = ping.HTTPS
	default:
		pinger = ping.NewTCPing()
		protocol = ping.TCP
	}

	timeoutDuration = time.Duration(args.Timeout) * time.Second
	intervalDuration = time.Duration(args.Timeout) * time.Second

	target := ping.Target{
		Timeout:  timeoutDuration,
		Interval: intervalDuration,
		Host:     args.Host,
		Port:     args.Port,
		Counter:  args.Counter,
		// Proxy:    proxy,
		Protocol: protocol,
	}
	
	pinger.SetTarget(&target)
	pingerDone := pinger.Start()
	pinger.Start()
	// TODO: сделать свою реализацию чтобы неотображать результат пинга
	select {
	case <-pingerDone:
		break
	// case <-sigs:
	// 	break
	}
		
	result := pinger.Result()
	if result.SuccessCounter > 1 {
		return nil
	}

	return fmt.Errorf(result.String())
}

// PingTCP ping tcp,http, https timeot in second, ping counter
// command line util https://github.com/cloverstd/tcping
func PingTCP(host string, timeout int, counter int) (err error) {
	var protocol ping.Protocol
	var timeoutDuration time.Duration
	var intervalDuration time.Duration
	// var (
	// 	// err    error
	// 	port   int
	// 	schema string
	// )
	// if len(args) == 2 {
	// 	port, err = strconv.Atoi(args[1])
	// 	if err != nil {
	// 		fmt.Println("port should be integer")
	// 		cmd.Usage()
	// 		return
	// 	}
		// schema = ping.TCP.String()
	// } else {
		// var matched bool
		errMess := fmt.Errorf("[Ping] not a valid uri: %v", host)
		schema, host, port, matched := ping.CheckURI(host)
		if !matched {
			// fmt.Println("not a valid uri")
			return errMess
		}
	// }
		timeoutDuration = time.Duration(timeout) * time.Second
	// if res, err := strconv.Atoi(timeout); err == nil {
		// timeoutDuration = time.Duration(res) * time.Microsecond
	// } else {
	// 	timeoutDuration, err = time.ParseDuration(timeout)
	// 	if err != nil {
	// 		// fmt.Println("parse timeout failed", err)
	// 		// cmd.Usage()
	// 		return fmt.Errorf("[Ping] %v", err.Error())
	// 	}
	// }

	// if res, err := strconv.Atoi(interval); err == nil {
		intervalDuration = time.Duration(timeout*2) * time.Second
	// } else {
	// 	intervalDuration, err = time.ParseDuration(interval)
	// 	if err != nil {
	// 		return fmt.Errorf("[Ping] parse interval failed: %v", err.Error())
	// 	}
	// }
	// if httpMode {
		protocol = ping.TCP
	// // } else {
	// 	protocol, err = ping.NewProtocol(schema)
	// 	if err != nil {
	// 		return fmt.Errorf("[Ping] %v", err.Error())
	// 	}
	// // }
	// if len(dnsServer) != 0 {
	// 	ping.UseCustomeDNS(dnsServer)
	// }

	parseHost, _ := FormatIP(host)
	target := ping.Target{
		Timeout:  timeoutDuration,
		Interval: intervalDuration,
		Host:     parseHost,
		Port:     port,
		Counter:  counter,
		// Proxy:    proxy,
		Protocol: protocol,
	}

	var pinger ping.Pinger
		switch protocol {
		case ping.TCP:
			pinger = ping.NewTCPing()
		case ping.HTTP, ping.HTTPS:
			// var httpMethod string
			// switch {
			// case httpHead:
			// 	httpMethod = "HEAD"
			// case httpPost:
			// 	httpMethod = "POST"
			// default:
			// 	httpMethod = "GET"
			// }
			httpMethod := "HEAD"
			pinger = ping.NewHTTPing(httpMethod)
		default:
			// fmt.Printf("schema: %s not support\n", schema)
			// cmd.Usage()
			return fmt.Errorf("[Ping] schema: %s not support", schema)
		}
		
		pinger.SetTarget(&target)
		pingerDone := pinger.Start()
		pinger.Start()
		// TODO: сделать свою реализацию чтобы неотображать результат пинга
		select {
		case <-pingerDone:
			break
		// case <-sigs:
		// 	break
		}
		
	result := pinger.Result()
	if result.SuccessCounter > 1 {
		return nil
	}

	return fmt.Errorf(result.String())
}

// FormatIP - trim spaces and format IP
//
// IP - the provided IP
//
// string - return "" if the input is neither valid IPv4 nor valid IPv6
//          return IPv4 in format like "192.168.9.1"
//          return IPv6 in format like "[2002:ac1f:91c5:1::bd59]"
func FormatIP(IP string) (string, error) {

	host := strings.Trim(IP, "[ ]")
	if parseIP := net.ParseIP(host); parseIP != nil {
		// valid ip
		if parseIP.To4() == nil {
			// ipv6
			host = fmt.Sprintf("[%s]", host)
		}
		return host, nil
	}
	return "", fmt.Errorf("Error IP format")
}