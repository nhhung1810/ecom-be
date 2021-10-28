package config

type Config struct {
	ConnString  string
	AllowOrigin []string
	OrderStatus map[int]string
}

var DefaultConfig = &Config{
	ConnString:  "host=localhost user=postgres password=admin dbname=ecom port=5432 sslmode=disable TimeZone=Asia/Shanghai",
	AllowOrigin: []string{"http://localhost:3000"},
	OrderStatus: map[int]string{
		0: "Pending",
		1: "Completed",
		2: "Cancel",
	},
}

var SizeArray = []string{
	"s",
	"m",
	"l",
	"xl",
	"xll",
}

var ColorArray = []string{
	"red",
	"yellow",
	"blue",
	"orange",
	"grey",
	"silver",
}
