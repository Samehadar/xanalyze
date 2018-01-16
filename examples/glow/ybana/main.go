package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	//"github.com/sniperkit/xanalyze/model"
	"github.com/sniperkit/xfilter/backend/goac"
	"github.com/sniperkit/xgraph/plugin/cayley/fact"
	// jsoniter "github.com/sniperkit/xutil/plugin/format/json"
	// simplejson "github.com/bitly/go-simplejson"
	// "github.com/k0kubun/pp"

	"github.com/chrislusf/glow/flow"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jamiealquiza/tachymeter"
	"github.com/manveru/faker"
	"github.com/nats-io/go-nats"
	"github.com/seiflotfy/cuckoofilter"
	"github.com/willf/bloom"
)

var (
	f    = flow.New()
	actx = NewAnaContext()
	// json    = jsoniter.ConfigCompatibleWithStandardLibrary
	fake    *faker.Faker
	blmflt  = bloom.New(100000, 5)
	cuckflt = cuckoofilter.NewDefaultCuckooFilter()
	dbh     *gorm.DB
	ac      *goac.AhoCorasick            = goac.NewAhoCorasick()
	dicts   map[string]*goac.AhoCorasick = make(map[string]*goac.AhoCorasick, 0)
	t                                    = tachymeter.New(&tachymeter.Config{Size: 600000})
)

// {"name":"admin-on-rest","owner":"marmelab","path":"src/mui/detail/Tab.js","remote_id":"63226588"}
type File struct {
	name      string `json:"name"`
	owner     string `json:"owner"`
	path      string `json:"path"`
	remote_id int    `json:"remote_id"`
}

type Manifest struct {
	manager string
	path    string
}

type ScanResult struct {
	match string
	tags  string
	group string
	start int
	end   int
}

type AnaContext struct {
	fctx    *flow.FlowContext
	mbcli   *MsgBusClient
	dbh     *gorm.DB
	ana     *Analizer
	blmflt  *bloom.BloomFilter
	cuckflt *cuckoofilter.CuckooFilter
	wflt    *goac.AhoCorasick
	grph    *fact.Fact
	fch     chan *EventR
	fch2    chan *File
}

func someTask(t *tachymeter.Tachymeter, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()

	// doSomeSlowDbCall()

	// Task we're timing added here.
	t.AddTime(time.Since(start))
}

type Trees struct {
	Entries []*File
}

func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}

/*
func ReadFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var n int64

	if fi, err := f.Stat(); err == nil {
		if size := fi.Size(); size < 1e9 {
			n = size
		}
	}
	return readAll(f, n+bytes.MinRead)
}
*/

// type Entries []*File

func main() {
	defer funcTrack(time.Now())

	wallTimeStart := time.Now()

	fake, _ = faker.New("en")

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

	// go actx.ana.run2()
	// Mock()
	// select {}

	wflt := goac.NewAhoCorasick()
	wflt.AddPatterns("google", []string{"google", "angular", "googlecloudplatform", "googlechrome", "golang", "gwtproject", "zxing", "v8"})
	wflt.AddPatterns("facebook", []string{"facebook", "facebookarchive", "boltsframework"})
	wflt.AddPatterns("postgres", []string{"postgres", "postgresql"})
	wflt.AddPatterns("elasticsearch", []string{"elastic", "elasticsearch"})
	wflt.AddPatterns("mongodb", []string{"mongodb", "mongo"})
	wflt.AddPatterns("zeromq", []string{"zeromq", "zmq", "0mq"})
	wflt.AddPatterns("kubernetes", []string{"kubernetes", "k8s"})
	wflt.AddPatterns("boilerplate", []string{"boilerplate", "seed"})
	wflt.AddPatterns("phantom", []string{"phantom", "phantomjs"})
	wflt.AddPatterns("twitter", []string{"twbs", "twitter", "bower", "flightjs"})
	wflt.AddPatterns("microsoft", []string{"microsoft", "dotnet", "aspnet", "exceptionless", "mono", "winjs"})

	wflt.AddPatterns("npm", []string{"package.json", "package-lock.json", "npm-shrinkwrap.json", "yarn.lock"})
	wflt.AddPatterns("Maven", []string{"pom.xml", "ivy.xml", "build.gradle"})
	wflt.AddPatterns("RubyGems", []string{"Gemfile", "Gemfile.lock", "gems.rb", "gems.locked", "*.gemspec"})
	wflt.AddPatterns("Packagist", []string{"composer.json", "composer.lock"})
	wflt.AddPatterns("PyPi", []string{"setup.py", "req*.txt", "req*.pip", "requirements/*.txt", "requirements/*.pip", "Pipfile", "Pipfile.lock"})
	wflt.AddPatterns("Nuget", []string{"packages.config", "Project.json", "Project.lock.json", "*.nuspec", "paket.lock", "*.csproj"})
	wflt.AddPatterns("Bower", []string{"bower.json"})
	wflt.AddPatterns("CPAN", []string{"META.json", "META.yml"})
	wflt.AddPatterns("CocoaPods", []string{"Podfile", "Podfile.lock", "*.podspec"})
	wflt.AddPatterns("Clojars", []string{"project.clj"})
	wflt.AddPatterns("Meteor", []string{"versions.json"})
	wflt.AddPatterns("CRAN", []string{"DESCRIPTION"})
	wflt.AddPatterns("Cargo", []string{"Cargo.toml", "Cargo.lock"})
	wflt.AddPatterns("Hex", []string{"mix.exs", "mix.lock"})
	wflt.AddPatterns("Swift", []string{"Package.swift"})
	wflt.AddPatterns("Pub", []string{"pubspec.yaml", "pubspec.lock"})
	wflt.AddPatterns("Carthage", []string{"Cartfile", "Cartfile.private", "Cartfile.resolved"})
	wflt.AddPatterns("Dub", []string{"dub.json", "dub.sdl"})
	wflt.AddPatterns("Julia", []string{"REQUIRE"})
	wflt.AddPatterns("Shards", []string{"shard.yml", "shard.lock"})
	wflt.AddPatterns("Go", []string{"glide.yaml", "glide.lock", "Godeps", "Godeps/Godeps.json", "vendor/manifest", "vendor/vendor.json"})
	wflt.AddPatterns("Elm", []string{"elm-package.json", "elm_dependencies.json", "elm-stuff/exact-dependencies.json"})

	wflt.Build()

	// fds := MockDataset(actx.fctx)
	// fds := MockTree(actx.fctx)

	fds := flow.New()
	fds.Source(func(out chan string) {

		// fileList := make([]*File, 0)
		// var fileList []*File

		file, err := ioutil.ReadFile("files.json")
		if err != nil {
			// log.Fatalf("File error: %v\n", err)
			os.Exit(1)
			return
		}

		var arrResult []map[string]interface{}
		err = json.Unmarshal(file, &arrResult)
		if err != nil {
			panic(err)
		}

		// fmt.Printf("%#v", keys)
		for _, entry := range arrResult {
			if entry["path"].(string) != "" {
				// pp.Println("*** entry.path=", entry["path"].(string))
				out <- entry["path"].(string)
			}
		}

		// close(out)

	}, 4).Map(func(p string) string {
		return p

	}).Filter(func(p string) bool {
		return p != ""

	}).Map(func(p string) (sr []*ScanResult) {
		// defer glowTrack(p, time.Now())
		start := time.Now()
		sr = make([]*ScanResult, 0)
		results := wflt.Scan(p)

		for _, result := range results {
			sr = append(sr, &ScanResult{
				tags:  string([]rune(p)[result.Start : result.End+1]),
				group: result.Group,
				start: result.Start,
				end:   result.End + 1,
			})
			log.Println("file.path", p, "match=", string([]rune(p)[result.Start:result.End+1]), ", group=", result.Group, ", start=", result.Start, ", end=", result.End+1)
		}

		t.AddTime(time.Since(start))
		return
	}) /*.Filter(func(sr []*ScanResult) bool {
		return len(sr) > 0
	})
		var wg sync.WaitGroup

		// Run tasks.
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go someTask(t, wg)
		}

		wg.Wait()
	*/

	/*
		.Map(func(file *File, out chan flow.KeyValue) {
			log.Println("Map().line=", file)
			for _, p := range file.path {
				out <- flow.KeyValue{p, 1}
			}

		}).Map(func(f File) (r ScanResult) {

		}).ReduceByKey(func(x int, y int) int {
			return x + y

		}).Map(func(tag string, count int) flow.KeyValue {
			return flow.KeyValue{count, tag}

		}).Sort(func(a, b int) bool {
			return a < b

		}).Map(func(count int, tag string) {
			fmt.Printf("%d %s\n", count, tag)

		})
	*/

	fds.Run()

	// When finished, set elapsed wall time.
	t.SetWallTime(time.Since(wallTimeStart))

	// Rate outputs will be accurate.
	fmt.Println(t.Calc().String())

}

func NewAnaContext() *AnaContext {
	this := &AnaContext{}

	this.fctx = flow.New()
	this.blmflt = bloom.New(100000, 5)
	this.cuckflt = cuckoofilter.NewDefaultCuckooFilter()

	// patterns with group
	this.wflt = goac.NewAhoCorasick()
	this.load()

	// graph
	this.grph = fact.NewFact("./xfacts.db")
	this.graphT()

	this.fch2 = make(chan *File, 0)
	this.fch = make(chan *EventR, 0)
	this.ana = NewAnalizer()

	return this
}

type Analizer struct {
	fctx *flow.FlowContext
}

func NewAnalizer() *Analizer {
	this := &Analizer{}
	this.fctx = flow.New()
	return this
}

var manifestFiles = []string{
	"package.json", "package-lock.json", "npm-shrinkwrap.json", "yarn.lock",
	"pom.xml", "ivy.xml", "build.gradle",
	"packages.config", "Project.json", "Project.lock.json",
	"Gemfile", "Gemfile.lock", "gems.rb", "gems.locked", "bower.json",
	"META.json", "META.yml",
	"Podfile", "Podfile.lock",
	"glide.yaml", "glide.lock", "Godeps", "Godeps/Godeps.json", "vendor/manifest", "vendor/vendor.json",
	"elm-package.json", "elm_dependencies.json", "elm-stuff/exact-dependencies.json",
	"shard.yml", "shard.lock",
	"Package.swift",
	"pubspec.yaml", "pubspec.lock",
	"Cartfile", "Cartfile.private", "Cartfile.resolved",
	"CMakeLists.txt", "HunterGate.cmake",
}

func ReadFile(pathToDoc string) ([]byte, error) {
	return ioutil.ReadFile(pathToDoc)
}

func ErrorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func FixIdSyntax(id string) string {
	return strings.Replace(id, "/", ".", 1)
}

func MockTree(context *flow.FlowContext) (ret *flow.Dataset) {

	fileList := make([]*File, 0)
	// var data []*File

	file, e := ioutil.ReadFile("./files.json")
	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	// log.Println("%s", string(file))

	// err := json.Unmarshal([]byte(gitTree), &data)
	// err := json.Unmarshal(file, &data)
	err := json.Unmarshal(file, &fileList)
	if err != nil {
		log.Printf("Unmarshal error: %v\n", err)
	}

	input := make(chan interface{})

	go func() {
		for _, data := range fileList {
			input <- data
		}
		close(input)
	}()

	ret = context.Channel(input)
	return
}

func tokenize(row []interface{}) error {
	if visit, ok := row[0].(map[interface{}]interface{}); ok {
		for k, v := range visit {
			if w, ok := k.(string); ok {
				println("key:", w)
			}
			if t, ok := v.(string); ok {
				println("v:", t)
			}
		}
	}
	return nil
}

func MockDataset(fctx *flow.FlowContext) (ret *flow.Dataset) {
	s1 := rand.NewSource(time.Now().Unix())
	r1 := rand.New(s1)
	fileList := make([]*File, 500000)
	input := make(chan *File)
	var count int

	for i := 0; i < 500000; i++ {
		key := r1.Intn(len(manifestFiles))
		manifest := fmt.Sprintf("%s/%s", strings.Replace(fake.URL(), "http://", "", -1), manifestFiles[key])
		fileList[count] = &File{path: manifest}
		count++
	}

	go func() {
		for _, data := range fileList {
			input <- data
		}
		close(input)
	}()

	ret = fctx.Channel(input)
	return
}

func MockChannel() { // (ret *flow.Dataset) {
	s1 := rand.NewSource(time.Now().Unix())
	r1 := rand.New(s1)
	fileList := make([]*File, 1000)
	var count int
	for i := 0; i < 100; i++ {
		key := r1.Intn(len(manifestFiles))
		manifest := manifestFiles[key]
		fileList[count] = &File{path: manifest}
		count++
	}

	go func() {
		for _, data := range fileList {
			actx.fch2 <- data
		}
		close(actx.fch2)
	}()
	return
}

func (this *Analizer) run2() {
	fds := this.fctx.Channel(actx.fch2)
	fds.Filter(func(line *File) bool {
		if line == nil {
			return false
		}
		log.Println("Filter() file.path=", line.path)
		if line.path == "" {
			return true
		} else {
			return false
		}

	}).Map(func(line *File, ch chan rune) {
		log.Println("Map().line=", line)
		for _, p := range line.path {
			ch <- p
		}
	}).Map(func(r rune) (int, int8) {
		return 1, 2
	}).Reduce(func(x flow.KeyValue, y flow.KeyValue) flow.KeyValue {
		log.Println("Reduce().line=", x, y)
		return flow.KeyValue{Key: 5, Value: 6}
	})
	fds.Run()
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

/*
func loadFiles(inputFile string) {
	var data []*Files
	err := json.Unmarshal([]byte(gitTree), &data)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	for k, v := range data.Entries {
		fileList[k] = v
	}
}
*/

/*
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
*/

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

func scanContent(input string) (res []*ScanResult) {
	results := ac.Scan(input)
	fmt.Println("Matches: ", input)
	for _, result := range results {
		fmt.Println("match=", string([]rune(input)[result.Start:result.End+1]), ", group=", result.Group, ", start=", result.Start, ", end=", result.End+1)
		res = append(res, &ScanResult{
			match: string([]rune(input)[result.Start : result.End+1]),
			group: result.Group,
			start: result.Start,
			end:   result.End + 1,
		})
	}
	return res
}

func (a *AnaContext) load() {

	a.wflt.AddPatterns("google", []string{"google", "angular", "googlecloudplatform", "googlechrome", "golang", "gwtproject", "zxing", "v8"})
	a.wflt.AddPatterns("facebook", []string{"facebook", "facebookarchive", "boltsframework"})
	a.wflt.AddPatterns("postgres", []string{"postgres", "postgresql"})
	a.wflt.AddPatterns("elasticsearch", []string{"elastic", "elasticsearch"})
	a.wflt.AddPatterns("mongodb", []string{"mongodb", "mongo"})
	a.wflt.AddPatterns("zeromq", []string{"zeromq", "zmq", "0mq"})
	a.wflt.AddPatterns("kubernetes", []string{"kubernetes", "k8s"})
	a.wflt.AddPatterns("boilerplate", []string{"boilerplate", "seed"})
	a.wflt.AddPatterns("phantom", []string{"phantom", "phantomjs"})
	a.wflt.AddPatterns("twitter", []string{"twbs", "twitter", "bower", "flightjs"})
	a.wflt.AddPatterns("microsoft", []string{"microsoft", "dotnet", "aspnet", "exceptionless", "mono", "winjs"})
	a.wflt.Build()

}

func (a *AnaContext) graphT() {

	// Declare truth
	a.grph.Let("cat").Has("name", "cat")
	a.grph.Let("dog").Has("name", "hou")
	a.grph.Let("cat").Has("white", "black")

	fmt.Println(a.grph.What("cat", "name"))
	fmt.Println(a.grph.WhoHas("name", "cat"))
	fmt.Println(a.grph.WhoHas("name", "woof"))
	fmt.Println(a.grph.What(a.grph.What("cat", "name"), "color"))
	fmt.Println(a.grph.Stringify(a.grph.WhoHas("name", "cat")))

	fmt.Println(a.grph.What("time"))
}
func loadPatterns() {
	ac := goac.NewAhoCorasick()
	ac.AddPatterns("google", []string{"google", "angular", "googlecloudplatform", "googlechrome", "golang", "gwtproject", "zxing", "v8"})
	ac.AddPatterns("facebook", []string{"facebook", "facebookarchive", "boltsframework"})
	ac.AddPatterns("postgres", []string{"postgres", "postgresql"})
	ac.AddPatterns("elasticsearch", []string{"elastic", "elasticsearch"})
	ac.AddPatterns("mongodb", []string{"mongodb", "mongo"})
	ac.AddPatterns("zeromq", []string{"zeromq", "zmq", "0mq"})
	ac.AddPatterns("kubernetes", []string{"kubernetes", "k8s"})
	ac.AddPatterns("boilerplate", []string{"boilerplate", "seed"})
	ac.AddPatterns("phantom", []string{"phantom", "phantomjs"})
	ac.AddPatterns("twitter", []string{"twbs", "twitter", "bower", "flightjs"})
	ac.AddPatterns("microsoft", []string{"microsoft", "dotnet", "aspnet", "exceptionless", "mono", "winjs"})
	ac.Build()

}

var packageManagers = []string{"atom", "bower", "cargo", "carthage", "clojars", "cocoapods", "cpan", "cran", "dub", "elm", "go", "hex", "julia", "maven", "meteor", "npm", "nuget", "packagist", "pub", "pypi", "rubygems", "shards", "swiftpm"}

/*
var packageManifests = map[string]map[string][]string{
	map[string]map[string]{"manifests": map[string][]string{
		"npm": []string{"package.json", "package-lock.json", "npm-shrinkwrap.json", "yarn.lock"},
		"Maven": []string{"pom.xml", "ivy.xml", "build.gradle"},
		"RubyGems": []string{"Gemfile", "Gemfile.lock", "gems.rb", "gems.locked", "*.gemspec"}},
		"Packagist": []string{"composer.json", "composer.lock"}},
		"PyPi": []string{"setup.py", "req*.txt", "req*.pip", "requirements/*.txt", "requirements/*.pip", "Pipfile", "Pipfile.lock"}},
		"Nuget": []string{"packages.config", "Project.json", "Project.lock.json", "*.nuspec", "paket.lock", "*.csproj"}},
		"Bower": []string{"bower.json"}},
		"CPAN": []string{"META.json", "META.yml"}},
		"CocoaPods": []string{"Podfile", "Podfile.lock", "*.podspec"}},
		"Clojars": []string{"project.clj"}},
		"Meteor": []string{"versions.json"}},
		"CRAN": []string{"DESCRIPTION"}},
		"Cargo": []string{"Cargo.toml", "Cargo.lock"}},
		"Hex": []string{"mix.exs", "mix.lock"}},
		"Swift": []string{"Package.swift"}},
		"Pub": []string{"pubspec.yaml", "pubspec.lock"}},
		"Carthage": []string{"Cartfile", "Cartfile.private", "Cartfile.resolved"}},
		"Dub": []string{"dub.json", "dub.sdl"}},
		"Julia": []string{"REQUIRE"}},
		"Shards": []string{"shard.yml", "shard.lock"}},
		"Go": []string{"glide.yaml", "glide.lock", "Godeps", "Godeps/Godeps.json", "vendor/manifest", "vendor/vendor.json"}},
		"Elm": []string{"elm-package.json", "elm_dependencies.json", "elm-stuff/exact-dependencies.json"}},
	},
	},
}
*/

/*
var searchKeywords = map[string]map[string][]string{
	"brands": map[string][]string{
		"google":    []string{"google", "angular", "googlecloudplatform", "googlechrome", "golang", "gwtproject", "zxing", "v8"},
		"twitter":   []string{"twbs", "twitter", "bower", "flightjs"},
		"facebook":  []string{"facebook", "facebookarchive", "boltsframework"},
		"github":    []string{"atom", "github"},
		"microsoft": []string{"microsoft", "dotnet", "aspnet", "exceptionless", "mono", "winjs"},
	},
	"keywords":
		"node":                []string{"node", "nodejs"}},
		"jquery":              []string{"jquery", "jq", "/^jq[\\-]?/"},
		"grunt":               []string{"grunt", "gruntjs"},
		"angular":             []string{"angular", "angularjs", "ng", "/^ng(?!inx)\\-]?/"},
		"ember":               []string{"emberjs", "ember"},
		"meteor":              []string{"meteor", "meteorjs"},
		"gulp":                []string{"gulp"},
		"express":             []string{"express", "expressjs"},
		"d3":                  []string{"d3"},
		"polymer":             []string{"polymer"},
		"ionic":               []string{"ionic"},
		"seajs":               []string{"seajs"},
		"yeoman":              []string{"yeoman"},
		"browserify":          []string{"browserify"},
		"requirejs":           []string{"requirejs"},
		"underscore":          []string{"underscore", "underscorejs"},
		"modernizr":           []string{"modernizr"},
		"phantom":             []string{"phantom", "phantomjs"},
		"metalsmith":          []string{"metalsmith"},
		"bootstrap":           []string{"bootstrap"},
		"django":              []string{"django"},
		"bottle":              []string{"bottlepy", "bottle"},
		"web2py":              []string{"web2py"},
		"webpy":               []string{"webpy"},
		"flask":               []string{"flask"},
		"ipython":             []string{"ipython"},
		"fabric":              []string{"fabric"},
		"celery":              []string{"celery"},
		"language/python":     []string{"python", "/^py/"},
		"language/ruby":       []string{"ruby"},
		"language/clojure":    []string{"clojure"},
		"language/lisp":       []string{"lisp"},
		"language/rust":       []string{"rust"},
		"language/erlang":     []string{"erlang"},
		"language/go":         []string{"golang", "go"},
		"language/javascript": []string{"javascript", "js"},
		"language/clojure":    []string{"coffeescript"},
		"language/php":        []string{"php"},
		"language/perl":       []string{"perl"},
		"language/swift":      []string{"swift"},
		"language/css":        []string{"css", "stylesheet"},
		"ios":                 []string{"ios"},
		"osx":                 []string{"osx"},
		"unix":                []string{"unix"},
		"android":             []string{"android"},
		"linux":               []string{"linux"},
		"windows":             []string{"windows"},
		"deprecated":          []string{"deprecated"},
		"pdf":                 []string{"pdf"},
		"polyfill":            []string{"polyfill"},
		"framework":           []string{"framework"},
		"dropbox":             []string{"dropbox"},
		"webkit":              []string{"webkit"},
		"sql":                 []string{"sql"},
		"svg":                 []string{"svg"},
		"boilerplate":         []string{"boilerplate", "seed"},
		"rails":               []string{"rails", "rails3"},
		"vim":                 []string{"vim", "vi"},
		"git":                 []string{"git"},
		"backbone":            []string{"backbone"},
		"docker":              []string{"docker"},
		"emacs":               []string{"emacs"},
		"redis":               []string{"redis"},
		"chrome":              []string{"chrome"},
		"sublime":             []string{"sublime"},
		"vagrant":             []string{"vagrant"},
		"wordpress":           []string{"wordpress", "/^wp\\-/"},
		"youtube":             []string{"youtube"},
		"apache":              []string{"apache"},
		"jekyll":              []string{"jekyll"},
		"puppet":              []string{"puppet"},
		"sass":                []string{"sass", "scss"},
		"nginx":               []string{"nginx"},
		"markdown":            []string{"markdown"},
		"elasticsearch":       []string{"elasticsearch"},
		"chef":                []string{"chef"},
		"mongodb":             []string{"mongodb", "mongo"},
		"cordova":             []string{"cordova"},
		"phonegap":            []string{"phonegap"},
		"ansible":             []string{"ansible"},
		"openshift":           []string{"openshift"},
		"mysql":               []string{"mysql"},
		"couchbase":           []string{"couchbase"},
		"firebase":            []string{"firebase"},
		"homebrew":            []string{"homebrew"},
		"openstack":           []string{"openstack"},
		"maven":               []string{"maven"},
		"hadoop":              []string{"hadoop"},
		"spark":               []string{"spark"},
		"jasmine":             []string{"jasmine"},
		"hubot":               []string{"hubot"},
		"jruby":               []string{"jruby"},
		"couchdb":             []string{"couchdb"},
		"travis":              []string{"travis"},
		"bash":                []string{"bash"},
		"coreos":              []string{"coreos"},
		"mustache":            []string{"mustache"},
		"zsh":                 []string{"zsh"},
		"jenkins":             []string{"jenkins"},
		"cassandra":           []string{"cassandra"},
		"statsd":              []string{"statsd"},
		"eclipse":             []string{"eclipse"},
		"knockout":            []string{"knockout"},
		"graphite":            []string{"graphite"},
		"textmate":            []string{"textmate"},
		"jed":                 []string{"jed"},
		"memcached":           []string{"memcached"},
		"mesos":               []string{"mesos"},
		"rabbitmq":            []string{"rabbitmq"},
		"firefox":             []string{"firefox", "ff"},
		"postgres":            []string{"postgres", "postgresql"},
		"selenium":            []string{"selenium"},
		"gems":                []string{"gems", "rubygems"},
		"zeromq":              []string{"zeromq", "zmq", "0mq"},
		"tmux":                []string{"tmux"},
		"cyanogenmod":         []string{"cyanogenmod"},
		"tornado":             []string{"tornado"},
		"octopress":           []string{"octopress"},
		"dokku":               []string{"dokku"},
		"karma":               []string{"karma"},
		"bitcoin":             []string{"bitcoin"},
		"handlebars":          []string{"handlebars"},
		"qt":                  []string{"qt"},
		"minecraft":           []string{"minecraft"},
		"unity":               []string{"unity"},
		"cocos2d":             []string{"cocos2d"},
		"openssl":             []string{"openssl"},
		"amqp":                []string{"amqp"},
		"logstash":            []string{"logstash"},
		"sqlite":              []string{"sqlite"},
		"v8":                  []string{"v8"},
		"fuse":                []string{"fuse"},
		"cocoa":               []string{"cocoa"},
		"curl":                []string{"curl"},
		"ffmpeg":              []string{"ffmpeg"},
		"hhvm":                []string{"hhvm"},
		"rake":                []string{"rake"},
		"drupal":              []string{"drupal"},
		"gevent":              []string{"gevent"},
		"nagios":              []string{"nagios"},
		"chromium":            []string{"chromium"},
		"jenkinsci":           []string{"jenkinsci"},
		"etcd":                []string{"etcd"},
		"kubernetes":          []string{"kubernetes"},
		"react":               []string{"react", "reactjs"},
	},
}
*/
