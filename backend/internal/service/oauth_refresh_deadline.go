package service

import (
	"context"
	"time"
)

// oauthRefreshContextErr checks the budget even when the deadline's cancellation
// callback has not run yet. It gates entry to persistence, not atomic DB commit.
func oauthRefreshContextErr(ctx context.Context, now func() time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if deadline, ok := ctx.Deadline(); ok {
		if now == nil {
			now = time.Now
		}
		if !now().Before(deadline) {
			return context.DeadlineExceeded
		}
	}
	return nil
}
