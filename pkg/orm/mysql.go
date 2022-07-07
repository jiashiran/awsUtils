package orm

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLDB(
	database string,
	username string,
	password string,
	host string,
	port int,
) DB {
	config := mysql.Config{
		User:                 username,
		Passwd:               password,
		DBName:               database,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", host, port),
		MultiStatements:      true,
		AllowNativePasswords: true,
		Params: map[string]string{
			"charset":   "utf8",
			"parseTime": "True",
		},
	}

	return &mySQLDB{
		dsn:    config.FormatDSN(),
		dbName: database,
	}
}

type mySQLDB struct {
	dsn    string
	db     *sql.DB
	dbName string
}

func (m *mySQLDB) Open() error {
	db, err := sql.Open("mysql", m.dsn)
	if err != nil {
		return err
	}

	m.db = db

	return nil
}

func (m *mySQLDB) Close() error {
	return m.db.Close()
}

func (m *mySQLDB) GetSchemaRows() (*sql.Rows, error) {
	query := `
	SELECT table_name,
				 column_name,
				 data_type,
				 character_maximum_length
	FROM   information_schema.columns
	WHERE  table_schema = ?`
	rows, err := m.db.Query(query, m.dbName)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (m *mySQLDB) DB() *sql.DB {
	return m.db
}

func (m *mySQLDB) ColumnNameForSelect(name string) string {
	return fmt.Sprintf("`%s`", name)
}

func (m *mySQLDB) EnableConstraints() error {
	_, err := m.db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	return err
}

func (m *mySQLDB) DisableConstraints() error {
	_, err := m.db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	return err
}

//---
type DB interface {
	Open() error
	Close() error
	GetSchemaRows() (*sql.Rows, error)
	DisableConstraints() error
	EnableConstraints() error
	ColumnNameForSelect(columnName string) string
	DB() *sql.DB
}

type Schema struct {
	Tables map[string]*Table
}

func (s *Schema) GetTable(name string) (*Table, error) {
	if table, ok := s.Tables[name]; ok {
		return table, nil
	}

	return nil, fmt.Errorf("table '%s' not found", name)
}

type Table struct {
	Name    string
	Columns []*Column
}

func (t *Table) HasColumn(name string) bool {
	_, _, err := t.GetColumn(name)
	return err == nil
}

func (t *Table) GetColumn(name string) (int, *Column, error) {
	for i, column := range t.Columns {
		if column.Name == name {
			return i, column, nil
		}
	}

	return -1, nil, fmt.Errorf("column '%s' not found", name)
}

type Column struct {
	Name     string
	Type     string
	MaxChars int64
}

func (c *Column) Compatible(other *Column) bool {
	if c.MaxChars == 0 && other.MaxChars == 0 {
		return true
	}

	if c.MaxChars > 0 && other.MaxChars > 0 {
		return c.MaxChars < other.MaxChars
	}

	return false
}

func (c *Column) Incompatible(other *Column) bool {
	return !c.Compatible(other)
}

func BuildSchema(db DB) (*Schema, error) {
	rows, err := db.GetSchemaRows()
	if err != nil {
		return nil, err
	}

	data := map[string][]*Column{}
	for rows.Next() {
		var (
			table    sql.NullString
			column   sql.NullString
			datatype sql.NullString
			maxChars sql.NullInt64
		)

		if err := rows.Scan(&table, &column, &datatype, &maxChars); err != nil {
			return nil, err
		}

		data[table.String] = append(data[table.String], &Column{
			Name:     column.String,
			Type:     datatype.String,
			MaxChars: maxChars.Int64,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate through schema rows: %s", err)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed closing rows: %s", err)
	}

	schema := &Schema{
		Tables: map[string]*Table{},
	}

	for k, v := range data {
		schema.Tables[k] = &Table{
			Name:    k,
			Columns: v,
		}
	}

	return schema, nil
}

func GetIncompatibleColumns(src, dst *Table) ([]*Column, error) {
	var incompatibleColumns []*Column
	for _, dstColumn := range dst.Columns {
		_, srcColumn, err := src.GetColumn(dstColumn.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to find column '%s/%s' in source schema: %s", dst.Name, dstColumn.Name, err)
		}

		if dstColumn.Incompatible(srcColumn) {
			incompatibleColumns = append(incompatibleColumns, dstColumn)
		}
	}

	return incompatibleColumns, nil
}

func GetIncompatibleRowIDs(db DB, src, dst *Table) ([]int, error) {
	columns, err := GetIncompatibleColumns(src, dst)
	if err != nil {
		return nil, fmt.Errorf("failed getting incompatible columns: %s", err)
	}

	if columns == nil {
		return nil, nil
	}

	limits := make([]string, len(columns))
	for i, column := range columns {
		limits[i] = fmt.Sprintf("LENGTH(%s) > %d", column.Name, column.MaxChars)
	}

	stmt := fmt.Sprintf("SELECT id FROM %s WHERE %s", src.Name, strings.Join(limits, " OR "))
	rows, err := db.DB().Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed getting incompatible row ids: %s", err)
	}

	var rowIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan row: %s", err)
		}
		rowIDs = append(rowIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return rowIDs, nil
}

func GetIncompatibleRowCount(db DB, src, dst *Table) (int64, error) {
	columns, err := GetIncompatibleColumns(src, dst)
	if err != nil {
		return 0, fmt.Errorf("failed getting incompatible columns: %s", err)
	}

	if columns == nil {
		return 0, nil
	}

	limits := make([]string, len(columns))
	for i, column := range columns {
		limits[i] = fmt.Sprintf("length(%s) > %d", column.Name, column.MaxChars)
	}

	stmt := fmt.Sprintf("SELECT count(1) FROM %s WHERE %s", src.Name, strings.Join(limits, " OR "))

	var count int64
	err = db.DB().QueryRow(stmt).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
