package argparser

import (
	"fmt"
	"container/list"
)

const (
	IntArg = iota
	FloatArg
	StringArg
)

type Arg struct {
	Name string
	ArgType int
	Required bool
	Default interface{}
}


func Parse(argspec []Arg, args map[string]interface{}) (bool, *list.List, map[string]interface{}) {
	var parsedArgs map[string]interface{} = make(map[string]interface{});
	var errors = list.New()

	for _, arg := range argspec {
		
		givenVal, ok := args[arg.Name]
		if !ok {
			errors.PushBack(fmt.Sprintf("Missing argument: %s (%s)", arg.Name, stringArgType(arg.ArgType)))
			continue
		}

		val := convertArgVal(arg.ArgType, givenVal)

		if val == nil {
			errors.PushBack(fmt.Sprintf("Invalid value for %s (expected %s): %v", arg.Name, stringArgType(arg.ArgType), givenVal))
			continue
		}

		parsedArgs[arg.Name] = val
	}

	if (errors.Len() > 0) {
		return false, errors, nil
	}

	return true, nil, parsedArgs
}


func convertArgVal(argType int, val interface{}) interface{} {
	defer func() {
			recover()
	}()
	switch argType {
	case IntArg:
		var floatval float64 = val.(float64)
		return int(floatval);
	case FloatArg:
		return val.(float64)
	case StringArg:
		return val.(string)
	}
	return nil
}


func stringArgType(argType int) string {
	switch argType {
	case IntArg: return "int"
	case FloatArg: return "float"
	case StringArg: return "string"
	}
	return "unknown"
}