// Event
package model

import (
	c "ProcessEngine/src/constant"
	logger "ProcessEngine/src/logger"
	"errors"
	"strconv"
	"strings"
)

type Event struct {
	FlowNode
}
type BoundaryEvent struct {
	Event
	AttachedToRef   Activity
	AttachedToRefId string
	CancelActivity  bool
}
type StartEvent struct {
	Event
}
type EndEvent struct {
	Event
}

/*------------- StartEvent ------------------*/
func (this *StartEvent) AddIncommingTransitions(trans *Transition) {
	this.IncommingFlows = append(this.IncommingFlows, trans)
}
func (this *StartEvent) AddOutgoingTransitions(trans *Transition) {
	this.OutgoingFlows = append(this.OutgoingFlows, trans)
}
func (this *StartEvent) GetIncommingTransitions() []*Transition {
	return this.IncommingFlows
}
func (this *StartEvent) GetOutgoingTransitions() []*Transition {
	return this.OutgoingFlows
}
func (this *StartEvent) GetElementId() string {
	return this.Id
}
func (this *StartEvent) Run(procExec *ProcessExecution) error {
	this.ExecuteState = c.EXECUTION_STATE_EXECUTING
	return this.Create(procExec)
}
func (this *StartEvent) Create(procExec *ProcessExecution) error {
	return this.Execute(procExec)
}
func (this *StartEvent) Execute(procExec *ProcessExecution) error {
	logger.GetLogger().WriteLog("<StartEvent> [INFO] : executing startevent "+this.Id, nil)
	return this.Leave(procExec)
}
func (this *StartEvent) Leave(procExec *ProcessExecution) error {
	for _, v := range this.OutgoingFlows {
		if node, ok := procExec.ProcInst.ProcDef.FlowElements[v.TargetRef].(ExecutionElement); ok {
			procExec.CurFlowNode = node
		}
	}
	// start node cannot point to itself
	if strings.EqualFold(this.Id, procExec.CurFlowNode.GetElementId()) {
		return errors.New("node repitition !")
	}
	this.ExecuteState = c.EXECUTION_STATE_SUSPEND
	return nil
}

/*------------- EndEvent ------------------*/
func (this *EndEvent) AddIncommingTransitions(trans *Transition) {
	this.IncommingFlows = append(this.IncommingFlows, trans)
}
func (this *EndEvent) AddOutgoingTransitions(trans *Transition) {
	this.OutgoingFlows = append(this.OutgoingFlows, trans)
}
func (this *EndEvent) GetIncommingTransitions() []*Transition {
	return this.IncommingFlows
}
func (this *EndEvent) GetOutgoingTransitions() []*Transition {
	return this.OutgoingFlows
}
func (this *EndEvent) GetElementId() string {
	return this.Id
}
func (this *EndEvent) Run(procExec *ProcessExecution) error {
	this.ExecuteState = c.EXECUTION_STATE_EXECUTING
	return this.Create(procExec)
}
func (this *EndEvent) Create(procExec *ProcessExecution) error {
	return this.Execute(procExec)
}
func (this *EndEvent) Execute(procExec *ProcessExecution) error {
	logger.GetLogger().WriteLog("<EndEvent> [INFO] : executing endevent "+this.Id, nil)
	return this.Leave(procExec)
}
func (this *EndEvent) Leave(procExec *ProcessExecution) error {
	// end process execution
	this.ExecuteState = c.EXECUTION_STATE_SUSPEND
	procExec.Finish(c.EXECUTION_FINISH_PROC_END)
	procPool := GetProcInstPool()
	id, err := strconv.Atoi(procExec.ProcInst.ProcDef.Id)
	if err != nil {
		return err
	} else {
		if pInst, ok := procPool.Pool[id]; ok {
			pInst.End()
			procPool.DeleteProcInst(pInst)
		} else {
			return errors.New("can not find process instance with id : " + strconv.Itoa(id))
		}
	}
	return nil
}
