package magicbytes

import "context"

func checkContextIsAlive(ctx context.Context, inner_ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case <-inner_ctx.Done():
		return false
	default:
	}

	return true
}
