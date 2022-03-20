package fast

import (
	"context"
	"fmt"
	"framey/assignment/fast"
	"framey/assignment/oututil"
	"framey/assignment/units"
	"log"

	"golang.org/x/sync/errgroup"
)

func download(m *fast.Manifest, client *fast.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), *dlTime)
	defer cancel()

	stream, finalize := proberPrinter(func(s units.BytesPerSecond) string {
		return formatSpeed("Download speed", s)
	})
	speed, err := m.ProbeDownloadSpeed(ctx, client, stream)
	if err != nil {
		log.Fatalf("Error probing download speed: %v", err)
		return
	}
	finalize(speed)
}

func upload(m *fast.Manifest, client *fast.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), *ulTime)
	defer cancel()

	stream, finalize := proberPrinter(func(s units.BytesPerSecond) string {
		return formatSpeed("Upload speed", s)
	})
	speed, err := m.ProbeUploadSpeed(ctx, client, stream)
	if err != nil {
		log.Fatalf("Error probing upload speed: %v", err)
		return
	}
	finalize(speed)
}

func proberPrinter(format func(units.BytesPerSecond) string) (
	stream chan units.BytesPerSecond,
	finalize func(units.BytesPerSecond),
) {
	p := oututil.StartPrinting()
	p.Println(format(units.BytesPerSecond(0)))

	stream = make(chan units.BytesPerSecond)
	var g errgroup.Group
	g.Go(func() error {
		for speed := range stream {
			p.Println(format(speed))
		}
		return nil
	})

	finalize = func(s units.BytesPerSecond) {
		g.Wait()
		p.Finalize(format(s))
	}
	return
}

func formatSpeed(prefix string, s units.BytesPerSecond) string {
	var i interface{}
	// Default return speed is in bytes.
	if *fmtBytes {
		i = s
	} else {
		i = s.BitsPerSecond()
	}
	return fmt.Sprintf("%s: %v", prefix, i)
}
