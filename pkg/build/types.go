package build

import (
	"github.com/golang/groupcache/singleflight"
	"sync"
)

type KanikoArgs struct {
	BuildArg         string
	Cache            bool
	Dockerfile       string
	Insecure         bool
	InsecurePull     bool
	InsecureRegistry bool
	SkipTlsVerify    bool
	Target           string
}

type Builder struct {
	Tag          string
	BuildContext string
	PodName      string
	Namespace    string
	Dockerfile   string
	KanikoArgs
}

type SyncStore[T any] struct {
	sf      singleflight.Group
	results syncMap[T]
}

type syncMap[T any] struct {
	sync.Map
}

type StoreError struct {
	message string
}
