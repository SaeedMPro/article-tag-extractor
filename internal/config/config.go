package config

type Config struct {
	Database Database
	Server   Server
}

type Database struct {
	URL        string
	DBName     string
	Collection string
}

type Server struct {
	Port string
}
