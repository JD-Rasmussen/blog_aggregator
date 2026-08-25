package main

import (
	"fmt"

	"github.com/jd-rasmussen/blog_aggregator/internal/config"
)

func main() {
	//fmt.Println("Hello, World!")
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	//fmt.Println("setting user to jdr")
	cfg.SetUser("jdr")

	cfg, err = config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	fmt.Println(cfg.Db_url, "\n", cfg.Current_user_name)

}
