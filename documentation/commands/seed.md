# Database Seeding Guide

This document explains how database seeding works in the **EduclyPro backend**, inspired by Laravel’s factory and seeder pattern.

---

## Overview

Seeding allows you to populate your database with fake or initial data for:

* Development
* Testing
* Demo environments

The system is composed of:

* **Factory** → Generates fake data
* **Seeder** → Inserts data into the database
* **Database Seeder** → Central place to run all seeders

---

## Structure

```bash
modules/
  students/
    factory.go   # Fake data generator
    seeder.go    # Inserts students into DB

shared/
  database/
    seeder.go    # Global seeder (like Laravel DatabaseSeeder)

cmd/
  seed/
    main.go      # CLI entry point
```

---

## Factory (Fake Data Generator)

File: `modules/students/factory.go`

The factory is responsible for generating fake student data using `gofakeit`.

### Example:

```go
func NewFakeStudent() Student {
	return Student{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email() + gofakeit.LetterN(3),
		Phone:     gofakeit.Phone(),
	}
}
```

### Generate multiple records:

```go
func NewFakeStudents(count int) []Student
```

---

## Seeder (Module Level)

File: `modules/students/seeder.go`

The seeder inserts generated data into the database.

### Example:

```go
func SeedStudents(db *gorm.DB, count int) error {
	students := NewFakeStudents(count)
	return db.Create(&students).Error
}
```

### Safe Seeding (Recommended)

To avoid duplicate inserts:

```go
var existing int64
db.Model(&Student{}).Count(&existing)

if existing > 0 {
	return nil
}
```

---

## Database Seeder (Global)

File: `shared/database/seeder.go`

This is the central place where all module seeders are executed.

### Example:

```go
func SeedAll(db *gorm.DB) {
	if err := students.SeedStudents(db, 20); err != nil {
		panic(err)
	}
}
```

---

## Running the Seeder

Command:

```bash
go run cmd/seed/main.go
```

### What happens:

1. Database connection is established
2. `SeedAll()` is executed
3. All module seeders run
4. Data is inserted into the database

---

## Default Behavior

* The number of records is defined inside `SeedAll`
* Example:

```go
students.SeedStudents(db, 20)
```

👉 This will create **20 students**

---

## Notes

* Emails are made unique to avoid database constraint errors
* Running the seeder multiple times may duplicate data unless protected
* Safe seeding is recommended for production-like environments

---

## Future Improvements

* `--fresh` flag (reset database before seeding)
* Per-module seeding (`students`, `teachers`, etc.)
* CLI arguments for dynamic counts
* Seeder integration in module generator

---

## Summary

| Concept        | Equivalent                |
| -------------- | ------------------------- |
| Factory        | Generate fake data        |
| Seeder         | Insert data               |
| DatabaseSeeder | Central execution         |
| CLI Command    | `go run cmd/seed/main.go` |

---
