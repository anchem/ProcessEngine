// EngineConfiguration
package model

type EngineConfiguration struct {
	EnginePath       string
	DBAddress        string
	DBType           string
	Username         string
	Password         string
	DBName           string
	MsgServerAddress string

	FileNats     string
	FileRcvTopic string

	TaskDbType     string
	TaskDbAddr     string
	TaskDbUsername string
	TaskDbPassword string
	TaskDbName     string

	EventDbType     string
	EventDbAddr     string
	EventDbUsername string
	EventDbPassword string
	EventDbName     string

	LogNatsAddr string
}
