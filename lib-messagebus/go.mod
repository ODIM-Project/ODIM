module github.com/ODIM-Project/ODIM/lib-messagebus

go 1.17

require (
	github.com/BurntSushi/toml v1.0.0
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20220426104855-9b203a83173f
	github.com/go-redis/redis/v8 v8.11.3
	github.com/satori/go.uuid v1.2.0
	github.com/segmentio/kafka-go v0.4.31
	github.com/sirupsen/logrus v1.8.1
)

require (
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.14.4 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
