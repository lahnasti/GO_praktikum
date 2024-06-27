package config

import "flag"

type Conifg struct {
	Addr   string
	DBAddr string
}

// Функция обработки флагов запуска
func ReadConfig() Conifg {
	var addr string
	var dbAddr string
	flag.StringVar(&addr, "addr", ":8080", "Server address") // mani.exe -help
	flag.StringVar(&dbAddr, "db", "postgres://nastya:pgspgs@localhost:5433/pgs?sslmode=disable", "database connection addres")
	flag.Parse()
	return Conifg{
		Addr:   addr,
		DBAddr: dbAddr,
	}
}
