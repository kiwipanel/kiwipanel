package controllers

import "fmt"

func ChangePassCode() {

}

func ChangePassWord() {

}

func ChangeTime() {
	//Future work: https://thuanbui.me/timezone-linux-crontab/
	command := "systemd-timesyncd"
	fmt.Println(command)
}
