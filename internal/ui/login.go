package ui

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/pn/kymar/internal/config"
	"github.com/pn/kymar/internal/db"
)

// ShowLoginScreen displays the login/connection screen
func ShowLoginScreen(w fyne.Window, onConnect func(db.ConnParams)) {
	// Load saved connections
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{Connections: []config.SavedConnection{}}
	}

	// Left sidebar with favorites
	favoritesHeader := widget.NewLabel("SAVED CONNECTIONS")
	favoritesHeader.TextStyle = fyne.TextStyle{Bold: true}

	// Saved connections list
	var favoritesList *widget.List
	favoritesList = widget.NewList(
		func() int { return len(cfg.Connections) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, o fyne.CanvasObject) {
			if id < len(cfg.Connections) {
				o.(*widget.Label).SetText(cfg.Connections[id].Name)
			}
		},
	)

	// Handle clicking on a saved connection
	favoritesList.OnSelected = func(id widget.ListItemID) {
		if id < len(cfg.Connections) {
			conn := cfg.Connections[id]
			onConnect(conn.Params)
		}
	}

	// Quick Connect section
	quickConnectHeader := widget.NewLabel("QUICK CONNECT")
	quickConnectHeader.TextStyle = fyne.TextStyle{Bold: true}

	// Callback to refresh the connections list
	refreshConnections := func() {
		cfg, _ = config.Load()
		favoritesList.Refresh()
	}

	// Left sidebar content
	sidebar := container.NewVBox(
		favoritesHeader,
		container.NewScroll(favoritesList),
		widget.NewSeparator(),
		quickConnectHeader,
	)

	// Connection tabs (TCP/IP, Socket, SSH)
	connectionTabs := container.NewAppTabs(
		container.NewTabItem("TCP/IP", createTCPIPTab(w, onConnect, cfg, refreshConnections)),
		container.NewTabItem("Socket", createSocketTab(w)),
		container.NewTabItem("SSH", createSSHTab(w, onConnect, cfg, refreshConnections)),
	)

	// Main connection area with centered form
	connectionArea := container.NewVBox(
		widget.NewLabel("Enter connection details below, or choose a saved connection"),
		widget.NewSeparator(),
		connectionTabs,
	)

	// Create the main layout (sidebar + connection area)
	mainLayout := container.NewHSplit(
		container.NewBorder(nil, nil, nil, nil, sidebar),
		container.NewBorder(nil, nil, nil, nil, connectionArea),
	)
	mainLayout.SetOffset(0.25) // 25% for sidebar, 75% for connection area

	w.SetContent(mainLayout)
}

func createTCPIPTab(w fyne.Window, onConnect func(db.ConnParams), cfg *config.Config, refreshConnections func()) *fyne.Container {
	// Database type selector
	dbType := widget.NewSelect([]string{"MySQL", "PostgreSQL"}, nil)
	dbType.SetSelected("MySQL")

	// Connection form fields
	name := widget.NewEntry()
	name.SetPlaceHolder("My Connection")

	host := widget.NewEntry()
	host.SetText("127.0.0.1")

	username := widget.NewEntry()
	username.SetText("root")

	password := widget.NewPasswordEntry()

	database := widget.NewEntry()
	database.SetPlaceHolder("database_name")

	port := widget.NewEntry()
	port.SetText("3306")

	saveConnection := widget.NewCheck("Save this connection", nil)

	// Set up database type change callback
	dbType.OnChanged = func(value string) {
		if value == "MySQL" {
			if port.Text == "5432" || port.Text == "" {
				port.SetText("3306")
			}
			if username.Text == "postgres" {
				username.SetText("root")
			}
		} else if value == "PostgreSQL" {
			if port.Text == "3306" || port.Text == "" {
				port.SetText("5432")
			}
			if username.Text == "root" {
				username.SetText("postgres")
			}
		}
	}

	// Connect button
	connectBtn := widget.NewButton("Connect", func() {
		selectedDBType := "mysql"
		if dbType.Selected == "PostgreSQL" {
			selectedDBType = "postgres"
		}

		p := db.ConnParams{
			DBType: selectedDBType,
			Host:   strings.TrimSpace(host.Text),
			User:   strings.TrimSpace(username.Text),
			Pass:   password.Text,
			DB:     strings.TrimSpace(database.Text),
			UseSSH: false,
		}
		p.Port, _ = strconv.Atoi(strings.TrimSpace(port.Text))

		// Save connection if checkbox is checked
		if saveConnection.Checked {
			connName := strings.TrimSpace(name.Text)
			if connName == "" {
				connName = p.Host + ":" + strconv.Itoa(p.Port)
			}

			savedConn := config.SavedConnection{
				Name:   connName,
				Params: p,
			}

			if err := cfg.AddConnection(savedConn); err != nil {
				dialog.ShowError(err, w)
			} else {
				refreshConnections()
			}
		}

		onConnect(p)
	})
	connectBtn.Importance = widget.HighImportance

	// Form layout
	form := container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Type:"), dbType,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Name:"), name,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Host:"), host,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Username:"), username,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Password:"), password,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Database:"), database,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Port:"), port,
		),
		widget.NewSeparator(),
		saveConnection,
		connectBtn,
	)

	// Center the form
	return container.NewBorder(
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		container.NewPadded(form),
	)
}

func createSocketTab(w fyne.Window) *fyne.Container {
	// Socket connection form (simplified for now)
	socketPath := widget.NewEntry()
	socketPath.SetText("/tmp/mysql.sock")

	username := widget.NewEntry()
	username.SetText("root")

	password := widget.NewPasswordEntry()

	database := widget.NewEntry()

	connectBtn := widget.NewButton("Connect", func() {
		dialog.ShowInformation("Socket Connection", "Socket connections not implemented yet", w)
	})
	connectBtn.Importance = widget.HighImportance

	form := container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Socket:"), socketPath,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Username:"), username,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Password:"), password,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Database:"), database,
		),
		widget.NewSeparator(),
		connectBtn,
	)

	return container.NewBorder(
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		container.NewPadded(form),
	)
}

func createSSHTab(w fyne.Window, onConnect func(db.ConnParams), cfg *config.Config, refreshConnections func()) *fyne.Container {
	// Database type selector
	dbType := widget.NewSelect([]string{"MySQL", "PostgreSQL"}, nil)
	dbType.SetSelected("MySQL")

	// Connection name
	name := widget.NewEntry()
	name.SetPlaceHolder("My SSH Connection")

	// SSH connection form
	host := widget.NewEntry()
	host.SetText("127.0.0.1")

	username := widget.NewEntry()
	username.SetText("root")

	password := widget.NewPasswordEntry()

	database := widget.NewEntry()

	port := widget.NewEntry()
	port.SetText("3306")

	// Set up database type change callback
	dbType.OnChanged = func(value string) {
		if value == "MySQL" {
			if port.Text == "5432" || port.Text == "" {
				port.SetText("3306")
			}
			if username.Text == "postgres" {
				username.SetText("root")
			}
		} else if value == "PostgreSQL" {
			if port.Text == "3306" || port.Text == "" {
				port.SetText("5432")
			}
			if username.Text == "root" {
				username.SetText("postgres")
			}
		}
	}

	// SSH fields
	sshHost := widget.NewEntry()
	sshHost.SetPlaceHolder("ssh.example.com")

	sshUser := widget.NewEntry()
	sshUser.SetPlaceHolder("ec2-user")

	sshPassword := widget.NewPasswordEntry()

	sshPort := widget.NewEntry()
	sshPort.SetText("22")

	saveConnection := widget.NewCheck("Save this connection", nil)

	connectBtn := widget.NewButton("Connect", func() {
		selectedDBType := "mysql"
		if dbType.Selected == "PostgreSQL" {
			selectedDBType = "postgres"
		}

		p := db.ConnParams{
			DBType:  selectedDBType,
			Host:    strings.TrimSpace(host.Text),
			User:    strings.TrimSpace(username.Text),
			Pass:    password.Text,
			DB:      strings.TrimSpace(database.Text),
			UseSSH:  true,
			SSHHost: strings.TrimSpace(sshHost.Text),
			SSHUser: strings.TrimSpace(sshUser.Text),
			SSHPass: sshPassword.Text,
		}
		p.Port, _ = strconv.Atoi(strings.TrimSpace(port.Text))
		p.SSHPort, _ = strconv.Atoi(strings.TrimSpace(sshPort.Text))

		// Save connection if checkbox is checked
		if saveConnection.Checked {
			connName := strings.TrimSpace(name.Text)
			if connName == "" {
				connName = p.Host + " via SSH"
			}

			savedConn := config.SavedConnection{
				Name:   connName,
				Params: p,
			}

			if err := cfg.AddConnection(savedConn); err != nil {
				dialog.ShowError(err, w)
			} else {
				refreshConnections()
			}
		}

		onConnect(p)
	})
	connectBtn.Importance = widget.HighImportance

	form := container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Type:"), dbType,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Name:"), name,
		),
		widget.NewLabel("Database Connection"),
		container.NewGridWithColumns(2,
			widget.NewLabel("Host:"), host,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Username:"), username,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Password:"), password,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Database:"), database,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("Port:"), port,
		),
		widget.NewSeparator(),
		widget.NewLabel("SSH Tunnel"),
		container.NewGridWithColumns(2,
			widget.NewLabel("SSH Host:"), sshHost,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("SSH User:"), sshUser,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("SSH Password:"), sshPassword,
		),
		container.NewGridWithColumns(2,
			widget.NewLabel("SSH Port:"), sshPort,
		),
		widget.NewSeparator(),
		saveConnection,
		connectBtn,
	)

	return container.NewBorder(
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		container.NewPadded(form),
	)
}
