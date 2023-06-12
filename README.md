## 项目架构

serializer: 序列化器，实现了serializer接口，有Marshal和Unmarshal两个方法

protocol: 协议约定,包括请求体，响应体

codec: \
编解码器模块，codec封装了serializer和TCP连接conn以及带缓冲的writer、reader;\
serverCodec实现:Close、ReadRequestHeader、ReadRequestBody、WriteResponse;\
clientCodec实现Close、ReadResponseHeader、ReadResponseBody、WriteRequest;\
io实现sendFrame和recvFrame将每个消息分为记录消息体长度消息头（利用binary/Uvarint）和消息体，从而解决粘包问题
其中封装了辅助函数read、write用于确定消息长度后进行读写\

cs: 定义客户端和服务端

service: 定义服务



## 项目启动方式

