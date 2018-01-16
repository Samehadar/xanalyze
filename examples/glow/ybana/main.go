package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	// "github.com/sniperkit/xanalyze/model"
	jsoniter "github.com/sniperkit/xutil/plugin/format/json"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/chrislusf/glow/flow"
	"github.com/comail/colog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/nats-io/go-nats"
	"github.com/seiflotfy/cuckoofilter"
	"github.com/willf/bloom"
)

//eg: {"name":"admin-on-rest","owner":"marmelab","path":"src/mui/detail/Tab.js","remote_id":"63226588"}
type File struct {
	name      string
	owner     string
	path      string
	remote_id int
}

type Manifest struct {
	manager string
	path    string
}

type AnaContext struct {
	fctx    *flow.FlowContext
	mbcli   *MsgBusClient
	dbh     *gorm.DB
	ana     *Analizer
	blmflt  *bloom.BloomFilter
	cuckflt *cuckoofilter.CuckooFilter
	fch     chan *EventR
}

func NewAnaContext() *AnaContext {
	this := &AnaContext{}

	this.fctx = flow.New()
	this.blmflt = bloom.New(100000, 5)
	this.cuckflt = cuckoofilter.NewDefaultCuckooFilter()
	this.fch = make(chan *EventR, 0)
	this.ana = NewAnalizer()

	return this
}

var (
	f    = flow.New()
	actx = NewAnaContext()
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Analizer struct {
	fctx *flow.FlowContext
}

func NewAnalizer() *Analizer {
	this := &Analizer{}
	this.fctx = flow.New()
	return this
}

func (this *Analizer) run() {
	fds := this.fctx.Channel(actx.fch)
	fds.Filter(func(line *EventR) bool {
		log.Println(line)
		return true
	}).Map(func(line *EventR, ch chan rune) {
		log.Println()
		for _, r := range line.Message {
			ch <- r
		}
	}).Map(func(r rune) (int, int8) {
		return 1, 2
	}).Reduce(func(x flow.KeyValue, y flow.KeyValue) flow.KeyValue {
		log.Println(x, y)
		return flow.KeyValue{Key: 5, Value: 6}
	})

	fds.Run()
}

// publish all messages to nats message bus
type MsgBusClient struct {
	nc *nats.Conn
}

func newMsgBusClient() *MsgBusClient {
	this := &MsgBusClient{}
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Println(err)
	}

	this.nc = nc
	return this
}

var blmflt = bloom.New(100000, 5)
var cuckflt = cuckoofilter.NewDefaultCuckooFilter()

var dbh *gorm.DB

// Msg is a structure used by Subscribers and PublishMsg().
type Msg struct {
	Subject string
	Reply   string
	Data    []byte
	// Sub     *Subscription
	next *Msg
}

type EventR struct {
	// gorm.Model
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Proto   string
	Nick    string
	Ident   string
	EType   string
	Message string
}

func messageHandler(m *nats.Msg) {
	// fmt.Printf("Received a message: %s\n", string(m.Data))
	jso, err := simplejson.NewJson(m.Data)
	if err != nil {
		log.Println(err)
		return
	}
	// TODO convert back to Event?

	switch jso.Get("Proto").MustString() {
	case "table":
		return
	}
	bfok := blmflt.TestAndAdd(m.Data)
	// log.Println(blmflt.TestAndAdd(m.Data))
	cfok := cuckflt.InsertUnique(m.Data)
	// log.Println(cuckflt.InsertUnique(m.Data))
	if bfok != !cfok {
		log.Println(bfok, cfok) // The result of the filter is inconsistent
	}
	if bfok == true || cfok == false {
		return // filtered
	}

	log.Printf("Received a message: %s\n", string(m.Data))
	if false { // use too much memory, about 1G
		nlp := NewSnowNLP(string(m.Data))
		log.Println(nlp.Sentiments(), nlp.Words())
	}

	evtrec := &EventR{}
	evtrec.EType = jso.Get("EType").MustString()
	evtrec.Proto = jso.Get("Proto").MustString()
	evtrec.Ident = jso.Get("Ident").MustString()

	switch jso.Get("Proto").MustString() {
	case "tox":
		evtrec.Nick = ""
		evtrec.Message = jso.Get("Args").GetIndex(0).MustString()
	case "irc":
		evtrec.Nick = jso.Get("Args").GetIndex(0).MustString()
		evtrec.Message = jso.Get("Args").GetIndex(2).MustString()
	}

	actx.fch <- evtrec

	dbh.Create(evtrec)
	if dbh.Error != nil {
		log.Println(dbh.Error)
	}
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func main() {
	db, err := gorm.Open("sqlite3", "ybana.db")
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v\n", db.DB().Driver())
	dbh = db
	actx.dbh = db
	defer db.Close()

	db.AutoMigrate(&EventR{})
	go actx.ana.run()

	/*
		files := f.TextFile(
			"files.json", 4,
		).Filter(func(line string) bool {
			return strings.Contains(line, "\"path\"")
		}).Map(func(line string, ch chan File) {
			var f File
			if strings.HasPrefix(line, ",") {
				line = after(line, ",")
			}
			err := json.Unmarshal([]byte(line), &f)
			if err != nil {
				fmt.Printf("parse source post error %v: %s\n", err, line)
				return
			}
			ch <- f
		}).Filter(func(f File) bool {
			return f.path != ""
		}).Map(func(f File) (m Manifest) {
			m.path = f.path

			// check for kewyords
			return
		})
	*/

	// f.Run()

	/*
		fileList := make([]model.TreeEntry, 1000)

		var data model.Tree
		err := json.Unmarshal([]byte(gitTree), &data)
		if err != nil {
			panic(err)
		}

		fmt.Println(data)

		for k, v := range data.Entries {
			fileList[k] = v
		}

		input := make(chan interface{})
		go func() {
			for _, data := range fileList {
				input <- data
			}
			close(input)
		}()

		// ret = context.Channel(input)
	*/

	/*
		 mbc := newMsgBusClient()
		sc, err := mbc.nc.Subscribe("yobotmsg", messageHandler)
		if err != nil {
			log.Println(err, sc.Type())
		}
	*/
	select {}
}

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	colog.Register()
}
