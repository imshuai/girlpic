package main

import "time"

//Comments define database struct for table comments
type Comments struct {
	ID         int64     `json:"id"          form:"id"`
	Auther     string    `json:"auther"      form:"auther"     binding:"required"`
	ReplyTo    int64     `json:"reply_to"    form:"reply_to"   binding:"required"`
	Content    string    `json:"content"     form:"content"    binding:"required"`
	CreateTime time.Time `json:"create_time" form:"create_time"`
}
