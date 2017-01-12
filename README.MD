# Chainecode

On this demo we create 4 peer with PBFT (Practical Byzantine Fault Tolerance) consensus. And we create 4 users and each want to move the 1 value on the key to the next key (and set 0 yo the last)

So the user is ab,bc,cd,da
 - ab: It waits for the value of "a" to be 1 to write the value 1 in "b" and write 0 in "a"
 - bc: It waits for the value of "b" to be 1 to write the value 1 in "c" and write 0 in "b"
 - cd: It waits for the value of "c" to be 1 to write the value 1 in "d" and write 0 in "c"
 - da: It waits for the value of "d" to be 1 to write the value 1 in "a" and write 0 in "d"
 
|       | a | b | c | d |
|-------|---|---|---|---|
| t0    | 1 | 0 | 0 | 0 |
| t1    | 0 | 1 | 0 | 0 |
| t2    | 0 | 0 | 1 | 0 |
| t3    | 0 | 0 | 0 | 1 |
| t4=t0 | 1 | 0 | 0 | 0 |

## Install

 We need docker and docker-compose 
 
   1. start peer
   
   It will download the docker "hyperledger/fabric-peer"
   ```
   # on deamon mode
   docker-compose up -d 
   ```
   
   2. start user
   
   It will download the docker "dngroup/user"
   ```
   ./users.sh
   ```
   
## edit

  We need docker and golang (v1.6)
   
   - maybe you need go dependency `github.com/hyperledger/fabric/`
        - `cd $GOPATH/src/github.com/hyperledger/`
        - `git clone --branch v0.6.1-preview --single-branch --depth 1 https://github.com/hyperledger/fabric.git`
   - edit [chainecode.go](/chainecode/chainecode.go) for change the chainecode
        - build for verify if working    `go build ./`
   - edit [user.go](/user/user.go) for change the chainecode
        - use the script to build the new docker image
        
   
## more information 

 to build https://github.com/hyperledger/fabric we need that
  
  ```
  cd $REPONAME
  cd /scripts/devenv && chmod +x setup.sh && sudo ./setup.sh
  cd /tmp
  git clone --branch v4.1 --single-branch --depth 1 https://github.com/facebook/rocksdb.git
  cd rocksdb
  make shared_lib
  sudo INSTALL_PATH=/usr/local make install-shared
  sudo ldconfig
  sudo mkdir /var/hyperledger/ && sudo chown $USER:$USER /var/hyperledger
  ```
  and finally you can make 
  ```
  cd $REPONAME
  make
  ```