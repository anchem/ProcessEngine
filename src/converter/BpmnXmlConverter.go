// BpmnXmlConverter
/**
 * 流程转换
 * 流程定义文件->类对象模型
 */
package converter

import (
	m "ProcessEngine/src/model"
	"encoding/xml"
	"errors"
	//	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type CXModel struct {
	XMLName xml.Name  `xml:"definitions"`
	Process CXProcess `xml:"process"`
}

type CXProcess struct {
	XMLName       xml.Name        `xml:"process"`
	Id            string          `xml:"id,attr"`
	Name          string          `xml:"name,attr"`
	IsExecutable  bool            `xml:"isExecutable,attr"`
	Documentation CXDoc           `xml:"documentation"`
	StartEvent    []CXStartEvent  `xml:"startEvent"`
	EndEvent      []CXEndEvent    `xml:"endEvent"`
	Transition    []CXTransition  `xml:"sequenceFlow"`
	UserTask      []CXUserTask    `xml:"userTask"`
	ServiceTask   []CXServiceTask `xml:"serviceTask"`
	ExclGateway   []CXExclGateway `xml:"exclusiveGateway"`
	ParaGateway   []CXParaGateway `xml:"parallelGateway"`
}
type CXDoc struct {
	XMLName xml.Name `xml:"documentation"`
	Doc     string   `xml:",chardata"`
}

type CXStartEvent struct {
	XMLName xml.Name `xml:"startEvent"`
	Id      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type CXEndEvent struct {
	XMLName xml.Name `xml:"endEvent"`
	Id      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type CXTransition struct {
	XMLName             xml.Name              `xml:"sequenceFlow"`
	Id                  string                `xml:"id,attr"`
	Name                string                `xml:"name,attr"`
	SourceRef           string                `xml:"sourceRef,attr"`
	TargetRef           string                `xml:"targetRef,attr"`
	ConditionExpression CXConditionExpression `xml:"conditionExpression"`
}

type CXUserTask struct {
	XMLName           xml.Name            `xml:"userTask"`
	Id                string              `xml:"id,attr"`
	Name              string              `xml:"name,attr"`
	Documentation     CXDoc               `xml:"documentation"`
	ExtensionElements CXExtentionElements `xml:"extensionElements"`
}

type CXServiceTask struct {
	XMLName           xml.Name            `xml:"serviceTask"`
	Id                string              `xml:"id,attr"`
	Name              string              `xml:"name,attr"`
	Documentation     CXDoc               `xml:"documentation"`
	ExtensionElements CXExtentionElements `xml:"extensionElements"`
	ActivityClass     string              `xml:"class,attr"`
}

type CXExclGateway struct {
	XMLName xml.Name `xml:"exclusiveGateway"`
	Id      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type CXParaGateway struct {
	XMLName xml.Name `xml:"parallelGateway"`
	Id      xml.Name `xml:"id,attr"`
	Name    xml.Name `xml:"name,attr"`
}

type CXConditionExpression struct {
	XMLName          xml.Name `xml:"conditionExpression"`
	Xsitype          string   `xml:"type,attr"`
	ConditionExpress string   `xml:",chardata"`
}

type CXExtentionElements struct {
	XMLName          xml.Name             `xml:"extensionElements"`
	FieldExtention   []CXFieldExtention   `xml:"field"`
	ActivityListener []CXActivityListener `xml:"taskListener"`
}

type CXActivityListener struct {
	XMLName        xml.Name           `xml:"taskListener"`
	Event          string             `xml:"event,attr"`
	Class          string             `xml:"class,attr"`
	FieldExtention []CXFieldExtention `xml:"field"`
}

type CXFieldExtention struct {
	XMLName     xml.Name      `xml:"field"`
	FieldName   string        `xml:"name,attr"`
	StringValue CXFieldString `xml:"string"`
}
type CXFieldString struct {
	XMLName xml.Name `xml:"string"`
	Value   string   `xml:",chardata"`
}

func ConvertXmlToBpmnModel(filename string, filepath string) (interface{}, error) {
	typeArr := []string{"xml", "bpmn"}
	strSplit := filename[strings.LastIndex(filename, ".")+1:]

	var flag = false
	for _, strr := range typeArr {
		if strings.EqualFold(strr, strSplit) {
			flag = true
			break
		}
	}
	if flag == false {
		err := errors.New("filetype is not xml")
		return nil, err
	}
	model := CXModel{}
	if strings.EqualFold(filename, "") {
		panic("File name is nil")
	}
	// read file and save
	file, err := os.Open(filepath + filename)
	defer closeFile(file)
	if err != nil {
		panic("File not found!")
	}
	modelStr, err := ioutil.ReadAll(file)
	if err != nil {
		panic("Read file error!")
	}
	err = xml.Unmarshal(modelStr, &model)
	if err != nil {
		panic(err)
	}
	return CreateProcess(&model)
}
func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Println("==============close file error : ", err)
	}
}

func CreateProcess(cModel *CXModel) (*m.Process, error) {

	proc := new(m.Process)
	cxproc := new(CXProcess)
	*cxproc = cModel.Process
	//	fmt.Println("------", string(len(cxproc.StartEvent)))
	slice := []m.ExecutionElement{}
	elementMap := make(map[string]interface{})
	sliceExten := []interface{}{}

	slice = CreateStartEvent(cxproc, slice, elementMap)
	CreateEndEvent(cxproc, elementMap)

	sliceExten = CreateUserTask(cxproc, sliceExten, elementMap)

	sliceExten = CreateServiceTask(cxproc, sliceExten, elementMap)

	CreateExclusiveGateway(cxproc, elementMap)
	// last create transitions
	CreateTransition(cxproc, elementMap)
	proc.PId = cxproc.Id
	proc.Name = cxproc.Name
	proc.Documentation = cxproc.Documentation.Doc
	proc.StartElements = slice
	//FlowElements->activity,transition,event,gateway
	proc.FlowElements = elementMap
	//ExtensionElements->ActivityListener/FiledExtension
	proc.ExtensionElements = sliceExten
	proc = initProcess(proc)
	//	testProc(*proc)
	return proc, nil
}
func initProcess(proc *m.Process) *m.Process {
	for _, element := range proc.FlowElements {
		if trans, ok := element.(*m.Transition); ok {
			if node, ok := proc.FlowElements[trans.SourceRef].(m.ConfigTransition); ok {
				node.AddOutgoingTransitions(trans)
			}
			if node, ok := proc.FlowElements[trans.TargetRef].(m.ConfigTransition); ok {
				node.AddIncommingTransitions(trans)
			}
		}
	}
	return proc
}

//startEvent加入列表Process.StartElements
func CreateStartEvent(cProcess *CXProcess, slice []m.ExecutionElement, elemmap map[string]interface{}) []m.ExecutionElement {
	for _, staEvnt := range cProcess.StartEvent {
		elem := new(m.StartEvent)
		elem.Id = staEvnt.Id
		elem.Name = staEvnt.Name
		slice = append(slice, elem)
		elemmap[elem.Id] = elem
	}

	return slice
}

//endEvent加入Map Process.FlowElements
func CreateEndEvent(cProcess *CXProcess, elemmap map[string]interface{}) error {
	for _, endEvnt := range cProcess.EndEvent {
		elem := new(m.EndEvent)
		elem.Id = endEvnt.Id
		elem.Name = endEvnt.Name
		elemmap[elem.Id] = elem
	}
	return nil
}

func CreateUserTask(cProcess *CXProcess, slice []interface{}, elemmap map[string]interface{}) []interface{} {
	for _, uTask := range cProcess.UserTask {
		elem := new(m.UserTask)
		elem.Id = uTask.Id
		elem.Name = uTask.Name
		elem.Documentation = uTask.Documentation.Doc
		cxExtentionelem := uTask.ExtensionElements
		cxActivityListener := cxExtentionelem.ActivityListener
		for _, listener := range cxActivityListener {
			actListener := new(m.ActivityListener)
			actListener.Event = listener.Event
			actListener.Implementation = listener.Class
			actListener.ImplementationType = "class"
			cxFiledExten := listener.FieldExtention
			for _, field := range cxFiledExten {
				filedEx := new(m.FieldExtension)
				filedEx.FieldName = field.FieldName
				filedEx.StringValue = field.StringValue.Value
				actListener.FieldExtensions = append(actListener.FieldExtensions, filedEx)
			}

			//BaseElement->ExtensionElements
			elem.ExtensionElements = append(elem.ExtensionElements, actListener)
			//FlowElement->executionListener
			elem.ExecutionListener = append(elem.ExecutionListener, actListener)
			elemmap[elem.Id] = elem
			break
		}
		cxFiledExten := cxExtentionelem.FieldExtention
		for _, field := range cxFiledExten {
			filedEx := new(m.FieldExtension)
			filedEx.FieldName = field.FieldName
			filedEx.StringValue = field.StringValue.Value
			elem.ExtensionElements = append(elem.ExtensionElements, filedEx)
		}
		elemmap[elem.Id] = elem
	}
	return slice
}

func CreateServiceTask(cProcess *CXProcess, slice []interface{}, elemmap map[string]interface{}) []interface{} {
	for _, sTask := range cProcess.ServiceTask {
		elem := new(m.ServiceTask)
		elem.Id = sTask.Id
		elem.Name = sTask.Name
		elem.Documentation = sTask.Documentation.Doc
		cxExtentionelem := sTask.ExtensionElements
		cxActivityListener := cxExtentionelem.ActivityListener
		activitiClass := &m.ExtensionAttribute{Name: "activiti:class", Value: sTask.ActivityClass}
		elem.ExtensionAttributes = append(elem.ExtensionAttributes, activitiClass)
		for _, listener := range cxActivityListener {
			actListener := new(m.ActivityListener)
			actListener.Event = listener.Event
			actListener.Implementation = listener.Class
			actListener.ImplementationType = "class"
			//BaseElement->ExtensionElements
			elem.ExtensionElements = append(elem.ExtensionElements, actListener)
			//FlowElement->executionListener
			elem.ExecutionListener = append(elem.ExecutionListener, actListener)
			cxFiledExten := listener.FieldExtention
			for _, field := range cxFiledExten {
				filedEx := new(m.FieldExtension)
				filedEx.FieldName = field.FieldName
				filedEx.StringValue = field.StringValue.Value
				elem.ExtensionElements = append(elem.ExtensionElements, filedEx)
			}
			elemmap[elem.Id] = elem
			break
		}
		cxFiledExten := cxExtentionelem.FieldExtention
		for _, field := range cxFiledExten {
			filedEx := new(m.FieldExtension)
			filedEx.FieldName = field.FieldName
			filedEx.StringValue = field.StringValue.Value
			elem.ExtensionElements = append(elem.ExtensionElements, filedEx)
		}
		elemmap[elem.Id] = elem
	}
	return slice
}

//transition
func CreateTransition(cProcess *CXProcess, elemmap map[string]interface{}) error {
	for _, squFlow := range cProcess.Transition {
		elem := new(m.Transition)
		elem.Id = squFlow.Id
		elem.Name = squFlow.Name
		elem.ConditionExpression = squFlow.ConditionExpression.ConditionExpress
		flag := m.CheckConditionExpression(elem.ConditionExpression)
		if flag == false {
			panic("condition expression Error")
		}
		elem.SourceRef = squFlow.SourceRef
		elem.TargetRef = squFlow.TargetRef
		elemmap[elem.Id] = elem
	}
	return nil
}

//exclusiveGateway
func CreateExclusiveGateway(cProcess *CXProcess, elemmap map[string]interface{}) error {
	for _, exclGateway := range cProcess.ExclGateway {
		elem := new(m.ExclusiveGateway)
		elem.Id = exclGateway.Id
		elem.Name = exclGateway.Name
		elemmap[elem.Id] = elem
	}
	return nil
}

func testProc(process m.Process) error {
	fmt.Println("process.name: " + process.Name)
	fmt.Println("process.Id: " + process.Id)

	for _, v := range process.FlowElements {
		switch obj := v.(type) {
		case *m.StartEvent:
			fmt.Println(obj.Id)
		case *m.ServiceTask:
			fmt.Println(obj.Id)
			fmt.Println(obj.Documentation)
			for _, v := range obj.ExtensionAttributes {
				fmt.Print(" [ extensionAttr:" + v.Name + " Value:" + v.Value + " ]")
			}
			fmt.Println()
			base := obj.ExtensionElements
			for _, v := range base {
				switch o := v.(type) {
				case *m.FieldExtension:
					fmt.Println("FieldExtension "+o.FieldName, " "+o.StringValue)
				case *m.ActivityListener:
					fmt.Println("ActivityListener "+o.Event, " "+o.Implementation)
				}

			}
		case *m.EndEvent:
			fmt.Println(obj.Id)
		case *m.UserTask:
			fmt.Println(obj.Id)
			base := obj.ExtensionElements
			for _, v := range base {
				switch o := v.(type) {
				case *m.FieldExtension:
					fmt.Println("FieldExtension "+o.FieldName, " "+o.StringValue)
				case *m.ActivityListener:
					fmt.Println("ActivityListener "+o.Event, " "+o.Implementation)
					for _, v := range o.ExtensionElements {
						if ext, ok := v.(*m.FieldExtension); ok {
							fmt.Println("	FieldExtension ", ext.FieldName, " "+ext.StringValue)
						}
					}
				}

			}
			for _, v := range obj.OutgoingFlows {
				fmt.Println(v.SourceRef, ":", v.TargetRef)
			}
		case *m.ExclusiveGateway:
			fmt.Println(obj.Id)
			for _, trans := range obj.OutgoingFlows {
				fmt.Println("outgoing: " + trans.SourceRef + " -> " + trans.TargetRef)
			}
		case *m.Transition:
			fmt.Print(obj.Id + " source:" + obj.SourceRef + " target:" + obj.TargetRef + " ConditionExpression: ")
			fmt.Println(obj.ConditionExpression)
		}
	}
	for _, value := range process.StartElements {
		fmt.Println("Start element :", value.GetElementId())
	}
	return nil
}
