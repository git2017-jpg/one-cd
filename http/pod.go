package http

import "net/http"

func podListHandler(c *MyContext) {
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
	podList, err := svc.Deployer.PodList(req.Cluster, req.Namespace, req.DeploymentName)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
		Data: podList,
	})
}

func podLogHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
		req      struct {
			Cluster      string `form:"cluster" binding:"required"`
			Namespace    string `form:"namespace" binding:"required"`
			PodName      string `form:"podName" binding:"required"`
			Container    string `form:"container" binding:"required"`
			SinceSeconds int64  `form:"sinceSeconds"`
			Previous     bool   `form:"previous"`
		}
	)
	defer httpErrorCommon(c, &httpCode, &err)
	if err = c.Bind(&req); err != nil {
		httpCode = http.StatusBadRequest
		return
	}
	log, err := svc.Deployer.PodLog(req.Cluster, req.Namespace, req.PodName, req.Container, req.SinceSeconds, req.Previous)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
		Data: log,
	})
}

func podEventsHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
		req      struct {
			Cluster   string `form:"cluster" binding:"required"`
			Namespace string `form:"namespace" binding:"required"`
			PodName   string `form:"podName" binding:"required"`
		}
	)
	defer httpErrorCommon(c, &httpCode, &err)
	if err = c.Bind(&req); err != nil {
		httpCode = http.StatusBadRequest
		return
	}
	events, err := svc.Deployer.PodEvents(req.Cluster, req.Namespace, req.PodName)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
		Data: events,
	})
}
