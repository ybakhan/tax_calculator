package main

type Config struct {
	Server struct {
		interviewServerHost string `yaml:"interviewServerHost"`
		InterviewServerPort string `yaml:"interviewServerPort"`
	} `yaml:"server"`
}
