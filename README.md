trajectory
==========

Distributed Process Profiler/Traceor


GO Environment
--------------

#### Install Go

    - http://golang.org/doc/install#install


#### Setup directory structure

Create a directory to house your go projects.

For example

```
mkdir -p $HOME/go
```

In your profile (bashrc, bash_profile, etc) setup your environment variables

```
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
```

Get Go and it's bin directorys in your path

```
export PATH=$PATH:$GOROOT/bin:$GOPATH:$GOPATH/bin
```


#### Add trajectory

```
mkdir -p $GOPATH/src/github.com/yourusername/
cd $GOPATH/src/github.com/yourusername
git clone git@github.com:lyddonb/trajectory.git
```


#### Install the trajectory deps

```
cd trajectory
go get ./...
```



Web Environment
---------------

#### Install Node

Install from http://nodejs.org/

```
make buildall
```


Create a login.json file with login info:

```
[
  {
    "Username": "uname",
    "Password": "pwd"
  }
]
```



#### Running Local (Redis in vagrant)


##### Start Redis in Vagrant

```
vagrant up
vagrant ssh
cd /vagrant
sudo make buildredis
sudo make runredis
exit
```

##### Local Build of JS/CSS

```
make buildjs
make buildcss
```

##### Run go server

```
go run main.go
```
