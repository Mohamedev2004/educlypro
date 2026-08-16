package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/make-module/main.go <module-name> [fields]")
		fmt.Println("Example: go run cmd/make-module/main.go course \"Title:string::min=3,Description:string,Price:float::min=0\"")
		return
	}

	module := strings.ToLower(os.Args[1])
	fieldsInput := ""
	if len(os.Args) >= 3 {
		fieldsInput = os.Args[2]
	}

	fields := parseFields(fieldsInput)
	folderName := plural(module)
	basePath := filepath.Join("modules", folderName)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		panic(err)
	}

	createFile(basePath, "model.go", modelTemplate(module, fields))
	createFile(basePath, "dto.go", dtoTemplate(module, fields))
	createFile(basePath, "repository.go", repositoryTemplate(module, fields))
	createFile(basePath, "service.go", serviceTemplate(module, fields))
	createFile(basePath, "handler.go", handlerTemplate(module, fields))
	createFile(basePath, "routes.go", routesTemplate(module, folderName))

	fmt.Println("✅ Module created:", folderName+"/  →  type: "+strings.Title(module))
}

type Field struct {
	Name       string
	Type       string
	Unique     bool
	Validation string
}

func plural(module string) string {
	if strings.HasSuffix(module, "y") {
		return strings.TrimSuffix(module, "y") + "ies"
	}
	if strings.HasSuffix(module, "s") || strings.HasSuffix(module, "x") ||
		strings.HasSuffix(module, "ch") || strings.HasSuffix(module, "sh") {
		return module + "es"
	}
	return module + "s"
}

func parseFields(input string) []Field {
	if input == "" {
		return []Field{}
	}
	parts := strings.Split(input, ",")
	var fields []Field
	for _, p := range parts {
		ft := strings.Split(strings.TrimSpace(p), ":")
		if len(ft) < 2 {
			continue
		}

		unique := false
		validation := ""

		if len(ft) >= 3 {
			third := strings.ToLower(ft[2])
			if third == "unique" {
				unique = true
			} else if third == "email" {
				validation = "email"
			}
		}

		if len(ft) >= 4 && ft[3] != "" {
			if validation != "" {
				validation = validation + "," + ft[3]
			} else {
				validation = ft[3]
			}
		}

		fields = append(fields, Field{
			Name:       ft[0],
			Type:       ft[1],
			Unique:     unique,
			Validation: validation,
		})
	}
	return fields
}

func createFile(path, name, content string) {
	full := filepath.Join(path, name)
	if _, err := os.Stat(full); err == nil {
		fmt.Println("⚠️  skipped:", name, "(already exists)")
		return
	}
	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		panic(err)
	}
	fmt.Println("✔ created:", full)
}

func modelTemplate(module string, fields []Field) string {
	typeName := strings.Title(module)
	fieldStr := ""
	for _, f := range fields {
		if f.Unique {
			fieldStr += fmt.Sprintf("    %s %s `gorm:\"uniqueIndex\" json:\"%s\"`\n", f.Name, mapFieldType(f.Type), strings.ToLower(f.Name))
		} else {
			fieldStr += fmt.Sprintf("    %s %s `json:\"%s\"`\n", f.Name, mapFieldType(f.Type), strings.ToLower(f.Name))
		}
	}
	return fmt.Sprintf(`package %s

import "gorm.io/gorm"

type %s struct {
	gorm.Model
%s}
`, plural(module), typeName, fieldStr)
}

func dtoTemplate(module string, fields []Field) string {
	typeName := strings.Title(module)
	createStr := ""
	updateStr := ""
	responseStr := ""
	for _, f := range fields {
		binding := "required"
		if f.Validation != "" {
			binding = "required," + f.Validation
		}

		updateBinding := "omitempty"
		if f.Validation != "" {
			updateBinding = "omitempty," + f.Validation
		}

		jsonTag := strings.ToLower(f.Name)

		createStr += fmt.Sprintf(
			"    %s %s `json:\"%s\" binding:\"%s\"`\n",
			f.Name, mapFieldType(f.Type), jsonTag, binding,
		)

		updateStr += fmt.Sprintf(
			"    %s %s `json:\"%s\" binding:\"%s\"`\n",
			f.Name, mapFieldType(f.Type), jsonTag, updateBinding,
		)

		responseStr += fmt.Sprintf(
			"    %s %s `json:\"%s\"`\n",
			f.Name, mapFieldType(f.Type), jsonTag,
		)
	}

	return fmt.Sprintf(`package %s

type Create%sRequest struct {
%s}

type Update%sRequest struct {
%s}

type %sResponse struct {
	ID uint `+"`json:\"id\"`"+`
%s}
`, plural(module), typeName, createStr,
		typeName, updateStr,
		typeName, responseStr)
}

func repositoryTemplate(module string, fields []Field) string {
	typeName := strings.Title(module)
	pkg := plural(module)

	return fmt.Sprintf(`package %s

import "gorm.io/gorm"

type Repository interface {
	Create(entity *%s) error
	GetByID(id uint) (*%s, error)
	GetAll() ([]%s, error)
	Update(entity *%s) error
	Delete(id uint) error
	Restore(id uint) (*`+typeName+`, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(entity *%s) error {
	return r.db.Create(entity).Error
}

func (r *repository) GetByID(id uint) (*%s, error) {
	var entity %s
	if err := r.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *repository) GetAll() ([]%s, error) {
	var entities []%s
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository) Update(entity *%s) error {
	return r.db.Save(entity).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&%s{}, id).Error
}

func (r *repository) Restore(id uint) (*`+typeName+`, error) {
	var entity `+typeName+`
	if err := r.db.Unscoped().First(&entity, id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Unscoped().Model(&entity).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}
`, pkg, typeName, typeName, typeName, typeName, typeName, typeName, typeName, typeName, typeName, typeName, typeName)
}

func serviceTemplate(module string, fields []Field) string {
	typeName := strings.Title(module)
	pkg := plural(module)

	reqToModel := ""
	for _, f := range fields {
		reqToModel += fmt.Sprintf("        %s: req.%s,\n", f.Name, f.Name)
	}

	modelToDTO := ""
	for _, f := range fields {
		modelToDTO += fmt.Sprintf("        %s: model.%s,\n", f.Name, f.Name)
	}

	updateFields := ""
	for _, f := range fields {
		updateFields += fmt.Sprintf("    if req.%s != %s {\n        model.%s = req.%s\n    }\n",
			f.Name, zeroValue(f.Type), f.Name, f.Name)
	}

	return `package ` + pkg + `

type Service interface {
	Create(req *Create` + typeName + `Request) (*` + typeName + `Response, error)
	GetByID(id uint) (*` + typeName + `Response, error)
	GetAll() ([]` + typeName + `Response, error)
	Update(id uint, req *Update` + typeName + `Request) (*` + typeName + `Response, error)
	Delete(id uint) error
	Restore(id uint) (*` + typeName + `Response, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(req *Create` + typeName + `Request) (*` + typeName + `Response, error) {
	model := ` + typeName + `{
` + reqToModel + `	}
	if err := s.repo.Create(&model); err != nil {
		return nil, err
	}
	return &` + typeName + `Response{
		ID: model.ID,
` + modelToDTO + `	}, nil
}

func (s *service) GetByID(id uint) (*` + typeName + `Response, error) {
	model, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &` + typeName + `Response{
		ID: model.ID,
` + modelToDTO + `	}, nil
}

func (s *service) GetAll() ([]` + typeName + `Response, error) {
	models, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]` + typeName + `Response, len(models))
	for i, model := range models {
		result[i] = ` + typeName + `Response{
			ID: model.ID,
` + modelToDTO + `		}
	}
	return result, nil
}

func (s *service) Update(id uint, req *Update` + typeName + `Request) (*` + typeName + `Response, error) {
	model, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
` + updateFields + `
	if err := s.repo.Update(model); err != nil {
		return nil, err
	}
	return &` + typeName + `Response{
		ID: model.ID,
` + modelToDTO + `	}, nil
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *service) Restore(id uint) (*` + typeName + `Response, error) {
	model, err := s.repo.Restore(id)
	if err != nil {
		return nil, err
	}
	return &` + typeName + `Response{
		ID: model.ID,
` + modelToDTO + `	}, nil
}
`
}

func handlerTemplate(module string, fields []Field) string {
	typeName := strings.Title(module)
	entityName := strings.ToLower(module)
	pkg := plural(module)

	return `package ` + pkg + `

import (
	"fmt"
	"net/http"

	"educlypro/shared/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(s Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Create(c *gin.Context) {
	var req Create` + typeName + `Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed",
			"errors":  utils.FormatValidationErrors(err),
		})
		return
	}

	res, err := h.svc.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create ` + entityName + `",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "` + entityName + ` created successfully",
		"data":    res,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	res, err := h.svc.GetByID(stringToUint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "` + entityName + ` not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

func (h *Handler) GetAll(c *gin.Context) {
	res, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch ` + entityName + `s",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req Update` + typeName + `Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed",
			"errors":  utils.FormatValidationErrors(err),
		})
		return
	}

	res, err := h.svc.Update(stringToUint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update ` + entityName + `",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "` + entityName + ` updated successfully",
		"data":    res,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(stringToUint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete ` + entityName + `",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "` + entityName + ` deleted successfully",
	})
}

func (h *Handler) Restore(c *gin.Context) {
	id := c.Param("id")

	res, err := h.svc.Restore(stringToUint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to restore ` + entityName + `",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "` + entityName + ` restored successfully",
		"data":    res,
	})
}

func stringToUint(s string) uint {
	var i uint
	fmt.Sscan(s, &i)
	return i
}
`
}

func routesTemplate(module, folderName string) string {
	typeName := strings.Title(module)
	return `package ` + folderName + `

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register` + typeName + `Routes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	group := rg.Group("/` + folderName + `")
	{
		group.POST("", h.Create)
		group.GET("/:id", h.GetByID)
		group.GET("", h.GetAll)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
		group.POST("/:id/restore", h.Restore)
	}
}
`
}

func mapFieldType(t string) string {
	switch t {
	case "string":
		return "string"
	case "int":
		return "int"
	case "float":
		return "float64"
	case "bool":
		return "bool"
	default:
		return "string"
	}
}

func zeroValue(t string) string {
	switch t {
	case "string":
		return `""`
	case "int", "float":
		return "0"
	case "bool":
		return "false"
	default:
		return "nil"
	}
}

func stringToUint(s string) uint {
	var i uint
	fmt.Sscan(s, &i)
	return i
}
