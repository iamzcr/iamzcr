package admin

import (
	"iamzcr/models"
	svc "iamzcr/services/admin"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) ListAdmins(c *gin.Context) {
	admins, err := h.adminSvc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": admins})
}

func (h *AdminHandler) GetAdmin(c *gin.Context) {
	id := parseInt(c.Param("id"))
	admin, err := h.adminSvc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": admin})
}

func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	if err := h.adminSvc.Create(&admin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": admin})
}

func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	id := parseInt(c.Param("id"))
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	if err := h.adminSvc.Update(id, data); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

func (h *AdminHandler) ChangeAdminPassword(c *gin.Context) {
	id := parseInt(c.Param("id"))
	var input struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	if err := h.adminSvc.ChangeAdminPassword(id, input.NewPassword); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "密码修改成功"})
}

func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if err := h.adminSvc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

func (h *AdminHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt("user_id")
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	if err := h.adminSvc.ChangePassword(userID, input.OldPassword, input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "密码修改成功"})
}

type AdminGroupHandler struct {
	svc *svc.AdminGroupService
}

func NewAdminGroupHandler(s *svc.AdminGroupService) *AdminGroupHandler {
	return &AdminGroupHandler{svc: s}
}

func (h *AdminGroupHandler) List(c *gin.Context) {
	groups, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": groups})
}

func (h *AdminGroupHandler) Get(c *gin.Context) {
	id := parseInt(c.Param("id"))
	group, err := h.svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": group})
}

func (h *AdminGroupHandler) Create(c *gin.Context) {
	var group models.AdminGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	if err := h.svc.Create(&group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": group})
}

func (h *AdminGroupHandler) Update(c *gin.Context) {
	id := parseInt(c.Param("id"))
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	if err := h.svc.Update(id, data); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "记录不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

func (h *AdminGroupHandler) Delete(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
