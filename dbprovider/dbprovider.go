package dbprovider

import (
	"database/sql"
	"fmt"
	"github.com/ShevchenkoVadim/helperlib/config"
	"github.com/ShevchenkoVadim/helperlib/sfotypes"
	"github.com/ShevchenkoVadim/helperlib/utils"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"net"
	"os"
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
	utils.LogWrapper("check TCP connect to DB")
	conn, err := net.DialTimeout("tcp",
		fmt.Sprint(mgr.DbConn.Server, ":", mgr.DbConn.Port), time.Second)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		utils.LogWrapper(err)
		for {
			utils.LogWrapper("check TCP connect to DB: wait to connect")
			time.Sleep(time.Second * 2)
			mgr.checkConnect()
		}
	} else {
		utils.LogWrapper("check TCP connect to DB: connect is ok")
		go mgr.writeToWaitChannel()
	}
}

func (mgr *DBManager) Connect() {
	if config.C.DBConn.Password != "" {
		utils.CreateNewCred("db_password", config.C.DBConn.Password)
	}
	dbPassword, err := utils.GetCred("db_password")
	if err != nil {
		log.Println("DB password is empty", err)
		os.Exit(1)
	}

	go mgr.checkConnect()
	<-mgr.WaitChannel
	utils.LogWrapper("DB is tried to connoent")
	//uri := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&connection+timeout=%d",
	//	mgr.DbConn.Server, mgr.DbConn.User, mgr.DbConn.Password, mgr.DbConn.Port, mgr.DbConn.Timeout, mgr.DbConn.Db)
	uri := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&connection+timeout=%d",
		config.C.DBConn.User, dbPassword, config.C.DBConn.Server,
		config.C.DBConn.Port, config.C.DBConn.Db, config.C.DBConn.Timeout)
	utils.LogWrapper(uri)

	db, err := sql.Open("sqlserver", uri)
	if err != nil {
		log.Fatal("DB connection is failed: ", err)
	}
	mgr.db = db
}

func (mgr *DBManager) SearchFileName(filePath string, userName string) (bool, error) {
	utils.LogWrapper("DB search file for user: ", userName, filePath)
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
	utils.LogWrapper("DB execute query with prepare statement: ", query)
	var stmt *sql.Stmt
	var err error

	for {
		stmt, err = mgr.db.Prepare(query)
		if err != nil {
			log.Println("DB error execute query: ", err)
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
	utils.LogWrapper("DB insert file: ", query)
	stmt, err2 := mgr.execWithPrepare(query, filePath, userName)
	return stmt, err2
}

func (mgr *DBManager) DeleteFileName(filePath string, userName string) (*sql.Stmt, error) {
	query := "DELETE FROM files WHERE filename=@p1 AND username=@p2"
	utils.LogWrapper("DB delete file: ", query)
	stmt, err2 := mgr.execWithPrepare(query, filePath, userName)
	return stmt, err2
}

func (mgr *DBManager) SelectQuery(query string) *sql.Row {
	return mgr.db.QueryRow(query)
}

func (mgr *DBManager) GetDb() *sql.DB {
	return mgr.db
}

func (mgr *DBManager) CreateTable() {
	query := "CREATE TABLE files([filename] [varchar](255) NOT NULL,[username] [varchar](255) NOT NULL,[taskstatus] [int] NULL) ON [PRIMARY]"
	utils.LogWrapper("Create table:", query)
	stmt, err := mgr.execWithPrepare(query)
	if err != nil {
		utils.LogWrapper(err)
		return
	}
	defer stmt.Close()
}

func (mgr *DBManager) DropTable() {
	query := "DROP TABLE files"
	utils.LogWrapper("Drop table:", query)
	stmt, err := mgr.execWithPrepare(query)
	if err != nil {
		utils.LogWrapper(err)
		return
	}
	defer stmt.Close()
}

func (mgr *DBManager) TuncateTable() {
	query := "TRUNCATE TABLE files"
	utils.LogWrapper("Truncate table:", query)
	stmt, err := mgr.execWithPrepare(query)
	if err != nil {
		utils.LogWrapper(err)
		return
	}
	defer stmt.Close()
}
