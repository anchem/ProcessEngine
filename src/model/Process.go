// Process
package model

import (
	c "ProcessEngine/src/constant"
	logger "ProcessEngine/src/logger"
	"strconv"
	"strings"
	"time"
)

type ProcessInstance struct {
	ProcDef    Process
	Count      uint64
	Executions map[uint64]*ProcessExecution
	Variables  map[string]*DataObject
	//	Transitions map[string]*Transition
	Ch           chan *ProcMsg
	IsExecutable bool
}
type ProcMsg struct {
	IdFromExecution uint64
	Cmd             byte
	Param           interface{}
}

func (this *ProcessInstance) initial() {
	this.Ch = make(chan *ProcMsg, 50)
	this.Executions = make(map[uint64]*ProcessExecution)
	this.Count = 0
	this.IsExecutable = true
	this.Variables = make(map[string]*DataObject)
}
func (this *ProcessInstance) AddProcInstVars(dataVar *DataObject) {
	if dataVar == nil || strings.EqualFold(dataVar.Name, "") {
		logger.GetLogger().WriteLog("<ProcessInstance> [ERROR] : add a nil var", nil)
	} else {
		this.Variables[dataVar.Name] = dataVar
	}
}
func (this *ProcessInstance) DeleteProcInstVars(dataVar *DataObject) {
	delete(this.Variables, dataVar.Name)
}

func (this *ProcessInstance) Start() {
	this.initial()
	crtTime := time.Now().Unix()
	for _, node := range this.ProcDef.StartElements {
		this.Count++
		exec := &ProcessExecution{Id: this.Count, ProcInst: this, CurFlowNode: node, IsExecutable: true, CreateTime: crtTime, ExecuteState: c.EXECUTION_RUN_STATE_INIT}
		this.Executions[this.Count] = exec
		go exec.Run()
	}
	var msg *ProcMsg
	for this.IsExecutable {
		msg = <-this.Ch
		switch msg.Cmd {
		case c.PROC_EXEC_MSG_TYPE_DONE:
			// finish execution
			delete(this.Executions, msg.IdFromExecution)
		case c.PROC_EXEC_MSG_TYPE_END:
			// end process instance
			delete(this.Executions, msg.IdFromExecution)
			this.End()
		case c.PROC_EXEC_MSG_TYPE_FORCE_STOP:
			delete(this.Executions, msg.IdFromExecution)
		case c.EXECUTION_FINISH_ERROR_STOP:
			logger.GetLogger().WriteLog("<ProcessInstance> [Error] : execution "+strconv.FormatUint(msg.IdFromExecution, 10)+" has occur an error", nil)
			delete(this.Executions, msg.IdFromExecution)
		case c.PROC_EXEC_MSG_TYPE_ADD:
			// add execution
			if exec, ok := msg.Param.(*ProcessExecution); ok {
				this.Count++
				exec.ProcInst = this
				exec.Id = this.Count
				exec.ExecuteState = c.EXECUTION_RUN_STATE_INIT
				exec.IsExecutable = true
				this.Executions[this.Count] = exec
				go exec.Run()
			}
		default:
			logger.GetLogger().WriteLog("<ProcessInstance> [ERROR] : unknown msg cmd in func Start()", nil)
		}
		if len(this.Executions) > 0 {

		} else {
			logger.GetLogger().WriteLog("<ProcessInstance> [INFO] : process instance will stop because it has no execution", nil)
			this.IsExecutable = false
			GetProcInstPool().DeleteProcInst(this)
		}
	}

}
func (this *ProcessInstance) End() {
	for _, exec := range this.Executions {
		exec.Finish(c.EXECUTION_FINISH_FORCE_STOP)
	}
	this.IsExecutable = false
	logger.GetLogger().WriteLog("<ProcessInstance> [INFO] : processInstance ended", nil)
}
