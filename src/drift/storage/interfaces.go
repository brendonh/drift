package storage

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
