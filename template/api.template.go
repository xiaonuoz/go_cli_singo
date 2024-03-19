package template

var ApiTemplate = `package ${TableName}

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)


type (
	${LocalName}API  struct{}
	I${Name}API interface {
		List(ctx context.Context, req *${TableName}.List${Name}Param) (resp *${TableName}.List${Name}Resp, err error)
		Create(ctx context.Context, req *${TableName}.Create${Name}Param) (err error)
		Update(ctx context.Context, param *${TableName}.Update${Name}Param) (err error)
		Delete(ctx context.Context, args *${TableName}.Delete${Name}Param) (err error)
	}
)

func New${Name}Svc() I${Name}API {
	return &${LocalName}API{}
}

// @Tags 组名
// @Summary 接口描述
// @Description
// @Accept application/json
// @Param Authorization header string true "授权访问的token"
// @Param data query ${TableName}.List${Name}Param true "请求参数"
// @Success 200 {object} ${TableName}.List${Name}Resp
// @Router 路由地址 [get]
func (api *${LocalName}API) List(ctx context.Context, param *${TableName}.List${Name}Param) (resp *${TableName}.List${Name}Resp, err error) {
	result, total, err := ${TableName}.GetList(param)
	if err != nil {
		return nil, errs.NewError(errcode.Code${Name}ListErr, err.Error())
	}

	return &${TableName}.List${Name}Resp{
		Data: result,

		Total:    int(total),
		PageSize: param.PageSize,
		PageNum:  param.PageNum,
	}, nil
}

// @Tags 组名
// @Summary 接口描述
// @Description
// @Accept application/json
// @Param Authorization header string true "授权访问的token"
// @Param data body ${TableName}.Create${Name}Param true "请求参数"
// @Success 200
// @Router 路由地址 [post]
func (api *${LocalName}API) Create(ctx context.Context, param *${TableName}.Create${Name}Param) (err error) {
	err = api.check(param)
	if err != nil {
		return errs.NewError(errcode.Code${Name}CreateErr, err.Error())
	}

	data, err := ${TableName}.Create(param)

	if err != nil {
		return errs.NewError(errcode.Code${Name}CreateErr, err.Error())
	}

	return nil
}

// @Tags 组名
// @Summary 接口描述
// @Description
// @Accept application/json
// @Param Authorization header string true "授权访问的token"
// @Param data body ${TableName}.Update${Name}Param true "请求参数"
// @Success 200
// @Router 路由地址 [put]
func (api *${LocalName}API) Update(ctx context.Context, param *${TableName}.Update${Name}Param) (err error) {
	if param.Id == 0 {
		return errs.NewError(errcode.Code${Name}UpdateErr, "id不能为0")
	}

	err = ${TableName}.Update(param)
	if err != nil {
		return errs.NewError(errcode.Code${Name}UpdateErr, err.Error())
	}

	return nil
}

// @Tags 组名
// @Summary 接口描述
// @Description
// @Accept application/json
// @Param Authorization header string true "授权访问的token"
// @Param data body ${TableName}.Delete${Name}Param true "请求参数"
// @Success 200
// @Router 路由地址 [delete]
func (api *${LocalName}API) Delete(ctx context.Context, param *${TableName}.Delete${Name}Param) (err error) {
	if param.Id == 0 {
		return errs.NewError(errcode.Code${Name}DeleteErr, "id不能为0")
	}

	err = ${TableName}.Delete(param)

	if err != nil {
		return errs.NewError(errcode.Code${Name}DeleteErr, err.Error())
	}

	return nil
}

// 校验参数
func (api *${LocalName}API) check(param *${TableName}.Create${Name}Param) error {

	return nil
}
`
