// See License.txt in main repository directory

// dbMysql is the plugin for MySQL and  contains the
// data and functions specific for SQLite to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.

package dbMysql

import (
	"../dbData"
)

var tds	= dbData.TypeDefns {
	{Name:"date", 		Html:"date", 		Sql:"DATE", 		Go:"string",	DftLen:0,},
	{Name:"datetime",	Html:"datetime",	Sql:"DATETIME",		Go:"string",	DftLen:0,},
	{Name:"email", 		Html:"email", 		Sql:"NVARCHAR", 	Go:"string",	DftLen:50,},
	{Name:"dec", 		Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"decimal", 	Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"int", 		Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"integer", 	Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"money", 		Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"number", 	Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"tel", 		Html:"tel",			Sql:"NVARCHAR",		Go:"string",	DftLen:19,},	//+nnn (nnn) nnn-nnnn
	{Name:"text", 		Html:"text",		Sql:"NVARCHAR",		Go:"string",	DftLen:0,},
	{Name:"time", 		Html:"time",		Sql:"TIME",			Go:"string",	DftLen:0,},
	{Name:"url", 		Html:"url",			Sql:"NVARCHAR",		Go:"string",	DftLen:50,},
}

func ImportString() string {
	return "\"github.com/go-sql-driver/mysql\""
}

func init() {
	pd :=  dbData.Plugin_Data{"mysql", &tds, ImportString, false, true}
	dbData.Register(&pd)
}

