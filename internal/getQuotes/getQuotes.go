package getQuotes

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"
	"os"
	"time"
)

type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Valutes []Valute `xml:"Valute"`
}

func GetQuotes() {
	if len(os.Args) != 3 {
		fmt.Println("Должно быть 2 аргумента: code, date")
		return
	}

	code := os.Args[1]
	inputDate, err := time.Parse("2006-01-02", os.Args[2])
	if err != nil {
		fmt.Println("Ошибка преобразования даты:", err)
		return
	}
	date := inputDate.Format("02/01/2006")
	url := fmt.Sprintf("https://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Ошибка создания запроса: %s\n", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка выполнения GET-запроса: %s\n", err)
		return
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return input, nil
	}

	valCurs := &ValCurs{}
	err = decoder.Decode(valCurs)
	if err != nil {
		fmt.Printf("Ошибка при разборе XML: %s\n", err)
		return
	}

	for _, valute := range valCurs.Valutes {
		if valute.CharCode == code {
			fmt.Printf("%s (%s): %s\n", valute.CharCode, valute.Name, valute.Value)
			return
		}
	}

	fmt.Printf("Валюта с кодом %s не найдена\n", code)
}
