package dbs

import (
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"time"
	"websocket_server/app/dbs/influxdb"
	"websocket_server/app/dbs/rdb"
	"websocket_server/util/logFile"
)

type Dbs interface {
	initCache()
	initRdb(log logFile.LogFile)
	initIdb(log logFile.LogFile)
	GetCache() *cache.Cache
	GetRdb() redis.UniversalClient
	GetIdb() HistoryDB
}

type HistoryDB interface {
	Close()
	Writer() api.WriteAPIBlocking
	Querier() api.QueryAPI
}

type dbs struct {
	Cache *cache.Cache
	Rdb   redis.UniversalClient
	Idb   HistoryDB
}

func NewDbs(log logFile.LogFile) Dbs {
	d := &dbs{}
	d.initCache()
	d.initRdb(log)
	d.initIdb(log)
	return d
}

func (d *dbs) initCache() {
	d.Cache = cache.New(5*time.Minute, 10*time.Minute)
}

func (d *dbs) initRdb(log logFile.LogFile) {
	d.Rdb = rdb.NewClient("redis")
	log.Info().Println("Redis Connection successful")
}

func (d *dbs) initIdb(log logFile.LogFile) {
	d.Idb = influxdb.NewInfluxdb("influxdb")
	log.Info().Println("InfluxDB Connection successful")
}

func (d *dbs) GetCache() *cache.Cache {
	return d.Cache
}

func (d *dbs) GetRdb() redis.UniversalClient {
	return d.Rdb
}

func (d *dbs) GetIdb() HistoryDB {
	return d.Idb
}
