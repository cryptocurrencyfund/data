package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/cryptocurrencyfund/data/util"
)

func usage() {
	fmt.Println("Get variable names")
	fmt.Println("./data [number of seconds]")
	er := errors.New("Wrong arguments")
	fmt.Println(er)
}

func main() {
	db := util.OpenDb()
	defer util.CloseDb(db)

	seconds := 5
	if len(os.Args) > 0 {
		seconds, _ = strconv.Atoi(os.Args[1])
	}
	fmt.Printf("Getting price every %d seconds\n", seconds)

	top100 := util.FetchPrices(db)
	dateString := util.DateString()
	util.SaveJSONToFile(dateString, top100)
	util.SyncGit(dateString)
}
