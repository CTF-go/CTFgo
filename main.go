package main

import (
	i "CTFgo/api/init"
)

func main() {
	r := i.SetupRouter()
	if err := r.Run(); err != nil {
		//fmt.Println("startup service failed, err:%v\n", err)
	}
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
