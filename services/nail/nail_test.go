package nail

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/evedata/internal/nsqhelper"
	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/sqlhelper"
	"github.com/antihax/evedata/services/hammer"
	nsq "github.com/nsqio/go-nsq"
	"golang.org/x/oauth2"
)

var (
	testWork = []redisqueue.Work{
		{Operation: "marketOrders", Parameter: int32(22)},
		{Operation: "structure", Parameter: int64(1000000017013)},
		{Operation: "structure", Parameter: int64(1000000025062)},
		{Operation: "structureOrders", Parameter: int64(1000000017013)},
		{Operation: "killmail", Parameter: []interface{}{"FAKEHASH", int32(56271)}},
		{Operation: "marketHistoryTrigger", Parameter: int32(1)},
		{Operation: "marketHistory", Parameter: []int32{1, 1}},
		{Operation: "marketOrders", Parameter: int32(1)},
		{Operation: "war", Parameter: int32(1)},
		{Operation: "alliance", Parameter: int32(1)},
		{Operation: "corporation", Parameter: int32(1)},
		{Operation: "character", Parameter: int32(1)},
		{Operation: "characterWalletTransactions", Parameter: []interface{}{int32(1), int32(1)}},
	}
	ham          *hammer.Hammer
	nailInstance *Nail
)

func TestMain(m *testing.M) {
	sql := sqlhelper.NewTestDatabase()

	redis := redigohelper.ConnectRedisTestPool()
	redConn := redis.Get()
	defer redConn.Close()
	redConn.Do("FLUSHALL")

	producer, err := nsqhelper.NewTestNSQProducer()
	if err != nil {
		log.Fatalln(err)
	}

	ham = hammer.NewHammer(redis, sql, producer, "123400", "faaaaaaake", "sofake", "faaaaaaake", "sofake")
	ham.ChangeBasePath("http://127.0.0.1:8080")
	ham.ChangeTokenPath("http://127.0.0.1:8080")
	ham.SetToken(1, 1, &oauth2.Token{
		RefreshToken: "fake",
		AccessToken:  "really fake",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	})

	nailInstance = NewNail(sql, nsqhelper.Test)

	go ham.Run()
	go nailInstance.Run()
	retCode := m.Run()
	time.Sleep(time.Second * 5)
	nailInstance.Close()
	ham.Close()
	redis.Close()
	sql.Close()

	os.Exit(retCode)
}

func TestQueue(t *testing.T) {
	err := ham.QueueWork(testWork)
	if err != nil {
		t.Fatal(err)
	}
}

func TestKillmails(t *testing.T) {
	m := &nsq.Message{
		ID: nsq.MessageID{0x30, 0x38, 0x62, 0x62, 0x32, 0x35, 0x61, 0x64, 0x33, 0x65, 0x62, 0x38, 0x35, 0x30, 0x30, 0x30},
		Body: []uint8{0xff, 0x95, 0xff, 0x8b, 0x3, 0x1, 0x1, 0x24, 0x47, 0x65, 0x74, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c,
			0x73, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x48, 0x61, 0x73,
			0x68, 0x4f, 0x6b, 0x1, 0xff, 0x8c, 0x0, 0x1, 0x7, 0x1, 0x9, 0x41, 0x74, 0x74, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x73, 0x1, 0xff, 0x90,
			0x0, 0x1, 0xa, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xc, 0x4b, 0x69, 0x6c, 0x6c, 0x6d,
			0x61, 0x69, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x92, 0x0, 0x1, 0x6, 0x4d, 0x6f, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1,
			0xd, 0x53, 0x6f, 0x6c, 0x61, 0x72, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x6, 0x56, 0x69, 0x63,
			0x74, 0x69, 0x6d, 0x1, 0xff, 0x94, 0x0, 0x1, 0x5, 0x57, 0x61, 0x72, 0x49, 0x64, 0x1, 0x4, 0x0, 0x0, 0x0, 0x3f, 0xff, 0x8f, 0x2,
			0x1, 0x1, 0x30, 0x5b, 0x5d, 0x65, 0x73, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x73,
			0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x48, 0x61, 0x73,
			0x68, 0x41, 0x74, 0x74, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x1, 0xff, 0x90, 0x0, 0x1, 0xff, 0x8e, 0x0, 0x0, 0xff, 0xc7, 0xff,
			0x8d, 0x3, 0x1, 0x1, 0x2a, 0x47, 0x65, 0x74, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x4b, 0x69, 0x6c, 0x6c,
			0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x48, 0x61, 0x73, 0x68, 0x41, 0x74,
			0x74, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x1, 0xff, 0x8e, 0x0, 0x1, 0x9, 0x1, 0xa, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x49,
			0x64, 0x1, 0x4, 0x0, 0x1, 0xb, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xd, 0x43, 0x6f,
			0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xa, 0x44, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x44,
			0x6f, 0x6e, 0x65, 0x1, 0x4, 0x0, 0x1, 0x9, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x9, 0x46, 0x69,
			0x6e, 0x61, 0x6c, 0x42, 0x6c, 0x6f, 0x77, 0x1, 0x2, 0x0, 0x1, 0xe, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x53, 0x74, 0x61,
			0x74, 0x75, 0x73, 0x1, 0x8, 0x0, 0x1, 0xa, 0x53, 0x68, 0x69, 0x70, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xc, 0x57,
			0x65, 0x61, 0x70, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x0, 0x0, 0x10, 0xff, 0x91, 0x5, 0x1, 0x1, 0x4, 0x54,
			0x69, 0x6d, 0x65, 0x1, 0xff, 0x92, 0x0, 0x0, 0x0, 0xff, 0xad, 0xff, 0x93, 0x3, 0x1, 0x1, 0x28, 0x47, 0x65, 0x74, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x48, 0x61, 0x73, 0x68, 0x56, 0x69, 0x63, 0x74, 0x69, 0x6d, 0x1, 0xff, 0x94, 0x0, 0x1, 0x8, 0x1, 0xa, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xb, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xd, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xb, 0x44, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x54, 0x61, 0x6b, 0x65, 0x6e, 0x1, 0x4, 0x0, 0x1, 0x9, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x5, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x1, 0xff, 0x9c, 0x0, 0x1, 0x8, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x9e, 0x0, 0x1, 0xa, 0x53, 0x68, 0x69, 0x70, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x0, 0x0, 0x3c, 0xff, 0x9b, 0x2, 0x1, 0x1, 0x2d, 0x5b, 0x5d, 0x65, 0x73, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x48, 0x61, 0x73, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x31, 0x1, 0xff, 0x9c, 0x0, 0x1, 0xff, 0x96, 0x0, 0x0, 0xff, 0x90, 0xff, 0x95, 0x3, 0x1, 0x1, 0x27, 0x47, 0x65, 0x74, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x48, 0x61, 0x73, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x31, 0x1, 0xff, 0x96, 0x0, 0x1, 0x6, 0x1, 0x4, 0x46, 0x6c, 0x61, 0x67, 0x1, 0x4, 0x0, 0x1, 0xa, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x5, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x1, 0xff, 0x9a, 0x0, 0x1, 0x11, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x1, 0x4, 0x0, 0x1, 0xf, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x44, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x1, 0x4, 0x0, 0x1, 0x9, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x74, 0x6f, 0x6e, 0x1, 0x4, 0x0, 0x0, 0x0, 0x3b, 0xff, 0x99, 0x2, 0x1, 0x1, 0x2c, 0x5b, 0x5d, 0x65, 0x73, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x48, 0x61, 0x73, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x1, 0xff, 0x9a, 0x0, 0x1, 0xff, 0x98, 0x0, 0x0, 0xff, 0x84, 0xff, 0x97, 0x3, 0x1, 0x1, 0x26, 0x47, 0x65, 0x74, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x48, 0x61, 0x73, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x1, 0xff, 0x98, 0x0, 0x1, 0x5, 0x1, 0x4, 0x46, 0x6c, 0x61, 0x67, 0x1, 0x4, 0x0, 0x1, 0xa, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x11, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x1, 0x4, 0x0, 0x1, 0xf, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x44, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x1, 0x4, 0x0, 0x1, 0x9, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x74, 0x6f, 0x6e, 0x1, 0x4, 0x0, 0x0, 0x0, 0x4a, 0xff, 0x9d, 0x3, 0x1, 0x1, 0x2a, 0x47, 0x65, 0x74, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x49, 0x64, 0x4b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x48, 0x61, 0x73, 0x68, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x9e, 0x0, 0x1, 0x3, 0x1, 0x1, 0x58, 0x1, 0x8, 0x0, 0x1, 0x1, 0x59, 0x1, 0x8, 0x0, 0x1, 0x1, 0x5a, 0x1, 0x8, 0x0, 0x0, 0x0, 0xff, 0x84, 0xff, 0x8c, 0x1, 0x1, 0x2, 0xfc, 0xb, 0x6b, 0xeb, 0x0, 0x1, 0xfd, 0x1e, 0x85, 0xe6, 0x1, 0xfe, 0x2c, 0xe2, 0x1, 0xfd, 0xf, 0x42, 0x46, 0x1, 0x1, 0x1, 0xfb, 0x40, 0x33, 0x33, 0xd3, 0xbf, 0x1, 0xfe, 0x8b, 0x62, 0x1, 0xfe, 0x18, 0x4, 0x0, 0x1, 0xfc, 0x6, 0xc3, 0x60, 0xfa, 0x1, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xcf, 0x9d, 0x95, 0x40, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x2, 0xfc, 0x3, 0x93, 0x9e, 0x40, 0x1, 0x1, 0xfc, 0x4a, 0x11, 0xbf, 0x74, 0x1, 0xfc, 0xb, 0xf, 0xea, 0xa2, 0x1, 0xfc, 0x64, 0x4c, 0x61, 0xae, 0x1, 0xfe, 0x2c, 0xe2, 0x2, 0x1, 0x1, 0x28, 0x1, 0xfe, 0x2e, 0xaa, 0x3, 0x2, 0x0, 0x1, 0x1, 0xfc, 0x1b, 0x52, 0x5a, 0x42, 0x1, 0xfc, 0x26, 0x14, 0x41, 0x42, 0x1, 0xfb, 0x20, 0x94, 0x7f, 0x39, 0x42, 0x0, 0x1, 0xfe, 0x8b, 0x28, 0x0, 0x0},
		Timestamp: 1508725921739026922,
		Attempts:  0x1,
	}
	err := nailInstance.killmailHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMarket(t *testing.T) {
	m := &nsq.Message{
		ID:        nsq.MessageID{0x30, 0x38, 0x62, 0x62, 0x32, 0x36, 0x63, 0x30, 0x62, 0x35, 0x33, 0x38, 0x35, 0x30, 0x30, 0x30},
		Body:      []uint8{0x33, 0xff, 0x99, 0x3, 0x1, 0x1, 0xc, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x1, 0xff, 0x9a, 0x0, 0x1, 0x2, 0x1, 0x6, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x1, 0xff, 0x9e, 0x0, 0x1, 0x8, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x1, 0x4, 0x0, 0x0, 0x0, 0x32, 0xff, 0x9d, 0x2, 0x1, 0x1, 0x23, 0x5b, 0x5d, 0x65, 0x73, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x73, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x32, 0x30, 0x30, 0x4f, 0x6b, 0x1, 0xff, 0x9e, 0x0, 0x1, 0xff, 0x9c, 0x0, 0x0, 0xff, 0xbd, 0xff, 0x9b, 0x3, 0x1, 0x1, 0x1d, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x73, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x32, 0x30, 0x30, 0x4f, 0x6b, 0x1, 0xff, 0x9c, 0x0, 0x1, 0xb, 0x1, 0x8, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0x4, 0x0, 0x1, 0xa, 0x49, 0x73, 0x42, 0x75, 0x79, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x1, 0x2, 0x0, 0x1, 0x6, 0x49, 0x73, 0x73, 0x75, 0x65, 0x64, 0x1, 0xff, 0x8c, 0x0, 0x1, 0xa, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x9, 0x4d, 0x69, 0x6e, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x1, 0x4, 0x0, 0x1, 0x7, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x5, 0x50, 0x72, 0x69, 0x63, 0x65, 0x1, 0x8, 0x0, 0x1, 0x6, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x1, 0xc, 0x0, 0x1, 0x6, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xc, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x1, 0x4, 0x0, 0x1, 0xb, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x1, 0x4, 0x0, 0x0, 0x0, 0x10, 0xff, 0x8b, 0x5, 0x1, 0x1, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x8c, 0x0, 0x0, 0x0, 0x46, 0xff, 0x9a, 0x1, 0x1, 0x1, 0xff, 0xb4, 0x2, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xcf, 0x5c, 0x52, 0xb9, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x1, 0xfc, 0x7, 0x27, 0x39, 0xbe, 0x1, 0x2, 0x1, 0xfb, 0x2, 0x27, 0x33, 0xea, 0xbe, 0x1, 0xfb, 0xc0, 0xcc, 0xcc, 0x23, 0x40, 0x1, 0x6, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x1, 0x44, 0x1, 0xfd, 0x27, 0x8d, 0x0, 0x1, 0xfd, 0x3d, 0x9, 0x0, 0x0, 0x1, 0x2c, 0x0},
		Timestamp: 1508726217514270341,
		Attempts:  0x1,
	}
	err := nailInstance.marketOrderHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStructure(t *testing.T) {
	m := &nsq.Message{
		ID:        nsq.MessageID{0x30, 0x38, 0x63, 0x30, 0x30, 0x31, 0x64, 0x38, 0x62, 0x65, 0x33, 0x38, 0x35, 0x30, 0x30, 0x30},
		Body:      []uint8{0x36, 0xff, 0x8d, 0x3, 0x1, 0x1, 0x9, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x1, 0xff, 0x8e, 0x0, 0x1, 0x2, 0x1, 0x9, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x1, 0xff, 0x90, 0x0, 0x1, 0xb, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x49, 0x44, 0x1, 0x4, 0x0, 0x0, 0x0, 0x64, 0xff, 0x8f, 0x3, 0x1, 0x1, 0x22, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x69, 0x76, 0x65, 0x72, 0x73, 0x65, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x73, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x49, 0x64, 0x4f, 0x6b, 0x1, 0xff, 0x90, 0x0, 0x1, 0x4, 0x1, 0x4, 0x4e, 0x61, 0x6d, 0x65, 0x1, 0xc, 0x0, 0x1, 0x8, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x92, 0x0, 0x1, 0xd, 0x53, 0x6f, 0x6c, 0x61, 0x72, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x6, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x0, 0x0, 0x48, 0xff, 0x91, 0x3, 0x1, 0x1, 0x28, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x69, 0x76, 0x65, 0x72, 0x73, 0x65, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x73, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x49, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x92, 0x0, 0x1, 0x3, 0x1, 0x1, 0x58, 0x1, 0x8, 0x0, 0x1, 0x1, 0x59, 0x1, 0x8, 0x0, 0x1, 0x1, 0x5a, 0x1, 0x8, 0x0, 0x0, 0x0, 0x2e, 0xff, 0x8e, 0x1, 0x1, 0x17, 0x56, 0x2d, 0x33, 0x59, 0x47, 0x37, 0x20, 0x56, 0x49, 0x20, 0x2d, 0x20, 0x54, 0x68, 0x65, 0x20, 0x43, 0x61, 0x70, 0x69, 0x74, 0x61, 0x6c, 0x1, 0x0, 0x1, 0xfc, 0x3, 0x93, 0x88, 0x1c, 0x0, 0x1, 0xfa, 0x1, 0xd1, 0xa9, 0x4a, 0xe3, 0xcc, 0x0},
		Timestamp: 1509067916560436388,
		Attempts:  0x1,
	}

	err := nailInstance.structureHandler(m)
	if err != nil {
		t.Fatal(err)
	}

	m = &nsq.Message{
		ID:        nsq.MessageID{0x30, 0x38, 0x63, 0x30, 0x30, 0x31, 0x64, 0x38, 0x62, 0x65, 0x37, 0x38, 0x35, 0x30, 0x30, 0x30},
		Body:      []uint8{0x36, 0xff, 0x8d, 0x3, 0x1, 0x1, 0x9, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x1, 0xff, 0x8e, 0x0, 0x1, 0x2, 0x1, 0x9, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x1, 0xff, 0x90, 0x0, 0x1, 0xb, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x49, 0x44, 0x1, 0x4, 0x0, 0x0, 0x0, 0x64, 0xff, 0x8f, 0x3, 0x1, 0x1, 0x22, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x69, 0x76, 0x65, 0x72, 0x73, 0x65, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x73, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x49, 0x64, 0x4f, 0x6b, 0x1, 0xff, 0x90, 0x0, 0x1, 0x4, 0x1, 0x4, 0x4e, 0x61, 0x6d, 0x65, 0x1, 0xc, 0x0, 0x1, 0x8, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x92, 0x0, 0x1, 0xd, 0x53, 0x6f, 0x6c, 0x61, 0x72, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x6, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x0, 0x0, 0x48, 0xff, 0x91, 0x3, 0x1, 0x1, 0x28, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x69, 0x76, 0x65, 0x72, 0x73, 0x65, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x73, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x49, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x92, 0x0, 0x1, 0x3, 0x1, 0x1, 0x58, 0x1, 0x8, 0x0, 0x1, 0x1, 0x59, 0x1, 0x8, 0x0, 0x1, 0x1, 0x5a, 0x1, 0x8, 0x0, 0x0, 0x0, 0x2e, 0xff, 0x8e, 0x1, 0x1, 0x17, 0x56, 0x2d, 0x33, 0x59, 0x47, 0x37, 0x20, 0x56, 0x49, 0x20, 0x2d, 0x20, 0x54, 0x68, 0x65, 0x20, 0x43, 0x61, 0x70, 0x69, 0x74, 0x61, 0x6c, 0x1, 0x0, 0x1, 0xfc, 0x3, 0x93, 0x88, 0x1c, 0x0, 0x1, 0xfa, 0x1, 0xd1, 0xa9, 0x4a, 0xa4, 0xea, 0x0},
		Timestamp: 1509067916561402790,
		Attempts:  0x1,
	}
	err = nailInstance.structureHandler(m)
	if err != nil {
		t.Fatal(err)
	}

	m = &nsq.Message{
		ID:        nsq.MessageID{0x30, 0x38, 0x62, 0x64, 0x39, 0x61, 0x37, 0x30, 0x65, 0x35, 0x37, 0x38, 0x35, 0x30, 0x30, 0x30},
		Body:      []uint8{0x39, 0xff, 0xa5, 0x3, 0x1, 0x1, 0xf, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x1, 0xff, 0xa6, 0x0, 0x1, 0x2, 0x1, 0x6, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x1, 0xff, 0xaa, 0x0, 0x1, 0xb, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x49, 0x44, 0x1, 0x4, 0x0, 0x0, 0x0, 0x39, 0xff, 0xa9, 0x2, 0x1, 0x1, 0x2a, 0x5b, 0x5d, 0x65, 0x73, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x73, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x73, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x49, 0x64, 0x32, 0x30, 0x30, 0x4f, 0x6b, 0x1, 0xff, 0xaa, 0x0, 0x1, 0xff, 0xa8, 0x0, 0x0, 0xff, 0xc4, 0xff, 0xa7, 0x3, 0x1, 0x1, 0x24, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x73, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x73, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x49, 0x64, 0x32, 0x30, 0x30, 0x4f, 0x6b, 0x1, 0xff, 0xa8, 0x0, 0x1, 0xb, 0x1, 0x8, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0x4, 0x0, 0x1, 0xa, 0x49, 0x73, 0x42, 0x75, 0x79, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x1, 0x2, 0x0, 0x1, 0x6, 0x49, 0x73, 0x73, 0x75, 0x65, 0x64, 0x1, 0xff, 0x8a, 0x0, 0x1, 0xa, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x9, 0x4d, 0x69, 0x6e, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x1, 0x4, 0x0, 0x1, 0x7, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x5, 0x50, 0x72, 0x69, 0x63, 0x65, 0x1, 0x8, 0x0, 0x1, 0x6, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x1, 0xc, 0x0, 0x1, 0x6, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xc, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x1, 0x4, 0x0, 0x1, 0xb, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x1, 0x4, 0x0, 0x0, 0x0, 0x10, 0xff, 0x89, 0x5, 0x1, 0x1, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x8a, 0x0, 0x0, 0x0, 0x46, 0xff, 0xa6, 0x1, 0x1, 0x1, 0xff, 0xb4, 0x2, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xcf, 0x5c, 0x52, 0xb9, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x1, 0xfc, 0x7, 0x27, 0x39, 0xbe, 0x1, 0x2, 0x1, 0xfb, 0x2, 0x27, 0x33, 0xea, 0xbe, 0x1, 0xfb, 0xc0, 0xcc, 0xcc, 0x23, 0x40, 0x1, 0x6, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x1, 0x44, 0x1, 0xfd, 0x27, 0x8d, 0x0, 0x1, 0xfd, 0x3d, 0x9, 0x0, 0x0, 0x1, 0x16, 0x0},
		Timestamp: 1508898755143235561,
		Attempts:  0x1,
	}

	// Hack to use an actual structureID.
	b := datapackages.MarketOrders{}
	gobcoder.GobDecoder(m.Body, &b)
	for i := range b.Orders {
		b.Orders[i].LocationId = 1000000017013
	}
	newb, _ := gobcoder.GobEncoder(b)
	m.Body = newb

	err = nailInstance.structureMarketHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMarketHistory(t *testing.T) {
	m := &nsq.Message{
		ID:        nsq.MessageID{0x30, 0x38, 0x63, 0x30, 0x31, 0x34, 0x65, 0x38, 0x32, 0x36, 0x66, 0x38, 0x35, 0x30, 0x30, 0x30},
		Body:      []uint8{0x40, 0xff, 0x95, 0x3, 0x1, 0x1, 0xd, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x1, 0xff, 0x96, 0x0, 0x1, 0x3, 0x1, 0x7, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x1, 0xff, 0x9a, 0x0, 0x1, 0x8, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x1, 0x4, 0x0, 0x1, 0x6, 0x54, 0x79, 0x70, 0x65, 0x49, 0x44, 0x1, 0x4, 0x0, 0x0, 0x0, 0x33, 0xff, 0x99, 0x2, 0x1, 0x1, 0x24, 0x5b, 0x5d, 0x65, 0x73, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x73, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x32, 0x30, 0x30, 0x4f, 0x6b, 0x1, 0xff, 0x9a, 0x0, 0x1, 0xff, 0x98, 0x0, 0x0, 0x72, 0xff, 0x97, 0x3, 0x1, 0x1, 0x1e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x73, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x32, 0x30, 0x30, 0x4f, 0x6b, 0x1, 0xff, 0x98, 0x0, 0x1, 0x6, 0x1, 0x7, 0x41, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x1, 0x8, 0x0, 0x1, 0x4, 0x44, 0x61, 0x74, 0x65, 0x1, 0xc, 0x0, 0x1, 0x7, 0x48, 0x69, 0x67, 0x68, 0x65, 0x73, 0x74, 0x1, 0x8, 0x0, 0x1, 0x6, 0x4c, 0x6f, 0x77, 0x65, 0x73, 0x74, 0x1, 0x8, 0x0, 0x1, 0xa, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x1, 0x4, 0x0, 0x1, 0x6, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x1, 0x4, 0x0, 0x0, 0x0, 0x33, 0xff, 0x96, 0x1, 0x1, 0x1, 0xfe, 0x15, 0x40, 0x1, 0xa, 0x32, 0x30, 0x31, 0x35, 0x2d, 0x30, 0x35, 0x2d, 0x30, 0x31, 0x1, 0xfb, 0xe0, 0x7a, 0x14, 0x15, 0x40, 0x1, 0xfb, 0xe0, 0xa3, 0x70, 0x14, 0x40, 0x1, 0xfe, 0x11, 0xb6, 0x1, 0xfb, 0x7, 0x94, 0x57, 0xf7, 0xa6, 0x0, 0x1, 0x2, 0x1, 0x2, 0x0},
		Timestamp: 1509073155786320746,
		Attempts:  0x1,
	}

	// Hack to have the latest date.
	b := datapackages.MarketHistory{}
	gobcoder.GobDecoder(m.Body, &b)
	for i := range b.History {
		b.History[i].Date = time.Now().Format("2006-01-02")
	}
	newb, _ := gobcoder.GobEncoder(b)
	m.Body = newb

	err := nailInstance.marketHistoryHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWar(t *testing.T) {
	m := &nsq.Message{
		ID:        nsq.MessageID{0x30, 0x38, 0x63, 0x32, 0x38, 0x64, 0x62, 0x33, 0x37, 0x65, 0x33, 0x38, 0x35, 0x30, 0x30, 0x30},
		Body:      []uint8{0xff, 0xa1, 0xff, 0xa7, 0x3, 0x1, 0x1, 0xe, 0x47, 0x65, 0x74, 0x57, 0x61, 0x72, 0x73, 0x57, 0x61, 0x72, 0x49, 0x64, 0x4f, 0x6b, 0x1, 0xff, 0xa8, 0x0, 0x1, 0xa, 0x1, 0x9, 0x41, 0x67, 0x67, 0x72, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x1, 0xff, 0xaa, 0x0, 0x1, 0x6, 0x41, 0x6c, 0x6c, 0x69, 0x65, 0x73, 0x1, 0xff, 0xae, 0x0, 0x1, 0x8, 0x44, 0x65, 0x63, 0x6c, 0x61, 0x72, 0x65, 0x64, 0x1, 0xff, 0x8e, 0x0, 0x1, 0x8, 0x44, 0x65, 0x66, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x1, 0xff, 0xb0, 0x0, 0x1, 0x8, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x1, 0xff, 0x8e, 0x0, 0x1, 0x2, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x6, 0x4d, 0x75, 0x74, 0x75, 0x61, 0x6c, 0x1, 0x2, 0x0, 0x1, 0xd, 0x4f, 0x70, 0x65, 0x6e, 0x46, 0x6f, 0x72, 0x41, 0x6c, 0x6c, 0x69, 0x65, 0x73, 0x1, 0x2, 0x0, 0x1, 0x9, 0x52, 0x65, 0x74, 0x72, 0x61, 0x63, 0x74, 0x65, 0x64, 0x1, 0xff, 0x8e, 0x0, 0x1, 0x7, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x1, 0xff, 0x8e, 0x0, 0x0, 0x0, 0x65, 0xff, 0xa9, 0x3, 0x1, 0x1, 0x15, 0x47, 0x65, 0x74, 0x57, 0x61, 0x72, 0x73, 0x57, 0x61, 0x72, 0x49, 0x64, 0x41, 0x67, 0x67, 0x72, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x1, 0xff, 0xaa, 0x0, 0x1, 0x4, 0x1, 0xa, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xd, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xc, 0x49, 0x73, 0x6b, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x1, 0x8, 0x0, 0x1, 0xb, 0x53, 0x68, 0x69, 0x70, 0x73, 0x4b, 0x69, 0x6c, 0x6c, 0x65, 0x64, 0x1, 0x4, 0x0, 0x0, 0x0, 0x25, 0xff, 0xad, 0x2, 0x1, 0x1, 0x16, 0x5b, 0x5d, 0x65, 0x73, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x61, 0x72, 0x73, 0x57, 0x61, 0x72, 0x49, 0x64, 0x41, 0x6c, 0x6c, 0x79, 0x1, 0xff, 0xae, 0x0, 0x1, 0xff, 0xac, 0x0, 0x0, 0x3f, 0xff, 0xab, 0x3, 0x1, 0x1, 0x10, 0x47, 0x65, 0x74, 0x57, 0x61, 0x72, 0x73, 0x57, 0x61, 0x72, 0x49, 0x64, 0x41, 0x6c, 0x6c, 0x79, 0x1, 0xff, 0xac, 0x0, 0x1, 0x2, 0x1, 0xa, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xd, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x0, 0x0, 0x10, 0xff, 0x8d, 0x5, 0x1, 0x1, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x8e, 0x0, 0x0, 0x0, 0x64, 0xff, 0xaf, 0x3, 0x1, 0x1, 0x14, 0x47, 0x65, 0x74, 0x57, 0x61, 0x72, 0x73, 0x57, 0x61, 0x72, 0x49, 0x64, 0x44, 0x65, 0x66, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x1, 0xff, 0xb0, 0x0, 0x1, 0x4, 0x1, 0xa, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xd, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xc, 0x49, 0x73, 0x6b, 0x44, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x1, 0x8, 0x0, 0x1, 0xb, 0x53, 0x68, 0x69, 0x70, 0x73, 0x4b, 0x69, 0x6c, 0x6c, 0x65, 0x64, 0x1, 0x4, 0x0, 0x0, 0x0, 0x28, 0xff, 0xa8, 0x1, 0x2, 0xfc, 0x75, 0x9e, 0xa6, 0x80, 0x0, 0x2, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xb8, 0x40, 0xda, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x1, 0x2, 0xfc, 0x77, 0x65, 0x3f, 0x36, 0x0, 0x2, 0xfe, 0xf, 0x2a, 0x0},
		Timestamp: 1509247096958751788,
		Attempts:  0x1,
	}
	err := nailInstance.warHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCharacter(t *testing.T) {
	m := &nsq.Message{
		ID:        nsq.MessageID{0x30, 0x38, 0x63, 0x33, 0x36, 0x35, 0x33, 0x63, 0x64, 0x66, 0x66, 0x38, 0x35, 0x30, 0x30, 0x30},
		Body:      []uint8{0x4e, 0xff, 0x9f, 0x3, 0x1, 0x1, 0x9, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x1, 0xff, 0xa0, 0x0, 0x1, 0x3, 0x1, 0x9, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x1, 0xff, 0xa2, 0x0, 0x1, 0x12, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x1, 0xff, 0xa6, 0x0, 0x1, 0xb, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x44, 0x1, 0x4, 0x0, 0x0, 0x0, 0xff, 0xb8, 0xff, 0xa1, 0x3, 0x1, 0x1, 0x1a, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x64, 0x4f, 0x6b, 0x1, 0xff, 0xa2, 0x0, 0x1, 0xa, 0x1, 0xa, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xa, 0x41, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x72, 0x79, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x8, 0x42, 0x69, 0x72, 0x74, 0x68, 0x64, 0x61, 0x79, 0x1, 0xff, 0x8c, 0x0, 0x1, 0xb, 0x42, 0x6c, 0x6f, 0x6f, 0x64, 0x6c, 0x69, 0x6e, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xd, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xb, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xc, 0x0, 0x1, 0x6, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x1, 0xc, 0x0, 0x1, 0x4, 0x4e, 0x61, 0x6d, 0x65, 0x1, 0xc, 0x0, 0x1, 0x6, 0x52, 0x61, 0x63, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xe, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1, 0x8, 0x0, 0x0, 0x0, 0x10, 0xff, 0x8b, 0x5, 0x1, 0x1, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x8c, 0x0, 0x0, 0x0, 0x44, 0xff, 0xa5, 0x2, 0x1, 0x1, 0x35, 0x5b, 0x5d, 0x65, 0x73, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x64, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x32, 0x30, 0x30, 0x4f, 0x6b, 0x1, 0xff, 0xa6, 0x0, 0x1, 0xff, 0xa4, 0x0, 0x0, 0x79, 0xff, 0xa3, 0x3, 0x1, 0x1, 0x2f, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x64, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x32, 0x30, 0x30, 0x4f, 0x6b, 0x1, 0xff, 0xa4, 0x0, 0x1, 0x4, 0x1, 0xd, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x9, 0x49, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x1, 0x2, 0x0, 0x1, 0x8, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x9, 0x53, 0x74, 0x61, 0x72, 0x74, 0x44, 0x61, 0x74, 0x65, 0x1, 0xff, 0x8c, 0x0, 0x0, 0x0, 0x79, 0xff, 0xa0, 0x1, 0x2, 0x26, 0x1, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xcc, 0xa3, 0x43, 0x5c, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x1, 0x6, 0x1, 0xfc, 0xd, 0x7, 0x91, 0xec, 0x2, 0x4, 0x6d, 0x61, 0x6c, 0x65, 0x1, 0xd, 0x43, 0x43, 0x50, 0x20, 0x42, 0x61, 0x72, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x1, 0x4, 0x0, 0x1, 0x2, 0x1, 0xfc, 0xa, 0xba, 0x95, 0x2, 0x1, 0x1, 0x1, 0xfe, 0x3, 0xe8, 0x1, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xcf, 0x2, 0x2b, 0x40, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x0, 0x1, 0xfc, 0xa, 0xba, 0x95, 0x4, 0x2, 0xfe, 0x3, 0xea, 0x1, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xcf, 0x29, 0xb8, 0x40, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x0, 0x1, 0xfc, 0x15, 0x85, 0xe9, 0x98, 0x0},
		Timestamp: 1509306343221144957,
		Attempts:  0x1,
	}
	err := nailInstance.characterHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAlliance(t *testing.T) {
	m := &nsq.Message{
		ID:        nsq.MessageID{0x30, 0x38, 0x63, 0x33, 0x36, 0x35, 0x33, 0x63, 0x64, 0x35, 0x66, 0x38, 0x35, 0x30, 0x30, 0x30},
		Body:      []uint8{0x4d, 0xff, 0x9b, 0x3, 0x1, 0x1, 0x8, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x1, 0xff, 0x9c, 0x0, 0x1, 0x3, 0x1, 0x8, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x1, 0xff, 0x9e, 0x0, 0x1, 0x14, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1, 0xff, 0x9a, 0x0, 0x1, 0xa, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x44, 0x1, 0x4, 0x0, 0x0, 0x0, 0x64, 0xff, 0x9d, 0x3, 0x1, 0x1, 0x18, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x4f, 0x6b, 0x1, 0xff, 0x9e, 0x0, 0x1, 0x4, 0x1, 0xc, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x1, 0xc, 0x0, 0x1, 0xb, 0x44, 0x61, 0x74, 0x65, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x65, 0x64, 0x1, 0xff, 0x8c, 0x0, 0x1, 0xc, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x43, 0x6f, 0x72, 0x70, 0x1, 0x4, 0x0, 0x1, 0x6, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x1, 0xc, 0x0, 0x0, 0x0, 0x10, 0xff, 0x8b, 0x5, 0x1, 0x1, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x8c, 0x0, 0x0, 0x0, 0xc, 0xff, 0x99, 0x2, 0x1, 0x2, 0xff, 0x9a, 0x0, 0x1, 0x4, 0x0, 0x0, 0x3e, 0xff, 0x9c, 0x1, 0x1, 0xe, 0x43, 0x20, 0x43, 0x20, 0x50, 0x20, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x1, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xcf, 0x2, 0x39, 0x50, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x1, 0xfc, 0xb, 0xb9, 0x97, 0xc2, 0x1, 0x7, 0x3c, 0x43, 0x20, 0x43, 0x20, 0x50, 0x3e, 0x0, 0x1, 0x1, 0xfc, 0xb, 0xae, 0xb9, 0x2, 0x1, 0x2, 0x0},
		Timestamp: 1509306343179632645,
		Attempts:  0x1,
	}
	err := nailInstance.allianceHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCorporation(t *testing.T) {
	m := &nsq.Message{
		ID:        nsq.MessageID{0x30, 0x38, 0x63, 0x33, 0x63, 0x31, 0x34, 0x65, 0x61, 0x63, 0x37, 0x38, 0x35, 0x30, 0x30, 0x30},
		Body:      []uint8{0x3c, 0xff, 0xab, 0x3, 0x1, 0x1, 0xb, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0xac, 0x0, 0x1, 0x2, 0x1, 0xd, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x1, 0x4, 0x0, 0x1, 0xb, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0xae, 0x0, 0x0, 0x0, 0xff, 0xcf, 0xff, 0xad, 0x3, 0x1, 0x1, 0x1e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x4f, 0x6b, 0x1, 0xff, 0xae, 0x0, 0x1, 0xb, 0x1, 0xa, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x5, 0x43, 0x65, 0x6f, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x16, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xc, 0x0, 0x1, 0xf, 0x43, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x1, 0xc, 0x0, 0x1, 0xc, 0x43, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x65, 0x1, 0xff, 0x8c, 0x0, 0x1, 0x9, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x7, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xc, 0x0, 0x1, 0xb, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x1, 0x4, 0x0, 0x1, 0x7, 0x54, 0x61, 0x78, 0x52, 0x61, 0x74, 0x65, 0x1, 0x8, 0x0, 0x1, 0x6, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x1, 0xc, 0x0, 0x1, 0x3, 0x55, 0x72, 0x6c, 0x1, 0xc, 0x0, 0x0, 0x0, 0x10, 0xff, 0x8b, 0x5, 0x1, 0x1, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x8c, 0x0, 0x0, 0x0, 0xff, 0xa2, 0xff, 0xac, 0x1, 0xfc, 0xd, 0x7, 0x91, 0xec, 0x1, 0x1, 0xfc, 0x33, 0xc4, 0x11, 0x16, 0x1, 0xfc, 0x15, 0x85, 0xe9, 0x98, 0x1, 0x3f, 0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2c, 0x20, 0x69, 0x74, 0x27, 0x73, 0x20, 0x62, 0x61, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x6c, 0x79, 0x20, 0x6a, 0x75, 0x73, 0x74, 0x20, 0x61, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x1, 0x5, 0x43, 0x20, 0x43, 0x20, 0x50, 0x1, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xb9, 0x3b, 0xf7, 0xb, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x1, 0xfc, 0x15, 0x85, 0xe9, 0x98, 0x2, 0xfe, 0x5, 0x20, 0x1, 0xfb, 0xe0, 0x4d, 0x62, 0xd0, 0x3f, 0x1, 0x5, 0x2d, 0x43, 0x43, 0x50, 0x2d, 0x1, 0x18, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x77, 0x77, 0x77, 0x2e, 0x65, 0x76, 0x65, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x0, 0x0},
		Timestamp: 1509331651099999982,
		Attempts:  0x1,
	}
	err := nailInstance.corporationHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCharWaller(t *testing.T) {
	m := &nsq.Message{
		ID:        nsq.MessageID{0x30, 0x38, 0x64, 0x32, 0x36, 0x64, 0x33, 0x31, 0x31, 0x34, 0x33, 0x38, 0x35, 0x30, 0x30, 0x30},
		Body:      []uint8{0x60, 0xff, 0xcb, 0x3, 0x1, 0x1, 0x1b, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1, 0xff, 0xcc, 0x0, 0x1, 0x3, 0x1, 0xc, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1, 0xff, 0xd0, 0x0, 0x1, 0xb, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x44, 0x1, 0x4, 0x0, 0x1, 0x10, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x44, 0x1, 0x4, 0x0, 0x0, 0x0, 0x44, 0xff, 0xcf, 0x2, 0x1, 0x1, 0x35, 0x5b, 0x5d, 0x65, 0x73, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x64, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x32, 0x30, 0x30, 0x4f, 0x6b, 0x1, 0xff, 0xd0, 0x0, 0x1, 0xff, 0xce, 0x0, 0x0, 0xff, 0xc5, 0xff, 0xcd, 0x3, 0x1, 0x1, 0x2f, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x49, 0x64, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x32, 0x30, 0x30, 0x4f, 0x6b, 0x1, 0xff, 0xce, 0x0, 0x1, 0xa, 0x1, 0x8, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x4, 0x44, 0x61, 0x74, 0x65, 0x1, 0xff, 0x84, 0x0, 0x1, 0x5, 0x49, 0x73, 0x42, 0x75, 0x79, 0x1, 0x2, 0x0, 0x1, 0xa, 0x49, 0x73, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x1, 0x2, 0x0, 0x1, 0xc, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6c, 0x52, 0x65, 0x66, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0xa, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x8, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x1, 0x4, 0x0, 0x1, 0xd, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x6, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x1, 0x4, 0x0, 0x1, 0x9, 0x55, 0x6e, 0x69, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x1, 0x8, 0x0, 0x0, 0x0, 0x10, 0xff, 0x83, 0x5, 0x1, 0x1, 0x4, 0x54, 0x69, 0x6d, 0x65, 0x1, 0xff, 0x84, 0x0, 0x0, 0x0, 0x3f, 0xff, 0xcc, 0x1, 0x1, 0x1, 0xfd, 0x1, 0xa8, 0x62, 0x1, 0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0xcf, 0x9f, 0xc4, 0x90, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x1, 0x1, 0x1, 0x1, 0x1, 0xfd, 0x2, 0x12, 0x64, 0x1, 0xfc, 0x7, 0x27, 0x80, 0xfe, 0x1, 0x2, 0x1, 0xfc, 0x93, 0x2c, 0x5, 0xa4, 0x1, 0xfe, 0x4, 0x96, 0x1, 0xfe, 0xf0, 0x3f, 0x0, 0x1, 0x2, 0x1, 0x2, 0x0},
		Timestamp: 1510364060741743066,
		Attempts:  0x1,
	}
	err := nailInstance.characterWalletTransactionConsumer(m)
	if err != nil {
		t.Fatal(err)
	}
}
