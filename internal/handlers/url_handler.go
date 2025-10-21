package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Hapaa16/go-shortener/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func base62Encode(n uint) string {
	if n == 0 {
		return string(alphabet[0])
	}
	var buf [11]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = alphabet[n%62]
		n /= 62
	}
	return string(buf[i:])
}

type UrlHandler struct {
	DB *gorm.DB
}

func NewUrlHandler(db *gorm.DB) *UrlHandler {
	return &UrlHandler{DB: db}
}

type shortenReq struct {
	Url string `json:"url"`
}

type shortenResp struct {
	ShortURL string `json:"short_url"`
}

func (h *UrlHandler) AccessUrl(c *gin.Context) {
	var in models.Url
	url := c.Param("url")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := h.DB.WithContext(ctx).
		Where("shortened_url = ?", url).
		Take(&in).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}
	h.DB.Model(models.Url{}).
		Where("shortened_url", url).
		Update("access_count", gorm.Expr("access_count + ?", 1))
	c.Redirect(http.StatusFound, in.OriginalUrl)
}
func (h *UrlHandler) UpdateUrl(c *gin.Context) {
	var req shortenReq

	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Url) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	url := c.Param("url")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var in models.Url

	if err := h.DB.WithContext(ctx).Where("shortened_url = ?", url).First(&in).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Url not found"})
		return
	}

	in.ShortenedUrl = req.Url
	h.DB.Save(&in)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *UrlHandler) DeleteUrl(c *gin.Context) {
	url := c.Param("url")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var in models.Url

	if err := h.DB.WithContext(ctx).Where("shortened_url = ?", url).First(&in).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Url not found"})
		return
	}
	h.DB.Delete(&in)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *UrlHandler) ShortenUrl(c *gin.Context) {
	var req shortenReq
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Url) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var existing models.Url
	if err := h.DB.WithContext(ctx).Where("original_url = ?", req.Url).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, shortenResp{
			ShortURL: fmt.Sprintf("http://localhost:8080/%s", existing.ShortenedUrl),
		})
		return
	}

	newUrl := models.Url{OriginalUrl: req.Url}
	if err := h.DB.WithContext(ctx).Create(&newUrl).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if err2 := h.DB.WithContext(ctx).Where("original_url = ?", req.Url).First(&existing).Error; err2 == nil {
				if existing.ShortenedUrl == "" {
					existing.ShortenedUrl = base62Encode(existing.ID)
					_ = h.DB.WithContext(ctx).
						Model(&existing).
						Update("shortened_url", existing.ShortenedUrl).Error
				}
				c.JSON(http.StatusOK, shortenResp{
					ShortURL: fmt.Sprintf("http://localhost:8080/%s", existing.ShortenedUrl),
				})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to shorten URL"})
		return
	}

	newUrl.ShortenedUrl = base62Encode(newUrl.ID)
	if err := h.DB.WithContext(ctx).
		Model(&newUrl).
		Update("shortened_url", newUrl.ShortenedUrl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save shortened URL"})
		return
	}

	c.JSON(http.StatusCreated, shortenResp{
		ShortURL: fmt.Sprintf("http://localhost:8080/%s", newUrl.ShortenedUrl),
	})
}
func (h *UrlHandler) GetStats(c *gin.Context) {
	url := c.Param("url")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var in models.Url

	if err := h.DB.WithContext(ctx).Where("shortened_url = ?", url).First(&in).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Url not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "result": in})

}
