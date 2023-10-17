package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type CountryInfo struct {
	Country string `json:"country"`
	Count   int    `json:"count"`
}

type Response struct {
	Status    string                 `json:"status"`
	Countries map[string]CountryInfo `json:"countries"`
}

func main() {
	color.Red(`
							
          _____                    _____            _____                    _____                    _____                    _____                    _____                    _____          
         /\    \                  /\    \          /\    \                  /\    \                  /\    \                  /\    \                  /\    \                  /\    \         
        /::\    \                /::\____\        /::\    \                /::\    \                /::\____\                /::\    \                /::\    \                /::\____\        
       /::::\    \              /:::/    /       /::::\    \              /::::\    \              /:::/    /               /::::\    \              /::::\    \              /::::|   |        
      /::::::\    \            /:::/    /       /::::::\    \            /::::::\    \            /:::/    /               /::::::\    \            /::::::\    \            /:::::|   |        
     /:::/\:::\    \          /:::/    /       /:::/\:::\    \          /:::/\:::\    \          /:::/    /               /:::/\:::\    \          /:::/\:::\    \          /::::::|   |        
    /:::/__\:::\    \        /:::/    /       /:::/__\:::\    \        /:::/__\:::\    \        /:::/____/               /:::/  \:::\    \        /:::/__\:::\    \        /:::/|::|   |        
   /::::\   \:::\    \      /:::/    /       /::::\   \:::\    \       \:::\   \:::\    \      /::::\    \              /:::/    \:::\    \      /::::\   \:::\    \      /:::/ |::|   |        
  /::::::\   \:::\    \    /:::/    /       /::::::\   \:::\    \    ___\:::\   \:::\    \    /::::::\    \   _____    /:::/    / \:::\    \    /::::::\   \:::\    \    /:::/  |::|___|______  
 /:::/\:::\   \:::\    \  /:::/    /       /:::/\:::\   \:::\    \  /\   \:::\   \:::\    \  /:::/\:::\    \ /\    \  /:::/    /   \:::\    \  /:::/\:::\   \:::\    \  /:::/   |::::::::\    \ 
/:::/  \:::\   \:::\____\/:::/____/       /:::/  \:::\   \:::\____\/::\   \:::\   \:::\____\/:::/  \:::\    /::\____\/:::/____/     \:::\____\/:::/  \:::\   \:::\____\/:::/    |:::::::::\____\
\::/    \:::\   \::/    /\:::\    \       \::/    \:::\  /:::/    /\:::\   \:::\   \::/    /\::/    \:::\  /:::/    /\:::\    \      \::/    /\::/    \:::\  /:::/    /\::/    / ~~~~~/:::/    /
 \/____/ \:::\   \/____/  \:::\    \       \/____/ \:::\/:::/    /  \:::\   \:::\   \/____/  \/____/ \:::\/:::/    /  \:::\    \      \/____/  \/____/ \:::\/:::/    /  \/____/      /:::/    / 
          \:::\    \       \:::\    \               \::::::/    /    \:::\   \:::\    \               \::::::/    /    \:::\    \                       \::::::/    /               /:::/    /  
           \:::\____\       \:::\    \               \::::/    /      \:::\   \:::\____\               \::::/    /      \:::\    \                       \::::/    /               /:::/    /   
            \::/    /        \:::\    \              /:::/    /        \:::\  /:::/    /               /:::/    /        \:::\    \                      /:::/    /               /:::/    /    
             \/____/          \:::\    \            /:::/    /          \:::\/:::/    /               /:::/    /          \:::\    \                    /:::/    /               /:::/    /     
                               \:::\    \          /:::/    /            \::::::/    /               /:::/    /            \:::\    \                  /:::/    /               /:::/    /      
                                \:::\____\        /:::/    /              \::::/    /               /:::/    /              \:::\____\                /:::/    /               /:::/    /       
                                 \::/    /        \::/    /                \::/    /                \::/    /                \::/    /                \::/    /                \::/    /        
                                  \/____/          \/____/                  \/____/                  \/____/                  \/____/                  \/____/                  \/____/         
                                                                                                                                                                                                
`)
	var codeList []string
	resp, err := http.Get("http://www.insecam.org/en/jsoncountries/")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	columnsPerRow := 6
	var i = 1
	var maxRowWidth int
	var rows []string

	for cc, country := range data.Countries {
		num := strconv.Itoa(i)
		row := fmt.Sprintf("%v: %v", num, country.Country)
		rows = append(rows, row)
		if len(row) > maxRowWidth {
			maxRowWidth = len(row)
		}
		codeList = append(codeList, cc)
		i++
	}

	for i, row := range rows {
		padding := strings.Repeat(" ", maxRowWidth-len(row))
		if i%columnsPerRow == 0 {
			fmt.Print(row + padding) 
		} else if i%columnsPerRow == 1 {
			fmt.Print(row)
		} else {
			fmt.Print(row + padding)
		}
		if i%columnsPerRow != 0 {
			fmt.Print("\t") 
		} else {
			fmt.Println() 
		}
	}

	color.HiBlue("\nSelect a country:")
	var countryCode string
	fmt.Scanln(&countryCode)
	num, e := strconv.Atoi(countryCode)
	if e != nil {
		fmt.Println(e)
		return
	}
	countryCode = codeList[num-1]

	resp, err = http.Get(fmt.Sprintf("http://www.insecam.org/en/bycountry/%s", countryCode))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var ips []string
	re := regexp.MustCompile(`http://\d+.\d+.\d+.\d+:\d+`)
	for _, match := range re.FindAllString(string(body), -1) {
		ips = append(ips, match)
	}

	fileName := fmt.Sprintf("%s.txt", countryCode)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for _, ip := range ips {
		fmt.Println(ip + "\n")
		file.WriteString(ip + "\n")
	}

	color.Green(fmt.Sprintf("IP addresses saved to file: %s", fileName))
}
