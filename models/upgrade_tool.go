package models

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"upgrade.maiyajia.com/services/mongo"
)

type Tool struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	DownloadURL string        `bson:"download_url" json:"download_url"`
	Name        string        `bson:"name" json:"name"`
	Relpath     string        `bson:"relpath" json:"relpath"` //course的URL
	Title       string        `bson:"title" json:"title"`
	Category    string        `bson:"category" json:"category"`
	Version     string        `bson:"version" json:"version"`
	Icon        string        `bson:"icon" json:"icon"`
	Purchased   bool          `bson:"purchased" json:"purchased"`
	Weight      int           `bson:"weight" json:"weight"` //权重
	Types       []string      `bson:"types" json:"types"`
	CreateTime  time.Time     `bson:"createTime" json:"createTime"`
}

// UpCourse 获取学习平台的课程
type UpCourse struct {
	ID      bson.ObjectId `bson:"_id" json:"_id"`
	Version string        `bson:"version" json:"version"`
	Name    string        `bson:"name" json:"name"` //课程名称
}

// Course 课程
type Course struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`                 //课程名称
	DownloadURL string        `bson:"download_url" json:"download_url"` //下载课程包的URL
	Version     string        `bson:"version" json:"version"`
	Icon        string        `bson:"icon" json:"icon"`           //icon的URL
	Category    string        `bson:"category" json:"category"`   //所属类别
	Desc        string        `bson:"desc" json:"desc"`           //课程描述
	Onsell      bool          `bson:"onsell" json:"onsell"`       //是否下载成功
	Purchased   bool          `bson:"purchased" json:"purchased"` //是否已经购买
	Relpath     string        `bson:"relpath" json:"relpath"`     //
	BasePath    string        `bson:"basepath" json:"base_path"`  //
	Lessions    []*Lession    `bson:"lessions" json:"lessions"`   //课时信息
	CreateTime  int64         `bson:"createTime" json:"createTime"`
	Browse      int           `bson:"browse" json:"browse"` //课程浏览量
}

// Lession 课时,即一节课
type Lession struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`                 //课节名称
	IconURL     string        `bson:"icon_url" json:"icon_url"`         //图标路径
	QuestionURL string        `bson:"question_url" json:"question_url"` //答题路径
	Contents    []*Content    `bson:"content" json:"content"`           //课节视频
	Tool        string        `bson:"tool" json:"tool"`
}

//Content 学习资源信息
type Content struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	VideoName string        `bson:"video_name" json:"video_name"` //视频名称
	VideoURL  string        `bson:"video_url" json:"video_url"`   //视频路径
	MdURL     string        `bson:"md_url" json:"md_url"`         //课时md路径
}

//LessListRes 课程内容列表
type LessListRes struct {
	Code    int       `json:"code"`
	Courses []*Course `json:"courses"`
}

//GetToolsInfo 获取所有工具信息，用于升级新增
func GetToolsInfo() ([]*Tool, error) {
	var tools []*Tool
	f := func(col *mgo.Collection) error {
		return col.Find(nil).Sort("-createTime").All(&tools)
	}
	err := mongo.Client.Do(beego.AppConfig.String("MongoDB"), "tool", f)
	logs.Info("Begin getTools in models:", err)
	return tools, err
}

//GetCoursesInfo 获取课程购买信息,从课程平台api获取已购买课程。
func GetCoursesInfo(url, productKey, productSerial string) ([]*Course, error) {
	var courses LessListRes
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logs.Error("NewRequest erro", err)
		return nil, err
	}
	req.Header.Set("productKey", productKey)
	req.Header.Set("productSerial", productSerial)
	resp, err := client.Do(req)
	if err != nil {
		logs.Error("resp Error", err)
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	courseJSON := buf.String()
	//logs.Info("courseJSON:", courseJSON)
	err = json.Unmarshal([]byte(courseJSON), &courses)
	if err != nil {
		logs.Info("unmarshal json err:", err)
		return nil, err
	}
	return courses.Courses, nil
}
