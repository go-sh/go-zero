package gogen

import (
	_ "embed"
	"fmt"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const defaultLogicPackage = "logic"

//go:embed handler.tpl
var handlerTemplate string

type handlerInfo struct {
	PkgName            string
	ImportPackages     string
	ImportHttpxPackage string
	HandlerName        string
	RequestType        string
	LogicName          string
	LogicType          string
	Call               string
	HasResp            bool
	HasRequest         bool
}

func genHandler(dir, rootPkg string, cfg *config.Config, group spec.Group, route spec.Route, withoutSuffix bool) error {
	handler := getHandlerName(route, withoutSuffix)
	handlerPath := getHandlerFolderPath(group, route)
	pkgName := handlerPath[strings.LastIndex(handlerPath, "/")+1:]
	logicName := defaultLogicPackage
	if handlerPath != handlerDir {
		handler = strings.Title(handler)
		logicName = pkgName
	}
	//return doGenToFile(dir, handler, cfg, group, route, handlerInfo{
	//	PkgName:        pkgName,
	//	ImportPackages: genHandlerImports(group, route, rootPkg, pkgName),
	//	HandlerName:    handler,
	//	RequestType:    util.Title(route.RequestTypeName()),
	//	LogicName:      logicName,
	//	LogicType:      strings.Title(getLogicName(route, withoutSuffix)),
	//	Call:           strings.Title(strings.TrimSuffix(handler, "Handler")),
	//	HasResp:        len(route.ResponseTypeName()) > 0,
	//	HasRequest:     len(route.RequestTypeName()) > 0,
	//})

	filename, err := format.FileNamingFormat(cfg.NamingFormat, handler)
	if err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          getHandlerFolderPath(group, route),
		filename:        filename + ".go",
		templateName:    "handlerTemplate",
		category:        category,
		templateFile:    handlerTemplateFile,
		builtinTemplate: handlerTemplate,
		data: map[string]any{
			"PkgName":        pkgName,
			"ImportPackages": genHandlerImports(group, route, rootPkg, pkgName),
			"HandlerName":    handler,
			"RequestType":    util.Title(route.RequestTypeName()),
			"LogicName":      logicName,
			"LogicType":      strings.Title(getLogicName(route, withoutSuffix)),
			"Call":           strings.Title(strings.TrimSuffix(handler, "Handler")),
			"HasResp":        len(route.ResponseTypeName()) > 0,
			"HasRequest":     len(route.RequestTypeName()) > 0,
			"HasDoc":         len(route.JoinedDoc()) > 0,
			"Doc":            getDoc(route.JoinedDoc()),
		},
	})
}

//func doGenToFile(dir, handler string, cfg *config.Config, group spec.Group,
//	route spec.Route, handleObj handlerInfo,
//) error {
//	filename, err := format.FileNamingFormat(cfg.NamingFormat, handler)
//	if err != nil {
//		return err
//	}
//
//	return genFile(fileGenConfig{
//		dir:             dir,
//		subdir:          getHandlerFolderPath(group, route),
//		filename:        filename + ".go",
//		templateName:    "handlerTemplate",
//		category:        category,
//		templateFile:    handlerTemplateFile,
//		builtinTemplate: handlerTemplate,
//		data: map[string]any{
//			"PkgName":        pkgName,
//			"ImportPackages": genHandlerImports(group, route, rootPkg),
//			"HandlerName":    handler,
//			"RequestType":    util.Title(route.RequestTypeName()),
//			"LogicName":      logicName,
//			"LogicType":      strings.Title(getLogicName(route)),
//			"Call":           strings.Title(strings.TrimSuffix(handler, "Handler")),
//			"HasResp":        len(route.ResponseTypeName()) > 0,
//			"HasRequest":     len(route.RequestTypeName()) > 0,
//			"HasDoc":         len(route.JoinedDoc()) > 0,
//			"Doc":            getDoc(route.JoinedDoc()),
//		},
//	})
//}

func genHandlers(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec, withoutSuffix bool) error {
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			if err := genHandler(dir, rootPkg, cfg, group, route, withoutSuffix); err != nil {
				return err
			}
		}
	}

	return nil
}

func genHandlerImports(group spec.Group, route spec.Route, parentPkg, pkgName string) string {
	imports := []string{
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, getLogicFolderPath(group, route))),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)),
	}
	if len(route.RequestTypeName()) > 0 {
		if pkgName != handler {
			imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir, pkgName)))
		} else {
			imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir)))

		}
	}

	return strings.Join(imports, "\n\t")
}

func getHandlerBaseName(route spec.Route) (string, error) {
	handler := route.Handler
	handler = strings.TrimSpace(handler)
	handler = strings.TrimSuffix(handler, "handler")
	handler = strings.TrimSuffix(handler, "Handler")

	return handler, nil
}

func getHandlerFolderPath(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return handlerDir
		}
	}

	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")

	return path.Join(handlerDir, folder)
}

func getHandlerName(route spec.Route, withoutSuffix bool) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}
	if !withoutSuffix {
		return handler + "Handler"
	}
	return handler
}

func getLogicName(route spec.Route, withoutSuffix bool) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}
	if !withoutSuffix {
		return handler + "Logic"
	}
	return handler
}
