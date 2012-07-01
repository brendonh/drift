package drift

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/ugorji/go-msgpack"
)

type RiakClient struct {
	httpc   *http.Client
	baseURL string
}

type Storable interface {
	StorageKey() string
}

func NewClient(baseURL string) *RiakClient {
	c := new(RiakClient)
	c.httpc = &http.Client{}
	c.baseURL = baseURL
	return c
}

func (client *RiakClient) Get(obj Storable) bool {
	structName := reflect.TypeOf(obj).Elem().Name()
	return client.GetKey(structName, obj.StorageKey(), &obj)
}

func (client *RiakClient) GetKey(bucket string, key string, target interface{}) bool {
	url := client.buildURL(bucket, key)

	fmt.Printf("URL: %s\n", url)

	resp, ok := client.getRaw(url)

	if !ok {
		return false
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err := msgpack.Unmarshal(body, &target, nil)

	if err != nil {
		fmt.Printf("Decode err: %#v\n", err)
		return false
	}

	return true
}

func (client *RiakClient) Put(obj Storable) bool {
	structName := reflect.TypeOf(obj).Elem().Name()
	return client.PutKey(structName, obj.StorageKey(), obj)
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

	url := client.buildURL(bucket, key)

	fmt.Printf("Put URL: %s\n", url)
	req, err := http.NewRequest("POST", url, w)
	req.Header.Add("Content-Type", "text/plain")
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
