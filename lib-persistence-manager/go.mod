module github.com/ODIM-Project/ODIM/lib-persistence-manager

go 1.17

require (
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20201201072448-9772421f1b55
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/sirupsen/logrus v1.4.2
)

require (
	github.com/google/uuid v1.1.2-0.20200519141726-cb32006e483f // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.16.0 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
