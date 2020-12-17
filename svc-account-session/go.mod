module github.com/ODIM-Project/ODIM/svc-account-session

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20201201072448-9772421f1b55
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	gopkg.in/go-playground/validator.v9 v9.30.0
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
