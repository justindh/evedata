package nail

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	nsq "github.com/nsqio/go-nsq"
)

// Nail handles storage of data from NSQ
type Nail struct {
	stop    chan bool
	wg      *sync.WaitGroup
	db      *sqlx.DB
	inQueue map[string]*nsq.Consumer
}

// NewNail creates a new storage engine
func NewNail(db *sqlx.DB, addresses []string) *Nail {
	n := &Nail{
		db:      db,
		wg:      &sync.WaitGroup{},
		stop:    make(chan bool),
		inQueue: make(map[string]*nsq.Consumer),
	}

	nsqcfg := nsq.NewConfig()

	for _, h := range handlers {
		c, err := nsq.NewConsumer(h.Topic, "nail", nsqcfg)
		if err != nil {
			log.Fatalln(err)
		}

		h.SpawnFunc(n, c)
		n.inQueue[h.Topic] = c

		err = c.ConnectToNSQLookupds(addresses)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return n
}

// Run the nail service
func (s *Nail) Run() {
	for {
		select {
		case <-s.stop:
			return
		}
	}
}

// Stop the nail service
func (s *Nail) Close() {
	close(s.stop)
	s.wg.Wait()
	for _, h := range s.inQueue {
		h.Stop()
	}
}

// Wrap handlers in a wait group we can properly account during shutdown.
func (s *Nail) wait(next nsq.Handler) nsq.Handler {
	return nsq.HandlerFunc(func(m *nsq.Message) error {
		s.wg.Add(1)
		defer s.wg.Done()
		err := next.HandleMessage(m)
		if err != nil {
			log.Println(err)
			m.Requeue(time.Second)
		} else {
			m.Finish()
		}
		return err
	})
}

type spawnFunc func(s *Nail, consumer *nsq.Consumer)

// Structure for handling routes
type nailHandler struct {
	Topic     string
	SpawnFunc spawnFunc
}

var handlers []nailHandler

func AddHandler(topic string, spawnFunc spawnFunc) {
	handlers = append(handlers, nailHandler{topic, spawnFunc})
}

// RetryTransaction on deadlocks
func RetryTransaction(tx *sqlx.Tx) error {
	for {
		err := tx.Commit()
		if err != nil {
			if !strings.Contains(err.Error(), "1213") {
				return err
			}
			time.Sleep(500 * time.Millisecond)
			continue
		} else {
			return err
		}
	}
}

// DoSQL executes a sql statement
func (s *Nail) DoSQL(stmt string, args ...interface{}) error {
	tx, err := s.db.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.Exec(stmt, args...)
	if err != nil {
		return err
	}

	err = RetryTransaction(tx)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
