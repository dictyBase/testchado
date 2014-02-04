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
    // Name of the database, might vary by implementation
    Database() string
    // Removes the active chado schema from the database
    DropSchema() error
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

// Loads the default fixture in the chado schema. The default fixture include.
//  1.List of default organisms.
//  2.Sequnence ontology(SO)
//  3.Relation ontology(RO)
func (dbh *DBHelper) LoadFixture() error {
    var c bytes.Buffer
    sqlx := dbh.dbhandler
    zpath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "dictybase", "testchado")
    zfile := filepath.Join(zpath, "preset.zip")
    zr, err := zip.OpenReader(zfile)
    if err != nil {
        return err
    }
    defer zr.Close()
    for _, f := range zr.File {
        if f.Name == "default.sql" {
            zc, err := f.Open()
            if err != nil {
                return err
            }
            _, err = io.Copy(&c, zc)
            if err != nil {
                return err
            }
            break
        }
    }
    _ = sqlx.Execf(c.String())
    return nil
}

func (dbh *DBHelper) LoadCustomFixture(fixture string) error {
    sqlx := dbh.dbhandler
    _, err := sqlx.LoadFile(fixture)
    return err
}

// The active database connection
func (dbh *DBHelper) DBHandle() *sqlx.DB {
    return dbh.dbhandler
}

// Name of datasource in a format understandable by database/sql package
func (dbh *DBHelper) DataSource() string {
    return dbh.dbsource
}
