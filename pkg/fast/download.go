package fast

import (
	"context"
	"framey/assignment/internal/prober"
	"framey/assignment/internal/prober/proberutil"
	"framey/assignment/internal/units"
	"io"
)

const (
	concurrentDownloadLimit = 12
	downloadBufferSize      = 4096
	downloadRepeats         = 5
)

var downloadSizes = []int{
	256, 1024, 4096,
	131_072, 1_048_576, 8_388_608, 16_777_216,
	33_554_432}

// ProbeDownloadSpeed Will probe download speed until enough samples are taken or ctx expires.
func (m *Manifest) ProbeDownloadSpeed(
	ctx context.Context,
	client *Client,
	stream chan<- units.BytesPerSecond,
) (units.BytesPerSecond, error) {
	grp := prober.NewGroup(concurrentDownloadLimit)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, size := range downloadSizes {
		for i := 0; i < downloadRepeats; i++ {
			for _, t := range m.m.Targets {
				url := putSizeIntoURL(t.URL, size)
				grp.Add(func() (prober.BytesTransferred, error) {
					return client.downloadFile(ctx, url)
				})
			}
		}
	}

	return proberutil.SpeedCollect(grp, stream)
}

func (c *Client) downloadFile(
	ctx context.Context,
	url string,
) (t prober.BytesTransferred, err error) {
	// Check early failure where context is already canceled.
	if err = ctx.Err(); err != nil {
		return
	}

	res, err := c.get(ctx, url)
	if err != nil {
		return t, err
	}
	defer res.Body.Close()

	var buf [downloadBufferSize]byte
	for {
		read, err := res.Body.Read(buf[:])
		t += prober.BytesTransferred(read)
		if err != nil {
			if err != io.EOF {
				return t, err
			}
			break
		}
	}
	return t, nil
}
