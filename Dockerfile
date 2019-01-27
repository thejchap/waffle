FROM golang:1.11.5
MAINTAINER thejchap
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ARG PKG=/go/src/github.com/thejchap/waffle
RUN mkdir -p ${PKG}
WORKDIR ${PKG}
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only
COPY . ${PKG}
RUN go build -o /go/bin/waffle ./cmd/waffle
EXPOSE 3000
CMD ["/go/bin/waffle"]
