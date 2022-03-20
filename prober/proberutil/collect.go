package proberutil

import (
	"framey/assignment/prober"
	"framey/assignment/units"
	"time"
)

func SpeedCollect(
	grp *prober.Group,
	stream chan<- units.BytesPerSecond,
) (units.BytesPerSecond, error) {
	start := time.Now()

	if stream != nil {
		inc := grp.GetIncremental()
		go func() {
			for b := range inc {
				d := float64(time.Since(start)) / float64(time.Second)
				stream <- units.BytesPerSecond(float64(b) / d)
			}
			close(stream)
		}()
	}

	b, err := grp.Collect()
	if err != nil {
		return units.BytesPerSecond(0), err
	} else {
		d := float64(time.Since(start)) / float64(time.Second)
		return units.BytesPerSecond(float64(b) / d), nil
	}
}
