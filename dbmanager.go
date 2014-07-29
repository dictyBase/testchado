package testchado

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

// Interface for managing the lifecycle of a chado database. Any backend should implement
// this interface
type DBManager interface {
	// Name of the database, might vary by implementation
	Database() string
	// The active database connection
	DBHandle() *sqlx.DB
	// An instance of gorm(https://github.com/jinzhu/gorm) ORM
	GormHandle() *gorm.DB
	// Name of datasource in a format understandable by database/sql package
	DataSource() string
	// Removes the active chado schema from the database
	DropSchema() error
	// Loads chado schema in the database
	DeploySchema() error
	// Reloads chado schema in the database
	ResetSchema() error
	// Return the content of chado schema for a particular backend
	SchemaDDL() (*bytes.Buffer, error)
	// Loads the default fixture in the chado schema. The default fixture include.
	//  1.List of default organisms.
	//  2.Sequnence ontology(SO)
	//  3.Relation ontology(RO)
	LoadDefaultFixture() error
	// Loads one of the preset fixture that comes bundled with testchado. Currently it could be one of
	// cvprop or eco.
	LoadPresetFixture(string) error
	// Loads a custom fixture in the test database. It accepts file containing sql statements to insert fixture.
	// The sql statements are generally series of INSERT statements one in a single line, however any other
	// accpetable forms are allowed as long as they are compatible with the backend.
	LoadCustomFixture(string) error
}

// A type that provides few helper attributes for implementing DBManager interface
// All backends are encouraged to embed this type in their implementation.
type DBHelper struct {
	driver          string
	dbsource        string
	dbhandler       *sqlx.DB
	hasLoadedSchema bool
	gormHandler     *gorm.DB
}

func currSrcDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("unable to retreive current src file path")
	}
	return filepath.Dir(filename)
}

// Return the content of chado schema for a particular backend
func (dbh *DBHelper) SchemaDDL() (*bytes.Buffer, error) {
	var c bytes.Buffer
	zfile := filepath.Join(currSrcDir(), "chado.zip")
	zr, err := zip.OpenReader(zfile)
	if err != nil {
		return &c, err
	}
	defer zr.Close()
	name := "chado." + dbh.Driver()
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
func (dbh *DBHelper) LoadDefaultFixture() error {
	if !dbh.hasLoadedSchema {
		return fmt.Errorf("chado schema is not loaded")
	}
	var c bytes.Buffer
	sqlx := dbh.dbhandler
	zfile := filepath.Join(currSrcDir(), "preset.zip")
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
	_ = sqlx.MustExec(c.String())
	return nil
}

func (dbh *DBHelper) LoadPresetFixture(name string) error {
	if !dbh.hasLoadedSchema {
		return fmt.Errorf("chado schema is not loaded")
	}
	var c bytes.Buffer
	sqlx := dbh.dbhandler
	zfile := filepath.Join(currSrcDir(), "preset.zip")
	zr, err := zip.OpenReader(zfile)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, f := range zr.File {
		if f.Name == name+".sql" {
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
	_ = sqlx.MustExec(c.String())
	return nil
}

func (dbh *DBHelper) LoadCustomFixture(file string) error {
	if !dbh.hasLoadedSchema {
		return fmt.Errorf("chado schema is not loaded")
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return err
	}
	//sqlx := dbh.dbhandler
	_, err := sqlx.LoadFile(dbh.dbhandler, file)
	return err
}

// The active database connection
func (dbh *DBHelper) DBHandle() *sqlx.DB {
	return dbh.dbhandler
}

// An instance of gorm(https://github.com/jinzhu/gorm) ORM
func (db *DBHelper) GormHandle() *gorm.DB {
	return db.gormHandler
}

func (dbh *DBHelper) Driver() string {
	return dbh.driver
}

// Name of datasource in a format understandable by database/sql package
func (dbh *DBHelper) DataSource() string {
	return dbh.dbsource
}

// Returns a new instance of DBManager.
// By default, it gives an instance of sqlite backend.
// If TC_DSOURCE env variable is set, returns a postgres backend.
func NewDBManager() DBManager {
	if CheckPostgresEnv() {
		return NewPostgresManager(GetDataSource())
	}
	return NewSQLiteManager()
}

func CheckPostgresEnv() bool {
	if len(os.Getenv("TC_DSOURCE")) > 0 {
		return true
	}
	return false
}

func GetDataSource() string {
	return os.Getenv("TC_DSOURCE")
}
