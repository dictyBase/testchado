package testchado

import (
	"bytes"
	"log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Generates a random string between a range(min and max) of length
func RandomString(min, max int) string {
	alphanum := []byte("abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UTC().UnixNano())
	size := min + rand.Intn(max-min)
	b := make([]byte, size)
	alen := len(alphanum)
	for i := 0; i < size; i++ {
		pos := rand.Intn(alen)
		b[i] = alphanum[pos]
	}
	return string(b)
}

// A type specific for postgresql backend
type Postgres struct {
	*DBHelper
	Schema string
}

// Get an instance of postgres DBManager.
// For details about datasource look here http://godoc.org/github.com/lib/pq
func NewPostgresManager(datasource string) *Postgres {
	gm, err := gorm.Open("postgres", datasource)
	if err != nil {
		log.Fatal(err)
	}
	gm.SingularTable(true)
	sqlx := sqlx.NewDb(gm.DB(), "postgres")
	schema := RandomString(9, 10)
	return &Postgres{&DBHelper{dbsource: datasource, driver: "postgres", dbhandler: sqlx, gormHandler: &gm}, schema}
}

func (postgres *Postgres) Database() string {
	return ""
}

func (postgres *Postgres) DeploySchema() error {
	schema := postgres.Schema
	// Setup the schema
	buff := bytes.NewBufferString("DROP SCHEMA IF EXISTS " + schema + " CASCADE;\n")
	_, err := buff.WriteString("CREATE SCHEMA " + schema + ";\n")
	if err != nil {
		return err
	}
	_, err = buff.WriteString("SET search_path TO " + schema + ";\n")
	if err != nil {
		return err
	}
	//Now get schema definition
	content, err := postgres.SchemaDDL()
	if err != nil {
		return err
	}
	buff.Write(content.Bytes())

	// Do everything in transaction
	tx := postgres.DBHandle().MustBegin()
	// Load schema in postgresql
	_ = tx.MustExec(buff.String())
	err = tx.Commit()
	if err != nil {
		return err
	}
	postgres.DBHelper.hasLoadedSchema = true
	return nil
}

func (postgres *Postgres) DropSchema() error {
	tx := postgres.DBHandle().MustBegin()
	_ = tx.MustExec("DROP SCHEMA IF EXISTS " + postgres.Schema + " CASCADE")
	err := tx.Commit()
	if err != nil {
		return err
	}
	postgres.Schema = RandomString(9, 10)
	postgres.DBHelper.hasLoadedSchema = false
	return nil
}

func (postgres *Postgres) ResetSchema() error {
	err := postgres.DropSchema()
	if err != nil {
		return err
	}
	err = postgres.DeploySchema()
	if err != nil {
		return err
	}
	postgres.DBHelper.hasLoadedSchema = true
	return nil
}
