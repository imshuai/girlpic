package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

var rdb *redis.Pool

type dbConfig struct {
	ServerName string `json:"server"`
	ServerPort int    `json:"port"`
	ServerDB   string `json:"db"`
	ServerUser string `json:"username"`
	ServerAuth string `json:"password"`
	Charset    string `json:"charset"`
	MaxIdle    int    `json:"max_idle_conn"`
	MaxActive  int    `json:"max_active_conn"`
}

func redisInit() error {
	rcfg, err := readConfig("db_config.json")
	if err != nil {
		return err
	}
	rdb = &redis.Pool{
		MaxIdle:     rcfg.MaxIdle,
		MaxActive:   rcfg.MaxActive,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", rcfg.ServerName+":"+strconv.Itoa(rcfg.ServerPort))
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			if rcfg.ServerAuth != "" {
				_, err := c.Do("AUTH", rcfg.ServerAuth)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			// 选择db
			c.Do("SELECT", rcfg.ServerDB)
			return c, nil
		},
	}
	conn := rdb.Get()
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		conn.Close()
		return err
	}
	conn.Close()
	return nil
}

func readConfig(filepath string) (*dbConfig, error) {
	byts, err := ioutil.ReadFile(filepath)
	if err != nil {
		return &dbConfig{}, err
	}
	c := &dbConfig{}
	err = json.Unmarshal(byts, c)
	if err != nil {
		return &dbConfig{}, err
	}
	return c, nil
}
