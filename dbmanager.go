package testchado

import (
    "archive/zip"
    "github.com/jmoiron/sqlx"
    "os"
    "path/filepath"
)

// Interface for managing the lifecycle of a chado database. Any backend should implement
// this interface
type DBManager interface {
    // The active database connection
    DBHandle() *sqlx.DB
    // Name of the database, might vary by implementation
    Database() string
    // Removes the active chado schema from the database
    DropSchema() error
    // Name of datasource in a format understandable by database/sql package
    DataSource() string
    // Loads chado schema in the database
    DeploySchema() error
    // Reloads chado schema in the database
    ResetSchema() error
}

// A type that provides few helper attributes for implementing DBManager interface
// All backends are encouraged to embed this type in their implementation.
type DBHelper struct {
    // Database user
    User string
    // Database password
    Password string
    // Database driver
    Driver    string
    dbsource  string
    dbhandler *sqlx.DB
}

// Give the full path to the chado schema for a particular backend
func (dbh *DBHelper) SchemaDDL() (fpath string, err error) {
    zpath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "dictybase", "testchado")
    zr, err := zip.OpenReader(zpath)
    if err != nil {
        return
    }
    defer zr.Close()
    name := dbh.Driver + ".chado"
    for _, f := range zr.File {
        if f.Name == name {
            fpath = filepath.Join(zpath, f.Name)
            break
        }
    }
    return
}
