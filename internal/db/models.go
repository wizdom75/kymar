package db

// ConnParams holds database connection parameters
type ConnParams struct {
	DBType  string // "mysql" or "postgres"
	Host    string
	Port    int
	User    string
	Pass    string
	DB      string
	UseSSH  bool
	SSHHost string
	SSHPort int
	SSHUser string
	SSHPass string
}
