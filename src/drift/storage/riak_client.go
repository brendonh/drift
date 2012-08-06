package storage

import (
	"bytes"
	"encoding/json"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"text/template"

	"github.com/ugorji/go-msgpack"
)

type RiakClient struct {
	httpc   *http.Client
	baseURL string
}

func NewRiakClient(baseURL string) StorageClient {
	return &RiakClient{ &http.Client{}, baseURL }
}

func (client *RiakClient) Get(obj Storable) bool {
	structName := reflect.TypeOf(obj).Elem().Name()
	return client.GetKey(structName, obj.StorageKey(), &obj)
}

func (client *RiakClient) GetKey(bucket string, key string, target interface{}) bool {
	url := client.buildURL(bucket, key)

	resp, ok := client.getRaw(url)

	if !ok {
		return false
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return client.Decode(body, target)
}


func (client *RiakClient) DecodeString(content string, target interface{}) bool {
	var bytes = bytes.NewBufferString(content).Bytes()
	return client.Decode(bytes, target)
}

func (client *RiakClient) Decode(content []byte, target interface{}) bool {
	var packed = make([]byte, base64.StdEncoding.DecodedLen(len(content)));

	_, err := base64.StdEncoding.Decode(packed, content)

	var temp interface{}
	msgpack.Unmarshal(packed, &temp, nil)

	err = msgpack.Unmarshal(packed, &target, nil)

	if err != nil {
		fmt.Printf("Decode err: %#v\n", err)
		return false
	}

	return true
}


func (client *RiakClient) Put(obj Storable) bool {
	structName := reflect.TypeOf(obj).Elem().Name()

	var key = obj.StorageKey()
	var ok bool

	if key == "" {
		key, ok = client.PutNew(structName, obj)
		if ok {
			obj.SetFromStorageKey(key)
			// Temp hack, hopefully
			client.Put(obj)
		}
	} else {
		_, ok = client.putRaw(structName, key, obj)
	}

	return ok
}

func (client *RiakClient) PutNew(bucket string, val interface{}) (string, bool) {
	resp, ok := client.putRaw(bucket, "", val)
	if !ok {
		return "", false
	}

	location := resp.Header["Location"][0]
	lastSlash := strings.LastIndex(location, "/")
	return location[lastSlash+1:], true
}

func (client *RiakClient) PutKey(bucket string, key string, val interface{}) bool {
	_, ok := client.putRaw(bucket, key, val)
	return ok
}

func (client *RiakClient) Keys(bucket string) ([]string, bool) {
	url := strings.Join([]string{client.baseURL, "buckets", bucket, "keys"}, "/") +
		"?keys=true"

	resp, ok := client.getRaw(url)

	if !ok {
		return nil, false
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result keyList
	err := json.Unmarshal(body, &result)

	if err != nil {
		fmt.Printf("Illegal response: %s\n", body)
		fmt.Printf("%#v\n", err)
		return nil, false
	}

	return result.Keys, true
}

func (client *RiakClient) Delete(bucket string, key string) bool {
	url := client.buildURL(bucket, key)
	req, err := http.NewRequest("DELETE", url, nil)
	resp, err := client.httpc.Do(req)

	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return false
	}

	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		fmt.Printf("Status: %d\n", resp.StatusCode)
		return false
	}

	return true
}


var indexQueryTemplate = `{
    "inputs": {
        "bucket": "{{.bucket}}",
        "index": "{{.index}}",
        "key": "{{.value}}"
    },
    "query": [
        {
            "map": {
                "language": "erlang",
                "module": "riak_kv_mapreduce",
                "function": "map_object_value"
            }
        }
    ]
}`


func (client *RiakClient) IndexLookup(obj Storable, index string) StorableIterator {
	var structType = reflect.TypeOf(obj)
	var structName = structType.Elem().Name()
	var indexValue = reflect.ValueOf(obj).Elem().FieldByName(index).String()
	fmt.Printf("Looking up %s => %s\n", index, indexValue)

	var tmpl = template.Must(template.New("indexQuery").Parse(indexQueryTemplate))

	var vars = make(map[string]string)
	vars["bucket"] = structName
	vars["index"] = index + "_bin"
	vars["value"] = indexValue

	var query = bytes.NewBufferString("")
	tmpl.Execute(query, vars)

	var url = client.baseURL + "/mapred"

	req, err := http.NewRequest("POST", url, query)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.httpc.Do(req)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return nil
	} 

	if resp.StatusCode != 200 {
		fmt.Printf("Status: %d\n", resp.StatusCode)
		fmt.Printf("%s\n", body)
		return nil
	}

	var blobs []string
	err = json.Unmarshal(body, &blobs)

	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return nil
	} 

	return NewBlobIterator(client, blobs)
}

type BlobIterator struct {
	client *RiakClient
	objs []string
	index int
}

func NewBlobIterator(client *RiakClient, blobs []string) *BlobIterator {
	return &BlobIterator{client, blobs, 0}
}

func (it *BlobIterator) Next(target Storable) bool {
	if it.index >= len(it.objs) {
		return false
	}

	it.client.DecodeString(it.objs[it.index], target)
	it.index += 1
	return true
}



type keyList struct {
	Keys []string `json:"keys"`
}

func (client *RiakClient) putRaw(bucket string, key string, val interface{}) (*http.Response, bool) {

	w := bytes.NewBufferString("")
	enc := msgpack.NewEncoder(w)
	err := enc.Encode(val)

	if err != nil {
		fmt.Printf("Err: %#v\n", err)
		return nil, false
	}

	b64 := base64.StdEncoding.EncodeToString(w.Bytes())

	url := client.buildURL(bucket, key)

	fmt.Printf("Put URL: %s\n", url)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(b64))
	req.Header.Add("Content-Type", "application/octet-stream")

	var structInfo = reflect.TypeOf(val).Elem()
	var structValue = reflect.ValueOf(val).Elem()
	for i := 0; i < structInfo.NumField(); i++ {
		var field = structInfo.Field(i)
		if field.Tag.Get("indexed") != "" {
			var indexKey = "x-riak-index-" + field.Name + "_bin"
			var indexVal = structValue.FieldByIndex([]int{i}).String()
			req.Header.Add(indexKey, indexVal)
		}
	}

	resp, err := client.httpc.Do(req)

	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return nil, false
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		fmt.Printf("Status: %d\n", resp.StatusCode)

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Body: %s\n", body)
		return nil, false
	}

	return resp, true

}

func (client *RiakClient) getRaw(url string) (*http.Response, bool) {
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.httpc.Do(req)

	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return nil, false
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Status: %d\n", resp.StatusCode)
		return nil, false
	}

	return resp, true
}

func (client *RiakClient) buildURL(bucket string, key string) string {
	return strings.Join([]string{client.baseURL, "riak", bucket, key}, "/")
}
