package mocks

import (
	"ptcg_trader/internal/redis"
	"ptcg_trader/pkg/repository"
	"ptcg_trader/pkg/service"
)

// ImplMockRepository return the targe mock interface
func ImplMockRepository() (*MockRepository, repository.Repositorier) {
	m := &MockRepository{}
	return m, m
}

// ImplMockRedis return the targe mock interface
func ImplMockRedis() (*MockRedis, redis.Redis) {
	m := &MockRedis{}
	return m, m
}

// ImplMockMatcher return the targe mock interface
func ImplMockMatcher() (*MockMatcher, service.Matcher) {
	m := &MockMatcher{}
	return m, m
}
