FROM golang:1.11.5
MAINTAINER thejchap
ARG PKG=/go/src/github.com/thejchap/waffle
ENV PKG ${PKG}
ADD . ${PKG}
WORKDIR ${PKG}
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go get ./cmd/waffle
RUN dep ensure
RUN go build -o main ./cmd/waffle
EXPOSE 3000
CMD ["/go/src/github.com/thejchap/waffle/main"]
