package testchado

import (
    "bytes"
    "regexp"
    "testing"
)

func TestPostgresManager(t *testing.T) {
    if !CheckPostgresEnv() {
        t.Skip("postgres environment variable TC_DSOURCE is not set")
    }

    ds := GetDataSource()
    dbm := NewPostgresManager(ds)
    if dbm.dbsource != ds {
        t.Errorf("should have %s dbsource\n", ds)
    }
    if dbm.Driver() != "postgres" {
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
    defer dbm.DropSchema()

    sqlx := dbm.DBHandle()
    type entries struct{ Counter int }
    e := entries{}
    err := sqlx.Get(&e, "SELECT count(table_name) counter FROM information_schema.tables WHERE table_schema = ($1)", dbm.Schema)
    if err != nil {
        t.Errorf("should have executed the query %s", err)
    }
    if e.Counter != 173 {
        t.Error("should have 173 tables")
    }

    type tbls struct{ Tname string }
    tbl := tbls{}
    err = sqlx.Get(&tbl, "SELECT table_name tname FROM information_schema.tables WHERE table_schema = ($1) AND table_name = ($2)", dbm.Schema, "feature")
    if err != nil {
        t.Errorf("should have executed the query %s", err)
    }
    if tbl.Tname != "feature" {
        t.Error("should have got feature table name")
    }

    if err = dbm.DropSchema(); err != nil {
        t.Errorf("should have dropped the schema: %s", err)
    }
    err = sqlx.Get(&e, "SELECT count(table_name) counter FROM information_schema.tables WHERE table_schema = ($1)", dbm.Schema)
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

func TestPostgresLoadDefaultFixture(t *testing.T) {
    if !CheckPostgresEnv() {
        t.Skip("postgres environment variable TC_DSOURCE is not set")
    }
    ds := GetDataSource()
    dbm := NewPostgresManager(ds)
    defer dbm.DropSchema()
    _ = dbm.DeploySchema()
    if err := dbm.LoadDefaultFixture(); err != nil {
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

func TestPostgresLoadPresetFixture(t *testing.T) {
    if !CheckPostgresEnv() {
        t.Skip("postgres environment variable TC_DSOURCE is not set")
    }
    ds := GetDataSource()
    dbm := NewPostgresManager(ds)
    defer dbm.DropSchema()
    _ = dbm.DeploySchema()
    if err := dbm.LoadPresetFixture("cvprop"); err != nil {
        t.Errorf("should have loaded fixture: %s", err)
    }

    type entries struct{ Counter int }
    e := entries{}
    sqlx := dbm.DBHandle()
    err := sqlx.Get(&e, "SELECT count(*) counter FROM cvterm")
    if err != nil {
        t.Errorf("should have executed the query %s", err)
    }
    if e.Counter != 13 {
        t.Error("should have 13 cvterms")
    }
}
