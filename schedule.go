package main

import (
	"fmt"
	"sort"
	"time"
)

var (
	//количество направлений в день
	count_direction_in_day = 3
	last_hour, last_minute = 0, 0
	duration_task_minute   = "0.5h"

	Work_schedule []Schedule
)

func init() {
	//
	start := time.Now()
	UpdateSchedule()
	end := time.Now()
	delta := end.Sub(start)
	fmt.Printf("Ok duration:%v\n", delta)
}

//расписание
type Schedule struct {
	Id        int
	Task      Task   //задача
	Start     string //начало задачи
	IsDone    bool   // задача в расписании выполнена?
	count_add int    //количество одной задачи в общем расписании.
}

func UpdateSchedule() {
	//
	Work_schedule = PrioritySchedule()
}

func less_2_1_duration(str_delta1, str_delta2 string) bool {
	//
	delta1, err := time.ParseDuration(str_delta1)
	if err != nil {
		fmt.Println(err)
	}
	delta2, err := time.ParseDuration(str_delta2)
	if err != nil {
		fmt.Println(err)
	}
	if delta1 >= delta2 {
		return true
	}
	return false
}

func scale(str_delta1, str_delta2 string) int {
	//
	delta1, err := time.ParseDuration(str_delta1)
	if err != nil {
		fmt.Println(err)
	}
	delta2, err := time.ParseDuration(str_delta2)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("int(delta1/delta2):%v\n", int(delta1/delta2))
	return int(delta1 / delta2)
}

func GetScheduleOne(id int) Schedule {
	//
	for _, val := range Work_schedule {
		if val.Id == id {
			return val
		}
	}
	return Schedule{}
}

func Done_Schedule(id int, isdone bool) {
	//
	for i, val := range Work_schedule {
		if val.Id == id {
			delta1, err := time.ParseDuration(Work_schedule[i].Task.Work_time)
			if err != nil {
				fmt.Println(err)
			}
			delta2, err := time.ParseDuration(duration_task_minute)
			if err != nil {
				fmt.Println(err)
			}
			if isdone {
				Work_schedule[i].Task.Work_time = (delta1 + delta2).String()
				temp := Task_map[Work_schedule[i].Task.Id]
				temp.DateDead = time.Now().Format("2006-01-02 15:04:05")
				temp.Work_time = (delta1 + delta2).String()
				Task_map[Work_schedule[i].Task.Id] = temp
				updateTaskList()
			} else {
				Work_schedule[i].Task.Work_time = (delta1 - delta2).String()
				temp := Task_map[Work_schedule[i].Task.Id]
				temp.Work_time = (delta1 - delta2).String()
				Task_map[Work_schedule[i].Task.Id] = temp
				updateTaskList()
			}
			if less_2_1_duration(Work_schedule[i].Task.Work_time, val.Task.Duration) {
				DoneTask(val.Task.Id)
				UpdateSchedule()
			}
		}
	}
	UpdateSchedule()
}

func printTask(s []*Task) {
	for _, o := range s {
		fmt.Printf("Id:%-8v Priority:(%v)\n", o.Id, o.Priority)
	}
}

func printDirection(s []*Direction) {
	for _, o := range s {
		fmt.Printf("Id:%-8v Priority:(%v)\n", o.Id, o.Priority)
	}
}

func is_date_in_WhenWork(datetime string, WhenWork []Period) bool {
	//
	dt, err := time.Parse("15:04", datetime)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	for _, val := range WhenWork {
		start, err := time.Parse("15:04", val.Start)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		stop, err := time.Parse("15:04", val.Stop)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		//fmt.Printf("1:%v  2:%v  3%v  itog:%v\n", start.Unix(), dt.Unix(), stop.Unix(), (start.Unix() <= dt.Unix()) && (dt.Unix() <= stop.Unix()))
		//fmt.Printf("1:%v  2:%v  3%v  itog:%v\n\n", val.Start, datetime, val.Stop, (start.Unix() <= dt.Unix()) && (dt.Unix() <= stop.Unix()))
		if (start.Unix() <= dt.Unix()) && (dt.Unix() <= stop.Unix()) {
			return true
		}
	}
	return false
}

func what_duration_direction_done(direction_id int, date time.Time) time.Duration {
	//
	var result time.Duration

	tasks := GetTasks(direction_id, true)
	for _, val := range tasks {
		date_dead, err := time.Parse("2006-01-02 15:04:05", val.DateDead)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		if date_dead.Unix() > time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).Unix() {
			if val.IsDone {
				dur_task, err := time.ParseDuration(val.Duration)
				if err != nil {
					fmt.Println(err)
				}
				result += dur_task
			} else {
				dur_task, err := time.ParseDuration(val.Work_time)
				if err != nil {
					fmt.Println(err)
				}
				result += dur_task
			}

		}
	}
	//fmt.Printf("ID : %d    duration: %v\n", direction_id, result)
	return result
}

func PrioritySchedule() []Schedule {
	//
	var result []Schedule
	var directions_list_filtering []Direction

	//fmt.Printf("len(Direction_list):%d\n", len(Direction_list))

	for _, val := range Direction_list {
		if val.IsWork &&
			(len(GetTasks(val.Id, false)) > 0) &&
			is_date_in_WhenWork(time.Now().Format("15:04"), val.WhenWork) {
			directions_list_filtering = append(directions_list_filtering, val)
		}
	}
	count_direction_in_day = len(directions_list_filtering)
	//fmt.Printf("len(directions_list_filtering):%d\n", len(directions_list_filtering))
	// if len(directions_list_filtering) < count_direction_in_day {
	// 	count_direction_in_day = len(directions_list_filtering)
	// } else {
	// 	if len(directions_list_filtering) <= 3 {
	// 		count_direction_in_day = len(directions_list_filtering)
	// 	}
	// }

	directions := convert_Directions(directions_list_filtering)
	sort.Sort(ByPriority_Directions{directions})
	work_direction := directions[0:count_direction_in_day]

	map_HowLongDay := make(map[int]time.Duration)

	for k := 0; k < 12; k++ {
		for i, direct := range work_direction {
			//
			howlongDay_direct, err := time.ParseDuration(direct.HowLongDay)
			if err != nil {
				fmt.Printf("howlongDay_direct_%v: %v\n", direct.Name, err)
			}
			map_HowLongDay[direct.Id] = what_duration_direction_done(direct.Id, time.Now())

			tasks := convert_Task(GetTasks(direct.Id, false))
			sort.Sort(ByPriority_Tasks{tasks})
			if k < len(tasks) {
				task := *tasks[k]
				if is_date_in_WhenWork(time.Now().Format("15:04"), GetDirection(task.Direction_Id).WhenWork) {
					duration_task, err := time.ParseDuration(task.Duration)
					if err != nil {
						fmt.Printf("duration_task: %v\n", err)
					}

					work_time_task, err := time.ParseDuration(task.Work_time)
					if err != nil {
						fmt.Printf("work_time_task: %v\n", err)
					}
					m30, _ := time.ParseDuration("30m")
					if map_HowLongDay[direct.Id]+m30 <= howlongDay_direct {
						result = append(result, Schedule{i + k*100 + 1,
							task,
							"1:1",
							false,
							scale((duration_task - work_time_task).String(), duration_task_minute)})
					}

					map_HowLongDay[direct.Id] += duration_task
				}
			}
		}
	}
	return result
}

func SeachLastDirectionTaskIsDone() int {
	//
	var id int
	var time_dt time.Time
	time_dt, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")

	for _, val := range Task_map {
		dt, eror := time.Parse("2006-01-02 15:04:05", val.DateDead)
		if eror != nil {
			fmt.Printf("Error: %v\n", eror)
			dt, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
		}
		if dt.Unix() > time_dt.Unix() {
			time_dt = dt
			id = val.Direction_Id
		}
	}
	if time_dt.Day() != time.Now().Day() {
		return 0
	}
	return id
}

func Get_Schedule() []Schedule {
	//
	UpdateSchedule()
	var result [][]Schedule

	last_id_direction := SeachLastDirectionTaskIsDone()
	//fmt.Printf("last_id_direction:%d\n", last_id_direction)
	temp := make(map[int]int)
	m15, _ := time.ParseDuration("15m")
	now := time.Now().Add(-1 * m15)
	hour_now, minute_now := now.Hour(), now.Minute()
	if minute_now > 35 {
		hour_now += 1
		minute_now = 0
	} else {
		minute_now = 30
	}
	for _, schedule := range Work_schedule {
		var temp_schedule []Schedule
		for i := 0; i < 24; i++ {
			if schedule.count_add > temp[schedule.Id] {
				temp[schedule.Id] += 1
				temp_schedule = append(temp_schedule, schedule)
			}
		}
		result = append(result, temp_schedule)
	}
	last_result := make([][]Schedule, count_direction_in_day)
	for k := 0; k < (len(Work_schedule) / count_direction_in_day); k++ {
		for i := 0; i < count_direction_in_day; i++ {
			last_result[i] = append(last_result[i], result[k*count_direction_in_day+i]...)
		}
	}
	var itog []Schedule
	var max int
	for _, val := range last_result {
		if len(val) > max {
			max = len(val)
		}
	}
	//fmt.Printf("last_do %d %d %v\n",last_id_direction, last_result[0][0].Task.Direction_Id, last_result )
	if len(last_result) > 1 && len(last_result[0]) > 0 && len(last_result[1]) > 0 {
		if last_result[0][0].Task.Direction_Id == last_id_direction {
			temp := last_result[1]
			last_result[1] = last_result[0]
			last_result[0] = temp
		}
	}
	//fmt.Printf("last_posle %d %d %v\n",last_id_direction, last_result[0][0].Task.Direction_Id, last_result )

	map_HowLongDay := make(map[int]time.Duration)
	delta, err := time.ParseDuration(duration_task_minute)
	if err != nil {
		fmt.Printf("duration_task_minute: %v\n", err)
	}

	dt := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour_now, minute_now, 0, 0, time.UTC)
	for i := 0; i < max; i++ {
		for j := 0; j < count_direction_in_day; j++ {
			if len(last_result[j]) > i {
				schedule := last_result[j][i]

				HowLong, err := time.ParseDuration(GetDirection(schedule.Task.Direction_Id).HowLongDay)
				if err != nil {
					fmt.Printf("HowLong: %v\n", err)
				}

				if HowLong > map_HowLongDay[schedule.Task.Direction_Id] && 
					is_date_in_WhenWork(dt.Format("15:04"), GetDirection(schedule.Task.Direction_Id).WhenWork){
					schedule.Start = fmt.Sprintf("%d:%d", dt.Hour(), dt.Minute())
					itog = append(itog, schedule)

					map_HowLongDay[schedule.Task.Direction_Id] += delta
					dt = dt.Add(delta)
				}
			}
		}
	}
	if len(itog) <1{
		fmt.Printf("itog:%v", itog)
	}
	return itog
}
