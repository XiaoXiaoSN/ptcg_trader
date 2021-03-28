package gormrepo_test

import (
	"context"
	"testing"

	"ptcg_trader/pkg/repository"
	"ptcg_trader/pkg/repository/gormrepo"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Test_Repository define testing case for database repository
func Test_Repository(t *testing.T) {
	suite.Run(t, new(repoTestSuite))
}

type repoTestSuite struct {
	suite.Suite

	ctx     context.Context
	Repo    repository.Repositorier
	SQLMock sqlmock.Sqlmock
}

func (s *repoTestSuite) SetupSuite() {
	// open database stub
	sqlDB, mock, err := sqlmock.New()
	s.Require().NoError(err)

	// open gorm DB
	dial := postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dial, &gorm.Config{})
	s.Require().NoError(err)

	// enable debug mode
	gormDB = gormDB.Debug()

	s.ctx = context.Background()
	s.SQLMock = mock
	s.Repo, err = gormrepo.NewRepository(gormrepo.GORMRepoParams{
		DB: gormDB,
	})
	s.Require().NoError(err)
}

func (s *repoTestSuite) SetupTest() {
}

func (s *repoTestSuite) TearDownTest() {
}

func (s *repoTestSuite) TearDownSuite() {
}
