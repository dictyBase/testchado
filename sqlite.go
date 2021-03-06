package testchado

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// A type specific for sqlite backend
type Sqlite struct {
	*DBHelper
}

// Get a in memory instance of sqlite DBManager
func NewSQLiteManager() *Sqlite {
	gm, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	gm.SingularTable(true)
	sqlx := sqlx.NewDb(gm.DB(), "sqlite3")
	return &Sqlite{&DBHelper{dbsource: ":memory:", driver: "sqlite3", dbhandler: sqlx, gormHandler: &gm}}
}

func (sqlite *Sqlite) Database() string {
	return ""
}

func (sqlite *Sqlite) DropSchema() error {
	dbh := sqlite.DBHandle()
	type table struct{ Name string }
	tbls := []table{}
	err := dbh.Select(&tbls, "SELECT name FROM sqlite_master where type = ?", "table")
	if err != nil {
		return err
	}
	tx := dbh.MustBegin()
	for _, tbl := range tbls {
		stmt := "DROP TABLE " + tbl.Name
		_ = tx.MustExec(stmt)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	sqlite.DBHelper.hasLoadedSchema = false
	return nil
}

func (sqlite *Sqlite) DeploySchema() error {
	dbh := sqlite.DBHandle()
	content, err := sqlite.SchemaDDL()
	if err != nil {
		return err
	}
	_ = dbh.MustExec(content.String())
	sqlite.DBHelper.hasLoadedSchema = true
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
	sqlite.DBHelper.hasLoadedSchema = true
	return nil
}
