package mongodb_official

import (
	"GIG/app/databases/index_manager"
	"github.com/revel/revel"
)

var MaxPool int
var PATH string
var DBNAME string

func LoadMongo() {
	MaxPool = revel.Config.IntDefault("mongo.maxPool", 0)
	PATH, _ = revel.Config.String("mongo.path")
	DBNAME, _ = revel.Config.String("mongo.database")
	CheckAndInitServiceConnection()

	// ensure db indexes
	index_manager.CreateDBIndexes(MongoOfficialIndexManager{})
}
