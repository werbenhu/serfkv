[English](README.md) | [简体中文](README-CN.md)

# SerfKV

使用Go写的分布式内存键值存储服务，使用[hashicorp/serf](https://github.com/hashicorp/serf) 做服务发现，提供HTTP接口设置和获取键值对。

## 安装

```shell
go github.com/werbenhu/serfkv
```

## 用法

```shell
serfkv
  -cluster-addr string
        memberlist主机:端口 (默认为":9601")
  -http-addr string
        HTTP主机:端口 (默认为":4001")
  -members string
        集群成员列表，例如127.0.0.1:9601,127.0.0.1:9602
```

### 创建集群

启动第一个节点":9601"
```shell
serfkv --cluster-addr=127.0.0.1:9601
```

Make a note of the local node address
```
================================
start serf cluster on 127.0.0.1:9601
================================
...
Listening and serving HTTP on :4001
```

启动第二个节点，并将第一个节点作为节点列表的一部分
```shell
serfkv --members=127.0.0.1:9601 --http-addr=:4002 --cluster-addr=127.0.0.1:9602
```

应该看到输出：
```
================================
start serf cluster on 127.0.0.1:9602
================================
...
Listening and serving HTTP on :4002
```

First node output will log the new connection
```shell
A node has joined: DESKTOP-SJCIADA-cgegdtvfb8nju677cl40
```

## HTTP接口

- /get - 获取一个键值对
- /set - 设置一个键值对
- /del - 删除一个键值对

http请求查询参数为key和val

```shell
# add
curl "http://localhost:4001/set?key=foo&val=bar"

# get
curl "http://localhost:4001/get?key=foo"

# delete
curl "http://localhost:4001/del?key=foo"
```
