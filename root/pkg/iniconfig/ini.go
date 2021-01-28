package iniconfig

import (
	"fmt"
	"sync"

	"github.com/go-ini/ini"
)

var iniFile *ini.File
var base *ini.Section
var thisSection *ini.Section
var loadOnce sync.Once
var appType string
var appId int32

func LoadINIConfig(app_type string, app_id int, path string) (err error) {
	loadOnce.Do(func() {
		appType = app_type
		appId = int32(app_id)

		iniFile, err = ini.Load(path)
		if err != nil {
			return
		}
		iniFile.BlockMode = false

		base, err = iniFile.GetSection("base")
		if err != nil {
			return
		}
		thisSection, err = iniFile.GetSection(fmt.Sprintf("%v_%v", app_type, app_id))
	})
	return err
}

func findKey(key string, section *ini.Section) *ini.Key {
	if base.HasKey(key) {
		return base.Key(key)
	}

	if section != nil && section.HasKey(key) {
		return section.Key(key)
	}
	return nil
}

func AppId() int32    { return appId }
func AppType() string { return appType }

func Int32(key string) int32 {
	keyValue := findKey(key, thisSection)
	if keyValue == nil {
		return 0
	}
	value, _ := keyValue.Int()
	return int32(value)
}

func Int64(key string) int64 {
	keyValue := findKey(key, thisSection)
	if keyValue == nil {
		return 0
	}
	value, _ := keyValue.Int()
	return int64(value)
}

func String(key string) string {
	keyValue := findKey(key, thisSection)
	if keyValue == nil {
		return ""
	}
	return keyValue.String()
}

func Bool(key string) bool {
	keyValue := findKey(key, thisSection)
	if keyValue == nil {
		return false
	}
	value, _ := keyValue.Bool()
	return value
}

func KeyCheck(keys ...string) (missingKeys []string) {
	for _, key := range keys {
		if findKey(key, thisSection) == nil {
			missingKeys = append(missingKeys, key)
		}
	}
	return
}
