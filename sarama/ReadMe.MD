# sarama 简单教程

## 本体 docker


## 安装命令行工具
### consumer 
go install github.com/IBM/sarama/tools/kafka-console-consumer@latest
### producer
go install github.com/IBM/sarama/tools/kafka-console-producer@latest

## 关于 Producer 的 AsyncProducer 和 SyncProducer
是用于向Kafka主题发送数据的两种生产者接口。它们的主要区别在于发送数据的方式和处理确认机制。

- 相同点：  
两者都用于将数据发送到Kafka主题。  
都支持批量发送数据，以减少网络开销。  
都提供了异步发送和接收确认的功能。    

- 不同点：  
AsyncProducer采用异步发送模式，不会等待发送确认。发送完成后，生产者将继续处理下一个消息。这种模式适用于对低延迟要求较高的场景。  
SyncProducer采用同步发送模式，会等待发送确认。在发送消息时，生产者会阻塞，直到收到确认或超时。这种模式适用于对数据一致性要求较高的场景。  


- 使用建议：  
对于低延迟场景，推荐使用AsyncProducer。  
对于高一致性需求场景，推荐使用SyncProducer。  
根据实际业务需求和系统架构，灵活选择适合的生产者接口。  

| |AsyncProducer|SyncProducer|
|---|---|---|
|异步发送|√|×|
|同步发送|×|√|
|等待确认|×|√|
|低延迟|√|×|
|数据一致性|×|√|  

- 补充说明:  
在实际应用中，根据业务需求和系统架构，可以选择合适的生产者接口。  
例如，在需要保证数据一致性的场景中，可以使用SyncProducer。  
而在低延迟的场景中，可以使用AsyncProducer。  
此外，还可以结合使用Kafka的消费者组和消费者偏移量管理，以确保高可用性和数据重传。