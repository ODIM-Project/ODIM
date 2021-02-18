module github.com/ODIM-Project/ODIM/svc-systems

go 1.13

require (
	github.com/Joker/hpp v1.0.0 // indirect
	github.com/ODIM-Project/ODIM/lib-dmtf v0.0.0-20201201072448-9772421f1b55
	github.com/ODIM-Project/ODIM/lib-persistence-manager v0.0.0-20201201072448-9772421f1b55
	github.com/ODIM-Project/ODIM/lib-rest-client v0.0.0-20201201072448-9772421f1b55
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20201201072448-9772421f1b55
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.5.1
	github.com/yudai/pp v2.0.1+incompatible // indirect
	gopkg.in/go-playground/validator.v9 v9.30.0
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
