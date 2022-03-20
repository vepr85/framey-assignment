package fast

import (
	"context"
	"framey/assignment/fast"
	"log"
)

func Main(args []string) {
	err := flagSet.Parse(args[1:])
	if err != nil {
		panic(err)
	}

	var client fast.Client

	ctx, cancel := context.WithTimeout(context.Background(), *cfgTime)
	defer cancel()

	m, err := fast.GetManifest(ctx, *urlCount)
	if err != nil {
		log.Fatalf("Error loading fast.com configuration: %v", err)
	}

	download(m, &client)
	upload(m, &client)
}
