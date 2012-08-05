package main

import (
	//"drift/storage"
	"drift/accounts"
	"drift/services"
	"fmt"

	"bytes"
	"encoding/json"

)

type Sector struct {
	X, Y int
	Name string
}

func (sector *Sector) StorageKey() string {
	return fmt.Sprintf("%d:%d", sector.X, sector.Y)
}

func main() {
	serviceCollection := services.NewServiceCollection()
	serviceCollection.AddService(accounts.GetService())

	blob := bytes.NewBufferString(
		`{"service": "accounts", "method": "register", "data": {"email": "brendonh3@gmail.com", "password": "test"}}`).Bytes()	

	var args map[string]interface{}
	err := json.Unmarshal(blob, &args)
	
	if err != nil {
		fmt.Printf("Oh no: %s\n", err)
		return
	}

	var response = serviceCollection.Handle(args)

	reply, _ := json.Marshal(response)

	fmt.Printf("Response: %s\n", reply)
}


// func main() {
// 	client := drift.NewRiakClient("http://localhost:8098")

// 	sector := Sector{0, 1, "Away"}

// 	ok := client.Put(&sector)

// 	if !ok {
// 		fmt.Printf("Write Failed\n")
// 		return
// 	}

// 	fmt.Printf("Ok\n")

// 	newSector := Sector{X: 0, Y: 1}
// 	ok = client.Get(&newSector)

// 	if !ok {
// 		fmt.Printf("Read Failed\n")
// 		return
// 	}

// 	fmt.Printf("%s\n", newSector.Name)
// }
