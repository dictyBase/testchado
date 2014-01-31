package testchado

import (
    "archive/zip"
    "bytes"
    "github.com/jmoiron/sqlx"
    "io"
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

// Return the content of chado schema for a particular backend
func (dbh *DBHelper) SchemaDDL() (*bytes.Buffer, error) {
    var c bytes.Buffer
    zpath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "dictybase", "testchado")
    zfile := filepath.Join(zpath, "chado.zip")
    zr, err := zip.OpenReader(zfile)
    if err != nil {
        return &c, err
    }
    defer zr.Close()
    name := "chado." + dbh.Driver
    for _, f := range zr.File {
        if f.Name == name {
            zc, err := f.Open()
            if err != nil {
                return &c, err
            }
            _, err = io.Copy(&c, zc)
            if err != nil {
                return &c, err
            }
            break
        }
    }
    return &c, nil
}
