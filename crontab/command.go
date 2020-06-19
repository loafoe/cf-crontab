package crontab

// Command implement job types
type Command interface {
	Run()
}