package db

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	mysql "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/pn/kymar/internal/ssh"
)

var (
	registerSSHOnce  sync.Once
	currentSSHDialer func(ctx context.Context, addr string) (net.Conn, error)
	defaultTCPDialer = func(ctx context.Context, addr string) (net.Conn, error) {
		d := &net.Dialer{Timeout: 5 * time.Second}
		return d.DialContext(ctx, "tcp", addr)
	}
)

// Connect establishes a database connection with the given parameters
func Connect(p ConnParams) (*sql.DB, func() error, error) {
	// Build the dialer used when DSN protocol is "ssh"
	var sshClose func() error = func() error { return nil }
	var forwardListener net.Listener
	var forwardDone chan struct{}

	if p.UseSSH {
		d, c, err := ssh.NewTunnelDialer(p.SSHHost, p.SSHPort, p.SSHUser, p.SSHPass)
		if err != nil {
			return nil, nil, err
		}
		sshClose = c

		// Wrap SSH tunnel d(network, addr) into a context-aware dialer
		currentSSHDialer = func(ctx context.Context, addr string) (net.Conn, error) {
			return d("tcp", addr)
		}

		// Register the "ssh" protocol once
		registerSSHOnce.Do(func() {
			mysql.RegisterDialContext("ssh", func(ctx context.Context, addr string) (net.Conn, error) {
				if currentSSHDialer != nil {
					return currentSSHDialer(ctx, addr)
				}
				return defaultTCPDialer(ctx, addr)
			})
		})

		// For PostgreSQL, lib/pq doesn't support custom dialers directly. Create a local TCP
		// forwarder over the SSH connection and connect to that.
		if p.DBType == "postgres" {
			ln, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				return nil, nil, err
			}
			forwardListener = ln
			forwardDone = make(chan struct{})
			remoteAddr := fmt.Sprintf("%s:%d", p.Host, p.Port)
			go func() {
				for {
					conn, err := ln.Accept()
					if err != nil {
						select {
						case <-forwardDone:
							return
						default:
							continue
						}
					}
					go func(c net.Conn) {
						defer c.Close()
						rc, err := d("tcp", remoteAddr)
						if err != nil {
							return
						}
						defer rc.Close()
						// Bi-directional copy
						go io.Copy(rc, c)
						io.Copy(c, rc)
					}(conn)
				}
			}()

			prevClose := sshClose
			sshClose = func() error {
				if forwardListener != nil {
					close(forwardDone)
					_ = forwardListener.Close()
				}
				return prevClose()
			}
		}
	} else {
		currentSSHDialer = nil
	}

	var dsn string
	var driverName string
	// Effective host/port (may be overridden by SSH forwarder for PostgreSQL)
	effectiveHost := p.Host
	effectivePort := p.Port
	if p.UseSSH && p.DBType == "postgres" && forwardListener != nil {
		if tcpAddr, ok := forwardListener.Addr().(*net.TCPAddr); ok {
			effectiveHost = "127.0.0.1"
			effectivePort = tcpAddr.Port
		}
	}

	if p.DBType == "mysql" {
		proto := "tcp"
		if p.UseSSH {
			proto = "ssh"
		}

		dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true&multiStatements=true",
			p.User, p.Pass, proto, p.Host, p.Port, p.DB)
		if p.DB == "" {
			dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/?parseTime=true&multiStatements=true",
				p.User, p.Pass, proto, p.Host, p.Port)
		}
		driverName = "mysql"
	} else { // PostgreSQL
		if p.DB == "" {
			dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
				effectiveHost, effectivePort, p.User, p.Pass)
		} else {
			dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				effectiveHost, effectivePort, p.User, p.Pass, p.DB)
		}
		driverName = "postgres"
	}

	dbh, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, nil, err
	}
	dbh.SetConnMaxLifetime(5 * time.Minute)
	dbh.SetMaxOpenConns(5)
	dbh.SetMaxIdleConns(2)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dbh.PingContext(ctx); err != nil {
		_ = dbh.Close()
		return nil, nil, err
	}

	return dbh, sshClose, nil
}
