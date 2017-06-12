package mysqlMgr

import (
	"database/sql"
	"testing"
)

//TestStmtExists test the helper function for existing statements
func TestStmtExists(t *testing.T) {
	testStore := store{}
	testStore.stmtMap = make(map[string]*sql.Stmt)
	testStore.stmtMap["test_query"] = nil
	if !testStore.stmtExists("test_query") {
		t.Fatal("stmtExists returned incorrect value")
	}
}

//TestConString checks if the connection string is being properly generated
func TestConString(t *testing.T) {
	sample := config{Database: "test", Server: "test", User: "test", Pw: "test"}
	e := sample.genConStr()
	if e != nil {
		t.Fatal("Test Con String was invalid:")
	}
}

//TestDefaults checks if the defaults are correctly applied to the configuration
func TestDefaults(t *testing.T) {
	sample := config{Database: "test", Server: "test", User: "test", Pw: "test"}
	e := sample.genConStr()
	if e != nil {
		t.Fatalf("Default port conn str returned: %s", e)
	}
	if sample.Port != "3306" {
		t.Fatal("Default port not set with empty port in config")
	}
	if sample.Driver != "mysql" {
		t.Fatal("Default driver not set with empty driver in config")
	}
}

//TestMissingConfigValues ensures that the inspection of the configuration
//returns the proper errors for missing values
func TestMissingConfigValues(t *testing.T) {
	sample := config{}
	e := sample.genConStr()
	if e == nil {
		t.Fatal("Missing User Value should have returned error, returned nil")
	}
	sample.User = "test"
	e = sample.genConStr()
	if e == nil {
		t.Fatal("Missing Server Value should have returned error, returned nil")
	}
	sample.Server = "test"
	e = sample.genConStr()
	if e == nil {
		t.Fatal("Missing Database Value should have returned error, returned nil")
	}
	sample.Database = "test"
	e = sample.genConStr()
	if e != nil {
		t.Fatal("Minimum configuration should have returned nil, returned error")
	}
}
