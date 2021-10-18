package config

type Config struct {
	ConnString  string
	AllowOrigin []string
}

var DefaultConfig = &Config{
	ConnString:  "host=localhost user=postgres password=admin dbname=ecom port=5432 sslmode=disable TimeZone=Asia/Shanghai",
	AllowOrigin: []string{"http://localhost:3000"},
}
