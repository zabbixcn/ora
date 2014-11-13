### Golang Oracle Database Driver ###

Package ora implements an Oracle database driver for the Go programming language.

An Oracle database may be accessed through the database/sql package or through the
ora package directly. database/sql offers connection pooling, thread safety,
a consistent API to multiple database technologies and a common set of Go types.
The ora package offers additional features including pointers, slices, nullable
types, numerics of various sizes, Oracle-specific types, Go return type configuration,
and Oracle abstractions such as environment, server and session.

The ora package is written with the Oracle Call Interface (OCI) C-language
libraries provided by Oracle. The OCI libraries are a standard for client
application communication and driver communication with Oracle databases.

The ora package has been verified to work with Oracle Enterprise 12c (12.1.0.1.0),
Windows 8.1 and 64-bit x86.

---
* [Installation](https://github.com/ranaian/ora#installation)
* [Data Types](https://github.com/ranaian/ora#data-types)
* [SQL Placeholder Syntax](https://github.com/ranaian/ora#sql-placeholder-syntax)
* [Working With The Sql Package](https://github.com/ranaian/ora#working-with-the-sql-package)
* [Working With The Oracle Package Directly](https://github.com/ranaian/ora#working-with-the-oracle-package-directly)
* [Test Database Setup](https://github.com/ranaian/ora#test-database-setup)
* [Limitations](https://github.com/ranaian/ora#limitations)
* [License](https://github.com/ranaian/ora#license)
* [API Reference](http://godoc.org/github.com/ranaian/ora#pkg-index)

---

##### Installation #####

Minimum requirements are Go 1.3 with CGO enabled, a GCC C compiler, and
Oracle 12c (12.1.0.1.0) or Oracle Instant Client (12.1.0.1.0).

Install Oracle 12c or Oracle Instant Client.

Set the CGO_CFLAGS and CGO_LDFLAGS environment variables to locate the OCI headers
and library. For example:

	// example OS environment variables for Oracle 12c on Windows
	CGO_CFLAGS=-Ic:/oracle/home/OCI/include/
	CGO_LDFLAGS=c:/oracle/home/BIN/oci.dll

CGO_CFLAGS identifies the location of the OCI header file. CGO_LDFLAGS identifies
the location of the OCI library. These locations will vary based on whether an Oracle
database is locally installed or whether the Oracle instant client libraries are
locally installed.

The ora package uses glog for logging. The ora package and glog package are available
from GitHub:

	go get github.com/ranaian/ora
	go get github.com/golang/glog

##### Data Types #####

The ora package supports all built-in Oracle data types. The supported Oracle
built-in data types are NUMBER, BINARY_DOUBLE, BINARY_FLOAT, FLOAT, DATE,
TIMESTAMP, TIMESTAMP WITH TIME ZONE, TIMESTAMP WITH LOCAL TIME ZONE,
INTERVAL YEAR TO MONTH, INTERVAL DAY TO SECOND, CHAR, NCHAR, VARCHAR, VARCHAR2,
NVARCHAR2, LONG, CLOB, NCLOB, BLOB, LONG RAW, RAW, ROWID and BFILE.
SYS_REFCURSOR is also supported.

Oracle does not provide a built-in boolean type. Oracle provides a single-byte
character type. A common practice is to define two single-byte characters which
represent `true` and `false`. The ora package adopts this approach. The oracle
package associates a Go `bool` value to a Go rune and sends and receives the
rune to a CHAR(1 BYTE) column or CHAR(1 CHAR) column.

The default `false` rune is zero `0`. The default `true` rune is one `1`. The `bool` rune
association may be configured or disabled when directly using the ora package
but not with the database/sql package.

##### SQL Placeholder Syntax #####

Within a SQL string a placeholder may be specified to indicate where a Go variable
is placed. The SQL placeholder is an Oracle identifier, from 1 to 30
characters, prefixed with a colon `:`. For example:

	// example Oracle placeholder uses a colon
	INSERT INTO T1 (C1) VALUES (:C1)

Placeholders within a SQL statement are bound by position. The actual name is not
used by the ora package driver e.g., placeholder names `:c1`, `:1`, or `:xyz` are
treated equally.

##### Working With The Sql Package #####

You may access an Oracle database through the database/sql package. The database/sql
package offers a consistent API across different databases, connection
pooling, thread safety and a set of common Go types. database/sql makes working
with Oracle straight-forward.

The ora package implements interfaces in the database/sql/driver package enabling
database/sql to communicate with an Oracle database. Using database/sql
ensures you never have to call the ora package directly.

When using database/sql, the mapping between Go types and Oracle types is immutable.
The Go-to-Oracle type mapping for database/sql is:

	Go type		Oracle type

	int64		NUMBER°, BINARY_DOUBLE, BINARY_FLOAT, FLOAT

	float64		NUMBER¹, BINARY_DOUBLE, BINARY_FLOAT, FLOAT

	time.Time	TIMESTAMP, TIMESTAMP WITH TIME ZONE, TIMESTAMP WITH LOCAL TIME ZONE, DATE

	string		CHAR², NCHAR, VARCHAR, VARCHAR2, NVARCHAR2, LONG, CLOB, NCLOB

	bool		CHAR(1 BYTE)³, CHAR(1 CHAR)³

	[]byte		BLOB, LONG RAW, RAW


	° A select-list column defined as an Oracle NUMBER with zero scale e.g.,
	NUMBER(10,0) is returned as an int64. Either int64 or float64 may be inserted
	into a NUMBER column with zero scale. float64 insertion will have its fractional
	part truncated.

	¹ A select-list column defined as an Oracle NUMBER with a scale greater than
	zero e.g., NUMBER(10,4) is returned as a float64. Either int64 or float64 may
	be inserted into a NUMBER column with a scale greater than zero.

	² A select-list column defined as an Oracle CHAR with a length greater than 1
	e.g., CHAR(2 BYTE) or CHAR(2 CHAR) is returned as a string. A Go string of any
	length up to the column max length may be inserted into the CHAR column.

	³ The Go bool value false is mapped to the zero rune '0'. The Go bool value
	true is mapped to the one rune '1'.

##### Working With The Oracle Package Directly #####

The ora package allows programming with pointers, slices, nullable types,
numerics of various sizes, Oracle-specific types, Go return type configuration, and
Oracle abstractions such as environment, server and session. When working with the
ora package directly, the API is slightly different than database/sql.

When using the ora package directly, the mapping between Go types and Oracle types
is mutable. The Go-to-Oracle type mapping for the ora package is:

	Go type							Oracle type

	int64, int32, int16, int8		NUMBER°, BINARY_DOUBLE, BINARY_FLOAT, FLOAT
	uint64, uint32, uint16, uint8
	Int64, Int32, Int16, Int8
	Uint64, Uint32, Uint16, Uint8
	*int64, *int32, *int16, *int8
	*uint64, *uint32, *uint16, *uint8
	[]int64, []int32, []int16, []int8
	[]uint64, []uint32, []uint16, []uint8
	[]Int64, []Int32, []Int16, []Int8
	[]Uint64, []Uint32, []Uint16, []Uint8

	float64, float32				NUMBER¹, BINARY_DOUBLE, BINARY_FLOAT, FLOAT
	Float64, Float32
	*float64, *float32
	[]float64, []float32
	[]Float64, []Float32

	time.Time						TIMESTAMP, TIMESTAMP WITH TIME ZONE,
	Time							TIMESTAMP WITH LOCAL TIME ZONE, DATE
	*time.Time
	[]time.Time
	[]Time

	string							CHAR², NCHAR, VARCHAR, VARCHAR2,
	String							NVARCHAR2, LONG, CLOB, NCLOB, ROWID
	*string
	[]string
	[]String

	bool							CHAR(1 BYTE)³, CHAR(1 CHAR)³
	Bool
	*bool
	[]bool
	[]Bool
	
	[]byte							BLOB, LONG RAW, RAW
	[][]byte
	Binary
	[]Binary

	IntervalYM						INTERVAL MONTH TO YEAR
	[]IntervalYM

	IntervalDS						INTERVAL DAY TO SECOND
	[]IntervalDS

	Bfile							BFILE

	° A select-list column defined as an Oracle NUMBER with zero scale e.g.,
	NUMBER(10,0) is returned as an int64 by default. Integer and floating point
	numerics may be inserted into a NUMBER column with zero scale. Inserting a
	floating point numeric will have its fractional part truncated.

	¹ A select-list column defined as an Oracle NUMBER with a scale greater than
	zero e.g., NUMBER(10,4) is returned as a float64 by default. Integer and
	floating point numerics may be inserted into a NUMBER column with a scale
	greater than zero.

	² A select-list column defined as an Oracle CHAR with a length greater than 1
	e.g., CHAR(2 BYTE) or CHAR(2 CHAR) is returned as a string. A Go string of any
	length up to the column max length may be inserted into the CHAR column.

	³ The Go bool value false is mapped to the zero rune '0'. The Go bool value
	true is mapped to the one rune '1'.

An example of using the ora package directly:

```go
package main

import (
	"fmt"
	"github.com/ranaian/ora"
)

func main() {
	// example usage of the ora package driver
	// connect to a server and open a session
	env, _ := ora.GetDrv().OpenEnv()
	defer env.Close()
	srv, err := env.OpenSrv("orcl")
	defer srv.Close()
	if err != nil {
		panic(err)
	}
	ses, err := srv.OpenSes("test", "test")
	defer ses.Close()
	if err != nil {
		panic(err)
	}

	// create table
	tableName := "t1"
	stmtTbl, err := ses.Prep(fmt.Sprintf("CREATE TABLE %v "+
		"(C1 NUMBER(19,0) GENERATED ALWAYS AS IDENTITY "+
		"(START WITH 1 INCREMENT BY 1), C2 VARCHAR2(48 CHAR))", tableName))
	defer stmtTbl.Close()
	if err != nil {
		panic(err)
	}
	rowsAffected, err := stmtTbl.Exec()
	if err != nil {
		panic(err)
	}
	fmt.Println(rowsAffected)

	// begin first transaction
	tx1, err := ses.StartTx()
	if err != nil {
		panic(err)
	}

	// insert record
	var id uint64
	str := "Go is expressive, concise, clean, and efficient."
	stmtIns, err := ses.Prep(fmt.Sprintf(
		"INSERT INTO %v (C2) VALUES (:C2) RETURNING C1 INTO :C1", tableName))
	defer stmtIns.Close()
	rowsAffected, err = stmtIns.Exec(str, &id)
	if err != nil {
		panic(err)
	}
	fmt.Println(rowsAffected)

	// insert nullable String slice
	a := make([]ora.String, 4)
	a[0] = ora.String{Value: "Its concurrency mechanisms make it easy to"}
	a[1] = ora.String{IsNull: true}
	a[2] = ora.String{Value: "It's a fast, statically typed, compiled"}
	a[3] = ora.String{Value: "One of Go's key design goals is code"}
	stmtSliceIns, err := ses.Prep(fmt.Sprintf(
		"INSERT INTO %v (C2) VALUES (:C2)", tableName))
	defer stmtSliceIns.Close()
	if err != nil {
		panic(err)
	}
	rowsAffected, err = stmtSliceIns.Exec(a)
	if err != nil {
		panic(err)
	}
	fmt.Println(rowsAffected)

	// fetch records
	stmtQuery, err := ses.Prep(fmt.Sprintf(
		"SELECT C1, C2 FROM %v", tableName))
	defer stmtQuery.Close()
	if err != nil {
		panic(err)
	}
	rset, err := stmtQuery.Query()
	if err != nil {
		panic(err)
	}
	for rset.Next() {
		fmt.Println(rset.Row[0], rset.Row[1])
	}
	if rset.Err != nil {
		panic(rset.Err)
	}

	// commit first transaction
	err = tx1.Commit()
	if err != nil {
		panic(err)
	}

	// begin second transaction
	tx2, err := ses.StartTx()
	if err != nil {
		panic(err)
	}
	// insert null String
	nullableStr := ora.String{IsNull: true}
	stmtTrans, err := ses.Prep(fmt.Sprintf(
		"INSERT INTO %v (C2) VALUES (:C2)", tableName))
	defer stmtTrans.Close()
	if err != nil {
		panic(err)
	}
	rowsAffected, err = stmtTrans.Exec(nullableStr)
	if err != nil {
		panic(err)
	}
	fmt.Println(rowsAffected)
	// rollback second transaction
	err = tx2.Rollback()
	if err != nil {
		panic(err)
	}

	// fetch and specify return type
	stmtCount, err := ses.Prep(fmt.Sprintf(
		"SELECT COUNT(C1) FROM %v WHERE C2 IS NULL", tableName), ora.U8)
	defer stmtCount.Close()
	if err != nil {
		panic(err)
	}
	rset, err = stmtCount.Query()
	if err != nil {
		panic(err)
	}
	row := rset.NextRow()
	if row != nil {
		fmt.Println(row[0])
	}
	if rset.Err != nil {
		panic(rset.Err)
	}

	// create stored procedure with sys_refcursor
	stmtProcCreate, err := ses.Prep(fmt.Sprintf(
		"CREATE OR REPLACE PROCEDURE PROC1(P1 OUT SYS_REFCURSOR) AS BEGIN "+
			"OPEN P1 FOR SELECT C1, C2 FROM %v WHERE C1 > 2 ORDER BY C1; "+
			"END PROC1;",
		tableName))
	defer stmtProcCreate.Close()
	rowsAffected, err = stmtProcCreate.Exec()
	if err != nil {
		panic(err)
	}

	// call stored procedure
	// pass *Rset to Exec to receive the results of a sys_refcursor
	stmtProcCall, err := ses.Prep("CALL PROC1(:1)")
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset := &ora.Rset{}
	rowsAffected, err = stmtProcCall.Exec(procRset)
	if err != nil {
		panic(err)
	}
	if procRset.IsOpen() {
		for procRset.Next() {
			fmt.Println(procRset.Row[0], procRset.Row[1])
		}
		if procRset.Err != nil {
			panic(procRset.Err)
		}
		fmt.Println(procRset.Len())
	}

	// Output:
	// 0
	// 1
	// 4
	// 1 Go is expressive, concise, clean, and efficient.
	// 2 Its concurrency mechanisms make it easy to
	// 3 <empty>
	// 4 It's a fast, statically typed, compiled
	// 5 One of Go's key design goals is code
	// 1
	// 1
	// 3 <empty>
	// 4 It's a fast, statically typed, compiled
	// 5 One of Go's key design goals is code
	// 3
}
```

Pointers may be used to capture out-bound values from a SQL statement such as
an insert or stored procedure call. For example, a numeric pointer captures an
identity value:

```go
// given:
// CREATE TABLE T1 (
// C1 NUMBER(19,0) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1),
// C2 VARCHAR2(48 CHAR))
var id int64
stmt, err = ses.Prep("INSERT INTO T1 (C2) VALUES ('GO') RETURNING C1 INTO :C1")
stmt.Exec(&id)
```

A `string` pointer captures an out parameter from a stored procedure:

```go
// given:
// CREATE OR REPLACE PROCEDURE PROC1 (P1 OUT VARCHAR2) AS BEGIN P1 := 'GO'; END PROC1;
var str string
stmt, err = ses.Prep("CALL PROC1(:1)")
stmt.Exec(&str)
```

Slices may be used to insert multiple records with a single insert statement:

```go
// insert one million rows with single insert statement
// given: CREATE TABLE T1 (C1 NUMBER)
values := make([]int64, 1000000)
for n, _ := range values {
	values[n] = int64(n)
}
rowsAffected, err := ses.PrepAndExec("INSERT INTO T1 (C1) VALUES (:C1)", values)
```

The ora package provides nullable Go types to support DML operations such as
insert and select. The nullable Go types provided by the ora package are `Int64`,
`Int32`, `Int16`, `Int8`, `Uint64`, `Uint32`, `Uint16`, `Uint8`, `Float64`, `Float32`, `Time`,
`IntervalYM`, `IntervalDS`, `String`, `Bool`, `Bytes` and `Bfile`. For example, you may insert
nullable strings and select nullable strings:

```go
// insert String slice
// given: CREATE TABLE T1 (C1 VARCHAR2(48 CHAR))
a := make([]ora.String, 5)
a[0] = ora.String{Value: "Go is expressive, concise, clean, and efficient."}
a[1] = ora.String{Value: "Its concurrency mechanisms make it easy to"}
a[2] = ora.String{IsNull: true}
a[3] = ora.String{Value: "It's a fast, statically typed, compiled"}
a[4] = ora.String{Value: "One of Go's key design goals is code"}
stmt, err = ses.Prep("INSERT INTO T1 (C1) VALUES (:C1)")
stmt.Exec(a)

// Specify OraS to Prep method to return ora.String values
// fetch records
stmt, err = ses.Prep("SELECT C1 FROM T1", OraS)
rst, err := stmt.Query()
for rst.Next() {
	fmt.Println(rst.Row[0])
}
```

The `Stmt.Prep` method is variadic accepting zero or more `GoColumnType`
which define a Go return type for a select-list column. For example, a `Prep`
call can be configured to return an `int64` and a nullable `Int64` from the same
column:

```go
// given: create table t1 (c1 number)
stmt, err = ses.Prep("SELECT C1, C1 FROM T1", ora.I64, ora.OraI64)
rst, err := stmt.Query()
for rst.Next() {
	fmt.Println(rst.Row[0], rst.Row[1])
}
```

Go numerics of various sizes are supported in DML operations. The ora package
supports `int64`, `int32`, `int16`, `int8`, `uint64`, `uint32`, `uint16`, `uint8`, `float64` and
`float32`. For example, you may insert a `uint16` and select numerics of various sizes:

```go
// insert uint16
// given: create table t1 (c1 number)
value := uint16(9)
stmt, err = ses.Prep("INSERT INTO T1 (C1) VALUES (:C1)")
stmt.Exec(value)

// select numerics of various sizes from the same column
stmt, err = ses.Prep(
	"SELECT C1, C1, C1, C1, C1, C1, C1, C1, C1, C1, FROM T1",
	ora.I64, ora.I32, ora.I16, ora.I8, ora.U64, ora.U32, ora.U16, ora.U8,
	ora.F64, ora.F32)
rst, err := stmt.Query()
row := rst.NextRow()
```

If a non-nullable type is defined for a nullable column returning null, the Go
type's zero value is returned.

GoColumnTypes defined by the ora package are:

	Go type		GoColumnType

	int64		I64

	int32		I32

	int16		I16

	int8		I8

	uint64		U64

	uint32		U32

	uint16		U16

	uint8		U8

	float64		F64

	Int64		OraI64

	Int32		OraI32

	Int16		OraI16

	Int8		OraI8

	Uint64		OraU64

	Uint32		OraU32

	Uint16		OraU16

	Uint8		OraU8

	Float64		OraF64

	Float32		OraF32

	time.Time	T

	Time		OraT

	string		S

	String		OraS

	bool		B

	Bool		OraB

	[]byte		Bin

	Binary		OraBin

	default°	D

	° D represents a default mapping between a select-list column and a Go type.
	The default mapping is defined in RsetConfig.

When `Stmt.Prep` doesn't receive a `GoColumnType`, or receives an incorrect `GoColumnType`,
the default value defined in `RsetConfig` is used.

There are two configuration structs, `StmtConfig` and `RsetConfig`.
`StmtConfig` configures various aspects of a `Stmt`. `RsetConfig` configures
various aspects of `Rset`, including the default mapping between an Oracle select-list
column and a Go type. `StmtConfig` may be set in an `Env`, `Srv`, `Ses`
and `Stmt`. `RsetConfig` may be set in a `StmtConfig`.

Setting `StmtConfig` on `Env`, `Srv`, `Ses`
or `Stmt` cascades the `StmtConfig` to all current and future descendent structs.
An `Env` may contain multiple `Srv`. A `Srv` may contain multiple `Ses`.
A `Ses` may contain multiple `Stmt`. A `Stmt` may contain multiple `Rset`.

	// setting StmtConfig cascades to descendent structs
	// Env -> Srv -> Ses -> Stmt -> Rset

Setting a `RsetConfig` on a `StmtConfig` does not cascade through descendent structs.
Configuration of `Stmt.Config` takes effect prior to calls to `Stmt.Exec` and
`Stmt.Query`; consequently, any updates to `Stmt.Config` after a call to `Stmt.Exec`
or `Stmt.Query` are not observed.

One configuration scenario may be to set a server's select statements to return nullable Go types by
default:

```go
sc := NewStmtConfig()
sc.Rset.SetNumberScaless(ora.OraI64)
sc.Rset.SetNumberScaled(ora.OraF64)
sc.Rset.SetBinaryDouble(ora.OraF64)
sc.Rset.SetBinaryFloat(ora.OraF64)
sc.Rset.SetFloat(ora.OraF64)
sc.Rset.SetDate(ora.OraT)
sc.Rset.SetTimestamp(ora.OraT)
sc.Rset.SetTimestampTz(ora.OraT)
sc.Rset.SetTimestampLtz(ora.OraT)
sc.Rset.SetChar1(ora.OraB)
sc.Rset.SetVarchar(ora.OraS)
sc.Rset.SetLong(ora.OraS)
sc.Rset.SetClob(ora.OraS)
sc.Rset.SetBlob(ora.OraBin)
sc.Rset.SetRaw(ora.OraBin)
sc.Rset.SetLongRaw(ora.OraBin)
srv, err := env.OpenSrv("orcl")
// setting the server StmtConfig will cascade to any open Sess, Stmts
// any new Ses, Stmt will receive this StmtConfig
// any new Rset will receive the StmtConfig.Rset configuration
srv.SetStmtConfig(sc)
```

Another scenario may be to configure the runes mapped to `bool` values:

```go
// update StmtConfig to change the FalseRune and TrueRune inserted into the database
// given: CREATE TABLE T1 (C1 CHAR(1 BYTE))

// insert 'false' record
var falseValue bool = false
stmt, err = ses.Prep("INSERT INTO T1 (C1) VALUES (:C1)")
stmt.Config.FalseRune = 'N'
stmt.Exec(falseValue)

// insert 'true' record
var trueValue bool = true
stmt, err = ses.Prep("INSERT INTO T1 (C1) VALUES (:C1)")
stmt.Config.TrueRune = 'Y'
stmt.Exec(trueValue)

// update RsetConfig to change the TrueRune
// used to translate an Oracle char to a Go bool
// fetch inserted records
stmt, err = ses.Prep("SELECT C1 FROM T1")
stmt.Config.Rset.TrueRune = 'Y'
rst, err := stmt.Query()
for rst.Next() {
	fmt.Println(rst.Row[0])
}
```

Oracle-specific types offered by the ora package are `Rset`, `IntervalYM`, `IntervalDS`, `Binary` and `Bfile`.
`Rset` represents an Oracle SYS_REFCURSOR. `IntervalYM` represents an Oracle INTERVAL YEAR TO MONTH.
`IntervalDS` represents an Oracle INTERVAL DAY TO SECOND. `Binary` represents an Oracle BLOB. And `Bfile` represents an Oracle BFILE. ROWID
columns are returned as strings and don't have a unique Go type.

`Rset` is used to obtain Go values from a SQL select statement. Methods `Rset.Next`,
`Rset.NextRow`, and `Rset.Len` are available. Fields `Rset.Row`, `Rset.Err`,
`Rset.Index`, and `Rset.ColumnNames` are also available. The `Next` method attempts to
load data from an Oracle buffer into `Row`, returning true when successful. When no data is available,
or if an error occurs, `Next` returns false setting `Row` to nil. Any error in `Next` is assigned to `Err`.
Calling `Next` increments `Index` and method `Len` returns the total number of rows processed. The `NextRow`
method is convenient for returning a single row. `NextRow` calls `Next` and returns `Row`.
`ColumnNames` returns the names of columns defined by the SQL select statement.

`Rset` has two usages. `Rset` may be returned from `Stmt.Query` when prepared with a SQL select
statement:

```go
// given: CREATE TABLE T1 (C1 NUMBER, C2, CHAR(1 BYTE), C3 VARCHAR2(48 CHAR))
stmt, err = ses.Prep("SELECT C1, C2, C3 FROM T1")
rst, err := stmt.Query()
for rst.Next() {
	fmt.Println(rst.Index, rst.Row[0], rst.Row[1], rst.Row[2])
}
```

Or, `*Rset` may be passed to `Stmt.Exec` when prepared with a stored procedure accepting
an OUT SYS_REFCURSOR parameter:

```go
// given:
// CREATE TABLE T1 (C1 NUMBER, C2 VARCHAR2(48 CHAR))
// CREATE OR REPLACE PROCEDURE PROC1(P1 OUT SYS_REFCURSOR) AS
// BEGIN OPEN P1 FOR SELECT C1, C2 FROM T1 ORDER BY C1; END PROC1;
stmt, err = ses.Prep("CALL PROC1(:1)")
rst := &ora.Rset{}
stmt.Exec(rst)
if rst.IsOpen() {
	for rst.Next() {
		fmt.Println(rst.Row[0], rst.Row[1])
	}
}
```

Stored procedures with multiple OUT SYS_REFCURSOR parameters enable a single `Exec` call to obtain
multiple `Rsets`:

```go
// given:
// CREATE TABLE T1 (C1 NUMBER, C2 VARCHAR2(48 CHAR))
// CREATE OR REPLACE PROCEDURE PROC1(P1 OUT SYS_REFCURSOR, P2 OUT SYS_REFCURSOR) AS BEGIN
// OPEN P1 FOR SELECT C1 FROM T1 ORDER BY C1; OPEN P2 FOR SELECT C2 FROM T1 ORDER BY C2;
// END PROC1;
stmt, err = ses.Prep("CALL PROC1(:1, :2)")
rst1 := &ora.Rset{}
rst2 := &ora.Rset{}
stmt.Exec(rst1, rst2)
// read from first cursor
if rst1.IsOpen() {
	for rst1.Next() {
		fmt.Println(rst1.Row[0])
	}
}
// read from second cursor
if rst2.IsOpen() {
	for rst2.Next() {
		fmt.Println(rst2.Row[0])
	}
}
```

The types of values assigned to `Row` may be configured in `StmtConfig.Rset`. For configuration
to take effect, assign `StmtConfig.Rset` prior to calling `Stmt.Query` or `Stmt.Exec`.

`Rset` prefetching may be controlled by `StmtConfig.PrefetchRowCount` and
`StmtConfig.PrefetchMemorySize`. `PrefetchRowCount` works in coordination with
`PrefetchMemorySize`. When `PrefetchRowCount` is set to zero only `PrefetchMemorySize` is used;
otherwise, the minimum of `PrefetchRowCount` and `PrefetchMemorySize` is used.
The default uses a `PrefetchMemorySize` of 134MB.

Opening and closing `Rsets` is managed internally. `Rset` does not have an Open method or Close method.

`IntervalYM` may be be inserted and selected:

```go
// insert IntervalYM slice
// given: create table t1 (c1 interval year to month)
a := make([]ora.IntervalYM, 5)
a[0] = ora.IntervalYM{Year: 1, Month: 1}
a[1] = ora.IntervalYM{Year: 99, Month: 9}
a[2] = ora.IntervalYM{IsNull: true}
a[3] = ora.IntervalYM{Year: -1, Month: -1}
a[4] = ora.IntervalYM{Year: -99, Month: -9}
stmt, err = ses.Prep("INSERT INTO T1 (C1) VALUES (:C1)")
stmt.Exec(a)

// query IntervalYM
stmt, err = ses.Prep("SELECT C1 FROM T1")
rst, err := stmt.Query()
for rst.Next() {
	fmt.Println(rst.Row[0])
}
```

`IntervalDS` may be be inserted and selected:

```go
// insert IntervalDS slice
// given: CREATE TABLE T1 (C1 INTERVAL DAY TO SECOND)
a := make([]ora.IntervalDS, 5)
a[0] = ora.IntervalDS{Day: 1, Hour: 1, Minute: 1, Second: 1, Nanosecond: 123456789}
a[1] = ora.IntervalDS{Day: 59, Hour: 59, Minute: 59, Second: 59, Nanosecond: 123456789}
a[2] = ora.IntervalDS{IsNull: true}
a[3] = ora.IntervalDS{Day: -1, Hour: -1, Minute: -1, Second: -1, Nanosecond: -123456789}
a[4] = ora.IntervalDS{Day: -59, Hour: -59, Minute: -59, Second: -59, Nanosecond: -123456789}
stmt, err = ses.Prep("INSERT INTO T1 (C1) VALUES (:C1)")
stmt.Exec(a)

// query IntervalDS
stmt, err = ses.Prep("SELECT C1 FROM T1")
rst, err := stmt.Query()
for rst.Next() {
	fmt.Println(rst.Row[0])
}
```

Transactions on an Oracle server are supported. DML statements auto-commit
unless a transaction has started:

```go
// given: create table t1 (c1 number)

// rollback
tx, err := ses.BeginTransaction()
stmt, err = ses.Prep("INSERT INTO T1 (C1) VALUES (3)")
stmt.Exec()
stmt, err = ses.Prep("INSERT INTO T1 (C1) VALUES (5)")
stmt.Exec()
tx.Rollback()

// commit
tx, err = ses.BeginTransaction()
stmt, err = ses.Prep("INSERT INTO T1 (C1) VALUES (7)")
stmt.Exec()
stmt, err = ses.Prep("INSERT INTO T1 (C1) VALUES (9)")
stmt.Exec()
tx.Commit()

// query records
stmt, err = ses.Prep("SELECT C1 FROM T1")
rst, err := stmt.Query()
for rst.Next() {
	fmt.Println(rst.Row[0])
}
```

The `Srv.Ping` method checks whether the client's connection to an
Oracle server is valid. A call to `Ping` requires an open Ses. `Ping`
will return a nil error when the connection is fine:

```go
// open a session before calling Ping
ses, _ := srv.OpenSes("username", "password")
err := srv.Ping()
if err == nil {
	fmt.Println("Ping sucessful")
}
```

The `Srv.Version` method is available to obtain the Oracle server version. A call
to `Version` requires an open Ses:

```go
// open a session before calling Version
ses, err := srv.OpenSes("username", "password")
version, err := srv.Version()
if version != "" && err == nil {
	fmt.Println("Received version from server")
}
```

Further code examples are available in the samples folder, example file and test files.

##### Test Database Setup #####

Tests are available and require some setup. Setup varies depending on whether
the Oracle server is configured as a container database or non-container database.
It's simpler to setup a non-container database. An example for each setup is
explained.

Non-container test database setup steps:

```sql
-- 1. login to an Oracle server with SqlPlus as sysdba:
sqlplus / as sysdba
```

```sql
-- 2. create a file for the test database use
CREATE TABLESPACE test_ts NOLOGGING DATAFILE 'test.dat' SIZE 100M AUTOEXTEND ON;
```

```sql
-- 3. create a test database
CREATE USER test IDENTIFIED BY test DEFAULT TABLESPACE test_ts;
```

```sql
-- 4. grant permissions to the database
GRANT CREATE SESSION, CREATE TABLE, CREATE SEQUENCE,
CREATE PROCEDURE, UNLIMITED TABLESPACE TO test;
```

```sql
-- 5. create OS environment variables
-- specify your_database_name; varies based on installation; may be 'orcl'
GO_ORA_DRV_TEST_DB = your_database_name
GO_ORA_DRV_TEST_USERNAME = test
GO_ORA_DRV_TEST_PASSWORD = test
```

Container test database setup steps:

```sql
-- 1. login to an Oracle server with SqlPlus as sysdba:
sqlplus / as sysdba
```

```sql
-- 2. create a test pluggable database and permissions
-- you will need to change the FILE_NAME_CONVERT file paths for your database installation
CREATE PLUGGABLE DATABASE go_driver_test
ADMIN USER test IDENTIFIED BY test
ROLES = (DBA)
FILE_NAME_CONVERT = ('d:\oracle\data\orcl\pdbseed\', 'd:\oracle\data\go_driver_test\');
```

```sql
-- 3. modify the pluggable database settings
ALTER PLUGGABLE DATABASE go_driver_test OPEN;
ALTER SESSION SET CONTAINER = go_driver_test;
GRANT DBA TO test;
```

```sql
-- 4. add new database service to the tnsnames.ora file:
-- located on your client machine in $ORACLE_HOME\network\admin\tnsnames.ora
GO_DRIVER_TEST =
  (DESCRIPTION =
	(ADDRESS = (PROTOCOL = TCP)(HOST = localhost)(PORT = 1521))
	(CONNECT_DATA =
	  (SERVER = DEDICATED)
	  (SERVICE_NAME = go_driver_test)
	)
  )
```

```sql
-- 5. create OS environment variables
GO_ORA_DRIVER_TEST_DB = go_driver_test
GO_ORA_DRIVER_TEST_USERNAME = test
GO_ORA_DRIVER_TEST_PASSWORD = test
```

Some helpful SQL maintenance statements:

```sql
-- delete all tables in a non-container database
BEGIN
FOR c IN (SELECT table_name FROM user_tables) LOOP
EXECUTE IMMEDIATE ('DROP TABLE "' || c.table_name || '" CASCADE CONSTRAINTS');
END LOOP;
END;
```

```sql
-- delete the non-container test database; use SqlPlus as sysdba
DROP USER test CASCADE;
```

Run the tests.

##### Limitations #####

database/sql method `Stmt.QueryRow` is not supported.

##### License #####

Copyright 2014 Rana Ian. All rights reserved.
Use of this source code is governed by The MIT License
found in the accompanying LICENSE file.
