# DBx

Database DSL Based on GORM / 基于 GORM 封装的一套 DSL

Maintainers: @will.huang, @chein.huang  
测试覆盖率：87.3%  
状态：可用，但考虑重构（脱离 GORM）

## Table of Contents

- [1. Design](#1-design)
- [2. Initialize](#2-initialize)
- [3. `Conn()` (`chain()`)](#3-conn)
- [4. About Option & Result](#4-about-option--result)
- [5. Query](#5-query)
- [6. Update](#6-update)
- [7. Create And Destroy](#7-create-and-destroy)
- [8. Transaction](#8-transaction)
- [9. `AfterCommit` Callback](#9-aftercommit-callback)
- [10. Other](#10-other)

## 1. Design

`dbx` 基于 GORM 封装了一套 DSL，目的是：
1. 隐藏、化简数据操作细节
2. 统一 CRUD 书写风格，并通过返回一致结构，与 API response 进行串联
3. 承载一些简单约定（自动化）
4. 填补 GORM 设计缺陷
5. [Detailed Features](#11-features)

`dbx` 设计了两类函数：
1. 一类是以 `*Chain{*gorm.DB}` 为接收者的链式方法  
    `dbx` 包外调用链式方法，由 `dbx.Conn()` 开启调用链，例如 `dbx.Conn().Where().UpdateBy()`  
    `dbx` 包内由 `dbx.chain()`开启调用链
2. 一类是可以直接调用的函数（下文称 **CRUD Interface**）），例如 `dbx.Create()`，CRUD Interface 将自动处理约定（例如 `DefaultScope`），且返回 `dbx.Result{}`
3. CRUD Interface 及部分链方法接收 variadic 参数 `opts ...dbx.Option` 以改变其行为，在下面不同例子中将详解

**注：下文中令 `should` 为 `dbx` import 时的别名**

Simple Usages
```go
var firmwares []Firmware{}
result := dbx.Where(&firmwares, should.LIKE{"product_type": "search"}, should.With{Count: true})
if result.Err != nil {
	return result.Err
}

err := dbx.UpdateBy(&firmwares, dbx.H{"version": 2}).Err
if err != nil {
	return err
}

var modules []Modules{}
err = dbx.Where(&modules, 
                should.PLAIN{"users.id IN (?)", []int{1, 2}}, 
                should.With{Join: []string{"Order", "User"}}, 
                should.Be{RelatedWith: &firmwares}).Err
if err != nil {
	return err
}

return firmwares, result.Count, modules
```

### 2. Initialize

```go
// in your initializers
func InitDBx() {
	dbx.Client = YourClient // initializer.DatabaseClient
	dbx.DefaultPage = 1
	dbx.DefaultRows = 10
	// dbx.DefaultLogger
	// dbx.DefaultLogFormat
}
```

### 3. `Conn()`

[source code](dbx.go)

该函数处理以下 options:

- `dbx.Conn(dbx.Opt{Tx: "begin"})` => `*Chain{*gorm.DB.Begin()}`
- `dbx.Conn(dbx.Opt{SaveAssoc: true})`
- `dbx.Conn(dbx.Opt{Model: user})` // shortcut to the left: `dbx.Model(user)`
- `dbx.Conn(dbx.Opt{Set: dbx.H{"trail:who": user}})`
- `dbx.Conn(dbx.Opt{Debug/UnLog/Logger/LogFormat})`

### 4. About Option & Result

source code: [Opt](option.go), [Result](result.go)

1. 为语义化需要，Opt 有两个类型别名 `Be` & `With`：
    ```
    should.Be{RelatedWith: ...}
    should.With{Count: true}
    ```
2. Opt 有用作 reduce 的两个函数 `OptsPack` & `OptsPackGet` 以及一个方法 `Merge`
3. 部分函数会返回 `Result`:
    ```go
    type Result struct {
    	Data  interface{} // 即入参中的 model 对象（指针）
    	Err   error       // 如果存在任意的 GORM error 则不为空
    	Total interface{} // 开启 Opt.Count 即触发 DB Count，返回 int64，否则为空
    	Tx    *Chain      // 若开启 dbx 事务，则不为空
    }
    ```
4. `Result` 有一些数据处理方法：
    1. `result.NotFound()`
    2. `result.Uniq()`
    3. `result.(Get)Ids()`

### 5. Query

#### 5.1 Query Chain

一般来说直接调用 (5.3) Query Interface 即可，不过这里也将 Interface 涉及的链方法文档化如下：

1. `Where(condition interface{})` [src](query.go)
    ```go
    // Conditions
    type EQ = map[string]interface{}
    type Combine []interface{}
 
    // http://gorm.io/docs/query.html#Plain-SQL
    type PLAIN []interface{}
    // http://gorm.io/docs/query.html#Or
    type OR []interface{}
    // http://gorm.io/docs/query.html#Not
    type NOT []interface{}
    type IN map[string]interface{}
    type LIKE map[string]interface{}
 
    // Usage
    dbx.Conn().Where("name = 'Williams'")
    dbx.Conn().Where(should.EQ{"id": 3})
    dbx.Conn().Where(should.PLAIN{"name = ? AND age = ?", "a", 10)
    dbx.Conn().Where(should.PLAIN{"name = ? AND ", []interface{}{"a", 10}})
    dbx.Conn().Where(should.NOT{"status = ?", "done"})
    dbx.Conn().Where(should.NOT{dbx.H{"status": "done"})
    dbx.Conn().Where(should.NOT{"status", []string{"done", "draft"}}) // NOT IN
 
    // Combine is a slice of dbx conditions
    dbx.Conn().Where(dbx.Combine{should.EQ{"id": 1}, should.OR{"name = ?", "Tony"}})
    dbx.Conn().Where(dbx.Combine{should.LIKE{"name": "Foo"}, should.IN{"id": []uint{1, 2, 3}}})
    ```
2. `Order(value interface{}, reorder = false)` [src](query.go)
    ```go
    dbx.Conn().Order("created_at DESC")
    // reorder 即忽略左边所有 Order
    dbx.Conn().Order("created_at DESC", true)
    ```
3. `Preload(value interface{}, withoutDefault = false)` [src](preload.go)
    ```go
    dbx.Conn().Preload("User")
    dbx.Conn().Preload("User.Orders", true)
    dbx.Conn().Preload([]string{"User", "User.Orders"})
    dbx.Conn().Preload(map[string][]interface{}{"User": []interface{}{
 	    func(db *gorm.DB) *gorm.DB {return db.Order("")}
    }})
    dbx.Conn().Preload(map[string][]interface{}{"User": []interface{}{"status = ?", "activated"}})
    ```
    说明：
    1. 上面 value 参数的三种类型即 Opt.Preload 选项支持的类型
    2. 无论是链式调用还是 CRUD Interface，preload 的对象（例如上例的 User）都会受配置的 DefaultScope 影响
    3. 传递 `map[string][]interface{}` 指定 preload 对象的 condition 会直接覆盖 DefaultScope
4. `Uniq(opt Opt)`：构建 `DISTINCT ON` 查询（仅支持 Pg）[src](query.go)
    ```go
    dbx.Conn().Uniq(dbx.Opt{UniqBy: "user_id"})
    dbx.Conn().Uniq(dbx.Opt{UniqBy: "user_id", UniqOrder: "created_at DESC"})
    ```
5. `Pagy(opt ...Opt)`：分页 [src](query.go)
    ```go
    dbx.Conn().Pagy() // 按照 Default Page & Rows(per page) 进行分页
    dbx.Conn().Pagy(dbx.Opt{Page: 2, Rows: 100})

    dbx.Conn().Unpagy() // 移除分页设置
    ```
6. `Unscoped()`： 移除包括 gorm 默认的 `deleted_at IS NULL` 在内的所有【左侧 scopes】[src](scope.go)
    ```go
    dbx.Conn().Unscoped()
    ```
7. `Scoping(obj interface{}, opts ...Opt)`：根据 Opt 对对象进行 unscope 或者附加 DefaultScope [src](scope.go)
    ```go
    // 等同于 Unscoped()
    dbx.Conn().Scoping(&user, dbx.Opt{Unscoped: true})
    // 等同于 Unscoped()
    dbx.Conn().Scoping(&user, dbx.Opt{WithDeleted: true, UnscopeDefault: true})
    // 先 Unscoped()，再执行 `model.Definitions` 中声明的 `User` 的 DefaultScope
    dbx.Conn().Scoping(&user, dbx.Opt{WithDeleted: true})
    
    // obj 也可以是一个 string (Definitions 的 key)
    dbx.Conn().Scoping("User")
    ```
8. `Joins(assocFieldName string, nestedAssocFieldNames ...string)` [src](query_preload.go)
    ```go
    dbx.Model(&user).Joins("Orders")
    dbx.ScopingModel(&user).Joins("Orders", "Products")
    dbx.ScopingModel(&user).Joins("Orders", "Products").Where(dbx.PLAIN{`products.name = 'abc'`})
    ```
    说明：
    1. 该方法通过（嵌套）读取 GORM 关联自动生成相应的 `LEFT JOIN` 子语句
    2. 若要 Join 的关联对象设置了 DefaultScope 中有关 `Join` 的部分，那么 `Joins` 方法会自动附加有关 SQL
    3. 与 Join 有关的 Where 目前只能通过上述 `PLAIN{}` 语句实现
    4. 其生成的 SQL 语句类似：
        ```sql
       SELECT "firmwares".* FROM "firmwares" /*
       */ LEFT JOIN "firm_permissions" ON "firm_permissions"."resource_id" = "firmwares"."id" AND "firm_permissions"."resource_type" = 'firmwares' /*
       */ LEFT JOIN "user_groups_and_firm_permissions" ON "user_groups_and_firm_permissions"."firm_permission_id" = "firm_permissions"."id" /*
       */ LEFT JOIN "user_groups" ON "user_groups_and_firm_permissions"."user_group_id" = "user_groups"."id" /*
       */ LEFT JOIN "users_and_user_groups" ON "users_and_user_groups"."user_group_id" = "user_groups"."id" /*
       */ LEFT JOIN "users" ON "users_and_user_groups"."user_id" = "users"."id" /*
       */ WHERE "firmwares"."deleted_at" IS NULL AND (("users"."id" = '1'))
        ```
  
### 5.2 Query Method

Chain 作为接收者，直接返回结果的方法：

1. `Count() (int64, error)` [src](query.go)
    ```go
    dbx.Model(&user).Where(should.IN{"id": []int64{1, 2, 3}}).Count()
    ```
2. `Pluck(fields interface{}, out interface{}) error` [src](query.go)
    ```go
    names := []string{}
    err := dbx.Model(&user).Pluck("name", &names)   
    err := dbx.Model(&user).Pluck([]string{"name", "age"}, &result) // Unsupported now
    ```
3. `Related(assoc interface{}, opts ...Opt) Result` [src](query_relation.go)
    ```go
    // will find out the users which associated with the user_groups
    dbx.Model(&user_groups).Related(&users)
    dbx.Model(&user_groups).Related(&users, dbx.Opt{Count: true})
    // AssocFieldName: 指定 User 模型在 UserGroup 声明中的字段名
    dbx.Model(&user_groups).Related(&users, dbx.Opt{AssocFieldName: "MyUsers"})
    ```
4. `FindOut(out ...interface{}) Result`: 等同于 GORM 的 `Find` [src](query.go)
    ```go
    dbx.Conn().Where().FindOut(&user)
    dbx.Model(&user).Where().FindOut()
    ```
5. `FindOrInitBy(condition interface{}) error` [src](query.go)
    ```go
    dbx.Model(&user).FindOrInitBy(should.LIKE{})
    ```
6. `RelatedWith(obj interface{}, opts ...Opt) Result` [src](query_relation.go)
    ```go
    // will find out the users which associated with the user_groups
    dbx.Model(&users).Related(&user_groups)
    ```

#### 5.3 Query Interface

（下列函数均会默认自动附加 DefaultScope）

1. `Where(out interface{}, condition interface{}, opts ...Opt) Result`: 查询多个 records [src](query.go)
    ```go
    // Basic Usage
    dbx.Where(&users, "name = 'Williams'")
    dbx.Where(&users, should.PLAIN{"name = ?", "a"})
    dbx.Where(&users, should.IN{"id": []uint{1, 2, 3}})
    dbx.Where(&users, should.EQ{"id": 3})
    dbx.Where(&users, should.LIKE{"name": "Will"})
    dbx.Where(&users, dbx.Combine{should.LIKE{"name": "T"}, should.IN{"id": []uint{1, 2, 3}}})
    
    // Related Usage
    dbx.Where(&users, should.LIKE{"account": "foo"}, should.Be{RelatedWith: &user_group})
    
    // With Options
    dbx.Where(&users, should.LIKE{"name": "Tony"}, should.With{Preload: "Orders"})
    dbx.Where(&users, should.LIKE{"name": "Tony"}, should.With{Order: "created_at DESC"})
    dbx.Where(&users, should.LIKE{"name": "Tony"}, should.With{Uniq: "group", UniqOrder: "created_at ASC"})
    dbx.Where(&users, should.LIKE{"name": "Tony"}, should.With{UnscopeDefault: true})
    
    // With Count
    dbx.Where(&users, should.LIKE{"name": "Tony"}, should.With{Count: true})
    dbx.Where(&users, should.LIKE{"name": "Tony"}, should.With{Page: 2, Count: true})
    dbx.Where(&users, should.LIKE{"account": "foo"}, should.Be{RelatedWith: &user_group}, should.With{Count: true})
    ```
2. `Find(out interface{}, condition interface{}, opts ...Opt) Result`: 查询一个 record [src](query.go)
    (alias as `FindBy`)
    ```go
    dbx.Find(&user, should.EQ{"id": 1})
    dbx.Find(&user, should.EQ{"name": "foo"}, should.With{Order: "created_at DESC"})
    ```
3. `First(out interface{}, n ...int) Result`: 简单取出前 n（默认 1）个 record [src](query.go)
    ```go
    dbx.First(&user)
    dbx.First(&users, 3)
    ```
4. `FindById(out interface{}, id interface{}, opts ...Opt) Result` [src](query.go)
    ```go
    dbx.FindById(&user, 1)
    ```
5. `Related(assoc interface{}, obj interface{}, opts ...Opt) Result` [src](query_relation.go)
    ```go
    dbx.Related(&users, &user_group)
    dbx.Related(&users, &user_group, should.With{Preload: "Orders"})
    ```
6. `FindOrInitBy(out, condition interface{}, opts ...Opt) Result` [src](query.go)
    ```go
    dbx.FindOrInitBy(&user, dbx.LIKE{})
    ```
7. `IsExists(value interface{}, condition interface{}, opts ...Opt) bool` [src](query.go)
    ```go
    dbx.IsExists(User{}, should.EQ{"name": "abc"})
    // value 可以为 string，但此时无法 attach DefaultScope
    dbx.IsExists("User", should.EQ{"deleted_at": nil, "name": "abc"})
    ```

### 6. Update

1. dbx 的更新方法既支持 struct 实例又支持 struct 切片（开发中）
2. 请注意更新前会自动校验 DefaultScope 中设置的 Uniqueness（除非 `dbx.Opt{SkipUniqValidate: true}`）
    // 即在操作执行前，自动查询并判断更新后否与数据库现有记录重复
3. 默认【不】在 update 时保存关联对象，除非设置了 `dbx.Opt{SaveAssoc: true}`
4. GORM: 使用`Update()`传递 struct，只会更新非空白字段

#### 6.1 Update By Chain

1. `Update(obj interface{}, opts ...Opt) Result` [src](update.go)
    ```go
    dbx.Conn().Where(should.EQ{}).Update(User{Name: "hello", Age: 18})
    ```
2. `UpdateBy(obj interface{}, values interface{}, opts ...Opt) Result` [src](update.go)
    ```go
    dbx.Conn().Where(should.EQ{}).UpdateBy(&user, dbx.H{"name": "hello", "age": 18})
    // values 也可以是 struct，会自动 marshal
    dbx.Conn().Where(should.EQ{}).UpdateBy(&user, Params{Name: "hello"})
    ```
3. `Decrement / Increment(obj interface{}, fields interface{}, opts ...Opt) Result` [src](update_extend.go)
    ```go
    // 更新 id 为 1 的固件 download_count 字段 +1
    dbx.Conn().Where(should.EQ{"id": 1}).Increment(&Firmware{}, "download_count", dbx.Opt{SkipCallback: true})
    // 更新 id 为 1 的固件 download_count 字段 +1, count 字段 +1
    dbx.Conn().Where(should.EQ{"id": 1}).Increment(&Firmware{}, []string{"download_count", "count"})
    // 更新 id 为 1 的固件 download_count 字段 +3, count 字段 +1
    dbx.Conn().Where(should.EQ{"id": 1}).Increment(&Firmware{}, map[string]int{"download_count": 3, "count": 1})
    ```

#### 6.2 Update Interface

1. `Update(obj interface{}, opts ...Opt) Result` [src](update.go)
    ```go
    dbx.Update(User{ID: 1, Name: "Test"})
    ```
2. `UpdateBy(obj interface{}, values interface{}, opts ...Opt) Result` [src](update.go)
    ```go
    dbx.UpdateBy(&user, dbx.H{"name": "Test"})
    ```
3. `Decrement/ Increment(obj interface{}, fields interface{}, opts ...Opt) Result` [src](update_extend.go)
    ```go
    dbx.Increment(&firmware, "download_count")
    ```

### 7. Create And Destroy

1. dbx 的增删方法既支持 struct 实例又支持 struct 切片（开发中）
2. 请注意 create 前会自动校验 DefaultScope 中设置的 Uniqueness（除非 `dbx.Opt{SkipUniqValidate: true}`）
    // 即在操作执行前，自动查询并判断更新后否与数据库现有记录重复
3. 默认【不】在 create 时保存关联对象，除非设置了 `dbx.Opt{SaveAssoc: true}`

#### 7.1 Create

```go
// Create(obj interface{}, opts ...Opt) Result
dbx.Create(&user)
// FirstOrCreate(obj interface{}, conditionObj interface{}, opts ...Opt) Result
dbx.FirstOrCreate(&user, User{Name: "Test"})
```

#### 7.2 Destroy

```go
// (c *Chain) Destroy(opts ...Opt) Result
dbx.Model(&user).Where(should.EQ{"id": 1}).Destroy()
// Destroy(obj interface{}, opts ...Opt) Result
dbx.Destroy(&user)
```

### 8. Transaction

[source code](transaction.go)

dbx 事务在 Query Interface 以及 CUD 所有函数中处理，即如果发生 error，自动 Rollback

有两种写法：  
第一种是闭包，无需 begin 以及 commit：
```go
result := dbx.Transaction(func(chain *dbx.Chain) error {
	err := dbx.Create(&user, should.With{Tx: chain}).Err
	if err != nil {
		return err
	}
	err = doSomething()
	if err != nil {
		return err
	}
	return nil
})
```
第二种写法需要手动 begin 以及 commit：
```go
created := dbx.Create(&user, should.With{Tx: "begin"})
if Fails(created) {
	return
}
updated := dbx.UpdateBy(&user, dbx.H{"activated": true}, should.With{Tx: created.Tx})
if Fails(updated) {
	return
}
result := updated.Commit() // or .Rollback()

// 或者在最后一个操作的 Opt 中指定 `TxCommit`
updated := dbx.UpdateBy(&user, dbx.H{"activated": true}, should.With{Tx: created.Tx, TxCommit: true})
```

### 9. `AfterCommit` Callback

TODO

### 10. Other

#### 10.1 Model

Model(obj) 用以指定后续操作的对象，可以视作 `new Model`. [source code](model.go)

```go
dbx.Conn().Model(&user).UpdateBy(dbx.H{Name: "Test"})
dbx.Conn().ScopingModel(&user).UpdateBy(dbx.H{Name: "Test"})

// 以下方法可以代替 Conn()
dbx.Model(&user).Where()
dbx.ScopingModel(&user).Where()
```

#### 10.2 Other Interface

1. `IdOf(obj interface{}) interface{}` [src](dbx.go)
    ```go
    dbx.IdOf(&user).(int64)
    ```
2. `IsNewRecord(obj interface{}) bool` [src](dbx.go)
    ```go
    dbx.IsNewRecord(User{})      // => true
    dbx.IsNewRecord(User{ID: 1}) // => false
    ```

#### 10.3 Logging

[source code](logging.go)

可以通过设置 Opt 有关 Log 的选项（或者配置 initializer 默认设置）改变 GORM 以及 dbx 的日志行为，包括：

1. Logger: 需要实现 [dbx.logger](logging.go)
    ```go
    // 默认
    gorm.Logger{log.New(os.Stdout, "\r\n", 0)}
    // 文件
    f, _ := os.Create("log/dbx.log")
    logger := log.New(f, "\r\n", 0)
    ```
2. LogFormat: 内置的实现有 `"json"`

### 11. Features

新增：
1. Joins
2. Increment / Decrement
3. Transaction Block & Auto-Rollback & Log
4. Uniq
5. Batch CUD Supporting in the same Interface (on development)
6. `AfterCommit`
7. Configurable Serialization
8. Helpers: IdOf...

增强：
1. Configurable DefaultScope: CRUD & Preload & Join
2. Uniqueness Auto-Validation
3. Configurable Logging
4. Condition DSL
