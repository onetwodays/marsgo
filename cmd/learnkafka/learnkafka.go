package main

import (
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"

	log "marsgo/utils/ut_log"
)

/*
type deal_t struct {
    timestamp float64
    code      string
    price     string
    amount    string
    side      int
    id        int
    dealMoney string
} */

var (
	side   int = 1
	dealId int = 1
)

const (
	APPNAME = "learnkafka.exe"
	VERSION = "0.0.1"
)

func getSide() int {
	if side == 1 {
		side = 2
	} else {
		side = 1
	}
	return side
}
func getdealId() int {
	dealId = dealId + 1
	return dealId
}

func init() {
	log.InitLog(APPNAME, VERSION)

}

func main() {
	defer log.Flush()

	producer, err := sarama.NewAsyncProducer([]string{"192.168.70.131:9092"}, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Error(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var enqueued, errors int

	rand.Seed(time.Now().UnixNano())
ProducerLoop:
	for {

		/*
		   deal:=&deal_t{
		       timestamp: float64(time.Now().UnixNano()/1000000.0),
		       code:      "BTCUSDT",
		       price:     strconv.Itoa(rand.Intn(400)),
		       amount:    strconv.Itoa(rand.Intn(1000)),
		       side:      getSide(),
		       id:        getdealId(),
		       dealMoney: "100",
		   }*/

		var builder strings.Builder
		builder.WriteString("[")
		//timestamp

		builder.WriteString(strconv.Itoa(int(time.Now().Unix())))

		//code
		builder.WriteString(`,"`)
		builder.WriteString("60000BTCEX0107NA")
		builder.WriteString(`",`)

		//price
		builder.WriteString(strconv.Itoa(rand.Intn(400)))
		builder.WriteString(",")

		//amount
		builder.WriteString(strconv.Itoa(rand.Intn(1400)))
		builder.WriteString(",")

		//side
		builder.WriteString(strconv.Itoa(getSide()))
		builder.WriteString(",")
		//dealid
		builder.WriteString(strconv.Itoa(getdealId()))

		//dealMoney
		builder.WriteString(`,"`)
		builder.WriteString("100")
		builder.WriteString(`"]`)

		text := builder.String()
		select {

		case producer.Input() <- &sarama.ProducerMessage{Topic: "t_me_deal", Partition: 0, Key: nil, Value: sarama.StringEncoder(text)}:
			enqueued++
			log.Info(text)
		case err := <-producer.Errors():
			log.Info("Failed to produce message", err)
			errors++
		case <-signals:
			break ProducerLoop
		}

		time.Sleep(1 * time.Second)
	}

	log.Info("Enqueued: %d; errors: %d\n", enqueued, errors)

}
