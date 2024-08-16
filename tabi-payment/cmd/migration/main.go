package main

import "tabi-payment/internal/functions/migration"

func main() {
	checkErr(migration.Run())
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
