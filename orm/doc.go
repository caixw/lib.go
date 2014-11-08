// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// e := orm.NewEngine("db")
//
// t:=e.Begin() // 事务
//
// e.Save(&User{}) // update or insert
// e.Find().Where().And().Fetch(&User)
//
// e.Update(&User{}) // 按主键更新。没有主键返回error,若是数组，则一个个删除。
// e.Delete(&User{}) // 按主键删除
// e.Insert(&User{}) // 按主键插入
// e.Select(&User{}) // 若是数组，则按每一条分开查找。不存在的删除。
//
// e.Create(&User{}, upgrade bool) // 创建表， 若存在，是否更新？
// e.Drop(&User{}/string) // 删除表
// e.Clear(&User{}/string) // 清空表
//
// e.Model(&User).Create()/Drop()/Clear()
//               .Delete().Where().And().Exec()
//               .Update().Where().And().Exec()
//               .Insert().Data().Exec()
//               .Find().Fetch(...)
//
// e.Builder().Create().AddIndex().AddColumn().Exec()
//            .Delete(table).Where().And().Exec()
//
// e.Stmts().Get("user_update_by_id").Exec(args...)
//          .Set("user_update_yb_id", stmt) // 添加或是更新
//          .Add("user_update_by_id", stmt) // 只添加，若存在则返回信息
//
// e.DB().Query(..)
//       .Exec(...)
//
package orm

type UserExample struct {
	Id int `orm:"-,index:id,unique,ai(0,2),pk(id)"`
}
