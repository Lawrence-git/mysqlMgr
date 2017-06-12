package mysqlMgr

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultDriver  = "mysql"
	defaultPort    = "3306"
	defaultMaxLife = 14400
)

//MysqlConn is the public structure that contains a channel
//controlling access to the underlying store
type MysqlConn struct {
	op chan (func(*store))
}

//Close closes the underlying channel and terminates the go-routine.
//This should be run as cleanup for each MysqlConn
func (mc MysqlConn) Close() {
	close(mc.op)
	return
}

//AddStmt adds statement q to the store with a key of k.
//Any error is returned as e
//To retrieve these, use GetStmt
func (mc MysqlConn) AddStmt(k string, q string) (e error) {
	if k == "" {
		e = fmt.Errorf("MysqlConn A key is required, key may not be empty")
		return
	}
	if q == "" {
		e = fmt.Errorf("MysqlConn A statement (query) is required and may not be empty")
		return
	}
	ech := make(chan error)
	mc.op <- func(c *store) {
		if c.stmtExists(k) {
			ech <- fmt.Errorf("A Statement with the name %s already exists in the registry.", k)
			return
		}
		var err error
		c.stmtMap[k], err = c.db.Prepare(q)
		if err != nil {
			ech <- fmt.Errorf("Error preparing statement %s : %s", q, err)
			return
		}
		ech <- nil
		return
	}
	return <-ech
}

//GetStmt returns the statement stores as key n
//Error is returned if the statement does not exist in store
func (mc MysqlConn) GetStmt(n string) (sm *sql.Stmt, e error) {
	if n == "" {
		e = fmt.Errorf("GetStmt got empty key, key is required")
		return
	}
	sch := make(chan *sql.Stmt)
	ech := make(chan error)
	mc.op <- func(c *store) {
		if !c.stmtExists(n) {
			ech <- fmt.Errorf("A Statement with the name %s does not exist in the store", n)
			sch <- nil
			return
		}
		sch <- c.stmtMap[n]
		ech <- nil
		return
	}
	return <-sch, <-ech
}

//New creates a new MysqlConn manager using the configuration passed to it
//in cfg. If there is an error creating the manager an error will be created.
//An error can be returned after the go routine is created, if this is the case,
//this method shuts down the running routing therefore, if this function returns
//any error, Close does not need to be called from the calling method/function.
func New(cfg io.Reader) (mc MysqlConn, err error) {
	cc := &config{}
	tb, e := ioutil.ReadAll(cfg)
	if e != nil {
		err = fmt.Errorf("MysqlConn Configuration Reader Error: %s", e)
		return
	}
	e = json.Unmarshal(tb, cc)
	if e != nil {
		err = fmt.Errorf("MysqlConn Configuration Error: %s", e)
		return
	}
	e = cc.genConStr()
	if e != nil {
		err = fmt.Errorf("MysqlConn Connection String Error: %s", e)
		return
	}
	db, e := sql.Open(cc.Driver, cc.dbConnStr)
	if e != nil {
		err = fmt.Errorf("MysqlConn SQL Initializaton Error: %s", e)
		return
	}
	mc.op = make(chan func(*store))
	//The listener starts here, any errors that are generated after this
	//point must close the listening channel to avoid the go-routine running
	//forever.
	go mc.loop()
	e = mc.setDb(db)
	if e != nil {
		err = fmt.Errorf("MysqlConn Database Connection Error: %s", e)
		//failing to close here would leave the go routine running.
		close(mc.op)
		return
	}
	return
}

//setDb sets and configures the sql.DB for a store
func (mc MysqlConn) setDb(d *sql.DB) (e error) {
	ler := make(chan error)
	mc.op <- func(c *store) {
		if e = d.Ping(); e != nil {
			e = fmt.Errorf("store - setDb: %s", e)
			ler <- e
			return
		}
		d.SetConnMaxLifetime(time.Second * defaultMaxLife)
		c.db = d
		ler <- nil
	}
	return <-ler
}

//loop initializes the statement map and creates the listening channel
func (mc MysqlConn) loop() {
	m := &store{}
	m.stmtMap = make(map[string]*sql.Stmt)
	for f := range mc.op {
		f(m)
	}
}

//store contains the database connector and
//map of prepared SQL statements
type store struct {
	stmtMap map[string]*sql.Stmt
	db      *sql.DB
}

//stmtExists performs checks to see if a statement with name
//sn exists in the Statement Map
func (s store) stmtExists(sn string) (b bool) {
	if _, ok := s.stmtMap[sn]; ok {
		return true
	}
	return false
}

//config contains the values used to configure/setup database connections.
type config struct {
	Database  string `json:"database"`
	Server    string `json:"server"`
	User      string `json:"user"`
	Pw        string `json:"pw"`
	Driver    string `json:"driver"`
	Port      string `json:"port"`
	dbConnStr string
}

//genConStr generates a connection string used to be passed to sql.Open
func (c *config) genConStr() (e error) {
	if c.User == "" {
		e = fmt.Errorf("Configuration missing Database User")
		return
	}
	if c.Server == "" {
		e = fmt.Errorf("Configuration missing Database Server [IP or URL]")
		return
	}
	if c.Database == "" {
		e = fmt.Errorf("Configuration missing Database Name")
		return
	}
	if c.Driver == "" {
		c.Driver = defaultDriver
	}
	if c.Port == "" {
		c.Port = defaultPort
	}
	c.dbConnStr = c.User + ":" + c.Pw + "@tcp(" + c.Server + ":" + c.Port + ")/" + c.Database
	return
}
