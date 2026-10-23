//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestTokenRefreshService_DeadlinePersistenceBoundary(t *testing.T) {
	for _, platform := range []string{PlatformGrok, PlatformGemini} {
		for _, unified := range []bool{false, true} {
			path := "compatibility"
			if unified {
				path = "unified"
			}
			for _, scenario := range []string{"before", "at", "after", "parent_deadline", "parent_cancel", "late_invalid_grant"} {
				t.Run(platform+"/"+path+"/"+scenario, func(t *testing.T) {
					account := grokPoolAccount(151)
					account.Platform = platform
					repo := &tokenRefreshAccountRepo{snapshotReads: true}
					repo.accountsByID = map[int64]*Account{account.ID: &account}
					scheduler := &tokenRefreshSchedulerCache{}
					now := time.Now()
					clock := func() time.Time { return now }
					svc := &TokenRefreshService{
						accountRepo: repo, refreshPolicy: DefaultBackgroundRefreshPolicy(),
						cfg:            &config.TokenRefreshConfig{MaxRetries: 1},
						schedulerCache: scheduler, deadlineNow: clock,
					}
					if unified {
						svc.refreshAPI = NewOAuthRefreshAPI(repo, nil)
						svc.refreshAPI.deadlineNow = clock
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
					defer cancel()
					calls := 0
					executor := &dynamicRefreshExecutor{
						canRefresh: true, cacheKey: "deadline-boundary",
						needsRefreshFunc: func() bool { return true },
						refreshFunc: func(attemptCtx context.Context, _ *Account) (map[string]any, error) {
							calls++
							deadline, ok := attemptCtx.Deadline()
							require.True(t, ok)
							now = deadline.Add(time.Nanosecond)
							switch scenario {
							case "before":
								now = deadline.Add(-time.Nanosecond)
							case "at":
								now = deadline
							case "parent_deadline":
								now, ok = ctx.Deadline()
								require.True(t, ok)
							case "parent_cancel":
								cancel()
							}
							if scenario != "parent_cancel" {
								require.NoError(t, attemptCtx.Err())
								select {
								case <-attemptCtx.Done():
									t.Fatal("deadline notification must still be pending")
								default:
								}
							}
							if scenario == "late_invalid_grant" {
								return nil, errors.New("invalid_grant")
							}
							return map[string]any{"access_token": "new-access", "refresh_token": "new-refresh"}, nil
						},
					}

					err := svc.refreshWithRetry(ctx, snapshotOAuthRefreshAccount(&account), executor, executor, time.Hour)

					require.Equal(t, 1, calls)
					require.Zero(t, repo.setErrorCalls)
					var timeoutErr *refreshAttemptTimeoutError
					if scenario == "before" {
						require.NoError(t, err)
						require.Equal(t, "new-refresh", repo.accountsByID[account.ID].GetCredential("refresh_token"))
						require.Equal(t, 1, scheduler.setAccountCalls)
						require.Zero(t, repo.setTempUnschedCalls)
						return
					}
					require.Zero(t, repo.updateCredentialsCalls)
					require.Zero(t, repo.conditionalSuccessCalls)
					require.Zero(t, scheduler.setAccountCalls)
					require.Equal(t, "refresh-token", repo.accountsByID[account.ID].GetCredential("refresh_token"))
					switch scenario {
					case "parent_deadline", "parent_cancel":
						expected := context.DeadlineExceeded
						if scenario == "parent_cancel" {
							expected = context.Canceled
						}
						require.ErrorIs(t, err, expected)
						require.False(t, errors.As(err, &timeoutErr))
						require.Zero(t, repo.setTempUnschedCalls)
						require.Zero(t, repo.conditionalTempCalls)
					default:
						require.ErrorAs(t, err, &timeoutErr)
						require.Equal(t, 1, repo.setTempUnschedCalls)
					}
				})
			}
		}
	}
}

type deadlineCASFailureRepo struct {
	*tokenRefreshAccountRepo
	beforeCAS func()
}

func (r *deadlineCASFailureRepo) UpdateGrokOAuthCredentialsIfUnchanged(ctx context.Context, id int64, expected map[string]any, proxyID *int64, credentials map[string]any) (bool, error) {
	r.beforeCAS()
	return r.tokenRefreshAccountRepo.UpdateGrokOAuthCredentialsIfUnchanged(ctx, id, expected, proxyID, credentials)
}

func TestTokenRefreshService_DeadlineDoesNotReclassifyUncertainPersistence(t *testing.T) {
	account := grokPoolAccount(152)
	now := time.Now()
	clock := func() time.Time { return now }
	var attemptDeadline time.Time
	repo := &deadlineCASFailureRepo{
		tokenRefreshAccountRepo: &tokenRefreshAccountRepo{snapshotReads: true, conditionalSuccessErr: context.DeadlineExceeded},
		beforeCAS:               func() { now = attemptDeadline.Add(time.Nanosecond) },
	}
	repo.accountsByID = map[int64]*Account{account.ID: &account}
	api := NewOAuthRefreshAPI(repo, nil)
	api.deadlineNow = clock
	svc := &TokenRefreshService{
		accountRepo: repo, refreshAPI: api, refreshPolicy: DefaultBackgroundRefreshPolicy(),
		cfg: &config.TokenRefreshConfig{MaxRetries: 3}, deadlineNow: clock,
	}
	calls := 0
	executor := &dynamicRefreshExecutor{
		canRefresh: true, cacheKey: "deadline-cas-failure", needsRefreshFunc: func() bool { return true },
		refreshFunc: func(ctx context.Context, _ *Account) (map[string]any, error) {
			calls++
			var ok bool
			attemptDeadline, ok = ctx.Deadline()
			require.True(t, ok)
			return map[string]any{"access_token": "new-access", "refresh_token": "new-refresh"}, nil
		},
	}

	err := svc.refreshWithRetry(context.Background(), &account, executor, executor, time.Hour)

	var containmentErr *providerCycleContainmentRefreshError
	var timeoutErr *refreshAttemptTimeoutError
	require.ErrorAs(t, err, &containmentErr)
	require.False(t, errors.As(err, &timeoutErr), "an uncertain commit must not become a retryable timeout")
	require.Equal(t, 1, calls)
	require.Equal(t, 1, repo.conditionalSuccessCalls)
	require.Zero(t, repo.setErrorCalls)
	require.Zero(t, repo.setTempUnschedCalls)
}
