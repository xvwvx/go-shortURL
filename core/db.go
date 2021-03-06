package core

import (
	"fmt"

	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"shortURL/base62"
	"sync/atomic"
	"time"
)

var (
	db *sql.DB

	globalLastId = uint64(0)

	insertErr   = errors.New("insert error")
	notFoundErr = errors.New("notfound")
)

func LinkDB() {
	var err error

	dataSourceName := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8mb4`, Conf.DBUsername, Conf.DBPassword, Conf.DBHost,
		Conf.DBDatabase)
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(1 * time.Hour)

	//保留 10000 及以内的id
	lastId := lastId()
	if lastId < 10000 {
		lastId = 10000
	}
	atomic.StoreUint64(&globalLastId, lastId)
}

//获取最后一个id
func lastId() (value uint64) {
	rows, err := db.Query("SELECT `id` FROM `short_url` ORDER BY	`id` DESC LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&value)
	}

	return
}

//插入数据
func insert(urlLong string) (suffix string, err error) {
	suffix, err = findWithURLLong(urlLong)
	if err != nil {
		return "", err
	} else if suffix != "" {
		return suffix, nil
	}

	stmt, err := db.Prepare("INSERT INTO `short_url`(`id`, `url_long`, `url_short`) VALUES(?, ?, ?)")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	id := atomic.AddUint64(&globalLastId, 1)
	suffix = base62.Encode(id)

	res, err := stmt.Exec(id, urlLong, suffix)
	if err != nil {
		return "", err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	if affect == 0 {
		return "", insertErr
	}

	return suffix, nil
}

//查询已存在的值
func findWithURLLong(urlLong string) (urlShort string, err error) {
	rows, err := db.Query("SELECT `url_short` FROM short_url WHERE `url_long`=? LIMIT 1", urlLong)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&urlShort)
	}

	return
}

func find(suffix string) (result string, err error) {
	id, err := base62.Decode(suffix)
	if err != nil || id == 0 || id > globalLastId {
		return "", notFoundErr
	}

	rows, err := db.Query("SELECT `url_long` FROM short_url WHERE `id`=? LIMIT 1", id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&result)
	}

	if result == "" {
		return "", notFoundErr
	}

	return

}
