// Logger
package logger

import (
	c "ProcessEngine/src/constant"
	//	"github.com/Sirupsen/logrus"
	//	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"log"
	//	"log/syslog"
	//	"os"
)

var _Logger *EngineLogger

//var _log *logrus.Logger

type EngineLogger struct {
	Name        string
	LogFileName string
	LogPath     string
	//	logFile     *os.File
	//	log         *log.Logger
}

func (lg *EngineLogger) WriteLog(str string, err error) {
	log.Println(str)
	if err != nil {
		log.Println(err)
	}
}
func initLogger() {
	_Logger = new(EngineLogger)
	_Logger.LogFileName = c.LOG_FILE_NAME
	_Logger.LogPath = c.LOG_FILE_PATH
	_Logger.Name = "DefaultEngineLogger"
	//newInstance()
}

func GetLogger() *EngineLogger {
	if _Logger == nil {
		initLogger()
	}
	return _Logger
}

//func newInstance() {
//	_log = logrus.New()
//	hook, err := logrus_syslog.NewSyslogHook("tcp", "10.128.48.239:6379", syslog.LOG_INFO, "")
//	_log.Info("new instance")
//	if err != nil {
//		_log.Errorf("Unable to connect to local syslog.")
//	}

//	_log.Hooks.Add(hook)

//	for _, level := range hook.Levels() {
//		if len(_log.Hooks[level]) != 1 {
//			_log.Errorf("SyslogHook was not added. The length of log.Hooks[%v]: %v", level, len(_log.Hooks[level]))
//		}
//	}
//}
