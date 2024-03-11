package template

var HandlerTemplate = `package http

import (
	"context"

	"github.com/gin-gonic/gin"
)

// 列表
func Make${Name}ListHandler(svc ${TableName}.I${Name}API) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &${TableName}.List${Name}Param{}
		err := BindQuery(c, req)
		if err != nil {
			transport.EncodeResponse(c, errors.NewError(_const.CodeDataError, "参数错误"), nil)
			return
		}
		ctx := context.Background()
		uuid, _ := c.GetQuery("uuid")
		if len(uuid) == 0 {
			transport.EncodeResponse(c, errors.NewError(_const.CodeDataError, "uuid不能为空"), nil)
			return
		}
		ctx = context.WithValue(ctx, "uuid", uuid)       //继续传递uid
		ctx = context.WithValue(ctx, "ip", c.ClientIP()) //继续传递ip
		req.Uuid = uuid
		resp, err := svc.List(ctx, req)
		transport.EncodeResponse(c, err, resp)
		c.Next()
	}
}

// 创建
func Make${Name}CreateHandler(svc ${TableName}.I${Name}API) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &${TableName}.Create${Name}Param{}
		err := BindJSON(c, &req)
		if err != nil {
			transport.EncodeResponse(c, errors.NewError(_const.CodeDataError, "参数错误"), nil)
			return
		}
		ctx := context.Background()
		uuid, _ := c.GetQuery("uuid")
		if len(uuid) == 0 {
			transport.EncodeResponse(c, errors.NewError(_const.CodeDataError, "uuid不能为空"), nil)
			return
		}
		ctx = context.WithValue(ctx, "uuid", uuid)       //继续传递uid
		ctx = context.WithValue(ctx, "ip", c.ClientIP()) //继续传递ip
		req.Uuid = uuid
		err = svc.Create(ctx, req)
		transport.EncodeResponse(c, err, nil)
		c.Next()
	}
}

// 修改
func Make${Name}UpdateHandler(svc ${TableName}.I${Name}API) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &${TableName}.Update${Name}Param{}
		err := BindJSON(c, &req)
		if err != nil {
			transport.EncodeResponse(c, errors.NewError(_const.CodeDataError, "参数错误"), nil)
			return
		}
		ctx := context.Background()
		uuid, _ := c.GetQuery("uuid")
		if len(uuid) == 0 {
			transport.EncodeResponse(c, errors.NewError(_const.CodeDataError, "uuid不能为空"), nil)
			return
		}
		ctx = context.WithValue(ctx, "uuid", uuid)       //继续传递uid
		ctx = context.WithValue(ctx, "ip", c.ClientIP()) //继续传递ip
		req.Uuid = uuid
		err = svc.Update(ctx, req)
		transport.EncodeResponse(c, err, nil)
		c.Next()
	}
}

// 删除
func Make${Name}DeleteHandler(svc ${TableName}.I${Name}API) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &${TableName}.Delete${Name}Param{}
		err := BindJSON(c, req)
		if err != nil {
			transport.EncodeResponse(c, errors.NewError(_const.CodeDataError, "参数错误"), nil)
			return
		}
		ctx := context.Background()
		uuid, _ := c.GetQuery("uuid")
		if len(uuid) == 0 {
			transport.EncodeResponse(c, errors.NewError(_const.CodeDataError, "uuid不能为空"), nil)
			return
		}
		ctx = context.WithValue(ctx, "uuid", uuid)       //继续传递uid
		ctx = context.WithValue(ctx, "ip", c.ClientIP()) //继续传递ip
		req.Uuid = uuid
		err = svc.Delete(ctx, req)
		transport.EncodeResponse(c, err, nil)
		c.Next()
	}
}
`
