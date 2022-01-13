module github.com/ODIM-Project/ODIM/svc-telemetry

go 1.17

require (
	github.com/ODIM-Project/ODIM/lib-dmtf v0.0.0-20201201072448-9772421f1b55
	github.com/ODIM-Project/ODIM/lib-persistence-manager v0.0.0-20210901061202-f84c396a018e // indirect
	github.com/ODIM-Project/ODIM/lib-rest-client v0.0.0-20210201172557-4fa2adafe1e3
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20210519055855-227d83cff80f
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.5.1
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/CloudyKit/jet/v3 v3.0.0 // indirect
	github.com/Shopify/goreferrer v0.0.0-20181106222321-ec9c9a553398 // indirect
	github.com/aymerick/raymond v2.0.3-0.20180322193309-b565731e1464+incompatible // indirect
	github.com/coreos/etcd v3.3.17+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eknkc/amber v0.0.0-20171010120322-cdade1c07385 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.2-0.20200519141726-cb32006e483f // indirect
	github.com/googleapis/gnostic v0.1.0 // indirect
	github.com/iris-contrib/blackfriday v2.0.0+incompatible // indirect
	github.com/iris-contrib/jade v1.1.3 // indirect
	github.com/iris-contrib/pongo2 v0.0.1 // indirect
	github.com/iris-contrib/schema v0.0.1 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kataras/golog v0.0.10 // indirect
	github.com/kataras/iris/v12 v12.1.8 // indirect
	github.com/kataras/pio v0.0.2 // indirect
	github.com/kataras/sitemap v0.0.5 // indirect
	github.com/klauspost/compress v1.13.4 // indirect
	github.com/microcosm-cc/bluemonday v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/ryanuber/columnize v2.1.0+incompatible // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/schollz/closestmatch v2.1.0+incompatible // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987 // indirect
	google.golang.org/grpc v1.41.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.51.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20191120175047-4206685974f2 // indirect
	k8s.io/api v0.18.5 // indirect
	k8s.io/apimachinery v0.18.5 // indirect
	k8s.io/client-go v0.18.5 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89 // indirect
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
	google.golang.org/grpc v1.41.0 => google.golang.org/grpc v1.26.0
)
