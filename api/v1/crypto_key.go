package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/service/gorm"
	"kama_chat_server/pkg/zlog"
)

// GetPublicKeyBundle 获取用户的公钥束
// @Summary 获取公钥束
// @Description 获取指定用户的公钥束，用于建立加密会话
// @Tags Crypto
// @Accept json
// @Produce json
// @Param user_id query string true "用户ID"
// @Success 200 {object} respond.PublicKeyBundleRespond
// @Router /crypto/getPublicKeyBundle [get]
func GetPublicKeyBundle(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "缺少 user_id 参数",
		})
		return
	}

	bundle, err := gorm.CryptoKeyService.GetPublicKeyBundle(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取公钥束成功",
		"data":    bundle,
	})
}

// GetOneTimePreKeyCount 获取剩余一次性预密钥数量
// @Summary 获取剩余OTP密钥数量
// @Description 查询当前用户剩余的一次性预密钥数量
// @Tags Crypto
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /crypto/getOneTimePreKeyCount [get]
func GetOneTimePreKeyCount(c *gin.Context) {
	// 从 JWT 中间件获取用户 ID
	uuid, exists := c.Get("uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权",
		})
		return
	}

	count, err := gorm.CryptoKeyService.GetOneTimePreKeyCount(uuid.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "查询成功",
		"data": gin.H{
			"count": count,
		},
	})
}

// ReplenishOneTimePreKeys 补充一次性预密钥
// @Summary 补充OTP密钥
// @Description 批量上传新的一次性预密钥
// @Tags Crypto
// @Accept json
// @Produce json
// @Param keys body []request.OneTimePreKeyData true "密钥列表"
// @Success 200 {object} map[string]interface{}
// @Router /crypto/replenishOneTimePreKeys [post]
func ReplenishOneTimePreKeys(c *gin.Context) {
	uuid, exists := c.Get("uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权",
		})
		return
	}

	var req struct {
		Keys []request.OneTimePreKeyData `json:"keys" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error("参数绑定失败: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}

	if err := gorm.CryptoKeyService.ReplenishOneTimePreKeys(uuid.(string), req.Keys); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "补充密钥失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "补充密钥成功",
	})
}

// RotateSignedPreKey 轮换签名预密钥
// @Summary 轮换签名预密钥
// @Description 更新用户的签名预密钥
// @Tags Crypto
// @Accept json
// @Produce json
// @Param signed_pre_key body request.SignedPreKeyData true "新的签名预密钥"
// @Success 200 {object} map[string]interface{}
// @Router /crypto/rotateSignedPreKey [post]
func RotateSignedPreKey(c *gin.Context) {
	uuid, exists := c.Get("uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "未授权",
		})
		return
	}

	var req struct {
		SignedPreKey request.SignedPreKeyData `json:"signed_pre_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error("参数绑定失败: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
		})
		return
	}

	if err := gorm.CryptoKeyService.RotateSignedPreKey(uuid.(string), req.SignedPreKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "轮换密钥失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "轮换密钥成功",
	})
}

