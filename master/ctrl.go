package master

import (
	"dister/protos"
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	svc *TaskService
}

func (ctrl *Controller) addTest(ctx *gin.Context) {
	var req protos.TaskData

	err := bindBody(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	task, err := ctrl.svc.AddTask(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "committed",
		"id": task,
	})
}

func bindBody(ctx *gin.Context, obj interface{}) error {
	err := ctx.ShouldBindBodyWith(obj, binding.JSON)
	if err != nil {
		return err
	}
	return nil
}
