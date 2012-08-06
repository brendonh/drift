package services

import (
	"container/list"
)

type APIData map[string]interface{}

type Method struct {
	Name string
	ArgSpec []Arg
	Handler func(APIData) (bool, APIData)
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
	handler func(APIData) (bool, APIData)) {

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

func (collection ServiceCollection) HandleRequest(request APIData) APIData {
	ok, resolutionErrors, args := Parse(requestArgSpec, request)
	if !ok {
		return ErrorResponse(ListToStringSlice(resolutionErrors))
	}

	return Response(collection.HandleCall(
		args["service"].(string), 
		args["method"].(string),
		args["data"].(APIData)))

}

func (collection ServiceCollection) HandleCall(
	serviceName string, 
	methodName string,
	data APIData,
    ) (bool, []string, APIData) {


	service, ok := collection.Services[serviceName]
	if !ok {
		return false, []string{"No such service"}, nil
	}

	method, ok := service.Methods[methodName]
	if !ok {
		return false, []string{"No such method"}, nil
	}

	ok, errors, args := Parse(method.ArgSpec, data)
	if !ok {
		return false, ListToStringSlice(errors), nil
	}

	ok, response := method.Handler(args)
	if !ok {
		return false, nil, response
	}

	return true, nil, response
}

func ListToStringSlice(l *list.List) []string {
	var slice = make([]string, l.Len())
	var i = 0
	for el := l.Front(); el != nil; el = el.Next() {
		slice[i] = el.Value.(string)
		i++
	}
	return slice
}


func Response(ok bool, errors []string, response APIData) APIData {
	if ok { 
		return SuccessResponse(response)
	}

	if errors != nil {
		return ErrorResponse(errors)
	}

	return FailureResponse(response)
}

func ErrorResponse(errors []string) APIData {
	var response = make(APIData)
	response["success"] = false
	response["reason"] = "call error"
	response["errors"] = errors
	return response
}

func SuccessResponse(data APIData) APIData {
	var response = make(APIData)
	response["success"] = true
	response["data"] = data
	return response
}

func FailureResponse(errors APIData) APIData {
	var response = make(APIData)
	response["success"] = false
	response["reason"] = "failure"
	response["errors"] = errors
	return response
}