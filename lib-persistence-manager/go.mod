module github.com/ODIM-Project/ODIM/lib-persistence-manager

go 1.17

require (
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gomodule/redigo v2.0.0+incompatible
)

require (
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/CloudyKit/jet/v6 v6.1.0 // indirect
	github.com/Shopify/goreferrer v0.0.0-20210630161223-536fa16abd6f // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/aymerick/raymond v2.0.2+incompatible // indirect
	github.com/eknkc/amber v0.0.0-20171010120322-cdade1c07385 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/flosch/pongo2/v4 v4.0.2 // indirect
	github.com/goccy/go-json v0.9.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/iris-contrib/jade v1.1.4 // indirect
	github.com/iris-contrib/schema v0.0.6 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kataras/blocks v0.0.5 // indirect
	github.com/kataras/golog v0.1.7 // indirect
	github.com/kataras/iris/v12 v12.2.0-alpha9 // indirect
	github.com/kataras/pio v0.0.10 // indirect
	github.com/kataras/sitemap v0.0.5 // indirect
	github.com/kataras/tunnel v0.0.3 // indirect
	github.com/klauspost/compress v1.14.4 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/microcosm-cc/bluemonday v1.0.18 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/schollz/closestmatch v2.1.0+incompatible // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/tdewolff/minify/v2 v2.10.0 // indirect
	github.com/tdewolff/parse/v2 v2.5.27 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/yosssi/ace v0.0.5 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220224211638-0e9765cccd65 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
	github.com/gomodule/redigo => github.com/gomodule/redigo v1.8.9

)
