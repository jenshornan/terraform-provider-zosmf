# TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=seb.se
NAMESPACE=seb
NAME=zosmf
BINARY=terraform-provider-${NAME}_${VERSION}
VERSION=0.3.0
OS_ARCH=linux_amd64
REPO7_PASSWORD=${REPO7_PASSWORD}

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_darwin_arm64
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

publish: build
	zip ${BINARY}_${OS_ARCH}.zip ${BINARY}
	curl -us44066:${REPO7_PASSWORD} -XPUT "https://repo7.sebank.se/artifactory/zos-terraform-provider/${NAMESPACE}/${NAME}/${VERSION}/${BINARY}_${OS_ARCH}.zip" -T ${BINARY}_${OS_ARCH}.zip

# test: 
#	go test -i $(TEST) || exit 1                                               
#	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

#testacc: 
#	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   