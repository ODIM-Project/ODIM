module github.com/ODIM-Project/ODIM/lib-persistence-manager

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200727091052-cb7db65624ce
	github.com/go-redis/redis/v8 v8.0.0-beta.10
	github.com/gomodule/redigo v2.0.0+incompatible
)
replace github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200727091052-cb7db65624ce => ../lib-utilities
