package v1

import (
	"net/http"

	"kama_chat_server/internal/service/oss"
	"kama_chat_server/pkg/zlog"

	"github.com/gin-gonic/gin"
)

// GetUploadCredential 获取 OSS 上传凭证
// @Summary 获取上传凭证
// @Description 获取阿里云 OSS 直传凭证，用于客户端直接上传加密文件
// @Tags OSS
// @Produce json
// @Param file_type query string false "文件类型：image/file"
// @Success 200 {object} oss.UploadCredential
// @Router /oss/upload-credential [get]
func GetUploadCredential(c *gin.Context) {
	// 从 JWT 中间件获取用户 ID
	uuid, exists := c.Get("uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权",
		})
		return
	}

	fileType := c.DefaultQuery("file_type", "file")

	credential, err := oss.OSSService.GenerateUploadCredential(uuid.(string), fileType)
	if err != nil {
		zlog.Error("生成上传凭证失败: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "生成上传凭证失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取上传凭证成功",
		"data":    credential,
	})
}

// GetDownloadURLRequest 获取下载 URL 请求
type GetDownloadURLRequest struct {
	OSSKey string `json:"oss_key" binding:"required"`
}

// GetDownloadURL 获取 OSS 下载签名 URL
// @Summary 获取下载URL
// @Description 获取阿里云 OSS 签名下载 URL，用于下载加密文件
// @Tags OSS
// @Accept json
// @Produce json
// @Param data body GetDownloadURLRequest true "OSS 文件键"
// @Success 200 {object} map[string]interface{}
// @Router /oss/download-url [post]
func GetDownloadURL(c *gin.Context) {
	// 从 JWT 中间件获取用户 ID
	_, exists := c.Get("uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权",
		})
		return
	}

	var req GetDownloadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 验证 OSS Key 格式
	if !oss.OSSService.ValidateOSSKey(req.OSSKey) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无效的文件路径",
		})
		return
	}

	// TODO: 可以添加权限验证，检查用户是否有权访问该文件
	// 例如：检查用户是否是消息的发送者或接收者

	downloadURL, err := oss.OSSService.GenerateDownloadURL(req.OSSKey)
	if err != nil {
		zlog.Error("生成下载 URL 失败: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "生成下载 URL 失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取下载 URL 成功",
		"data": gin.H{
			"url": downloadURL,
		},
	})
}
