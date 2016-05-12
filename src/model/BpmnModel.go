// BpmnModel
package model

type BpmnModel struct {
	Proc Process
}
type Process struct {
	BaseElement
	PId           string
	Name          string
	PFile         string
	Documentation string
	FlowElements  map[string]interface{}
	StartElements []ExecutionElement
}
type ConfigTransition interface {
	AddIncommingTransitions(*Transition)
	AddOutgoingTransitions(*Transition)
	GetIncommingTransitions() []*Transition
	GetOutgoingTransitions() []*Transition
}

type ExecutionElement interface {
	GetElementId() string
	Run(*ProcessExecution) error
	Create(*ProcessExecution) error
	Execute(*ProcessExecution) error
	Leave(*ProcessExecution) error
}
