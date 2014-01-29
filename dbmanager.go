package testchado

import (
    "archive/zip"
    "database/sql"
    "os"
    "path/filepath"
)

type DBManager interface {
    DBHandle() *sql.DB
    Database() string
    DropSchema() error
    CreateDatabase() error
    DropDatabase() error
    GetClient() (string, error)
    DeployByClient(string) error
    DeployByDB() error
    DataSource() string
    DeploySchema() error
    ResetSchema() error
}

type DBHelper struct {
    User      string
    Password  string
    Driver    string
    dbsource  string
    dbhandler *sql.DB
}

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
            fpath := filepath.Join(zpath, f.Name)
            break
        }
    }
    return
}
