package mocks

import (
	"ptcg_trader/internal/redis"
	repository "ptcg_trader/pkg/repository"
	"ptcg_trader/pkg/service"
)

// NewMockRepository new MockRepository instance
func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

// ImplMockRepository return the targe mock interface
func ImplMockRepository() (*MockRepository, repository.Repositorier) {
	m := &MockRepository{}
	return m, m
}

// NewMockRedis new MockRedis instance
func NewMockRedis() *MockRedis {
	return &MockRedis{}
}

// ImplMockRedis return the targe mock interface
func ImplMockRedis() (*MockRedis, redis.Redis) {
	m := &MockRedis{}
	return m, m
}

// NewMockMatcher new MockMatcher instance
func NewMockMatcher() *MockMatcher {
	return &MockMatcher{}
}

// ImplMockMatcher return the targe mock interface
func ImplMockMatcher() (*MockMatcher, service.Matcher) {
	m := &MockMatcher{}
	return m, m
}
