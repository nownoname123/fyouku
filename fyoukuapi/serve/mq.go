package serve

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Callback func(msg string)

func Connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://test:123456@175.178.212.4:5672/")
	return conn, err
}

// Publish 发送端函数
func Publish(exchange string, queueName string, body string) error {
	conn, err := Connect()
	if err != nil {
		log.Println("连接失败:", err)
		return err
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Println("关闭失败：", err)
		}
	}(conn)

	//创建通道
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)
	//创建队列：
	q, err := channel.QueueDeclare(
		queueName,
		true, //是否持久化
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	//发送消息
	err = channel.Publish(exchange, q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent, //设置为可持久化
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	return err
}

// Consumer 接收者方法
func Consumer(exchange string, queueName string, callback Callback) {
	conn, err := Connect()
	if err != nil {
		log.Println("连接失败:", err)
		fmt.Println("连接失败", err)
		return
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Println("关闭失败：", err)
		}
	}(conn)
	//创建通道
	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		log.Println("consumer channel error：", err)
		return
	}
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)

	//创建队列：
	q, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	//从队列中获取数据
	msgs, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			s := BytesToString(&(d.Body))
			callback(*s)
			err := d.Ack(false)
			if err != nil {
				return
			} //设置为手动应答
		}
	}()
	fmt.Println("waitting for messages")
	<-forever
	//用无缓冲channel起到一直阻塞的作用，消费者需要一直处在监控队列的状态
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}

func PublishEx(exchange string, types string, routingKey string, body string) error {
	conn, err := Connect()
	defer func(conn *amqp.Connection) {
		_ = conn.Close()
	}(conn)
	if err != nil {
		return err
	}
	//创建通道
	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		log.Println("consumer channel error：", err)
		return err
	}
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)
	err = channel.ExchangeDeclare(
		exchange,
		types,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.Publish(exchange, routingKey, false, false, amqp.Publishing{

		ContentType: "text/plain",

		DeliveryMode: amqp.Persistent,

		Body: []byte(body),
	})
	return err
}

func ConsumerEx(exchange string, types string, routingKey string, callback Callback) {
	conn, err := Connect()
	defer func(conn *amqp.Connection) {
		_ = conn.Close()
	}(conn)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		return
	}
	//创建通道
	channel, err := conn.Channel()
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		return
	}
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)
	err = channel.ExchangeDeclare(
		exchange,
		types,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		return
	}
	q, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	//绑定队列和交换机
	err = channel.QueueBind(
		q.Name,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		return
	}
	msgs, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		return
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			s := BytesToString(&(d.Body))
			callback(*s)
			_ = d.Ack(false)
		}
	}()
	fmt.Println("wait for message")
	<-forever
}
