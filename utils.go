package goboleto

import (
	"time"
	"fmt"
)

func Date2_html(t time.Time) (ret string) {
	if t.IsZero() {
		ret = ""
	} else {
		ret = fmt.Sprintf("%04d-%02d-%02d",
			t.Year(),
			t.Month(),
			t.Day())
	}
	return ret
}

func Date2_str_br(t time.Time) (ret string) {
	//a := "2015-11-27"
	//ret := fmt.Sprintf("%2s/%2s/%4s", a[8:10], a[5:7], a[0:4])
	if t.IsZero() {
		ret = ""
	} else {
		ret = fmt.Sprintf("%02d/%02d/%04d",
			t.Day(),
			t.Month(),
			t.Year())
	}
	return ret
}

func Hoje() (ret time.Time) {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	hoje := Date2_html(time.Now())

	ret, err := time.ParseInLocation("2006-01-02 15:04", hoje+" 00:00", loc)

	if err != nil {
		fmt.Println(err.Error())
	}

	return
}


func Str2Date(data string) (ret time.Time) {
	loc, _ := time.LoadLocation("America/Sao_Paulo")

	ret, err := time.ParseInLocation("2006-01-02", data, loc)

	if err != nil {
		fmt.Println(err.Error())
	}

	return
}


