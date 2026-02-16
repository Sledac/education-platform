package main

import (
	"k8s-connect/internal/k8s"
	"log"
)

func main() {


	lab_mng,err := k8s.CreateLabManager("test-user", "count-pods")

	if err != nil {
		log.Fatal("Can't create LabManager")
	}

	cli_connect, err := lab_mng.CreateLabSession("test-user")

	if err != nil {
		log.Fatal("Can't create LabSession")
	}
	log.Print(cli_connect)
	err = lab_mng.Checker("test-user")
	err = lab_mng.DeleteNamespace("test-user")

}
