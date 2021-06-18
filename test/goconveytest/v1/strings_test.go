package v1

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	success                = ""
	needExactValues        = "This assertion requires exactly %d comparison values (you provided %d)."
	needNonEmptyCollection = "This assertion requires at least 1 comparison value (you provided 0)."
)

func TestStringSliceEqual(t *testing.T) {
	Convey("TestStringSliceEqual should return true when a != nil  && b != nil", t, func() {
		a := []string{"hello", "goconvey"}
		b := []string{"hello", "goconvey"}
		So(StringSliceEqual(a, b), ShouldBeTrue)
	})

	Convey("TestStringSliceEqual should return true when a ＝= nil  && b ＝= nil", t, func() {
		So(StringSliceEqual(nil, nil), ShouldBeTrue)
	})

	Convey("TestStringSliceEqual should return false when a ＝= nil  && b != nil", t, func() {
		a := []string(nil)
		b := []string{}
		So(StringSliceEqual(a, b), ShouldBeFalse)
	})

	Convey("TestStringSliceEqual should return false when a != nil  && b != nil", t, func() {
		a := []string{"hello", "world"}
		b := []string{"hello", "goconvey"}
		So(StringSliceEqual(a, b), ShouldBeFalse)
	})
}

func TestStringSliceEqualNested(t *testing.T) {
	Convey("TestStringSliceEqualNested", t, func() {
		Convey("should return true when a != nil  && b != nil", func() {
			a := []string{"hello", "goconvey"}
			b := []string{"hello", "goconvey"}
			So(StringSliceEqual(a, b), ShouldBeTrue)
		})

		Convey("should return true when a ＝= nil  && b ＝= nil", func() {
			So(StringSliceEqual(nil, nil), ShouldBeTrue)
		})

		Convey("should return false when a ＝= nil  && b != nil", func() {
			a := []string(nil)
			b := []string{}
			So(StringSliceEqual(a, b), ShouldBeFalse)
		})

		Convey("should return false when a != nil  && b != nil", func() {
			a := []string{"hello", "world"}
			b := []string{"hello", "goconvey"}
			So(StringSliceEqual(a, b), ShouldBeFalse)
		})
	})
}

// 定制断言函数
// So 的函数原型
// func So(actual interface{}, assert assertion, expected ...interface{})
// assertion 原型 返回值为 "" 时表示断言成功，否则表示失败
// type assertion func(actual interface{}, expected ...interface{}) string

func ShouldSummerBeComing(actual interface{}, expected ...interface{}) string {
	if actual == "summer" && expected[0] == "coming" {
		return ""
	} else {
		return "summer is not coming!"
	}
}

// SkipConvey 函数表明相应的闭包函数将不被执行
// SkipSo 函数表明相应的断言将不被执行
func TestSummer(t *testing.T) {
	SkipConvey("TestSummer", t, func() {
		So("summer", ShouldSummerBeComing, "coming")
		SkipSo("winter", ShouldSummerBeComing, "coming")
	})
}
