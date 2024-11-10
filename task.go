package main

/*CUSTOM ITEM*/
type Task struct {
	status      status
	title       string
	description string
}

func NewTask(status status, title, description string) Task {
	return Task{
		status:      status,
		title:       title,
		description: description,
	}
}

// implement the List.Item interface

func (t Task) FilterValue() string {
	return t.title
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// Getters
func (t Task) Title() string       { return t.title }
func (t Task) Description() string { return t.description }
