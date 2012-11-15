package main

import (
	"fmt"
	"time"
)

type Info_Direction struct{
	Id int
	Name string
	Duration string
}

type Statistic struct{
	Date string
	List_Info_Direction []Info_Direction
	Sum_Duration string
}

func print_date(dt time.Time)string {
	//
	var days int
	now :=  time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

	if dt.Unix() == now.Unix(){
		return "today"
	}else{
		delta := now.Sub(dt)
		days = int(delta.Hours())/24		
	}	
	return fmt.Sprintf("%d days ago", int(days))
}

func Get_Statistic()[]Statistic {
	//
	var result []Statistic
	for i:=0;i<5;i+=1{
		//
		map_direct_duration := make(map[string]time.Duration)
		dt_now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

		//dur_day, _ := time.ParseDuration(fmt.Sprintf("%dh", 24 * i))
		dt_now = dt_now.AddDate(0, 0, -1*i)
		//fmt.Printf("date:%v\n", dt_now)

		for _, val := range Task_map{
			DateDead, err := time.Parse("2006-01-02 15:04:05",val.DateDead)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			if DateDead.Unix() > dt_now.Unix() && DateDead.Unix() < dt_now.AddDate(0, 0, 1).Unix(){
				//fmt.Printf("DateDead:%v\n", DateDead)
				if val.IsDone{
					dur_task, err := time.ParseDuration(val.Duration)
					if err != nil {
						fmt.Printf("%v\n", err)
					}
					map_direct_duration[val.Direction] += dur_task
				}else{
					dur_task, err := time.ParseDuration(val.Work_time)
					if err != nil {
						fmt.Printf("%v\n", err)
					}
					if dur_task.Seconds() > 0{
						map_direct_duration[val.Direction] += dur_task
					}
				}
			}
		}

		var list_info []Info_Direction
		var sum time.Duration

		for name, val := range map_direct_duration{
			var info Info_Direction
			info.Name = name
			info.Duration = pretty_duration(val.String())
			list_info = append(list_info, info)
			sum += val
		}

		var sts Statistic
		sts.Date = print_date(dt_now)
		sts.List_Info_Direction = list_info
		sts.Sum_Duration = pretty_duration(sum.String())
		//fmt.Printf("Statistic:%v\n", sts)
		if len(list_info) > 0{
			result = append(result, sts)
		}
	}
	return result
}