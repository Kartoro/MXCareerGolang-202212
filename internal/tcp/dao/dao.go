package dao

import (
	"MXCareerGolang-202212/config"
	"MXCareerGolang-202212/internal/model"
	"MXCareerGolang-202212/internal/util"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", config.DSN)
	if err != nil {
		log.Fatal("connect mysql failed, err: " + err.Error())
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	log.Println("mysql connected")
}

func CreateUser(userName, nickName, password string) (err error) {
	password = util.SHA256(password)
	_, err = db.Exec("INSERT INTO user (username, nickname, password, avatar) VALUES (?,?,?,?)", userName, nickName, password, config.DefaultImagePath)
	if err != nil {
		log.Println("insert user failed", err)
		return err
	}
	return nil
}

func LoginAuth(userName string, password string) (bool, error) {
	var user model.User
	rows, err := db.Query("SELECT * FROM user WHERE username = ? LIMIT 1", userName)
	if rows == nil || err != nil {
		log.Println(err)
		return false, err
	}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.UserName, &user.NickName, &user.Password, &user.Avatar)
		if err != nil {
			log.Println(err)
			return false, err
		}
	}
	defer rows.Close()
	pwd := util.SHA256(password)
	if user.Password == pwd {
		return true, nil
	}
	return false, nil
}

func GetProfile(userName string) (model.User, bool) {
	var user model.User
	rows, err := db.Query("SELECT * FROM user WHERE username = ? LIMIT 1", userName)
	if err != nil {
		log.Println(err)
		return model.User{}, false
	}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.UserName, &user.NickName, &user.Password, &user.Avatar)
		if err != nil {
			log.Println(err)
			return model.User{}, false
		}
	}
	defer rows.Close()
	return user, true
}

func UpdateNickName(userName, nickName string) (bool, error) {
	_, err := db.Exec("UPDATE user SET nickname = ? WHERE username = ? LIMIT 1", nickName, userName)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

func UpdateAvatar(userName, avatar string) (bool, error) {
	_, err := db.Exec("UPDATE user SET avatar = ? WHERE username = ? LIMIT 1", avatar, userName)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}
