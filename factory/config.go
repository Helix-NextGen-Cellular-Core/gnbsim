// SPDX-FileCopyrightText: 2022 Great Software Laboratory Pvt. Ltd
// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0

/*
 * gNBSim Configuration Factory
 */

package factory

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/omec-project/gnbsim/common"
	gnbctx "github.com/omec-project/gnbsim/gnodeb/context"
	profctx "github.com/omec-project/gnbsim/profile/context"

	"sort"
)

const (
	GNBSIM_EXPECTED_CONFIG_VERSION string = "1.0.0"
	GNBSIM_DEFAULT_CONFIG_PATH            = "/gnbsim/config/gnb.conf"
)

type Config struct {
	Info          *Info          `yaml:"info"`
	Configuration *Configuration `yaml:"configuration"`
	Logger        *Logger        `yaml:"logger"`
}

type Info struct {
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

type Configuration struct {
	Gnbs                     map[string]*gnbctx.GNodeB   `yaml:"gnbs"`
	CustomProfiles           map[string]*profctx.Profile `yaml:"customProfiles"`
	Profiles                 []*profctx.Profile          `yaml:"profiles"`
	SingleInterface          bool                        `yaml:"singleInterface"`
	ExecInParallel           bool                        `yaml:"execInParallel"`
	StartUeatOnce            bool                        `yaml:"startUeatOnce"`
	StartImsiDelay           int                         `yaml:"startImsiDelay"`
	RealTrace                bool                        `yaml:"realTrace"`
	RealTraceFileSR          string                      `yaml:"realTraceFileSR"`
	RealTraceFileS1          string                      `yaml:"realTraceFileS1"`
	RealTraceSimulationTime  string                      `yaml:"realTraceSimulationTime"`
	PauseAtStart             bool                        `yaml:"pauseAtStart"` // if set to true it will not wait at beginning - inverted
	EventRetry               bool                        `yaml:"eventRetry"`
	EventRetryTime           int                         `yaml:"eventRetryTime"`
	Server                   HttpServer                  `yaml:"httpServer"`
	GoProfile                ProfileServer               `yaml:"goProfile"`
	RunConfigProfilesAtStart bool                        `yaml:"runConfigProfilesAtStart"`
}

type ProfileServer struct {
	Enable bool `yaml:"enable"`
	Port   int  `yaml:"port"`
}

type HttpServer struct {
	Enable bool   `yaml:"enable"`
	IpAddr string `yaml:"ipAddr"`
	Port   string `yaml:"port"`
}

type Logger struct {
	LogLevel string `yaml:"logLevel"`
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (c *Config) Validate() (err error) {

	fmt.Printf("c.Configuration.StartUeatOnce", c.Configuration.StartUeatOnce, "\n")
	fmt.Printf("c.Configuration.StartImsiDelay", c.Configuration.StartImsiDelay, "\n")
	fmt.Printf("c.Configuration.RealTrace", c.Configuration.RealTrace, "\n")
	fmt.Printf("c.Configuration.RealTraceFileSR", c.Configuration.RealTraceFileSR, "\n")
	fmt.Printf("c.Configuration.RealTraceFileS1", c.Configuration.RealTraceFileS1, "\n")
	fmt.Printf("c.Configuration.RealTraceSimulationTime", c.Configuration.RealTraceSimulationTime, "\n")
	fmt.Printf("c.Configuration.PauseAtStart if set to true it will not wait at beginning - inverted", c.Configuration.PauseAtStart, "\n")
	fmt.Printf("c.Configuration.EventRetry", c.Configuration.EventRetry, "\n")
	fmt.Printf("c.Configuration.EventRetryTime", c.Configuration.EventRetryTime, "\n")

	if c.Info == nil {
		return fmt.Errorf("Info field missing")
	}

	if c.Configuration == nil {
		return fmt.Errorf("Configuration field missing")
	}

	if len(c.Configuration.Gnbs) == 0 {
		return fmt.Errorf("no gnbs configured")
	}
	if c.Configuration.GoProfile.Enable == true {
		if c.Configuration.GoProfile.Port == 0 {
			c.Configuration.GoProfile.Port = 5000
		}
	}

	if c.Configuration.Server.IpAddr == "POD_IP" {
		c.Configuration.Server.IpAddr = os.Getenv("POD_IP")
	}

	if c.Configuration.SingleInterface == true {
		for _, gnb := range c.Configuration.Gnbs {
			if gnb.GnbN3Ip == "POD_IP" {
				gnb.GnbN3Ip = os.Getenv("POD_IP")
			}
		}
	}

	if len(c.Configuration.Profiles) == 0 {
		return fmt.Errorf("no profile information available")
	}

	keys := make([]int, 0, len(c.Configuration.CustomProfiles))
	if len(c.Configuration.CustomProfiles) != 0 {
		for key, v := range c.Configuration.CustomProfiles {
			it := v.Iterations
			v.PIterations = make(map[string]*profctx.PIterations)
			numKey, _ := strconv.Atoi(key[14:])
			keys = append(keys, numKey)
			for _, v1 := range it {
				if len(v1.Next) == 0 {
					v1.Next = "quit" // default value
				}
				PIter := &profctx.PIterations{Name: v1.Name, NextItr: v1.Next, Repeat: v1.Repeat}
				PIter.ProcMap = make(map[int]common.ProcedureType)
				PIter.WaitMap = make(map[int]int)
				PIter.WaitDistMap = make(map[int]*profctx.DelayDist)
				PIter.WaitMap[0] = 0
				if len(v1.First) > 0 {
					x := strings.Fields(v1.First)
					if c.Configuration.RealTrace == false {
						PIter.ProcMap[1] = common.GetProcId(x[0])
						PIter.WaitMap[1], err = strconv.Atoi(x[1])
					} else {
						PIter.SetIterationParam(1, x)
					}
				}
				if len(v1.Second) > 0 {
					x := strings.Fields(v1.Second)
					if c.Configuration.RealTrace == false {
						PIter.ProcMap[2] = common.GetProcId(x[0])
						PIter.WaitMap[2], err = strconv.Atoi(x[1])
					} else {
						PIter.SetIterationParam(2, x)
					}
				}
				if len(v1.Third) > 0 {
					x := strings.Fields(v1.Third)
					if c.Configuration.RealTrace == false {
						PIter.ProcMap[3] = common.GetProcId(x[0])
						PIter.WaitMap[3], err = strconv.Atoi(x[1])
					} else {
						PIter.SetIterationParam(3, x)
					}
				}
				if len(v1.Fourth) > 0 {
					x := strings.Fields(v1.Fourth)
					if c.Configuration.RealTrace == false {
						PIter.ProcMap[4] = common.GetProcId(x[0])
						PIter.WaitMap[4], err = strconv.Atoi(x[1])
					} else {
						PIter.SetIterationParam(4, x)
					}
				}
				if len(v1.Fifth) > 0 {
					x := strings.Fields(v1.Fifth)
					if c.Configuration.RealTrace == false {
						PIter.ProcMap[5] = common.GetProcId(x[0])
						PIter.WaitMap[5], err = strconv.Atoi(x[1])
					} else {
						PIter.SetIterationParam(5, x)
					}
				}
				if len(v1.Sixth) > 0 {
					x := strings.Fields(v1.Sixth)
					if c.Configuration.RealTrace == false {
						PIter.ProcMap[6] = common.GetProcId(x[0])
						PIter.WaitMap[6], err = strconv.Atoi(x[1])
					} else {
						PIter.SetIterationParam(6, x)
					}
				}
				if len(v1.Seventh) > 0 {
					x := strings.Fields(v1.Seventh)
					if c.Configuration.RealTrace == false {
						PIter.ProcMap[7] = common.GetProcId(x[0])
						PIter.WaitMap[7], err = strconv.Atoi(x[1])
					} else {
						PIter.SetIterationParam(1, x)
					}
				}
				v.PIterations[v1.Name] = PIter // add iterations in the custom profile
			}
		}

		sort.Ints(keys)

		for _, v := range keys {
			c.Configuration.Profiles = append(c.Configuration.Profiles,
				c.Configuration.CustomProfiles["customProfiles"+strconv.Itoa(v)])
		}

	}

	return nil
}

func (c *Configuration) GetGNodeB(name string) (*gnbctx.GNodeB, error) {
	var err error
	gnb, ok := c.Gnbs[name]
	if !ok {
		err = fmt.Errorf("no corresponding gNodeB found for:%v", name)
	}
	return gnb, err
}
