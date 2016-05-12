// BpmnjsonConverter
/**
 * 流程转换
 * 流程定义Json模型->类对象模型
 */
package converter

import (
	c "ProcessEngine/src/constant"
	m "ProcessEngine/src/model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

func ConvertJsonToBpmnModel(filename string, filepath string) (interface{}, error) {

	fileJson, err := readJsonFile(filename, filepath)
	if err != nil {
		return nil, err
	}
	process := new(JSProcess)
	mProc := new(m.Process)

	err = json.Unmarshal(fileJson, &process)
	if err != nil {
		return nil, err
	}

	return CreateJSProcess(process, mProc)
}
func readJsonFile(filename string, filepath string) ([]byte, error) {
	if filename == "" {
		return []byte{}, errors.New("converter find filename is null")
	}
	typeArr := []string{"json"}
	strSplit := filename[strings.LastIndex(filename, ".")+1:]

	var flag = false
	for _, strr := range typeArr {
		if strings.EqualFold(strr, strSplit) {
			flag = true
			break
		}
	}
	if flag == false {
		err := errors.New("filetype is not json")
		return []byte{}, err
	}

	file, err := os.Open(filepath + filename)
	defer closeFile(file)
	if err != nil {
		err := errors.New("Json file not found")
		return []byte{}, err
	}
	modelStr, err := ioutil.ReadAll(file)
	if err != nil {
		err := errors.New("Read Json file error")
		return []byte{}, err
	}

	return modelStr, nil
}

func CreateJSProcess(jProc *JSProcess, mProc *m.Process) (*m.Process, error) {

	for _, node := range jProc.Node {
		mProc.PId = node.ProcessId
		mProc.Name = node.ProcessName
		mProc.Documentation = node.ProcessDes
		break
	}

	slice := make([]m.ExecutionElement, 0)
	elementMap := make(map[string]interface{})

	slice = CreateJSStartEvent(jProc, slice, elementMap)

	CreateJSEndEvent(jProc, elementMap)

	CreateJSUserTask(jProc, elementMap)

	CreateJSServiceTask(jProc, elementMap)

	CreateJSGateway(jProc, elementMap)

	// last create transitions
	CreateJSTransition(jProc, elementMap)

	mProc.StartElements = slice
	//FlowElements->activity,transition,event,gateway
	mProc.FlowElements = elementMap
	//ExtensionElements->ActivityListener/FiledExtension
	mProc = initJSProcess(mProc)
	testJSProc(*mProc)

	return mProc, nil
}

func initJSProcess(proc *m.Process) *m.Process {
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
func CreateJSStartEvent(jProcess *JSProcess, slice []m.ExecutionElement, elemmap map[string]interface{}) []m.ExecutionElement {

	for _, node := range jProcess.Node {
		if strings.EqualFold(node.Item, "start") {
			elem := new(m.StartEvent)
			elem.Id = strconv.Itoa(node.Key)
			elem.Name = node.Text
			slice = append(slice, elem)
			elemmap[elem.Id] = elem
			elem.ExecuteState = c.EXECUTION_STATE_SUSPEND
		}
	}
	return slice
}

//endEvent加入Map Process.FlowElements
func CreateJSEndEvent(jProcess *JSProcess, elemmap map[string]interface{}) error {

	for _, node := range jProcess.Node {

		if strings.EqualFold(node.Item, "End") {
			elem := new(m.EndEvent)
			elem.Id = strconv.Itoa(node.Key)
			elem.Name = node.Text
			elemmap[elem.Id] = elem
			elem.ExecuteState = c.EXECUTION_STATE_SUSPEND
		}

	}
	return nil
}

func CreateJSUserTask(jProcess *JSProcess, elemmap map[string]interface{}) {
	for _, uTask := range jProcess.Node {
		if strings.EqualFold(uTask.Category, "eveTask") {
			elem := new(m.UserTask)
			elem.Id = strconv.Itoa(uTask.Key)
			elem.Name = uTask.Text
			elem.ExecuteState = c.EXECUTION_STATE_SUSPEND
			//Documentation json
			docuJson, _ := json.Marshal(uTask.Documentation)
			elem.Documentation = string(docuJson)

			actListener := new(m.ActivityListener)
			actListener.Event = uTask.Event
			actListener.Implementation = uTask.Class
			actListener.ImplementationType = "class"
			//BaseElement->ExtensionElements
			elem.ExtensionElements = append(elem.ExtensionElements, actListener)
			//FlowElement->executionListener
			elem.ExecutionListener = append(elem.ExecutionListener, actListener)

			elemmap[elem.Id] = elem
		}
	}
}

func CreateJSServiceTask(jProcess *JSProcess, elemmap map[string]interface{}) {
	for _, sTask := range jProcess.Node {
		if strings.EqualFold(sTask.Category, "task") || strings.EqualFold(sTask.Category, "exTask") {
			elem := new(m.ServiceTask)
			elem.Id = strconv.Itoa(sTask.Key)
			elem.Name = sTask.Text
			elem.ExecuteState = c.EXECUTION_STATE_SUSPEND
			//Documentation json
			docuJson, _ := json.Marshal(sTask.Documentation)
			elem.Documentation = string(docuJson)

			activitiClass := &m.ExtensionAttribute{Name: "activiti:class", Value: sTask.Class}
			elem.ExtensionAttributes = append(elem.ExtensionAttributes, activitiClass)
			actListener := new(m.ActivityListener)
			actListener.Event = sTask.Event
			actListener.Implementation = sTask.Class
			actListener.ImplementationType = "class"
			//BaseElement->ExtensionElements
			elem.ExtensionElements = append(elem.ExtensionElements, actListener)
			//FlowElement->executionListener
			elem.ExecutionListener = append(elem.ExecutionListener, actListener)

			elemmap[elem.Id] = elem
		}
	}
}

//transition
func CreateJSTransition(jProcess *JSProcess, elemmap map[string]interface{}) error {
	for _, squFlow := range jProcess.LinkData {
		elem := new(m.Transition)
		elem.Id = squFlow.Id
		elem.Name = squFlow.Name
		elem.ConditionExpression = squFlow.ConditionValue
		elem.SourceRef = strconv.Itoa(squFlow.SourceRef)
		elem.TargetRef = strconv.Itoa(squFlow.TargetRef)
		elemmap[elem.SourceRef+elem.TargetRef] = elem
		if squFlow.IsDefault {
			// find gateway, add defaultFlow !!!
			gateway := elemmap[elem.SourceRef]
			if gtw, ok := gateway.(*m.ExclusiveGateway); ok {
				gtw.DefaultFlow = squFlow.Id
			}
		}
	}
	return nil
}

//exclusiveGateway
func CreateJSGateway(jProcess *JSProcess, elemmap map[string]interface{}) error {
	for _, gGateway := range jProcess.Node {
		if strings.EqualFold(gGateway.Category, "exGateway") {
			elem := new(m.ExclusiveGateway)
			elem.Id = strconv.Itoa(gGateway.Key)
			elem.Name = gGateway.Text
			wmx := new(sync.Mutex)
			elem.SetMutex(wmx)
			elem.HasPass = make(map[int64]bool)
			elem.SameRoutineNum = make(map[int64]int)
			elem.ExecuteState = c.EXECUTION_STATE_SUSPEND
			elemmap[elem.Id] = elem
		} else if strings.EqualFold(gGateway.Category, "paGateway") {
			elem := new(m.ParallelGateway)
			elem.Id = strconv.Itoa(gGateway.Key)
			elem.Name = gGateway.Text
			wmx := new(sync.Mutex)
			elem.SetMutex(wmx)
			elem.ExecuteState = c.EXECUTION_STATE_SUSPEND
			elemmap[elem.Id] = elem
		} else {
			fmt.Println("unknown gateway!")
		}
	}
	return nil
}

func testJSProc(process m.Process) error {
	fmt.Println("---------test process-------------")
	fmt.Println("process.name: " + process.Name)
	fmt.Println("process.PId: " + process.PId + "\n")
	for _, v := range process.StartElements {
		fmt.Println("process start:", v.GetElementId())
	}

	for _, v := range process.FlowElements {
		switch obj := v.(type) {
		case *m.StartEvent:
			fmt.Println("====Start Event==== ")
			fmt.Print("Id: ", obj.Id)
			fmt.Print("Name: ", obj.Name)
			fmt.Println("")
		case *m.ServiceTask:
			fmt.Println("====Service task==== ")
			fmt.Print("Id: ", obj.Id)
			fmt.Print(" Documentation: ", obj.Documentation)
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
			fmt.Println("")
		case *m.EndEvent:
			fmt.Println("====End event====")
			fmt.Print("Id: ", obj.Id)
			fmt.Println("")
		case *m.UserTask:
			fmt.Println("====User task====")
			fmt.Print("Id: ", obj.Id)
			fmt.Println(" Documentation: ", obj.Documentation)
			base := obj.ExtensionElements
			for _, v := range base {
				switch o := v.(type) {
				case *m.FieldExtension:
					fmt.Println("FieldExtension "+o.FieldName, " "+o.StringValue)
				case *m.ActivityListener:
					fmt.Println("ActivityListener "+o.Event, " "+o.Implementation)
				}

			}
			fmt.Println("")
		case *m.ExclusiveGateway:
			fmt.Println("====ExclusiveGateway====")
			fmt.Print("Id: ", obj.Id)
			for _, trans := range obj.OutgoingFlows {
				fmt.Println("outgoing: " + trans.SourceRef + " -> " + trans.TargetRef)
			}
			fmt.Println("")
		case *m.Transition:
			fmt.Println("====Transition====")
			fmt.Print("Id: ", obj.Id+" source:"+obj.SourceRef+" target:"+obj.TargetRef+" ConditionExpression: ")
			fmt.Println(obj.ConditionExpression)
			fmt.Println("")
		}
	}
	return nil
}
