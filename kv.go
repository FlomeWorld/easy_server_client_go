package kv_go

import (
	"encoding/asn1"
	"sync"
	"net/rpc"
	"net"
	"log"
	"fmt"
)

const (
	OK = "OK"
	ErrNoKey = "ErrNoKey"
)

type Err string

type PutArgs struct {
	Key string
	Value string
}

type PutReply struct {
	Err Err
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Err Err
	Value string
}

type KV struct {
	mu sync.Mutex
	keyValue map[string] string
}


func (kv *KV) Get(args *GetArgs,reply *GetReply) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	reply.Err = "OK"
	val,ok := kv.keyValue[args.Key]

	if ok {
		reply.Err = OK
		reply.Value = val

	}else{
		reply.Err = ErrNoKey
		reply.Value = ""
	}
	return nil

}

func (kv *KV) Put(args *PutArgs,reply *PutReply) error {
	kv.mu.Lock()
	defer  kv.mu.Unlock()


	kv.keyValue[args.Key] = args.Value
	reply.Err = OK
	return nil

}

func server(){
	kv := new(KV)

	kv.keyValue = map[string]string{}
	rpcs := rpc.NewServer()
	rpcs.Register(kv)
	l,e := net.Listen("tcp",":1234")
	if e!= nil{
		log.Fatal("listen error:",e)

	}

	go func() {
		for{
			conn,err :=l.Accept()
			if err == nil {
				go rpcs.ServeConn(conn)
			}else{
				break
			}
		}
		l.Close()
		fmt.Printf("Server done\n")
	}()


}


//Client

func Dial() *rpc.Client {
	client,err := rpc.Dial("tcp",":1234")

	if err != nil{
		log.Fatal("dialing:",err)
	}
	return client
}

func Get(key string)string{
	client := Dial()
	args := &GetArgs{"subject"}
	reply := GetReply{"",""}
	err := client.Call("KV.Get",args,&reply)
	if err != nil{
		log.Fatal("error:",err)
	}
	client.Close()
	return reply.Value

}

func Put(key string,val string){
	client := Dial()
	args := &PutArgs{"subject","6.824"}
	reply := PutReply{""}
	err := client.Call("KV.Put",args,&reply)
	if err != nil {
		log.Fatal("error:", err)
	}

	client.Close()

}


//MAIN


func main()  {
	server()

	Put("subject","6.824")
	fmt.Printf("Put subject done\n")
	fmt.Printf("Get %s\n",Get("subject"))
}