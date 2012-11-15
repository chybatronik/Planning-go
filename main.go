package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main_page(w http.ResponseWriter, r *http.Request) {
	//
	content, err := ioutil.ReadFile("client/main.html")
	if err != nil {
		fmt.Printf("Error get_param:%v\n", err)
	}
	//w.Header().Set("Cache-control", "public, max-age=0")
	w.Write(content)
}

func direction(w http.ResponseWriter, r *http.Request) {
	//
	defer SaveJSON(namefile_store)
	//
	var resposible []byte

	fmt.Printf("%v  %v\n", r.Method, r.URL)
	url := fmt.Sprintf("%s", r.URL)
	if url == "/direction/" {
		switch r.Method {
		case "GET":
			resposible, _ = json.Marshal(GetDirections())

		case "PUT":
			body, _ := ioutil.ReadAll(r.Body)
			NewDirection(fmt.Sprintf("%s", body))
		}
	} else {

		switch r.Method {
		case "GET":
			temp := strings.Split(url, "/")
			item := temp[2]
			id, err := strconv.Atoi(item)
			if err != nil {
				fmt.Printf("Error strconv:%v\n", err)
			}

			fmt.Printf("url:%v  id:%d \n", url, id)
			val := GetDirection(id)
			resposible, _ = json.Marshal(val)

		case "POST":
			temp := strings.Split(url, "/")
			item := temp[2]
			id, err := strconv.Atoi(item)
			if err != nil {
				fmt.Printf("Error strconv:%v\n", err)
			}

			var p Direction
			body, _ := ioutil.ReadAll(r.Body)
			//fmt.Printf("url:%v  id:%d body:%s\n", url, id, body)
			err = json.Unmarshal(body, &p)
			if err != nil {
				fmt.Printf("Error:%v\n", err)
			}
			SetDirection(id, p)
		case "DELETE":
			temp := strings.Split(url, "/")
			item := temp[2]
			id, err := strconv.Atoi(item)
			if err != nil {
				fmt.Printf("Error strconv:%v\n", err)
			}
			DelDirection(id)
		}
	}
	w.Header().Set("Cache-control", "no-cache, must-revalidate")
	w.Write(resposible)
}

func task(w http.ResponseWriter, r *http.Request) {
	//
	defer SaveJSON(namefile_store)
	//
	var resposible []byte
	var executed bool
	var direction int

	fmt.Printf("%v  %v\n", r.Method, r.URL)

	u, err := url.Parse(fmt.Sprintf("%s", r.URL))
	if err != nil {
		fmt.Println(err)
	}
	q := u.Query()

	//fmt.Printf("Path:%v\n", u.Path)
	//fmt.Printf("Query.executed:%v\n", q.Get("executed"))
	if len(q) > 0 {
		if len(q.Get("executed")) > 0{
			executed, err = strconv.ParseBool(q.Get("executed"))
			if err != nil {
				fmt.Printf("Error  strconv.ParseBool:%v\n", err)
			}
			fmt.Printf("executed:%v\n", executed)
		}

		if len(q.Get("direction")) > 0{
			direction64, err := strconv.ParseInt(q.Get("direction"), 10, 32)
			direction = int(direction64)
			if err != nil {
				fmt.Printf("Error  strconv.ParseBool:%v\n", err)
			}
			fmt.Printf("direction:%v\n", direction)
		}
	}

	mas_url := strings.Split(u.Path, "/")[1:]

	fmt.Printf("mas_url:%v %d \n", mas_url, len(mas_url))

	if mas_url[0] == "task" && direction == 0 && mas_url[1] == ""{
		switch r.Method {
		case "GET":
			resposible, _ = json.Marshal(Task_list)

		case "PUT":
			var p Task
			body, _ := ioutil.ReadAll(r.Body)
			//fmt.Println(fmt.Sprintf("%s", body))
			err := json.Unmarshal(body, &p)
			if err != nil {
				fmt.Printf("Error:%v\n", err)
			}
			NewTask(p.Name, p.Direction_Id, p.Duration, p.PriorityTask)
		}
	} else {

		switch r.Method {
		case "GET":
			//temp := strings.Split(url, "/")		
			if direction > 0 {
				val := GetTasks(direction, executed)
				if len(val) != 0 {
					resposible, _ = json.Marshal(val)
				}
			} else {
				item := mas_url[1]
				if item != "label" {
					id, err := strconv.Atoi(item)
					if err != nil {
						fmt.Printf("Error strconv:%v\n", err)
					}

					val := GetTask(id)
					resposible, _ = json.Marshal(val)
				} else {
					val := GetTaskLabel()
					resposible, _ = json.Marshal(val)
				}
			}
		case "POST":
			item := mas_url[1]
			id, err := strconv.Atoi(item)
			if err != nil {
				fmt.Printf("Error strconv:%v\n", err)
			}

			var p Task
			body, _ := ioutil.ReadAll(r.Body)
			err = json.Unmarshal(body, &p)
			if err != nil {
				fmt.Printf("Error:%v\n", err)
			}
			SetTask(id, p)
		case "DELETE":
			item := mas_url[1]
			id, err := strconv.Atoi(item)
			if err != nil {
				fmt.Printf("Error strconv:%v\n", err)
			}
			DelTask(id)
		}
	}
	w.Header().Set("Cache-control", "no-cache, must-revalidate")
	w.Write(resposible)
}

func schedule(w http.ResponseWriter, r *http.Request) {
	//
	var resposible []byte
	var err error

	fmt.Printf("schedule:%v  %v\n", r.Method, r.URL)
	url := fmt.Sprintf("%s", r.URL)
	if url == "/schedule/" {
		switch r.Method {
		case "GET":
			//Get_Schedule()
			//asd:=Get_Schedule()
			//fmt.Printf("schedule:%v\n", len(Get_Schedule()))
			resposible, err = json.Marshal(Get_Schedule())
			if err != nil {
				fmt.Printf("Json errror:%v\n", err)
			}
		}
	} else {

		switch r.Method {
		case "GET":
			temp := strings.Split(url, "/")
			item := temp[2]
			id, err := strconv.Atoi(item)
			if err != nil {
				fmt.Printf("Error strconv:%v\n", err)
			}
			//fmt.Printf("schedule GetScheduleOne:url:%v  id:%d \n", url, id)
			val := GetScheduleOne(id)
			resposible, _ = json.Marshal(val)

		case "POST":
			temp := strings.Split(url, "/")
			item := temp[2]
			id, err := strconv.Atoi(item)
			if err != nil {
				fmt.Printf("Error strconv:%v\n", err)
			}

			var p Schedule
			body, _ := ioutil.ReadAll(r.Body)
			//fmt.Printf("schedule Done_Schedule:url:%v  id:%d body:%s\n", url, id, body)
			err = json.Unmarshal(body, &p)
			if err != nil {
				fmt.Printf("Error:%v\n", err)
			}
			Done_Schedule(id, p.IsDone)
		}
	}
	//_ = Get_Schedule()
	w.Header().Set("Cache-control", "no-cache, must-revalidate")
	w.Write(resposible)

}

func statistic(w http.ResponseWriter, r *http.Request) {
	//
	var resposible []byte
	var err error

	fmt.Printf("statistic:%v  %v\n", r.Method, r.URL)
	url := fmt.Sprintf("%s", r.URL)
	if url == "/statistic/" {
		switch r.Method {
		case "GET":
			resposible, err = json.Marshal(Get_Statistic())
			if err != nil {
				fmt.Printf("Json errror:%v\n", err)
			}
		}
	}
	w.Header().Set("Cache-control", "no-cache, must-revalidate")
	w.Write(resposible)
}

func backup() {
	//
	for {
		SaveJSON(namefile_backup)
		time.Sleep(time.Second * 60 * 5)
	}
}

func main() {
	//
	runtime.GOMAXPROCS(runtime.NumCPU())
	//
	go backup()
	//
	http.Handle("/client/",
		http.StripPrefix("/client", http.FileServer(http.Dir("client"))))
	http.HandleFunc("/", main_page)
	http.HandleFunc("/direction/", direction)
	http.HandleFunc("/task/", task)
	http.HandleFunc("/schedule/", schedule)
	http.HandleFunc("/statistic/", statistic)
	http.ListenAndServe(":8080", nil)
}
