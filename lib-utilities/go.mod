module github.com/ODIM-Project/ODIM/lib-utilities

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-persistence-manager v0.0.0-20201201072448-9772421f1b55
	github.com/coreos/etcd v3.3.17+incompatible
	github.com/fsnotify/fsnotify v1.4.7
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.2-0.20200519141726-cb32006e483f
	github.com/kataras/iris/v12 v12.1.8
	github.com/micro/go-micro v1.13.2
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/net v0.0.0-20210510120150-4163338589ed
	google.golang.org/grpc v1.24.0
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.5
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
