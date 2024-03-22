package {{.packageName}}

import (
	"context"
    "errors"
	{{.imports}}

    "github.com/golifes/go-sqlbuilder"
    "google.golang.org/grpc/status"
)

const (
	ExistCode           = 1000020 //已存在
	CreateCode          = 1000021 //添加失败
	FindOneCode         = 1000022 //查询详情失败
	FindCode            = 1000023 //查询列表失败
	UpdateCode          = 1000024 //更新失败
	BatchUpdateCode     = 1000025 //批量更新失败
	DeleteSoftCode      = 1000026 //删除失败(软)
	BatchDeleteSoftCode = 1000027 //批量删除失败(软)
	DeleteCode          = 1000028 //删除失败
	BatchDeleteCode     = 1000029 //批量删除失败
	DeleteCacheCode     = 1000030 //删除缓存失败
	LimitMaxCode        = 1000031 //操作超过最大限制数

)

type entity *model.AccessMethod

func _fineOne_(ctx context.Context, svc *svc.ServiceContext, id any, builder *sqlbuilder.SelectBuilder) (obj entity, err error) {
	conn := svc.Entity.ReadAccessMethod
	if id != nil {
		obj, err = conn.FindOne(ctx, id)
	} else {
		obj, err = conn.FindOneByArgs(ctx, *builder)
	}

	return obj, err
}

func findOne(ctx context.Context, svc *svc.ServiceContext, id any, builder *sqlbuilder.SelectBuilder) (obj entity, err error) {
	obj, err = _fineOne_(ctx, svc, id, builder)
	if err != nil {
		return nil, status.Errorf(FindOneCode, err.Error())
	}
	return obj, nil
}

func create(ctx context.Context, svc *svc.ServiceContext, builder *sqlbuilder.SelectBuilder) (obj entity, err error) {
	obj, err = _fineOne_(ctx, svc, nil, builder)
	if err == nil {
		return nil, status.Errorf(ExistCode, model.Exist.Error())
	} else if err != nil && !errors.Is(err, model.NotFound) {
		return nil, status.Errorf(CreateCode, err.Error())
	}
	return obj, nil
}

func updateOne(ctx context.Context, svc *svc.ServiceContext, id any) (obj entity, err error) {
	obj, err = findOne(ctx, svc, id, nil)
	if err != nil {
		return nil, err
	} else if obj.Status != constant.NormalStatus || obj.IsDeleted != constant.IsNotDeleted {
		return nil, status.Errorf(FindOneCode, model.NotAvailable.Error())
	}
	return obj, err
}

func updateStatus(ctx context.Context, svc *svc.ServiceContext, id any, f int64) (obj entity, err error) {
	obj, err = findOne(ctx, svc, id, nil)
	if err != nil {
		return nil, err
	} else if obj.Status == f || obj.IsDeleted == constant.IsDeleted {
		return nil, status.Errorf(UpdateCode, model.NotAvailable.Error())
	}
	return obj, nil
}

func deleteSoft(ctx context.Context, svc *svc.ServiceContext, id any) (obj entity, err error) {
	obj, err = findOne(ctx, svc, id, nil)
	if err != nil {
		return nil, err
	} else if obj.IsDeleted == constant.IsDeleted {
		return nil, status.Errorf(FindOneCode, model.NotModified.Error())
	}
	return obj, nil
}
