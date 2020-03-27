package main

import (
	"flag"
	"fmt"
	"os"
)

type Person struct {
	firstName   string
	lastName    string
	city        string
	address     string
	email       string
	phoneNumber int64
}

func displayusage() string {
	return string("usage: ccli [options] <command> <subcommand(s)> [parameters]\nFor help try:\n\nccli help\nccli <command> help\nccli <command> <subcommand> help")
}

func displayErrorUsage(s string) string {
	return string("usage: of command [" + s + "] is not supported\nFor help try:\n\nccli help\nccli <command> help\nccli <command> <subcommand> help")
}




func main() {

	// TODO: Create/Find local storage path to store contact info

	// pseudo
	// If $CCLI_DATA_PATH exists
	//		*check if contacts.json|csv|txt 
	//		if contacts.json|csv|txt exists
	//			1. data := load(contacts.json|csv|txt)
	//			2. convert(data, contactList)
	// else
	//		create $CCLI_DATA_PATH




	// holds contact info
	contactList := []Person

	// sub commands for cli
	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)

	//flags for create command
	createFirstNamePtr := createCommand.String("fname", "0", "-fname - field value for first name (REQUIRED])")
	createLastNamePtr := createCommand.String("lname", "0", "-lname - sets value for last name (REQUIRED)")
	createCityPtr := createCommand.String("city", "0", "-c - sets value for city name")
	createAddrPtr := createCommand.String("addr", "0", "-addr - sets value for address")
	createEmailPtr := createCommand.String("email", "0", "-email - sets value for email (REQUIRED)")
	createPhonePtr := createCommand.Int64("pnum", -1, "-pnum - sets value for phoneNumber (OPTIONAL)")

	//flags for create command
	// listAllPtr := flag.String("la", "0", "-la - lists all")

	// exists if sub command wasn't sepcificed
	if len(os.Args) < 2 {
		fmt.Println(displayusage())
		os.Exit(1)
	}

	flag.Parse()

	//switch statement to for what command and sub commands to run
	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[2:])
	case "create":
		createCommand.Parse(os.Args[2:])
	case "help":
		helpCommand.Parse(os.Args[2:])
	default:
		fmt.Println(displayErrorUsage(string(os.Args[1])))
		os.Exit(1)
	}

	if createCommand.Parsed() {
		// verify parsed content can execute
		if *createFirstNamePtr == "0" {
			flag.Usage()
			os.Exit(1)
		}
		if *createLastNamePtr == "0" {
			flag.Usage()
			os.Exit(1)
		}
		if *createEmailPtr == "0" {
			flag.Usage()
			os.Exit(1)
		}
		if *createPhonePtr == -1 {
			flag.Usage()
			os.Exit(1)
		}

		p := Person{*createFirstNamePtr, *createLastNamePtr, *createCityPtr, *createAddrPtr, *createEmailPtr, *createPhonePtr}

		contactList = append(contactList, p)

	}
	if listCommand.Parsed() {
		fmt.Println(*listCommand)
	}
	if helpCommand.Parsed() {
		fmt.Println(*helpCommand)
	}
	// check for parsed commands and it's tails
	// debug
	fmt.Printf("createFirstNamePtr: %s\ncreateLastNamePtr: %s\ncreateCityPtr: %s\ncreateAddrPtr: %s\ncreateEmailPtr: %s\ncreatePhonePtr: %s\n", *createFirstNamePtr, *createLastNamePtr, *createCityPtr, *createAddrPtr, *createEmailPtr, *createPhonePtr)


	// TODO: Save the changes (if any back to the data structure)
	
	return

}
