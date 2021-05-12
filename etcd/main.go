package main

import 	"github.com/coreos/etcd/clientv3"

type ETCDConfig struct {
	client clientv3.Client
}