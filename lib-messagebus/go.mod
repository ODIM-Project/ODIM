module github.com/ODIM-Project/ODIM/lib-messagebus

go 1.19

require (
	github.com/BurntSushi/toml v1.0.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/satori/go.uuid v1.2.0
	github.com/segmentio/kafka-go v0.4.31
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/klauspost/compress v1.14.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
