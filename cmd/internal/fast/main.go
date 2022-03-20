package fast

import (
	"context"
	fast2 "framey/assignment/pkg/fast"
	"log"
)

func Main(args []string) {
	err := flagSet.Parse(args[1:])
	if err != nil {
		panic(err)
	}

	var client fast2.Client

	ctx, cancel := context.WithTimeout(context.Background(), *cfgTime)
	defer cancel()

	m, err := fast2.GetManifest(ctx, *urlCount)
	if err != nil {
		log.Fatalf("Error loading fast.com configuration: %v", err)
	}

	download(m, &client)
	upload(m, &client)
}
