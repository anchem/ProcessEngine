// Gateway
package model

import (
	c "ProcessEngine/src/constant"
	logger "ProcessEngine/src/logger"
	"errors"
	"strings"
	"sync"
	"time"
)

type Gateway struct {
	FlowNode
	wmutex      sync.Mutex
	DefaultFlow string
}
type ExclusiveGateway struct {
	HasPass        map[int64]bool
	SameRoutineNum map[int64]int
	Gateway
}
type ParallelGateway struct {
	Gateway
	NumOfDone int // runtime
}

/*------------- ExclusiveGateway ------------------*/
func (this *ExclusiveGateway) SetMutex(m *sync.Mutex) {
	this.wmutex = *m
}
func (this *ExclusiveGateway) AddIncommingTransitions(trans *Transition) {
	this.IncommingFlows = append(this.IncommingFlows, trans)
}
func (this *ExclusiveGateway) AddOutgoingTransitions(trans *Transition) {
	this.OutgoingFlows = append(this.OutgoingFlows, trans)
}
func (this *ExclusiveGateway) GetIncommingTransitions() []*Transition {
	return this.IncommingFlows
}
func (this *ExclusiveGateway) GetOutgoingTransitions() []*Transition {
	return this.OutgoingFlows
}
func (this *ExclusiveGateway) GetElementId() string {
	return this.Id
}
func (this *ExclusiveGateway) Run(procExec *ProcessExecution) error {
	this.ExecuteState = c.EXECUTION_STATE_EXECUTING
	if len(this.IncommingFlows) > 1 {
		// run only one routine
		this.wmutex.Lock()
		execTime := procExec.CreateTime
		if val, ok := this.HasPass[execTime]; ok {
			if val {
				// stop this execution
				this.SameRoutineNum[execTime]--
				if this.SameRoutineNum[execTime] < 1 {
					delete(this.HasPass, execTime)
					delete(this.SameRoutineNum, execTime)
				}
				procExec.Finish(c.EXECUTION_FINISH_DONE)
			} else {
				return errors.New("unknown process pass status")
			}
		} else {
			this.HasPass[execTime] = true
			num := 0
			for _, exec := range procExec.ProcInst.Executions {
				if exec.CreateTime == procExec.CreateTime {
					num++
				}
			}
			num--
			this.SameRoutineNum[execTime] = num
		}
		this.wmutex.Unlock()
	}
	return this.Create(procExec)
}
func (this *ExclusiveGateway) Create(procExec *ProcessExecution) error {
	return this.Execute(procExec)
}
func (this *ExclusiveGateway) Execute(procExec *ProcessExecution) error {
	logger.GetLogger().WriteLog("<ExclusiveGateway> [INFO] : executing ExclusiveGateway "+this.Id, nil)
	return this.Leave(procExec)
}
func (this *ExclusiveGateway) Leave(procExec *ProcessExecution) error {
	this.ExecuteState = c.EXECUTION_STATE_SUSPEND
	if len(this.OutgoingFlows) == 1 {
		if node, ok := procExec.ProcInst.ProcDef.FlowElements[this.OutgoingFlows[0].TargetRef].(ExecutionElement); ok {
			procExec.CurFlowNode = node
		}
	} else if len(this.OutgoingFlows) > 1 {
	CHECK:
		for _, transition := range this.OutgoingFlows {
			if CalculateConditionExpression(transition.ConditionExpression, procExec.ProcInst.Variables) {
				if node, ok := procExec.ProcInst.ProcDef.FlowElements[transition.TargetRef].(ExecutionElement); ok {
					procExec.CurFlowNode = node
					break CHECK
				}
			}
		}
		// no condition expression is true, do default flow
		if strings.EqualFold(procExec.CurFlowNode.GetElementId(), this.Id) {
			if !strings.EqualFold(this.DefaultFlow, "") {
				for _, trans := range this.OutgoingFlows {
					if strings.EqualFold(trans.Id, this.DefaultFlow) {
						if node, ok := procExec.ProcInst.ProcDef.FlowElements[trans.TargetRef].(ExecutionElement); ok {
							procExec.CurFlowNode = node
						}
					}
				}
			}
		}
	} else {
		return errors.New("no outgoing flows !")
	}
	if strings.EqualFold(procExec.CurFlowNode.GetElementId(), this.Id) {
		return errors.New("node repitition !")
	}
	return nil
}

/*------------- ParallelGateway ------------------*/
func (this *ParallelGateway) SetMutex(m *sync.Mutex) {
	this.wmutex = *m
}
func (this *ParallelGateway) AddIncommingTransitions(trans *Transition) {
	this.IncommingFlows = append(this.IncommingFlows, trans)
}
func (this *ParallelGateway) AddOutgoingTransitions(trans *Transition) {
	this.OutgoingFlows = append(this.OutgoingFlows, trans)
}
func (this *ParallelGateway) GetIncommingTransitions() []*Transition {
	return this.IncommingFlows
}
func (this *ParallelGateway) GetOutgoingTransitions() []*Transition {
	return this.OutgoingFlows
}
func (this *ParallelGateway) GetElementId() string {
	return this.Id
}
func (this *ParallelGateway) Run(procExec *ProcessExecution) error {
	switch this.ExecuteState {
	case c.EXECUTION_STATE_EXECUTING:
		this.wmutex.Lock()
		this.NumOfDone++
		this.wmutex.Unlock()
	case c.EXECUTION_STATE_SUSPEND:
		this.wmutex.Lock()
		this.ExecuteState = c.EXECUTION_STATE_EXECUTING
		this.NumOfDone = 0
		this.NumOfDone++
		this.wmutex.Unlock()
	default:
		return errors.New("unknown ExecuteState")
	}
	return this.Create(procExec)
}
func (this *ParallelGateway) Create(procExec *ProcessExecution) error {
	return this.Execute(procExec)
}
func (this *ParallelGateway) Execute(procExec *ProcessExecution) error {
	logger.GetLogger().WriteLog("<ParallelGateway> [INFO] : executing parallel gateway "+this.Id, nil)

	if len(this.IncommingFlows) == this.NumOfDone {
		return this.Leave(procExec)
	} else {
		procExec.Finish(c.EXECUTION_FINISH_DONE)
		return nil
	}
}
func (this *ParallelGateway) Leave(procExec *ProcessExecution) error {
	this.ExecuteState = c.EXECUTION_STATE_SUSPEND
	count := 0
	curTime := time.Now().Unix()
	for _, transition := range this.OutgoingFlows {
		count++
		if count < 2 {
			if node, ok := procExec.ProcInst.ProcDef.FlowElements[transition.TargetRef].(ExecutionElement); ok {
				procExec.CurFlowNode = node
				procExec.CreateTime = curTime
			}
		} else {
			// add new go_routine to run
			if node, ok := procExec.ProcInst.ProcDef.FlowElements[transition.TargetRef].(ExecutionElement); ok {
				newExec := &ProcessExecution{CurFlowNode: node, CreateTime: curTime}
				msg := &ProcMsg{Cmd: c.PROC_EXEC_MSG_TYPE_ADD, Param: newExec, IdFromExecution: procExec.Id}
				procExec.ProcInst.Ch <- msg
			}
		}
	}
	return nil
}
