package rate_limit

type IRateLimiter interface {
	Allow(userid string) (bool, error)
}
