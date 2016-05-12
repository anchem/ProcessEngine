// DBConnector
package db

import (
	c "ProcessEngine/src/constant"
	logger "ProcessEngine/src/logger"
	m "ProcessEngine/src/model"
	"ProcessEngine/src/util"
	"bytes"
	"database/sql"
	"encoding/gob"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	LOG_DB_ERROR  = "<DBConnector> [ERROR] : "
	LOG_DB_WARING = "<DBConnector> [WARNING] : "
	LOG_DB_INFO   = "<DBConnector> [INFO] : "
)

var _DBCnt *DBConnector

type DBConnector struct {
	Sqldb  *sql.DB
	DBType string
	DBUrl  string
}

func GetDBConnector() *DBConnector {
	if _DBCnt == nil {
		initDBConnector()
	}
	return _DBCnt
}
func initDBConnector() {
	_DBCnt = new(DBConnector)
	_DBCnt.DBType = m.GetContext().Config.DBType
	_DBCnt.DBUrl = m.GetContext().Config.Username + ":" + m.GetContext().Config.Password + "@tcp(" + m.GetContext().Config.DBAddress + ")/" + m.GetContext().Config.DBName + "?charset=utf8"
	db1, err := sql.Open(_DBCnt.DBType, _DBCnt.DBUrl)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"init error, cannot open database", err)
		panic(err)
	}
	_DBCnt.Sqldb = db1
}

func (dbCnt *DBConnector) SaveProcDef(file string, process *m.Process, mode byte) byte {
	var dbQuery string
	switch mode {
	case c.PROCENG_DEPLOY_MDOE_DEFAULT:
		// check whether processDef exists
		rows, err := dbCnt.Sqldb.Query("SELECT * FROM process_definition WHERE proc_id=?", process.Id)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"save process def query error", err)
			return c.DB_EXEC_FAILED
		}
		if rows != nil && rows.Next() {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"repitition process definition", err)
			rows.Close()
			return c.DB_EXEC_FAILED
		}
		defer rows.Close()
		dbQuery = "INSERT process_definition SET proc_id=?,proc_name=?,proc_file=?,proc_def=?"
		stmt, err := dbCnt.Sqldb.Prepare(dbQuery)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"insert error", err)
			return c.DB_EXEC_FAILED
		}
		defer stmt.Close()
		buf := bytes.NewBuffer(nil)
		registerGob()
		enc := gob.NewEncoder(buf)
		err1 := enc.Encode(process)
		if err1 != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"gob error", err1)
			return c.DB_EXEC_FAILED
		}
		res, err := stmt.Exec(process.Id, process.Name, file, buf.Bytes())
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"exec error", err)
			return c.DB_EXEC_FAILED
		}
		if result, _ := res.RowsAffected(); result == 1 {
			return c.DB_EXEC_SUCCESS
		} else {
			return c.DB_EXEC_FAILED
		}
	case c.PROCENG_DEPLOY_MODE_UPDATE:
		// check whether processDef exists
		rows, err := dbCnt.Sqldb.Query("SELECT * FROM process_definition WHERE proc_id=?", process.Id)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"query process error with id "+process.Id, err)
			return c.DB_EXEC_FAILED
		}
		defer rows.Close()
		if rows.Next() {
			dbQuery = "UPDATE process_definition SET proc_name=?,proc_file=?,proc_def=? WHERE proc_id=?"
			stmt, err := dbCnt.Sqldb.Prepare(dbQuery)
			if err != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"insert error", err)
				return c.DB_EXEC_FAILED
			}
			defer stmt.Close()
			buf := bytes.NewBuffer(nil)
			registerGob()
			enc := gob.NewEncoder(buf)
			err1 := enc.Encode(process)
			if err1 != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"gob error", err1)
				return c.DB_EXEC_FAILED
			}
			res, err := stmt.Exec(process.Name, file, buf.Bytes(), process.Id)
			if err != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"exec error", err)
				return c.DB_EXEC_FAILED
			}
			if result, _ := res.RowsAffected(); result == 1 {
				return c.DB_EXEC_SUCCESS
			} else {
				return c.DB_EXEC_FAILED
			}
		} else {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"can not find process with id "+process.Id, err)
			return c.DB_EXEC_FAILED
		}
	default:
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"unknown deploy mode in DBconnector", nil)
		return c.DB_EXEC_FAILED
	}
	return c.DB_EXEC_FAILED
}
func (dbCnt *DBConnector) SaveProcDefWithRestrict(file string, process *m.Process, rstct []string) (int, byte) {
	var dbQuery string
	if len(rstct) > 0 {
		for _, v := range rstct {
			if strings.EqualFold(process.PId, v) {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"repitition process definition", nil)
				return 0, c.DB_EXEC_FAILED_RPT
			}
		}
	}
	dbQuery = "INSERT process_definition SET proc_id=?,proc_name=?,proc_desc=?,proc_file=?,proc_def=?,create_time=?"
	stmt, err := dbCnt.Sqldb.Prepare(dbQuery)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"insert error", err)
		return 0, c.DB_EXEC_FAILED
	}
	defer stmt.Close()
	buf := bytes.NewBuffer(nil)
	registerGob()
	enc := gob.NewEncoder(buf)
	err1 := enc.Encode(process)
	if err1 != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"gob error", err1)
		return 0, c.DB_EXEC_FAILED
	}
	res, err := stmt.Exec(process.PId, process.Name, process.Documentation, file, buf.Bytes(), time.Now().Unix())
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"exec error", err)
		return 0, c.DB_EXEC_FAILED
	}
	if result, _ := res.RowsAffected(); result == 1 {
		pid, _ := res.LastInsertId()
		return int(pid), c.DB_EXEC_SUCCESS
	} else {
		return 0, c.DB_EXEC_FAILED
	}

	return 0, c.DB_EXEC_FAILED
}
func (dbCnt *DBConnector) UpdateProcDefById(file string, process *m.Process, id int) byte {

	// check whether processDef exists
	rows, err := dbCnt.Sqldb.Query("SELECT * FROM process_definition WHERE id=?", id)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"query process error with id "+strconv.Itoa(id), err)
		return c.DB_EXEC_FAILED
	}
	if rows.Next() {
		stmt, err := dbCnt.Sqldb.Prepare("UPDATE process_definition SET proc_id=?,proc_name=?,proc_desc=?,proc_file=?,proc_def=? WHERE id=?")
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"insert error", err)
			return c.DB_EXEC_FAILED
		}
		defer stmt.Close()
		buf := bytes.NewBuffer(nil)
		registerGob()
		enc := gob.NewEncoder(buf)
		err1 := enc.Encode(process)
		if err1 != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"gob error", err1)
			return c.DB_EXEC_FAILED
		}
		res, err := stmt.Exec(process.PId, process.Name, process.Documentation, file, buf.Bytes(), id)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"exec error", err)
			return c.DB_EXEC_FAILED
		}
		if result, _ := res.RowsAffected(); result == 1 {
			return c.DB_EXEC_SUCCESS
		} else {
			return c.DB_EXEC_FAILED
		}
	} else {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"can not find process with id "+strconv.Itoa(id), err)
		return c.DB_EXEC_FAILED
	}
	defer rows.Close()
	return c.DB_EXEC_FAILED
}
func registerGob() {
	gob.Register(&m.StartEvent{})
	gob.Register(&m.EndEvent{})
	gob.Register(&m.Transition{})
	gob.Register(&m.ServiceTask{})
	gob.Register(&m.UserTask{})
	gob.Register(&m.ExclusiveGateway{})
	gob.Register(&m.ParallelGateway{})
	gob.Register(&m.FieldExtension{})
	gob.Register(&m.ActivityListener{})
}

func (this *DBConnector) CreateProcessDefById(key int) (*m.Process, byte) {
	proc := new(m.Process)
	var bufInt []uint8
	err := this.Sqldb.QueryRow("SELECT id,proc_def FROM process_definition WHERE id=?", key).Scan(&proc.Id, &bufInt)
	var buf bytes.Buffer
	buf = *bytes.NewBuffer(bufInt)
	switch {
	case err == sql.ErrNoRows:
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"No processDef with that ID.", err)
		return nil, c.DB_EXEC_FAILED
	default:
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"Unknown Error", err)
			return nil, c.DB_EXEC_FAILED
		}
	}
	registerGob()
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&proc)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"Decode Error", err)
		return nil, c.DB_EXEC_FAILED
	}
	return proc, c.DB_EXEC_SUCCESS
}

func (this *DBConnector) QueryProcessById(id int) (*m.ProcessTrs, byte) {
	rows, err := this.Sqldb.Query("SELECT p.id,p.proc_id,p.proc_name,p.proc_desc,p.proc_file,p.create_time,p.multiple_instance,p.is_excutable FROM process_definition AS p  WHERE p.id=? ", id)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"query error", err)
		return nil, c.DB_EXEC_FAILED
	}
	defer rows.Close()
	if rows.Next() {
		proc := new(m.ProcessTrs)
		if err := rows.Scan(&proc.Id, &proc.ProcId, &proc.ProcName, &proc.ProcDesc, &proc.ProcFile, &proc.CreateTime, &proc.IsMultiple, &proc.IsExecutable); err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"scan row error", err)
			return nil, c.DB_EXEC_FAILED
		}
		return proc, c.DB_EXEC_SUCCESS
	} else {
		return nil, c.DB_EXEC_FAILED
	}
}

func (this *DBConnector) CreateProcessDefByKey(key string) (*m.Process, byte) {
	var proc m.Process
	var bufInt []uint8
	err := this.Sqldb.QueryRow("SELECT proc_def FROM process_definition WHERE proc_id=?", key).Scan(&bufInt)
	var buf bytes.Buffer
	buf = *bytes.NewBuffer(bufInt)
	switch {
	case err == sql.ErrNoRows:
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"No processDef with that ID.", err)
		return nil, c.DB_EXEC_FAILED
	default:
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"Unknown Error", err)
			return nil, c.DB_EXEC_FAILED
		}
	}
	registerGob()
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&proc)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"Decode Error", err)
		return nil, c.DB_EXEC_FAILED
	}
	return &proc, c.DB_EXEC_SUCCESS
}
func (this *DBConnector) DeleteProcessById(pid int) byte {
	var file string
	err := this.Sqldb.QueryRow("SELECT proc_file FROM process_definition WHERE id=?", pid).Scan(&file)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"can not find process with id "+strconv.Itoa(pid), err)
		return c.DB_EXEC_FAILED
	}
	stmt, err := this.Sqldb.Prepare("DELETE FROM process_definition WHERE id=?")
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"delete Error", err)
		return c.DB_EXEC_FAILED
	}
	defer stmt.Close()
	if result, err := stmt.Exec(pid); err == nil {
		if _, err := result.RowsAffected(); err == nil {
			// delete file
			cur, _ := os.Getwd()
			filepath := cur + util.GetSp() + c.CONFIG_RES_PATH + file
			err = os.Remove(filepath)
			if err != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"remove file error", err)
				return c.DB_EXEC_FAILED
			}
			path := filepath[:strings.LastIndex(filepath, util.GetSp())]
			err = os.Remove(path)
			if err != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"remove file catalog error", err)
				return c.DB_EXEC_FAILED
			}
			return c.DB_EXEC_SUCCESS
		} else {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"delete Error", err)
			return c.DB_EXEC_FAILED
		}
	} else {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"delete Error", err)
		return c.DB_EXEC_FAILED
	}

}

/*=============== related with user =================*/
func (this *DBConnector) CheckUser(uname string, pwd string) (byte, string) {
	rows, err := this.Sqldb.Query("SELECT * FROM user_info")
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"checkuser query error", err)
		return 1, ""
	}
	defer rows.Close()
	for rows.Next() {
		var uid int
		var username string
		var password string
		err = rows.Scan(&uid, &username, &password)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"checkuser scan error", err)
			return 1, ""
		}
		if strings.EqualFold(username, uname) && strings.EqualFold(password, pwd) {
			return 0, strconv.Itoa(uid)
		}
	}
	return 1, ""
}
func (dbCnt *DBConnector) SaveProcDefByUser(file string, filetype byte, process *m.Process, mode byte, userId string, pId int) byte {
	uId, err := strconv.Atoi(strings.TrimSpace(userId))
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"userId parse error", err)
		return c.DB_EXEC_FAILED
	}
	var dbQuery string
	switch mode {
	case c.PROCENG_DEPLOY_MDOE_DEFAULT:
		// check whether processDef exists
		rows, err := dbCnt.Sqldb.Query("SELECT * FROM process_definition AS p INNER JOIN rlt_user_proc AS r ON r.proc_id=p.id WHERE r.user_id=? AND p.proc_id=?", uId, process.Id)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"insert process query error", err)
			return c.DB_EXEC_FAILED
		}
		if rows != nil && rows.Next() {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"repitition process definition", err)
			rows.Close()
			return c.DB_EXEC_FAILED_RPT
		}
		defer rows.Close()
		dbQuery = "INSERT process_definition SET proc_id=?,proc_name=?,proc_desc=?,proc_file_type=?,proc_file=?,proc_def=?,create_time=?"
		stmt, err := dbCnt.Sqldb.Prepare(dbQuery)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"insert error", err)
			return c.DB_EXEC_FAILED
		}
		defer stmt.Close()
		buf := bytes.NewBuffer(nil)
		registerGob()
		enc := gob.NewEncoder(buf)
		err1 := enc.Encode(process)
		if err1 != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"gob error", err1)
			return c.DB_EXEC_FAILED
		}
		res, err := stmt.Exec(process.PId, process.Name, process.Documentation, filetype, file, buf.Bytes(), time.Now().Unix())
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"exec error", err)
			return c.DB_EXEC_FAILED
		}
		if result, _ := res.RowsAffected(); result == 1 {
			// add user rlt
			pId, err := res.LastInsertId()
			if err != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"can not get process Id after insert", err)
				return c.DB_EXEC_FAILED
			}
			dbQuery = "INSERT rlt_user_proc SET user_id=?,proc_id=?"
			stmtx, err := dbCnt.Sqldb.Prepare(dbQuery)
			if err != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"insert error", err)
				return c.DB_EXEC_FAILED
			}
			defer stmtx.Close()
			resx, err := stmtx.Exec(uId, pId)
			if err != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"exec error", err)
				return c.DB_EXEC_FAILED
			}
			if result, _ := resx.RowsAffected(); result == 1 {
				return c.DB_EXEC_SUCCESS
			} else {
				return c.DB_EXEC_FAILED
			}
			return c.DB_EXEC_SUCCESS
		} else {
			return c.DB_EXEC_FAILED
		}
	case c.PROCENG_DEPLOY_MODE_UPDATE:
		// check whether processDef exists
		//		qStr := "SELECT p.id FROM process_definition AS p INNER JOIN rlt_user_proc AS r ON r.proc_id=p.id WHERE r.user_id=? AND p.proc_id=?"
		rows, err := dbCnt.Sqldb.Query("SELECT id,proc_file FROM process_definition WHERE id=?", pId)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"query process error with id "+process.Id, err)
			return c.DB_EXEC_FAILED
		}
		defer rows.Close()
		if rows.Next() {
			var upId int
			var fPath string
			rows.Scan(&upId, &fPath)
			dbQuery = "UPDATE process_definition SET proc_id=?,proc_name=?,proc_desc=?,proc_file=?,proc_file_type=?,proc_def=? WHERE id=?"
			stmt, err := dbCnt.Sqldb.Prepare(dbQuery)
			if err != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"insert error", err)
				return c.DB_EXEC_FAILED
			}
			defer stmt.Close()
			buf := bytes.NewBuffer(nil)
			registerGob()
			enc := gob.NewEncoder(buf)
			err1 := enc.Encode(process)
			if err1 != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"gob error", err1)
				return c.DB_EXEC_FAILED
			}
			res, err := stmt.Exec(process.PId, process.Name, process.Documentation, file, filetype, buf.Bytes(), upId)
			if err != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"exec error", err)
				return c.DB_EXEC_FAILED
			}
			if result, _ := res.RowsAffected(); result == 1 {
				// delete file
				cur, _ := os.Getwd()
				filepath := cur + util.GetSp() + c.CONFIG_RES_PATH + fPath
				err = os.Remove(filepath)
				if err != nil {
					logger.GetLogger().WriteLog(LOG_DB_ERROR+"remove file error", err)
					return c.DB_EXEC_FAILED
				}
				path := filepath[:strings.LastIndex(filepath, util.GetSp())]
				err = os.Remove(path)
				if err != nil {
					logger.GetLogger().WriteLog(LOG_DB_ERROR+"remove file catalog error", err)
					return c.DB_EXEC_FAILED
				}
				return c.DB_EXEC_SUCCESS
			} else {
				return c.DB_EXEC_FAILED
			}
		} else {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"can not find process with id "+process.Id, err)
			return c.DB_EXEC_FAILED
		}
		defer rows.Close()
	default:
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"unknown deploy mode in DBconnector", nil)
		return c.DB_EXEC_FAILED
	}
	return c.DB_EXEC_FAILED
}

func (dbCnt *DBConnector) QueryProcByInternalUser(userId string) ([]m.Process, byte) {
	uId, err := strconv.Atoi(userId)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"userId parse error", err)
		return nil, c.DB_EXEC_FAILED
	}
	rows, err := dbCnt.Sqldb.Query("SELECT p.id,p.proc_id,p.proc_name,p.proc_file FROM process_definition AS p INNER JOIN rlt_user_proc AS r ON r.proc_id=p.id WHERE r.user_id=? ", uId)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"query error", err)
		return nil, c.DB_EXEC_FAILED
	}
	defer rows.Close()
	var pList []m.Process
	for rows.Next() {
		var proc m.Process
		if err := rows.Scan(&proc.PId, &proc.Id, &proc.Name, &proc.PFile); err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"scan row error", err)
			return nil, c.DB_EXEC_FAILED
		}
		pList = append(pList, proc)
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"scan row error", err)
		return nil, c.DB_EXEC_FAILED
	}
	return pList, c.DB_EXEC_SUCCESS
}
func (dbCnt *DBConnector) QueryProcByUser(uid int) ([]*m.ProcessTrs, byte) {
	listTrs := make([]*m.ProcessTrs, 0)
	rows, err := dbCnt.Sqldb.Query("SELECT p.id,p.proc_id,p.proc_name,p.proc_desc,p.proc_file,p.proc_file_type,p.create_time,p.multiple_instance,p.is_excutable FROM process_definition AS p INNER JOIN rlt_user_proc AS r ON r.proc_id=p.id WHERE r.user_id=? ", uid)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"query error", err)
		return listTrs, c.DB_EXEC_FAILED
	}
	defer rows.Close()
	for rows.Next() {
		proc := new(m.ProcessTrs)
		if err := rows.Scan(&proc.Id, &proc.ProcId, &proc.ProcName, &proc.ProcDesc, &proc.ProcFile, &proc.FileType, &proc.CreateTime, &proc.IsMultiple, &proc.IsExecutable); err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"scan row error", err)
			return listTrs, c.DB_EXEC_FAILED
		}
		proc.Running = false
		listTrs = append(listTrs, proc)
	}
	if err := rows.Err(); err != nil {
		logger.GetLogger().WriteLog(LOG_DB_ERROR+"scan row error", err)
		return listTrs, c.DB_EXEC_FAILED
	}

	return listTrs, c.DB_EXEC_SUCCESS
}
func (dbCnt *DBConnector) QueryProc(list []int) ([]*m.ProcessTrs, byte) {
	listTrs := make([]*m.ProcessTrs, 0)
	for _, v := range list {
		rows, err := dbCnt.Sqldb.Query("SELECT id,proc_id,proc_name,proc_desc,proc_file,create_time,multiple_instance,is_excutable FROM process_definition WHERE id=? ", v)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"query error", err)
			return listTrs, c.DB_EXEC_FAILED
		}
		defer rows.Close()
		for rows.Next() {
			proc := new(m.ProcessTrs)
			if err := rows.Scan(&proc.Id, &proc.ProcId, &proc.ProcName, &proc.ProcDesc, &proc.ProcFile, &proc.CreateTime, &proc.IsMultiple, &proc.IsExecutable); err != nil {
				logger.GetLogger().WriteLog(LOG_DB_ERROR+"scan row error", err)
				return listTrs, c.DB_EXEC_FAILED
			}
			proc.Running = false
			listTrs = append(listTrs, proc)
		}
		if err := rows.Err(); err != nil {
			logger.GetLogger().WriteLog(LOG_DB_ERROR+"scan row error", err)
			return listTrs, c.DB_EXEC_FAILED
		}
	}
	return listTrs, c.DB_EXEC_SUCCESS
}
