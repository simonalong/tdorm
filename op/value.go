package op

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	goleTime "github.com/simonalong/gole/time"
	"github.com/simonalong/gole/util"
)

func ValueToStr(data interface{}) string {
	if data == nil {
		return ""
	}

	dataTem := data
	dataType := reflect.TypeOf(data)
	if dataType.Kind() == reflect.Ptr {
		dataType = dataType.Elem()
		dataTem = reflect.ValueOf(data).Elem().Interface()
	}
	if util.IsStringType(dataType) {
		if strings.HasPrefix(data.(string), "'") && strings.HasSuffix(data.(string), "'") {
			return fmt.Sprintf("%v", data)
		} else if strings.Contains(data.(string), "now") {
			return fmt.Sprintf("%v", data)
		} else {
			return fmt.Sprintf("'%v'", data)
		}
	} else if util.IsNumberType(dataType) || util.IsBoolType(dataType) {
		return fmt.Sprintf("%v", data)
	} else if util.IsTimeType(dataType) {
		if dataType == reflect.TypeOf(time.Time{}) {
			return fmt.Sprintf("'%v'", goleTime.TimeToStringYmdHmsS(dataTem.(time.Time)))
		} else if dataType == reflect.TypeOf(time.Duration(0)) {
			return fmt.Sprintf("%v", dataTem.(time.Duration).String())
		}
	}
	return fmt.Sprintf("%v", dataTem)
}
