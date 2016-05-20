package logger

import (
	"fmt"
	"github.com/ttacon/chalk"
)

func Log(moduleName string, message string) {
	log(chalk.ResetColor, "LOG", moduleName, message)
}

func Info(moduleName string, message string) {
	log(chalk.Blue, "INFO", moduleName, message)
}

func Warn(moduleName string, message string) {
	log(chalk.Magenta, "WARN", moduleName, message)
}

// Just alias, sometimes typing this
func Warning(moduleName string, message string) {
	log(chalk.Magenta, "WARN", moduleName, message)
}

func Error(moduleName string, message string) {
	log(chalk.Red, "ERROR", moduleName, message)
}

func log(color chalk.Color, level, moduleName string, message string) {
	fmt.Println(color, "[" + moduleName + "][" + level + "]", message, chalk.Reset)
}