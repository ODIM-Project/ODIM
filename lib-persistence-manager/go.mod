module github.com/ODIM-Project/ODIM/lib-persistence-manager

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200925145026-eac0549d2f51
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gomodule/redigo v2.0.0+incompatible
)
replace github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200925145026-eac0549d2f51 => ../lib-utilities
