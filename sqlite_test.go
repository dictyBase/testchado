package testchado

import (
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

}
