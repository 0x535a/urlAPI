package cmd

import (
	"github.com/spf13/cobra"
	"zhongxin/internal/database"
	"zhongxin/internal/model"
	"zhongxin/util"
)

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "admin management",
	Run: func(cmd *cobra.Command, args []string) {
		adminCmd(args)
	},
}
var name, phone, wxid string

func adminCmd(args []string) {
	Init()
	defer Release()

	db := database.GetLocalDB()
	admin := model.User{}
	if len(args) < 1 {
		util.Log.Error("please provide action: add, edit, delete")
		return
	}
	if name == "" {
		util.Log.Error("please provide name, at least")
		return
	}

	var err error
	switch args[0] {
	case "add":
		admin.Type = "admin"
		admin.ID = util.NewUUID()
		admin.Name = name
		admin.Phone = phone
		admin.WXID = wxid
		if err = db.NewUser(admin); err != nil {
			util.Log.Error(err)
			return
		}
	case "edit":
		if admin, _, err = db.GetUserByName(name); err != nil {
			util.Log.Error(err)
			return
		}
		if phone != "" {
			admin.Phone = phone
		}
		if wxid != "" {
			admin.WXID = wxid
		}
		if err = db.UpdateUser(admin); err != nil {
			util.Log.Error(err)
			return
		}
	case "del":
		if err = db.DeleteUserByName(name); err != nil {
			util.Log.Error(err)
			return
		}
	case "get":
		if admin, _, err = db.GetUserByName(name); err != nil {
			util.Log.Error(err)
			return
		}
	default:
		util.Log.Error("unknown action: ", args[0])
	}
	util.Log.Infof("operation success: %s", args[0])
	util.Log.Info(admin)
}

func init() {
	RootCmd.AddCommand(AdminCmd)
}
