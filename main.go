package main

import "context"
import "flag"
import "fmt"
import "log"
import "net/url"
import "time"

import "google.golang.org/api/sheets/v4"
import "google.golang.org/api/option"

func main() {

	username := flag.String("username", "", "Peloton email or username")
	password := flag.String("password", "", "Peloton password")

	flag.Parse()

	if *username == "" || *password == "" {
		log.Println("A username and password input is required")
		log.Fatalf("Use format 'hcotf-cli -username YOURUSERNAME -password YOURPASSWORD'")
	}

	currentTime := time.Now()
	currentDay := currentTime.Format("1/02/2006")

	fmt.Printf("Today is: %s\n", currentDay)

	ps := NewPelotonSession(*username, *password)
	pt := NewPelotonStack(ps.Session)
	//ps.GetClass(m["classId"][0])
	//pt.GetStack()

	apiKey, err := Asset("data/config/api.key")
	if err != nil {
		log.Fatalf("Unable to get api key", err)
	}

	sheetId, err := Asset("data/config/sheet.id")
	if err != nil {
		log.Fatalf("Unable to get Google Sheet ID", err)
	}

	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithAPIKey(string(apiKey)))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := string(sheetId)
	readRange := "MONTHLY_LIST_VIEW!C8:I200"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		log.Fatal("No data found.")
	}

	fmt.Println()

	fmt.Println("Today's classes:")

	curClassNames := []string{}
	curClassIds := []string{}

	for _, row := range resp.Values {
		if row[0] == currentDay && len(row) >= 7 {
			m, err := url.ParseQuery(row[6].(string))
			if err != nil {
				log.Printf("Unable to parse row from today: %s - %s", row[6].(string), err)
			} else {
				fmt.Printf("\t %s %s\n", row[1], m["classId"][0])
				curClassNames = append(curClassNames, row[1].(string))
				curClassIds = append(curClassIds, m["classId"][0])
			}
		}
	}

	if len(curClassNames) == 0 {
		log.Print("Unable to find any classes for today. Exiting without making any changes")
		return
	}

	fmt.Println()

	log.Println("Clear current stack")
	pt.ClearStack()

	for i := range curClassNames {
		log.Printf("Add %s to stack\n", curClassNames[i])
		pClass, err := ps.GetClass(curClassIds[i])
		if err != nil {
			log.Fatalf("Unable to find class: %s", curClassIds[i])
		}
		pt.AddClassToStack(pClass.JoinTokens.OnDemand)
	}
}
