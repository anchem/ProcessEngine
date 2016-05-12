// Execution
package model

import (
	c "ProcessEngine/src/constant"
	logger "ProcessEngine/src/logger"
)

type ProcessExecution struct {
	Id           uint64
	ProcInst     *ProcessInstance
	ExecuteState byte
	IsExecutable bool
	CreateTime   int64
	CurFlowNode  ExecutionElement
}

func (this *ProcessExecution) Run() {
	this.ExecuteState = c.EXECUTION_RUN_STATE_EXECUTING
	for this.IsExecutable {
		err := this.CurFlowNode.Run(this)
		if err != nil {
			logger.GetLogger().WriteLog("<ProcessExecution> [ERROR] : runtime error", err)
			this.IsExecutable = false
			this.Finish(c.EXECUTION_RUN_STATE_ERROR_STOP)
		}
	}

}
func (this *ProcessExecution) Finish(flag byte) {
	this.IsExecutable = false
	this.ExecuteState = c.EXECUTION_RUN_STATE_DONE
	msg := &ProcMsg{IdFromExecution: this.Id}
	switch flag {
	case c.EXECUTION_RUN_STATE_ERROR_STOP:
		msg.Cmd = c.EXECUTION_FINISH_ERROR_STOP
	case c.EXECUTION_FINISH_DONE:
		msg.Cmd = c.PROC_EXEC_MSG_TYPE_DONE
	case c.EXECUTION_FINISH_PROC_END:
		msg.Cmd = c.PROC_EXEC_MSG_TYPE_END
	case c.EXECUTION_FINISH_FORCE_STOP:
		msg.Cmd = c.PROC_EXEC_MSG_TYPE_FORCE_STOP
	default:
		logger.GetLogger().WriteLog("<ProcessExecution> [ERROR] : unknown flag in func Finish()", nil)
	}
	this.ProcInst.Ch <- msg
}
