package model

type Result struct {
	Res    byte
	Reason string
}

type JsonDef struct {
	UserId    string
	FileName  string
	ProcessId string
	ProcDef   string
	TrsReason string
}
type ProcessTrs struct {
	Id           int
	ProcId       string
	ProcName     string
	ProcDesc     string
	ProcFile     string
	FileType     byte
	CreateTime   int64
	NumInstance  int64
	IsMultiple   bool
	IsExecutable bool
	Running      bool
}
