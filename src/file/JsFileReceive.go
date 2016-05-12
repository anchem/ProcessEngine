package file

import (
	"errors"
	"fmt"
	//	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"strings"
	//	_file "protos/File"
	"encoding/json"
	"time"
)

//var mReceivers map[string]*FileBuilder
//func init(){
//	Topics = make(map[uint32]string)
//}
//fileReceiver singleton
var _fileReceiver *fileReceiver

func NewFileReceiver() *fileReceiver {
	if _fileReceiver == nil {
		_fileReceiver = new(fileReceiver)
		_fileReceiver.CPLch = make(chan *UserFile, 100)
	}
	return _fileReceiver
}

type fileReceiver struct {
	consumer   *nsq.Consumer
	mReceivers map[string]*FileBuilder
	CPLch      chan *UserFile
}

type JsFileHandler struct {
	nsqaddr string
	label   string
	ch      chan *FileShard
}

func (h *JsFileHandler) HandleMessage(msg *nsq.Message) error {
	//	log.Println("new slice")
	//	log.Println(msg.Body)
	fShard := &FileShard{}
	err := json.Unmarshal(msg.Body, fShard)
	if err != nil {
		log.Println("JSON decode err: ", err)
		return nil
	}
	h.ch <- fShard
	return nil
}

func (this *fileReceiver) Receive() {
	cfg := nsq.NewConfig()
	cfg.ClientID = "proc engine web file"
	cfg.Hostname = "windows 7"
	cfg.HeartbeatInterval = time.Second * 30
	cfg.DialTimeout = time.Second * 30
	cfg.WriteTimeout = time.Second * 15
	//	cfg.OutputBufferSize = 128 * 1024
	//	cfg.OutputBufferTimeout = time.Second * 60
	handler1 := &JsFileHandler{nsqaddr: "120.76.41.114:4150", label: "handler1", ch: make(chan *FileShard, 1000)}

	cfg.ClientID = "proc engine web file"
	this.consumer, _ = nsq.NewConsumer("procFile", "channel1", cfg)
	this.consumer.AddHandler(handler1)
	this.consumer.ConnectToNSQD(handler1.nsqaddr)
	defer this.consumer.Stop()
	//初始化mReceivers
	this.mReceivers = make(map[string]*FileBuilder)

	rcv := &FileShard{}
	//	var fWrite *os.File
	//	var err error

	for {
		//1.在mReceivers中根据FileId找到FileReceiver
		//		fmt.Println("newslice!")
		rcv = <-handler1.ch

		err, fReceiver := this.GetReceiver(rcv)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		//2.写入buf, 修改Completeness
		_, ok := fReceiver.Completeness[rcv.ShardIndex]
		if !ok {
			n, err := fReceiver.Writer.WriteAt(rcv.Content, rcv.ShardIndex*rcv.ShardSize)
			if err != nil {
				log.Println(n, ", ", err.Error())
			}
			fReceiver.Completeness[rcv.ShardIndex] = 1
			log.Println("File:", rcv.FileName, "ShardIndex:", rcv.ShardIndex, "(total:", rcv.ShardCount, ")", " Progress:", len(fReceiver.Completeness)*100/int(rcv.ShardCount), "%")
		}

		//3.计算完成度
		//当前文件传输完成：
		if len(fReceiver.Completeness) == int(rcv.ShardCount) {
			fReceiver.Writer.Close()
			log.Println("FILE: ", rcv.FileName, "transfer completeness!")
			delete(this.mReceivers, rcv.FileName+fmt.Sprint(rcv.FileSize))
			userFile := &UserFile{UserId: rcv.UserId, FileName: rcv.FileName}
			this.CPLch <- userFile
		}

	}

}
func (this *fileReceiver) GetReceiver(rcv *FileShard) (err error, fReceiver *FileBuilder) {
	fReceiver, ok := this.mReceivers[rcv.FileName+fmt.Sprint(rcv.FileSize)]
	cur, _ := os.Getwd()
	path := cur[:strings.LastIndex(cur, "\\")] + "/res/"
	if !ok {
		//查看文件是否存在
		_, err1 := os.Stat(path + rcv.FileName)
		if err1 != nil {
			//文件或目录不存在
			fWrite, err := os.Create(path + rcv.FileName)
			if err != nil {
				log.Println(err.Error())
				return err, fReceiver
			}
			//新建一个文件接收对象
			fReceiver = &FileBuilder{
				Writer:       fWrite,
				FileId:       rcv.FileId,
				FileName:     rcv.FileName,
				FileSize:     rcv.FileSize,
				ShardCount:   rcv.ShardCount,
				ShardSize:    rcv.ShardSize,
				Completeness: make(map[int64]int),
			}
			this.mReceivers[path+rcv.FileName+fmt.Sprint(rcv.FileSize)] = fReceiver

			log.Println("GetReceiver: NEW ", fReceiver)
			return nil, fReceiver
		} else {
			//文件已存在
			err = errors.New("file exists already!")
			return err, fReceiver
		}

	}
	//	log.Println("GetReceiver: ", fReceiver)
	return nil, fReceiver
}
