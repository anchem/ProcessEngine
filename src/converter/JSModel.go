/**
 * 流程转换
 * 流程定义Json模型
 */
package converter

import ()

type JSProcess struct {
	Node     []JSNodeData   `json:"nodeDataArray"`
	LinkData []JSTransition `json:"linkDataArray"`
}

type JSNodeData struct {
	ProcessId     string                 `json:"process_id"`
	ProcessName   string                 `json:"process_name"`
	ProcessDes    string                 `json:"process_des"`
	Category      string                 `json:"category"`
	Item          string                 `json:"item"`
	Key           int                    `json:"key"`
	Text          string                 `json:"text"`
	EventType     int                    `json:"eventType"`
	TaskType      int                    `json:"taskType"`
	Event         string                 `json:"event"`
	Class         string                 `json:"class"`
	Documentation map[string]interface{} `json:"documentation"`
}

type JSTransition struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	ConditionType  string `json:"condition_type"`
	ConditionValue string `json:"condition"`
	SourceRef      int    `json:"from"`
	TargetRef      int    `json:"to"`
	IsDefault      bool   `json:"isDefault"`
}

type JSGateway struct {
	Id   int    `json:"key"`
	Name string `json:"text"`
}
