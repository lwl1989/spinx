package cmd

import "fmt"

func ShowHelp() {
	fmt.Println("Welcome use Ttpush!")
	fmt.Println("===================")
	fmt.Println("")
	fmt.Println("Usage: command install | remove | start | stop | status")
	fmt.Println("")
	fmt.Println("parmas:")
	fmt.Printf("\t%s\n", "-c <filepath>:")
	fmt.Printf("\t\t%s\n", "input config file path,content like this:")
	fmt.Printf("\t\t%s\n", `{
		  "server_port":"8080",
		  "max_ttl":2419200,
		  "api_key":"AAAA_1dLSps:APA91......ZHrCUioe-vx6wFvDXfnoh9h",
		  "notify_callback":"http://localhost:8000/fcm/notify",
		  "log_file":"/tmp/",
		  "proxy":"socket5 proxy url like 127.0.0.1:1080",
		  "notification":{
			"title":"",
			"body":"",
			"icon":"icon url like http://xxx.ico",
			"uri":"click_action like https://www.google.com or any schema://"
		  }
		}`)
	fmt.Println("")
	fmt.Printf("\t%s\n", "-h <show>:")
	fmt.Printf("\t\t%s\n", "help command and list commands")
	fmt.Println("")
	fmt.Printf("\t%s\n", "-d <bool>:")
	fmt.Printf("\t\t%s\n", "true or false set It's daemon?")
}
