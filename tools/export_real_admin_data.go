package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type resultSet struct {
	Orders   []map[string]any `json:"orders"`
	Admins   []map[string]any `json:"admins"`
	Roles    []map[string]any `json:"roles"`
	Menus    []map[string]any `json:"menus"`
	Cities   []map[string]any `json:"cities"`
	CarTypes []map[string]any `json:"carTypes"`
}

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	outPath := os.Getenv("OUT_PATH")
	if dsn == "" || outPath == "" {
		panic("MYSQL_DSN and OUT_PATH are required")
	}

	db, err := sql.Open("mysql", dsn)
	must(err)
	defer db.Close()
	must(db.Ping())

	data := resultSet{
		Orders:   query(db, "passenger_orders", []string{"order_id", "passenger_name", "passenger_phone", "driver_id", "car_type", "start_address", "end_address", "status", "estimated_price", "final_price", "created_at"}, "created_at DESC", 200),
		Admins:   queryAdmins(db),
		Roles:    queryRoles(db),
		Menus:    query(db, "admin_menus", []string{"id", "parent_id", "name", "path", "icon", "sort"}, "sort ASC, id ASC", 500),
		Cities:   query(db, "city_configs", []string{"id", "city_code", "city_name", "status"}, "id ASC", 500),
		CarTypes: query(db, "car_type_configs", []string{"id", "type_name", "base_price", "km_price", "minute_price", "status"}, "id ASC", 500),
	}

	file, err := os.Create(outPath)
	must(err)
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	must(enc.Encode(data))
	fmt.Println(outPath)
}

func queryAdmins(db *sql.DB) []map[string]any {
	if !tableExists(db, "admin_users") {
		return []map[string]any{}
	}

	sqlText := `SELECT u.id, u.username, u.role_id, COALESCE(r.name, '') AS role_name, u.status, u.last_login_time
FROM admin_users u
LEFT JOIN admin_roles r ON r.id = u.role_id
ORDER BY u.id ASC
LIMIT 500`
	rows, err := db.Query(sqlText)
	if err != nil {
		return []map[string]any{}
	}
	defer rows.Close()
	return rowsToMaps(rows)
}

func queryRoles(db *sql.DB) []map[string]any {
	roles := query(db, "admin_roles", []string{"id", "name", "remark"}, "id ASC", 500)
	if len(roles) == 0 || !tableExists(db, "admin_role_menus") {
		for _, role := range roles {
			role["menu_ids"] = []any{}
			role["menu_count"] = 0
		}
		return roles
	}

	roleMenus := query(db, "admin_role_menus", []string{"role_id", "menu_id"}, "role_id ASC, menu_id ASC", 10000)
	menuMap := map[string][]any{}
	for _, item := range roleMenus {
		roleID := fmt.Sprint(item["role_id"])
		menuMap[roleID] = append(menuMap[roleID], item["menu_id"])
	}
	for _, role := range roles {
		menuIDs := menuMap[fmt.Sprint(role["id"])]
		if menuIDs == nil {
			menuIDs = []any{}
		}
		role["menu_ids"] = menuIDs
		role["menu_count"] = len(menuIDs)
	}
	return roles
}

func query(db *sql.DB, table string, wanted []string, orderBy string, limit int) []map[string]any {
	if !tableExists(db, table) {
		return []map[string]any{}
	}

	cols := existingColumns(db, table, wanted)
	if len(cols) == 0 {
		return []map[string]any{}
	}

	sqlText := fmt.Sprintf("SELECT %s FROM %s", strings.Join(backtick(cols), ", "), table)
	if orderBy != "" {
		sqlText += " ORDER BY " + orderBy
	}
	if limit > 0 {
		sqlText += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := db.Query(sqlText)
	if err != nil {
		return []map[string]any{}
	}
	defer rows.Close()

	out := rowsToMaps(rows)
	for _, row := range out {
		for _, key := range wanted {
			if _, ok := row[key]; !ok {
				row[key] = nil
			}
		}
	}
	return out
}

func rowsToMaps(rows *sql.Rows) []map[string]any {
	cols, err := rows.Columns()
	must(err)

	out := []map[string]any{}
	for rows.Next() {
		values := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}
		must(rows.Scan(ptrs...))

		row := map[string]any{}
		for i, col := range cols {
			switch v := values[i].(type) {
			case []byte:
				row[col] = string(v)
			default:
				row[col] = v
			}
		}
		out = append(out, row)
	}
	must(rows.Err())
	return out
}

func tableExists(db *sql.DB, table string) bool {
	var name string
	err := db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", table).Scan(&name)
	return err == nil
}

func existingColumns(db *sql.DB, table string, wanted []string) []string {
	rows, err := db.Query("SELECT column_name FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = ?", table)
	if err != nil {
		return nil
	}
	defer rows.Close()

	found := map[string]bool{}
	for rows.Next() {
		var col string
		must(rows.Scan(&col))
		found[col] = true
	}

	var cols []string
	for _, col := range wanted {
		if found[col] {
			cols = append(cols, col)
		}
	}
	return cols
}

func backtick(cols []string) []string {
	out := make([]string, len(cols))
	for i, col := range cols {
		out[i] = "`" + strings.ReplaceAll(col, "`", "``") + "`"
	}
	return out
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
