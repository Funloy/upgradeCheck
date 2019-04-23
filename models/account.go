package models

import (
	"time"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"upgrade.maiyajia.com/services/mongo"
)

// Account 产品账户
type Account struct {
	ID             bson.ObjectId `bson:"_id" json:"-"`
	ProductKey     string        `bson:"productKey" json:"productKey"`
	ProductSerial  string        `bson:"productSerial" json:"productSerial"`
	ProductOS      string        `bson:"productOS" json:"productOS"`
	ProductVersion string        `bson:"productVersion" json:"productVersion"`
	Name           string        `bson:"name" json:"name"`
	Address        string        `form:"address"`
	Notes          string        `form:"notes"`
	Locked         bool          `bson:"locked" json:"locked"`         // 账户是否锁住
	CreateTime     time.Time     `bson:"createTime" json:"createTime"` //创建时间
}

//FindAccount 通过给定查询条件，返回查询到的账号信息
func FindAccount(query interface{}) (*Account, error) {
	var account *Account
	f := func(col *mgo.Collection) error {
		return col.Find(query).One(&account)
	}
	err := mongo.Client.Do(beego.AppConfig.String("MongoDB"), "account", f)
	return account, err
}

// AccountExists 通过给定查询条件，检查是否存在账号
func AccountExists(query interface{}) bool {
	f := func(col *mgo.Collection) error {
		return col.Find(query).Select(bson.M{"_id": 1}).One(&bson.M{})
	}
	if err := mongo.Client.Do(beego.AppConfig.String("MongoDB"), "account", f); err == mgo.ErrNotFound {
		return false
	}
	return true
}

// AccountUpgraded 客户产品升级后，更新客户账户的信息
func AccountUpgraded(key, serial, os, version string) error {
	query := bson.M{"productKey": key, "productSerial": serial}
	update := bson.M{"$set": bson.M{"productOS": os, "productVersion": version}}
	f := func(col *mgo.Collection) error {
		return col.Update(query, update)
	}
	return mongo.Client.Do(beego.AppConfig.String("MongoDB"), "account", f)

}
