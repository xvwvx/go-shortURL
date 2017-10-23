package core

import (
	"flag"

	"github.com/Terry-Mao/goconf"
)

func init() {
	flag.StringVar(&confFile, "c", "./shortURL.conf", " set shortURL config file path")
}

var (
	Conf     *Config
	confFile string
)

type Config struct {
	RootURL string `goconf:"base:root_url"`
	Token   string `goconf:"base:token"`
	Procs   int    `goconf:"base:procs"`

	DBHost      string `goconf:"base:db.host"`
	DBUsername  string `goconf:"base:db.username"`
	DBPassword  string `goconf:"base:db.password"`
	DBDatabase  string `goconf:"base:db.database"`
	DBTablename string `goconf:"base:db.tablename"`
}

func LoadConfg() error {
	Conf = &Config{

		DBHost:      "127.0.0.1:3306",
		DBUsername:  "root",
		DBPassword:  "123456",
		DBTablename: "test",
	}

	gconf := goconf.New()
	if err := gconf.Parse(confFile); err != nil {
		return err
	}

	if err := gconf.Unmarshal(Conf); err != nil {
		return err
	}

	return nil
}
