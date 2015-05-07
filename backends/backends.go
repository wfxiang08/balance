package backends

import (
	"log"
)

type backend struct {
	hostname string
}

func (b *backend) String() string {
	return b.hostname
}

type Backend interface {
	String() string
}

type Backends interface {
	Choose() Backend
	Len() int
	Add(string)
	Remove(string)
}

type Factory func([]string) Backends

var factories = make(map[string]Factory)

func Build(algorithm string, specs []string) Backends {
	// 从factories中取出对应的算法
	factory, found := factories[algorithm]
	if !found {
		log.Fatalf("balance algorithm %s not supported", algorithm)
	}
	// 创建Factory
	return factory(specs)
}
