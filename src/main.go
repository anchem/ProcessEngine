// ProcessEnginePlatform project main.go
package test

/*=============== test json converter ===================*/
//import (
//	"converter"
//)

//func main() {
//	jsonData := ""
//	_, err := converter.ConvertJsonToBpmnModel(jsonData)
//	if err != nil {
//		panic(err)
//	}
//}

/*=============== test standalone ===================*/
//import (
//	pe "ProcessEngine/src/processengine"
//	//	"log"
//)

//func main() {
//	engine := pe.GetProcessEngine()
//	engine.Run()
//	//	pids := make([]string, 0) //myProcess
//	//	result := engine.DeployWithRestrict("testProcess.bpmn", pe.DEPLOY_FILE_XML, pe.DEPLOY_MODE_DEFAULT, pids, 0)
//	//	switch result {
//	//	case pe.DEPLOY_SUCCESS:
//	//		log.Println("deploy process success")
//	//	case pe.DEPLOY_FAILED_ERROR:
//	//		log.Println("deploy process error")
//	//	case pe.DEPLOY_FAILED_REPITITION:
//	//		log.Println("deploy process repitition")
//	//	}
//}

/*=============== test string ===================*/
//import (
//	"fmt"
//)

//func main() {
//	plistOrg := make([]int, 0)
//	plistOrg = append(plistOrg, 1)
//	plistOrg = append(plistOrg, 2)
//	fmt.Println(plistOrg)
//}

/*=============== test web app ===================*/
//import (
//	w "webclient"
//)

//func main() {
//	w.ServeHttp()
//}

/*=============== test process Api ===================*/
//import (
//	model "ProcessEngine/src/model"
//	pe "ProcessEngine/src/processengine"
//	"code.google.com/p/goprotobuf/proto"
//	"fmt"
//	"github.com/nsqio/go-nsq"
//	garage "protos/Garage"
//	procDef "protos/ProcessDefinition"
//	"strconv"
//	"strings"
//	"time"
//)

//func main() {
//	procEng := pe.GetProcessEngine()
//	result := procEng.Deploy("GarageDoor.bpmn", pe.DEPLOY_FILE_XML, pe.DEPLOY_MODE_UPDATE)
//	switch result {
//	case pe.DEPLOY_SUCCESS:
//		fmt.Println("deploy process success")
//		//		procEng.StartProcessByKey("GarageDoorId")
//		//		testProcEng()
//	case pe.DEPLOY_FAILED_ERROR:
//		fmt.Println("deploy process error")
//	case pe.DEPLOY_FAILED_REPITITION:
//		fmt.Println("deploy process repitition")
//	}

//}
//func testProcEng() {
//	procEng := pe.GetProcessEngine()
//	fmt.Println(procEng.Name)
//	cfg := nsq.NewConfig()
//	cfg.ClientID = "DefaultNsqClient1"
//	cfg.HeartbeatInterval = time.Second * 30
//	cfg.DialTimeout = time.Second * 30
//	cfg.WriteTimeout = time.Second * 15
//	consumer, err := nsq.NewConsumer("processEngine", "processEngine", cfg)
//	if err != nil {
//		panic("create consumer error!")
//	}
//	handler1 := &Handler{Nsqaddr: "120.76.41.114:4150", Ch: make(chan *procDef.ProcessDefinition, 100)}
//	consumer.AddHandler(handler1)
//	consumer.ConnectToNSQD(handler1.Nsqaddr)
//	defer consumer.Stop()
//	rcv := &procDef.ProcessDefinition{}
//	for {
//		rcv = <-handler1.Ch
//		procEng.StartProcessByKeyWithVars("GarageDoorId", rcv)
//		var pRcvTopic, phoneId, cameraId, garageId string
//		for _, v := range rcv.GetDefItems() {
//			switch {
//			case strings.EqualFold(v.GetName(), "phoneReceiveTopic"):
//				pRcvTopic = v.GetValue()
//			case strings.EqualFold(v.GetName(), "phoneDeviceId"):
//				phoneId = v.GetValue()
//			case strings.EqualFold(v.GetName(), "cameraDeviceId"):
//				cameraId = v.GetValue()
//			case strings.EqualFold(v.GetName(), "garageDeviceId"):
//				garageId = v.GetValue()
//			}
//		}
//		// pacakaging
//		gMsg := new(garage.Garage)
//		gMsg.InstanceId = proto.String("")
//		devIdUint, err := strconv.ParseUint(phoneId, 10, 32)
//		if err != nil {
//			panic("parse devId error")
//		}
//		gMsg.DeviceId = proto.Uint32(uint32(devIdUint))
//		msgCmdUint, err := strconv.ParseUint("06010002", 16, 32)
//		if err != nil {
//			panic("parse msgCmd error")
//		}
//		gMsg.Msgcmd = proto.Uint32(uint32(msgCmdUint))

//		gProcDefRes := new(procDef.ProcessDefResp)
//		gProcDefRes.Success = proto.Bool(true)
//		gProcDefRes.CameraId = proto.String(cameraId)
//		gProcDefRes.GarageId = proto.String(garageId)
//		gProcDefRes.ErrorCode = proto.Uint32(0)
//		gpdRes, err := proto.Marshal(gProcDefRes)
//		if err != nil {
//			panic("marshal processDefinitionRes error")
//		}
//		gMsg.Content = gpdRes

//		gGrg, err := proto.Marshal(gMsg)
//		if err != nil {
//			panic("marshal garage msg error")
//		}
//		producer, err := nsq.NewProducer(model.GetContext().Config.NSQAddress, cfg)
//		producer.Publish(pRcvTopic, gGrg)
//	}
//}

//type Handler struct {
//	Nsqaddr string
//	Ch      chan *procDef.ProcessDefinition
//}

//func (h *Handler) HandleMessage(message *nsq.Message) error {
//	body := &procDef.ProcessDefinition{}
//	err := proto.Unmarshal(message.Body, body)
//	if err != nil {
//		return err
//	}
//	h.Ch <- body
//	return nil
//}

/*=============== test string ===================*/
//import (
//	"fmt"
//	"strings"
//)

//func main() {
//	str := "asdaljs.bpmn.xml"
//	strarr := []string{"json"}
//	strSplit := str[strings.LastIndex(str, ".")+1:]
//	fmt.Println(strings.TrimSpace(strSplit))
//	for _, strr := range strarr {
//		if strings.EqualFold(strr, strSplit) {
//			fmt.Println("ok")
//		}
//	}
//}

/*=============== test interface ===================*/

//type A struct {
//	Arr map[string]interface{}
//}
//type B struct {
//	Id  string
//	Num int
//}

//func (this *B) ChangeNum(num int) {
//	this.Num = num
//}

//type CanChangeNum interface {
//	ChangeNum(num int)
//}

//func main() {
//	a := new(A)
//	ma := make(map[string]interface{})
//	a.Arr = ma
//	b1 := &B{Id: "1", Num: 1}
//	ma["1"] = b1
//	b2 := &B{Id: "2", Num: 2}
//	ma["2"] = b2
//	if value, ok := a.Arr["2"].(CanChangeNum); ok {
//		value.ChangeNum(4)
//	}
//	for _, v := range a.Arr {
//		if value, ok := v.(*B); ok {
//			fmt.Println(value.Id, value.Num)
//		}
//	}
//}
