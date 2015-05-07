package main

import (
	"fmt"
	"io"
	"log"
	"net"

	BA "github.com/darkhelmet/balance/backends"
	"github.com/gonuts/commander"
)

func copy(wc io.WriteCloser, r io.Reader) {
	defer wc.Close()
	io.Copy(wc, r)
}

func handleConnection(us net.Conn, backend BA.Backend) {
	if backend == nil {
		log.Printf("no backend available for connection from %s", us.RemoteAddr())
		us.Close()
		return
	}

	// 连接到后端
	ds, err := net.Dial("tcp", backend.String())
	if err != nil {
		log.Printf("failed to dial %s: %s", backend, err)
		us.Close()
		return
	}

	// 数据拷贝
	// Ignore errors
	// 这个地方需要进一步细化，每次解析出一个命令之后，再选择connection(dest)
	go copy(ds, us)
	go copy(us, ds)
}

func tcpBalance(bind string, backends BA.Backends) error {
	log.Println("using tcp balancing")

	// balance tcp -bind :4000 localhost:4001 localhost:4002
	// 1. 首先绑定到某个ip/port
	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return fmt.Errorf("failed to bind: %s", err)
	}

	log.Printf("listening on %s, balancing %d backends", bind, backends.Len())

	// 2. 来请求了，则交给handleConnection
	//    backends随机选一个后端的服务
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept: %s", err)
			continue
		}
		go handleConnection(conn, backends.Choose())
	}

	return err
}

// 自动调用
func init() {
	fs := newFlagSet("tcp")

	cmd.Subcommands = append(cmd.Subcommands, &commander.Command{
		UsageLine: "tcp [options] <backend> [<more backends>]",
		Short:     "performs tcp based load balancing",
		Flag:      *fs,
		Run:       balancer(tcpBalance),
	})
}
