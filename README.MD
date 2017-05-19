# Contant Fabric 

## python prerequisites

```
backports.ssl-match-hostname==3.4.0.2
cached-property==1.3.0
chardet==2.3.0
cryptography==1.2.3
docker-compose==1.8.0
docker-py==1.9.0
dockerpty==0.4.1
docopt==0.6.2
enum34==1.1.2
funcsigs==0.4
functools32==3.2.3.post2
idna==2.0
ipaddress==1.0.16
jsonschema==2.5.1
mock==1.3.0
ndg-httpsclient==0.4.0
pbr==1.8.0
pyasn1==0.1.9
pyOpenSSL==0.15.1
PyYAML==3.11
requests==2.9.1
six==1.10.0
texttable==0.8.1
urllib3==1.13.1
websocket-client==0.18.0
```
## docker images prerequisites

pull the following images

```
REPOSITORY                      TAG                 IMAGE ID            CREATED             SIZE
dngroup/content-contract-user   latest              cc451b6df2f6        3 months ago        16.3 MB
dngroup/content-contract-cp     latest              465c57f9bf38        3 months ago        16.3 MB
dngroup/content-contract-te     latest              df3db7b135ba        3 months ago        16.3 MB
nherbaut/fabric-peer            latest              21cb00fb27f4        7 months ago        1.42 GB

```
## Install

 We need docker and docker-compose 
 
   1. start peer
   
   It will download the docker "hyperledger/fabric-peer"
   ```
   # on deamon mode
   cd $REPONAME
   docker-compose up -d 
   ```
   2. send chaincode
   
   ```
   curl --request POST \
  --url http://localhost:7050/chaincode \
  --header 'cache-control: no-cache' \
  --header 'content-type: application/json' \
  --header 'postman-token: 356d6fb9-b49c-e62a-2d9b-6dc170309041' \
  --data '{\n  "jsonrpc": "2.0",\n  "method": "deploy",\n  "params": {\n    "type": 1,\n    "chaincodeID": {\n      "path": "https://github.com/dngroup/content-fabric/content-contract-cc"\n    },\n    "ctorMsg": {\n      "function": "init",\n      "args": [\n      ]\n    }\n  },\n  "id": 1\n}'
   ```
   
   2. start cp and te
   
   ```
   ./start $NUMBER_INSTANCE $CHAINE_ID
   ```
   
   2. start user
   
   ````
   ./content-contract-cc/content-contract-cc  -userId lala -contentId bbb.mp4 -chaincodeid 4c166ef7d5220329f658ac55b4a5d3daabbf08cf10f47db5fe3aeb980ac55b2e5a42ac97fb0797f460ad770995cb1c23cec830470d1a3223a1cb885bb3580ae0 -time-max 30
   ````
   
  Working in progress
   
## edit

  We need docker and golang (v1.6)
   
   - maybe you need go dependency `github.com/hyperledger/fabric/`
        - `cd $GOPATH/src/github.com/hyperledger/`
        - `git clone --branch v0.6.1-preview --single-branch --depth 1 https://github.com/hyperledger/fabric.git`
   - edit [content-contract-cc.go](content-contract-cc/content-contract-cc.go) for change the chainecode
        - build for verify if working    `go build ./`
   - edit [content-contract-user.go](content-contract-user/content-contract-user.go) for change the user
        - use the script to build the new docker image
   - edit [content-contract-cp.go](content-contract-cp/content-contract-cp.go) for change the CP
        - use the script to build the new docker image
   - edit [content-contract-te.go](content-contract-te/content-contract-te.goo) for change the TE
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
