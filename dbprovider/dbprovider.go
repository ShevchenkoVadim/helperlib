package dbprovider

import (
	"database/sql"
	"fmt"
	"github.com/ShevchenkoVadim/helperlib/sfotypes"
	"github.com/ShevchenkoVadim/helperlib/utils"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"net"
	"time"
)

type Manager interface {
	InsertFileName(filePath string, userName string) error
	SearchFileName(filePath string, userName string) (bool, error)
	DeleteFileName(filePath string, userName string) error
}

type DBManager struct {
	DbConn      sfotypes.DBConnection
	db          *sql.DB
	WaitChannel chan bool
}

func (mgr *DBManager) writeToWaitChannel() {
	func() {
		mgr.WaitChannel <- true
	}()
}

func (mgr *DBManager) checkConnect() {
	utils.LogWrapper("checkConnect")
	conn, err := net.DialTimeout("tcp",
		fmt.Sprint(mgr.DbConn.Server, ":", mgr.DbConn.Port), time.Second)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		utils.LogWrapper(err)
		for {
			utils.LogWrapper("checkConnect: Wait to connect")
			time.Sleep(time.Second * 2)
			mgr.checkConnect()
		}
	} else {
		utils.LogWrapper("checkConnect Ok")
		go mgr.writeToWaitChannel()
	}
}

func (mgr *DBManager) Connect() {
	utils.LogWrapper("DB Connect")
	go mgr.checkConnect()
	<-mgr.WaitChannel
	uri := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;connection timeout=%d",
		mgr.DbConn.Server, mgr.DbConn.User, mgr.DbConn.Password, mgr.DbConn.Port, mgr.DbConn.Timeout)
	fmt.Println(uri)
	db, err := sql.Open("sqlserver", uri)
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}
	mgr.db = db
}

func (mgr *DBManager) SearchFileName(filePath string, userName string) (bool, error) {
	utils.LogWrapper("DB SearchFileName")
	var name string
	var err error
	for {
		err = mgr.db.QueryRow("SELECT filename FROM files WHERE filename=@p1 AND username=@p2", filePath, userName).Scan(&name)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return false, nil
			}
			go mgr.checkConnect()
			<-mgr.WaitChannel
			continue
		}
		break
	}
	return true, nil
}

func (mgr *DBManager) execWithPrepare(query string, args ...any) (*sql.Stmt, error) {
	utils.LogWrapper("execWithPrepare")
	var stmt *sql.Stmt
	var err error

	for {
		stmt, err = mgr.db.Prepare(query)
		if err != nil {
			log.Println("ERROR: ", err)
			go mgr.checkConnect()
			<-mgr.WaitChannel
			continue
		}
		break
	}
	_, err2 := stmt.Exec(args...)
	return stmt, err2
}

func (mgr *DBManager) InsertFileName(filePath string, userName string) (*sql.Stmt, error) {
	query := "INSERT INTO files (filename, username) VALUES (@p1, @p2)"
	utils.LogWrapper(query)
	stmt, err2 := mgr.execWithPrepare(query, filePath, userName)
	return stmt, err2
}

func (mgr *DBManager) DeleteFileName(filePath string, userName string) (*sql.Stmt, error) {
	query := "DELETE FROM files WHERE filename=@p1 AND username=@p2"
	stmt, err2 := mgr.execWithPrepare(query, filePath, userName)
	return stmt, err2
}
