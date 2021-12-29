FROM golang:1.17.2
RUN apt update && apt install

RUN mkdir /go/src/app

WORKDIR /go/src/app

RUN go get \
    #github.com/uudashr/gopkgs/v2/cmd/gopkgs \
    github.com/ramya-rao-a/go-outline
#github.com/nsf/gocode \
#github.com/acroca/go-symbols \
#github.com/fatih/gomodifytags \
#github.com/josharian/impl \
#github.com/haya14busa/goplay/cmd/goplay \
#golang.org/x/lint/golint \

RUN GOBIN=/tmp/ go install github.com/go-delve/delve/cmd/dlv@master
RUN go install golang.org/x/tools/gopls@latest

RUN mv /tmp/dlv $GOPATH/bin/dlv-dap
