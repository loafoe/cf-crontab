package server

import (
	"crypto/subtle"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/philips-labs/cf-crontab/crontab"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strconv"
)

// authCheck verifies basic auth
func authCheck(username, password string, c echo.Context) (bool, error) {
	if subtle.ConstantTimeCompare([]byte(username), []byte(username)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(password)) == 1 {
		return true, nil
	}
	return false, nil
}

func entriesDeleteHandler(state *crontab.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		stringID := c.Param("entryID")
		entryID, err := strconv.Atoi(stringID)
		if err != nil {
			return err
		}
		err = state.DeleteEntry(entryID)
		if err != nil {
			return err
		}
		return c.String(http.StatusNoContent, "")
	}
}

func entriesGetHandler(state *crontab.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, state.Entries(), "  ")
	}
}

func entriesPostHandler(state *crontab.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var newEntries []crontab.Task
		if err := c.Bind(&newEntries); err != nil {
			return err
		}
		state.AddEntries(newEntries)
		return c.JSONPretty(http.StatusOK, newEntries, "  ")
	}
}

// Start starts the cf-crontab server
func Start() {
	// Config
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.SetEnvPrefix("cf-crontab")
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()

	// Cron
	state := crontab.NewState()

	entries, err := crontab.LoadFromEnv()
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		return
	}
	state.AddEntries(entries)
	state.StartCron()

	// Echo
	e := echo.New()
	e.GET("/entries", entriesGetHandler(state))
	e.POST("/entries", entriesPostHandler(state))
	e.DELETE("/entries/:entryID", entriesDeleteHandler(state))
	usePort := os.Getenv("PORT")
	if usePort == "" {
		usePort = "8080"
	}
	_ = e.Start(":" + usePort)
}
