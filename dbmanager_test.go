package testchado

import (
    "regexp"
    "testing"
)

func TestNewDBManager(t *testing.T) {
    dbm := NewDBManager()
    if _, err := regexp.MatchString(`postgres|sqlite3`, dbm.DataSource()); err != nil {
        t.Errorf("should have matched one of the driver %s", err)
    }
}
