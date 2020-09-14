naivechain
A naive and simple implementation of blockchains.

Build And Run
Download and compile go get -v github.com/kofj/naivechain

Start First Node

naivechain -peers ""
Start Second Node

naivechain -api :3002 -p2p :6002 -peers ws://localhost:6001
HTTP API
query blocks

curl http://localhost:3001/blocks

mine block

curl -H "Content-type:application/json" --data '{"data" : "Some data to the first block"}' http://localhost:3001/mine_block

add peer

curl -H "Content-type:application/json" --data '{"peer" : "ws://localhost:6002"}' http://localhost:3001/add_peer

query peers

curl http://localhost:3001/peers

我产生一个区块,广播出去,广播的内容是这一个区块
peer收到区块后,发现是一个区块,也广播一个消息,查询所有的区块,然后替换自己的区块链,替换完成之后,把最新的区块广播出去
