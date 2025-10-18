package ui

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/pn/kymar/internal/db"
)

// ShowMainInterface displays the main database query interface
func ShowMainInterface(w fyne.Window, dbh *sql.DB, closer func() error, connParams db.ConnParams, onDisconnect func()) {
	// Table list state
	var tableNames []string

	// Sort state tracking
	var currentTable string
	var sortColumn string
	var sortDirection string // "ASC" or "DESC"

	// Query editor
	query := widget.NewMultiLineEntry()
	query.SetPlaceHolder("-- Enter your SQL query here\n-- Example: SELECT * FROM users LIMIT 10;\n-- Tip: Use Cmd+Enter or click 'Run Query' to execute")

	// Table model state
	var headers []string
	var rows [][]string

	// Results table
	table := widget.NewTable(
		func() (int, int) {
			if len(headers) == 0 {
				return 1, 1
			}
			return len(rows) + 1, len(headers)
		},
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("")
			lbl.Wrapping = fyne.TextTruncate // Truncate text that's too long
			return lbl
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			lbl := o.(*widget.Label)
			if id.Row == 0 {
				// Header row styling
				if id.Col < len(headers) {
					headerText := headers[id.Col]
					// Add sort indicator if this column is being sorted
					if currentTable != "" && sortColumn == headers[id.Col] {
						if sortDirection == "ASC" {
							headerText += " ▲"
						} else {
							headerText += " ▼"
						}
					}
					lbl.SetText(headerText)
				} else {
					lbl.SetText("")
				}
				lbl.TextStyle = fyne.TextStyle{Bold: true}
				lbl.Alignment = fyne.TextAlignCenter
				lbl.Wrapping = fyne.TextTruncate
				return
			}
			// Data rows
			rowIdx := id.Row - 1
			if rowIdx < len(rows) && id.Col < len(rows[rowIdx]) {
				lbl.TextStyle = fyne.TextStyle{}
				lbl.Alignment = fyne.TextAlignLeading
				lbl.Wrapping = fyne.TextTruncate
				lbl.SetText(rows[rowIdx][id.Col])
			}
		},
	)

	// Set intelligent column widths based on content
	setupTableColumns := func() {
		if len(headers) > 0 {
			for i, header := range headers {
				// Calculate width based on header and content
				minWidth := float32(100) // Minimum 100px
				maxWidth := float32(300) // Maximum 300px for readability

				// Base width on header length
				headerWidth := float32(len(header) * 8) // ~8px per character

				// Check first few rows for content width
				contentWidth := headerWidth
				checkRows := len(rows)
				if checkRows > 10 {
					checkRows = 10 // Only check first 10 rows for performance
				}

				for j := 0; j < checkRows; j++ {
					if j < len(rows) && i < len(rows[j]) {
						cellWidth := float32(len(rows[j][i]) * 7) // ~7px per character for data
						if cellWidth > contentWidth {
							contentWidth = cellWidth
						}
					}
				}

				// Set width with min/max bounds
				width := contentWidth + 20 // Add padding
				if width < minWidth {
					width = minWidth
				}
				if width > maxWidth {
					width = maxWidth
				}

				table.SetColumnWidth(i, width)
			}
		}
	}

	// Left sidebar (like Sequel Pro)
	var tablesHeader *widget.Label
	if connParams.DBType == "mysql" && connParams.DB == "" {
		tablesHeader = widget.NewLabel("DATABASES")
	} else {
		tablesHeader = widget.NewLabel("TABLES")
	}
	tablesHeader.TextStyle = fyne.TextStyle{Bold: true}

	// Table list widget - declare early so it can be used in fetchTables
	var tableList *widget.List

	// Fetch tables function - defined early so it can be used in callbacks
	fetchTables := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var queryStr string
		if connParams.DBType == "mysql" {
			// For MySQL, if no database is selected, show all databases instead
			if connParams.DB == "" {
				queryStr = "SHOW DATABASES"
			} else {
				queryStr = "SHOW TABLES"
			}
		} else { // PostgreSQL
			queryStr = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'"
		}

		tableRows, err := dbh.QueryContext(ctx, queryStr)
		if err != nil {
			fmt.Printf("Error fetching tables: %v\n", err)
			tableNames = nil
			if tableList != nil {
				tableList.Refresh()
			}
			return
		}
		defer tableRows.Close()

		var tables []string
		for tableRows.Next() {
			var tableName string
			if err := tableRows.Scan(&tableName); err != nil {
				fmt.Printf("Error scanning table name: %v\n", err)
				continue
			}
			tables = append(tables, tableName)
		}

		if err := tableRows.Err(); err != nil {
			fmt.Printf("Error iterating table rows: %v\n", err)
		}

		fmt.Printf("Found %d tables: %v\n", len(tables), tables)
		tableNames = tables
		if tableList != nil {
			tableList.Refresh()
		}
	}

	// Status widget - declare early so it can be used in run function
	status := widget.NewLabel("🟢 Connected")
	status.TextStyle = fyne.TextStyle{Bold: true}

	// Run query function - define early so it can be used in table selection callback
	var run func()
	run = func() {
		if dbh == nil {
			dialog.ShowInformation("Not connected", "Database connection lost.", w)
			return
		}
		q := strings.TrimSpace(query.Text)
		if q == "" {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		start := time.Now()
		headers = nil
		rows = rows[:0]

		// Decide exec vs query
		lower := strings.ToLower(q)
		if strings.HasPrefix(lower, "select") || strings.HasPrefix(lower, "show") || strings.HasPrefix(lower, "desc") {
			r, err := dbh.QueryContext(ctx, q)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			defer r.Close()
			cols, _ := r.Columns()
			headers = cols
			vals := make([]sql.RawBytes, len(cols))
			scanArgs := make([]any, len(cols))
			for i := range vals {
				scanArgs[i] = &vals[i]
			}
			count := 0
			for r.Next() {
				if err := r.Scan(scanArgs...); err != nil {
					dialog.ShowError(err, w)
					return
				}
				out := make([]string, len(cols))
				for i, v := range vals {
					if v == nil {
						out[i] = "NULL"
					} else {
						out[i] = string(v)
					}
				}
				rows = append(rows, out)
				count++
				if count%200 == 0 {
					table.Refresh()
				}
			}
			if err := r.Err(); err != nil {
				dialog.ShowError(err, w)
				return
			}
			table.Refresh()
			setupTableColumns()
			status.SetText(fmt.Sprintf("🟢 Connected | %d row(s) in %v", len(rows), time.Since(start)))
			return
		}
		res, err := dbh.ExecContext(ctx, q)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		affected, _ := res.RowsAffected()
		headers = []string{"Result"}
		rows = [][]string{{fmt.Sprintf("OK, %d row(s) affected", affected)}}
		table.Refresh()
		setupTableColumns()
		status.SetText(fmt.Sprintf("🟢 Connected | Done in %v", time.Since(start)))
	}

	// Make column headers clickable for sorting (defined after run function)
	table.OnSelected = func(id widget.TableCellID) {
		if id.Row == 0 && currentTable != "" && id.Col < len(headers) {
			// Clicked on a header - toggle sort
			clickedColumn := headers[id.Col]

			if sortColumn == clickedColumn {
				// Toggle direction
				if sortDirection == "ASC" {
					sortDirection = "DESC"
				} else {
					sortDirection = "ASC"
				}
			} else {
				// New column - default to ASC
				sortColumn = clickedColumn
				sortDirection = "ASC"
			}

			// Regenerate and run the query with ORDER BY
			var sqlQuery string
			if connParams.DBType == "mysql" {
				sqlQuery = fmt.Sprintf("SELECT * FROM `%s` ORDER BY `%s` %s LIMIT 100;",
					currentTable, sortColumn, sortDirection)
			} else { // PostgreSQL
				sqlQuery = fmt.Sprintf("SELECT * FROM \"%s\" ORDER BY \"%s\" %s LIMIT 100;",
					currentTable, sortColumn, sortDirection)
			}

			query.SetText(sqlQuery)
			run()

			// Deselect the cell
			table.UnselectAll()
		}
	}

	// Now initialize the table list widget
	tableList = widget.NewList(
		func() int { return len(tableNames) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, o fyne.CanvasObject) {
			if id < len(tableNames) {
				o.(*widget.Label).SetText(tableNames[id])
			}
		},
	)
	tableList.OnSelected = func(id widget.ListItemID) {
		if id < len(tableNames) {
			itemName := tableNames[id]

			// Check if we're showing databases or tables
			if connParams.DBType == "mysql" && connParams.DB == "" {
				// We're showing databases, so switch to that database and show its tables
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				useQuery := fmt.Sprintf("USE `%s`", itemName)
				_, err := dbh.ExecContext(ctx, useQuery)
				if err != nil {
					dialog.ShowError(fmt.Errorf("failed to switch to database %s: %v", itemName, err), w)
					return
				}

				// Update connection params to reflect the selected database
				connParams.DB = itemName

				// Update the header to show "TABLES"
				tablesHeader.SetText("TABLES")

				// Refresh the table list to show tables from the selected database
				fetchTables()

				// Show a success message in the query editor
				query.SetText(fmt.Sprintf("-- Switched to database: %s\n-- Tables are now listed in the sidebar", itemName))
			} else {
				// We're showing tables, generate a SELECT statement
				currentTable = itemName

				// Try to find an ID column for default sorting
				var idColumn string
				commonIDColumns := []string{"id", "ID", "Id", itemName + "_id", "pk"}

				// We'll default to the first column if we can't find an ID
				// The actual ID detection will happen after we get the columns
				sortColumn = ""
				sortDirection = "ASC"

				// Generate query with ORDER BY for common ID columns
				var sqlQuery string
				if connParams.DBType == "mysql" {
					// Try to use id column if it exists
					for _, col := range commonIDColumns {
						idColumn = col
						break
					}
					if idColumn != "" {
						sqlQuery = fmt.Sprintf("SELECT * FROM `%s` ORDER BY `%s` ASC LIMIT 100;", itemName, idColumn)
						sortColumn = idColumn
					} else {
						sqlQuery = fmt.Sprintf("SELECT * FROM `%s` LIMIT 100;", itemName)
					}
					query.SetText(sqlQuery)
					run()
				} else { // PostgreSQL
					for _, col := range commonIDColumns {
						idColumn = col
						break
					}
					if idColumn != "" {
						sqlQuery = fmt.Sprintf("SELECT * FROM \"%s\" ORDER BY \"%s\" ASC LIMIT 100;", itemName, idColumn)
						sortColumn = idColumn
					} else {
						sqlQuery = fmt.Sprintf("SELECT * FROM \"%s\" LIMIT 100;", itemName)
					}
					query.SetText(sqlQuery)
					run()
				}
			}
		}
	}

	// Buttons
	runBtn := widget.NewButton("▶ Run Query", nil)
	runBtn.Importance = widget.HighImportance

	disconnectBtn := widget.NewButton("Disconnect", func() {
		if dbh != nil {
			_ = dbh.Close()
		}
		_ = closer()
		onDisconnect()
	})

	runBtn.OnTapped = run

	// Keyboard shortcuts
	s := &desktop.CustomShortcut{
		KeyName:  fyne.KeyReturn,
		Modifier: fyne.KeyModifierSuper,
	}
	w.Canvas().AddShortcut(s, func(sc fyne.Shortcut) {
		if w.Canvas().Focused() == query {
			run()
		}
	})

	// Initial fetch of tables/databases
	fetchTables()

	tableListContainer := container.NewVScroll(tableList)

	sidebar := container.NewBorder(
		container.NewVBox(tablesHeader, widget.NewSeparator()),
		container.NewVBox(widget.NewSeparator(), disconnectBtn),
		nil, nil,
		tableListContainer,
	)

	// Query editor area
	queryHeader := widget.NewLabel("SQL Query")
	queryHeader.TextStyle = fyne.TextStyle{Bold: true}

	queryToolbar := container.NewHBox(
		queryHeader,
		layout.NewSpacer(),
		runBtn,
	)

	queryArea := container.NewBorder(queryToolbar, nil, nil, nil, query)

	// Results area
	resultsHeader := widget.NewLabel("Query Results")
	resultsHeader.TextStyle = fyne.TextStyle{Bold: true}

	resultsArea := container.NewBorder(
		resultsHeader,
		status,
		nil, nil,
		container.NewVScroll(table),
	)

	// Main content area
	mainContent := container.NewVSplit(queryArea, resultsArea)
	mainContent.SetOffset(0.3)

	// Overall layout
	root := container.NewHSplit(sidebar, mainContent)
	root.SetOffset(0.2) // Narrow sidebar like Sequel Pro

	w.SetContent(root)
}
