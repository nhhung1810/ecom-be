package config

type Config struct {
	ConnString  string
	AllowOrigin []string
	OrderStatus map[int]string
	SecretKey   string
}

var DefaultConfig = &Config{
	ConnString:  "host=localhost user=postgres password=admin dbname=ecom port=5432 sslmode=disable TimeZone=Asia/Shanghai",
	AllowOrigin: []string{"http://localhost:3000"},
	OrderStatus: map[int]string{
		0: "Pending",
		1: "Completed",
		2: "Cancel",
	},
	SecretKey: "EcaLf2vYAe1GtT369eD6jtfxA0iXC6HlPj1meCE/oro=",
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

var BrandArray = []string{
	"zara",
	"hm",
	"ninomax",
}
