module github.com/ODIM-Project/ODIM/svc-api

go 1.13

require (
	github.com/Joker/jade v1.0.0 // indirect
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20210506103851-66c53837fd0f
	github.com/flosch/pongo2 v0.0.0-20200913210552-0d938eb266f3 // indirect
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/iris-contrib/formBinder v5.0.0+incompatible // indirect
	github.com/kataras/iris v11.1.1+incompatible
	github.com/kataras/iris/v12 v12.1.9-0.20200616210209-a85c83b70ad0
	github.com/sirupsen/logrus v1.4.2
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
