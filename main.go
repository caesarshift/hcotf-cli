package main

import "context"
import "fmt"
import "log"
import "net/url"
import "os"
import "time"

import "github.com/levigross/grequests"
import "google.golang.org/api/sheets/v4"
import "google.golang.org/api/option"

/*
"C:\git\bin\go-bindata.exe" -o data.go data/...
"C:\git\bin\go-bindata.exe" -debug -o data.go data/...
*/

type GetClassDetailsJSON struct {
	Title      string                     `json:"title"`
	JoinTokens GetClassRideJoinTokensJSON `json:"join_tokens"`
}

type GetClassRideJoinTokensJSON struct {
	OnDemand string `json:"on_demand"`
}

type ViewUserStackResponseJSON struct {
	ViewUserStackData ViewUserStackDataJSON `json:"data"`
}

type ViewUserStackDataJSON struct {
	ViewUserStack ViewUserStackJSON `json:"ViewUserStack"`
}

type ViewUserStackJSON struct {
	NumClasses int `json:"numClasses"`
}

type PelotonStack struct {
	Endpoint string
	Headers  map[string]string
	Session  *grequests.Session
}

func NewPelotonStack(session *grequests.Session) *PelotonStack {
	ps := PelotonStack{
		Endpoint: "https://gql-graphql-gateway.prod.k8s.onepeloton.com/graphql",
		Headers: map[string]string{
			"Content-Type":     "application/json",
			"User-Agent":       "peloton-scheduler",
			"Peloton-Platform": "web",
		},
		Session: session,
	}
	return &ps
}

type PelotonSession struct {
	DefaultBaseUrl string
	DefaultHeaders map[string]string
	Session        *grequests.Session
}

func NewPelotonSession(username string, password string) *PelotonSession {
	ps := PelotonSession{
		DefaultBaseUrl: "https://api.onepeloton.com",
		DefaultHeaders: map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   "peloton-scheduler",
		},
		Session: grequests.NewSession(nil),
	}

	authUrl := fmt.Sprintf("%s/%s", ps.DefaultBaseUrl, "auth/login")
	log.Printf("Authenticating to %s", authUrl)

	resp, _ := ps.Session.Post(
		authUrl,
		&grequests.RequestOptions{
			JSON: map[string]string{
				"username_or_email": username,
				"password":          password,
			},
			Headers: ps.DefaultHeaders,
		},
	)

	if resp.Error != nil {
		fmt.Printf("Unable to make request", resp.Error)
	}
	if resp.Ok {
		log.Printf("Successfully created Peloton Session: %d", resp.StatusCode)
	} else {
		if resp.StatusCode == 401 {
			fmt.Printf("Invalid username or password: %d - %s", resp.StatusCode, "")
		} else {
			fmt.Printf("Unable to login and create a Peloton session: %d - %s", resp.StatusCode, "")
		}
	}

	return &ps
}

func (ps *PelotonSession) GetClass(classId string) (GetClassDetailsJSON, error) {
	ro := &grequests.RequestOptions{UserAgent: "web", DisableCompression: false}
	resp, _ := ps.Session.Get(fmt.Sprintf("%s/api/ride/%s", ps.DefaultBaseUrl, classId), ro)

	if !resp.Ok {
		if resp.StatusCode == 401 {
			return GetClassDetailsJSON{}, fmt.Errorf("Invalid username or password: %d - %s", resp.StatusCode, "")
		} else {
			return GetClassDetailsJSON{}, fmt.Errorf("Unable to login and create a Peloton session: %d - %s", resp.StatusCode, resp.String())
		}
	}

	// log.Printf("Found join token: %s", resp.String())

	classjson := &GetClassDetailsJSON{}

	err := resp.JSON(classjson)
	if err != nil {
		return GetClassDetailsJSON{}, err
	}
	//log.Printf("%s: %s", classjson.Title, classjson.JoinTokens.OnDemand)
	return *classjson, nil
}

func (ps *PelotonStack) GetStack() {
	query, err := Asset("data/queries/ViewUserStack.graphql")
	if err != nil {
		log.Fatalf("Unable to get ViewUserStack: ", err)
	}

	resp, _ := ps.Session.Post(
		ps.Endpoint,
		&grequests.RequestOptions{
			JSON: map[string]string{
				"operationName": "ViewUserStack",
				"query":         string(query),
			},
			Headers: ps.Headers,
		},
	)

	if !resp.Ok {
		log.Fatalf("Unable to get user stack: %d - %s", resp.StatusCode, resp.String())
	}

	viewuserstackjson := &ViewUserStackResponseJSON{}

	err = resp.JSON(viewuserstackjson)
	if err != nil {
		log.Printf("Unable to serialize stack details JSON blob", err)
	} else {
		log.Printf("Number of classes in current stack: %d", viewuserstackjson.ViewUserStackData.ViewUserStack.NumClasses)
	}
}

type ModifyStackRequestJSON struct {
	OperationName string                          `json:"operationName"`
	Variables     ModifyStackRequestVariablesJSON `json:"variables"`
	Query         string                          `json:"query"`
}

type ModifyStackRequestVariablesJSON struct {
	Input ModifyStackRequestInputJSON `json:"input"`
}

type ModifyStackRequestInputJSON struct {
	PelotonClassIdList []string `json:"pelotonClassIdList"`
}

func (ps *PelotonStack) ClearStack() {
	query, err := Asset("data/queries/ModifyStack.graphql")
	if err != nil {
		log.Fatalf("Unable to get ModifyStack.graphql: ", err)
	}

	clearStackRequest := ModifyStackRequestJSON{
		OperationName: "ModifyStack",
		Variables: ModifyStackRequestVariablesJSON{
			Input: ModifyStackRequestInputJSON{
				PelotonClassIdList: []string{},
			},
		},
		Query: string(query),
	}
	resp, _ := ps.Session.Post(
		ps.Endpoint,
		&grequests.RequestOptions{
			JSON:    clearStackRequest,
			Headers: ps.Headers,
		},
	)

	if !resp.Ok {
		log.Fatalf("Unable to clear user stack: %d - %s", resp.StatusCode, resp.String())
	}
}

type AddClassToStackRequestJSON struct {
	OperationName string                              `json:"operationName"`
	Variables     AddClassToStackRequestVariablesJSON `json:"variables"`
	Query         string                              `json:"query"`
}

type AddClassToStackRequestVariablesJSON struct {
	Input AddClassToStackRequestInputJSON `json:"input"`
}

type AddClassToStackRequestInputJSON struct {
	PelotonClassId string `json:"pelotonClassId"`
}

func (ps *PelotonStack) AddClassToStack(joinToken string) {
	query, err := Asset("data/queries/AddClassToStack.graphql")
	if err != nil {
		log.Fatalf("Unable to get AddClassToStack.graphql: ", err)
	}

	addClassToStackRequest := AddClassToStackRequestJSON{
		OperationName: "AddClassToStack",
		Variables: AddClassToStackRequestVariablesJSON{
			Input: AddClassToStackRequestInputJSON{
				PelotonClassId: joinToken,
			},
		},
		Query: string(query),
	}

	resp, _ := ps.Session.Post(
		ps.Endpoint,
		&grequests.RequestOptions{
			JSON:    addClassToStackRequest,
			Headers: ps.Headers,
		},
	)

	if !resp.Ok {
		log.Fatalf("Unable to add class to user stack: %d - %s", resp.StatusCode, resp.String())
	}

	//fmt.Println("RESP:", resp.String())
}

func main() {
	currentTime := time.Now()
	currentDay := currentTime.Format("1/02/2006")

	fmt.Printf("Today is: %s\n", currentDay)

	username := os.Getenv("PELOTON_USERNAME")
	password := os.Getenv("PELOTON_PASSWORD")
	ps := NewPelotonSession(username, password)
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
