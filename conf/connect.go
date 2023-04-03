package conf

import (
	"gorm.io/gorm"
)

func (app *BaseApp) GetDB() *gorm.DB {
	if !app.initialized {
		err := app.Initialize()
		if err != nil {
			panic(err)
		}
	}
	//// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", app.Config.DBUser, app.Config.DBPass, app.Config.DBHost, app.Config.DBPort, app.Config.DBName)
	//database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	panic("Could not connect to the database!")
	//}
	return GetDB()
}
