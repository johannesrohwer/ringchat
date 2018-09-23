# Ringchat v0.1
Ringchat is a rather minimal example for the usage of gRPC in Golang.
Multiple distributed instances of the Ringchat application
can form a ring structure and exchange broadcast text messages
with each other.

## Usage
To build and run Ringchat as a master node clone the repository
to your `$GOPATH` and execute the following in a terminal:
```
$ make
$ ./ringchat --master true
```

If you wish to connect to this master, run:
```
$ make
$ ./ringchat --master-host 192.168.2.100 --master-port 9999
```

## Dependencies
- gRPC
- github.com/google/uuid

## Known Bugs and Limitations
- No reconnect features; a broken ring cannot be fixed
- Ring structure is rigid; arbitrary insertion of new nodes is not possible
- Docker support currently missing