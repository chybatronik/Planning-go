package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
	"strings"
	"sort"
)

var (
	Direction_map  map[int]Direction
	Direction_list []Direction
	Direction_id   int

	Task_map  map[int]Task
	Task_list []Task
	Task_id   int

	namefile_store  = "tasks.json"
	namefile_backup = "tasks_backup.json"
)

//период с датой начала и конца
type Period struct {
	Start string
	Stop  string
}

//направление
type Direction struct {
	Id              int
	Name            string
	Priority        int      //приоритет
	WhenWork        []Period //в какое время задачи должны выполняться
	DurationDone    string   //сколько время потрачено на задачи
	DurationOstalos string   //сколько время осталось на задачи
	IsWork       bool     //тру - направление отключено
	HowLongDay      string   //сколько времени в день работать над направлением
}

//повторение
type Repeat struct{
	WhatWeekyDayRestore []string // по каким дням 
	Work_time    string //отработанное время над задачей
	LastRestore     int //в какой день послендний раз востанавливали
}

//задача
type Task struct {
	Id           int
	Name         string
	Duration     string //продолжительность задачи
	Direction_Id int    //id направления
	Direction    string //название направления
	Label        bool   //задача
	PriorityTask int    //приоритет задачи
	Priority     int    //приоритет итоговы = задача + направления
	IsDone       bool   //задача выполнена?
	Work_time    string //отработанное время над задачей
	DateDead     string //время выполнения задачи или время выполнения этапа отработки
	ActiveRepeat bool   //работает повторение
	Repeat Repeat // параметры повторения
}

//для сохранения объекто в файл
type Enum_key_string struct {
	Task_map      map[string]Task
	Direction_map map[string]Direction
}

//для вытаскивания объекто из файла
type Enum_key_int struct {
	Task_map      map[int]Task
	Direction_map map[int]Direction
}

func init() {
	//	
	Direction_map = make(map[int]Direction)
	Task_map = make(map[int]Task)
	LoadJSON()
}

func LoadJSON() {
	//

	data, error := ioutil.ReadFile(namefile_store)
	if error != nil {
		fmt.Printf("error: %v", error)
	}

	format := new(Enum_key_string)

	error = json.Unmarshal(data, format)
	if error != nil {
		fmt.Printf("error: %v", error)
	}
	form := map_str_int(*format)
	Task_map = form.Task_map
	Direction_map = form.Direction_map
	updateDirectionList()
	updateTaskList()
	RestoreTask()
	//fmt.Println(map_str_int(*format))
}

func map_int_str(data Enum_key_int) Enum_key_string {
	//
	var result Enum_key_string
	temp_Task_map := make(map[string]Task)
	temp_Direction_map := make(map[string]Direction)

	for i, val := range data.Task_map {
		temp_Task_map[strconv.Itoa(i)] = val
	}
	for i, val := range data.Direction_map {
		temp_Direction_map[strconv.Itoa(i)] = val
	}

	result.Task_map = temp_Task_map
	result.Direction_map = temp_Direction_map
	return result
}

func map_str_int(data Enum_key_string) Enum_key_int {
	//
	var result Enum_key_int
	temp_Task_map := make(map[int]Task)
	temp_Direction_map := make(map[int]Direction)

	for i, val := range data.Task_map {
		d, err := strconv.Atoi(i)
		if err != nil {
			fmt.Println(err)
		}
		temp_Task_map[d] = val
	}
	for i, val := range data.Direction_map {
		d, err := strconv.Atoi(i)
		if err != nil {
			fmt.Println(err)
		}
		temp_Direction_map[d] = val
	}

	result.Task_map = temp_Task_map
	result.Direction_map = temp_Direction_map
	return result
}

func SaveJSON(namefile string) {
	//	
	var result Enum_key_int

	result.Task_map = Task_map
	result.Direction_map = Direction_map

	b, err := json.Marshal(map_int_str(result))
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(namefile, b, 0666)
	if err != nil {
		fmt.Printf("Unexpected error: %s\n", err)
	}
}

func pretty_duration(s string) string{
	//
	res := strings.Replace(s, "h0m", "h", -1)
	res = strings.Replace(res, "m0s", "m", -1)
	res = strings.Replace(res, "h0s", "h", -1)
	res = strings.Replace(res, "h", "h ", -1)
	res = strings.Replace(res, "m", "m ", -1)
	return res
}

func calc_Duration_Done_All(Direction_id int) (string, string) {
	//
	var delta_work time.Duration
	var delta_all time.Duration

	tasks_work := GetTasks(Direction_id, true)
	for _, val := range tasks_work {
		var delta time.Duration
		var err error

		if val.IsDone{
			delta, err = time.ParseDuration(val.Duration)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}else{
			delta, err = time.ParseDuration(val.Work_time)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
		delta_work += delta
	}

	tasks_all := GetTasks(Direction_id, false)
	for _, val := range tasks_all {
		delta, err := time.ParseDuration(val.Duration)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		delta_all += delta
	}
	return pretty_duration(delta_work.String()), pretty_duration(delta_all.String())
}

func updateDirectionList() {
	//
	Direction_list = []Direction{}
	for _, val := range Direction_map {
		val.DurationDone, val.DurationOstalos = calc_Duration_Done_All(val.Id)
		//fmt.Printf("%v    %v\n", val.DurationWork, val.DurationAll)
		Direction_list = append(Direction_list, val)
	}
}

func updateTaskList() {
	//
	Task_list = []Task{}
	for _, val := range Task_map {
		if !val.IsDone {
			Task_list = append(Task_list, val)
		}
	}
}

func NewTask(name string, direction_id int, duration string, PriorityTask int) Task {
	//
	Task_id += 1
	var res Task
	_, ok := Task_map[Task_id]
	if !ok {
		res = Task{Task_id, name, duration,
			direction_id,
			GetDirection(direction_id).Name,
			false,
			PriorityTask,
			GetDirection(direction_id).Priority + PriorityTask,
			false,
			"0h",
			time.Now().Format("2006-01-02 15:04:05"),
			false, 
			Repeat{[]string{""}, "", -1}}
		Task_map[Task_id] = res
		//fmt.Printf("Task_id:%d\n", Task_id)
		updateTaskList()
	} else {
		return NewTask(name, direction_id, duration, PriorityTask)
	}
	return res
}

func GetTaskLabel() []Task {
	//
	var res []Task
	for _, val := range Task_map {
		if val.Label && !val.IsDone {
			res = append(res, val)
		}
	}
	if len(res) == 0 {
		return []Task{}
	}
	return res
}

func GetTasks(id int, executed bool) []Task {
	//
	var res []Task
	for _, val := range Task_map {
		if executed {
			if val.Direction_Id == id {
				res = append(res, val)
			}
		} else {
			if val.Direction_Id == id && !val.IsDone {
				res = append(res, val)
			}
		}

	}
	tasks := convert_Task(res)
	sort.Sort(ByPriority_Tasks{tasks})
	return convert_back_Task(tasks)
}

func SetTask(id int, task Task) {
	task.Priority = GetDirection(task.Direction_Id).Priority + task.PriorityTask
	task.Direction = GetDirection(task.Direction_Id).Name
	tmp, _ := Task_map[id]
	if task.IsDone && !tmp.IsDone {
		task.DateDead = time.Now().Format("2006-01-02 15:04:05")
	}
	Task_map[id] = task
	updateTaskList()
}

func GetTask(id int) Task {
	//
	res, _ := Task_map[id]
	return res
}

func DelTask(id int) {
	//
	delete(Task_map, id)
	updateTaskList()
}

func GetDirections() []Direction {
	// /
	updateDirectionList()
	directions := convert_Directions(Direction_list)
	sort.Sort(ByPriority_Directions{directions})
	return convert_back_Directions(directions)
}

func GetDirection(id int) Direction {
	//
	res, err := Direction_map[id]
	if !err {
		return Direction{}
	}
	return res
}

func DoneTask(id int) {
	//
	temp := Task_map[id]
	temp.IsDone = true
	temp.DateDead = time.Now().Format("2006-01-02 15:04:05")
	Task_map[id] = temp
	updateTaskList()
}

func UpdatePriorityTask(direct Direction) {
	//
	fmt.Printf("%v\n", direct)
	for i, val := range Task_map {
		if val.Direction_Id == direct.Id {
			val.Priority = direct.Priority + val.PriorityTask
			val.Direction = direct.Name
			Task_map[i] = val
			//fmt.Printf("val:%v\n", val)
		}
	}
	updateTaskList()
	UpdateSchedule()
}

func SetDirection(id int, direct Direction) {
	//
	Direction_map[id] = direct
	updateDirectionList()
	UpdatePriorityTask(direct)
}

func DelDirection(id int) {
	delete(Direction_map, id)
	for _, val := range Task_map {
		if val.Direction_Id == id {
			delete(Task_map, val.Id)
		}
	}
	updateTaskList()
	updateDirectionList()
}

func find_max_priority_direction() int {
	//
	max := -100000000
	for _, val := range Direction_map {
		if val.Priority > max {
			max = val.Priority
		}
	}
	return max
}

func NewDirection(name string) Direction {
	Direction_id += 1
	var res Direction
	_, ok := Direction_map[Direction_id]
	if !ok {
		res = Direction{Direction_id, name, find_max_priority_direction() + 100,
			[]Period{Period{"9:00", "17:00"}}, "", "", false, "24h"}
		Direction_map[Direction_id] = res
		//fmt.Printf("%d\n", Direction_id)
		updateDirectionList()
	} else {
		return NewDirection(name)
	}
	return res
}

func conv_str_day_to_int(str string)int{
	result := map[string]int{
		"Sunday"  : int(time.Sunday),
	    "Monday"  : int(time.Monday),
	    "Tuesday" : int(time.Tuesday),
	    "Wednesday":int(time.Wednesday),
	    "Thursday" :int(time.Thursday),
	    "Friday"   :int(time.Friday),
	    "Saturday" :int(time.Saturday)}
	return result[str]	
}

func when_weekDay_last(mass_day []int)int {
	//
	datetime := time.Now()
	day := int(datetime.Weekday())
	var last_val int
	last_val = -1
	sort.Ints(mass_day)
	for i, val := range mass_day{
		if i != 0{
			if last_val > day && day < val{
				return last_val
			}
			if val > day && (len(mass_day)-1) == i{
				return val
			}
		}else{
			if val == day{
				return val
			}
			if day < val{
				return mass_day[len(mass_day)-1]
			}
		}
		last_val = val
	}
	return -1
}

func RestoreTask() {
	//
	for id, val := range Task_map{
		if val.ActiveRepeat && val.IsDone{
			//datetime := time.Now()

			//days := strings.Replace(val.Repeat.WhatWeekyDayRestore, " ", "", -1) 
			mass_day := val.Repeat.WhatWeekyDayRestore
			var mass_day_int []int
			for _, val := range mass_day{
				mass_day_int = append(mass_day_int, conv_str_day_to_int(val))
			}

			if val.Repeat.LastRestore == -1{
				val.Work_time = "0h"
				val.IsDone = false
				val.Repeat.Work_time = val.Duration
				val.Repeat.LastRestore = when_weekDay_last(mass_day_int)
				fmt.Println(val)
				Task_map[id] = val
			}else{
				if val.Repeat.LastRestore != when_weekDay_last(mass_day_int){

					work_time, err := time.ParseDuration(val.Repeat.Work_time)
					if err != nil {
						fmt.Println(err)
					}

					durat, err := time.ParseDuration(val.Duration)
					if err != nil {
						fmt.Println(err)
					}


					val.Work_time = "0h"
					val.IsDone = false
					val.Repeat.Work_time = (durat+work_time).String()
					val.Repeat.LastRestore = when_weekDay_last(mass_day_int)

					Task_map[id] = val
					fmt.Println(val)
				}
			}

		}
	}
}