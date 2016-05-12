// ProcessEngine
package processengine

import (
	db "ProcessEngine/src/db"
	logger "ProcessEngine/src/logger"
	model "ProcessEngine/src/model"
	"ProcessEngine/src/util"
	c "constant"
	"converter"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	procDef "protos/ProcessDefinition"
	"strconv"
	"strings"
)

const (
	_ = iota
	DEPLOY_SUCCESS
	DEPLOY_FAILED_REPITITION
	DEPLOY_FAILED_ERROR
	DEPLOY_FILE_XML
	DEPLOY_FILE_JSON
	DEPLOY_MODE_DEFAULT
	DEPLOY_MODE_UPDATE
	QUERY_SUCCESS
	QUERY_FAILED
)
const (
	LOG_ENGINE_ERROR  = "<ProcessEngine> [ERROR] : "
	LOG_ENGINE_WARING = "<ProcessEngine> [WARNING] : "
	LOG_ENGINE_INFO   = "<ProcessEngine> [INFO] : "
)

type ProcessEngine struct {
	Name     string
	deployer *converter.Converter
}

var _engine *ProcessEngine

func GetProcessEngine() *ProcessEngine {
	if _engine == nil {
		_engine = new(ProcessEngine)
		_engine.init()
	}
	return _engine
}
func (this *ProcessEngine) init() {
	this.Name = "DefaultEngine"
	this.deployer = new(converter.Converter)
	model.InitContext()
}

func (this *ProcessEngine) Run() {
	// standalone mode
	serveWebSite()
}

func (this *ProcessEngine) Deploy(filename string, filetype byte, mode byte) byte {
	if strings.EqualFold(filename, "") {
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"func Deploy() receive nil filename", nil)
		return DEPLOY_FAILED_ERROR
	}
	cur, _ := os.Getwd()
	filepath := cur[:strings.LastIndex(cur, util.GetSp())] + util.GetSp() + c.CONFIG_RES_PATH + util.GetSp()
	switch mode {
	case DEPLOY_MODE_DEFAULT:
		var tp byte
		switch filetype {
		case DEPLOY_FILE_XML:
			tp = c.CONVERTER_FILE_TYPE_XML
		case DEPLOY_FILE_JSON:
			tp = c.CONVERTER_FILE_TYPE_JSON
		default:
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"unknown file type", nil)
			return DEPLOY_FAILED_ERROR
		}
		process, err := this.deployer.ConvertToBpmnModel(filename, filepath, tp)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter error", err)
			return DEPLOY_FAILED_ERROR
		}
		if proc, ok := process.(*model.Process); ok {
			result := db.GetDBConnector().SaveProcDef(filename, proc, c.PROCENG_DEPLOY_MDOE_DEFAULT)
			switch result {
			case c.DB_EXEC_SUCCESS:
				return DEPLOY_SUCCESS
			case c.DB_EXEC_FAILED_RPT:
				return DEPLOY_FAILED_REPITITION
			default:
				logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"deploy process db error", nil)
				return DEPLOY_FAILED_ERROR
			}
		}
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter convert wrong type", nil)
		return c.PROCENG_DEPLOY_FAILED
	case DEPLOY_MODE_UPDATE:
		var tp byte
		switch filetype {
		case DEPLOY_FILE_XML:
			tp = c.CONVERTER_FILE_TYPE_XML
		case DEPLOY_FILE_JSON:
			tp = c.CONVERTER_FILE_TYPE_JSON
		default:
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"unknown file type", nil)
			return DEPLOY_FAILED_ERROR
		}
		process, err := this.deployer.ConvertToBpmnModel(filename, filepath, tp)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter error", err)
			return DEPLOY_FAILED_ERROR
		}
		if proc, ok := process.(*model.Process); ok {
			result := db.GetDBConnector().SaveProcDef(filename, proc, c.PROCENG_DEPLOY_MODE_UPDATE)
			switch result {
			case c.DB_EXEC_SUCCESS:
				return DEPLOY_SUCCESS
			default:
				logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"deploy process db error", nil)
				return DEPLOY_FAILED_ERROR
			}
		}
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter convert wrong type", nil)
		return c.PROCENG_DEPLOY_FAILED
	default:
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"unknown deploy mode", nil)
		return DEPLOY_FAILED_ERROR
	}
	return DEPLOY_FAILED_ERROR

}

// ========================= op id ==============================
func (this *ProcessEngine) DeployWithRestrict(filename string, filetype byte, mode byte, pIds []string, id int) (int, byte) {
	if strings.EqualFold(filename, "") {
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"func Deploy() receive nil filename", nil)
		return 0, DEPLOY_FAILED_ERROR
	}
	cur, _ := os.Getwd()
	filepath := cur[:strings.LastIndex(cur, "\\")] + "\\" + c.CONFIG_RES_PATH + "\\"
	switch mode {
	case DEPLOY_MODE_DEFAULT:
		var tp byte
		switch filetype {
		case DEPLOY_FILE_XML:
			tp = c.CONVERTER_FILE_TYPE_XML
		case DEPLOY_FILE_JSON:
			tp = c.CONVERTER_FILE_TYPE_JSON
		default:
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"unknown file type", nil)
			return 0, DEPLOY_FAILED_ERROR
		}
		process, err := this.deployer.ConvertToBpmnModel(filename, filepath, tp)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter error", err)
			return 0, DEPLOY_FAILED_ERROR
		}
		if proc, ok := process.(*model.Process); ok {
			istId, result := db.GetDBConnector().SaveProcDefWithRestrict(filename, proc, pIds)
			switch result {
			case c.DB_EXEC_SUCCESS:
				return istId, DEPLOY_SUCCESS
			case c.DB_EXEC_FAILED_RPT:
				return 0, DEPLOY_FAILED_REPITITION
			default:
				logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"deploy process db error", nil)
				return 0, DEPLOY_FAILED_ERROR
			}
		}
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter convert wrong type", nil)
		return 0, c.PROCENG_DEPLOY_FAILED
	case DEPLOY_MODE_UPDATE:
		var tp byte
		switch filetype {
		case DEPLOY_FILE_XML:
			tp = c.CONVERTER_FILE_TYPE_XML
		case DEPLOY_FILE_JSON:
			tp = c.CONVERTER_FILE_TYPE_JSON
		default:
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"unknown file type", nil)
			return 0, DEPLOY_FAILED_ERROR
		}
		process, err := this.deployer.ConvertToBpmnModel(filename, filepath, tp)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter error", err)
			return 0, DEPLOY_FAILED_ERROR
		}
		if proc, ok := process.(*model.Process); ok {
			result := db.GetDBConnector().UpdateProcDefById(filename, proc, id)
			switch result {
			case c.DB_EXEC_SUCCESS:
				return 0, DEPLOY_SUCCESS
			default:
				logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"deploy process db error", nil)
				return 0, DEPLOY_FAILED_ERROR
			}
		}
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter convert wrong type", nil)
		return 0, c.PROCENG_DEPLOY_FAILED
	default:
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"unknown deploy mode", nil)
		return 0, DEPLOY_FAILED_ERROR
	}
	return 0, DEPLOY_FAILED_ERROR

}
func (this *ProcessEngine) DeleteProcessById(key int) error {
	flag := db.GetDBConnector().DeleteProcessById(key)
	if flag == c.DB_EXEC_FAILED {
		return errors.New("delete process by id failed")
	}
	return nil
}
func (this *ProcessEngine) StopProcessById(key int) error {
	procPool := model.GetProcInstPool()
	if pInst, ok := procPool.Pool[key]; ok {
		pInst.End()
		procPool.DeleteProcInst(pInst)
	} else {
		return errors.New("can not find process instance with id : " + strconv.Itoa(key))
	}
	return nil
}
func (this *ProcessEngine) StartProcessById(key int) error {
	//	if strings.EqualFold(key, "") {
	//		err := errors.New("process key can not be nil")
	//		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"process key is nil", err)
	//		return err
	//	}
	procInst, err := this.createProcessInstanceById(key)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"process instance create failed", err)
		return err
	}
	model.GetProcInstPool().AddProcInst(procInst)
	go procInst.Start()
	return nil
}

// ========================= end op id ==============================
func (this *ProcessEngine) StartProcessByKey(key string) error {
	if strings.EqualFold(key, "") {
		err := errors.New("process key can not be nil")
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"process key is nil", err)
		return err
	}
	procInst, err := this.createProcessInstanceByKey(key)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"process instance create failed", err)
		return err
	}
	model.GetProcInstPool().AddProcInst(procInst)
	go procInst.Start()
	return nil
}
func (this *ProcessEngine) StartProcessByKeyWithVars(key string, itf interface{}) error {
	if strings.EqualFold(key, "") {
		err := errors.New("process key can not be nil")
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"process key is nil", err)
		return err
	}
	procInst, err := this.createProcessInstanceByKey(key)
	if err != nil {
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"process instance create failed", err)
		return err
	}
	// init process with Vars
	if procD, ok := itf.(*procDef.ProcessDefinition); ok {
		if procD.GetIsTemplate() {
			// need proc template
			logger.GetLogger().WriteLog(LOG_ENGINE_INFO+"processDef is "+procD.String(), nil)
			if err := replaceProcessData(procD, &procInst.ProcDef); err != nil {
				logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"replace process data error", err)
				return err
			} else {
				logger.GetLogger().WriteLog(LOG_ENGINE_INFO+"process instance with template created successfuly", nil)
			}
		} else {
			// unuse
			logger.GetLogger().WriteLog(LOG_ENGINE_INFO+"process instance without template", nil)
		}
		model.GetProcInstPool().AddProcInst(procInst)
		go procInst.Start()
		return nil
	} else {
		err := errors.New("var format type error")
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"process instance with vars can't start", err)
		return err
	}
}

func replaceProcessData(procd *procDef.ProcessDefinition, process *model.Process) error {
	for _, elmt := range process.FlowElements {
		switch e := elmt.(type) {
		case *model.ServiceTask:
			for _, ext := range e.ExtensionElements {
				if field, ok := ext.(*model.FieldExtension); ok {
					replaceFieldExtension(field, procd)
				}
			}
		case *model.UserTask:
			for _, tl := range e.ExecutionListener {
				for _, field := range tl.FieldExtensions {
					replaceFieldExtension(field, procd)
				}
			}
		}
	}
	return nil
}
func replaceFieldExtension(field *model.FieldExtension, procd *procDef.ProcessDefinition) {
	for _, item := range procd.GetDefItems() {
		if strings.EqualFold(field.StringValue, item.GetName()) {
			field.StringValue = item.GetValue()
		}
	}
}
func (this *ProcessEngine) createProcessInstanceById(key int) (*model.ProcessInstance, error) {
	//	if strings.EqualFold(key, "") {
	//		err := errors.New("process key can not be nil")
	//		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"process key is nil", err)
	//		return nil, err
	//	}
	// create process definition and bound it to process instance
	procDef, result := db.GetDBConnector().CreateProcessDefById(key)
	switch result {
	case c.DB_EXEC_SUCCESS:
		procInst := new(model.ProcessInstance)
		procInst.ProcDef = *procDef
		return procInst, nil
	case c.DB_EXEC_FAILED:
		return nil, errors.New("process instance create error")
	}
	return nil, errors.New("unknown error")
}
func (this *ProcessEngine) createProcessInstanceByKey(key string) (*model.ProcessInstance, error) {
	if strings.EqualFold(key, "") {
		err := errors.New("process key can not be nil")
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"process key is nil", err)
		return nil, err
	}
	// create process definition and bound it to process instance
	procDef, result := db.GetDBConnector().CreateProcessDefByKey(key)
	switch result {
	case c.DB_EXEC_SUCCESS:
		procInst := new(model.ProcessInstance)
		procInst.ProcDef = *procDef
		return procInst, nil
	case c.DB_EXEC_FAILED:
		return nil, errors.New("process instance create error")
	}
	return nil, errors.New("unknown error")
}

/*=============== related with user =================*/
func (this *ProcessEngine) CheckUser(uname string, pwd string) (byte, string) {
	return db.GetDBConnector().CheckUser(uname, pwd)
}
func (this *ProcessEngine) DeployWithUser(filename string, filetype byte, mode byte, userId string, pId int) byte {
	uId := strings.TrimSpace(userId)
	//	if strings.EqualFold(uId, "") || strings.EqualFold(filename, "") {
	//		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"func DeployWithUser() receive nil filename or userId", nil)
	//		return DEPLOY_FAILED_ERROR
	//	}
	cur, _ := os.Getwd()
	filepath := cur + util.GetSp() + c.CONFIG_RES_PATH
	switch mode {
	case DEPLOY_MODE_DEFAULT:
		var tp byte
		switch filetype {
		case DEPLOY_FILE_XML:
			tp = c.CONVERTER_FILE_TYPE_XML
		case DEPLOY_FILE_JSON:
			tp = c.CONVERTER_FILE_TYPE_JSON
		default:
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"unknown file type", nil)
			return DEPLOY_FAILED_ERROR
		}
		process, err := this.deployer.ConvertToBpmnModel(filename, filepath, tp)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter error", err)
			return DEPLOY_FAILED_ERROR
		}
		if proc, ok := process.(*model.Process); ok {
			result := db.GetDBConnector().SaveProcDefByUser(filename, tp, proc, c.PROCENG_DEPLOY_MDOE_DEFAULT, uId, pId)
			switch result {
			case c.DB_EXEC_SUCCESS:
				return DEPLOY_SUCCESS
			case c.DB_EXEC_FAILED_RPT:
				return DEPLOY_FAILED_REPITITION
			default:
				logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"deploy process db error", nil)
				return DEPLOY_FAILED_ERROR
			}
		}
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter convert wrong type", nil)
		return c.PROCENG_DEPLOY_FAILED
	case DEPLOY_MODE_UPDATE:
		var tp byte
		switch filetype {
		case DEPLOY_FILE_XML:
			tp = c.CONVERTER_FILE_TYPE_XML
		case DEPLOY_FILE_JSON:
			tp = c.CONVERTER_FILE_TYPE_JSON
		default:
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"unknown file type", nil)
			return DEPLOY_FAILED_ERROR
		}
		process, err := this.deployer.ConvertToBpmnModel(filename, filepath, tp)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter error", err)
			return DEPLOY_FAILED_ERROR
		}
		if proc, ok := process.(*model.Process); ok {
			result := db.GetDBConnector().SaveProcDefByUser(filename, tp, proc, c.PROCENG_DEPLOY_MODE_UPDATE, uId, pId)
			switch result {
			case c.DB_EXEC_SUCCESS:
				return DEPLOY_SUCCESS
			default:
				logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"deploy process db error", nil)
				return DEPLOY_FAILED_ERROR
			}
		}
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"converter convert wrong type", nil)
		return c.PROCENG_DEPLOY_FAILED
	default:
		logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"unknown deploy mode", nil)
		return DEPLOY_FAILED_ERROR
	}
	return DEPLOY_FAILED_ERROR
}

func (this *ProcessEngine) QueryProcessFileById(pId int) (string, byte) {
	proc, flag := db.GetDBConnector().QueryProcessById(pId)
	if flag == c.DB_EXEC_SUCCESS {
		// read file
		cur, _ := os.Getwd()
		filepath := cur + util.GetSp() + c.CONFIG_RES_PATH
		file, err := os.Open(filepath + proc.ProcFile)
		defer file.Close()
		if err != nil {
			return "Open file error" + fmt.Sprint(err.Error()), QUERY_FAILED
		}
		modelStr, err := ioutil.ReadAll(file)
		if err != nil {
			return "Read file error" + fmt.Sprint(err.Error()), QUERY_FAILED
		}
		return string(modelStr), QUERY_SUCCESS
	} else {
		return "query process file db execution failed", QUERY_FAILED
	}
}

type ProcDefFront struct {
	Id       string
	Key      string
	ProcName string
	ProcFile string
}

func (this *ProcessEngine) QueryProcessDefByUser(userId string) (string, byte) {
	// success : "[]" or "[{\"key\":\"process id\",\"procName\":\"procName\",\"procFile\":\"ok\"}]"
	// failed : "reason"
	procList, flag := db.GetDBConnector().QueryProcByInternalUser(userId)
	switch flag {
	case c.DB_EXEC_SUCCESS:
		var pList []ProcDefFront
		for _, value := range procList {
			p := &ProcDefFront{Id: value.PId, Key: value.Id, ProcName: value.Name, ProcFile: value.PFile}
			pList = append(pList, *p)
		}
		b, err := json.Marshal(pList)
		if err != nil {
			logger.GetLogger().WriteLog(LOG_ENGINE_ERROR+"QueryProcessDefByUser() marshal json error", err)
			return "marshal json error", QUERY_FAILED
		}
		return string(b), QUERY_SUCCESS
	case c.DB_EXEC_FAILED:
		return "server db query failed", QUERY_FAILED
	default:
		return "server unknown dbresult", QUERY_FAILED
	}
	return "error", QUERY_FAILED
}
