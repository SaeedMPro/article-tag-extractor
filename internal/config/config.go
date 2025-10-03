package config

type Config struct {
	Database Database
	Server   Server
}

type Database struct {
	URI        string
	DBName     string
}

type Server struct {
	GRPCPort string
}
