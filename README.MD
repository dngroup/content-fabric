# Contant Fabric 


## Install

 We need docker and docker-compose 
 
   1. start peer
   
   It will download the docker "hyperledger/fabric-peer"
   ```
   # on deamon mode
   cd $REPONAME
   docker-compose up -d 
   ```
   
   2. start user
   
  Working in progress
   
## edit

  We need docker and golang (v1.6)
   
   - maybe you need go dependency `github.com/hyperledger/fabric/`
        - `cd $GOPATH/src/github.com/hyperledger/`
        - `git clone --branch v0.6.1-preview --single-branch --depth 1 https://github.com/hyperledger/fabric.git`
   - edit [chainecode.go](content-contract-cc/chainecode.go) for change the chainecode
        - build for verify if working    `go build ./`
   - edit [user.go](content-contract-user/user.go) for change the chainecode
        - use the script to build the new docker image
        
   
## more information 

 to build  [https://github.com/hyperledger/fabric]() we need that
  
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