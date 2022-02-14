package cmd

type Task struct {
	Name      string
	Completed bool
}

func (t *Task) complete() {
	t.Completed = true
}

func NewTask(name string) *Task {
	return &Task{Name: name, Completed: false}
}
