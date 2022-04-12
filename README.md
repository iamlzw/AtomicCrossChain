# AtomicCrossChain
### HLF与以太坊之间的原子跨链实现

本文是HLF与以太坊之间的原子跨链实现

参考https://docs.google.com/presentation/d/1QBP0EbJ94lp2bjRpRqS_CzzvvmY-gOkUNGCTmF7jKsI/edit#slide=id.g57af0cd0aa_2_32

假设Alice与Bob想在没有信任的第三方参与的场景下进行资产交易

Alice向Bob转1ETH，Bob向Alice转1个HLF上的Token(这里假设1Token与1ETH价值相同的)，那么应当怎样实现？

参考文章中给了一个HLF与以太坊之间的原子跨链实现。

首先有两条链，HLF以及以太坊，分别部署了HTLC合约。

首先Alice选择一个key并计算H(key),之后将H(key)告诉Bob,并与Bob商议一个TimeLock(就是说这笔交易在某个时间之前是有效的)。

Alice在以太坊上创建一笔交易向Bob转1ETH，这笔交易在特定情况下会完全执行(输入正确的key值)。否则处于半执行状态，即这1ETH既没有在Alice的账户中，也没有在Bob的账户中， 通常由智能合约或多重签名钱包持有资产。当Bob在截至时间之前输入正确的key,那么这1ETH将转入Bob的账户中。

Bob在HLF上创建一笔交易向Alice转1 HLF Token,这笔交易在特定情况下会完全执行(输入正确的key值)。否则处于半执行状态，即这1 HLF Token既没有在Bob的账户中，也没有在Alice的账户中， 通常由智能合约或多重签名钱包持有资产。当Alice在截至时间之前输入正确的key,那么这1 Token将转入Alice的账户中。

### 部署hlf网络

hlf网络基于1.4.2版本部署。

克隆仓库

```sh
git clone https://github.com/iamlzw/AtomicCrossChain.git
```

启动hlf网络

```sh
cd AtomicCrossChain/hlf
./start.sh
## 在部署智能合约之前需要将合约复制到gopath下，否则会安装合约失败
cp -r chaincode ${GOPATH}/src
./init.sh
```

测试合约方法

```sh
## 设置环境变量
export PATH=$(pwd)/bin:$PATH
export FABRIC_CFG_PATH=$(pwd)/config/

export CORE_PEER_TLS_ENABLED=true

export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$(pwd)/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=$(pwd)/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051

export ORDERER_CA_PATH=$(pwd)/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

## 初始化账本
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile ${ORDERER_CA_PATH} -C mychannel -n mycc -c '{"Args":["InitLedger"]}'

## 执行转账
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile ${ORDERER_CA_PATH} -C mychannel -n mycc -c '{"Args":["transfer","Bob","Alice","50"]}'

## 查询Alice资产
peer chaincode query -C mychannel -n mycc -c '{"Args":["GetBalance","Alice"]}'

## 查询Bob资产
peer chaincode query -C mychannel -n mycc -c '{"Args":["GetBalance","Bob"]}'
```

### 部署以太坊测试网络并部署htlc合约

参考

[以太坊入门(一)环境安装](http://lifegoeson.cn/2022/03/17/以太坊入门(一)环境安装/)

[以太坊入门(二)启动一个测试网络](http://lifegoeson.cn/2022/03/17/以太坊入门(二)启动一个测试网络/)

[以太坊入门(三) 在本地部署remix](http://lifegoeson.cn/2022/03/17/以太坊入门(三)在本地部署remix/)

将Ethereum路径下的HTLC合约通过remix部署到以太坊测试网络上。这里我们将测试网络的第一个账户作为Alice的地址，第二个账户作为Bob的地址。三个参数分别是_TOBOB(Bob的地址)，HASHLOCK(key的hash值，我们这里选择htlc作为key，则hashlock为ef8c49cab8ca567b21b21ead803e0ff1238c13ba53ece34349bbd497f9bbb121,这个key的值目前只有Alice知道),TIMEOUT(交易的有效时间,值为1647783730，也是就是截止到2022/03/20 13:42:10)

![image.png](http://lifegoeson.cn:8888/images/2022/03/19/image.png)

至此hlf和以太坊网络部署完成，合约也成功部署。下面我们模拟Alice与Bob在没有信任的第三方参与的场景下进行资产交易。

### Alice创建交易

Alice在以太坊上创建一笔交易向Bob转1ETH。在remix上将VALUE设置为1ETH,点击`Transact`

### Bob创建交易

Bob在hlf上创建一笔交易向Alice转1Token,打开hlf服务器终端

```sh
## 设置环境变量
export PATH=$(pwd)/bin:$PATH
export FABRIC_CFG_PATH=$(pwd)/config/

export CORE_PEER_TLS_ENABLED=true

export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$(pwd)/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=$(pwd)/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051

export ORDERER_CA_PATH=$(pwd)/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

## 创建交易,这里的hashlock是由Alice告诉Bob的，目前Bob不知道该hashlock的key是什么。
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile ${ORDERER_CA_PATH} -C mychannel -n mycc -c '{"Args":["conditional","htlc111", "Bob","Alice","1", "ef8c49cab8ca567b21b21ead803e0ff1238c13ba53ece34349bbd497f9bbb121", "2022-11-12T11:45:26.371Z"]}'
```

### Alice在hlf上确认交易

打开hlf服务器终端

```sh
## 设置环境变量
export PATH=$(pwd)/bin:$PATH
export FABRIC_CFG_PATH=$(pwd)/config/

export CORE_PEER_TLS_ENABLED=true

export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$(pwd)/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=$(pwd)/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051

export ORDERER_CA_PATH=$(pwd)/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

## 创建交易,这里的hashlock是由Alice告诉Bob的，目前Bob不知道该hashlock的key是什么。
peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile ${ORDERER_CA_PATH} -C mychannel -n mycc -c '{"Args":["Commit","htlc111","htlc"]}'
## 至此Alice获取到Bob向其转的1 Token，而Bob也通过hlf上的交易监听事件获取到hashlock的key值
```

### Bob在以太坊上确认交易

在remix上执行commit方法，参数是通过hlf上的交易监听事件获取到hashlock的key值，也就是htlc,至此Bob获取到Alice转给其的1ETH。
