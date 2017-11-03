
- install go (refer https://golang.org/doc/install)
	+ download the archive file
	+ extract to "/usr/local" directory
	+ add GO_HOME=/usr/local/go/bin to PATH environment variable.

- Setup go local enviroment
	+ add GOPATH to enviroment available. Refer: https://github.com/golang/go/wiki/Setting-GOPATH
	  example: GOPATH=~/go
	+ add GOBIN=$GOPATH/bin to environment available.

- Setup project
	+ mkdir $GOPATH/bin
	+ mkdir $GOPATH/src & cd $_
	+ go get github.com/gin-gonic/gin
	+ go get github.com/kardianos/govendor
	+ mkdir -p hackathon/battleship && cd "$_"
	+ checkout our project
- how to run:
	+ go run main.go