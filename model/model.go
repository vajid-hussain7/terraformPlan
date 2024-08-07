package model

type ResourseChanges struct {
	Resourses []Resourses
}

type Resourses struct {
	Type   string
	Action string
}
