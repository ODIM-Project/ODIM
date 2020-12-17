module github.com/ODIM-Project/ODIM/svc-task

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-messagebus v0.0.0-20201201072448-9772421f1b55
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20201201072448-9772421f1b55
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2-0.20200519141726-cb32006e483f
	github.com/micro/go-micro v1.13.2
	github.com/satori/go.uuid v1.2.0
	github.com/satori/uuid v1.2.0
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
