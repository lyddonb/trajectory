FROM ubuntu

#RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list

RUN apt-get -y update
RUN apt-get -y upgrade

RUN apt-get install -y --force-yes build-essential
RUN apt-get install -y --force-yes wget
RUN apt-get install -y --force-yes g++
RUN apt-get install -y --force-yes make

# REDIS Install
RUN wget http://download.redis.io/releases/redis-2.8.8.tar.gz --no-check-certificate
RUN tar -C /usr/local -xzf redis-2.8.8.tar.gz
WORKDIR /usr/local/redis-2.8.8
RUN make
ENTRYPOINT ["/usr/local/redis-2.8.8/src/redis-server"]
EXPOSE 6379

CMD ["--bind", "0.0.0.0"]
