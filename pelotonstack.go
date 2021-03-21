package main

import "log"

import "github.com/levigross/grequests"

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
