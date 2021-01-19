go mod download
go mod vendor
go build .
if [ $? -eq 0 ]; then
    echo Build for plugin ur service $i is Successful !!!!
else
    echo Build for plugin ur service $i is Failed !!!!
fi
