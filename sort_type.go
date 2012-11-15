package main

type Tasks []*Task

func (s Tasks) Len() int      { return len(s) }
func (s Tasks) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByPriority_Tasks struct{ Tasks }

func (s ByPriority_Tasks) Less(i, j int) bool { return s.Tasks[i].Priority > s.Tasks[j].Priority }

func convert_Task(tasks []Task) []*Task {
	//
	var result []*Task
	for i, _ := range tasks {
		result = append(result, &tasks[i])
	}
	return result
}

func convert_back_Task(tasks []*Task) []Task {
	//
	var result []Task
	for i, _ := range tasks {
		result = append(result, *tasks[i])
	}
	return result
}

type Directions []*Direction

func (s Directions) Len() int      { return len(s) }
func (s Directions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByPriority_Directions struct{ Directions }

func (s ByPriority_Directions) Less(i, j int) bool {
	return s.Directions[i].Priority > s.Directions[j].Priority
}

func convert_Directions(directions []Direction) []*Direction {
	//
	var result []*Direction
	for i, _ := range directions {
		result = append(result, &directions[i])
	}
	return result
}
func convert_back_Directions(directions []*Direction) []Direction {
	//
	var result []Direction
	for i, _ := range directions {
		result = append(result, *directions[i])
	}
	return result
}
