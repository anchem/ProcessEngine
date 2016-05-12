// Element
package model

type BaseElement struct {
	Id                  string
	ExtensionElements   []interface{}
	ExtensionAttributes []*ExtensionAttribute
}
type ExtensionAttribute struct {
	Name            string
	Value           string
	NamespacePrefix string
	Namespace       string
}
type ExtensionElement struct {
	BaseElement
}

type FlowElement struct {
	BaseElement
	Name              string
	Documentation     string
	ExecutionListener []*ActivityListener
}
type DataObject struct {
	FlowElement
	ValueType   string
	ValueString string
	Value       interface{}
}
type Transition struct {
	FlowElement
	ConditionExpression string
	SourceRef           string
	TargetRef           string
	SkipExpression      string
}
type FlowNode struct {
	FlowElement
	DataObjects    []*DataObject
	IncommingFlows []*Transition
	OutgoingFlows  []*Transition
	ExecuteState   byte // runtime
}
type ActivityListener struct {
	BaseElement
	Event              string
	ImplementationType string // class
	Implementation     string
	FieldExtensions    []*FieldExtension
}
type FieldExtension struct {
	BaseElement
	FieldName   string
	StringValue string
	Expression  string
}
