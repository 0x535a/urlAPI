package database

import (
	"encoding/json"
	"sync"
	"urlAPI/util"

	"github.com/pkg/errors"
)

const appSettingsKey = "app"

type appSettingsStore struct {
	mu       sync.RWMutex
	settings util.AppSettings
}

var SettingsStore = appSettingsStore{}

var SkipAppSettingsSync bool

func (store *appSettingsStore) Get() util.AppSettings {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.settings
}

func (store *appSettingsStore) Replace(settings util.AppSettings) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.settings = settings
}

func initAppSettings() error {
	setting := AppSetting{Key: appSettingsKey}
	list, err := setting.Read()
	if err == nil && len(list.AppSettingList) > 0 {
		var settings util.AppSettings
		if err := json.Unmarshal([]byte(list.AppSettingList[0].Value), &settings); err != nil {
			return errors.WithStack(err)
		}
		SettingsStore.Replace(util.NormalizeSettings(settings))
		return nil
	}

	settings := util.MigrateLegacySettings(SettingMap)
	return SaveAppSettings(settings)
}

func SaveAppSettings(settings util.AppSettings) error {
	settings = util.NormalizeSettings(settings)
	value, err := json.Marshal(settings)
	if err != nil {
		return errors.WithStack(err)
	}
	setting := AppSetting{
		Key:     appSettingsKey,
		Version: settings.Version,
		Value:   string(value),
	}
	if err := setting.Update(); err != nil {
		return errors.WithStack(err)
	}
	SettingsStore.Replace(settings)
	return nil
}

func syncAppSettingsFromLegacy() error {
	return SaveAppSettings(util.MigrateLegacySettings(SettingMap))
}

func (setting *AppSetting) Create() error {
	return errors.WithStack(db.Create(setting).Error)
}

func (setting *AppSetting) Update() error {
	return errors.WithStack(db.Save(setting).Error)
}

func (setting *AppSetting) Read() (*DBList, error) {
	var settings []AppSetting
	err := db.Where("key = ?", setting.Key).Find(&settings).Error
	if len(settings) == 0 {
		err = errors.WithStack(errors.New("AppSetting not found"))
	}
	ret := DBList{
		AppSettingList: settings,
	}
	return &ret, errors.WithStack(err)
}

func (setting *AppSetting) Delete() error {
	return errors.WithStack(db.Delete(setting).Error)
}
