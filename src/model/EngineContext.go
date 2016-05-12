// EngineConfiguration
package model

import (
	c "ProcessEngine/src/constant"
	"ProcessEngine/src/util"
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
)

var _Context *Context

type Context struct {
	Config *EngineConfiguration
}
type CXEngineConfiguration struct {
	XMLName    xml.Name           `xml:"Configuration"`
	Properties []CXEngineProperty `xml:"property"`
}
type CXEngineProperty struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

func InitContext() {
	_Context = new(Context)
	config, err := ReadConfigFile()
	if err != nil {
		panic(err)
	}
	_Context.Config = config
}
func GetContext() *Context {
	if _Context == nil {
		InitContext()
	}
	return _Context
}
func ReadConfigFile() (*EngineConfiguration, error) {
	// convert process config file
	config := new(EngineConfiguration)
	cur, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	config.EnginePath = cur
	file, err := os.Open(config.EnginePath + util.GetSp() + c.CONFIG_CONF_PATH + util.GetSp() + c.CONFIG_CFG_FILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	engineConfigStr, err := ioutil.ReadAll(file)
	e := CXEngineConfiguration{}
	err = xml.Unmarshal(engineConfigStr, &e)
	if err != nil {
		return nil, err
	}
	for _, prop := range e.Properties {
		switch {
		case strings.EqualFold(prop.Name, "dbType"):
			config.DBType = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "username"):
			config.Username = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "password"):
			config.Password = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "dbName"):
			config.DBName = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "msgAddress"):
			config.MsgServerAddress = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "dbAddress"):
			config.DBAddress = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "fileNats"):
			config.FileNats = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "fileRcvTopic"):
			config.FileRcvTopic = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "TdbType"):
			config.TaskDbType = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "TdbAddr"):
			config.TaskDbAddr = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "TdbUsername"):
			config.TaskDbUsername = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "TdbPassword"):
			config.TaskDbPassword = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "TdbName"):
			config.TaskDbName = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "EdbType"):
			config.EventDbType = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "EdbAddr"):
			config.EventDbAddr = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "EdbUsername"):
			config.EventDbUsername = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "EdbPassword"):
			config.EventDbPassword = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "EdbName"):
			config.EventDbName = strings.TrimSpace(prop.Value)
		case strings.EqualFold(prop.Name, "logMsgAddress"):
			config.LogNatsAddr = strings.TrimSpace(prop.Value)
		}
	}
	return config, nil
}
