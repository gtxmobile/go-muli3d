package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)
var privateToken string= "_hg5q_KprH2qicrwi7E-"
type assignee struct {
	assignee_id	int	`json:"assignee_id"`
}
func updateAssignee(username string,projectID int64,mriid int) bool{

	     values := map[string]int{"assignee_id": 574}
			d ,_:= json.Marshal(values)

	     reqURL := fmt.Sprintf("http://git-internal.nie.netease.com/api/v4/projects/%d/merge_requests/%d", projectID,mriid)
	     log.Println(reqURL)
	     //resp,err := http.PostForm(reqURL, values)
	     //body := strings.NewReader(values.Encode())

	     resp,err := doJsonRequest(http.MethodPut,reqURL,"application/json; charset=utf-8",bytes.NewBuffer(d),nil)
	     fmt.Println(resp,err)
	     return true
 }
type GitUser struct {
    ID      int     `json:"id"`
    Username    string  `json:"username"`
    Name        string  `json:"name"`
}

 func getUserByName(username string) GitUser{
     reqURL := fmt.Sprintf("http://git-internal.nie.netease.com/api/v4/users?username=%s", username)
     println(reqURL)

     user := []GitUser{}
     resp1,err := doJsonRequest("GET", reqURL, "", nil,&user)
     //bs := string(body)

     if err != nil{
         return GitUser{Username:"no user"}
     }
     fmt.Println(resp1)
     println(len(user))
     println(user[0].Name,user[0].Username)
     return user[0]
 }

func getUserByID(userid int) string{
	reqURL := fmt.Sprintf("http://git-internal.nie.netease.com/api/v4/users/%d", userid)
	println(reqURL)
	user := GitUser{}
	resp1,err := doJsonRequest("GET", reqURL, "", nil,&user)
	//bs := string(body)
	if err != nil{
		return "No User"
		}
	println(resp1)
	println(user.Name,user.Username)
	return user.Username
}


func doJsonRequest(method, urlStr string, bodyType string, body io.Reader, data interface{}) (resp *http.Response, err error) {
	if privateToken == "" {
		return nil, errors.New("missing --private-token")
	}

	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return
	}

	req.Header.Set("PRIVATE-TOKEN", privateToken)
	if bodyType != "" {
		req.Header.Set("Content-Type", bodyType)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	//bs,_ := ioutil.ReadAll(resp.Body)
	//println(resp.StatusCode)
	//fmt.Println(string(bs))

	if resp.StatusCode/100 == 2 {
		if data != nil {
			d := json.NewDecoder(resp.Body)
			err = d.Decode(data)
			fmt.Println(err)
		}
	}else{
		err = errors.New(resp.Status)
	}
	return
}

func sub(a interface{},b interface{})interface{}{
	return a


}
func main() {
	//fmt.Println("Hello, %s",sub(2,1))
	//fmt.Println("Hello, %s",sub(3.2,1.1))
	tm,_ := time.Parse(time.RFC3339,"2012-01-03T23:35:21+02:00")
	//tm2,_ := time.Parse(time.RFC3339,"2012-01-03T23:36:21+02:00")
	fmt.Println(time.Now())
	fmt.Println(tm.Unix()-time.Now().Unix())
	//updateAssignee("cydn1579",692,25)
	//getUserByID(574)
	getUserByName("cydn1579")
}
