package configs

import "os"

var (
	Profile  string
	Port     string
	Test     bool
	LogLevel string
	Apis     string
)

type Vars struct {
	Profile string
}

func NewVars() *Vars {
	return &Vars{
		Profile: Profile,
	}
}

func init() {
	Profile = os.Getenv("env")
	if len(Profile) == 0 {
		Profile = "dev"
	}
}
