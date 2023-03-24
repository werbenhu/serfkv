package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/werbenhu/serfkv/cluster"
)

var (
	nodes       = flag.String("members", "", "seeds memberlist of cluster,such as 192.168.0.1:9601,192.168.0.2:9601")
	clusterAddr = flag.String("cluster-addr", ":9601", "memberlist host:port")
	httpAddr    = flag.String("http-addr", ":4001", "http host:port")
	srv         *cluster.Server
)

func init() {
	flag.Parse()
}

func del(c *gin.Context) {
	key := c.Query("key")

	if err := srv.Delete(key, true); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "success")
}

func get(c *gin.Context) {
	key := c.Query("key")

	val, err := srv.Get(key)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "success, val:%s", val.(string))
}

func set(c *gin.Context) {
	key := c.Query("key")
	val := c.Query("val")

	if err := srv.Set(key, val, true); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "success")
}

func main() {
	var err error
	var members []string
	if len(*nodes) > 0 {
		members = strings.Split(*nodes, ",")
	}

	if srv, err = cluster.New(&cluster.Options{
		Address: *clusterAddr,
		Members: members,
	}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/get", get)
	r.GET("/set", set)
	r.GET("/del", del)
	r.Run(*httpAddr)
}
