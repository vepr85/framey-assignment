package fast

import (
	"context"
	internal2 "framey/assignment/pkg/fast/internal"
)

type Manifest struct {
	m *internal2.Manifest
}

func GetManifest(ctx context.Context, urls int) (*Manifest, error) {
	tok, err := internal2.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	mi, err := internal2.GetManifest(ctx, tok, urls)
	if err != nil {
		return nil, err
	}
	return &Manifest{mi}, nil
}
