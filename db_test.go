package main

// This contains mostly integration tests for db.go

import (
	"fmt"
	"testing"
	"time"

	"os/exec"

	"github.com/stretchr/testify/suite"
	null "gopkg.in/guregu/null.v3"
)

const (
	testDbType  = "postgres"
	testDSN     = "host=127.0.0.1 port=5432 user=postgres password=crimsonsux sslmode=disable"
	createDbSql = "CREATE DATABASE %s ;"
	dropDbSql   = "DROP DATABASE %s ;"
)

type dbTestSuite struct {
	suite.Suite
	p      *postgresDb
	dbName string
}

func (s *dbTestSuite) Test_CreateRealm() {
	realmName := null.StringFrom("crimson")
	realmTitle := null.StringFrom("sux")
	r, err := s.p.CreateRealm(realmName, realmTitle)
	s.Require().NotNil(r)
	s.Require().NoError(err)
	s.Require().Equal(realmName.String, r.Name)
	s.Require().Equal(realmTitle, r.Title)
	s.Require().NotZero(r.ID)

	var testRealm Realm
	s.Require().NoError(s.p.db.Get(&testRealm, fmt.Sprintf("SELECT * from realm where id = %d", r.ID)))
	s.Require().Equal(realmName.String, testRealm.Name)
	s.Require().Equal(realmTitle, testRealm.Title)
	s.Require().Equal(r.ID, testRealm.ID)
}

func (s *dbTestSuite) Test_CreatePlayer() {
	r, _ := s.p.CreateRealm(null.StringFrom("testName"), null.StringFrom("testTitle"))
	s.Require().NotZero(r.ID)
	playerName := null.StringFrom("crimson")
	realmID := null.IntFrom(int64(r.ID))
	p, err := s.p.CreatePlayer(playerName, realmID)
	s.Require().NotNil(p)
	s.Require().NoError(err)
	s.Require().Equal(playerName.String, p.Name)
	s.Require().Equal(realmID.Int64, int64(p.RealmID))
	s.Require().NotZero(p.ID)

	var testPlayer Player
	s.Require().NoError(s.p.db.Get(&testPlayer, fmt.Sprintf("SELECT * from player where id = %d", p.ID)))
	s.Require().Equal(playerName.String, testPlayer.Name)
	s.Require().Equal(realmID.Int64, int64(p.RealmID))
	s.Require().Equal(p.ID, testPlayer.ID)
}

func (s *dbTestSuite) Test_CreateSession() {
	r, _ := s.p.CreateRealm(null.StringFrom("testName"), null.StringFrom("testTitle"))
	s.Require().NotZero(r.ID)
	realmID := null.IntFrom(int64(r.ID))
	sessionName := null.StringFrom("christmas poker night 2017")
	p1, _ := s.p.CreatePlayer(null.StringFrom("p1"), realmID)
	p2, _ := s.p.CreatePlayer(null.StringFrom("p2"), realmID)
	p3, _ := s.p.CreatePlayer(null.StringFrom("p3"), realmID)
	playerSessions := []PlayerSession{
		PlayerSession{
			PlayerID: p1.ID,
			Buyin:    null.IntFrom(500),
			Walkout:  null.IntFrom(1250),
		},
		PlayerSession{
			PlayerID: p2.ID,
			Buyin:    null.IntFrom(1000),
			Walkout:  null.IntFrom(0),
		},
		PlayerSession{
			PlayerID: p3.ID,
			Buyin:    null.IntFrom(1500),
			Walkout:  null.IntFrom(1),
		},
	}
	now := null.TimeFrom(time.Now())
	ps, err := s.p.CreateSession(realmID, sessionName, now, playerSessions)
	s.Require().NotNil(ps)
	s.Require().NoError(err)
	s.Require().NotZero(ps.ID)
	s.Require().Equal(sessionName, ps.Name)
	s.Require().Equal(now.Time.Unix(), ps.Time.Unix())
	s.Require().Equal(realmID.Int64, int64(ps.RealmID))

	dbPS := PlayerSession{}
	rows, err := s.p.db.Queryx("SELECT * FROM player_session where session_id = $1 order by player_id", ps.ID)
	s.Require().NoError(err)
	for i := 0; rows.Next(); i++ {
		err := rows.StructScan(&dbPS)
		s.Require().NoError(err)
		playerSessions[i].CreatedAt = dbPS.CreatedAt
		playerSessions[i].SessionID = ps.ID
		s.Require().Equal(playerSessions[i], dbPS)
	}

	t := Transfer{}
	rows, err = s.p.db.Queryx("SELECT * FROM transfer where session_id = $1 order by player_id", ps.ID)
	s.Require().NoError(err)
	for i := 0; rows.Next(); i++ {
		err := rows.StructScan(&t)
		s.Require().NoError(err)
		s.Require().EqualValues(playerSessions[i].Walkout.Int64-playerSessions[i].Buyin.Int64, t.Amount)
	}
}

func (s *dbTestSuite) SetupTest() {
	s.p.db.MustExec("DELETE FROM player_session")
	s.p.db.MustExec("DELETE FROM transfer")
	s.p.db.MustExec("DELETE FROM session")
	s.p.db.MustExec("DELETE FROM player")
	s.p.db.MustExec("DELETE FROM realm")
}

// SetupSuite provisions a fresh test db and runs migrations on it. Also inits the dbTestSuite struct
func (s *dbTestSuite) SetupSuite() {
	setupDb := mustInitDB(testDbType, testDSN)
	s.dbName = fmt.Sprintf("dbtest_%d", time.Now().Unix())
	setupDb.db.MustExec(fmt.Sprintf(createDbSql, s.dbName))
	testDSNWithDb := fmt.Sprintf("%s database=%s", testDSN, s.dbName)

	// wish goose had cleaner way to programmatically migrate with .sql files
	migration := exec.Command("goose", "-dir", "migrations", testDbType, testDSNWithDb, "up")
	s.Require().NoError(migration.Run(), "migration failed")
	setupDb.Close()
	// setup test db connection to be used for rest of tests
	s.p = mustInitDB(testDbType, testDSNWithDb)
}

// Drop the test db
func (s *dbTestSuite) TearDownSuite() {
	s.p.Close()
	setupDb := mustInitDB(testDbType, testDSN)
	setupDb.db.MustExec(fmt.Sprintf(dropDbSql, s.dbName))
	setupDb.Close()
}

func Test_DbTestSuite(t *testing.T) {
	// since this is mostly integration tests, skip if short mode
	if testing.Short() {
		t.SkipNow()
	}
	ts := new(dbTestSuite)
	//defer ts.TearDownSuite()

	suite.Run(t, ts)
}