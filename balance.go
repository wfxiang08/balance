package main

import (
	"log"
	"os"

	BA "github.com/darkhelmet/balance/backends"
	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
)

var cmd = &commander.Command{
	Short: "load balance tcp, http, and https connections to multiple backends",
}

func ensureBind(bindFlag *flag.Flag) string {
	if bindFlag == nil {
		log.Fatalln("bind flag not defined")
	}

	bind, ok := bindFlag.Value.Get().(string)
	if !ok {
		log.Fatalln("bind flag must be defined as a string")
	}

	if bind == "" {
		log.Fatalln("specify the address to listen on with -bind")
	}

	return bind
}

func buildBackends(balanceFlag *flag.Flag, backends []string) BA.Backends {
	if balanceFlag == nil {
		log.Fatalln("balance flag not defined")
	}

	balance, ok := balanceFlag.Value.Get().(string)
	if !ok {
		log.Fatalln("balance flag must be defined as a string")
	}

	if balance == "" {
		log.Fatalln("specify the balancing algorithm with -balance")
	}

	return BA.Build(balance, backends)
}

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	// 参数如何处理呢?
	fs.String("bind", "", "the address to listen on")
	fs.String("balance", "round-robin", "the balancing algorithm to use")
	return fs
}

//
// balancer如何实现呢?
//
// func tcpBalance(bind string, backends BA.Backends) error
func balancer(f func(string, BA.Backends) error) func(*commander.Command, []string) error {

	return func(cmd *commander.Command, args []string) error {
		// 取出 bind, backends, 并且将它和f绑定起来
		bind := ensureBind(cmd.Flag.Lookup("bind"))

		// backends的
		backends := buildBackends(cmd.Flag.Lookup("balance"), args)
		return f(bind, backends)
	}
}

func main() {
	// 注意cmd在上面被实例化
	// balance tcp -bind :4000 localhost:4001 localhost:4002
	// balance http -bind :4000 localhost:4001 localhost:4002
	//
	// 调用子命令，并且将参数传递给子命令
	err := cmd.Dispatch(os.Args[1:])
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
