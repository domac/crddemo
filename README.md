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
