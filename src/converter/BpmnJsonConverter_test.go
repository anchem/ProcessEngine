package converter

import (
//	"fmt"
	"testing"
)

func Test_ImmediateTask(t *testing.T){
var Testjson = `{ "class": "go.GraphLinksModel",
  "nodeDataArray": [ 
{"process_id":0, "process_name":"bpmn", "process_des":"", "key":101, "category":"event", "text":"Start", "eventType":1, "eventDimension":1, "item":"start", "loc":"189.984375 186"},
{"process_id":0, "process_name":"bpmn", "process_des":"", "key":132, "category":"activity", "text":"User Task", "taskType":2, "item":"User task", "documentation":{"task_id":"132", "task_name":"", "group_id":"immediate_g", "description":"���ͳ���״̬��ѯ��������", "user_id":"1", "user_pass":"159357qw", "type":"", "cron":{"placeholder":"30 2 1 20 10 *", "type":"text", "id":"cron"}, "command":"", "AppDid":"xxxxxxxx", "AppCore":2, "CoreCode":"8888888888", "body":{"DeviceId":"garageDeviceId", "CmdId":903503, "SubDevId":"", "ArgInt32":[ 255,255,255 ], "ArgDouble":null, "ArgString":null, "ArgByte":null}}, "loc":"462.984375 195", "boundaryEventArray":[ {"portId":"be0", "eventType":3, "eventDimension":5, "color":"white", "alignmentIndex":0},{"portId":"be1", "eventType":10, "eventDimension":5, "color":"white", "alignmentIndex":1} ]},
{"process_id":0, "process_name":"bpmn", "process_des":"", "key":133, "category":"activity", "text":"Service\nTask", "taskType":6, "item":"service task", "documentation":{"task_id":"134", "task_name":"Service Task", "group_id":"immediate_g", "description":"���ͳ���״̬��ѯ��������", "user_id":"1", "user_pass":"159357qw", "type":"immediate", "cron":"30 2 1 20 10 *", "command":"0x02010003", "AppDid":"xxxxxxxx", "AppCore":2, "CoreCode":"8888888888", "body":{"DeviceId":"garageDeviceId", "CmdId":903503, "SubDevId":"", "ArgInt32":[ 255,255,255 ], "ArgDouble":null, "ArgString":null, "ArgByte":null}}, "loc":"702.984375 204", "boundaryEventArray":[ {"portId":"be0", "eventType":2, "eventDimension":5, "color":"white", "alignmentIndex":0} ]},
{"process_id":0, "process_name":"bpmn", "process_des":"", "key":134, "category":"activity", "text":"Script Task", "taskType":4, "isSequential":true, "item":"Script Task", "documentation":{"task_id":"134", "task_name":"Script Task", "group_id":"immediate_g", "description":"���ͳ���״̬��ѯ��������", "user_id":"1", "user_pass":"159357qw", "type":"immediate", "cron":"30 2 1 20 10 *", "command":"0x02010003", "AppDid":"xxxxxxxx", "AppCore":2, "CoreCode":"8888888888", "body":{"DeviceId":"garageDeviceId", "CmdId":903503, "SubDevId":"", "ArgInt32":[ 255,255,255 ], "ArgDouble":null, "ArgString":null, "ArgByte":null}}, "loc":"875.984375 207", "boundaryEventArray":[ {"portId":"be0", "eventType":7, "eventDimension":5, "color":"white", "alignmentIndex":0} ]},
{"process_id":0, "process_name":"bpmn", "process_des":"", "key":201, "category":"gateway", "text":"Parallel", "gatewayType":1, "loc":"596.984375 62"},
{"process_id":0, "process_name":"bpmn", "process_des":"", "key":204, "category":"gateway", "text":"Exclusive", "gatewayType":4, "loc":"585.984375 328"}
 ],
  "linkDataArray": [ 
{"from":101, "to":132, "points":[211.484375,186,221.484375,186,302.234375,186,302.234375,195,382.984375,195,402.984375,195]},
{"from":132, "to":133, "points":[522.984375,195,532.984375,195,577.984375,195,577.984375,204,622.984375,204,642.984375,204]},
{"from":133, "to":134, "points":[762.984375,204,772.984375,204,784.484375,204,784.484375,207,795.984375,207,815.984375,207]}
 ]}`

	ConvertJsonToBpmnModel(Testjson)
	
	
//	resultMap := make(map[string]string)
//	event.Init("mysql", "event_bus", "127.0.0.1:3306", "root", "root", "nats://127.0.0.1:4222")
//	resultMap = event.Listen(Testjson)
	
//	fmt.Println("resultMap: ",resultMap)
	
//	for k, v := range resultMap {  
//    fmt.Println(k, v)  
//} 
}