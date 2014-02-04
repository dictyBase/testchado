package testchado

import (
    "regexp"
    "testing"
)

func TestNewChadoSchema(t *testing.T) {
    dbm := NewChadoSchema()
    if _, err := regexp.MatchString(`postgres|sqlite3`, dbm.DataSource()); err != nil {
        t.Errorf("should have matched one of the driver %s", err)
    }
}
