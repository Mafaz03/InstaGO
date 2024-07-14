package main

import (
	"fmt"
	"github.com/Mafaz03/InstaGO"
)

func main() {
	pic, err := GetProfilePicture("nigg.pablo")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(pic)

}
