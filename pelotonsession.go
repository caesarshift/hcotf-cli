package main

import "fmt"
import "log"

import "github.com/levigross/grequests"

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
		log.Printf("Unable to make request", resp.Error)
	}
	if resp.Ok {
		log.Printf("Successfully created Peloton Session")
	} else {
		if resp.StatusCode == 401 {
			log.Printf("Invalid username or password: %d - %s", resp.StatusCode, "")
		} else {
			log.Printf("Unable to login and create a Peloton session: %d - %s", resp.StatusCode, "")
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
