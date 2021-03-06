// See License.txt in main repository directory

// dbGener provides the functions to generate the go statements
// necessary to access and manipulate the SQL databases defined
// by the user. The problem that it tries to solve is that while
// SQL is supposed to be a universal language. It unfortunately
// is not and each type of database manager must be handled slightly
// differently.

// We give this package access to user defined JSON and the ap-
// propriate plugin for the data being processed. Between those
// two resources, it must generate the go code.

package dbGener

import (
	"fmt"
	"genapp/pkg/genSqlAppGo/dbJson"
	"genapp/pkg/genSqlAppGo/dbPlugin"
	"genapp/pkg/sharedData"
	"log"
	"strings"

	"github.com/2kranki/go_util"
)

//============================================================================
//                        	Interface Support
//============================================================================

// dbGener uses interfaces to determine what a plugin can do or not do and when it
// should be called.  If the plugin does not support a particular interface, then
// dbSql will perform default logic to handle the situation.
//
// The reason for all this is that even though Go uses a "common" interface for
// accessing SQL Servers. The SQL, itself, can vary.  Although SQL is supposed to
// to be a standard, it is not consistently implemented unforturnately.
//
// Functions that return a full SQL statement must return a slice of strings even
// if there is only one statement ever generated.  That is because some servers
// such as Microsoft's SQL Server may not do anything until an additional statement
// is issued such as "go".

//----------------------------------------------------------------------------
//               Exec/Query Error Processing Interface Support
//----------------------------------------------------------------------------

type GenExecErrorChecker interface {
	GenExecErrorCheck(db *dbJson.Database) string
}

type GenQueryErrorChecker interface {
	GenQueryErrorCheck(db *dbJson.Database) string
}

//----------------------------------------------------------------------------
//                        	Database SQL Interface Support
//----------------------------------------------------------------------------

type GenDatabaseCreateStmter interface {
	GenDatabaseCreateStmt(db *dbJson.Database) string
}

type GenDatabaseDeleteStmter interface {
	GenDatabaseDeleteStmt(db *dbJson.Database) string
}

//----------------------------------------------------------------------------
//                        	Table SQL Interface Support
//----------------------------------------------------------------------------

type GenTableCountStmter interface {
	GenTableCountStmt(tb *dbJson.DbTable) string
}

type GenTableCreateStmter interface {
	GenTableCreateStmt(tb *dbJson.DbTable) string
}

type GenTableDeleteStmter interface {
	GenTableDeleteStmt(tb *dbJson.DbTable) string
}

//----------------------------------------------------------------------------
//                        	Row SQL Interface Support
//----------------------------------------------------------------------------

// GenRowLimiter defines the interface for generating the LIMIT n option on
// SELECT.  LIMIT is used in general SQL, but not supported by T-SQL (Microsoft).
type GenRowLimiter interface {
	GenRowLimit(tb *dbJson.DbTable, n string) string
}

// GenRowOffseter defines the interface for generating the OFFSET n option on
// SELECT.  OFFSET has a different grammar on T-SQL (Microsoft).
type GenRowOffseter interface {
	GenRowOffset(tb *dbJson.DbTable, n string) string
}

type GenRowDeleteStmter interface {
	GenRowDeleteStmt(tb *dbJson.DbTable) string
}

type GenRowFindStmter interface {
	GenRowFindStmt(tb *dbJson.DbTable) string
}

type GenRowFirstStmter interface {
	GenRowFirstStmt(tb *dbJson.DbTable) string
}

type GenRowInsertStmter interface {
	GenRowInsertStmt(tb *dbJson.DbTable) string
}

type GenRowLastStmter interface {
	GenRowLastStmt(tb *dbJson.DbTable) string
}

type GenRowNextStmter interface {
	GenRowNextStmt(tb *dbJson.DbTable) string
}

type GenRowPageStmter interface {
	GenRowPageStmt(tb *dbJson.DbTable) string
}

type GenRowPrevStmter interface {
	GenRowPrevStmt(tb *dbJson.DbTable) string
}

type GenRowUpdateStmter interface {
	GenRowUpdateStmt(tb *dbJson.DbTable) string
}

//----------------------------------------------------------------------------
//                        	Form Interface Support
//----------------------------------------------------------------------------

type GenFormDataDisplayer interface {
	GenFormDataDisplay(tb *dbJson.DbTable) string
}

type GenFormDataKeyGetter interface {
	GenFormDataKeyGet(tb *dbJson.DbTable) string
}

type GenFormDataKeyser interface {
	GenFormDataKeys(tb *dbJson.DbTable) string
}

//----------------------------------------------------------------------------
//                        	Miscellaneous Interface Support
//----------------------------------------------------------------------------

// GenPlaceHolderer includes the methods responsible for generating place holders
// for SQL statements.
// Note: Some SQL Drivers use '?' as a placeholder and others use "$nn" where
// nn is a number starting at 1..(number of columns). Since all placeholders
// are consistant per driver, we combined the interface methods into one inter-
// face.

type GenPlaceHolderer interface {
	// GenDataPlaceHolder generates the string for table columns when a list of them
	// is involved such as used in RowInsert().  Example: "?, ?, ?"
	GenDataPlaceHolder(tb *dbJson.DbTable) string

	// GenKeySearchPlaceHolder generates the string for multiple keys when an expression
	// is involved such as used in RowFind(). The expression is constrolled by rel which
	// defines it. Possibilities are '=', '<', '>', ... and will apply to all keys in the
	// table. Example: "key1 = $1 AND key2 = $2"
	GenKeySearchPlaceHolder(tb *dbJson.DbTable, rel string) string

	// GenKeysPlaceHolder generates the string for multiple keys when a list of key
	// is involved such as used in RowFind().  Example: "$1, $2, $3"
	GenKeysPlaceHolder(tb *dbJson.DbTable) string
}

//----------------------------------------------------------------------------
//               Exec/Query Error Processing Interface Support
//----------------------------------------------------------------------------

func GenExecErrorCheck(db *dbJson.Database) string {
	var str util.StringBuilder
	var intr GenExecErrorChecker
	var ok bool

	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenExecErrorChecker)
	if ok {
		return intr.GenExecErrorCheck(db)
	}

	return str.String()
}

func GenQueryErrorCheck(db *dbJson.Database) string {
	var str util.StringBuilder
	var intr GenQueryErrorChecker
	var ok bool

	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenQueryErrorChecker)
	if ok {
		return intr.GenQueryErrorCheck(db)
	}

	return str.String()
}

//----------------------------------------------------------------------------
//						Global Database Support Functions
//----------------------------------------------------------------------------

func GenDatabaseCreateStmt(db *dbJson.Database) string {
	var str strings.Builder
	var intr GenDatabaseCreateStmter
	var ok bool

	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenDatabaseCreateStmter)
	if ok {
		return intr.GenDatabaseCreateStmt(db)
	}

	str.WriteString("\tstr.WriteStringf(\"CREATE DATABASE IF NOT EXISTS %s;\", dbName)\n")

	return str.String()
}

func GenDatabaseDeleteStmt(db *dbJson.Database) string {
	var str strings.Builder
	var intr GenDatabaseDeleteStmter
	var ok bool

	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenDatabaseDeleteStmter)
	if ok {
		return intr.GenDatabaseDeleteStmt(db)
	}

	str.WriteString("\tstr.WriteStringf(\"DROP DATABASE IF EXISTS %s;\", dbName)\n")

	return str.String()
}

//----------------------------------------------------------------------------
//						Global Table Support Functions
//----------------------------------------------------------------------------

func GenTableCountStmt(t *dbJson.DbTable) string {
	var str strings.Builder
	var intr GenTableCountStmter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenTableCountStmter)
	if ok {
		return intr.GenTableCountStmt(t)
	}

	fmt.Fprintf(&str, "SELECT COUNT(*) FROM %s%s;\\n", db.Schema, t.Name)

	return str.String()
}

func GenTableCreateStmt(t *dbJson.DbTable) string {
	var str strings.Builder
	var intr GenTableCreateStmter
	var ok bool
	var hasKeys bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenTableCreateStmter)
	if ok {
		return intr.GenTableCreateStmt(t)
	}

	fmt.Fprintf(&str, "CREATE TABLE IF NOT EXISTS %s%s (\\n", db.Schema, t.Name)
	for i, _ := range t.Fields {
		var cm string
		var f *dbJson.DbField
		var ft string
		var incr string
		var nl string
		var pk string
		var sp string

		f = &t.Fields[i]
		cm = ""
		if i != (len(t.Fields) - 1) {
			cm = ","
		} else {
			if hasKeys {
				cm = ","
			}
		}

		td := f.Typ
		if td == nil {
			log.Fatalln("Error - Could not find Type definition for field,",
				f.Name, "type:", f.TypeDefn)
		}
		tdd := f.Typ.SqlType()

		if f.Len > 0 {
			if f.Dec > 0 {
				ft = fmt.Sprintf("%s(%d,%d)", tdd, f.Len, f.Dec)
			} else {
				ft = fmt.Sprintf("%s(%d)", tdd, f.Len)
			}
		} else {
			ft = tdd
		}
		nl = " NOT NULL"
		if f.Nullable {
			nl = ""
		}
		incr = ""
		if f.Incr {
			if db.SqlType == "postgres" {
				ft = "SERIAL"
			} else if db.SqlType == "sqlite" {
				incr = " AUTOINCREMENT"
			} else {
				incr = " AUTO_INCREMENT"
			}
		}
		pk = ""
		if f.KeyNum > 0 {
			if t.KeyCount() == 1 && i == 0 && db.SqlType == "sqlite" {
				pk = " PRIMARY KEY"
			} else {
				hasKeys = true
			}
		}
		sp = ""
		if len(f.SQLParms) > 0 {
			sp = " " + f.SQLParms
		}

		fmt.Fprintf(&str, "\\t%s\\t%s%s%s%s%s%s\\n", f.Name, ft, nl, pk, incr, cm, sp)
	}
	if hasKeys {
		wrk := fmt.Sprintf("\\tCONSTRAINT PK_%s PRIMARY KEY(%s)\\n", t.Name, t.KeysList("", ""))
		str.WriteString(wrk)
	}
	str.WriteString(")")
	if len(t.SQLParms) > 0 {
		str.WriteString(",\\n")
		for _, l := range t.SQLParms {
			fmt.Fprintf(&str, "%s\\n", l)
		}
	}
	str.WriteString(";\\n")

	return str.String()
}

func GenTableDeleteStmt(t *dbJson.DbTable) string {
	var str strings.Builder
	var intr GenTableDeleteStmter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenTableDeleteStmter)
	if ok {
		return intr.GenTableDeleteStmt(t)
	}

	fmt.Fprintf(&str, "DROP TABLE IF EXISTS %s%s;\\n", db.Schema, t.Name)

	return str.String()
}

//----------------------------------------------------------------------------
//						Global Row Support Functions
//----------------------------------------------------------------------------

func GenRowDeleteStmt(t *dbJson.DbTable) string {
	var str strings.Builder
	var intr GenRowDeleteStmter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenRowDeleteStmter)
	if ok {
		return intr.GenRowDeleteStmt(t)
	}

	//TODO: Finish Row Delete SQL
	fmt.Fprintf(&str, "DELETE FROM %s%s WHERE %s;\\n", db.Schema, t.Name, GenKeySearchPlaceHolder(t, "="))

	return str.String()
}

func GenRowFindStmt(t *dbJson.DbTable) string {
	var str strings.Builder
	var intr GenRowFindStmter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenRowFindStmter)
	if ok {
		return intr.GenRowFindStmt(t)
	}

	fmt.Fprintf(&str, "SELECT * FROM %s%s WHERE %s;\\n", db.Schema, t.Name, GenKeySearchPlaceHolder(t, "="))

	return str.String()
}

func GenRowFirstStmt(t *dbJson.DbTable) string {
	var str util.StringBuilder
	var intr GenRowFirstStmter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenRowFirstStmter)
	if ok {
		return intr.GenRowFirstStmt(t)
	}

	str.WriteStringf("SELECT * FROM %s%s ORDER BY %s %s;\\n",
		db.Schema, t.Name, t.KeysList("", " ASC"),
		GenRowLimit(t, "1"))

	return str.String()
}

func GenRowInsertStmt(t *dbJson.DbTable) string {
	var str util.StringBuilder
	var intr GenRowInsertStmter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenRowInsertStmter)
	if ok {
		return intr.GenRowInsertStmt(t)
	}

	str.WriteStringf("INSERT INTO %s%s (%s) VALUES (%s);\\n",
		db.Schema, t.Name, t.InsertNameList(""), GenDataPlaceHolder(t))

	return str.String()
}

func GenRowLastStmt(t *dbJson.DbTable) string {
	var str util.StringBuilder
	var intr GenRowLastStmter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenRowLastStmter)
	if ok {
		return intr.GenRowLastStmt(t)
	}

	str.WriteStringf("SELECT * FROM %s%s ORDER BY %s %s;\\n",
		db.Schema, t.Name, t.KeysList("", " DESC"),
		GenRowLimit(t, "1"))

	return str.String()
}

// GenRowLimit defines the interface for generating the LIMIT n option on
// SELECT.  LIMIT is used in general SQL, but not supported by T-SQL (Microsoft).
func GenRowLimit(t *dbJson.DbTable, n string) string {
	var str util.StringBuilder
	var intr GenRowLimiter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenRowLimiter)
	if ok {
		return intr.GenRowLimit(t, n)
	}

	str.WriteStringf("LIMIT %s", n)

	return str.String()
}

func GenRowNextStmt(t *dbJson.DbTable) string {
	var str util.StringBuilder
	var intr GenRowNextStmter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenRowNextStmter)
	if ok {
		return intr.GenRowNextStmt(t)
	}

	str.WriteStringf("SELECT * FROM %s%s WHERE %s ORDER BY %s %s;\\n",
		db.Schema, t.Name, GenKeySearchPlaceHolder(t, ">"), t.KeysList("", " ASC"),
		GenRowLimit(t, "1"))

	return str.String()
}

// GenRowOffset defines the interface for generating the OFFSET n option on
// SELECT.  OFFSET has a slightly different grammar on T-SQL (Microsoft).
func GenRowOffset(t *dbJson.DbTable, n string) string {
	var str util.StringBuilder
	var intr GenRowOffseter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenRowOffseter)
	if ok {
		return intr.GenRowOffset(t, n)
	}

	str.WriteStringf("OFFSET %s", n)

	return str.String()
}

func GenRowPageStmt(t *dbJson.DbTable) string {
	var str util.StringBuilder
	var intr GenRowPageStmter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenRowPageStmter)
	if ok {
		return intr.GenRowPageStmt(t)
	}

	str.WriteStringf("SELECT * FROM %s%s ORDER BY %s %s %s;\\n",
		db.Schema, t.Name, t.KeysList("", " ASC"),
		GenRowLimit(t, "?"), GenRowOffset(t, "?"))

	return str.String()
}

func GenRowPrevStmt(t *dbJson.DbTable) string {
	var str util.StringBuilder
	var intr GenRowPrevStmter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenRowPrevStmter)
	if ok {
		return intr.GenRowPrevStmt(t)
	}

	str.WriteStringf("SELECT * FROM %s%s WHERE %s ORDER BY %s %s;\\n",
		db.Schema, t.Name, GenKeySearchPlaceHolder(t, "<"), t.KeysList("", " DESC"),
		GenRowLimit(t, "1"))

	return str.String()
}

func GenRowUpdateStmt(t *dbJson.DbTable) string {
	var str util.StringBuilder
	var intr GenRowUpdateStmter
	var ok bool

	db := t.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenRowUpdateStmter)
	if ok {
		return intr.GenRowUpdateStmt(t)
	}

	//TODO: Finish Row Update SQL
	str.WriteStringf("INSERT INTO %s%s ([[.Table.CreateInsertStr]]) VALUES ([[.Table.CreateValueStr]]);\\n",
		db.Schema, t.Name)

	return str.String()
}

//----------------------------------------------------------------------------
//						Global Form Functions
//----------------------------------------------------------------------------

func GenFormDataDisplay(tb *dbJson.DbTable) string {
	var str strings.Builder
	var lbl string
	var m string
	var intr GenFormDataDisplayer
	var ok bool
	var keys []string
	var err error

	db := tb.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenFormDataDisplayer)
	if ok {
		return intr.GenFormDataDisplay(tb)
	}

	// Put non-hidden fields in a table to align columns
	str.WriteString("<table>\n")
	for _, f := range tb.Fields {

		if !f.Hidden {
			tdd := f.Typ.Html
			if len(f.Label) > 0 {
				lbl = strings.Title(f.Label)
			} else {
				lbl = strings.Title(f.Name)
			}
			switch f.Typ.GoType() {
			case "float64":
				m = "m=\"0\" step=\"0.01\" "
			default:
				m = ""
			}
			fmt.Fprintf(&str, "\t<tr><td><label>%s</label></td> <td><input type=\"%s\" name=\"%s\" id=\"%s\" %svalue=\"{{.Rcd.%s}}\"></td></tr>\n",
				lbl, tdd, f.TitledName(), f.TitledName(), m, f.TitledName())
		}
	}
	str.WriteString("</table>\n")

	// Process Hidden fields outside of the table
	for _, f := range tb.Fields {
		if f.Hidden {
			//tdd := f.Typ.Html
			if len(f.Label) > 0 {
				lbl = strings.Title(f.Label)
			} else {
				lbl = strings.Title(f.Name)
			}
			switch f.Typ.GoType() {
			case "float64":
				m = "m=\"0\" step=\"0.01\" "
			default:
				m = ""
			}
			fmt.Fprintf(&str, "\t<input type=\"hidden\" name=\"%s\" id=\"%s\" %svalue=\"{{.Rcd.%s}}\">\n",
				f.TitledName(), f.TitledName(), m, f.TitledName())
		}
	}

	// Process the key fields
	if keys, err = tb.Keys(); err != nil {
		panic("GenFormDataDisplay: error getting keys!")
	}
	for i, fn := range keys {
		f := tb.FindField(fn)
		if f == nil {
			panic(fmt.Sprintf("GenFormDataDisplay: error finding key: %s!", fn))
		}
		//tdd := f.Typ.Html
		if len(f.Label) > 0 {
			lbl = strings.Title(f.Label)
		} else {
			lbl = strings.Title(f.Name)
		}
		switch f.Typ.GoType() {
		case "float64":
			m = "m=\"0\" step=\"0.01\" "
		default:
			m = ""
		}
		fmt.Fprintf(&str, "<input type=\"hidden\" id=\"key%d\" name=\"key%d\"%svalue=\"{{.Rcd.%s}}\">\n",
			i, i, m, f.TitledName())
	}

	return str.String()
}

func GenFormDataKeyGet(tb *dbJson.DbTable) string {
	var str strings.Builder
	var intr GenFormDataKeyGetter
	var ok bool
	var keys []string
	var err error

	db := tb.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenFormDataKeyGetter)
	if ok {
		return intr.GenFormDataKeyGet(tb)
	}

	// Process the key fields
	if keys, err = tb.Keys(); err != nil {
		panic("GenFormDataDisplay: error getting keys!")
	}
	for i, _ := range keys {
		fmt.Fprintf(&str, "\t\t\tkey%d = document.getElementById(\"key%d\").value\n", i, i)
	}

	return str.String()
}

func GenFormDataKeys(tb *dbJson.DbTable) string {
	var err error
	var str strings.Builder
	var intr GenFormDataKeyser
	var ok bool
	var keys []string

	db := tb.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenFormDataKeyser)
	if ok {
		return intr.GenFormDataKeys(tb)
	}

	// Process the key fields
	if keys, err = tb.Keys(); err != nil {
		panic("GenFormDataDisplay: error getting keys!")
	}
	if len(keys) > 0 {
		str.WriteString("\"?\"")
	}
	for i, _ := range keys {
		fmt.Fprintf(&str, "+\"key=\"+key%d", i)
		//tdd := f.Typ.Html
		if i < len(keys)-1 {
			str.WriteString("+\",\"+")
		}
	}

	return str.String()
}

//----------------------------------------------------------------------------
//                        Miscellaneous Interface Support
//----------------------------------------------------------------------------

// GenDataPlaceHolder generates the string for table columns when a list of them
// is involved such as used in RowInsert().  Example: "?, ?, ?"
func GenDataPlaceHolder(tb *dbJson.DbTable) string {
	var str util.StringBuilder
	var cnt int
	var intr GenPlaceHolderer
	var ok bool

	db := tb.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenPlaceHolderer)
	if ok {
		return intr.GenDataPlaceHolder(tb)
	}

	// Accumulate field name count.
	for _, f := range tb.Fields {
		if !f.Incr {
			cnt++
		}
	}

	for i := 0; i < cnt; i++ {
		cm := ", "
		if i == cnt-1 {
			cm = ""
		}
		str.WriteStringf("?%s", cm)
		//str.WriteStringf("$%d%s", i+1, cm)
	}

	return str.String()
}

// GenKeySearchPlaceHolder generates the string for multiple keys when an expression
// is involved such as used in RowFind(). The expression will always be '=' and will
// apply to all keys in the table. Example: "key1 = $1 AND key2 = $2"
func GenKeySearchPlaceHolder(tb *dbJson.DbTable, rel string) string {
	var intr GenPlaceHolderer
	var ok bool

	db := tb.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenPlaceHolderer)
	if ok {
		return intr.GenKeySearchPlaceHolder(tb, rel)
	}

	insertStr := ""
	keys, _ := tb.Keys()
	for i, _ := range keys {
		cm := " AND "
		if i == len(keys)-1 {
			cm = ""
		}
		insertStr += fmt.Sprintf("%s %s ?%s", keys[i], rel, cm)
	}

	return insertStr
}

// GenKeysPlaceHolder generates the string for multiple keys when a list of key
// is involved such as used in RowFind().  Example: "?, ?, ?"
func GenKeysPlaceHolder(tb *dbJson.DbTable) string {
	var intr GenPlaceHolderer
	var ok bool

	db := tb.DB
	pluginData := db.Plugin.(dbPlugin.PluginData)
	plugin := pluginData.Plugin
	intr, ok = plugin.(GenPlaceHolderer)
	if ok {
		return intr.GenKeysPlaceHolder(tb)
	}

	insertStr := ""
	keys, _ := tb.Keys()
	for i := 0; i < len(keys); i++ {
		cm := ", "
		if i == len(tb.Fields)-1 {
			cm = ""
		}
		insertStr += fmt.Sprintf("?%s", cm)
		//insertStr += fmt.Sprintf("$%d%s", i+1, cm)
	}
	return insertStr
}

//----------------------------------------------------------------------------
//							Global Support Functions
//----------------------------------------------------------------------------

// init() is called before main(). Here we define the functions that will be
// used in the templates.
func init() {
	sharedData.SetFunc("GenExecErrorCheck", GenExecErrorCheck)
	sharedData.SetFunc("GenQueryErrorCheck", GenQueryErrorCheck)
	sharedData.SetFunc("GenDatabaseCreateStmt", GenDatabaseCreateStmt)
	sharedData.SetFunc("GenDatabaseDeleteStmt", GenDatabaseDeleteStmt)
	sharedData.SetFunc("GenTableCountStmt", GenTableCountStmt)
	sharedData.SetFunc("GenTableCreateStmt", GenTableCreateStmt)
	sharedData.SetFunc("GenTableDeleteStmt", GenTableDeleteStmt)
	sharedData.SetFunc("GenRowDeleteStmt", GenRowDeleteStmt)
	sharedData.SetFunc("GenRowFindStmt", GenRowFindStmt)
	sharedData.SetFunc("GenRowFirstStmt", GenRowFirstStmt)
	sharedData.SetFunc("GenRowInsertStmt", GenRowInsertStmt)
	sharedData.SetFunc("GenRowLastStmt", GenRowLastStmt)
	sharedData.SetFunc("GenRowNextStmt", GenRowNextStmt)
	sharedData.SetFunc("GenRowPageStmt", GenRowPageStmt)
	sharedData.SetFunc("GenRowPrevStmt", GenRowPrevStmt)
	sharedData.SetFunc("GenRowUpdateStmt", GenRowUpdateStmt)
	sharedData.SetFunc("GenFormDataDisplay", GenFormDataDisplay)
	sharedData.SetFunc("GenFormDataKeyGet", GenFormDataKeyGet)
	sharedData.SetFunc("GenFormDataKeys", GenFormDataKeys)
}
