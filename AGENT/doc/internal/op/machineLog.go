package op

import (
	"fmt"
	"strconv"
	"time"
	"zhongxin/cmd/flags"
	"zhongxin/internal/conf"
	_error "zhongxin/internal/error"
	"zhongxin/internal/model"
	"zhongxin/util"

	"github.com/go-co-op/gocron/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func InitTask() {
	var err error
	scheduler, err = gocron.NewScheduler()
	if err != nil {
		util.Log.Error(err.Error())
		return
	}

	job, err = scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(uint(conf.Conf.RemoteDatabase.SyncTimeDay.Hour), uint(conf.Conf.RemoteDatabase.SyncTimeDay.Minute), 0),
				gocron.NewAtTime(uint(conf.Conf.RemoteDatabase.SyncTimeNight.Hour), uint(conf.Conf.RemoteDatabase.SyncTimeNight.Minute), 0),
			),
		),
		gocron.NewTask(func() {
			util.Log.Infof("Running machine log sync task at %s", util.TimeNow().Format("2006-01-02 15:04:05"))
			if err, _ := SyncMachineLog(); err != nil {
				util.Log.Error("failed to sync machine log:", err.Error())
			}
		}),
	)

	scheduler.Start()
}

func SyncMachineLog() (error, int) {
	if flags.NoRemote {
		fmt.Println("Skip sync machine log because No Remote DB is set")
		return nil, 0
	}
	pasIDs, err := db.GetAllMachinePasID()
	if err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}

	pasToMacDBs, _, err := remoteDB.GetAllMachinePasID()
	if err != nil {
		return errors.WithStack(err), _error.ConvertGormError(err)
	}
	pasToMac := make(map[int]int)
	for _, pasToMacDB := range pasToMacDBs {
		pasToMac[pasToMacDB.PasID] = pasToMacDB.MacID
	}

	for _, pasID := range pasIDs {
		pasIDInt, _ := strconv.Atoi(pasID)
		macID := pasToMac[pasIDInt]
		if flags.Verbose {
			fmt.Printf("Syncing machine log for PasID: %s, MacID is %d\n", pasID, macID)
		}
		machine, _, err := db.GetMachineByID(pasID)
		if err != nil {
			return errors.WithStack(err), _error.ConvertGormError(err)
		}
		if machine.LastSyncTime == 0 || machine.LastSyncTime > util.TimeNow().Unix() {
			machine.LastSyncTime = util.TimeNow().Unix()
		}
		if flags.Verbose {
			fmt.Println("Machine Info: ", machine)
		}
		machineLogs, _, err := remoteDB.GetMachineONLogByFilter([]clause.Expression{
			clause.Eq{Column: "MacID", Value: macID},
			clause.Gte{Column: "BeginTime", Value: util.NaiveLocalToNaiveUTC(util.
				FindNthPreviousTime(8, 0, 0, conf.MachineLogSyncLength, machine.LastSyncTime))},
			clause.Lte{Column: "BeginTime", Value: util.NaiveLocalToNaiveUTC(util.TimeNow())},
			//clause.Gte{Column: "BeginTime", Value: time.Date(2025, 8, 30, 20, 0, 0, 0, time.UTC)},
			//clause.Lte{Column: "BeginTime", Value: time.Date(2025, 9, 5, 20, 0, 0, 0, time.UTC)},
			// 2025-8-31 8:00:00 +0800 CST - 1756598400
			// 2025-8-31 20:00:00 +0800 CST -1756641600
			// 2025-9-1 8:00:00 +0800 CST -  1756684800
			// 2025-9-1 20:00:00 +0800 CST - 1756728000
			// 2025-9-2  8:00:00 +0800 CST - 1756771200
			// 2025-9-2 20:00:00 +0800 CST - 1756814400
			// 2025-9-3  8:00:00 +0800 CST - 1756857600
			// 2025-9-3 20:00:00 +0800 CST - 1756900800
			// 2025-9-4  8:00:00 +0800 CST - 1756944000
			// 2025-9-4 20:00:00 +0800 CST - 1756987200
			// 2025-9-5  8:00:00 +0800 CST - 1757030400
		})
		if err != nil {
			return errors.WithStack(err), _error.ConvertGormError(err)
		}
		if flags.Verbose {
			fmt.Println("Fetched machine logs: ", machineLogs)
		}

		machine.LastSyncTime = util.TimeNow().Unix()
		if err, errCode := UpdateMachine(machine); err != nil {
			return errors.WithStack(err), errCode
		}
		workPeriods, exist, err := db.GetWorkPeriodsByClauses([]clause.Expression{
			clause.Eq{Column: "machineID", Value: pasID},
			clause.Eq{Column: "isMachineDataMerged", Value: false},
		})
		if err != nil {
			return errors.WithStack(err), _error.ConvertGormError(err)
		}
		if flags.Verbose {
			fmt.Println("Fetched work periods: ", workPeriods)
		}

		workPeriodStartTime, workPeriodEndTime := []int64{}, []int64{}
		machinePeriodStartTime, machinePeriodEndTime := []int64{}, []int64{}
		for index, workPeriod := range workPeriods {
			workPeriodStartTime = append(workPeriodStartTime, workPeriod.StartTime)
			workPeriodEndTime = append(workPeriodEndTime, workPeriod.EndTime)
			if workPeriodEndTime[index] == 0 {
				workPeriodEndTime[index] = util.TimeNow().Unix()
			}
		}
		for _, machineLog := range machineLogs {
			machinePeriodStartTime = append(machinePeriodStartTime, util.NaiveUTCToNaiveLocal(machineLog.BeginTime).Unix())
			machinePeriodEndTime = append(machinePeriodEndTime, util.NaiveUTCToNaiveLocal(machineLog.EndTime).Unix())
		}
		if !exist || len(workPeriods) == 0 || len(machineLogs) == 0 {
			continue
		}
		if flags.Verbose {
			fmt.Println("Work Period Start Times: ", workPeriodStartTime)
			fmt.Println("Work Period End Times: ", workPeriodEndTime)
			fmt.Println("Machine Period Start Times: ", machinePeriodStartTime)
			fmt.Println("Machine Period End Times: ", machinePeriodEndTime)
			fmt.Println()
		}

		mergedWorkPeriodStartTime, mergedWorkPeriodEndTime := util.MergeMachineLogToWorkPeriods(workPeriodStartTime, workPeriodEndTime, machinePeriodStartTime, machinePeriodEndTime)
		if flags.Dev {
			util.Log.Infof("Syncing machine log for PasID: %s, MacID is %d\n", pasID, macID)
			util.Log.Infof("Fetch %d work periods and %d machine logs\n", len(workPeriods), len(machineLogs))
			util.Log.Infof("Get %d merged work periods\n", len(mergedWorkPeriodStartTime))
			util.Log.Infoln()
		}
		if flags.Verbose {
			for index, startTime := range mergedWorkPeriodStartTime {
				fmt.Println("For WP ", index)
				fmt.Println("Work Period Start Times: ", workPeriodStartTime[index])
				fmt.Println("Work Period End Times: ", workPeriodEndTime[index])
				fmt.Println("Merged Work Period Start Times: ", startTime)
				fmt.Println("Merged Work Period End Times: ", mergedWorkPeriodEndTime[index])
				fmt.Println()
			}
		}

		for index, workPeriod := range workPeriods {
			mergedMachinePeriods := make([]model.MachinePeriod, len(mergedWorkPeriodStartTime[index]))
			for i, startTime := range mergedWorkPeriodStartTime[index] {
				mergedMachinePeriods[i] = model.MachinePeriod{
					StartTime: startTime,
					EndTime:   mergedWorkPeriodEndTime[index][i],
				}
			}
			workPeriod.MachineONPeriods = append(workPeriod.MachineONPeriods, mergedMachinePeriods...)
			if time.Unix(workPeriod.EndTime, 0).Before(util.FindNthPreviousTime(8, 0, 0, 1, machine.LastSyncTime)) && workPeriod.EndTime != 0 {
				uniqueStartTime, uniqueEndTime := []int64{}, []int64{}
				for _, machineONPeriod := range workPeriod.MachineONPeriods {
					uniqueStartTime = append(uniqueStartTime, machineONPeriod.StartTime)
					uniqueEndTime = append(uniqueEndTime, machineONPeriod.EndTime)
				}
				uniqueStartTime, uniqueEndTime = util.MergeInWorkPeriod(uniqueStartTime, uniqueEndTime)
				mergedMachinePeriods = make([]model.MachinePeriod, len(uniqueStartTime))
				validTimeSeconds := int64(0)
				for i, startTime := range uniqueStartTime {
					mergedMachinePeriods[i] = model.MachinePeriod{
						StartTime: startTime,
						EndTime:   uniqueEndTime[i],
					}
					validTimeSeconds += uniqueEndTime[i] - startTime
				}
				workPeriod.MachineONPeriods = mergedMachinePeriods
				workPeriod.ValidTimeSeconds = validTimeSeconds
				workPeriod.IsMachineDataMerged = true
			}
			if err := db.UpdateWorkPeriod(workPeriod); err != nil {
				return errors.WithStack(err), _error.ConvertGormError(err)
			}
		}
		if err := db.UpdateMachine(machine); err != nil {
			return errors.WithStack(err), _error.ConvertGormError(err)
		}
	}
	return nil, 0
}

func TestRemoteConnection(pasID int) error {
	machine, _, err := db.GetMachineByID(strconv.Itoa(pasID))
	if err != nil {
		return errors.WithStack(err)
	}

	pasToMacDBs, _, err := remoteDB.GetAllMachinePasID()
	if err != nil {
		return errors.WithStack(err)
	}
	pasToMac := make(map[int]int)
	for _, pasToMacDB := range pasToMacDBs {
		pasToMac[pasToMacDB.PasID] = pasToMacDB.MacID
	}
	macID := pasToMac[pasID]
	fmt.Println("Machine Info: ", machine)
	fmt.Println("MacID: ", macID)

	machineLogs, _, err := remoteDB.GetMachineONLogByFilter([]clause.Expression{
		clause.Eq{Column: "MacID", Value: macID},
		clause.Gte{Column: "BeginTime", Value: util.NaiveLocalToNaiveUTC(util.
			FindNthPreviousTime(8, 0, 0, conf.MachineLogSyncLength, util.TimeNow().Unix()))},
		clause.Lte{Column: "BeginTime", Value: util.NaiveLocalToNaiveUTC(util.TimeNow())},
		//clause.Gte{Column: "BeginTime", Value: time.Date(2025, 8, 30, 20, 0, 0, 0, time.UTC)},
		//clause.Lte{Column: "BeginTime", Value: time.Date(2025, 9, 5, 20, 0, 0, 0, time.UTC)},
		// 2025-8-31 8:00:00 +0800 CST - 1756598400
		// 2025-8-31 20:00:00 +0800 CST -1756641600
		// 2025-9-1 8:00:00 +0800 CST -  1756684800
		// 2025-9-1 20:00:00 +0800 CST - 1756728000
		// 2025-9-2  8:00:00 +0800 CST - 1756771200
		// 2025-9-2 20:00:00 +0800 CST - 1756814400
		// 2025-9-3  8:00:00 +0800 CST - 1756857600
		// 2025-9-3 20:00:00 +0800 CST - 1756900800
		// 2025-9-4  8:00:00 +0800 CST - 1756944000
		// 2025-9-4 20:00:00 +0800 CST - 1756987200
		// 2025-9-5  8:00:00 +0800 CST - 1757030400
	})
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println(machineLogs)
	return nil
}

func ResetInvalidWorkPeriods() error {
	workPeriods, _, err := db.GetWorkPeriodsByClauses([]clause.Expression{
		clause.Eq{Column: "isMachineDataMerged", Value: true},
		clause.Eq{Column: "validTimeSeconds", Value: 0},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	for _, workPeriod := range workPeriods {
		workPeriod.IsMachineDataMerged = false
		if err := db.UpdateWorkPeriod(workPeriod); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func ResetAllWorkPeriods() error {
	workPeriods, _, err := db.GetWorkPeriodsByClauses([]clause.Expression{
		clause.Eq{Column: "isMachineDataMerged", Value: true},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	for _, workPeriod := range workPeriods {
		workPeriod.IsMachineDataMerged = false
		if err := db.UpdateWorkPeriod(workPeriod); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
