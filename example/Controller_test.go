package example

import (
	"github.com/zhuxiujia/easy_mvc"
	"net/http"
	"os"
	"testing"
	"time"
)

type TestUserVO struct {
	Name string
}

type TestController struct {
	easy_mvc.Controller `path:"/api"`
	//登录接口案例,返回值默认转json，如果要返回其他东西，请在参数里加上 request *http.Request 把content-type 改了，然后可以自行处理（或者直接兼容标准库func(writer http.ResponseWriter, request *http.Request)）
	Login func(phone string, pwd string, age *int) interface{} `path:"/login" arg:"phone,pwd,age" note:"phone:手机号,pwd:密码,age:年龄"`
	//兼容go标准库http案例,可以无返回值
	Login2 func(writer http.ResponseWriter, request *http.Request)             `path:"/login2" arg:"w,r"`
	Login3 func(writer http.ResponseWriter, request *http.Request) interface{} `path:"/login3" arg:"w,r"`
	Login4 func(phone string, pwd string, request *http.Request) interface{}   `path:"/login4" arg:"phone,pwd,r"`

	UserInfo func() interface{} `path:"/api/login2"`
}

func (it TestController) New() TestController {
	it.Login = func(phone string, pwd string, age *int) interface{} {
		println("do Login")
		return TestUserVO{}
	}
	it.UserInfo = func() interface{} {
		return TestUserVO{}
	}
	it.Login2 = func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("fuck"))
	}
	it.Init(&it)
	return it
}

func TestController_Init(t *testing.T) {
	//自定义一个全局错误处理器 到调用链中
	easy_mvc.RegisterGlobalErrorHandleChan(&easy_mvc.HttpErrorHandle{
		Func: func(err interface{}, w http.ResponseWriter, r *http.Request) {
			if err != nil {
				println(err.(error).Error())
				w.Write([]byte(err.(error).Error()))
			}
		},
		Name: "custom",
	})

	//初始化 控制器
	TestController{}.New()

	//你也可以使用标准库的api（使用标准库不经过easy_mvc）
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("yes"))
	})

	go func() {
		var exitTimeSec = 30
		for i := 0; i < exitTimeSec; i++ {
			time.Sleep(time.Second)
			println("测试倒计时：", exitTimeSec-i, "秒后退出！")
		}
		os.Exit(0) //30秒后退出测试，
	}()

	//使用标准库启动http
	http.ListenAndServe("127.0.0.1:8080", nil)
}
