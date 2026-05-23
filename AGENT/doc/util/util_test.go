package util

import (
	"fmt"
	"testing"
	//"time"
)

func TestMergeMachineLogToWorkPeriods(t *testing.T) {
	//workPeriodStartTime := []int64{
	//	1, 8, 27, 15, 18, 12, 5,
	//}
	//workPeriodEndTime := []int64{
	//	3, 11, 31, 16, 25, 13, 6,
	//}
	//
	//machinePeriodStartTime := []int64{
	//	1, 7, 30, 22, 29, 14, 4, 21, 17, 24,
	//}
	//machinePeriodEndTime := []int64{
	//	2, 10, 32, 26, 33, 19, 9, 23, 20, 28,
	//}
	workPeriodStartTime := []int64{
		1760143474,
		1760185994,
		1760186208,
		1760227833,
		1760314127,
		1760325992,
		1760326762,
		1760326817,
		1760326909,
	}
	workPeriodEndTime := []int64{
		1760347076,
		1760186001,
		1760186217,
		1760270382,
		1760325930,
		1760326485,
		1760326773,
		1760326849,
		1760347076,
	}
	machinePeriodStartTime := []int64{
		1760140802,
		1760143670,
		1760147014,
		1760147796,
		1760184001,
		1760187325,
		1760190189,
		1760218079,
		1760227202,
		1760240370,
		1760243554,
		1760250915,
		1760251133,
		1760258450,
		1760258707,
		1760262468,
		1760268485,
		1760270401,
		1760271679,
		1760293304,
		1760302833,
		1760307558,
		1760310329,
		1760311657,
		1760313602,
		1760326317,
		1760326348,
		1760326647,
		1760333126,
		1760333539,
		1760337650,
	}
	machinePeriodEndTime := []int64{
		1760142889,
		1760145845,
		1760147211,
		1760184000,
		1760187149,
		1760189598,
		1760210441,
		1760227201,
		1760238471,
		1760242813,
		1760250006,
		1760250919,
		1760256645,
		1760258488,
		1760261242,
		1760268041,
		1760270400,
		1760270570,
		1760286952,
		1760298364,
		1760303998,
		1760309682,
		1760310958,
		1760313601,
		1760326151,
		1760326335,
		1760326629,
		1760327139,
		1760333138,
		1760334544,
		1760338010,
	}

	retStartTime, retEndTime := MergeMachineLogToWorkPeriods(workPeriodStartTime, workPeriodEndTime, machinePeriodStartTime, machinePeriodEndTime)
	for index, startTime := range retStartTime {
		fmt.Println(index + 1)
		fmt.Println(startTime)
		fmt.Println(retEndTime[index])
		fmt.Println()
	}
}

func TestMergeInWorkPeriod(t *testing.T) {
	startTime := []int64{
		1, 2, 3, 6, 8, 11, 12, 14, 17, 19, 9,
	}
	endTime := []int64{
		5, 4, 7, 10, 10, 13, 15, 16, 20, 20, 10,
	}
	mergedStartTime, mergedEndTime := MergeInWorkPeriod(startTime, endTime)
	fmt.Println(mergedStartTime)
	fmt.Println(mergedEndTime)
}

func TestFindNthPreviousTime(t *testing.T) {

	fmt.Println(FindNthPreviousTime(10, 0, 2, 1, TimeNow().Unix()))
}

func TestRemoveLeadingZeros(t *testing.T) {
	fmt.Println(RemoveLeadingZeros("099"))
}
