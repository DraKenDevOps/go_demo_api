// ============================================================================
// GORM RAW QUERY - QUICK REFERENCE
// ============================================================================

/*
Two Main Methods:
=================

1. db.Raw(sql) - For SELECT queries (returns results)
2. db.Exec(sql) - For INSERT/UPDATE/DELETE (modifies data)
*/

// ============================================================================
// METHOD 1: RAW QUERY (SELECT)
// ============================================================================
```go
package main

import "gorm.io/gorm"

type User struct {
	ID       uint
	Username string
	Email    string
	Status   string
}

// Simple SELECT
func SelectExample(db *gorm.DB) {
	var users []User
	db.Raw("SELECT * FROM users").Scan(&users)
}

// SELECT with WHERE
func SelectWithWhere(db *gorm.DB) {
	var users []User
	db.Raw("SELECT * FROM users WHERE status = ?", "ACTIVE").Scan(&users)
}

// SELECT with multiple parameters
func SelectWithMultipleParams(db *gorm.DB) {
	var users []User
	db.Raw(
		"SELECT * FROM users WHERE status = ? AND created_at > ?",
		"ACTIVE",
		"2024-01-01",
	).Scan(&users)
}

// SELECT single value
func SelectSingleValue(db *gorm.DB) {
	var count int64
	db.Raw("SELECT COUNT(*) FROM users WHERE status = ?", "ACTIVE").Scan(&count)
	println(count)
}

// SELECT into struct
func SelectIntoStruct(db *gorm.DB) {
	var result struct {
		Username string
		Email    string
		Count    int
	}
	db.Raw(
		"SELECT username, email, COUNT(*) as count FROM users GROUP BY username, email LIMIT 1",
	).Scan(&result)
}

// ============================================================================
// METHOD 2: EXEC (INSERT/UPDATE/DELETE)
// ============================================================================

// RAW INSERT
func InsertExample(db *gorm.DB) {
	result := db.Exec(
		"INSERT INTO users (username, email, status, created_at) VALUES (?, ?, ?, NOW())",
		"john_doe",
		"john@example.com",
		"ACTIVE",
	)

	rowsAffected := result.RowsAffected
	err := result.Error
}

// RAW UPDATE
func UpdateExample(db *gorm.DB) {
	result := db.Exec(
		"UPDATE users SET email = ?, updated_at = NOW() WHERE user_id = ?",
		"newemail@example.com",
		1,
	)

	if result.RowsAffected == 0 {
		println("User not found")
	}
}

// RAW DELETE (Soft)
func SoftDeleteExample(db *gorm.DB) {
	result := db.Exec(
		"UPDATE users SET deleted_at = NOW() WHERE user_id = ?",
		1,
	)
}

// RAW DELETE (Hard)
func HardDeleteExample(db *gorm.DB) {
	result := db.Exec(
		"DELETE FROM users WHERE user_id = ?",
		1,
	)
}
```
// ============================================================================
// SIDE-BY-SIDE: RAW vs ORM
// ============================================================================

/*
Task: Get all active users

RAW QUERY:
==========
var users []User
db.Raw("SELECT * FROM users WHERE status = ?", "ACTIVE").Scan(&users)

ORM:
====
var users []User
db.Where("status = ?", "ACTIVE").Find(&users)


Task: Count users by status

RAW QUERY:
==========
var result struct {
    Status string
    Count  int
}
db.Raw("SELECT status, COUNT(*) as count FROM users GROUP BY status").Scan(&result)

ORM:
====
db.Model(&User{}).Select("status", "COUNT(*) as count").
   Group("status").Scan(&result)


Task: Update multiple fields

RAW QUERY:
==========
db.Exec("UPDATE users SET email = ?, status = ?, updated_at = NOW() WHERE user_id = ?",
        email, status, userID)

ORM:
====
db.Model(&User{}).Where("user_id = ?", userID).
   Updates(map[string]interface{}{"email": email, "status": status})


Task: Complex JOIN with aggregation

RAW QUERY:
==========
db.Raw(`
    SELECT u.user_id, u.username, COUNT(p.post_id) as post_count
    FROM users u
    LEFT JOIN posts p ON u.user_id = p.user_id
    WHERE u.status = ?
    GROUP BY u.user_id
`, "ACTIVE").Scan(&results)

ORM:
====
db.Preload("Posts").Where("status = ?", "ACTIVE").Find(&users)
(Then count posts manually)
*/

// ============================================================================
// PARAMETER BINDING METHODS
// ============================================================================

/*
Method 1: Positional Placeholders (?)
======================================
db.Raw("SELECT * FROM users WHERE id = ? AND status = ?", 1, "ACTIVE").Scan(&user)

Works with: MySQL, PostgreSQL, SQLite
✓ Simple
✓ Clean


Method 2: Named Parameters
===========================
import "database/sql"

db.Raw(
    "SELECT * FROM users WHERE id = @id AND status = @status",
    sql.Named("id", 1),
    sql.Named("status", "ACTIVE"),
).Scan(&user)

✓ Clearer for complex queries
✗ More verbose


Method 3: Struct Binding
========================
type UserFilter struct {
    ID     int    `gorm:"column:id"`
    Status string `gorm:"column:status"`
}

filter := UserFilter{ID: 1, Status: "ACTIVE"}
db.Raw("SELECT * FROM users WHERE id = @ID AND status = @Status", filter).Scan(&user)

✓ Type-safe
✓ Clean for filters
*/

// ============================================================================
// COMMON PATTERNS
// ============================================================================

/*
Pattern 1: SELECT with COUNT (Pagination)
===========================================
var users []User
var total int64

db.Raw("SELECT COUNT(*) FROM users WHERE status = ?", "ACTIVE").Scan(&total)
db.Raw(
    "SELECT * FROM users WHERE status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
    "ACTIVE",
    pageSize,
    offset,
).Scan(&users)


Pattern 2: SELECT with aggregation
===================================
var results []struct {
    Status string
    Count  int64
}

db.Raw("SELECT status, COUNT(*) as count FROM users GROUP BY status").Scan(&results)


Pattern 3: SELECT with JOIN
============================
var results []struct {
    Username string
    PostCount int
}

db.Raw(`
    SELECT u.username, COUNT(p.post_id) as post_count
    FROM users u
    LEFT JOIN posts p ON u.user_id = p.user_id
    GROUP BY u.user_id
`).Scan(&results)


Pattern 4: INSERT with RETURNING (PostgreSQL)
===============================================
var newUser User
db.Raw(
    "INSERT INTO users (username, email) VALUES (?, ?) RETURNING *",
    "john",
    "john@example.com",
).Scan(&newUser)


Pattern 5: BATCH INSERT
========================
db.Exec(
    "INSERT INTO users (username, email) VALUES (?, ?), (?, ?), (?, ?)",
    "user1", "user1@example.com",
    "user2", "user2@example.com",
    "user3", "user3@example.com",
)


Pattern 6: Transaction with Raw
================================
tx := db.BeginTx(ctx, nil)
tx.Exec("INSERT INTO users ...")
tx.Exec("UPDATE profiles ...")
if err := tx.Commit().Error; err != nil {
    log.Fatal(err)
}


Pattern 7: Check if row exists
================================
var exists bool
db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = ?)", userID).Scan(&exists)


Pattern 8: Get last insert ID (MySQL)
======================================
var id int64
db.Raw("SELECT LAST_INSERT_ID()").Scan(&id)
*/

// ============================================================================
// ERROR HANDLING
// ============================================================================

/*
Check db.Raw() Error:
=====================
var users []User
if err := db.Raw("SELECT * FROM users").Scan(&users).Error; err != nil {
    log.Println("Query error:", err)
}


Check db.Exec() Error:
======================
result := db.Exec("UPDATE users SET status = ? WHERE id = ?", "ACTIVE", 1)
if result.Error != nil {
    log.Println("Exec error:", result.Error)
}

// Also check rows affected
if result.RowsAffected == 0 {
    log.Println("No rows affected")
}


Full error handling:
====================
result := db.Exec("INSERT INTO users (username) VALUES (?)", username)
if result.Error != nil {
    if strings.Contains(result.Error.Error(), "Duplicate") {
        // Handle duplicate key error
    }
    log.Println("Insert failed:", result.Error)
    return
}

if result.RowsAffected != 1 {
    log.Println("Expected 1 row affected, got:", result.RowsAffected)
}
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================

/*
✓ DO:

1. Always use parameters for user input
   db.Raw("SELECT * FROM users WHERE id = ?", userID)

2. Check for errors
   if result.Error != nil { ... }

3. Use transactions for multiple operations
   tx := db.BeginTx(ctx, nil)

4. Validate parameters before executing
   if userID < 1 { return error }

5. Use prepared statements for repeated queries
   stmt := db.Session(&gorm.Session{PrepareStmt: true})

6. Log raw queries in development
   fmt.Println("Query:", db.Statement.SQL.String())


✗ DON'T:

1. String concatenation for user input
   ❌ db.Raw("SELECT * FROM users WHERE id = " + userID)

2. Ignore errors
   ❌ db.Raw("SELECT ...").Scan(&users)

3. Use Raw for simple queries
   ❌ db.Raw("SELECT * FROM users").Scan(&users)
   ✓ db.Find(&users)

4. Forget WHERE clause in UPDATE/DELETE
   ❌ db.Exec("DELETE FROM users")
   ✓ db.Exec("DELETE FROM users WHERE id = ?", id)

5. Use User input without validation
   ❌ db.Raw("ORDER BY " + sortField)
   ✓ Validate sortField against whitelist
*/

// ============================================================================
// QUICK REFERENCE TABLE
// ============================================================================

/*
Operation          | Method      | Example
===================+=============+============================================
SELECT             | Raw()       | db.Raw("SELECT * FROM users").Scan(&users)
SELECT with WHERE  | Raw()       | db.Raw("SELECT * WHERE id = ?", id).Scan()
SELECT count       | Raw()       | db.Raw("SELECT COUNT(*) FROM users").Scan()
INSERT             | Exec()      | db.Exec("INSERT INTO users (...) VALUES (?...)")
UPDATE             | Exec()      | db.Exec("UPDATE users SET ... WHERE ...")
DELETE (soft)      | Exec()      | db.Exec("UPDATE users SET deleted_at = ...")
DELETE (hard)      | Exec()      | db.Exec("DELETE FROM users WHERE ...")
Transaction        | BeginTx()   | tx := db.BeginTx(ctx, nil)
Rollback           | Rollback()  | tx.Rollback()
Commit             | Commit()    | tx.Commit()
*/

// ============================================================================
// COMPLETE WORKING EXAMPLE
// ============================================================================

/*
package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

// GET /users - List all active users
```go
func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []User
	
	if err := h.db.Raw(
		"SELECT * FROM users WHERE status = ? ORDER BY created_at DESC",
		"ACTIVE",
	).Scan(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{"data": users})
}
```
// POST /users - Create user
```go
func (h *UserHandler) CreateUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	result := h.db.Exec(
		"INSERT INTO users (username, email, status, created_at) VALUES (?, ?, ?, NOW())",
		input.Username,
		input.Email,
		"ACTIVE",
	)
	
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	
	c.JSON(201, gin.H{"status": "success"})
}
```
// PUT /users/:id - Update user
```go
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	result := h.db.Exec(
		"UPDATE users SET email = ?, updated_at = NOW() WHERE user_id = ?",
		input.Email,
		userID,
	)
	
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	
	c.JSON(200, gin.H{"status": "success"})
}
```
// DELETE /users/:id - Delete user
```go
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	
	result := h.db.Exec(
		"UPDATE users SET deleted_at = NOW() WHERE user_id = ?",
		userID,
	)
	
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	
	c.JSON(200, gin.H{"status": "success"})
}
```
*/