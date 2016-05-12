// ProcessInstancePool
package model

import (
	"strconv"
)

var _ProcInstPool *ProcessInstancePool

type ProcessInstancePool struct {
	Pool map[int]*ProcessInstance
}

func initProcessPool() {
	_ProcInstPool = new(ProcessInstancePool)
	_ProcInstPool.Pool = make(map[int]*ProcessInstance)
}
func GetProcInstPool() *ProcessInstancePool {
	if _ProcInstPool == nil {
		initProcessPool()
	}
	return _ProcInstPool
}

func (this *ProcessInstancePool) AddProcInst(procInst *ProcessInstance) {
	id, err := strconv.Atoi(procInst.ProcDef.Id)
	if err != nil {
		panic(err)
	}
	_ProcInstPool.Pool[id] = procInst
}
func (this *ProcessInstancePool) DeleteProcInst(procInst *ProcessInstance) {
	id, err := strconv.Atoi(procInst.ProcDef.Id)
	if err != nil {
		panic(err)
	}
	delete(this.Pool, id)
}
