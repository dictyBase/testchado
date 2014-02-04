package testchado

import (
    "bytes"
    "os"
    "regexp"
    "testing"
)

func CheckPostgresEnv() bool {
    if len(os.Getenv("TC_DSOURCE")) > 0 {
        return true
    }
    return false
}

func GetDataSource() string {
    return os.Getenv("TC_DSOURCE")
}

func TestPostgresManager(t *testing.T) {
    if !CheckPostgresEnv() {
        t.Skip("postgres environment variable TC_DSOURCE is not set")
    }

    ds := GetDataSource()
    dbm := NewPostgresManager(ds)
    if dbm.dbsource != ds {
        t.Errorf("should have %s dbsource\n", ds)
    }
    if dbm.Driver != "postgres" {
        t.Error("should have postgres driver")
    }
    if dbm.dbhandler == nil {
        t.Error("should have dbhandler instance")
    }
    if dbm.Database() == "something" {
        t.Error("should not have any database name")
    }
    if dbm.dbsource != dbm.DataSource() {
        t.Error("should have identical datasource")
    }
    if _, err := regexp.MatchString(`^\w+$`, dbm.Schema); err != nil {
        t.Errorf("should have matched schema %s name\n", dbm.Schema)
    }
}

func TestPostgresSchemaPath(t *testing.T) {
    if !CheckPostgresEnv() {
        t.Skip("postgres environment variable TC_DSOURCE is not set")
    }
    ds := GetDataSource()
    dbm := NewPostgresManager(ds)
    content, err := dbm.SchemaDDL()
    if err != nil {
        t.Errorf("Should not throw any error: %s", err)
    }
    if !bytes.Contains(content.Bytes(), []byte("feature")) {
        t.Error("should have contain feature")
    }
}

func TestPostgresSchemaCRUD(t *testing.T) {
    if !CheckPostgresEnv() {
        t.Skip("postgres environment variable TC_DSOURCE is not set")
    }
    ds := GetDataSource()
    dbm := NewPostgresManager(ds)
    if err := dbm.DeploySchema(); err != nil {
        t.Errorf("error %s: should have deployed the chado schema", err)
    }

    sqlx := dbm.DBHandle()
    type tbls struct{ Name string }
    tbl := tbls{}
    err := sqlx.Get(&tbl, "SELECT table_name name FROM information_schema.tables where table_schema = ? and table_name = ?", dbm.Schema, "feature")
    if err != nil {
        t.Errorf("should have executed the query %s", err)
    }
    if tbl.Name != "feature" {
        t.Error("should have got feature table name")
    }

    type entries struct{ Counter int }
    e := entries{}
    err = sqlx.Get(&e, "SELECT count(table_name) counter FROM information_schema.tables where table_schema = ?", dbm.Schema)
    if err != nil {
        t.Errorf("should have executed the query %s", err)
    }
    if e.Counter != 174 {
        t.Error("should have 174 tables")
    }

    if err = dbm.DropSchema(); err != nil {
        t.Errorf("should have dropped the schema: %s", err)
    }
    err = sqlx.Get(&e, "SELECT count(table_name) counter FROM information_schema.tables where table_schema = ?", dbm.Schema)
    if err != nil {
        t.Errorf("should have executed the query %s", err)
    }
    if e.Counter != 0 {
        t.Error("should have 0 tables")
    }

    _ = dbm.DeploySchema()
    if err = dbm.ResetSchema(); err != nil {
        t.Errorf("should have reset the schema: %s", err)
    }
}

func TestPostgresLoadFixture(t *testing.T) {
    if !CheckPostgresEnv() {
        t.Skip("postgres environment variable TC_DSOURCE is not set")
    }
    ds := GetDataSource()
    dbm := NewPostgresManager(ds)
    _ = dbm.DeploySchema()
    if err := dbm.LoadFixture(); err != nil {
        t.Errorf("should have loaded fixture: %s", err)
    }

    type entries struct{ Counter int }
    e := entries{}
    sqlx := dbm.DBHandle()
    err := sqlx.Get(&e, "SELECT count(*) counter FROM organism")
    if err != nil {
        t.Errorf("should have executed the query %s", err)
    }
    if e.Counter != 12 {
        t.Error("should have 12 organisms")
    }

    query := `
     SELECT count(cvterm.cvterm_id) counter from CVTERM join CV on CV.CV_ID=CVTERM.CV_ID
     WHERE CV.NAME = 'sequence'
    `
    err = sqlx.Get(&e, query)
    if err != nil {
        t.Errorf("should have executed the query %s", err)
    }
    if e.Counter != 286 {
        t.Error("should have 286 sequence ontology term")
    }
}
