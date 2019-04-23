package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"upgrade.maiyajia.com/services/mongo"
)

var OS_ARCH = [4]string{"linux/docker", "linux/amd64", "windows/amd64", "windows/386"}

// Product 产品信息
type Product struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Version   string        `bson:"version" json:"version"`
	Assets    []*Asset      `bson:"assets" json:"assets"`
	Changelog string        `bson:"changelog" json:"changelog"`
	Date      time.Time     `bson:"date" json:"date" `
}

// Asset 产品安装包信息
type Asset struct {
	OS     string `bson:"os" json:"os"`
	Source string `bson:"source" json:"source"`
	Hash   string `bson:"hash" json:"hash"`
}

//QueryUpgradeProduct 版本检测是否为最新，最新则返回最新的产品信息
func QueryUpgradeProduct(version string) (*Product, bool) {
	//获取最新发布的产品
	var product *Product
	f := func(col *mgo.Collection) error {
		return col.Find(nil).Sort("-date").Limit(1).One(&product)
	}
	mongo.Client.Do(beego.AppConfig.String("MongoDB"), "product", f)

	//读取最新版本号，并进行比较
	prodVersion := product.Version
	isUpgrade := compareVersion(version, prodVersion)
	// 没有最新版本，不需要升级
	if !isUpgrade {
		return nil, isUpgrade
	}

	return product, true

}

//QueryUpgradeTool 工具检测是否为最新版本，最新则返回最新的产品信息
func QueryUpgradeTool(id bson.ObjectId, neme, version string) (*Tool, bool, error) {
	//获取最新发布的产品
	var tool *Tool
	f := func(col *mgo.Collection) error {
		return col.Find(bson.M{"_id": id}).Sort("-date").Limit(1).One(&tool)
		//return col.Find(bson.M{"name": name}).One(&tool)
	}
	err := mongo.Client.Do(beego.AppConfig.String("MongoDB"), "tool", f)

	//读取最新版本号，并进行比较
	prodVersion := tool.Version
	isUpgrade := compareVersion(version, prodVersion)
	// 没有最新版本，不需要升级
	if !isUpgrade {
		return nil, isUpgrade, err
	}

	return tool, true, err

}

//QueryUpgradeCourse 工具检测是否为最新版本，最新则返回最新的产品信息
func QueryUpgradeCourse(id bson.ObjectId, version string) (*Course, bool, error) {
	//获取最新发布的产品
	var course *Course
	f := func(col *mgo.Collection) error {
		return col.Find(bson.M{"_id": id}).Sort("-date").Limit(1).One(&course)
		//return col.Find(bson.M{"name": name}).One(&tool)
	}
	err := mongo.Client.Do(beego.AppConfig.String("MongoDB"), "course", f)
	//读取最新版本号，并进行比较
	prodVersion := course.Version
	isUpgrade := compareVersion(version, prodVersion)
	// 没有最新版本，不需要升级
	if !isUpgrade {
		return nil, isUpgrade, err
	}
	return course, true, err
}

//新老版本比较, 返回true表示有新的版本
func compareVersion(oldVersion, newVersion string) bool {
	old := strings.SplitN(oldVersion, ".", 3)
	new := strings.SplitN(newVersion, ".", 3)

	for i := 0; i < 3; i++ {
		old, _ := strconv.Atoi(old[i])
		new, _ := strconv.Atoi(new[i])
		if new > old {
			return true
		} else if old > new {
			return false
		}
	}
	return false
}
