package services

import (
	. "drift/common"
	"fmt"
	"container/list"
)

const (
	IntArg = iota
	FloatArg
	StringArg
	NestedArg
    RawArg
)

func Parse(argspec []APIArg, args APIData) (
	bool, *list.List, APIData) {

	var parsedArgs APIData = make(APIData);
	var errors = list.New()

	for _, arg := range argspec {
		
		givenVal, ok := args[arg.Name]
		if !ok {
			if arg.Default != nil {
				parsedArgs[arg.Name] = arg.Default
			} else {
				errors.PushBack(fmt.Sprintf(
					"Missing argument: %s (%s)", 
					arg.Name, stringArgType(arg.ArgType)))
			}
			continue
		}

		ok, conversionErrors, val := convertArgVal(arg, givenVal)

		if !ok {
			if conversionErrors != nil {
				for e := conversionErrors.Front(); e != nil; e = e.Next() {
					errors.PushBack(fmt.Sprintf(
						"In %s: %s", arg.Name, e.Value))
				}
			} else {
				errors.PushBack(fmt.Sprintf(
					"Invalid value for %s (expected %s): %v", 
					arg.Name, stringArgType(arg.ArgType), givenVal))
			}
			continue
		}

		parsedArgs[arg.Name] = val
	}

	if (errors.Len() > 0) {
		return false, errors, nil
	}

	return true, nil, parsedArgs
}


func convertArgVal(arg APIArg, val interface{}) (
	bool, *list.List, interface{}) {
	defer func() {
			recover()
	}()
	switch arg.ArgType {
	case IntArg:
		var floatval float64 = val.(float64)
		return true, nil, int(floatval);
	case FloatArg:
		return true, nil, val.(float64)
	case StringArg:
		return true, nil, val.(string)
	case NestedArg:
		spec := arg.Extra.([]APIArg)
		nest := val.(APIData)
		return Parse(spec, nest)
	case RawArg:
		return true, nil, val
	}
	return false, nil, nil
}


func stringArgType(argType int) string {
	switch argType {
	case IntArg: return "int"
	case FloatArg: return "float"
	case StringArg: return "string"
	case NestedArg: return "nested"
	case RawArg: return "raw"
	}
	return "unknown"
}