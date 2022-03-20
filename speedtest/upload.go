package speedtest

import (
	"context"
	"crypto/rand"
	"fmt"
	"framey/assignment/prober"
	"framey/assignment/prober/proberutil"
	"framey/assignment/units"
	"io"
	"strings"
)

const (
	concurrentUploadLimit = concurrentDownloadLimit
	uploadRepeats         = downloadRepeats * 25

	safeChars = "0123456789abcdefghijklmnopqrstuv"
)

var uploadSizes = []int{1000 * 1000 / 4, 1000 * 1000 / 2}

// Will probe upload speed until enough samples are taken or ctx expires.
func (s Server) ProbeUploadSpeed(
	ctx context.Context,
	client *Client,
	stream chan<- units.BytesPerSecond,
) (units.BytesPerSecond, error) {
	grp := prober.NewGroup(concurrentUploadLimit)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := range uploadSizes {
		for j := 0; j < uploadRepeats; j++ {
			size := uploadSizes[i]
			grp.Add(func() (t prober.BytesTransferred, err error) {
				err = client.uploadFile(ctx, s.URL, size)
				if err == nil {
					t = prober.BytesTransferred(size)
				}
				return
			})
		}
	}

	return proberutil.SpeedCollect(grp, stream)
}

type safeReader struct {
	in io.Reader
}

func (r safeReader) Read(p []byte) (n int, err error) {
	n, err = r.in.Read(p)
	for i := 0; i < n; i++ {
		p[i] = safeChars[int(p[i])%len(safeChars)]
	}
	return n, err
}

func (c *Client) uploadFile(
	ctx context.Context,
	url string,
	size int,
) error {
	res, err := c.post(ctx, url, "application/x-www-form-urlencoded",
		io.MultiReader(
			strings.NewReader("content1="),
			io.LimitReader(&safeReader{rand.Reader}, int64(size-9))))
	if err != nil {
		return fmt.Errorf("upload to %q failed: %v", url, err)
	}
	defer res.Body.Close()

	return nil
}
