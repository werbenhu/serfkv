# SerfKV

A distributed in-memory key-value store built using [hashicorp/serf](https://github.com/hashicorp/serf) with HTTP API

## Install

```shell
go github.com/werbenhu/serfkv
```

## Usage

```shell
serfkv
  -cluster-addr string
        memberlist host:port (default ":9601")
  -http-addr string
        http host:port (default ":4001")
  -members string
        seeds memberlist of cluster,such as 127.0.0.1:9601,127.0.0.1:9602
```

### Create Cluster

Start first node ":9601"
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

Start second node with first node as part of the nodes list
```shell
serfkv --members=127.0.0.1:9601 --http-addr=:4002 --cluster-addr=127.0.0.1:9602
```

You should see the output
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

## HTTP API

- /get - get a value
- /set - set a value
- /del - delete a value

Query params expected are `key` and `val`

```shell
# add
curl "http://localhost:4001/set?key=foo&val=bar"

# get
curl "http://localhost:4001/get?key=foo"

# delete
curl "http://localhost:4001/del?key=foo"
```
