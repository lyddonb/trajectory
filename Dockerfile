FROM ubuntu

#RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list

RUN apt-get -y update
RUN apt-get -y upgrade

RUN apt-get install -y --force-yes \
    software-properties-common \
    python-software-properties \
    python \
    build-essential \
    wget \
    git-core 

RUN add-apt-repository ppa:chris-lea/node.js
RUN apt-get -y update
RUN apt-get install -y --force-yes nodejs

ENV GO_VERSION 1.3
ENV OS linux
ENV ARCH amd64

ENV HOME /root
ENV GOROOT /usr/local/go
ENV GOPATH $HOME/go
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin:$GOPATH

ENV TRAJ github.com/lyddonb/trajectory

RUN mkdir -p $GOPATH

RUN wget http://golang.org/dl/go$GO_VERSION.$OS-$ARCH.tar.gz --no-check-certificate

RUN tar -C /usr/local -xzf go$GO_VERSION.$OS-$ARCH.tar.gz

WORKDIR /root/go

RUN go get $TRAJ

WORKDIR /root/go/src/github.com/lyddonb/trajectory

RUN make buildall

CMD ["make", "dockerrun"]

#EXPOSE 4180
EXPOSE 1301
EXPOSE 3001
