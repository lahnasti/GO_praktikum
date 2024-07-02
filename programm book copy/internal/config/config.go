package config

import (
	"flag"
	"os"
)

type Conifg struct {
	Addr   string
	DBAddr string
}

const (
	defaultAddr = ":8080"
	defaultDbDSN = "postgres://nastya:pgspgs@localhost:5455/books"
)

// Функция обработки флагов запуска
func ReadConfig() Conifg {
	var addr string
	var dbAddr string
	flag.StringVar(&addr, "addr", defaultAddr, "Server address") // mani.exe -help
	flag.StringVar(&dbAddr, "db", defaultDbDSN, "database connection addres")
	flag.Parse()

	
	if temp := os.Getenv("SERVER_ADDR"); temp != "" {
		if addr == defaultAddr {
		addr = temp
		}
	}
	if temp := os.Getenv("DB_DSN"); temp != "" {
		if dbAddr == defaultDbDSN {
		dbAddr = temp
		}
	}


	return Conifg{
		Addr:   addr,
		DBAddr: dbAddr,
	}
}
