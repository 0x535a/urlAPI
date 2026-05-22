package op

import (
	"os"
	"time"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
	"urlAPI/util"
)

func ExecuteCachedTask(task model.Task, filter TaskQueueFilter, skipDB bool, process func(*model.Task) (GenerateResult, error)) (model.Task, GenerateResult, error) {
	preparedTask, returnURL, cached := prepareTask(task, filter)
	if cached {
		if !skipDB {
			util.ErrorPrinter(db.CreateTask(&preparedTask))
		}
		return preparedTask, returnURL, nil
	}

	returnURL, err := process(&preparedTask)
	finishTask(&preparedTask, filter, &returnURL, skipDB)
	return preparedTask, returnURL, err
}

func SaveTask(task model.Task, skipDB bool) {
	if !skipDB {
		util.ErrorPrinter(db.CreateTask(&task))
	}
}

func prepareTask(task model.Task, filter TaskQueueFilter) (model.Task, GenerateResult, bool) {
	settings := database.SettingsStore.Get()
	expiredTime := settings.Text.CacheMinutes
	switch filter.Type {
	case "img.gen":
		expiredTime = settings.Image.CacheMinutes
	case "web.img":
		expiredTime = settings.Web.CacheMinutes
	}

	TaskCounter.Mu.RLock()
	cachedTask, ok := TaskQueue.Queue[filter]
	TaskCounter.Mu.RUnlock()

	if ok {
		for cachedTask.Running {
			time.Sleep(time.Second)
			TaskCounter.Mu.RLock()
			cachedTask = TaskQueue.Queue[filter]
			TaskCounter.Mu.RUnlock()
		}
		time.Sleep(time.Millisecond)
		if !cachedTask.Running &&
			time.Since(cachedTask.DB.Time) <= time.Duration(expiredTime)*time.Minute &&
			cachedTask.DB.Status == "success" {
			id := task.UUID
			task = cachedTask.DB
			task.UUID = id
			task.Time = time.Now()
			task.Temp = "Yes"
			return task, cachedTask.Return, true
		}
		os.Remove(ImgPath + cachedTask.DB.UUID + ".png")
	}

	task.Temp = "No"
	for {
		TaskCounter.Mu.RLock()
		value, exists := TaskCounter.Counter[filter.API]
		TaskCounter.Mu.RUnlock()
		if !exists || value <= 2 {
			break
		}
		time.Sleep(time.Second)
	}

	TaskCounter.Mu.RLock()
	item := TaskQueue.Queue[filter]
	TaskCounter.Mu.RUnlock()
	item.Running = true

	TaskCounter.Mu.Lock()
	TaskQueue.Queue[filter] = item
	TaskCounter.Counter[filter.API]++
	TaskCounter.Mu.Unlock()

	return task, GenerateResult{}, false
}

func finishTask(task *model.Task, filter TaskQueueFilter, result *GenerateResult, skipDB bool) {
	TaskCounter.Mu.Lock()
	TaskCounter.Counter[filter.API]--
	TaskCounter.Mu.Unlock()

	TaskCounter.Mu.RLock()
	item := TaskQueue.Queue[filter]
	TaskCounter.Mu.RUnlock()

	if !util.PngChecker(ImgPath+task.UUID+".png") && task.Temp != "Yes" && task.Status == "success" {
		task.Status = "failed"
		task.Return = "Invalid Image File"
		result.URL = "download?img=empty"
	}

	item.Running = false
	if task.Status == "success" {
		item.DB = *task
		item.Return = *result
	}

	TaskCounter.Mu.Lock()
	TaskQueue.Queue[filter] = item
	TaskCounter.Mu.Unlock()

	SaveTask(*task, skipDB)
}
