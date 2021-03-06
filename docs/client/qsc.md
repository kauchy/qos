# QSC命令行工具

[QSC](../spec/txs/qsc.md)工具包含以下命令:

* `qoscli tx create-qsc`: 创建联盟币，发放联盟币。
* `qoscli tx issue-qsc`: 发行联盟币
* `qoscli query qsc`: 查询qsc信息


## create

> 创建QSC需要申请[CA]()

1. 创建QSC

```
$ qoscli tx create-qsc --help
create qsc

Usage:
  qoscli tx create-qsc [flags]

Flags:
      --async                 broadcast transactions asynchronously
      --chain-id string       Chain ID of tendermint node
      --creator string        name or address of creator
      --desc string           description
      --extrate string        extrate: qos:qscxxx (default "1:280.0000")
  -h, --help                  help for create-qsc
      --indent                add indent to json response
      --max-gas int           gas limit to set per tx
      --node string           <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --nonce int             account nonce to sign the tx
      --nonce-node string     tcp://<host>:<port> to tendermint rpc interface for some chain to query account nonce
      --qcp                   enable qcp mode. send qcp tx
      --qcp-blockheight int   qcp mode flag. original tx blockheight, blockheight must greater than 0
      --qcp-extends string    qcp mode flag. qcp tx extends info
      --qcp-from string       qcp mode flag. qcp tx source chainID
      --qcp-seq int           qcp mode flag.  qcp in sequence
      --qcp-signer string     qcp mode flag. qcp tx signer key name
      --qcp-txindex int       qcp mode flag. original tx index
      --qsc.crt string        path of CA(qsc)
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- creator       创建账号
- qsc.crt       证书位置
- accounts      初始发放地址币值集合，[addr1],[amount];[addr2],[amount2],...，eg：address1vkl6nc6eedkxwjr5rsy2s5jr7qfqm487wu95w7,100;address1vkl6nc6eedkxwjr5rsy2s5jr7qfqm487wu95w7,100。
该参数可为空，即只创建联盟币

> 可以通过`qoscli keys import`导入*creator*账户

```
$ qoscli tx create-qsc --creator qosInitAcc --qsc.crt "qsc.crt"
```

2. 查询QOS绑定的QSCs

```
$ qoscli query store --path /store/qsc/subspace --data qsc
```

## query
```
$ qoscli query qsc --help
query qsc info by name

Usage:
  qoscli query qsc [qsc] [flags]

Flags:
      --chain-id string   Chain ID of tendermint node
      --height int        block height to query, omit to get most recent provable block
  -h, --help              help for qsc
      --indent            add indent to json response
      --node string       <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --trust-node        Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- qsc

```
$ qoscli query qsc QSC
```

## issue

```
$ qoscli tx issue-qsc --help
issue qsc

Usage:
  qoscli tx issue-qsc [flags]

Flags:
      --amount int            coin amount send to banker (default 100000)
      --async                 broadcast transactions asynchronously
      --banker string         address or name of banker
      --chain-id string       Chain ID of tendermint node
  -h, --help                  help for issue-qsc
      --indent                add indent to json response
      --max-gas int           gas limit to set per tx
      --node string           <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --nonce int             account nonce to sign the tx
      --nonce-node string     tcp://<host>:<port> to tendermint rpc interface for some chain to query account nonce
      --qcp                   enable qcp mode. send qcp tx
      --qcp-blockheight int   qcp mode flag. original tx blockheight, blockheight must greater than 0
      --qcp-extends string    qcp mode flag. qcp tx extends info
      --qcp-from string       qcp mode flag. qcp tx source chainID
      --qcp-seq int           qcp mode flag.  qcp in sequence
      --qcp-signer string     qcp mode flag. qcp tx signer key name
      --qcp-txindex int       qcp mode flag. original tx index
      --qsc-name string       qsc name
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- qsc-name  qsc名
- banker    banker账户名
- amount    qsc币值

> 可以通过`qoscli keys import QSCBanker --file ~/banker.pri` 使用banker的私钥文件导入*QSCBanker*账户


导入QSCBanker密钥

```
$ qoscli keys import QSCBanker --file ~/banker.pri
> Enter a passphrase for your key:
> Repeat the passphrase:
```

发放联盟币

```
$ qoscli tx issue-qsc --qsc-name=QSC --banker=QSCBanker --amount=10000
```

查询账户信息:

```
$ qoscli query account QSCBanker
```

