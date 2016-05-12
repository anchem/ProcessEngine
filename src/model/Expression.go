// Expression
package model

import (
	logger "ProcessEngine/src/logger"
	"regexp"
	"strconv"
	"strings"
)

//判断表达式是否合
func CheckConditionExpression(cdnExpr string) bool {
	if strings.EqualFold(cdnExpr, "") {
		return true
	} else {
		strSplit := cdnExpr[strings.Index(cdnExpr, "${")+2 : strings.LastIndex(cdnExpr, "}")]
		reg := "^[A-Za-z]\\w*\\s*[!>=<]+\\s*\\w+"
		m, _ := regexp.Match(reg, []byte(strings.TrimSpace(strSplit)))
		return m
	}
}
func CalculateConditionExpression(expr string, varMap map[string]*DataObject) bool {
	if expr == "" {
		return false
	}
	splReg := "[!>=<]+"
	strSplit := expr[strings.Index(expr, "${")+2 : strings.LastIndex(expr, "}")]
	r, _ := regexp.Compile(splReg)
	splStr := string(r.Find([]byte(strSplit)))
	vars := strings.Split(strSplit, splStr)
	var1 := strings.TrimSpace(vars[0])
	var2 := strings.TrimSpace(vars[1])
	// find process instance vars
	rstFlag := false
	varInProc := new(DataObject)
LOOKUP:
	for key, temp := range varMap {
		if strings.EqualFold(key, var1) {
			varInProc.Name = temp.Name
			varInProc.ValueType = temp.ValueType
			varInProc.ValueString = temp.ValueString
			varInProc.Value = temp.Value
			rstFlag = true
			break LOOKUP
		}
	}
	if rstFlag == false {
		logger.GetLogger().WriteLog("<Expression> [ERROR] : can not find variable in process instance:"+var1, nil)
		return false
	}
	// calculate
	switch {
	case strings.EqualFold(varInProc.ValueType, "int"):
		oprt, err := strconv.Atoi(varInProc.ValueString)
		if err != nil {
			logger.GetLogger().WriteLog("<Expression> [ERROR] : varInProc wrong type", err)
			return false
		}
		result, err := strconv.Atoi(var2)
		if err != nil {
			logger.GetLogger().WriteLog("<Expression> [ERROR] : expr: var2 wrong type", err)
			return false
		}
		switch {
		case strings.EqualFold("!=", splStr):
			return oprt != result
		case strings.EqualFold("==", splStr):
			return oprt == result
		case strings.EqualFold("<", splStr):
			return oprt < result
		case strings.EqualFold(">", splStr):
			return oprt > result
		case strings.EqualFold("<=", splStr):
			return oprt <= result
		case strings.EqualFold(">=", splStr):
			return oprt >= result
		default:
			logger.GetLogger().WriteLog("<Expression> [ERROR] : unknown error calculating int", nil)
			return false
		}
	case strings.EqualFold(varInProc.ValueType, "long"):
		oprt, err := strconv.ParseInt(varInProc.ValueString, 10, 64)
		if err != nil {
			logger.GetLogger().WriteLog("<Expression> [ERROR] : varInProc wrong type", err)
			return false
		}
		result, err := strconv.ParseInt(var2, 10, 64)
		if err != nil {
			logger.GetLogger().WriteLog("<Expression> [ERROR] : var2 wrong type", err)
			return false
		}
		switch {
		case strings.EqualFold("!=", splStr):
			return oprt != result
		case strings.EqualFold("==", splStr):
			return oprt == result
		case strings.EqualFold("<", splStr):
			return oprt < result
		case strings.EqualFold(">", splStr):
			return oprt > result
		case strings.EqualFold("<=", splStr):
			return oprt <= result
		case strings.EqualFold(">=", splStr):
			return oprt >= result
		default:
			logger.GetLogger().WriteLog("<Expression> [ERROR] : unknown error calculating long", nil)
			return false
		}
	case strings.EqualFold(varInProc.ValueType, "float"):
		return false
	case strings.EqualFold(varInProc.ValueType, "bool"):
		oprt, err := strconv.ParseBool(varInProc.ValueString)
		if err != nil {
			logger.GetLogger().WriteLog("<Expression> [ERROR] : varInProc wrong type", err)
			return false
		}
		result, err := strconv.ParseBool(var2)
		if err != nil {
			logger.GetLogger().WriteLog("<Expression> [ERROR] : var2 wrong type", err)
			return false
		}
		switch {
		case strings.EqualFold("!=", splStr):
			return oprt != result
		case strings.EqualFold("==", splStr):
			return oprt == result
		default:
			logger.GetLogger().WriteLog("<Expression> [ERROR] : unknown error calculating bool", nil)
			return false
		}
	default:
		switch {
		case strings.EqualFold("!=", splStr):
			return !strings.EqualFold(varInProc.ValueString, var2)
		case strings.EqualFold("==", splStr):
			return strings.EqualFold(varInProc.ValueString, var2)
		default:
			logger.GetLogger().WriteLog("<Expression> [ERROR] : unknown error calculating string", nil)
			return false
		}
	}
	return false
}
