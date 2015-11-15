all: install

checkgopath:
	@echo "Checking if project is at ${GOPATH}"
	@for miniotestspath in $(echo ${GOPATH} | sed 's/:/\n/g'); do if [ ! -d ${miniotestspath}/src/github.com/minio/minio-tests ]; then echo "Project not found in ${miniotestspath}, please follow instructions provided at https://github.com/minio/minio/blob/master/CONTRIBUTING.md#setup-your-minio-github-repository" && exit 1; fi done

verifiers: checkgopath

gomake-all: verifiers
	@GO15VENDOREXPERIMENT=1 go build -o $(GOPATH)/bin/minio-tests github.com/minio/minio-tests/cmd

install: gomake-all

clean:
	@rm -fv cover.out
	@rm -fv mc
