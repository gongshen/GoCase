### 单机Redis（efficiency）
```shell
# 获取锁
SET key val NX PX 100
# 释放锁
if GET key.val == val then DEL key
```

### ETCD（correctness）

> 对于每个key来说，都需要对key配置Lease，租约到期，key就删掉了；
每个key都带有一个revision号，这个是单调递增的；
只有根据prefix获取最小的revision才能获得锁；
若获取失败，监听pre-revision的DELETE事件或者Lease过期，自己获得锁；

```go
grant, err := cli.Grant(context.TODO(), 2)
if err != nil {
	return err
}
sessoin, err := concurrency.NewSession(cli, concurrency.WithLease(grant.ID))
if err != nil {
	return err
}
mutex := concurrency.NewMutex(sessoin, "key")
if err := mutex.Lock(context.TODO()); err != nil {
		return err
}
// todo 业务处理
_ = mutex.Unlock(context.TODO())
```

[基于Redis的分布式锁到底安全吗](