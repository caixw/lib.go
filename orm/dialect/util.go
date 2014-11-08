// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package dialect

// 生成limit n offset m语句，支持以下数据库：
// MySQL, H2, HSQLDB, Postgres, SQLite3
func mysqlLimit(limit, offset int) (string, []interface{}) {
	return " LIMIT ? OFFSET ? ", []interface{}{limit, offset}
}

// 生成limit n offset m语句，支持以下数据库：
// Derby, SQL Server 2012, Oracle 12c, the SQL:2008 standard
func oracleLimit(limit, offset int) (string, []interface{}) {
	return " OFFSET ? ROWS FETCH NEXT ? ROWS ONLY ", []interface{}{offset, limit}
}
