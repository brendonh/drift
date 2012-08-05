package storage

type StorageClient interface {
	Get(Storable) bool
	GetKey(bucket string, key string, target interface{}) bool

	Put(Storable) bool
	PutNew(bucket string, val interface{}) (string, bool)
	PutKey(bucket string, key string, val interface{}) bool

	Delete(bucket string, key string) bool
	
	Keys(bucket string) ([]string, bool)
}

type Storable interface {
	StorageKey() string
}
