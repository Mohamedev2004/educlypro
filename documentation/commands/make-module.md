---

## Generated Layers

---

### 1. Model (model.go)

Defines the database structure using GORM.
```go
type Course struct {
    gorm.Model
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}
```

**Features:**
* Includes `gorm.Model` (ID, CreatedAt, UpdatedAt, DeletedAt)
* Supports `uniqueIndex` via `:unique`
* Soft delete included automatically via `DeletedAt`

---

### 2. DTO (dto.go)

Separates API input/output from database models.

#### Create Request — all fields required + validation
```go
type CreateCourseRequest struct {
    Title       string  `json:"title"       binding:"required,min=3"`
    Description string  `json:"description" binding:"required"`
    Price       float64 `json:"price"       binding:"required,min=0"`
}
```

#### Update Request — all fields optional, no binding
```go
type UpdateCourseRequest struct {
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}
```

#### Response
```go
type CourseResponse struct {
    ID          uint    `json:"id"`
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}
```

---

### 3. Repository (repository.go)

Handles all database operations. Only layer that talks to GORM.
```go
type Repository interface {
    Create(entity *Course) error
    GetByID(id uint) (*Course, error)
    GetAll() ([]Course, error)
    Update(entity *Course) error
    Delete(id uint) error
    Restore(id uint) (*Course, error)
}
```

| Method | SQL equivalent |
| ---------- | ------------------ |
| `Create` | INSERT |
| `GetByID` | SELECT WHERE id=? |
| `GetAll` | SELECT all |
| `Update` | UPDATE (full save) |
| `Delete` | soft DELETE |
| `Restore` | UPDATE deleted_at = NULL |

**Soft delete:** GORM sets `deleted_at` instead of removing the row. All queries automatically exclude soft-deleted records.

**Restore:** Uses `Unscoped()` to find soft-deleted records and clears `deleted_at` to bring them back.

---

### 4. Service (service.go)

Contains business logic. Only talks to Repository.

**Responsibilities:**
* Convert DTO → Model
* Call repository methods
* Convert Model → DTO response

**Partial update — only non-zero fields are applied:**
```go
if req.Title != "" {
    model.Title = req.Title
}
if req.Price != 0 {
    model.Price = req.Price
}
```

---

### 5. Handler (handler.go)

Handles HTTP requests using Gin. Only talks to Service.

| Method | Endpoint | Handler | Status |
| ------ | --------- | ------- | ------ |
| POST | /courses | `Create` | 201 |
| GET | /courses/:id | `GetByID` | 200/404 |
| GET | /courses | `GetAll` | 200 |
| PUT | /courses/:id | `Update` | 200 |
| DELETE | /courses/:id | `Delete` | 200 |
| POST | /courses/:id/restore | `Restore` | 200 | `{ message: "restored", data: {...} }` |

Validation errors from binding return `400 Bad Request` automatically.

---

### 6. Routes (routes.go)

Registers all routes into a Gin router group.
```go
func RegisterCourseRoutes(rg *gin.RouterGroup, db *gorm.DB)
```

Wire it in `main.go`:
```go
api := r.Group("/api")
courses.RegisterCourseRoutes(api, database.DB)
students.RegisterStudentRoutes(api, database.DB)
teachers.RegisterTeacherRoutes(api, database.DB)
```

---

## Data Flow
```text
Client Request
     ↓
Handler  — validates JSON binding, returns 400 on error
     ↓
Service  — business logic, partial updates, mapping
     ↓
Repository — GORM queries only
     ↓
Database
     ↑
Response back to client
```

---

## Validation Reference

| Syntax | Generated binding tag | Use case |
| ------------------- | ----------------------------- | --------------------- |
| `Email:string:unique:email` | `binding:"required,email"` | Email format |
| `Phone:string::len=10` | `binding:"required,len=10"` | Exact length |
| `Password:string::min=6` | `binding:"required,min=6"` | Minimum length |
| `Age:int::min=1` | `binding:"required,min=1"` | Minimum value |
| `Price:float::min=0` | `binding:"required,min=0"` | Non-negative |
| `Name:string` | `binding:"required"` | Required only |

---

## Utility Functions

| Function | Purpose |
| --------------- | ------------------------------------------ |
| `plural()` | Converts singular module name to plural for folder and package |
| `parseFields()` | Parses field string into `[]Field` |
| `mapFieldType()` | Maps string type name to Go type |
| `zeroValue()` | Returns zero value for partial update check |
| `stringToUint()` | Converts route param string to uint |
| `createFile()` | Creates file only if it doesn't exist |

---

## Features Summary

* Full CRUD generation
* Clean architecture (Handler → Service → Repository)
* Always write singular in command — plural handled automatically
* DTO separation from models
* GORM integration with soft delete
* Gin HTTP handlers with proper status codes
* Consistent response format: `{ message, data }` on Create, Update, Delete, Restore
* Validation via gin binding tags (required, email, min, len)
* Safe file creation — never overwrites existing files

---

## Current Limitations

* No pagination or filtering
* No relationships (foreign keys, hasMany, belongsTo)
* No Swagger/OpenAPI doc generation
* No middleware per-route support

---

## Possible Improvements

* Add pagination (`page`, `limit` query params)
* Add filtering and search
* Add foreign key / relation support
* Generate Swagger docs automatically
* Add per-route middleware injection