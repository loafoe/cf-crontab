package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/philips-labs/cf-crontab/crontab"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strconv"
)

type ErrResponse struct {
	Message string `json:"message"`
	Code int `json:"code"`
}
func entriesDeleteHandler(state *crontab.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		stringID := c.Param("entryID")
		entryID, err := strconv.Atoi(stringID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Message: "invalid entry",
				Code: http.StatusBadRequest,
			})
		}
		err = state.DeleteEntry(entryID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Message: err.Error(),
				Code: http.StatusBadRequest,
			})
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
	viper.SetEnvPrefix("cf_crontab")
	viper.SetDefault("secret", "")
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
	secret := viper.GetString("secret")
	if secret == "" {
		fmt.Printf("secret is required\n")
		return
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(HSDPValidator(crontab.SharedKey, secret))
	e.GET("/entries", entriesGetHandler(state))
	e.POST("/entries", entriesPostHandler(state))
	e.DELETE("/entries/:entryID", entriesDeleteHandler(state))
	usePort := os.Getenv("PORT")
	if usePort == "" {
		usePort = "8080"
	}
	_ = e.Start(":" + usePort)
}
