# crddemo : Kubernetes 自定义控制器Demo例子

---- 

## 基本目录结构

```bash
├── code-gen.sh
├── controller.go
├── crd
│   └── mydemo.yaml
├── example
│   └── example-mydemo.yaml
├── go.mod
├── go.sum
├── main.go
└── pkg
    └── apis
        └── crddemo
            ├── register.go
            └── v1
                ├── doc.go
                ├── register.go
                ├── types.go
```

## 执行生成Kubernetes客户端代码

```
sh code-gen.sh
```

执行后，会在`crddemo/pkg`目录生成对应的代码

```bash
client
├── clientset
│   └── versioned
│       ├── clientset.go
│       ├── doc.go
│       ├── fake
│       │   ├── clientset_generated.go
│       │   ├── doc.go
│       │   └── register.go
│       ├── scheme
│       │   ├── doc.go
│       │   └── register.go
│       └── typed
│           └── crddemo
│               └── v1
│                   ├── crddemo_client.go
│                   ├── doc.go
│                   ├── fake
│                   │   ├── doc.go
│                   │   ├── fake_crddemo_client.go
│                   │   └── fake_mydemo.go
│                   ├── generated_expansion.go
│                   └── mydemo.go
├── informers
│   └── externalversions
│       ├── crddemo
│       │   ├── interface.go
│       │   └── v1
│       │       ├── interface.go
│       │       └── mydemo.go
│       ├── factory.go
│       ├── generic.go
│       └── internalinterfaces
│           └── factory_interfaces.go
└── listers
    └── crddemo
        └── v1
            ├── expansion_generated.go
            └── mydemo.go
```

## 编译项目

```bash
$ make 

... ...

gofmt -w .
go test -v . 
?       github.com/domac/crddemo        [no test files]
mkdir -p releases
GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags "-s -w" -v -o releases/crddemo *.go
github.com/golang/groupcache/lru
k8s.io/apimachinery/third_party/forked/golang/json
k8s.io/apimachinery/pkg/util/mergepatch
k8s.io/kube-openapi/pkg/util/proto
k8s.io/client-go/tools/record/util
k8s.io/apimachinery/pkg/util/strategicpatch
k8s.io/client-go/tools/record
command-line-arguments
go clean -i
... ...
```