package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Person struct {
	FirstName   string
	LastName    string
	City        string
	Address     string
	Email       string
	PhoneNumber int64
}

func displayusage() string {
	return string("usage: ccli [options] <command> <subcommand(s)> [parameters]\nFor help try:\n\nccli help\nccli <command> help\nccli <command> <subcommand> help")
}

func displayErrorUsage(s string) string {
	return string("usage: of command [" + s + "] is not supported\nFor help try:\n\nccli help\nccli <command> help\nccli <command> <subcommand> help")
}

func convertJSONToContact(json []Person, p *[]Person) {
	*p = json
}

// func encodeTemplate(template *Person, p Person)

func loadFile(fname string) []Person {
	jsonFile, err := os.Open(fname)

	//	check if opening file throws an error
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer jsonFile.Close()
	var persons []Person
	rawBytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(rawBytes), &persons)

	return persons
}

func save(contacts []Person, savePath string) {
	// anon function
	// person := func(f string, l string, c string, a string, e string, p int64) *Person {
	// 	return &Person{
	// 		FirstName:   f,
	// 		LastName:    l,
	// 		City:        c,
	// 		Address:     a,
	// 		Email:       e,
	// 		PhoneNumber: p,
	// 	}
	// }

	// scontact := []Person

	// for _, p := range contacts {

	// }
	data, _ := json.Marshal(contacts)

	_ = ioutil.WriteFile(savePath, data, 0644)
}

func fileExists(path string) bool {
	src, err := os.Stat(path)
	if os.IsNotExist(err) && src == nil {
		return false
	}
	return true
}

func dirExists(dname string, homedir string) bool {
	path := homedir + "/" + dname
	fmt.Println("Path:", path)
	src, err := os.Stat(path)

	if os.IsNotExist(err) && src.Name() != dname {
		// create new dir
		newDir := os.MkdirAll(path, 0755)
		if newDir != nil {
			panic(err)
			return false
		}
		return true
	}
	return true
}

func main() {
	fileName := "contacts.json"
	dataDirName := ".ccli/"
	homePath := os.Getenv("HOME")
	// holds contact info
	var contactList []Person

	// checks if home path exists
	if homePath == "" {
		log.Fatal("Error, no envvar for $HOME")
		os.Exit(1)
	}

	fmt.Printf("homePath: %s\n", homePath)
	fmt.Println("dirExists:", dirExists(dataDirName, string(homePath)))

	if dirExists(dataDirName, string(homePath)) {
		// check for data file to load
		pathToFile := homePath + "/" + dataDirName + fileName
		fmt.Println("fileExists:", fileExists(pathToFile))
		if fileExists(pathToFile) {
			data := loadFile(pathToFile)
			if data != nil {
				convertJSONToContact(data, &contactList)
			}
			fmt.Println("contactList:", contactList)

		}
	}

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
	listAllPtr := flag.String("la", "0", "-la - lists all")

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
		// verify parsed content can execute
		if *listAllPtr == "0" {
			flag.Usage()
			os.Exit(1)
		}

	}
	if helpCommand.Parsed() {
		fmt.Println(*helpCommand)
	}
	// check for parsed commands and it's tails
	// debug
	// fmt.Printf("createFirstNamePtr: %s\ncreateLastNamePtr: %s\ncreateCityPtr: %s\ncreateAddrPtr: %s\ncreateEmailPtr: %s\ncreatePhonePtr: %s\n", *createFirstNamePtr, *createLastNamePtr, *createCityPtr, *createAddrPtr, *createEmailPtr, *createPhonePtr)

	savePath := homePath + "/" + dataDirName + fileName
	// TODO: Save the changes (if any back to the data structure)
	save(contactList, savePath)

	return

}
