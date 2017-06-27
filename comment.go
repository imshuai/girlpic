package main

import "time"

//Comments define database struct for table comments
type Comments struct {
	ID         int64     `json:"id"                                   form:"id"`
	Auther     string    `xorm:"varchar(20)"       json:"auther"      form:"auther"     binding:"required"`
	ReplyTo    int64     `xorm:"notnull default 0" json:"reply_to"    form:"reply_to"   binding:"required"`
	Content    string    `xorm:"text"              json:"content"     form:"content"    binding:"required"`
	CreateTime time.Time `xorm:"created"           json:"create_time" form:"create_time"`
}
