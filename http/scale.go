package http

import "net/http"

func scaleHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
		req      struct {
			Cluster        string `form:"cluster" binding:"required"`
			Namespace      string `form:"namespace" binding:"required"`
			DeploymentName string `form:"deploymentName" binding:"required"`
			Replicas       int32  `form:"replicas" binding:"required"`
		}
	)
	defer httpErrorCommon(c, &httpCode, &err)
	if err = c.Bind(&req); err != nil {
		httpCode = http.StatusBadRequest
		return
	}
	scale, err := svc.Deployer.UpdateScale(req.Cluster, req.Namespace, req.DeploymentName, req.Replicas)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
		Data: scale,
	})
}

func getScaleHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
		req      struct {
			Cluster        string `form:"cluster" binding:"required"`
			Namespace      string `form:"namespace" binding:"required"`
			DeploymentName string `form:"deploymentName" binding:"required"`
		}
	)
	defer httpErrorCommon(c, &httpCode, &err)
	if err = c.Bind(&req); err != nil {
		httpCode = http.StatusBadRequest
		return
	}
	scale, err := svc.Deployer.GetScale(req.Cluster, req.Namespace, req.DeploymentName)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
		Data: scale,
	})
}
