package main

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
)

func main() {
	id := uuid.Must(uuid.NewV4())

	fmt.Printf("uuid info: %s\n", id)

	fmt.Printf("uuid info: %s\n", uuid.Must(uuid.NewV4()).String())

	tmp := map[string]uuid.UUID{
		"5a3a7138-44a3-11ec-8a23-93b5625893d3": uuid.FromStringOrNil("5a66aad2-44a3-11ec-87fd-53ddb650bd33"),
	}
	fmt.Printf("tmp: %v\n", tmp)

	m, err := json.Marshal(tmp)
	if err != nil {
		fmt.Printf("Could not marshal the tmp. err: %v", err)
	}
	fmt.Printf("marshal: %v\n", m)

	res := map[string]uuid.UUID{}
	if err := json.Unmarshal(m, &res); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	for k, v := range res {
		fmt.Printf("res. key: %s, value: %s\n", k, v.String())
	}

}
