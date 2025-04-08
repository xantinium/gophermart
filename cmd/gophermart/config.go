package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/xantinium/gophermart/internal/tools"
)

type appArgs struct {
	Addr            string
	DatabaseConnStr string
	AccrualHost     string
}

func parseAppArgs() appArgs {
	address := new(netAddress)
	flag.Var(address, "a", "address of gophermart server in form <host:port>")
	databaseConnStr := flag.String("d", "", "connection string for postgresql")
	accrualHost := flag.String("r", "", "base url of accural system")

	flag.Parse()

	args := appArgs{
		Addr:            address.String(),
		DatabaseConnStr: *databaseConnStr,
		AccrualHost:     *accrualHost,
	}

	envArgs := parseAppArgsFromEnv()

	if envArgs.Addr.Exists {
		args.Addr = envArgs.Addr.Value
	}
	if envArgs.DatabaseConnStr.Exists {
		args.DatabaseConnStr = envArgs.DatabaseConnStr.Value
	}
	if envArgs.AccrualHost.Exists {
		args.AccrualHost = envArgs.AccrualHost.Value
	}

	return args
}

type appEnvArgs struct {
	Addr            tools.StrEnvVar
	DatabaseConnStr tools.StrEnvVar
	AccrualHost     tools.StrEnvVar
}

func parseAppArgsFromEnv() appEnvArgs {
	return appEnvArgs{
		Addr:            tools.GetStrFromEnv("RUN_ADDRESS"),
		DatabaseConnStr: tools.GetStrFromEnv("DATABASE_URI"),
		AccrualHost:     tools.GetStrFromEnv("ACCRUAL_SYSTEM_ADDRESS"),
	}
}

type netAddress struct {
	Host   string
	Port   int
	Parsed bool
}

func (a netAddress) String() string {
	if a.Parsed {
		return fmt.Sprintf("%s:%d", a.Host, a.Port)
	}

	return "localhost:8080"
}

func (a *netAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("invalid address format")
	}

	host := hp[0]

	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}

	a.Host = host
	a.Port = port
	a.Parsed = true

	return nil
}
