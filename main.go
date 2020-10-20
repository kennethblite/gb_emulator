package main

import ("fmt"
	"encoding/json")


type Flag struct{
	Z string
	N string
	H string
	C string
}

type Instruction struct{
	Name string
	Group string
	TCycleBranch int
	TCycleBranch int
	Length int
	Flags Flag
	TimingNoBranch
}

func main(){

		var jsonBlob = []byte(`[
	{"Name": "Platypus", "Order": "Monotremata"},
	{"Name": "Quoll",    "Order": "Dasyuromorphia"}
]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)
	fmt.Println("hello world")
}
