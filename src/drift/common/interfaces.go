package common


// ------------------------------------------
// Server
// ------------------------------------------

type ServerContext interface {
	Storage() StorageClient
	API() API
}

type Endpoint interface {
	Start() bool
	Stop() bool
}



// ------------------------------------------
// API
// ------------------------------------------

const (
	IntArg = iota
	FloatArg
	StringArg
	NestedArg
    RawArg
)

type APIArg struct {
	Name string
	ArgType int
	Required bool
	Default interface{}
	Extra interface{}
}

type APIMethod struct {
	Name string
	ArgSpec []APIArg
	Handler APIHandler
}

type APIData map[string]interface{}

type APIHandler func(APIData, ServerContext) (bool, APIData)



// ------------------------------------------
// Services
// ------------------------------------------

type APIService interface {
	Name() string
	AddMethod(string, []APIArg, APIHandler)
	FindMethod(string) *APIMethod
}

type API interface {
	AddService(APIService)
	HandleRequest(APIData, ServerContext) APIData
	HandleCall(string, string, APIData, ServerContext) (bool, []string, APIData)
}



// ------------------------------------------
// Storage
// ------------------------------------------

type StorageClient interface {
	GenerateID() string

	Get(Storable) bool
	Put(Storable) bool
	IndexLookup(obj Storable, results interface{}, index string) bool

	GetKey(bucket string, key string, target interface{}) bool
	PutNew(bucket string, val interface{}) (string, bool)
	PutKey(bucket string, key string, val interface{}) bool

	Delete(bucket string, key string) bool
	
	Keys(bucket string) ([]string, bool)
}


type Storable interface {
	StorageKey() string
}
