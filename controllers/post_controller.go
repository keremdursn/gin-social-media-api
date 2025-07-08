package controllers

import (
	"net/http"
	"strconv"

	"gin-blog-api/database"
	"gin-blog-api/models"
	"gin-blog-api/utils"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	userID := uint(userIDVal.(float64))

	title := c.PostForm("title")
	content := c.PostForm("content")

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Başlık ve içerik gerekli"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dosya alınamadı"})
		return
	}

	files := form.File["images"] // input name="images"

	var imageURLs []string

	for _, file := range files {
		// Dosyayı geçici kaydet
		tempPath := "/tmp/" + file.Filename
		if err := c.SaveUploadedFile(file, tempPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya kaydedilemedi"})
			return
		}

		url, err := utils.UploadImage(tempPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary yüklemesi başarısız"})
			return
		}
		imageURLs = append(imageURLs, url)
	}

	post := models.Post{
		Title:     title,
		Content:   content,
		UserID:    userID,
		ImageURLs: imageURLs,
	}

	if err := database.DB.Db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gönderi oluşturulamadı"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func GetAllPosts(c *gin.Context) {
	var posts []models.Post
	search := c.Query("search")

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	dbQuery := database.DB.Db.Preload("Likes").Where("is_active = ?", true).Limit(limit).Offset(offset)

	if search != "" {
		dbQuery = dbQuery.Where("title ILIKE ?", "%"+search+"%")
	}

	if err := dbQuery.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gönderiler alınamadı"})
		return
	}

	var response []gin.H
	for _, post := range posts {
		response = append(response, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"user_id":    post.UserID,
			"like_count": len(post.Likes),
			"created_at": post.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetPostsByUser(c *gin.Context) {
	userID := c.Param("id")

	var posts []models.Post
	if err := database.DB.Db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı gönderileri alınamadı"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func GetPostByID(c *gin.Context) {
	postID := c.Param("id")

	var post models.Post
	err := database.DB.Db.Preload("Comments").Preload("Likes").
		First(&post, postID).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gönderi bulunamadı"})
		return
	}

	// Like sayısını ayrı döndürmek istersen:
	response := gin.H{
		"id":       post.ID,
		"title":    post.Title,
		"content":  post.Content,
		"user_id":  post.UserID,
		"likes":    len(post.Likes),
		"comments": post.Comments,
	}

	c.JSON(http.StatusOK, response)
}

func UpdatePost(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	userID := uint(userIDVal.(float64))

	postID := c.Param("id")

	var post models.Post
	if err := database.DB.Db.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gönderi bulunamadı"})
		return
	}

	// Kullanıcı sahibi değilse hata
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu gönderiyi güncelleyemezsiniz"})
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}

	// Güncelle
	post.Title = input.Title
	post.Content = input.Content

	if err := database.DB.Db.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	postIDParam := c.Param("id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz post ID"})
		return
	}

	var post models.Post
	if err := database.DB.Db.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post bulunamadı"})
		return
	}

	// Sadece sahibi silebilir
	userIDVal, _ := c.Get("userID")
	userID := uint(userIDVal.(float64))
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Yetkiniz yok"})
		return
	}

	post.IsActive = false
	if err := database.DB.Db.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Post arşivlenemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post başarıyla arşivlendi"})
}

func GetFeed(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	userID := uint(userIDVal.(float64))

	// Önce takip edilen kullanıcıları al
	var following []uint
	database.DB.Db.
		Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Pluck("following_id", &following)

	if len(following) == 0 {
		c.JSON(http.StatusOK, []string{}) // kimseyi takip etmiyorsa boş liste
		return
	}

	var posts []models.Post
	if err := database.DB.Db.
		Preload("Likes").
		Where("user_id IN ?", following).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Feed alınamadı"})
		return
	}

	var response []gin.H
	for _, post := range posts {
		response = append(response, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"user_id":    post.UserID,
			"like_count": len(post.Likes),
			"created_at": post.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetLikedPosts(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş gerekli"})
		return
	}
	userID := uint(userIDVal.(float64))

	var likes []models.Like
	if err := database.DB.Db.Where("user_id = ?", userID).Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Beğenilen gönderiler alınamadı"})
		return
	}

	// Like edilen post ID'lerini topla
	var postIDs []uint
	for _, like := range likes {
		postIDs = append(postIDs, like.PostID)
	}

	var posts []models.Post
	if err := database.DB.Db.Preload("Likes").Where("id IN ?", postIDs).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gönderiler alınamadı"})
		return
	}

	var response []gin.H
	for _, post := range posts {
		response = append(response, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"user_id":    post.UserID,
			"like_count": len(post.Likes),
			"created_at": post.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}
