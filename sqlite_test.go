package testchado

import (
    "bytes"
    "testing"
)

func TestSQLiteManager(t *testing.T) {
    dbm := NewSQLiteManager()

    if dbm.dbsource != ":memory:" {
        t.Error("should have :memory: dbsource")
    }
    if dbm.Driver() != "sqlite3" {
        t.Error("should have sqlite3 driver")
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
}

func TestSQLiteSchemaPath(t *testing.T) {
    dbm := NewSQLiteManager()
    content, err := dbm.SchemaDDL()

    if err != nil {
        t.Errorf("Should not throw any error: %s", err)
    }
    if !bytes.Contains(content.Bytes(), []byte("feature")) {
        t.Error("should have contain feature")
    }
}

func TestSQLiteSchemaCRUD(t *testing.T) {
    dbm := NewSQLiteManager()
    if err := dbm.DeploySchema(); err != nil {
        t.Errorf("error %s: should have deployed the chado schema", err)
    }

    if !dbm.hasLoadedSchema {
        t.Error("should have been set after schema deployment")
    }

    sqlx := dbm.DBHandle()
    type tbls struct{ Name string }
    tbl := tbls{}
    err := sqlx.Get(&tbl, "SELECT name FROM sqlite_master where type = ? and tbl_name = ?", "table", "feature")
    if err != nil {
        t.Errorf("should have executed the query %s", err)
    }
    if tbl.Name != "feature" {
        t.Error("should have got feature table name")
    }

    type entries struct{ Counter int }
    e := entries{}
    err = sqlx.Get(&e, "SELECT count(name) counter FROM sqlite_master where type = ?", "table")
    if err != nil {
        t.Errorf("should have executed the query %s", err)
    }
    if e.Counter != 174 {
        t.Error("should have 172 tables")
    }

    if err = dbm.DropSchema(); err != nil {
        t.Errorf("should have dropped the schema: %s", err)
    }
    if dbm.hasLoadedSchema {
        t.Error("should not have been set after schema deployment")
    }
    err = sqlx.Get(&e, "SELECT count(name) counter FROM sqlite_master where type = ?", "table")
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
    if !dbm.hasLoadedSchema {
        t.Error("should have been set after schema re-deployment")
    }
}

func TestSQLiteLoadDefaultFixture(t *testing.T) {
    dbm := NewSQLiteManager()
    if err := dbm.LoadDefaultFixture(); err == nil {
        t.Error("should have not loaded default fixture")
    }
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

func TestSQLiteLoadPresetFixture(t *testing.T) {
    dbm := NewSQLiteManager()
    if err := dbm.LoadPresetFixture("cvprop"); err == nil {
        t.Error("should have not loaded preset fixture")
    }
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
