package speedtest

import (
	"context"
	"fmt"
	speedtest2 "framey/assignment/pkg/speedtest"
	"log"
)

// Loads the list of servers and exits the program on failure.
//
func listServers(
	ctx context.Context,
	client *speedtest2.Client,
) []speedtest2.Server {
	servers, err := client.LoadAllServers(ctx)
	if err != nil {
		log.Fatalf("Failed to load server list: %v\n", err)
	}
	if len(servers) == 0 {
		log.Fatalf("No servers found somehow...")
	}
	if len(srvBlk) != 0 {
		servers = pruneBlockedServers(servers)
	}
	return servers
}

func pruneBlockedServers(servers []speedtest2.Server) []speedtest2.Server {
	n := make([]speedtest2.Server, len(servers)-len(srvBlk))[:0]
	for _, s := range servers {
		var i bool
		for _, b := range srvBlk {
			if s.ID == b {
				i = true
			}
		}
		if !i {
			n = append(n, s)
		}
	}
	return n
}

// Iterates through the list of server and prints them out.
//
func printServers(client *speedtest2.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), *cfgTime)
	defer cancel()

	for _, s := range listServers(ctx, client) {
		fmt.Println(s)
	}
}
