package admin

import (
	"fmt"
	"iamzcr/models"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type AttachService struct{}

func NewAttachService() *AttachService {
	return &AttachService{}
}

func (s *AttachService) Upload(file *multipart.FileHeader, assetDir string) (*models.Attach, error) {
	uploadDir := assetDir
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", time.Now().Unix(), time.Now().UnixNano()%100000, ext)
	savePath := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("文件读取失败")
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return nil, fmt.Errorf("文件保存失败")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("文件写入失败")
	}

	link := "/asset/" + filename
	cdnLink := ""

	var cdn models.Website
	if err := models.DB.Where("`key` = ?", "cdn_url").First(&cdn).Error; err == nil && cdn.Value != "" {
		cdnLink = cdn.Value + "/asset/" + filename
	}

	attach := models.Attach{
		Name:       file.Filename,
		Link:       link,
		CdnLink:    cdnLink,
		Path:       savePath,
		Status:     1,
		Type:       1,
		CreateTime: int(time.Now().Unix()),
	}

	if err := models.DB.Create(&attach).Error; err != nil {
		return nil, err
	}

	return &attach, nil
}

func (s *AttachService) List(page, pageSize int) ([]models.Attach, int64, error) {
	var attaches []models.Attach
	var total int64

	models.DB.Model(&models.Attach{}).Count(&total)
	err := models.DB.Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&attaches).Error

	return attaches, total, err
}

func (s *AttachService) Get(id int) (*models.Attach, error) {
	var attach models.Attach
	if err := models.DB.First(&attach, id).Error; err != nil {
		return nil, err
	}
	return &attach, nil
}

func (s *AttachService) Create(attach *models.Attach) error {
	attach.CreateTime = int(time.Now().Unix())
	return models.DB.Create(attach).Error
}

func (s *AttachService) Update(id int, data map[string]interface{}) error {
	var attach models.Attach
	if err := models.DB.First(&attach, id).Error; err != nil {
		return err
	}
	delete(data, "id")
	delete(data, "create_time")
	data["update_time"] = int(time.Now().Unix())
	return models.DB.Model(&attach).Updates(data).Error
}

func (s *AttachService) Delete(id int) error {
	return models.DB.Delete(&models.Attach{}, id).Error
}
