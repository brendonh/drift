package storage

type StorageClient interface {
	Get(Storable) bool
	Put(Storable) bool
	IndexLookup(Storable, string) StorableIterator

	GetKey(bucket string, key string, target interface{}) bool
	PutNew(bucket string, val interface{}) (string, bool)
	PutKey(bucket string, key string, val interface{}) bool

	Delete(bucket string, key string) bool
	
	Keys(bucket string) ([]string, bool)
}

type Storable interface {
	StorageKey() string
	SetFromStorageKey(string)
}


type StorableIterator interface {
	Next(target Storable) bool
}
