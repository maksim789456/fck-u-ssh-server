# Fck U Ssh Server
Small ssh server to print message to who try to connecting to 22 port

## Build
```sh
go mod download
go build
```

## Run
Generate dummy keys for ssh server
```sh
ssh-keygen -t rsa -b 4096 -C "dummy@ssh.com" -f id_rsa -q -N ""
```

Run
```sh
./dummy_ssh_server
```
