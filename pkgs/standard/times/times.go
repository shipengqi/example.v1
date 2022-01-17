package main

import (
	"fmt"
	"time"
)

func main() {
	dateTime := time.Now()
	fmt.Println(dateTime) // 2022-01-17 11:43:58.0866166 +0800 CST m=+0.010161001

	year := time.Now().Year()
	fmt.Println(year) // 2022

	month := time.Now().Month()
	fmt.Println(month) // January

	day := time.Now().Day()
	fmt.Println(day) // 17

	hour := time.Now().Hour()
	fmt.Println(hour) // 11

	minute := time.Now().Minute()
	fmt.Println(minute) // 43

	second := time.Now().Second()
	fmt.Println(second) // 58

	nanosecond := time.Now().Nanosecond()
	fmt.Println(nanosecond) // 89864800

	// timestamp
	timeUnix := time.Now().Unix()
	timeUnixMilli := time.Now().UnixMilli()
	timeUnixMicro := time.Now().UnixMicro()
	timeUnixNano := time.Now().UnixNano()

	fmt.Println(timeUnix) // 1642391038
	fmt.Println(timeUnixMilli) // 1642391038089
	fmt.Println(timeUnixMicro) // 1642391038089864
	fmt.Println(timeUnixNano) // 1642391038089864800

	// format timestamp
	fmt.Println(time.Now().Format("2006-01-02 15:04:05")) // 2022-01-17 11:43:58

	// int64 to go timestamp
	var timeUnix2 int64 = 1562555859
	fmt.Println(time.Unix(timeUnix2,0)) // 2019-07-08 11:17:39 +0800 CST
	fmt.Println(time.Unix(timeUnix2, 0).Format("2006-01-02 15:04:05")) // 2019-07-08 11:17:39

	// string to timestamp
	t := time.Date(2014, 1, 7, 5, 50, 4, 0, time.Local).Unix()
	fmt.Println(t) // 1389045004

	calculateTime()
}

func calculateTime() {
	// 获取今天 0 点 0 时 0 分的时间戳
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	fmt.Println(startTime) // 2022-01-17 00:00:00 +0800 CST
	fmt.Println(startTime.Format("2006/01/02 15:04:05")) // 2022/01/17 00:00:00

	// 获取今天 23:59:59 秒的时间戳
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())
	fmt.Println(endTime) // 2022-01-17 23:59:59 +0800 CST
	fmt.Println(endTime.Format("2006/01/02 15:04:05")) // 2022/01/17 23:59:59


	// 获取 1 分钟之前的时间
	m, _ := time.ParseDuration("-1m")
	result := currentTime.Add(m)
	fmt.Println(result) // 2022-01-17 11:51:32.393872 +0800 CST m=-59.990519799
	fmt.Println(result.Format("2006/01/02 15:04:05")) // 2022/01/17 11:51:32

    // 获取 1 小时之前的时间
	m, _ = time.ParseDuration("-1h")
	result = currentTime.Add(m)
	fmt.Println(result) // 2022-01-17 10:52:32.393872 +0800 CST m=-3599.990519799
	fmt.Println(result.Format("2006/01/02 15:04:05")) // 2022/01/17 10:52:32

	// 获取 1 分钟之后的时间
	m, _ = time.ParseDuration("1m")
	result = currentTime.Add(m)
	fmt.Println(result) // 2022-01-17 11:53:32.393872 +0800 CST m=+60.009480201
	fmt.Println(result.Format("2006/01/02 15:04:05")) // 2022/01/17 11:53:32

	// 获取 1 小时之后的时间
	m, _ = time.ParseDuration("1h")
	result = currentTime.Add(m)
	fmt.Println(result) // 2022-01-17 12:52:32.393872 +0800 CST m=+3600.009480201
	fmt.Println(result.Format("2006/01/02 15:04:05")) // 2022/01/17 12:52:32


	// 计算两个时间戳间隔时间
	afterTime, _ := time.ParseDuration("1h")
	result = currentTime.Add(afterTime)

	beforeTime, _ := time.ParseDuration("-1h")
	result2 := currentTime.Add(beforeTime)

	m = result.Sub(result2)
	fmt.Printf("%v 分钟 \n", m.Minutes()) // 120 分钟
	fmt.Printf("%v 小时 \n", m.Hours()) // 2 小时
	fmt.Printf("%v 天\n", m.Hours()/24) // 0.08333333333333333 天

	// 判断一个时间是否在一个时间之后
	stringTime, _ := time.Parse("2006-01-02 15:04:05", "2019-12-12 12:00:00")
	beforeOrAfter := stringTime.After(time.Now())

	// 2019-12-12 12:00:00 在当前时间之前!
	if beforeOrAfter {
		fmt.Println("2019-12-12 12:00:00 在当前时间之后!")
	} else {
		fmt.Println("2019-12-12 12:00:00 在当前时间之前!")
	}

	// 判断一个时间相比另外一个时间过去了多久
	startTime = time.Now()
	time.Sleep(time.Second * 5)

	fmt.Println("离现在过去了：", time.Since(startTime)) // 离现在过去了： 5.0081632s
}
