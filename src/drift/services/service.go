package services

import (
	"container/list"
)

type Method struct {
	Name string
	ArgSpec []Arg
	Handler func(map[string]interface{}) (bool, map[string]interface{})
}

type Service struct {
	Name string
	Methods map[string]Method
}

type ServiceCollection struct {
	Services map[string]Service
}

func NewService(name string) *Service {
	return &Service{ 
		Name: name,
	    Methods: make(map[string]Method),
	}
}

func NewServiceCollection() *ServiceCollection {
	return &ServiceCollection{ 
		Services: make(map[string]Service),
	}
}

func (service *Service) AddMethod(
	name string, 
	argSpec []Arg,
	handler func(map[string]interface{}) (bool, map[string]interface{})) {

	service.Methods[name] = Method{
		Name: name,
		ArgSpec: argSpec,
		Handler: handler, 
	}
}

func (collection *ServiceCollection) AddService(service *Service) {
	collection.Services[service.Name] = *service
}

var requestArgSpec = []Arg {
	Arg{Name: "service", ArgType: StringArg},
	Arg{Name: "method", ArgType: StringArg},
	Arg{Name: "data", ArgType: RawArg},
}

func (collection ServiceCollection) Handle(request map[string]interface{}) map[string]interface{} {
	ok, errors, args := Parse(requestArgSpec, request)
	if !ok {
		return ErrorResponse(ListToSlice(errors))
	}

	service, ok := collection.Services[args["service"].(string)]
	if !ok {
		return ErrorResponse([]interface{}{"No such service"})
	}

	method, ok := service.Methods[args["method"].(string)]
	if !ok {
		return ErrorResponse([]interface{}{"No such method"})
	}

	ok, errors, args = Parse(method.ArgSpec, args["data"].(map[string]interface{}))
	if !ok {
		return ErrorResponse(ListToSlice(errors))
	}

	ok, response := method.Handler(args)
	if !ok {
		return FailureResponse(response)
	}

	return SuccessResponse(response)
}

func ListToSlice(l *list.List) []interface{} {
	var slice = make([]interface{}, l.Len())
	var i = 0
	for el := l.Front(); el != nil; el = el.Next() {
		slice[i] = el.Value
		i++
	}
	return slice
}


func ErrorResponse(errors []interface{}) map[string]interface{} {
	var response = make(map[string]interface{})
	response["success"] = false
	response["reason"] = "call error"
	response["errors"] = errors
	return response
}

func SuccessResponse(data map[string]interface{}) map[string]interface{} {
	var response = make(map[string]interface{})
	response["success"] = true
	response["data"] = data
	return response
}

func FailureResponse(errors map[string]interface{}) map[string]interface{} {
	var response = make(map[string]interface{})
	response["success"] = false
	response["reason"] = "failure"
	response["errors"] = errors
	return response
}