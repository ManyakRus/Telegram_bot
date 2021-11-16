package main

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/xuri/excelize/v2"
	"gopkg.in/ini.v1"
)

type settingsINI struct {
	api_key           string
	filename_telefon  string
	column_name       int
	column_phone      int
	column_post       int
	column_cell_phone int
	column_adress     int
	column_email      int
	first_line_number int
}

type User struct {
	FIO        string
	Phone      string
	Cell_phone string
	Post       string
	Adress     string
	Email      string
}

var settings settingsINI
var MassUsers []User

func main() {
	LoadINI()

	LoadExcel()

	StartBot()

}

func LoadExcel() {
	MassUsers = make([]User, 0)

	Excel, err := excelize.OpenFile(settings.filename_telefon)

	if err != nil {
		log.Fatal(err)
	}

	//defer func() {
	//	if err = Excel.; err != nil {
	//		fmt.Println(err)
	//	}
	//}()

	//for sheet0,_ range Excel.Sheet {
	//	sheet1 := Excel.Sheet
	//}

	sheet1 := Excel.WorkBook.Sheets.Sheet[0].Name

	rows, err := Excel.GetRows(sheet1)
	for y, row := range rows {
		if y < (settings.first_line_number - 1) {
			continue
		}

		newUser := User{}
		for x, value := range row {
			if x == (settings.column_name - 1) {
				newUser.FIO = value
			}

			if x == (settings.column_phone - 1) {
				newUser.Phone = value
			}

			if x == (settings.column_cell_phone - 1) {
				newUser.Cell_phone = value
			}

			if x == (settings.column_post - 1) {
				newUser.Post = value
			}

			if x == (settings.column_adress - 1) {
				newUser.Adress = value
			}

			if x == (settings.column_email - 1) {
				newUser.Email = value
			}

			//fmt.PrintLn(value, "\t", y,x)
		}
		MassUsers = append(MassUsers, newUser)
	}

	//c1, err := f.GetCellValue("Sheet1", "A1")
	//
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Total users: " + strconv.Itoa(len(MassUsers)))
}

func StartBot() {
	bot, err := tgbotapi.NewBotAPI(settings.api_key)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	var Users []User
	Users = make([]User, 0)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		TextFromUser := update.Message.Text

		Text := "Not understand. Fill name or phone number. "
		if valid.IsInt(TextFromUser) == true {
			Users = FindUser_by_Phone(TextFromUser)
			if len(Users) == 0 {
				Users = FindUser_by_CellPhone(TextFromUser)
			}
		} else if HaveAt(TextFromUser) == true {
			Users = FindUser_by_Email(TextFromUser)
		} else if valid.IsInt(TextFromUser) == false && HaveNumbers(TextFromUser) == false {
			Users = FindUser_by_Name(TextFromUser)
			if len(Users) == 0 {
				Users = FindUser_by_Post(TextFromUser)
			}
		} else {
			Users = FindUser_by_Adress(TextFromUser)
		}

		if len(Users) > 0 {
			Text = ""
			for _, User1 := range Users {
				Text = Text + User1.String() + "\n"
			}
		}

		if len(Text) > 2000 {
			Text = Text[0:2000]
			Text = Text + "\n" + "..."
			Text = Text + "\n" + "..."
			Text = Text + "\n" + "..."
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, Text)
		msg.ReplyToMessageID = update.Message.MessageID
		//msg.Entities = append(msg.Entities, )
		//append(msg.Entities, )

		bot.Send(msg)
	}

}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func LoadINI() {
	settings = settingsINI{}

	ProgramDir := ProgramDir()
	filename := "SettingsMain.txt"
	flagOK, _ := FileExists(filename)
	if flagOK == false {
		filename = "SettingsM.txt"
	}
	cfg, err := ini.Load(ProgramDir + filename)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	SectionMain := cfg.Section("Main")

	settings.api_key = SectionMain.Key("api_key").String()

	filename = "TelefonMain.xlsx"
	flagOK, _ = FileExists(filename)
	if flagOK == false {
		filename = SectionMain.Key("filename_telefon").String()
	}
	settings.filename_telefon = filename

	settings.column_name, err = SectionMain.Key("column_name").Int()
	if err != nil {
		fmt.Println("Fail to read column_name !")
		os.Exit(1)
	}

	settings.column_phone, err = SectionMain.Key("column_phone").Int()
	if err != nil {
		fmt.Println("Fail to read column_phone !")
		os.Exit(1)
	}

	settings.column_post, err = SectionMain.Key("column_post").Int()
	if err != nil {
		fmt.Println("Fail to read column_post !")
		os.Exit(1)
	}

	settings.column_cell_phone, err = SectionMain.Key("column_cell_phone").Int()
	if err != nil {
		fmt.Println("Fail to read column_cell_phone !")
		os.Exit(1)
	}

	settings.column_adress, err = SectionMain.Key("column_adress").Int()
	if err != nil {
		fmt.Println("Fail to read column_adress !")
		os.Exit(1)
	}

	settings.first_line_number, err = SectionMain.Key("first_line_number").Int()
	if err != nil {
		fmt.Println("Fail to read first_line_number !")
		os.Exit(1)
	}

	settings.column_email, err = SectionMain.Key("column_email").Int()
	if err != nil {
		fmt.Println("Fail to read column_email !")
		os.Exit(1)
	}
}

func ProgramDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		dir = ""
	}
	dir = dir + "\\"
	return dir
}

func (u User) String() string {
	Otvet := ""

	Otvet = Otvet +
		"Name: " + u.FIO + "\n" +
		"Phone: " + u.Phone + "\n" +
		"Cell phone: " + u.Cell_phone + "\n" +
		//		"Cell phone: " + "<a href=\"tel:" + u.Cell_phone + "\">" + "\n" +
		"Post: " + u.Post + "\n" +
		"Adress: " + u.Adress + "\n" +
		"Email: " + u.Email + "\n"

	return Otvet
}

func FindUser_by_Phone(Phone string) []User {
	var Otvet []User
	Otvet = make([]User, 0)

	for _, User1 := range MassUsers {
		if strings.Contains(User1.Phone, Phone) == true {
			Otvet = append(Otvet, User1)
			//break
		}
	}

	return Otvet
}

func FindUser_by_CellPhone(Phone string) []User {
	var Otvet []User
	Otvet = make([]User, 0)

	for _, User1 := range MassUsers {
		if strings.Contains(User1.Cell_phone, Phone) == true {
			Otvet = append(Otvet, User1)
			//break
		}
	}

	return Otvet
}

func FindUser_by_Name(Name string) []User {
	var Otvet []User
	Otvet = make([]User, 0)

	for _, User1 := range MassUsers {
		if strings.Contains(strings.ToLower(User1.FIO), strings.ToLower(Name)) == true {
			Otvet = append(Otvet, User1)
			//break
		}
	}

	return Otvet
}

func FindUser_by_Post(Name string) []User {
	var Otvet []User
	Otvet = make([]User, 0)

	for _, User1 := range MassUsers {
		if strings.Contains(strings.ToLower(User1.Post), strings.ToLower(Name)) == true {
			Otvet = append(Otvet, User1)
			//break
		}
	}

	return Otvet
}

func FindUser_by_Adress(Name string) []User {
	var Otvet []User
	Otvet = make([]User, 0)

	Name = DeleteNumbers(Name)
	Name = strings.Trim(Name, " ")

	for _, User1 := range MassUsers {
		if strings.Contains(strings.ToLower(User1.Adress), strings.ToLower(Name)) == true {
			Otvet = append(Otvet, User1)
			//break
		}
	}

	return Otvet
}

func FindUser_by_Email(Name string) []User {
	var Otvet []User
	Otvet = make([]User, 0)

	Name = strings.Trim(Name, " ")
	Name = strings.ReplaceAll(Name, "@", "")

	for _, User1 := range MassUsers {
		if strings.Contains(strings.ToLower(User1.Email), strings.ToLower(Name)) == true {
			Otvet = append(Otvet, User1)
			//break
		}
	}

	return Otvet
}

func HaveNumbers(s string) bool {
	Otvet := false

	var s1 string
	for _, s0 := range s {
		s1 = string(s0)
		if (s1 == "0" || s1 == "1" || s1 == "2" || s1 == "3" || s1 == "4" || s1 == "5" || s1 == "6" || s1 == "7" || s1 == "8" || s1 == "9") == true {
			Otvet = true
			break
		}

	}

	return Otvet
}

func HaveAt(s string) bool {
	Otvet := false

	var s1 string
	for _, s0 := range s {
		s1 = string(s0)
		if s1 == "@" {
			Otvet = true
			break
		}

	}

	return Otvet
}

func DeleteNumbers(s string) string {
	Otvet := ""

	var s1 string
	for _, s0 := range s {
		s1 = string(s0)
		if (s1 == "0" || s1 == "1" || s1 == "2" || s1 == "3" || s1 == "4" || s1 == "5" || s1 == "6" || s1 == "7" || s1 == "8" || s1 == "9") == false {
			Otvet = Otvet + s1
		}

	}

	return Otvet
}
