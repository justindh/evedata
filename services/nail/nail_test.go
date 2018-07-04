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
		{Operation: "characterWalletTransactions", Parameter: []int32{int32(1), int32(1)}},
		{Operation: "characterWalletJournal", Parameter: []int32{int32(1), int32(1)}},
		{Operation: "characterAssets", Parameter: []int32{int32(1), int32(1)}},
		{Operation: "loyaltyStore", Parameter: int32(1000001)},
		{Operation: "characterNotifications", Parameter: []int32{int32(1), int32(1)}},
		{Operation: "characterAuthOwner", Parameter: []int32{int32(1), int32(1)}},
		{Operation: "charSearch", Parameter: "SomeDOtherude"},
	}
	ham          *hammer.Hammer
	nailInstance *Nail
)

func TestMain(m *testing.M) {
	sql := sqlhelper.NewTestDatabase()

	redis := redigohelper.ConnectRedisTestPool()
	ledis := redigohelper.ConnectLedisTestPool()
	redConn := redis.Get()
	defer redConn.Close()
	redConn.Do("FLUSHALL")

	producer, err := nsqhelper.NewTestNSQProducer()
	if err != nil {
		log.Fatalln(err)
	}

	ham = hammer.NewHammer(redis, ledis, sql, producer, "sofake", "faaaaaaake", "sofake")
	ham.ChangeBasePath("http://127.0.0.1:8080")
	ham.ChangeTokenPath("http://127.0.0.1:8080")
	ham.SetToken(1, 1, &oauth2.Token{
		RefreshToken: "fake",
		AccessToken:  "really fake",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	})

	nailInstance = NewNail(redis, sql, nsqhelper.Test)

	go ham.Run()
	go nailInstance.Run()
	retCode := m.Run()
	time.Sleep(time.Second * 5)
	nailInstance.Close()
	ham.Close()
	redis.Close()
	sql.Close()

	time.Sleep(time.Second)

	os.Exit(retCode)
}

func TestQueue(t *testing.T) {
	err := ham.QueueWork(testWork, redisqueue.Priority_Low)
	if err != nil {
		t.Fatal(err)
	}
}

func TestKillmails(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x35, 0x64, 0x33, 0x38, 0x62, 0x36, 0x61, 0x33, 0x34, 0x31, 0x30, 0x30, 0x30, 0x31}, Body: []uint8{0x42, 0x2, 0x0, 0x0, 0x2, 0x68, 0x61, 0x73, 0x68, 0x0, 0x9, 0x0, 0x0, 0x0, 0x46, 0x41, 0x4b, 0x45, 0x48, 0x41, 0x53, 0x48, 0x0, 0x3, 0x6b, 0x69, 0x6c, 0x6c, 0x0, 0x24, 0x2, 0x0, 0x0, 0x10, 0x6b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x69, 0x64, 0x0, 0x7d, 0xb0, 0x61, 0x3, 0x9, 0x6b, 0x69, 0x6c, 0x6c, 0x6d, 0x61, 0x69, 0x6c, 0x74, 0x69, 0x6d, 0x65, 0x0, 0x0, 0x2a, 0x62, 0xed, 0x57, 0x1, 0x0, 0x0, 0x3, 0x76, 0x69, 0x63, 0x74, 0x69, 0x6d, 0x0, 0x16, 0x1, 0x0, 0x0, 0x10, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x51, 0xf5, 0x87, 0x5, 0x10, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0xd7, 0x30, 0x26, 0x32, 0x10, 0x61, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x69, 0x64, 0x0, 0xba, 0xdf, 0x8, 0x25, 0x10, 0x66, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x74, 0x61, 0x6b, 0x65, 0x6e, 0x0, 0x71, 0x16, 0x0, 0x0, 0x10, 0x73, 0x68, 0x69, 0x70, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0x94, 0x45, 0x0, 0x0, 0x4, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x0, 0x76, 0x0, 0x0, 0x0, 0x3, 0x30, 0x0, 0x6e, 0x0, 0x0, 0x0, 0x10, 0x69, 0x74, 0x65, 0x6d, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0x55, 0x17, 0x0, 0x0, 0x12, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x64, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x12, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x74, 0x6f, 0x6e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x66, 0x6c, 0x61, 0x67, 0x0, 0x14, 0x0, 0x0, 0x0, 0x4, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x26, 0x0, 0x0, 0x0, 0x1, 0x78, 0x0, 0x63, 0x5e, 0x1e, 0xfd, 0x1a, 0x52, 0x5a, 0x42, 0x1, 0x79, 0x0, 0x7c, 0x73, 0xe9, 0x7, 0x26, 0x14, 0x41, 0x42, 0x1, 0x7a, 0x0, 0x76, 0x8b, 0xb4, 0x20, 0x94, 0x7f, 0x39, 0x42, 0x0, 0x0, 0x4, 0x61, 0x74, 0x74, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x73, 0x0, 0xa6, 0x0, 0x0, 0x0, 0x3, 0x30, 0x0, 0x9e, 0x0, 0x0, 0x0, 0x10, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x80, 0xf5, 0xb5, 0x5, 0x10, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0xf3, 0x42, 0xf, 0x0, 0x10, 0x61, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x66, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x23, 0xa1, 0x7, 0x0, 0x1, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x0, 0x0, 0x0, 0x0, 0x40, 0x33, 0x33, 0xd3, 0xbf, 0x8, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x62, 0x6c, 0x6f, 0x77, 0x0, 0x1, 0x10, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x64, 0x6f, 0x6e, 0x65, 0x0, 0x71, 0x16, 0x0, 0x0, 0x10, 0x73, 0x68, 0x69, 0x70, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0xb1, 0x45, 0x0, 0x0, 0x10, 0x77, 0x65, 0x61, 0x70, 0x6f, 0x6e, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0x2, 0xc, 0x0, 0x0, 0x0, 0x0, 0x10, 0x73, 0x6f, 0x6c, 0x61, 0x72, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x69, 0x64, 0x0, 0x20, 0xcf, 0xc9, 0x1, 0x10, 0x6d, 0x6f, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x77, 0x61, 0x72, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1520130891061767331}
	err := nailInstance.killmailHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMarket(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x62, 0x63, 0x63, 0x38, 0x37, 0x62, 0x38, 0x31, 0x30, 0x30, 0x30, 0x31}, Body: []uint8{0xd8, 0x0, 0x0, 0x0, 0x4, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x0, 0xbd, 0x0, 0x0, 0x0, 0x3, 0x30, 0x0, 0xb5, 0x0, 0x0, 0x0, 0x12, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x69, 0x64, 0x0, 0x5f, 0xf5, 0x99, 0x13, 0x1, 0x0, 0x0, 0x0, 0x10, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0x22, 0x0, 0x0, 0x0, 0x12, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0xdf, 0x9c, 0x93, 0x3, 0x0, 0x0, 0x0, 0x0, 0x10, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x0, 0x80, 0x84, 0x1e, 0x0, 0x10, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x0, 0x80, 0xc6, 0x13, 0x0, 0x10, 0x6d, 0x69, 0x6e, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x70, 0x72, 0x69, 0x63, 0x65, 0x0, 0xcd, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0x23, 0x40, 0x8, 0x69, 0x73, 0x62, 0x75, 0x79, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x0, 0x0, 0x10, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x5a, 0x0, 0x0, 0x0, 0x9, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x0, 0xa8, 0x4a, 0x76, 0xee, 0x56, 0x1, 0x0, 0x0, 0x2, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x0, 0x7, 0x0, 0x0, 0x0, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x0, 0x0, 0x0, 0x10, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513904375874297977, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.marketOrderHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStructure(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x62, 0x65, 0x34, 0x34, 0x35, 0x32, 0x38, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0x9b, 0x0, 0x0, 0x0, 0x3, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x0, 0x76, 0x0, 0x0, 0x0, 0x2, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0x18, 0x0, 0x0, 0x0, 0x56, 0x2d, 0x33, 0x59, 0x47, 0x37, 0x20, 0x56, 0x49, 0x20, 0x2d, 0x20, 0x54, 0x68, 0x65, 0x20, 0x43, 0x61, 0x70, 0x69, 0x74, 0x61, 0x6c, 0x0, 0x10, 0x73, 0x6f, 0x6c, 0x61, 0x72, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x69, 0x64, 0x0, 0xe, 0xc4, 0xc9, 0x1, 0x10, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x26, 0x0, 0x0, 0x0, 0x1, 0x78, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x79, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x7a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x12, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x69, 0x64, 0x0, 0x75, 0x52, 0xa5, 0xd4, 0xe8, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513904783724688237, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.structureHandler(m)
	if err != nil {
		t.Fatal(err)
	}

	m = &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x62, 0x65, 0x30, 0x65, 0x32, 0x66, 0x34, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0x9b, 0x0, 0x0, 0x0, 0x3, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x0, 0x76, 0x0, 0x0, 0x0, 0x2, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0x18, 0x0, 0x0, 0x0, 0x56, 0x2d, 0x33, 0x59, 0x47, 0x37, 0x20, 0x56, 0x49, 0x20, 0x2d, 0x20, 0x54, 0x68, 0x65, 0x20, 0x43, 0x61, 0x70, 0x69, 0x74, 0x61, 0x6c, 0x0, 0x10, 0x73, 0x6f, 0x6c, 0x61, 0x72, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x69, 0x64, 0x0, 0xe, 0xc4, 0xc9, 0x1, 0x10, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x26, 0x0, 0x0, 0x0, 0x1, 0x78, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x79, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x7a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x12, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x69, 0x64, 0x0, 0xe6, 0x71, 0xa5, 0xd4, 0xe8, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513904725594086304, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err = nailInstance.structureHandler(m)
	if err != nil {
		t.Fatal(err)
	}
	m = &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x62, 0x65, 0x36, 0x65, 0x66, 0x34, 0x63, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0xdf, 0x0, 0x0, 0x0, 0x4, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x0, 0xbd, 0x0, 0x0, 0x0, 0x3, 0x30, 0x0, 0xb5, 0x0, 0x0, 0x0, 0x12, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x69, 0x64, 0x0, 0x5f, 0xf5, 0x99, 0x13, 0x1, 0x0, 0x0, 0x0, 0x10, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0x22, 0x0, 0x0, 0x0, 0x12, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x28, 0x5b, 0xa6, 0xb7, 0xed, 0x0, 0x0, 0x0, 0x10, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x0, 0x80, 0x84, 0x1e, 0x0, 0x10, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x0, 0x80, 0xc6, 0x13, 0x0, 0x10, 0x6d, 0x69, 0x6e, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x70, 0x72, 0x69, 0x63, 0x65, 0x0, 0xcd, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0x23, 0x40, 0x8, 0x69, 0x73, 0x62, 0x75, 0x79, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x0, 0x0, 0x10, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x5a, 0x0, 0x0, 0x0, 0x9, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x0, 0xa8, 0x4a, 0x76, 0xee, 0x56, 0x1, 0x0, 0x0, 0x2, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x0, 0x7, 0x0, 0x0, 0x0, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x0, 0x0, 0x0, 0x12, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x69, 0x64, 0x0, 0x75, 0x52, 0xa5, 0xd4, 0xe8, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513904829501863472, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	// Hack to use an actual structureID.
	b := datapackages.StructureOrders{}
	gobcoder.GobDecoder(m.Body, &b)
	for i := range b.Orders {
		b.Orders[i].LocationId = 1000000017013
	}
	b.StructureID = 1000000017013
	newb, _ := gobcoder.GobEncoder(b)
	m.Body = newb

	err = nailInstance.structureMarketHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMarketHistory(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x63, 0x38, 0x65, 0x32, 0x65, 0x35, 0x38, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0xa0, 0x0, 0x0, 0x0, 0x4, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x0, 0x78, 0x0, 0x0, 0x0, 0x3, 0x30, 0x0, 0x70, 0x0, 0x0, 0x0, 0x2, 0x64, 0x61, 0x74, 0x65, 0x0, 0xb, 0x0, 0x0, 0x0, 0x32, 0x30, 0x31, 0x35, 0x2d, 0x30, 0x35, 0x2d, 0x30, 0x31, 0x0, 0x12, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x0, 0xdb, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x12, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x0, 0xd3, 0xfb, 0x2b, 0xca, 0x3, 0x0, 0x0, 0x0, 0x1, 0x68, 0x69, 0x67, 0x68, 0x65, 0x73, 0x74, 0x0, 0x14, 0xae, 0x47, 0xe1, 0x7a, 0x14, 0x15, 0x40, 0x1, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x15, 0x40, 0x1, 0x6c, 0x6f, 0x77, 0x65, 0x73, 0x74, 0x0, 0x71, 0x3d, 0xa, 0xd7, 0xa3, 0x70, 0x14, 0x40, 0x0, 0x0, 0x10, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x10, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513907702771383981, Attempts: 0x1, NSQDAddress: "localhost:4150"}

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
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x63, 0x61, 0x61, 0x33, 0x30, 0x65, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0x2e, 0x1, 0x0, 0x0, 0x10, 0x69, 0x64, 0x0, 0x95, 0x7, 0x0, 0x0, 0x9, 0x64, 0x65, 0x63, 0x6c, 0x61, 0x72, 0x65, 0x64, 0x0, 0x0, 0xb8, 0x26, 0xab, 0xfc, 0x0, 0x0, 0x0, 0x9, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x0, 0x0, 0x28, 0xd3, 0xed, 0x7c, 0xc7, 0xff, 0xff, 0x9, 0x72, 0x65, 0x74, 0x72, 0x61, 0x63, 0x74, 0x65, 0x64, 0x0, 0x0, 0x28, 0xd3, 0xed, 0x7c, 0xc7, 0xff, 0xff, 0x9, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x0, 0x0, 0x28, 0xd3, 0xed, 0x7c, 0xc7, 0xff, 0xff, 0x8, 0x6d, 0x75, 0x74, 0x75, 0x61, 0x6c, 0x0, 0x0, 0x8, 0x6f, 0x70, 0x65, 0x6e, 0x66, 0x6f, 0x72, 0x61, 0x6c, 0x6c, 0x69, 0x65, 0x73, 0x0, 0x0, 0x3, 0x61, 0x67, 0x67, 0x72, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x0, 0x4f, 0x0, 0x0, 0x0, 0x10, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x40, 0x53, 0xcf, 0x3a, 0x10, 0x61, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x73, 0x68, 0x69, 0x70, 0x73, 0x6b, 0x69, 0x6c, 0x6c, 0x65, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x69, 0x73, 0x6b, 0x64, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x64, 0x65, 0x66, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x0, 0x4f, 0x0, 0x0, 0x0, 0x10, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x9b, 0x9f, 0xb2, 0x3b, 0x10, 0x61, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x73, 0x68, 0x69, 0x70, 0x73, 0x6b, 0x69, 0x6c, 0x6c, 0x65, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x69, 0x73, 0x6b, 0x64, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x61, 0x6c, 0x6c, 0x69, 0x65, 0x73, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513908183977700436, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.warHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCharacter(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x63, 0x62, 0x39, 0x36, 0x34, 0x32, 0x38, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0x93, 0x1, 0x0, 0x0, 0x3, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x0, 0xc9, 0x0, 0x0, 0x0, 0x2, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0xe, 0x0, 0x0, 0x0, 0x43, 0x43, 0x50, 0x20, 0x42, 0x61, 0x72, 0x74, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x0, 0x2, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x10, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0xf6, 0xc8, 0x83, 0x6, 0x10, 0x61, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9, 0x62, 0x69, 0x72, 0x74, 0x68, 0x64, 0x61, 0x79, 0x0, 0x60, 0x47, 0x92, 0x4b, 0x4c, 0x1, 0x0, 0x0, 0x2, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x0, 0x5, 0x0, 0x0, 0x0, 0x6d, 0x61, 0x6c, 0x65, 0x0, 0x10, 0x72, 0x61, 0x63, 0x65, 0x69, 0x64, 0x0, 0x2, 0x0, 0x0, 0x0, 0x10, 0x62, 0x6c, 0x6f, 0x6f, 0x64, 0x6c, 0x69, 0x6e, 0x65, 0x69, 0x64, 0x0, 0x3, 0x0, 0x0, 0x0, 0x10, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x72, 0x79, 0x69, 0x64, 0x0, 0x13, 0x0, 0x0, 0x0, 0x1, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x66, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x0, 0x95, 0x0, 0x0, 0x0, 0x3, 0x30, 0x0, 0x45, 0x0, 0x0, 0x0, 0x9, 0x73, 0x74, 0x61, 0x72, 0x74, 0x64, 0x61, 0x74, 0x65, 0x0, 0x0, 0x1a, 0x4c, 0x8e, 0x55, 0x1, 0x0, 0x0, 0x10, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x81, 0x4a, 0x5d, 0x5, 0x8, 0x69, 0x73, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x0, 0x1, 0x10, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x64, 0x0, 0xf4, 0x1, 0x0, 0x0, 0x0, 0x3, 0x31, 0x0, 0x45, 0x0, 0x0, 0x0, 0x9, 0x73, 0x74, 0x61, 0x72, 0x74, 0x64, 0x61, 0x74, 0x65, 0x0, 0x0, 0xe2, 0xca, 0x28, 0x56, 0x1, 0x0, 0x0, 0x10, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x82, 0x4a, 0x5d, 0x5, 0x8, 0x69, 0x73, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x0, 0x0, 0x10, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x64, 0x0, 0xf5, 0x1, 0x0, 0x0, 0x0, 0x0, 0x10, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x51, 0xf5, 0x87, 0x5, 0x0}, Timestamp: 1513908445117343967, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.characterHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAlliance(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x63, 0x62, 0x63, 0x31, 0x38, 0x32, 0x38, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0xa2, 0x0, 0x0, 0x0, 0x3, 0x61, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x0, 0x61, 0x0, 0x0, 0x0, 0x2, 0x61, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0xf, 0x0, 0x0, 0x0, 0x43, 0x20, 0x43, 0x20, 0x50, 0x20, 0x41, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x0, 0x2, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x0, 0x8, 0x0, 0x0, 0x0, 0x3c, 0x43, 0x20, 0x43, 0x20, 0x50, 0x3e, 0x0, 0x10, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x63, 0x6f, 0x72, 0x70, 0x0, 0xe1, 0xcb, 0xdc, 0x5, 0x9, 0x64, 0x61, 0x74, 0x65, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x65, 0x64, 0x0, 0x80, 0x8, 0x83, 0x8e, 0x55, 0x1, 0x0, 0x0, 0x0, 0x4, 0x61, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x0, 0xc, 0x0, 0x0, 0x0, 0x10, 0x30, 0x0, 0x81, 0x5c, 0xd7, 0x5, 0x0, 0x10, 0x61, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513908491556595870, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.allianceHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCorporation(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x63, 0x62, 0x65, 0x38, 0x65, 0x39, 0x38, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0x45, 0x1, 0x0, 0x0, 0x10, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x81, 0x5c, 0xd7, 0x5, 0x3, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x20, 0x1, 0x0, 0x0, 0x2, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0x6, 0x0, 0x0, 0x0, 0x43, 0x20, 0x43, 0x20, 0x50, 0x0, 0x2, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x0, 0x6, 0x0, 0x0, 0x0, 0x2d, 0x43, 0x43, 0x50, 0x2d, 0x0, 0x10, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x0, 0x90, 0x2, 0x0, 0x0, 0x10, 0x63, 0x65, 0x6f, 0x69, 0x64, 0x0, 0xcc, 0xf4, 0xc2, 0xa, 0x10, 0x61, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x69, 0x64, 0x0, 0x8b, 0x8, 0xe2, 0x19, 0x2, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x40, 0x0, 0x0, 0x0, 0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2c, 0x20, 0x69, 0x74, 0x27, 0x73, 0x20, 0x62, 0x61, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x6c, 0x79, 0x20, 0x6a, 0x75, 0x73, 0x74, 0x20, 0x61, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x0, 0x1, 0x74, 0x61, 0x78, 0x72, 0x61, 0x74, 0x65, 0x0, 0x0, 0x0, 0x0, 0xe0, 0x4d, 0x62, 0xd0, 0x3f, 0x9, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x64, 0x61, 0x74, 0x65, 0x0, 0xf8, 0x2a, 0x10, 0x80, 0x0, 0x1, 0x0, 0x0, 0x10, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x69, 0x64, 0x0, 0xcc, 0xf4, 0xc2, 0xa, 0x2, 0x75, 0x72, 0x6c, 0x0, 0x19, 0x0, 0x0, 0x0, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x77, 0x77, 0x77, 0x2e, 0x65, 0x76, 0x65, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x0, 0x2, 0x66, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513908533864555539, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.corporationHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCharWallet(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x63, 0x64, 0x34, 0x39, 0x65, 0x30, 0x34, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0xe6, 0x0, 0x0, 0x0, 0x4, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x0, 0xac, 0x0, 0x0, 0x0, 0x3, 0x30, 0x0, 0xa4, 0x0, 0x0, 0x0, 0x12, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0xd2, 0x2, 0x96, 0x49, 0x0, 0x0, 0x0, 0x0, 0x9, 0x64, 0x61, 0x74, 0x65, 0x0, 0x80, 0xfa, 0xea, 0xf5, 0x57, 0x1, 0x0, 0x0, 0x10, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0x4b, 0x2, 0x0, 0x0, 0x12, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x7f, 0xc0, 0x93, 0x3, 0x0, 0x0, 0x0, 0x0, 0x1, 0x75, 0x6e, 0x69, 0x74, 0x70, 0x72, 0x69, 0x63, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf0, 0x3f, 0x10, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x0, 0x1, 0x0, 0x0, 0x0, 0x10, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x69, 0x64, 0x0, 0x31, 0xd4, 0x0, 0x0, 0x8, 0x69, 0x73, 0x62, 0x75, 0x79, 0x0, 0x1, 0x8, 0x69, 0x73, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x0, 0x1, 0x12, 0x6a, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6c, 0x72, 0x65, 0x66, 0x69, 0x64, 0x0, 0x32, 0x9, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x10, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513908912856541045, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.characterWalletTransactionConsumer(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCharAssets(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x63, 0x64, 0x36, 0x65, 0x38, 0x36, 0x34, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0xc0, 0x0, 0x0, 0x0, 0x4, 0x61, 0x73, 0x73, 0x65, 0x74, 0x73, 0x0, 0x8c, 0x0, 0x0, 0x0, 0x3, 0x30, 0x0, 0x84, 0x0, 0x0, 0x0, 0x10, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0xbc, 0xd, 0x0, 0x0, 0x10, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x0, 0x0, 0x0, 0x0, 0x0, 0x12, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x8f, 0x92, 0x93, 0x3, 0x0, 0x0, 0x0, 0x0, 0x2, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x74, 0x79, 0x70, 0x65, 0x0, 0x8, 0x0, 0x0, 0x0, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x12, 0x69, 0x74, 0x65, 0x6d, 0x69, 0x64, 0x0, 0xc3, 0x51, 0xa5, 0xd4, 0xe8, 0x0, 0x0, 0x0, 0x2, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x66, 0x6c, 0x61, 0x67, 0x0, 0x7, 0x0, 0x0, 0x0, 0x48, 0x61, 0x6e, 0x67, 0x61, 0x72, 0x0, 0x8, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x74, 0x6f, 0x6e, 0x0, 0x1, 0x0, 0x0, 0x10, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x10, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513908952207343919, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.characterAssetsConsumer(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoyaltyStore(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x63, 0x63, 0x32, 0x31, 0x63, 0x61, 0x38, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0x6, 0x1, 0x0, 0x0, 0x10, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x41, 0x42, 0xf, 0x0, 0x4, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x0, 0xe7, 0x0, 0x0, 0x0, 0x3, 0x30, 0x0, 0x5d, 0x0, 0x0, 0x0, 0x10, 0x6f, 0x66, 0x66, 0x65, 0x72, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x10, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x10, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x0, 0x1, 0x0, 0x0, 0x0, 0x10, 0x6c, 0x70, 0x63, 0x6f, 0x73, 0x74, 0x0, 0x64, 0x0, 0x0, 0x0, 0x1, 0x69, 0x73, 0x6b, 0x63, 0x6f, 0x73, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x31, 0x0, 0x7f, 0x0, 0x0, 0x0, 0x10, 0x6f, 0x66, 0x66, 0x65, 0x72, 0x69, 0x64, 0x0, 0x2, 0x0, 0x0, 0x0, 0x10, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0xd3, 0x4, 0x0, 0x0, 0x10, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x0, 0xa, 0x0, 0x0, 0x0, 0x10, 0x6c, 0x70, 0x63, 0x6f, 0x73, 0x74, 0x0, 0x64, 0x0, 0x0, 0x0, 0x1, 0x69, 0x73, 0x6b, 0x63, 0x6f, 0x73, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x8f, 0x40, 0x4, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x0, 0x27, 0x0, 0x0, 0x0, 0x3, 0x30, 0x0, 0x1f, 0x0, 0x0, 0x0, 0x10, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0xd2, 0x4, 0x0, 0x0, 0x10, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x0, 0xa, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513908594937510512, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.loyaltyStoreHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotifications(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x63, 0x64, 0x62, 0x31, 0x34, 0x63, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0xd, 0x1, 0x0, 0x0, 0x4, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x0, 0xd2, 0x0, 0x0, 0x0, 0x3, 0x30, 0x0, 0xca, 0x0, 0x0, 0x0, 0x12, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x0, 0x13, 0x0, 0x0, 0x0, 0x49, 0x6e, 0x73, 0x75, 0x72, 0x61, 0x6e, 0x63, 0x65, 0x50, 0x61, 0x79, 0x6f, 0x75, 0x74, 0x4d, 0x73, 0x67, 0x0, 0x10, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x69, 0x64, 0x0, 0xc4, 0x42, 0xf, 0x0, 0x2, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x74, 0x79, 0x70, 0x65, 0x0, 0xc, 0x0, 0x0, 0x0, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x0, 0x9, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x0, 0x0, 0x9c, 0x83, 0xea, 0x5d, 0x1, 0x0, 0x0, 0x8, 0x69, 0x73, 0x72, 0x65, 0x61, 0x64, 0x0, 0x1, 0x2, 0x74, 0x65, 0x78, 0x74, 0x0, 0x3f, 0x0, 0x0, 0x0, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x3a, 0x20, 0x33, 0x37, 0x33, 0x31, 0x30, 0x31, 0x36, 0x2e, 0x34, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x34, 0x5c, 0x6e, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x44, 0x3a, 0x20, 0x31, 0x30, 0x32, 0x34, 0x38, 0x38, 0x31, 0x30, 0x32, 0x31, 0x36, 0x36, 0x33, 0x5c, 0x6e, 0x70, 0x61, 0x79, 0x6f, 0x75, 0x74, 0x3a, 0x20, 0x31, 0x5c, 0x6e, 0x0, 0x0, 0x0, 0x10, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x10, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513909023903629650, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.characterNotificationsHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCharJournal(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x30, 0x34, 0x64, 0x36, 0x37, 0x36, 0x38, 0x30, 0x38, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0xf1, 0x1, 0x0, 0x0, 0x4, 0x6a, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6c, 0x0, 0xbc, 0x1, 0x0, 0x0, 0x3, 0x30, 0x0, 0xb4, 0x1, 0x0, 0x0, 0x9, 0x64, 0x61, 0x74, 0x65, 0x0, 0x80, 0xfa, 0xea, 0xf5, 0x57, 0x1, 0x0, 0x0, 0x12, 0x72, 0x65, 0x66, 0x69, 0x64, 0x0, 0xd2, 0x2, 0x96, 0x49, 0x0, 0x0, 0x0, 0x0, 0x2, 0x72, 0x65, 0x66, 0x74, 0x79, 0x70, 0x65, 0x0, 0xf, 0x0, 0x0, 0x0, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x0, 0x10, 0x66, 0x69, 0x72, 0x73, 0x74, 0x70, 0x61, 0x72, 0x74, 0x79, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x66, 0x69, 0x72, 0x73, 0x74, 0x70, 0x61, 0x72, 0x74, 0x79, 0x74, 0x79, 0x70, 0x65, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x10, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x70, 0x61, 0x72, 0x74, 0x79, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x70, 0x61, 0x72, 0x74, 0x79, 0x74, 0x79, 0x70, 0x65, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x1, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x10, 0x74, 0x61, 0x78, 0x72, 0x65, 0x63, 0x69, 0x65, 0x76, 0x65, 0x72, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x74, 0x61, 0x78, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x65, 0x78, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x66, 0x6f, 0x0, 0xcd, 0x0, 0x0, 0x0, 0x12, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x12, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x6e, 0x70, 0x63, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x10, 0x6e, 0x70, 0x63, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x64, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x65, 0x64, 0x73, 0x68, 0x69, 0x70, 0x74, 0x79, 0x70, 0x65, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x63, 0x6f, 0x72, 0x70, 0x6f, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x61, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x63, 0x65, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x6a, 0x6f, 0x62, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x74, 0x69, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x10, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1513911434674320784, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.characterWalletJournalConsumer(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthOwner(t *testing.T) {
	m := &nsq.Message{ID: nsq.MessageID{0x30, 0x39, 0x34, 0x39, 0x38, 0x61, 0x32, 0x65, 0x32, 0x61, 0x63, 0x31, 0x30, 0x30, 0x30, 0x30}, Body: []uint8{0xa0, 0x0, 0x0, 0x0, 0x3, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x0, 0x6d, 0x0, 0x0, 0x0, 0x4, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x0, 0x2c, 0x0, 0x0, 0x0, 0x2, 0x30, 0x0, 0x9, 0x0, 0x0, 0x0, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x0, 0x2, 0x31, 0x0, 0x10, 0x0, 0x0, 0x0, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x0, 0x0, 0x4, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x61, 0x74, 0x68, 0x71, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x4, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x61, 0x74, 0x62, 0x61, 0x73, 0x65, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x4, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x61, 0x74, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x10, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x64, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0}, Timestamp: 1518745909632219803, Attempts: 0x1, NSQDAddress: "localhost:4150"}
	err := nailInstance.characterAuthOwnerHandler(m)
	if err != nil {
		t.Fatal(err)
	}
}
