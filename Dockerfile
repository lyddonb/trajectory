FROM ubuntu

RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list

RUN apt-get -y update
RUN apt-get -y upgrade

RUN apt-get install -y wget
RUN apt-get install -y git
RUN apt-get install -y make


# GO Install
RUN wget https://go.googlecode.com/files/go1.2.1.linux-amd64.tar.gz --no-check-certificate

RUN tar -C /usr/local -xzf go1.2.1.linux-amd64.tar.gz

ENV GOROOT /usr/local/go

RUN mkdir -p /opt/go/trajectory/{bin, src, pkg}
RUN mkdir -p /opt/go/trajectory/src/github.com/lyddonb

ENV GOPATH /opt/go
ENV GOBIN /opt/go/bin
ENV PATH $PATH:/opt/go/bin:/usr/local/go/bin

ADD . /opt/go/src/github.com/lyddonb/trajectory
WORKDIR /opt/go/src/github.com/lyddonb/trajectory

RUN go get github.com/garyburd/redigo/redis
RUN go install trajectory.go

EXPOSE 1200

ENTRYPOINT ["/opt/trajectory/bin/trajectory"]
