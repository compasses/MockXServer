# MockXServer
[![Build Status](https://travis-ci.org/compasses/MockXServer.svg?branch=master)](https://travis-ci.org/compasses/MockXServer)
## Release

## 说明
Auto-record/ Auto-replay API请求；支持oneline 和 offline 两种工作模式，支持http和https方式。
RunMode为online时，作用相当于一个API Proxy，并自动记录请求/响应记录到一个文件数据库：ReplayDB；
RunMode为offline时，作用类似于MockServer，默认从ReplayDB中请求响应，如果找不到对应的响应，然后会根据OfflineHandler配置项是否从offline handler中请求对应的响应；
使用场景常在开发、测试中使用，因能够自动Record和Replay会节省很多准备真实环境的时间。
特别在两个团队分前后台开发的时候，把后台服务直接replay或通过offline handler simulate出来，两个团队之间只进行API接口编程，有助于提升开发效率。

## PACT record & Json 文件生成
1. 通过请求http(s)://MockXServer_IP:Port/pact生成pact文件，目前只支持1.1 的版本
2. 通过请求http(s)://MockXServer_IP:Port/json生成json文件，里面记录了当前record的所有请求/响应记录。

## offline
1.	RESTFul资源服务器，作为离线使用需要完成正常的所有功能。
2. 	自带存储，需要存储可能重复使用的信息。保证功能的完备性

## online
1. 类似API的proxy，直接将请求路由给目的服务器。
2. 目前是主要记录请求和响应信息，http相关的能够完整的记录下来。
3. 支持https，所以这个proxy可以将请求做二次处理，并能记录详细信息。

## config.json说明
1. RunMode， 表示offline和online；
2. TLS，是否起用https模式；
3. RemoteServer，online时，请求的远程服务器；
4. ListenOn, APIService的监听地址;
5. OfflineHandler, 当runMode为offline时，是否启用offline handler响应请求；
5. LogFile, 日志文件名字。留空的话直接打印到命令行窗口。

## 运行
在 Mac或Linux 上执行make 即可，在windows上执行 *sh buil.sh*。可执行文件位于cmd目录下。

## 简单使用说明
1. 打开online模式，运行mock server，会将请求、响应消息信息保存到ReplayDB文件中；
2. 分享给其他开发者ReplayDB，只需配置为offline模式，所有online模式中记录的API都能正常的响应。

### example
一个Get请求的记录：
```
        "/api/CreditCards/v1/getCardsByCustomerId/120016694099968": {
            "GET": [
                {
                    "request": null,
                    "response": {
                        "200": "[{\"id\":120041944596480,\"last4Digits\":\"0005\",\"cardType\":\"AMERICANEXPRESS\",\"nameOnCard\":\"Jet 1\",\"expiryYear\":\"2222\",\"expiryMonth\":\"12\",\"customerId\":120016694099968,\"creationTime\":null,\"updateTime\":null},{\"id\":120043117346816,\"last4Digits\":\"5904\",\"cardType\":\"DINERSCLUB\",\"nameOnCard\":\"Jet 2\",\"expiryYear\":\"2222\",\"expiryMonth\":\"12\",\"customerId\":120016694099968,\"creationTime\":null,\"updateTime\":null},{\"id\":120045496467456,\"last4Digits\":\"1117\",\"cardType\":\"DISCOVER\",\"nameOnCard\":\"Jet 3\",\"expiryYear\":\"2222\",\"expiryMonth\":\"12\",\"customerId\":120016694099968,\"creationTime\":null,\"updateTime\":null},{\"id\":120048522469376,\"last4Digits\":\"4444\",\"cardType\":\"MASTERCARD\",\"nameOnCard\":\"Jet 4\",\"expiryYear\":\"2322\",\"expiryMonth\":\"12\",\"customerId\":120016694099968,\"creationTime\":null,\"updateTime\":null},{\"id\":120049814609920,\"last4Digits\":\"5100\",\"cardType\":\"MASTERCARD\",\"nameOnCard\":\"Jet 5\",\"expiryYear\":\"2323\",\"expiryMonth\":\"11\",\"customerId\":120016694099968,\"creationTime\":null,\"updateTime\":null}]"
                    }
                }
            ]
        }
```

两次Post请求的记录：
```
"/api/CreditCardCheckout/v1/checkout": {
      "POST": [
        {
          "request": {
            "amount": 100,
            "currencyCode": "USD",
            "id": 118650906992640,
            "paymentAccountId": 118628324237312
          },
          "response": {
            "200": "{\"pnref\":\"A10AA261D8C8\",\"paymentResponse\":{\"requestId\":\"FE60350123E15D3A8DF6D44BE67110CF\",\"result\":0,\"respMsg\":\"Approved\",\"status\":true,\"authCode\":\"040PNI\",\"avsAddr\":null,\"avsZip\":null,\"preFpsMsg\":null,\"postFpsMsg\":null,\"transError\":null,\"pnref\":\"A10AA261D8C8\"}}"
          }
        },
        {
          "request": {
            "amount": 222,
            "currencyCode": "USD",
            "id": 118650906992640,
            "paymentAccountId": 118628324237312
          },
          "response": {
            "200": "{\"pnref\":\"A70AA0C6F56C\",\"paymentResponse\":{\"requestId\":\"62B4918913F4A4FE7BB826888BC7B29A\",\"result\":0,\"respMsg\":\"Approved\",\"status\":true,\"authCode\":\"537PNI\",\"avsAddr\":null,\"avsZip\":null,\"preFpsMsg\":null,\"postFpsMsg\":null,\"transError\":null,\"pnref\":\"A70AA0C6F56C\"}}"
          }
        }
      ]
    }
```
上面的结果是通过访问**http(s)://MockXServer_IP:Port/json**获得的json记录。


## 结构图：
![architecture](./architecture.PNG)

整体结构较为简单，MiddleWare层封装了offline和online。offline和online会各自访问DB进行读写。

### Third party lib
1. [httprouter](http://godoc.org/github.com/julienschmidt/httprouter)
2. [boltDB](http://godoc.org/github.com/boltdb/bolt)
3. [PACT-go](https://github.com/SEEK-Jobs/pact-go)
