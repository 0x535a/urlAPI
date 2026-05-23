package op

import (
	"github.com/pkg/errors"
	"reflect"
	"sort"
	"urlAPI/internal/model"
	"urlAPI/util"
)

func fetchTask(info *Session) error {
	var taskGetter model.Task
	v := reflect.ValueOf(&taskGetter).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("json")
		if tag == info.TaskCatagory && tag != "time" {
			v.Field(i).Set(reflect.ValueOf(info.TaskBy))
		}
	}
	if info.TaskCatagory == "time" {
		taskGetter.Time = util.GetDate(info.TaskBy)
	}
	taskDBList, err := db.ReadTask(taskGetter)
	if err != nil {
		return errors.WithStack(err)
	}
	taskList := taskDBList.TaskList
	if len(taskList) == 0 {
		info.TaskMaxPage = 0
		info.TaskData = nil
		return nil
	}
	info.TaskMaxPage = ((len(taskList) - 1) / 100) + 1
	sort.Slice(taskList, func(i, j int) bool {
		return taskList[i].Time.After(taskList[j].Time)
	})
	switch {
	case info.TaskPage == -1:
		info.TaskData = taskList
	default:
		page := info.TaskPage
		if page < 1 {
			page = 1
		}
		start := (page - 1) * 100
		if start >= len(taskList) {
			info.TaskData = nil
			return nil
		}
		end := start + 100
		if end > len(taskList) {
			end = len(taskList)
		}
		info.TaskData = taskList[start:end]
	}
	return nil
}
