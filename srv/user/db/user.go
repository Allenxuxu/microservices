package db

import (
	"fmt"
	pb "github.com/Allenxuxu/microservices/srv/user/proto/user"

	"github.com/micro/go-log"

	config "github.com/micro/go-config"
	"github.com/micro/go-config/source/consul"

	"github.com/jinzhu/gorm"
	//初始化数据库驱动
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type dbInfo struct {
	Address      string `json:"address"`
	Port         int    `json:"port"`
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
	DbName       string `json:"db_name"`
}

func init() {
	consulSource := consul.NewSource()
	conf := config.NewConfig()

	// Load file source
	err := conf.Load(consulSource)
	if err != nil {
		log.Fatal(err)
	}
	var v dbInfo
	err = conf.Get("micro", "config", "database", "user").Scan(&v)
	if err != nil {
		log.Fatal(err)
	}

	log.Log(v)
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		v.UserName, v.UserPassword, v.Address, v.Port, v.DbName))
	if err != nil {
		log.Fatal("failed to connect database：", err)
	}

	db.AutoMigrate(&pb.User{})
	db.Model(&pb.User{}).AddUniqueIndex("uIndex_email", "email")
	db.Model(&pb.User{}).AddUniqueIndex("uIndex_tel", "tel")
}

// CreateUser 在数据库中创建一个用户
func CreateUser(user *pb.User) error {
	return db.Create(user).Error
}

// DelUser 删除用户
func DelUser(user *pb.User) error {
	return db.Delete(user).Error
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(user *pb.User) error {
	return db.Model(user).Updates(*user).Error
}

// GetByID 通过id取用户信息
func GetByID(id string) (pb.User, error) {
	var user pb.User
	err := db.Where("id = ?", id).Find(&user).Error
	return user, err
}

// GetByTel 通过电话获取用户信息
func GetByTel(tel string) (pb.User, error) {
	var user pb.User
	err := db.Where("tel = ?", tel).Find(&user).Error
	return user, err
}

// GetByEmail 通过邮箱获取用户信息
func GetByEmail(email string) (pb.User, error) {
	var user pb.User
	err := db.Where("email = ?", email).Find(&user).Error
	return user, err
}

// GetAllUsers 获取所有用户信息
func GetAllUsers() ([]*pb.User, error) {
	var users []*pb.User
	err := db.Find(&users).Error
	return users, err
}
