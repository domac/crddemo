# crddemo : Kubernetes 自定义控制器Demo例子


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

## 执行程序文件

把crddemo放到kubernetes集群中，或者本地也行，只要能访问到apiserver和具备kubeconfig

可以看到，程序运行的时候，一开始会报错。这是因为，此时 Mydemo 对象的 CRD 还没有被创建出来，所以 Informer 去 APIServer 里获取 Mydemos 对象时，并不能找到 Mydemo 这个 API 资源类型的定义

```
$  ./crddemo --kubeconfig=/data/svr/projects/kubernetes/config/kubectl.kubeconfig --master=http://127.0.0.1:8080  -alsologtostderr=true 
I0308 12:23:18.494507   27426 controller.go:79] Setting up mydemo event handlers
I0308 12:23:18.494829   27426 controller.go:105] Starting Mydemo control loop
I0308 12:23:18.494840   27426 controller.go:108] Waiting for informer caches to sync
E0308 12:23:18.496902   27426 reflector.go:178] github.com/domac/crddemo/pkg/client/informers/externalversions/factory.go:117: Failed to list *v1.Mydemo: the server could not find the requested resource (get mydemos.crddemo.k8s.io)
E0308 12:23:18.497477   27426 reflector.go:178] github.com/domac/crddemo/pkg/client/informers/externalversions/factory.go:117: Failed to list *v1.Mydemo: the server could not find the requested resource (get mydemos.crddemo.k8s.io)
E0308 12:23:21.604508   27426 reflector.go:178] github.com/domac/crddemo/pkg/client/informers/externalversions/factory.go:117: Failed to list *v1.Mydemo: the server could not find the requested resource (get mydemos.crddemo.k8s.io)
E0308 12:23:26.932293   27426 reflector.go:178] github.com/domac/crddemo/pkg/client/informers/externalversions/factory.go:117: Failed to list *v1.Mydemo: the server could not find the requested resource (get mydemos.crddemo.k8s.io)

... ...
```

接下来，我们执行我们自定义资源的定义文件：

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mydemos.crddemo.k8s.io
spec:
  group: crddemo.k8s.io
  version: v1
  names:
    kind: Mydemo
    plural: mydemos
  scope: Namespaced
```

执行创建

```
$  kubectl apply -f crd/mydemo.yaml 
customresourcedefinition.apiextensions.k8s.io/mydemos.crddemo.k8s.io created
```

此时，观察crddemo的日志输出，可以看到Controller的日志恢复了正常，控制循环启动成功

```
I0308 12:30:29.956263   28282 controller.go:113] Starting workers
I0308 12:30:29.956307   28282 controller.go:118] Started workers
```

然后，我们可以对我们的Mydemo对象进行增删改查操作了。

首先，新建一个自定义资源对象

example-mydemo.yaml 

```yaml
apiVersion: crddemo.k8s.io/v1
kind: Mydemo
metadata:
  name: example-mydemo
spec:
  ip: "127.0.0.1"
  port: 8080
```

执行创建

```
$  kubectl apply -f example-mydemo.yaml 
mydemo.crddemo.k8s.io/example-mydemo created
```

创建成功够，看k8s集群是否成功存储起来

```
$  kubectl get Mydemo                   
NAME             AGE
example-mydemo   2s
```

这时候，查看一下控制器的输出：

```
I0308 12:31:24.983663   28282 controller.go:216] [DemoCRD] Try to process mydemo: &v1.Mydemo{TypeMeta:v1.TypeMeta{Kind:"", APIVersion:""}, ObjectMeta:v1.ObjectMeta{Name:"example-mydemo", GenerateName:"", Namespace:"default", SelfLink:"/apis/crddemo.k8s.io/v1/namespaces/default/mydemos/example-mydemo", UID:"8a6d17f7-17f3-4a1d-8250-bb092678ae7e", ResourceVersion:"10818363", Generation:1, CreationTimestamp:v1.Time{Time:time.Time{wall:0x0, ext:63719238684, loc:(*time.Location)(0x1e566c0)}}, DeletionTimestamp:(*v1.Time)(nil), DeletionGracePeriodSeconds:(*int64)(nil), Labels:map[string]string(nil), Annotations:map[string]string{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"crddemo.k8s.io/v1\",\"kind\":\"Mydemo\",\"metadata\":{\"annotations\":{},\"name\":\"example-mydemo\",\"namespace\":\"default\"},\"spec\":{\"ip\":\"127.0.0.1\",\"port\":8080}}\n"}, OwnerReferences:[]v1.OwnerReference(nil), Finalizers:[]string(nil), ClusterName:"", ManagedFields:[]v1.ManagedFieldsEntry(nil)}, Spec:v1.MydemoSpec{Ip:"127.0.0.1", Port:8080}} ...
I0308 12:31:24.983844   28282 controller.go:174] Successfully synced 'default/example-mydemo'
I0308 12:31:24.983893   28282 event.go:278] Event(v1.ObjectReference{Kind:"Mydemo", Namespace:"default", Name:"example-mydemo", UID:"8a6d17f7-17f3-4a1d-8250-bb092678ae7e", APIVersion:"crddemo.k8s.io/v1", ResourceVersion:"10818363", FieldPath:""}): type: 'Normal' reason: 'Synced' Mydemo synced successfully
```

可以看到，我们上面创建 example-mydemo.yaml 的操作，触发了 EventHandler 的`添加`事件，从而被放进了工作队列。紧接着，控制循环就从队列里拿到了这个对象，并且打印出了正在`处理`这个 Mydemo 对象的日志。


我们这时候，尝试修改资源，对对应的port属性进行修改

```yaml
apiVersion: crddemo.k8s.io/v1
kind: Mydemo
metadata:
  name: example-mydemo
spec:
  ip: "127.0.0.1"
  port: 9090
```

手段执行修改

```
$  kubectl apply -f example-mydemo.yaml 
```

此时，crddemo新增出来的日志如下：

```
I0308 12:32:05.663044   28282 controller.go:216] [DemoCRD] Try to process mydemo: &v1.Mydemo{TypeMeta:v1.TypeMeta{Kind:"", APIVersion:""}, ObjectMeta:v1.ObjectMeta{Name:"example-mydemo", GenerateName:"", Namespace:"default", SelfLink:"/apis/crddemo.k8s.io/v1/namespaces/default/mydemos/example-mydemo", UID:"8a6d17f7-17f3-4a1d-8250-bb092678ae7e", ResourceVersion:"10818457", Generation:2, CreationTimestamp:v1.Time{Time:time.Time{wall:0x0, ext:63719238684, loc:(*time.Location)(0x1e566c0)}}, DeletionTimestamp:(*v1.Time)(nil), DeletionGracePeriodSeconds:(*int64)(nil), Labels:map[string]string(nil), Annotations:map[string]string{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"crddemo.k8s.io/v1\",\"kind\":\"Mydemo\",\"metadata\":{\"annotations\":{},\"name\":\"example-mydemo\",\"namespace\":\"default\"},\"spec\":{\"ip\":\"127.0.0.1\",\"port\":9080}}\n"}, OwnerReferences:[]v1.OwnerReference(nil), Finalizers:[]string(nil), ClusterName:"", ManagedFields:[]v1.ManagedFieldsEntry(nil)}, Spec:v1.MydemoSpec{Ip:"127.0.0.1", Port:9080}} ...
I0308 12:32:05.663179   28282 controller.go:174] Successfully synced 'default/example-mydemo'
I0308 12:32:05.663208   28282 event.go:278] Event(v1.ObjectReference{Kind:"Mydemo", Namespace:"default", Name:"example-mydemo", UID:"8a6d17f7-17f3-4a1d-8250-bb092678ae7e", APIVersion:"crddemo.k8s.io/v1", ResourceVersion:"10818457", FieldPath:""}): type: 'Normal' reason: 'Synced' Mydemo synced successfully
```

可以看到，这一次，Informer 注册的“更新”事件被触发，更新后的 Mydemo 对象的 Key 被添加到了工作队列之中。

所以，接下来控制循环从工作队列里拿到的 Mydemo 对象，与前一个对象是不同的：它的 ResourceVersion 的值从 10818363 变成了 10818457 ；而 Spec 里的Port字段，则变成了 9080。最后，我再把这个对象删除掉：

```
$  kubectl delete -f example-mydemo.yaml 
mydemo.crddemo.k8s.io "example-mydemo" deleted
```


这一次，在控制器的输出里，我们就可以看到，Informer 注册的“删除”事件被触发,输出如下：

```
W0308 12:33:08.494755   28282 controller.go:203] DemoCRD: default/example-mydemo does not exist in local cache, will delete it from Mydemo ...
I0308 12:33:08.495793   28282 controller.go:206] [DemoCRD] Deleting mydemo: default/example-mydemo ...
I0308 12:33:08.495808   28282 controller.go:174] Successfully synced 'default/example-mydemo'

```
然后，k8s集群的资源也被清除了：

```
$  kubectl get Mydemo                    
No resources found in default namespace.
```

以上就是使用自定义控制器的基本开发流程
