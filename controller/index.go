package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/angao/gin-xorm-admin/db"
	"github.com/angao/gin-xorm-admin/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// IndexController is handle home page request
type IndexController struct {
}

// Home is handle "/" request
func (IndexController) Home(c *gin.Context) {
	session := sessions.Default(c)
	var user *models.UserRole
	var err error

	userID, ok := session.Get("user_id").(int64)
	if ok {
		var userDao db.UserDao
		var menuDao db.MenuDao
		user, err = userDao.GetUserRole(userID)
		if err != nil {
			r.HTML(c.Writer, http.StatusInternalServerError, "index.html", gin.H{
				"username": user.User.Name,
			})
			return
		}
		roleIDs := strings.Split(user.User.RoleId, ",")
		menus, err := menuDao.GetMenuByRoleIds(roleIDs)
		if err != nil {
			r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		menus = buildTree(menus)
		r.HTML(c.Writer, http.StatusOK, "index.html", gin.H{
			"username": user.User.Name,
			"rolename": user.Role.Name,
			"menus":    menus,
		})
	}
}

// BlackBoard is handle "/blackboard"
func (IndexController) BlackBoard(c *gin.Context) {
	var noticeDao db.NoticeDao
	notices, err := noticeDao.List()
	if err != nil {
		r.JSON(c.Writer, http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	r.HTML(c.Writer, http.StatusOK, "blackboard.html", gin.H{
		"noticeList": notices,
	})
}

// buildTree 生成菜单树结构
func buildTree(menuNodes []models.Menu) []models.Menu {
	if len(menuNodes) == 0 {
		return menuNodes
	}
	var menus []models.Menu
	for _, menu := range menuNodes {
		for _, sub := range menuNodes {
			if sub.Pcode == fmt.Sprintf("%d", menu.Id) {
				menu.Children = append(menu.Children, sub)
			}
		}
		if menu.Pcode == "0" {
			menus = append(menus, menu)
		}
	}

	return menus
}
