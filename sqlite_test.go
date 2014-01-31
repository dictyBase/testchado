package testchado

import (
    "os"
    "path/filepath"
    "testing"
)

func TestSQLiteManager(t *testing.T) {
    dbm := NewSQLiteManager()

    if dbm.dbsource != ":memory:" {
        t.Error("should have :memory: dbsource")
    }
    if dbm.Driver != "sqlite3" {
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
    fpath, err := dbm.SchemaDDL()

    if err != nil {
        t.Error("Should not throw any error")
    }
    expath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "dictybase", "testchado", "chado."+dbm.Driver)
    if fpath != expath {
        t.Error("should have returned the correct schema file path")
    }
}
