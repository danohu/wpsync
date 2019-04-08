package main

import (
	"fmt"
	"os"

	"github.com/ttacon/chalk"
)

type Logger struct {
	DebugLevel bool
	Quiet      bool
}

func (l Logger) Debug(a ...interface{}) {
	if l.DebugLevel {
		fmt.Println(chalk.Cyan, a, chalk.Reset)
	}
}

func (l Logger) Info(a ...interface{}) {
	if !l.Quiet {
		fmt.Println(chalk.Green, a, chalk.Reset)
	}
}

func (l Logger) Warn(a ...interface{}) {
	fmt.Println(chalk.Yellow, a, chalk.Reset)
}

func (l Logger) Fatal(a ...interface{}) {
	fmt.Println(chalk.Red, a, chalk.Reset)
	os.Exit(1)
}
