# go-blockchain
区块链


# 模块
- abci：tendermint core连接的server
- databse：保存tendmint core提交的数据
- types：结构定义
- rpc：对外接口

# 共识
- pbft（使用tendermint集成）

# 语言
- golang

# 目的

一步步教你开发区块链技术，通过最简单的方式实现一个自己区块链demo。

#  用法

1. 启动当前项目
2. 启动tendermint core
3. 发送请求：`curl http://localhost:46657/broadcast_tx_commit?tx=\"abcd\"`

