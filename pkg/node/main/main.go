package main

import (
	"flag"
	"os"
	"sync"
	"teamide/pkg/base"
	"teamide/pkg/node"
)

var (
	waitGroupForStop sync.WaitGroup
)

func main() {
	for _, v := range os.Args {
		if v == "-version" || v == "-v" {
			println(base.GetVersion())
			return
		}
	}

	var id string
	var address string
	var token string
	var connAddress string
	var connToken string
	flag.StringVar(&id, "id", "", "节点ID，不可变更，需要唯一")
	flag.StringVar(&address, "address", "", "节点启动监听地址")
	flag.StringVar(&token, "token", "", "节点Token，用于验证")
	flag.StringVar(&connAddress, "connAddress", "", "上层节点连接地址")
	flag.StringVar(&connToken, "connToken", "", "上层节点连接Token")

	//解析
	flag.Parse()

	if id == "" {
		flag.Usage()
		panic("请设置 -id")
	}
	//if address == "" {
	//	flag.Usage()
	//	panic("请设置 -address")
	//}
	if token == "" {
		flag.Usage()
		panic("请设置 -token")
	}
	if connAddress != "" && connToken == "" {
		flag.Usage()
		panic("请设置 -connToken")
	}

	server := &node.Server{}
	server.Start()
	localNode := &node.LocalNode{
		Id:          id,
		BindAddress: address,
		BindToken:   token,
		ConnAddress: connAddress,
		ConnToken:   connToken,
	}
	println("启动节点 [" + id + "][" + address + "] 开始")
	server.AddLocalNode(localNode)
	println("启动节点 [" + id + "][" + address + "] 成功")

	waitGroupForStop.Add(1)

	waitGroupForStop.Wait()
}
