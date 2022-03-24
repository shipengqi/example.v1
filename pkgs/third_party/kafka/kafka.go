package kafka

import (
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
)

type Config struct {
	*sarama.Config
	Brokers []string
}

type Kafka struct {
	Config *Config
}

func New() *Kafka {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return &Kafka{Config: nil}
}

func (k *Kafka) NewConsumer() (sarama.Consumer, error) {
	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer(k.Config.Brokers, k.Config.Config)
	if err != nil {
		return nil, err
	}
	partitions, err := consumer.Partitions("")
	if err != nil {
		panic(err)
	}
	for partition := range partitions {
		pc, err := consumer.ConsumePartition("", int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}
		defer pc.AsyncClose()
		wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
		wg.Wait()
		consumer.Close()
	}

	return consumer, nil
}

func (k *Kafka) NewProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(k.Config.Brokers, k.Config.Config)
	if err != nil {
		return nil, err
	}
	defer producer.Close()
	var value string
	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder("key"),
		Topic: "go-test",
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Send Message Fail")
	}
	fmt.Printf("Partion = %d, offset = %d\n", partition, offset)

	return producer, nil
}
