// Converter
package converter

import (
	c "ProcessEngine/src/constant"
	"errors"
	"strings"
)

type Converter struct {
}

func (cvt *Converter) ConvertToBpmnModel(filename string, filepath string, filetype byte) (interface{}, error) {
	// identify which converter is needed
	if strings.EqualFold(filename, "") || strings.EqualFold(filepath, "") {
		return nil, errors.New("empty params")
	}
	switch filetype {
	case c.CONVERTER_FILE_TYPE_XML:
		return ConvertXmlToBpmnModel(filename, filepath)
	case c.CONVERTER_FILE_TYPE_JSON:
		return ConvertJsonToBpmnModel(filename, filepath)
	default:
		return nil, errors.New("unknown file type")
	}
	return nil, errors.New("unknown error")
}
