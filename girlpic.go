package main

import "github.com/garyburd/redigo/redis"
import "encoding/json"

//GirlPic define database struct for table girl_pic
type GirlPic struct {
	URL    string `json:"url"`
	Like   int    `json:"like"`
	Unlike int    `json:"unlike"`
}

func (gp GirlPic) serialize() string {
	byts, err := json.Marshal(gp)
	if err != nil {
		errorlog("pic serialize fail with error:", err)
		return ""
	}
	return string(byts)
}

func deserialize(str string) GirlPic {
	gp := GirlPic{}
	err := json.Unmarshal([]byte(str), &gp)
	if err != nil {
		errorlog("pic deserialize fail with error:", err)
		return GirlPic{}
	}
	return gp
}

type tPic struct {
	ID int `json:"id"`
	GirlPic
}

func getPicNum(key string) int {
	conn := rdb.Get()
	defer conn.Close()
	count, err := redis.Int(conn.Do("llen", key))
	//debug("pic count:", count)
	if err != nil {
		errorlog("get length of", key, "list fail with error:", err)
		return 0
	}
	return count
}

func getPics(page int) []tPic {
	num := getPicNum(listChecked)
	if num == 0 {
		debug("number of pics is:", num)
		return nil
	}
	var start int
	if start = (page - 1) * numPerPage; start < 0 || start > num {
		start = 0
	}
	end := start + numPerPage - 1
	conn := rdb.Get()
	defer conn.Close()
	objs, err := redis.Strings(conn.Do("LRANGE", listChecked, start, end))
	//debug("objs:", objs)
	if err != nil {
		errorlog("get pics from", listChecked, "fail with error:", err)
		return make([]tPic, 0)
	}
	pics := make([]tPic, numPerPage)
	for i, v := range objs {
		obj := deserialize(v)
		t := tPic{}
		t.ID = start + i
		t.Like = obj.Like
		t.Unlike = obj.Unlike
		t.URL = obj.URL
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
	//debug("urls:", urls)
	if err != nil {
		return err
	}
	if len(urls) > 0 {
		for _, url := range urls {
			pic := GirlPic{
				URL: url,
			}
			conn.Send("RPUSH", listChecked, pic.serialize())
		}
		_, err := conn.Do("")
		//debug("data:", d)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPicWaitReview() []tPic {
	err := review()
	if err != nil {
		errorlog("move checked pics to reviewed list fail with error:", err)
	}
	pics := make([]tPic, 0)
	count := getPicNum(listUnchecked)
	if count == 0 {
		return pics
	}
	conn := rdb.Get()
	defer conn.Close()
	for i := 0; i < count; i++ {
		if i == numPerPage {
			break
		}
		t, err := redis.String(conn.Do("RPOPLPUSH", listUnchecked, listTemp))
		if err != nil {
			errorlog("move pics to temp list fail with error:", err, "and notviewed pics count:", count)
			return make([]tPic, 0)
		}
		tp := tPic{
			ID: func() int {
				if count > numPerPage {
					return numPerPage - i - 1
				}
				return count - i - 1
			}(),
			GirlPic: GirlPic{
				URL: t,
			},
		}
		pics = append(pics, tp)
	}
	return pics
}
