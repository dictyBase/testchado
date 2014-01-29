package testchado

import (
    _ "github.com/cybersiddhu/go-sqlite3"
    "github.com/jmoiron/sqlx"
    "log"
)

type Sqlite struct {
    *DBHelper
}

func NewDBManager() *Sqlite {
    dbh, err := sqlx.Open("sqlite3", ":memory:")
    if err != nil {
        log.Fatal(err)
    }
    return &Sqlite{&DBHelper{dbsource: ":memory:", Driver: "sqlite3", dbhandler: dbh}}
}

func (sqlite *Sqlite) DBHandle() *sqlx.DB {
    return sqlite.DBHelper.dbhandler
}

func (sqlite *Sqlite) Database() string {
    return ""
}

func (sqlite *Sqlite) DataSource() string {
    return sqlite.DBHelper.dbsource
}

func (sqlite *Sqlite) DropSchema() error {
    dbh := sqlite.DBHelper.dbhandler
    type table struct{ Name string }
    tbls := []table{}
    err := dbh.Select(&tbls, "SELECT name FROM sqlite_master where type = ?", "table")
    if err != nil {
        return err
    }
    tx := dbh.MustBegin()
    for _, tbl := range tbls {
        _ = tx.Execf("DROP table ?", tbl)
    }
    err = tx.Commit()
    if err != nil {
        return err
    }
    return nil
}

func (sqlite *Sqlite) DeploySchema() error {
    dbh := sqlite.DBHelper.dbhandler
    schema, err := sqlite.SchemaDDL()
    if err != nil {
        return err
    }
    _, err = dbh.LoadFile(schema)
    if err != nil {
        return err
    }
    return nil
}

func (sqlite *Sqlite) DropDatabase() error {
    return nil
}

func (sqlite *Sqlite) CreateDatabase() error {
    return nil
}

func (sqlite *Sqlite) ResetSchema() error {
    err := sqlite.DropSchema()
    if err != nil {
        return err
    }
    err = sqlite.DeploySchema()
    if err != nil {
        return err
    }
    return nil
}
