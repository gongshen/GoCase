package main

import "fmt"

type Read struct {
	key string
	reply chan<- string
}

type Write struct {
	key string
	value string
}

var hole=make(chan interface{})

func deepspace()  {
	m:=map[string]string{}
	for{
		switch t:=(<-hole).(type) {
		case Read:
			t.reply <- m[t.key]+"from Mars."
		case Write:
			m[t.key]=t.value
		}
	}
}

func main()  {
	go deepspace()
	hole<-Write{"Name","Miracle "}
	home :=make(chan string)
	hole<-Read{"Name",home}
	fmt.Println(<-home)
}