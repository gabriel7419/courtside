package nba

import (
	"sync"
	"time"

	"github.com/gabriel7419/courtside/internal/api"
)

// CacheConfig holds TTL and size settings for the response cache.
type CacheConfig struct {
	MatchesTTL      time.Duration
	MatchDetailsTTL time.Duration
	LiveMatchesTTL  time.Duration
	MaxMatchesCache int
	MaxDetailsCache int
}

// DefaultCacheConfig returns sensible defaults for the NBA client.
// Live game data changes frequently; finished game data is permanent.
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		MatchesTTL:      30 * time.Second, // scoreboard updates frequently
		MatchDetailsTTL: 10 * time.Second, // live box score data
		LiveMatchesTTL:  10 * time.Second, // live game list
		MaxMatchesCache: 10,
		MaxDetailsCache: 50,
	}
}

type cachedMatches struct {
	matches   []api.Match
	expiresAt time.Time
}

type cachedDetails struct {
	details   *api.MatchDetails
	expiresAt time.Time
}

// ResponseCache provides thread-safe caching for NBA API responses.
type ResponseCache struct {
	config       CacheConfig
	matchesMu    sync.RWMutex
	matchesCache map[string]cachedMatches // key: "YYYY-MM-DD"
	detailsMu    sync.RWMutex
	detailsCache map[int]cachedDetails // key: gameID
	liveMu       sync.RWMutex
	liveCache    *cachedMatches
}

// NewResponseCache creates a new cache with the given configuration.
func NewResponseCache(config CacheConfig) *ResponseCache {
	return &ResponseCache{
		config:       config,
		matchesCache: make(map[string]cachedMatches),
		detailsCache: make(map[int]cachedDetails),
	}
}

// Matches retrieves cached games for a date key, or nil if expired/absent.
func (c *ResponseCache) Matches(dateKey string) []api.Match {
	c.matchesMu.RLock()
	defer c.matchesMu.RUnlock()
	cached, ok := c.matchesCache[dateKey]
	if !ok || time.Now().After(cached.expiresAt) {
		return nil
	}
	return cached.matches
}

// SetMatches stores games in cache with TTL.
func (c *ResponseCache) SetMatches(dateKey string, matches []api.Match) {
	c.matchesMu.Lock()
	defer c.matchesMu.Unlock()
	if len(c.matchesCache) >= c.config.MaxMatchesCache {
		c.evictOldestMatches()
	}
	c.matchesCache[dateKey] = cachedMatches{
		matches:   matches,
		expiresAt: time.Now().Add(c.config.MatchesTTL),
	}
}

// Details retrieves cached game details, or nil if expired/absent.
func (c *ResponseCache) Details(gameID int) *api.MatchDetails {
	c.detailsMu.RLock()
	defer c.detailsMu.RUnlock()
	cached, ok := c.detailsCache[gameID]
	if !ok || time.Now().After(cached.expiresAt) {
		return nil
	}
	return cached.details
}

// SetDetails stores game details in cache. Finished games are cached permanently.
func (c *ResponseCache) SetDetails(gameID int, details *api.MatchDetails) {
	c.detailsMu.Lock()
	defer c.detailsMu.Unlock()
	if len(c.detailsCache) >= c.config.MaxDetailsCache {
		c.evictOldestDetails()
	}
	ttl := c.config.MatchDetailsTTL
	if details != nil && details.Status == api.MatchStatusFinished {
		ttl = 24 * time.Hour // finished games never change
	}
	c.detailsCache[gameID] = cachedDetails{
		details:   details,
		expiresAt: time.Now().Add(ttl),
	}
}

// ClearDetails removes a specific game from the details cache (force refresh).
func (c *ResponseCache) ClearDetails(gameID int) {
	c.detailsMu.Lock()
	defer c.detailsMu.Unlock()
	delete(c.detailsCache, gameID)
}

// LiveMatches retrieves cached live games, or nil if expired/absent.
func (c *ResponseCache) LiveMatches() []api.Match {
	c.liveMu.RLock()
	defer c.liveMu.RUnlock()
	if c.liveCache == nil || time.Now().After(c.liveCache.expiresAt) {
		return nil
	}
	return c.liveCache.matches
}

// SetLiveMatches stores live games in cache with TTL.
func (c *ResponseCache) SetLiveMatches(matches []api.Match) {
	c.liveMu.Lock()
	defer c.liveMu.Unlock()
	c.liveCache = &cachedMatches{
		matches:   matches,
		expiresAt: time.Now().Add(c.config.LiveMatchesTTL),
	}
}

// ClearLive invalidates the live games cache.
func (c *ResponseCache) ClearLive() {
	c.liveMu.Lock()
	defer c.liveMu.Unlock()
	c.liveCache = nil
}

func (c *ResponseCache) evictOldestMatches() {
	now := time.Now()
	var oldestKey string
	var oldestTime time.Time
	first := true
	for key, cached := range c.matchesCache {
		if now.After(cached.expiresAt) {
			delete(c.matchesCache, key)
			continue
		}
		if first || cached.expiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = cached.expiresAt
			first = false
		}
	}
	if len(c.matchesCache) >= c.config.MaxMatchesCache && oldestKey != "" {
		delete(c.matchesCache, oldestKey)
	}
}

func (c *ResponseCache) evictOldestDetails() {
	now := time.Now()
	var oldestKey int
	var oldestTime time.Time
	first := true
	for key, cached := range c.detailsCache {
		if now.After(cached.expiresAt) {
			delete(c.detailsCache, key)
			continue
		}
		if first || cached.expiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = cached.expiresAt
			first = false
		}
	}
	if len(c.detailsCache) >= c.config.MaxDetailsCache && oldestKey != 0 {
		delete(c.detailsCache, oldestKey)
	}
}
