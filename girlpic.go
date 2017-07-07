package main

import "github.com/garyburd/redigo/redis"

//Pic 定义图片数据
type Pic struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func getPicNum(key string) int {
	conn := rdb.Get()
	defer conn.Close()
	count, err := redis.Int(conn.Do("ZCARD", key))
	//debug("pic count:", count)
	if err != nil {
		errorlog("get length of", key, "list fail with error:", err)
		return 0
	}
	return count
}

func getPics(page int) []Pic {
	num := getPicNum(listChecked)
	if num == 0 {
		return make([]Pic, 0)
	}
	var start int
	if start = (page - 1) * numPerPage; start < 0 || start > num {
		start = 0
	}
	end := start + numPerPage - 1
	conn := rdb.Get()
	defer conn.Close()
	objs, err := redis.Strings(conn.Do("ZRANGE", listChecked, start, end))
	//debug("start:", start, "end:", end)
	if err != nil {
		errorlog("get pics from", listChecked, "fail with error:", err)
		return make([]Pic, 0)
	}
	pics := make([]Pic, numPerPage)
	for i, v := range objs {
		t := Pic{}
		t.ID = start + i
		t.URL = v
		pics[i] = t
	}
	return pics
}

func review() error {
	conn := rdb.Get()
	defer conn.Close()
	_, err := conn.Do("LREM", listTemp, 0, "deleted")
	if err != nil {
		return err
	}
	urls, err := redis.Strings(conn.Do("lrange", listTemp, 0, -1))
	if err != nil {
		return err
	}
	conn.Do("del", listTemp)
	if len(urls) > 0 {
		for _, url := range urls {
			conn.Send("ZADD", listChecked, 0, url)
		}
		_, err := conn.Do("")
		if err != nil {
			return err
		}
	}
	return nil
}

func getPicWaitReview() []Pic {
	err := review()
	if err != nil {
		errorlog("move checked pics to reviewed list fail with error:", err)
	}
	pics := make([]Pic, 0)
	count := getPicNum(listUnchecked)
	if count == 0 {
		return pics
	}
	conn := rdb.Get()
	defer conn.Close()
	urls, err := redis.Strings(conn.Do("ZRANGE", listUnchecked, 0, func() int {
		if count >= numPerPage {
			return numPerPage - 1
		}
		return -1
	}()))
	if err != nil {
		errorlog("move pics to temp list fail with error:", err, "and notviewed pics count:", count)
		return make([]Pic, 0)
	}
	conn.Do("ZREMRANGEBYRANK", listUnchecked, 0, func() int {
		if count >= numPerPage {
			return numPerPage - 1
		}
		return -1
	}())
	for i, v := range urls {

		conn.Do("RPUSH", listTemp, v)
		tp := Pic{
			ID:  i,
			URL: v,
		}
		pics = append(pics, tp)
	}
	return pics
}
