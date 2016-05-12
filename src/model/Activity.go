// Activity
package model

import (
	event "EventBus/event"
	c "ProcessEngine/src/constant"
	logger "ProcessEngine/src/logger"
	t "TaskScheduler/task"
	"strconv"
	//	"encoding/json"
	//	"time"
)

type Activity struct {
	FlowNode
	DefaultFlow                  string
	BoundaryEvents               []BoundaryEvent
	FailedJobRetryTimeCycleValue string
}
type SubProcess struct {
	Activity
	FlowElements []FlowElement
	DataObjects  []DataObject
}
type Task struct {
	Activity
}
type ServiceTask struct {
	Task
}
type UserTask struct {
	Task
}

/*------------- ServiceTask ------------------*/
func (this *ServiceTask) AddIncommingTransitions(trans *Transition) {
	this.IncommingFlows = append(this.IncommingFlows, trans)
}
func (this *ServiceTask) AddOutgoingTransitions(trans *Transition) {
	this.OutgoingFlows = append(this.OutgoingFlows, trans)
}
func (this *ServiceTask) GetIncommingTransitions() []*Transition {
	return this.IncommingFlows
}
func (this *ServiceTask) GetOutgoingTransitions() []*Transition {
	return this.OutgoingFlows
}

func (this *ServiceTask) GetElementId() string {
	return this.Id
}
func (this *ServiceTask) Run(procExec *ProcessExecution) error {
	this.ExecuteState = c.EXECUTION_STATE_EXECUTING
	return this.Create(procExec)
}
func (this *ServiceTask) Create(procExec *ProcessExecution) error {
	return this.Execute(procExec)
}
func (this *ServiceTask) Execute(procExec *ProcessExecution) error {
	logger.GetLogger().WriteLog("<ServiceTask> [INFO] : executing ServiceTask "+this.Id, nil)
	err := t.Do(this.Documentation)
	if err != nil {
		return err
	}
	return this.Leave(procExec)
}
func (this *ServiceTask) Leave(procExec *ProcessExecution) error {
	this.ExecuteState = c.EXECUTION_STATE_SUSPEND
	for _, v := range this.OutgoingFlows {
		if node, ok := procExec.ProcInst.ProcDef.FlowElements[v.TargetRef].(ExecutionElement); ok {
			procExec.CurFlowNode = node
		}
	}
	return nil
}

/*------------- UserTask ------------------*/
func (this *UserTask) AddIncommingTransitions(trans *Transition) {
	this.IncommingFlows = append(this.IncommingFlows, trans)
}
func (this *UserTask) AddOutgoingTransitions(trans *Transition) {
	this.OutgoingFlows = append(this.OutgoingFlows, trans)
}
func (this *UserTask) GetIncommingTransitions() []*Transition {
	return this.IncommingFlows
}
func (this *UserTask) GetOutgoingTransitions() []*Transition {
	return this.OutgoingFlows
}
func (this *UserTask) GetElementId() string {
	return this.Id
}
func (this *UserTask) Run(procExec *ProcessExecution) error {
	this.ExecuteState = c.EXECUTION_STATE_EXECUTING
	return this.Create(procExec)
}
func (this *UserTask) Create(procExec *ProcessExecution) error {
	return this.Execute(procExec)
}
func (this *UserTask) Execute(procExec *ProcessExecution) error {
	logger.GetLogger().WriteLog("<UserTask> [INFO] : executing UserTask "+this.Id, nil)
	//	varMap := event.Listen(this.Documentation)
	event.SubscribeOnce(this.Documentation, func(str string) {
		varMap := event.JsonToMap(str)
		logger.GetLogger().WriteLog("Event receive json : "+str, nil)
		for k, v := range varMap {
			_, err := strconv.Atoi(v)
			varType := "string"
			if err != nil {
				logger.GetLogger().WriteLog("<ProcessVar> [INFO] : ", err)
			} else {
				varType = "int"
			}
			if val, ok := procExec.ProcInst.Variables[k]; ok {
				//do something here
				val.ValueString = v
				val.ValueType = varType
				procExec.ProcInst.Variables[k] = val
			} else {
				dataVar := new(DataObject)
				dataVar.Name = k
				dataVar.ValueType = varType
				dataVar.ValueString = v
				procExec.ProcInst.AddProcInstVars(dataVar)
			}
		}
	})
	//	for _, taskListener := range this.ExecutionListener {
	//		task := &ReceiveTypeTask{Event: taskListener.Event, Class: taskListener.Implementation, Fields: taskListener.FieldExtensions}
	//		err := task.DoTask(procExec)
	//		if err != nil {
	//			return err
	//		}
	//	}
	return this.Leave(procExec)
}
func (this *UserTask) Leave(procExec *ProcessExecution) error {
	this.ExecuteState = c.EXECUTION_STATE_SUSPEND
	for _, v := range this.OutgoingFlows {
		if node, ok := procExec.ProcInst.ProcDef.FlowElements[v.TargetRef].(ExecutionElement); ok {
			procExec.CurFlowNode = node
		}
	}
	return nil
}
