//go:generate sh api.sh

package api

import (
	"context"
)

func (f *faas) Push(ctx context.Context, msg *PushRequest) (*Response, error) {

	return okResponse(), nil
}
