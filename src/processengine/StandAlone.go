// ProcessEngine
package processengine

import (
	"File-Uploader/receiver"
	cn "ProcessEngine/src/constant"
	db "ProcessEngine/src/db"
	"ProcessEngine/src/model"
	"ProcessEngine/src/util"
	t "TaskScheduler/task"
	"Website/app/models"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	//	nf "nats-file-transfer/file"
	"nats-cli-go/msg"
	"os"
	"strconv"
)

func serveWebSite() {
	//	msg.InitCfg(model.GetContext().Config.MsgServerAddress)
	cli, err := msg.NewClient(model.GetContext().Config.MsgServerAddress)
	if err != nil {
		panic(err)
	}
	config := model.GetContext().Config
	t.Init(config.TaskDbType, config.TaskDbName, config.TaskDbAddr, config.TaskDbUsername, config.TaskDbPassword, config.LogNatsAddr, config.LogNatsAddr)
	//	event.Init(config.EventDbType, config.EventDbName, config.EventDbAddr, config.EventDbUsername, config.EventDbPassword, config.MsgServerAddress, config.LogNatsAddr)
	//defer nc.Close()
	//	go servAddProc(nc)
	go servProcList(cli)
	go servStartProc(cli)
	go servStopProc(cli)
	go servDeleteProc(cli)
	go servFileRcv()
	go servOlEdit(cli)
	go servOlUpdQry(cli)
	runtime.Goexit()
}
func handleRst(status byte, reason string) []byte {
	rst := new(model.Result)
	switch status {
	case cn.TRANS_RESULT_FAILED:
		rst.Reason = reason
	}
	rst.Res = status
	rstTrs, err := json.Marshal(rst)
	if err != nil {
		panic(err)
	}
	return rstTrs
}

func handleUpd(p model.JsonDef) []byte {
	rstTrs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return rstTrs
}
func servOlUpdQry(cli *msg.Client) {
	cli.Subscribe(cn.TPC_PROC_UPDQRY, func(m []byte) []byte {
		var procJson model.JsonDef
		err := json.Unmarshal(m, &procJson)
		if err != nil {
			panic(err)
		}
		pId, err := strconv.Atoi(procJson.ProcessId)
		if err != nil {
			procJson.TrsReason = fmt.Sprint(err.Error())
			return handleUpd(procJson)
		}
		fileStr, flag := GetProcessEngine().QueryProcessFileById(pId)
		if flag != QUERY_SUCCESS {
			procJson.TrsReason = "Query Failed:" + fileStr
			return handleUpd(procJson)
		}
		procJson.ProcDef = fileStr
		procJson.TrsReason = ""
		return handleUpd(procJson)
	})
}
func servOlEdit(cli *msg.Client) {
	cli.Subscribe(cn.TPC_PROC_OLEDIT, func(m []byte) []byte {
		var procJson model.JsonDef
		err := json.Unmarshal(m, &procJson)
		if err != nil {
			return handleRst(cn.TRANS_RESULT_FAILED, "json unmarshal error")
		}
		var flag byte
		if procJson.ProcessId == "" || procJson.ProcessId == "0" {
			// Add
			// save json to file
			err, fileName := util.SaveUniqueFile(procJson.FileName+".json", []byte(procJson.ProcDef))
			if err != nil {
				return handleRst(cn.TRANS_RESULT_FAILED, "saveFile Error　: "+err.Error())
			}
			flag = deployProcess(fileName, DEPLOY_FILE_JSON, DEPLOY_MODE_DEFAULT, procJson.UserId, 0)
		} else {
			// Update
			pid, err := strconv.Atoi(procJson.ProcessId)
			if err != nil {
				return handleRst(cn.TRANS_RESULT_FAILED, "process id can not be nil :"+err.Error())
			}
			err, fileName := util.SaveUniqueFile(procJson.FileName+".json", []byte(procJson.ProcDef))
			if err != nil {
				return handleRst(cn.TRANS_RESULT_FAILED, "saveFile Error : "+err.Error())
			}
			flag = deployProcess(fileName, DEPLOY_FILE_JSON, DEPLOY_MODE_UPDATE, procJson.UserId, pid)
		}
		switch flag {
		case DEPLOY_SUCCESS:
			return handleRst(cn.TRANS_RESULT_SUCCESS, "")
		case DEPLOY_FAILED_ERROR:
			return handleRst(cn.TRANS_RESULT_FAILED, "deploy failed")
		case DEPLOY_FAILED_REPITITION:
			return handleRst(cn.TRANS_RESULT_FAILED, "deploy failed repitition")
		}
		return handleRst(cn.TRANS_RESULT_FAILED, "unknown error")
	})
}
func servFileRcv() {
	cur, _ := os.Getwd()
	filepath := cur + util.GetSp() + cn.CONFIG_RES_PATH + util.GetSp()
	var call receiver.CallBack
	call = fileSaveMessage
	r, err := receiver.NewFileReceiver(filepath, model.GetContext().Config.FileNats, model.GetContext().Config.FileRcvTopic, &call)
	if err != nil {
		log.Fatal(err)
	}
	go r.ExecuteWithoutResp()
}
func fileSaveMessage(extra map[string]interface{}, path string, name string) {
	cur, _ := os.Getwd()
	filepath := cur + util.GetSp() + cn.CONFIG_RES_PATH + path + util.GetSp() + name
	log.Println(filepath)
	log.Println("===[userID]:", extra["UserId"], "===[ProcessId]:", extra["ProcessId"])
	// add process definition
	if val, ok := extra["UserId"]; ok {
		if uid, ok := val.(string); ok {
			var mode byte
			mode = DEPLOY_MODE_DEFAULT
			var pid int = 0
			if v, ok := extra["ProcessId"]; ok {
				if pv, ok := v.(string); ok {
					pNum, err := strconv.Atoi(pv)
					if err != nil {
						log.Println("Error : cannot parse ProcessId type string")
					} else {
						if pNum > 0 {
							pid = pNum
							mode = DEPLOY_MODE_UPDATE
						}
					}
				}
			}
			deployProcess(path+util.GetSp()+name, DEPLOY_FILE_XML, mode, uid, pid)
		}
	}
}
func deployProcess(filename string, filetype byte, mode byte, userId string, pId int) byte {
	result := GetProcessEngine().DeployWithUser(filename, filetype, mode, userId, pId)
	switch result {
	case DEPLOY_SUCCESS:
		log.Println("deploy process success")
	case DEPLOY_FAILED_ERROR:
		log.Println("deploy process error")
	case DEPLOY_FAILED_REPITITION:
		log.Println("deploy process repitition")
	}
	return result
}
func servDeleteProc(cli *msg.Client) {
	cli.Subscribe(cn.TPC_PROC_DELETE, func(m []byte) []byte {
		var p model.ProcessTrs
		err := json.Unmarshal(m, &p)
		if err != nil {
			panic(err)
		}
		err = GetProcessEngine().DeleteProcessById(p.Id)
		rst := new(model.Result)
		if err != nil {
			rst.Res = cn.TRANS_RESULT_FAILED
			rst.Reason = "delete process failed with ID : " + strconv.Itoa(p.Id)
		}
		rst.Res = cn.TRANS_RESULT_SUCCESS
		rstTrs, err := json.Marshal(rst)
		if err != nil {
			panic(err)
		}
		return rstTrs
	})
}
func servStopProc(cli *msg.Client) {
	cli.Subscribe(cn.TPC_PROC_STOP, func(m []byte) []byte {
		var p model.ProcessTrs
		err := json.Unmarshal(m, &p)
		if err != nil {
			panic(err)
		}
		err = GetProcessEngine().StopProcessById(p.Id)
		rst := new(model.Result)
		if err != nil {
			rst.Res = cn.TRANS_RESULT_FAILED
			rst.Reason = "stop process failed with ID : " + strconv.Itoa(p.Id)
		}
		rst.Res = cn.TRANS_RESULT_SUCCESS
		rstTrs, err := json.Marshal(rst)
		if err != nil {
			panic(err)
		}
		return rstTrs
	})
}
func servStartProc(cli *msg.Client) {
	cli.Subscribe(cn.TPC_PROC_START, func(m []byte) []byte {
		var p model.ProcessTrs
		err := json.Unmarshal(m, &p)
		if err != nil {
			panic(err)
		}
		err = GetProcessEngine().StartProcessById(p.Id)
		rst := new(model.Result)
		if err != nil {
			rst.Res = cn.TRANS_RESULT_FAILED
			rst.Reason = "start process failed with ID : " + strconv.Itoa(p.Id)
		}
		rst.Res = cn.TRANS_RESULT_SUCCESS
		rstTrs, err := json.Marshal(rst)
		if err != nil {
			panic(err)
		}
		return rstTrs
	})
}

func servAddProc(cli *msg.Client) {
	cli.Subscribe(cn.TPC_PROC_ADD, func(m []byte) []byte {
		var rlt []models.ProcUserRlt
		err := json.Unmarshal(m, &rlt)
		if err != nil {
			panic(err)
		}
		log.Println(rlt)
		//		pids := make([]string, 0)
		//		flag := true
		//		for _, v := range p {
		//			id, result := GetProcessEngine().DeployWithRestrict(v.ProcFile, DEPLOY_FILE_XML, DEPLOY_MODE_DEFAULT, pids, 0)
		//			switch result {
		//			case DEPLOY_SUCCESS:
		//				log.Println("deploy process success")
		//			case DEPLOY_FAILED_ERROR:
		//				log.Println("deploy process error")
		//				flag = false
		//			case DEPLOY_FAILED_REPITITION:
		//				log.Println("deploy process repitition")
		//				flag = false
		//			}
		//			if flag == true {
		//			} else {
		//			}
		//		}
		//		listTrs := make([]model.ProcessTrs, len(p))
		//		for k, v := range listTrs {

		//		}
		proc := new(model.ProcessTrs)
		proc.Id = 1
		procTrs, err := json.Marshal(proc)
		if err != nil {
			panic(err)
		}
		return procTrs
	})
}
func servProcList(cli *msg.Client) {
	cli.Subscribe(cn.TPC_PROC_LIST, func(m []byte) []byte {
		var rlt models.ProcUserRlt
		err := json.Unmarshal(m, &rlt)
		if err != nil {
			panic(err)
		}
		log.Println(rlt)
		pl, flag := db.GetDBConnector().QueryProcByUser(rlt.UserId)
		//		plist := make([]int, 0)
		//		var buf bytes.Buffer
		//		buf = *bytes.NewBuffer(msg.Data)
		//		dec := gob.NewDecoder(&buf)
		//		err := dec.Decode(&plist)
		//		if err != nil {
		//			panic(err)
		//		}
		//		log.Println(plist)
		//		pl, flag := db.GetDBConnector().QueryProc(plist)
		//		if flag != cn.DB_EXEC_SUCCESS {
		//			log.Println("process list failed")
		//		}
		for _, v := range pl {
			procPool := model.GetProcInstPool()
			if _, ok := procPool.Pool[v.Id]; ok {
				v.Running = true
			}
		}
		if flag != cn.DB_EXEC_SUCCESS {
			log.Println("process list failed")
		}
		plistTrs, err := json.Marshal(pl)
		if err != nil {
			panic(err)
		}
		log.Println(pl)
		return plistTrs
	})
}
