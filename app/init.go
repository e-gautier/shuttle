package app

import (
	"github.com/revel/revel"
	"database/sql"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"shuttle/app/models"
	"github.com/revel/revel/logger"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

const SQLITE_EXT = ".sqlite3"
var DB *gorp.DbMap
var RootLog = logger.New()
var AppLog = RootLog.New("module", "app")

func InitDB() {
	driver, _ := revel.Config.String("db.driver")
	dbName, _ := revel.Config.String("db.name")
	db, err := sql.Open(driver, dbName + SQLITE_EXT)
	CheckErr(err)
	DB = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	DB.AddTableWithName(models.Resource{}, "resources").SetKeys(true, "Id").AddIndex("long_url", "string", []string{"long_url"}).SetUnique(true)
	err = DB.CreateTablesIfNotExists()
	CheckErr(err)
	DB.CreateIndex()
	// TODO change for a gorp auto_increment 1000
	_, err = DB.Exec("INSERT OR IGNORE INTO resources(id, long_url) values (999, '')")
	CheckErr(err)
	_, err = DB.Exec("DELETE FROM resources WHERE id = 999")
	CheckErr(err)
}

func CheckErr(err error) {
	if err != nil {
		AppLog.Fatalf(err.Error())
	}
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
