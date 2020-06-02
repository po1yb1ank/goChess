package database

import (
	//"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)
var database *cache.Cache
var dataBaseStatus bool

func SetDataBase()  {
	database = cache.New(5*time.Minute, 10*time.Minute)
	dataBaseStatus = true
}
func DataBaseStatus() bool {
	return dataBaseStatus
}