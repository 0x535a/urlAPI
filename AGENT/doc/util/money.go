package util

import (
	"fmt"
	"sort"
)

type timePeriod struct {
	id        int
	startTime int64
	endTime   int64
}

// 将机器的运行时间段合并到工作时间段中
func MergeMachineLogToWorkPeriods(workPeriodStartTime, workPeriodEndTime, machinePeriodStartTime, machinePeriodEndTime []int64) ([][]int64, [][]int64) {
	workPeriods := make([]timePeriod, len(workPeriodStartTime))
	for index, startTime := range workPeriodStartTime {
		workPeriods[index].id = index
		workPeriods[index].startTime = startTime
		workPeriods[index].endTime = workPeriodEndTime[index]
	}
	machinePeriods := make([]timePeriod, len(machinePeriodStartTime))
	for index, startTime := range machinePeriodStartTime {
		machinePeriods[index].id = index
		machinePeriods[index].startTime = startTime
		machinePeriods[index].endTime = machinePeriodEndTime[index]
	}

	sort.Slice(machinePeriods[:], func(i, j int) bool {
		if machinePeriods[i].startTime == machinePeriods[j].startTime {
			return machinePeriods[i].endTime < machinePeriods[j].endTime
		}
		return machinePeriods[i].startTime < machinePeriods[j].startTime
	})

	fmt.Println(workPeriods)
	fmt.Println(machinePeriods)

	retStartTime := make([][]int64, len(workPeriods))
	retEndTime := make([][]int64, len(workPeriods))
	var mergeStartTime, mergeEndTime int64
	for _, workPeriod := range workPeriods {
		for _, machinePeriod := range machinePeriods {
			switch {
			case machinePeriod.endTime > workPeriod.startTime &&
				machinePeriod.endTime <= workPeriod.endTime &&
				machinePeriod.startTime <= workPeriod.startTime:
				//case 2
				mergeStartTime = workPeriod.startTime
				mergeEndTime = machinePeriod.endTime
			case machinePeriod.startTime > workPeriod.startTime &&
				machinePeriod.endTime <= workPeriod.endTime:
				//case 3
				mergeStartTime = machinePeriod.startTime
				mergeEndTime = machinePeriod.endTime
			case machinePeriod.startTime > workPeriod.startTime &&
				machinePeriod.startTime <= workPeriod.endTime &&
				machinePeriod.endTime > workPeriod.endTime:
				//case 4
				mergeStartTime = machinePeriod.startTime
				mergeEndTime = workPeriod.endTime
			case machinePeriod.startTime <= workPeriod.startTime &&
				machinePeriod.endTime > workPeriod.endTime:
				//case 6
				mergeStartTime = workPeriod.startTime
				mergeEndTime = workPeriod.endTime
			default:
				continue
			}

			retStartTime[workPeriod.id] = append(retStartTime[workPeriod.id], mergeStartTime)
			retEndTime[workPeriod.id] = append(retEndTime[workPeriod.id], mergeEndTime)
		}
	}
	return retStartTime, retEndTime
}

// 合并，去重
func MergeInWorkPeriod(startTime, endTime []int64) ([]int64, []int64) {
	if len(startTime) == 0 || len(endTime) == 0 || len(startTime) != len(endTime) {
		return []int64{}, []int64{}
	}
	workPeriods := make([]timePeriod, len(startTime))
	for index, workPeriod := range startTime {
		workPeriods[index].startTime = workPeriod
		workPeriods[index].endTime = endTime[index]
	}
	sort.Slice(workPeriods[:], func(i, j int) bool {
		if workPeriods[i].startTime == workPeriods[j].startTime {
			return workPeriods[i].endTime < workPeriods[j].endTime
		}
		return workPeriods[i].startTime < workPeriods[j].startTime
	})

	retStartTime := make([]int64, 0)
	retEndTime := make([]int64, 0)
	currentStartTime := workPeriods[0].startTime
	currentEndTime := workPeriods[0].endTime

	for _, workPeriod := range workPeriods {
		if workPeriod.startTime > currentEndTime {
			retStartTime = append(retStartTime, currentStartTime)
			retEndTime = append(retEndTime, currentEndTime)
			currentStartTime = workPeriod.startTime
			currentEndTime = workPeriod.endTime
		}
		currentStartTime = min(currentStartTime, workPeriod.startTime)
		currentEndTime = max(currentEndTime, workPeriod.endTime)
	}
	retStartTime = append(retStartTime, currentStartTime)
	retEndTime = append(retEndTime, currentEndTime)
	return retStartTime, retEndTime
}
