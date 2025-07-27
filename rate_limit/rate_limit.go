package rate_limit

// the reason I am having this interface is because when we do testing
// we could mock the ratelimit service
// and for rate limit service we could also have different DAO implementation if we decide not to use REDIS
type IRateLimiter interface {
	Allow(userid string) (bool, error)
}
