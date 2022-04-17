FROM alpine:3.11.6

RUN apk update &&\
  apk upgrade &&\
  apk add ca-certificates &&\
  apk add make &&\
  apk add git &&\
  apk add go

ENV GOPATH=/root/go
RUN mkdir -p /root/go/src/github.com/mzk622/go-links/backend/
ADD . /root/go/src/github.com/mzk622/go-links/backend/
WORKDIR /root/go/src/github.com/mzk622/go-links/backend/
RUN make build
RUN make test

FROM alpine:3.11.6

RUN apk update &&\
  apk upgrade &&\
  apk add ca-certificates

RUN addgroup mzk &&\
  adduser mzk -h /home/mzk -s /bin/bash -G mzk -D &&\
  chmod 700 /home/mzk &&\
  passwd mzk -d

WORKDIR /home/mzk
RUN mkdir bin

COPY --from=0 /root/go/src/github.com/mzk622/go-links/backend bin/go-links-api

RUN chown -R mzk:mzk /home/mzk
USER mzk

EXPOSE 8080

CMD bin/go-links-api
