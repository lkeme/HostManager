package wxSQLite3

import (
	sqlite3 "HostManager/backend/pkg/wxSQLite3/go-wxsqlite3"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// sqlite3Header defines the header string used by SQLite 3.
var sqlite3Header = []byte("SQLite format 3\000")

// IsEncrypted returns true, if the database with the given filename is
// encrypted, and false otherwise.
// If the database header cannot be read properly an error is returned.
func IsEncrypted(filename string) (bool, error) {
	// open file
	db, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer db.Close()
	// read header
	var header [16]byte
	n, err := db.Read(header[:])
	if err != nil {
		return false, err
	}
	if n != len(header) {
		return false, errors.New("go-wxsqlite3: could not read full header")
	}
	// SQLCipher encrypts also the header, the file is encrypted if the read
	// header does not equal the header string used by SQLite 3.
	encrypted := !bytes.Equal(header[:], sqlite3Header)
	return encrypted, nil
}

// ResetDBKey Reset database password
func ResetDBKey(dbName string, oldKey string, newKey string) error {
	sr := sqlite3.Cpp{}
	err := sr.Open(dbName, oldKey)
	if err != nil {
		return errors.New(fmt.Sprintf("open err: %v", err))
	}

	err = sr.ReKey(newKey)
	if err != nil {
		return errors.New(fmt.Sprintf("rekey err: %v", err))
	}
	sr.Close()
	return nil
}

// wxSqlite3 默认aes128加密

var (
	ErrConstraintsNotImplemented = errors.New("constraints not implemented on sqlite, consider using DisableForeignKeyConstraintWhenMigrating, more details https://github.com/go-gorm/gorm/wiki/GORM-V2-Release-Note-Draft#all-new-migrator")
)

// DriverName is the default driver name for SQLite.
const DriverName = "sqlite3"

type Dialector struct {
	DriverName string
	DSN        string
	Conn       gorm.ConnPool
}

func Open(dsn string) gorm.Dialector {
	return &Dialector{DSN: dsn}
}

func (dialector Dialector) Name() string {
	return "sqlite"
}

func (dialector Dialector) Initialize(db *gorm.DB) (err error) {
	if dialector.DriverName == "" {
		dialector.DriverName = DriverName
	}

	// register callbacks
	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
		LastInsertIDReversed: true,
	})

	if dialector.Conn != nil {
		db.ConnPool = dialector.Conn
	} else {
		db.ConnPool, err = sql.Open(dialector.DriverName, dialector.DSN)
		if err != nil {
			return err
		}
	}

	for k, v := range dialector.ClauseBuilders() {
		db.ClauseBuilders[k] = v
	}
	return
}

func (dialector Dialector) ClauseBuilders() map[string]clause.ClauseBuilder {
	return map[string]clause.ClauseBuilder{
		"INSERT": func(c clause.Clause, builder clause.Builder) {
			if insert, ok := c.Expression.(clause.Insert); ok {
				if stmt, ok := builder.(*gorm.Statement); ok {
					stmt.WriteString("INSERT ")
					if insert.Modifier != "" {
						stmt.WriteString(insert.Modifier)
						stmt.WriteByte(' ')
					}

					stmt.WriteString("INTO ")
					if insert.Table.Name == "" {
						stmt.WriteQuoted(stmt.Table)
					} else {
						stmt.WriteQuoted(insert.Table)
					}
					return
				}
			}

			c.Build(builder)
		},
		"LIMIT": func(c clause.Clause, builder clause.Builder) {
			if limit, ok := c.Expression.(clause.Limit); ok {
				var lmt = -1
				if limit.Limit != nil && *limit.Limit >= 0 {
					lmt = *limit.Limit
				}
				if lmt >= 0 || limit.Offset > 0 {
					builder.WriteString("LIMIT ")
					builder.WriteString(strconv.Itoa(lmt))
				}
				if limit.Offset > 0 {
					builder.WriteString(" OFFSET ")
					builder.WriteString(strconv.Itoa(limit.Offset))
				}
			}
		},
		"FOR": func(c clause.Clause, builder clause.Builder) {
			if _, ok := c.Expression.(clause.Locking); ok {
				// SQLite3 does not support row-level locking.
				return
			}
			c.Build(builder)
		},
	}
}

func (dialector Dialector) DefaultValueOf(field *schema.Field) clause.Expression {
	if field.AutoIncrement {
		return clause.Expr{SQL: "NULL"}
	}

	// doesn't work, will raise error
	return clause.Expr{SQL: "DEFAULT"}
}

func (dialector Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	return Migrator{migrator.Migrator{Config: migrator.Config{
		DB:                          db,
		Dialector:                   dialector,
		CreateIndexAfterCreateTable: true,
	}}}
}

func (dialector Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	writer.WriteByte('?')
}

func (dialector Dialector) QuoteTo(writer clause.Writer, str string) {
	writer.WriteByte('`')
	if strings.Contains(str, ".") {
		for idx, str := range strings.Split(str, ".") {
			if idx > 0 {
				writer.WriteString(".`")
			}
			writer.WriteString(str)
			writer.WriteByte('`')
		}
	} else {
		writer.WriteString(str)
		writer.WriteByte('`')
	}
}

func (dialector Dialector) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, nil, `"`, vars...)
}

func (dialector Dialector) DataTypeOf(field *schema.Field) string {
	switch field.DataType {
	case schema.Bool:
		return "numeric"
	case schema.Int, schema.Uint:
		if field.AutoIncrement && !field.PrimaryKey {
			// https://www.sqlite.org/autoinc.html
			return "integer PRIMARY KEY AUTOINCREMENT"
		} else {
			return "integer"
		}
	case schema.Float:
		return "real"
	case schema.String:
		return "text"
	case schema.Time:
		return "datetime"
	case schema.Bytes:
		return "blob"
	}

	return string(field.DataType)
}

func (dialectopr Dialector) SavePoint(tx *gorm.DB, name string) error {
	tx.Exec("SAVEPOINT " + name)
	return nil
}

func (dialectopr Dialector) RollbackTo(tx *gorm.DB, name string) error {
	tx.Exec("ROLLBACK TO SAVEPOINT " + name)
	return nil
}

// ========================================================================================
// ======================    Migrator   ===================================================
// ========================================================================================

type Migrator struct {
	migrator.Migrator
}

func (m *Migrator) RunWithoutForeignKey(fc func() error) error {
	var enabled int
	m.DB.Raw("PRAGMA foreign_keys").Scan(&enabled)
	if enabled == 1 {
		m.DB.Exec("PRAGMA foreign_keys = OFF")
		defer m.DB.Exec("PRAGMA foreign_keys = ON")
	}

	return fc()
}

func (m Migrator) HasTable(value interface{}) bool {
	var count int
	m.Migrator.RunWithValue(value, func(stmt *gorm.Statement) error {
		return m.DB.Raw("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", stmt.Table).Row().Scan(&count)
	})
	return count > 0
}

func (m Migrator) DropTable(values ...interface{}) error {
	return m.RunWithoutForeignKey(func() error {
		values = m.ReorderModels(values, false)
		tx := m.DB.Session(&gorm.Session{})

		for i := len(values) - 1; i >= 0; i-- {
			if err := m.RunWithValue(values[i], func(stmt *gorm.Statement) error {
				return tx.Exec("DROP TABLE IF EXISTS ?", clause.Table{Name: stmt.Table}).Error
			}); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func (m Migrator) HasColumn(value interface{}, name string) bool {
	var count int
	m.Migrator.RunWithValue(value, func(stmt *gorm.Statement) error {
		if field := stmt.Schema.LookUpField(name); field != nil {
			name = field.DBName
		}

		if name != "" {
			m.DB.Raw(
				"SELECT count(*) FROM sqlite_master WHERE type = ? AND tbl_name = ? AND (sql LIKE ? OR sql LIKE ? OR sql LIKE ?)",
				"table", stmt.Table, `%"`+name+`" %`, `%`+name+` %`, "%`"+name+"`%",
			).Row().Scan(&count)
		}
		return nil
	})
	return count > 0
}

// ColumnTypes return columnTypes []gorm.ColumnType and execErr error
func (m Migrator) ColumnTypes(value interface{}) ([]gorm.ColumnType, error) {
	columnTypes := make([]gorm.ColumnType, 0)
	execErr := m.RunWithValue(value, func(stmt *gorm.Statement) (err error) {
		var (
			sqls   []string
			sqlDDL *ddl
		)

		if err := m.DB.Raw("SELECT sql FROM sqlite_master WHERE type IN ? AND tbl_name = ? AND sql IS NOT NULL order by type = ? desc", []string{"table", "index"}, stmt.Table, "table").Scan(&sqls).Error; err != nil {
			return err
		}

		if sqlDDL, err = parseDDL(sqls...); err != nil {
			return err
		}

		rows, err := m.DB.Session(&gorm.Session{}).Table(stmt.Table).Limit(1).Rows()
		if err != nil {
			return err
		}
		defer func() {
			err = rows.Close()
		}()

		var rawColumnTypes []*sql.ColumnType
		rawColumnTypes, err = rows.ColumnTypes()
		if err != nil {
			return err
		}

		for _, c := range rawColumnTypes {
			columnType := migrator.ColumnType{SQLColumnType: c}
			for _, column := range sqlDDL.columns {
				if column.NameValue.String == c.Name() {
					column.SQLColumnType = c
					columnType = column
					break
				}
			}
			columnTypes = append(columnTypes, columnType)
		}

		return err
	})

	return columnTypes, execErr
}

func (m Migrator) AlterColumn(value interface{}, name string) error {
	return m.RunWithoutForeignKey(func() error {
		return m.RunWithValue(value, func(stmt *gorm.Statement) error {
			if field := stmt.Schema.LookUpField(name); field != nil {
				var (
					createSQL    string
					newTableName = stmt.Table + "__temp"
				)

				m.DB.Raw("SELECT sql FROM sqlite_master WHERE type = ? AND tbl_name = ? AND name = ?", "table", stmt.Table, stmt.Table).Row().Scan(&createSQL)

				if reg, err := regexp.Compile("(`|'|\"| )" + field.DBName + "(`|'|\"| ) .*?,"); err == nil {
					tableReg, err := regexp.Compile(" ('|`|\"| )" + stmt.Table + "('|`|\"| ) ")
					if err != nil {
						return err
					}

					createSQL = tableReg.ReplaceAllString(createSQL, fmt.Sprintf(" `%v` ", newTableName))
					createSQL = reg.ReplaceAllString(createSQL, fmt.Sprintf("`%v` ?,", field.DBName))

					var columns []string
					columnTypes, _ := m.DB.Migrator().ColumnTypes(value)
					for _, columnType := range columnTypes {
						columns = append(columns, fmt.Sprintf("`%v`", columnType.Name()))
					}

					return m.DB.Transaction(func(tx *gorm.DB) error {
						queries := []string{
							createSQL,
							fmt.Sprintf("INSERT INTO `%v`(%v) SELECT %v FROM `%v`", newTableName, strings.Join(columns, ","), strings.Join(columns, ","), stmt.Table),
							fmt.Sprintf("DROP TABLE `%v`", stmt.Table),
							fmt.Sprintf("ALTER TABLE `%v` RENAME TO `%v`", newTableName, stmt.Table),
						}
						for _, query := range queries {
							if err := tx.Exec(query, m.FullDataTypeOf(field)).Error; err != nil {
								return err
							}
						}
						return nil
					})
				} else {
					return err
				}
			} else {
				return fmt.Errorf("failed to alter field with name %v", name)
			}
		})
	})
}

func (m Migrator) DropColumn(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if field := stmt.Schema.LookUpField(name); field != nil {
			name = field.DBName
		}

		var (
			createSQL    string
			newTableName = stmt.Table + "__temp"
		)

		m.DB.Raw("SELECT sql FROM sqlite_master WHERE type = ? AND tbl_name = ? AND name = ?", "table", stmt.Table, stmt.Table).Row().Scan(&createSQL)

		if reg, err := regexp.Compile("(`|'|\"| )" + name + "(`|'|\"| ) .*?,"); err == nil {
			tableReg, err := regexp.Compile(" ('|`|\"| )" + stmt.Table + "('|`|\"| ) ")
			if err != nil {
				return err
			}

			createSQL = tableReg.ReplaceAllString(createSQL, fmt.Sprintf(" `%v` ", newTableName))
			createSQL = reg.ReplaceAllString(createSQL, "")

			var columns []string
			columnTypes, _ := m.DB.Migrator().ColumnTypes(value)
			for _, columnType := range columnTypes {
				if columnType.Name() != name {
					columns = append(columns, fmt.Sprintf("`%v`", columnType.Name()))
				}
			}

			return m.DB.Transaction(func(tx *gorm.DB) error {
				queries := []string{
					createSQL,
					fmt.Sprintf("INSERT INTO `%v`(%v) SELECT %v FROM `%v`", newTableName, strings.Join(columns, ","), strings.Join(columns, ","), stmt.Table),
					fmt.Sprintf("DROP TABLE `%v`", stmt.Table),
					fmt.Sprintf("ALTER TABLE `%v` RENAME TO `%v`", newTableName, stmt.Table),
				}
				for _, query := range queries {
					if err := tx.Exec(query).Error; err != nil {
						return err
					}
				}
				return nil
			})
		} else {
			return err
		}
	})
}

func (m Migrator) CreateConstraint(interface{}, string) error {
	return ErrConstraintsNotImplemented
}

func (m Migrator) DropConstraint(interface{}, string) error {
	return ErrConstraintsNotImplemented
}

func (m Migrator) HasConstraint(value interface{}, name string) bool {
	var count int64
	m.RunWithValue(value, func(stmt *gorm.Statement) error {
		m.DB.Raw(
			"SELECT count(*) FROM sqlite_master WHERE type = ? AND tbl_name = ? AND (sql LIKE ? OR sql LIKE ? OR sql LIKE ?)",
			"table", stmt.Table, `%CONSTRAINT "`+name+`" %`, `%CONSTRAINT `+name+` %`, "%CONSTRAINT `"+name+"`%",
		).Row().Scan(&count)

		return nil
	})

	return count > 0
}

func (m Migrator) CurrentDatabase() (name string) {
	var null interface{}
	m.DB.Raw("PRAGMA database_list").Row().Scan(&null, &name, &null)
	return
}

func (m Migrator) BuildIndexOptions(opts []schema.IndexOption, stmt *gorm.Statement) (results []interface{}) {
	for _, opt := range opts {
		str := stmt.Quote(opt.DBName)
		if opt.Expression != "" {
			str = opt.Expression
		}

		if opt.Collate != "" {
			str += " COLLATE " + opt.Collate
		}

		if opt.Sort != "" {
			str += " " + opt.Sort
		}
		results = append(results, clause.Expr{SQL: str})
	}
	return
}

func (m Migrator) CreateIndex(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if idx := stmt.Schema.LookIndex(name); idx != nil {
			opts := m.BuildIndexOptions(idx.Fields, stmt)
			values := []interface{}{clause.Column{Name: idx.Name}, clause.Table{Name: stmt.Table}, opts}

			createIndexSQL := "CREATE "
			if idx.Class != "" {
				createIndexSQL += idx.Class + " "
			}
			createIndexSQL += "INDEX ?"

			if idx.Type != "" {
				createIndexSQL += " USING " + idx.Type
			}
			createIndexSQL += " ON ??"

			if idx.Where != "" {
				createIndexSQL += " WHERE " + idx.Where
			}

			return m.DB.Exec(createIndexSQL, values...).Error
		}

		return fmt.Errorf("failed to create index with name %v", name)
	})
}

func (m Migrator) HasIndex(value interface{}, name string) bool {
	var count int
	m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if idx := stmt.Schema.LookIndex(name); idx != nil {
			name = idx.Name
		}

		if name != "" {
			m.DB.Raw(
				"SELECT count(*) FROM sqlite_master WHERE type = ? AND tbl_name = ? AND name = ?", "index", stmt.Table, name,
			).Row().Scan(&count)
		}
		return nil
	})
	return count > 0
}

func (m Migrator) RenameIndex(value interface{}, oldName, newName string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		var sql string
		m.DB.Raw("SELECT sql FROM sqlite_master WHERE type = ? AND tbl_name = ? AND name = ?", "index", stmt.Table, oldName).Row().Scan(&sql)
		if sql != "" {
			return m.DB.Exec(strings.Replace(sql, oldName, newName, 1)).Error
		}
		return fmt.Errorf("failed to find index with name %v", oldName)
	})
}

func (m Migrator) DropIndex(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if idx := stmt.Schema.LookIndex(name); idx != nil {
			name = idx.Name
		}

		return m.DB.Exec("DROP INDEX ?", clause.Column{Name: name}).Error
	})
}
