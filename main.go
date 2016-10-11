package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

//PermissionsPost is the go structure for the json body of a POST to /permissions
type PermissionsPost struct {
	Subject     string `json:"_subject"`
	Permissions struct {
		Read   string `json:"read"`
		Update string `json:"update"`
	} `json:"_permissions"`
	Self struct {
		Update string `json:"update"`
		Read   string `json:"read"`
		Delete string `json:"delete"`
	} `json:"_self"`
	PermissionsHeirs struct {
		Add    string `json:"add"`
		Read   string `json:"read"`
		Remove string `json:"remove"`
	} `json:"_permissionsHeirs"`
	TestData bool `json:"test-data"`
}

//Permissions object
type Permissions struct {
	Read   string `json:"read"`
	Update string `json:"update"`
}

//Self object
type Self struct {
	Update string `json:"update"`
	Read   string `json:"read"`
	Delete string `json:"delete"`
}

//PermissionsHeirs object
type PermissionsHeirs struct {
	Add    string `json:"add"`
	Read   string `json:"read"`
	Remove string `json:"remove"`
}

//TestClaim object
type TestClaim struct {
	Jti       string   `json:"jti"`
	Sub       string   `json:"sub"`
	Scope     []string `json:"scope"`
	ClientID  string   `json:"client_id"`
	Cid       string   `json:"cid"`
	Azp       string   `json:"azp"`
	GrantType string   `json:"grant_type"`
	UserID    string   `json:"user_id"`
	Origin    string   `json:"origin"`
	UserName  string   `json:"user_name"`
	Email     string   `json:"email"`
	AuthTime  int      `json:"auth_time"`
	RevSig    string   `json:"rev_sig"`
	Iat       int      `json:"iat"`
	Exp       int      `json:"exp"`
	Iss       string   `json:"iss"`
	Zid       string   `json:"zid"`
	Aud       []string `json:"aud"`

	jwt.StandardClaims
}

func main() {

	httpClient := &http.Client{}

	tokenBytes, err := ioutil.ReadFile("token1.txt")
	if err != nil {
		fmt.Printf("Got Err at ReadFile: %v\n", err)
	}
	tokenString := string(tokenBytes)

	token, err := jwt.ParseWithClaims(tokenString, &TestClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	claims := token.Claims.(*TestClaim)

	userID := claims.Iss + "#" + claims.Sub
	fmt.Printf("UserID: %v\n", userID)

	tempPermissions := PermissionsPost{
		Subject: "http://apigee.com/o/acme",
		Permissions: Permissions{
			Read:   userID,
			Update: userID,
		},
		Self: Self{
			Update: userID,
			Read:   userID,
			Delete: userID,
		},
		PermissionsHeirs: PermissionsHeirs{
			Add:    userID,
			Read:   userID,
			Remove: userID,
		},
	}

	jsonPermissions := new(bytes.Buffer)
	json.NewEncoder(jsonPermissions).Encode(tempPermissions)

	req1, err := http.NewRequest("POST", "http://localhost:8080/permissions", jsonPermissions)
	if err != nil {
		fmt.Printf("Got Err at NewRequest: %v\n", err)
	}
	req1.Header.Add("Accept", "application/json")
	req1.Header.Add("Authorization", "Bearer "+tokenString)

	fmt.Printf("Making POST Request\n")
	resp1, err := httpClient.Do(req1)
	if err != nil {
		fmt.Printf("Got Err at httpClient.Do(req1): %v\n", err)
	}
	fmt.Printf("Response Status: %v\n", resp1.Status)

	//Make a GET request to the permissions
	orgURL := "http://localhost:8080/permissions?" + tempPermissions.Subject
	req2, err := http.NewRequest("GET", orgURL, nil)
	req2.Header.Add("Accept", "application/json")
	req2.Header.Add("Authorization", "Bearer "+tokenString)

	fmt.Printf("Making GET Request\n")

	resp2, err := httpClient.Do(req2)
	if err != nil {
		fmt.Printf("Got Err at httpClient.Do(req2): %v\n", err)
	}
	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		fmt.Printf("Got Err at ioutil.ReadAll: %v\n", err)
	}
	fmt.Printf("Response Status: %v\n", resp2.Status)

	var jsonBlob interface{}
	err = json.Unmarshal(body, &jsonBlob)
	jsonMap := jsonBlob.(map[string]interface{})
	fmt.Printf("JSON:\n%v\n", jsonMap)

}
