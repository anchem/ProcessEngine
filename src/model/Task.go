// Task
package model

import (
	cmd "ProcessEngine/src/command"
	"errors"
	"fmt"
	//	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"
	garage "ProcessEngine/src/protos/Garage"
	procCmd "ProcessEngine/src/protos/ProcCmd"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	Nsqaddr string
	Ch      chan *garage.OpMap
}

func (h *Handler) HandleMessage(message *nsq.Message) error {
	body := &garage.OpMap{}
	err := proto.Unmarshal(message.Body, body)
	if err != nil {
		fmt.Println("handler error->", err)
		return nil
	}
	h.Ch <- body
	return nil
}

type SendTypeTask struct {
	Class  string
	Fields []interface{}
}
type ReceiveTypeTask struct {
	Event  string
	Class  string
	Fields []*FieldExtension
}

func (this *SendTypeTask) DoTask(procExec *ProcessExecution) error {
	cfg := nsq.NewConfig()
	cfg.ClientID = "DefaultNsqClient1"
	cfg.HeartbeatInterval = time.Second * 30
	cfg.DialTimeout = time.Second * 30
	cfg.WriteTimeout = time.Second * 15
	producer, err := nsq.NewProducer(GetContext().Config.MsgServerAddress, cfg)
	if err != nil {
		return err
	}
	defer producer.Stop()
	var topic string = ""
	pCmd := new(procCmd.ProcCmd)
	for _, value := range this.Fields {
		if field, ok := value.(*FieldExtension); ok {
			switch {
			case strings.EqualFold(field.FieldName, "topic"):
				topic = field.StringValue
			case strings.EqualFold(field.FieldName, "deviceId"):
				deviceId, err := strconv.ParseUint(strings.TrimSpace(field.StringValue), 10, 32)
				if err != nil {
					return err
				}
				pCmd.DeviceId = proto.Uint32(uint32(deviceId))
			case strings.EqualFold(field.FieldName, "commandCode"):
				codeString := field.StringValue[strings.Index(field.FieldName, "0x")+3:]
				opCode, err := strconv.ParseUint(strings.TrimSpace(codeString), 16, 32)
				if err != nil {
					return err
				}
				pCmd.OperationCode = proto.Uint32(uint32(opCode))
			}
		}
	}
	gCtt, err := proto.Marshal(pCmd)
	if err != nil {
		return err
	}
	gMsg := new(garage.Garage)
	gMsg.InstanceId = proto.String("defaultId")
	gMsg.DeviceId = proto.Uint32(666)
	gMsg.Msgcmd = proto.Uint32(cmd.DEVSIM_IN_MESSAGE)
	gMsg.Content = gCtt
	gProto, err := proto.Marshal(gMsg)
	if err != nil {
		return err
	}
	if strings.EqualFold(topic, "") {
		return errors.New("field has no topic")
	} else {
		err = producer.Publish(topic, gProto)
	}
	if err != nil {
		return err
	}
	return nil
}
func (this *ReceiveTypeTask) DoTask(procExec *ProcessExecution) error {
	//	dataVar := &DataObject{ValueString: "true"}
	//	dataVar.Name = "valid"
	//	procExec.ProcInst.AddProcInstVars(dataVar)
	//	var topic string = ""
	//	for _, field := range this.Fields {
	//		switch {
	//		case strings.EqualFold(field.FieldName, "topic"):
	//			topic = field.StringValue
	//		}
	//	}
	//	if strings.EqualFold(topic, "") {
	//		return errors.New("")
	//	} else {
	//		cfg := nsq.NewConfig()
	//		cfg.ClientID = "DefaultNsqClient1"
	//		cfg.HeartbeatInterval = time.Second * 30
	//		cfg.DialTimeout = time.Second * 30
	//		cfg.WriteTimeout = time.Second * 15
	//		fmt.Println("topic->", topic)
	//		consumer, err := nsq.NewConsumer(topic, topic, cfg)
	//		if err != nil {
	//			return err
	//		}
	//		handler1 := &Handler{Nsqaddr: "120.76.41.114:4150", Ch: make(chan *garage.OpMap, 100)}
	//		consumer.AddHandler(handler1)
	//		consumer.ConnectToNSQD(handler1.Nsqaddr)
	//		defer consumer.Stop()
	//		rcv := &garage.OpMap{}
	//	RECEIVE:
	//		for {
	//			rcv = <-handler1.Ch
	//			//			vrMap := &garage.OpMap{}
	//			//			proto.Unmarshal(rcv.GetContent(), vrMap)
	//			fmt.Println(rcv.Values)
	//			for _, item := range rcv.GetValues() {
	//				dataVar := &DataObject{ValueString: item.GetValue()}
	//				dataVar.Name = item.GetName()
	//				procExec.ProcInst.AddProcInstVars(dataVar)
	//			}
	//			break RECEIVE
	//		}
	//	}
	return nil
}
