FROM golang
MAINTAINER Evey Quirk

ADD . /go/src/github.com/evq/mdns-liaison

RUN cd /go/src/github.com/evq/mdns-liaison && go get
RUN cd /go/src/github.com/evq/mdns-liaison && go install

CMD /go/bin/mdns-liaison

EXPOSE 53
