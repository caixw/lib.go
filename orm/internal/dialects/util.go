// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package dialects

// mysq系列数据库分页语法的实现。支持以下数据库：
// MySQL, H2, HSQLDB, Postgres, SQLite3
func mysqlLimitSQL(limit, offset int) (string, []interface{}) {
	return " LIMIT ? OFFSET ? ", []interface{}{limit, offset}
}

// oracle系列数据库分页语法的实现。支持以下数据库：
// Derby, SQL Server 2012, Oracle 12c, the SQL 2008 standard
func oracleLimitSQL(limit, offset int) (string, []interface{}) {
	return " OFFSET ? ROWS FETCH NEXT ? ROWS ONLY ", []interface{}{offset, limit}
}
