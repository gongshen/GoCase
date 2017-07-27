package main

import (
	"github.com/go-yaml/yaml"
	"log"
	"fmt"
)

const data=`
a: Easy!
b:
 c: 2
 d: [3,4]
`

type T struct {
	A string
	B struct{
		C int	`yaml:"c"`
		D []int	`yaml:",flow"`
	}
}

func main()  {
	t:=T{}
	//解码
	err:=yaml.Unmarshal([]byte(data),&t)
	if err!=nil{
		log.Fatalf("error:%v",err)
	}
	fmt.Printf("--- d:\n%v\n\n",t)
	c,err:=yaml.Marshal(&t)
	if err!=nil{
		log.Fatalf("error:%v",err)
	}
	fmt.Printf("--- c:\n%v\n",string(c))

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- d2:\n%v\n\n", m)

	c2, err := yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- c2:\n%s\n", string(c2))
}

