package postgres_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-courier/sqlx/v2"
	. "github.com/onsi/gomega"
	"github.com/saitofun/qkit/conf/postgres"
)

var (
	dbName   string
	dbUser   string
	dbPasswd string
)

func init() {
	dbName = os.Getenv("PG_TEST_DB_NAME")
	dbUser = os.Getenv("PG_TEST_DB_USER")
	dbPasswd = os.Getenv("PG_TEST_DB_PASSWD")
}

func TestEndpoint(t *testing.T) {
	var (
		pg        = &postgres.Endpoint{Database: &sqlx.Database{Name: dbName}}
		masterURL = []byte(fmt.Sprintf("postgresql://%s:%s@127.0.0.1/%s?sslmode=disable", dbUser, dbPasswd, dbName))
		slaveURL  = []byte(fmt.Sprintf("postgresql://%s:%s@localhost/%s?sslmode=disable", dbUser, dbPasswd, dbName))
	)

	NewWithT(t).Expect(pg.Master.UnmarshalText(masterURL)).To(BeNil())
	NewWithT(t).Expect(pg.Slave.UnmarshalText(slaveURL)).To(BeNil())

	pg.SetDefault()
	pg.Init()

	{
		row, err := pg.QueryContext(context.Background(), "SELECT 1")
		NewWithT(t).Expect(err).To(BeNil())
		_ = row.Close()
	}

	{
		row, err := postgres.SwitchSlave(pg).QueryContext(context.Background(), "SELECT 1")
		NewWithT(t).Expect(err).To(BeNil())
		_ = row.Close()
	}

	NewWithT(t).Expect(pg.UseSlave()).NotTo(Equal(pg.DB))
	NewWithT(t).Expect(pg.LivenessCheck()).To(
		Equal(map[string]string{
			pg.Master.Host(): "ok",
			pg.Slave.Host():  "ok",
		}),
	)
}