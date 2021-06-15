module github.com/ODIM-Project/ODIM/svc-telemetry

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-dmtf v0.0.0-20201201072448-9772421f1b55
	github.com/ODIM-Project/ODIM/lib-rest-client v0.0.0-20210201172557-4fa2adafe1e3
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20210519055855-227d83cff80f
	github.com/ODIM-Project/ODIM/svc-update v0.0.0-20210201172557-4fa2adafe1e3
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.5.1
)

replace (
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
