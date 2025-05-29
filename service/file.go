package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_logistics/common"
	"go_logistics/model/entity"
	"strconv"
)

const (
	MaxFileSize = 16 << 10
)

func UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	fileTypeStr := c.PostForm("fileType")
	if fileTypeStr == "" || err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	defer file.Close()
	fileTypeInt, err := strconv.Atoi(fileTypeStr)
	if err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	size := header.Size
	if size > MaxFileSize {
		common.ErrorResponse(c, common.ServerError("仅支持16KB以下的文件！"))
		return
	}

	fileName := header.Filename
	fileContentType := header.Header.Get("Content-Type")
	fileData := make([]byte, size)
	_, err = file.Read(fileData)
	if err != nil {
		common.ErrorResponse(c, common.ServerError("文件读取错误！"))
		return
	}

	err = entity.InsertFile(c.Request.Context(), &entity.File{
		FileName:    fileName,
		FileType:    entity.BusinessType(fileTypeInt),
		FileSize:    size,
		ContentType: fileContentType,
		FileData:    fileData,
	})
	if err != nil {
		common.ErrorResponse(c, common.ServerError("文件上传失败！"))
		return
	}

	common.SuccessResponse(c)
}

func DeleteFile(c *gin.Context) {
	fileId := c.Query("fileId")
	if fileId == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	err := entity.DeleteFile(c.Request.Context(), fileId)
	if err != nil {
		common.ErrorResponse(c, common.ServerError("文件删除失败！"))
		return
	}
	common.SuccessResponse(c)
}

func GetFileList(c *gin.Context) {
	var dto entity.FindFileListDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	files, err := entity.GetFileList(c.Request.Context(), dto)
	if err != nil {
		common.ErrorResponse(c, common.ServerError("文件列表获取失败！"))
		return
	}
	common.SuccessResponseWithData(c, files)
}

func DownloadFile(c *gin.Context) {
	fileId := c.Query("fileId")
	if fileId == "" {
		common.ErrorResponse(c, common.ParamError)
		return
	}
	file, err := entity.GetFileById(c.Request.Context(), fileId)
	if err != nil {
		common.ErrorResponse(c, common.ServerError("文件获取失败！"))
		return
	}

	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
	c.Header("Content-Length", strconv.Itoa(len(file.FileData)))

	c.Data(200, file.ContentType, file.FileData)
}
