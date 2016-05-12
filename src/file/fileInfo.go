package file

import (
	"os"
	//	"encoding/json"
)

//message FileSlice {
//	required string file_id = 1;
//	required uint32 queue_index = 2;
//	required bytes slice_content = 3;

//	required bool is_head = 7 [default = false];
//	required bool is_tail = 8 [default = false];
//}

//message FileInfo {
//	required string path = 1;
//	required string name = 2;
//	optional uint32 size = 3;
//}

//type IotCmdBody struct {
//	DeviceId  string    `json:"device_id"`
//	CmdId     int       `json:"cmd_id"`
//	SubDevId  string    `json:"sub_dev_id"`
//	ArgInt32  []int     `json:"arg_int32"`
//	ArgDouble []float64 `json:"arg_double"`
//	ArgString []string  `json:"arg_string"`
//}
type FileShard struct {
	UserId     string `json:"userId"`
	FileId     string `json:"fileId"`
	FileName   string `json:"fileName"`
	FileSize   uint32 `json:"fileSize"`
	ShardCount uint32 `json:"shardCount"`
	ShardIndex int64  `json:"shardIndex"`
	ShardSize  int64  `json:"shardSize"`

	Content []byte `json:"content"`
	//	Content json.RawMessage `json:"content"`
	//	Content string `json:"content"`
}

type FileBuilder struct {
	Writer       *os.File
	FileId       string
	FileName     string
	FileSize     uint32
	ShardCount   uint32
	ShardSize    int64
	Completeness map[int64]int
}
type UserFile struct {
	UserId   string
	FileName string
}
