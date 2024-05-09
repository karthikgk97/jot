echo "Performing go mod tidy"
go mod tidy

echo "Building Jot"
go build

echo "Copying jot to /usr/bin"
sudo cp jot /usr/local/bin/jot
