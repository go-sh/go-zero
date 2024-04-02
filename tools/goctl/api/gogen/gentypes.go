package gogen

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

const typesFile = "types"
const empty = ""

//go:embed types.tpl
var typesTemplate string

// BuildTypes gen types to string
func BuildTypes(types []spec.Type) (string, error) {
	var builder strings.Builder
	first := true
	for _, tp := range types {
		if first {
			first = false
		} else {
			builder.WriteString("\n\n")
		}
		if err := writeType(&builder, tp); err != nil {
			return "", apiutil.WrapErr(err, "Type "+tp.Name()+" generate error")
		}
	}

	return builder.String(), nil
}

func Foreach(obj spec.Type, keys map[string][]string) {
	dfs := obj.(spec.DefineStruct)
	for _, d := range dfs.Members {
		switch d.Type.(type) {
		case spec.ArrayType:
			info := d.Type.(spec.ArrayType)
			switch info.Value.(type) {
			case spec.DefineStruct:
				v := info.Value.(spec.DefineStruct)
				keys[obj.Name()] = append(keys[obj.Name()], v.RawName)
				Foreach(v, keys)
			}
		case spec.DefineStruct:
			v := d.Type.(spec.DefineStruct)
			keys[obj.Name()] = append(keys[obj.Name()], v.RawName)
			Foreach(v, keys)
		}
	}
}

func genTypes(dir string, cfg *config.Config, api *spec.ApiSpec) error {

	dup := map[string]struct{}{}
	dfGroup := make(map[string]string)
	for _, obj := range api.Service.Groups {
		group := empty
		if v, ok := obj.Annotation.Properties["group"]; ok {
			group = v
		}
		for _, route := range obj.Routes {
			if route.ResponseType != nil {
				respKey := route.ResponseType.Name() + ":" + group
				if _, ok := dup[respKey]; !ok {
					dfGroup[respKey] = group
					dup[respKey] = struct{}{}
				}
			}
			if route.RequestType != nil {
				reqKey := route.RequestType.Name() + ":" + group
				if _, ok := dup[reqKey]; !ok {
					dfGroup[reqKey] = group
					dup[reqKey] = struct{}{}
				}
			}
		}
	}
	for _, obj := range api.Types {
		keys := make(map[string][]string) //当前结构体的所有子结构体名称名称
		Foreach(obj, keys)
		for _, arrays := range keys {
			for _, arr := range arrays {
				key := obj.Name() + ":"
				for k, g := range dfGroup {
					if strings.Contains(k, key) {
						dfGroup[arr+":"+g] = g
					}
				}
			}
		}

	}
	typeGroup := make(map[string][]spec.Type)
	for _, obj := range api.Types {
		key := obj.Name() + ":"
		for k, v := range dfGroup {
			//if strings.Contains(k, key) {
			//	typeGroup[v] = append(typeGroup[v], obj)
			//}

			arr := strings.Split(k, ":")
			//if len(arr) >= 2 {
			//	fmt.Println("key", k, arr[0], arr[1], obj.Name())
			//}
			if arr[0]+":" == key {
				typeGroup[v] = append(typeGroup[v], obj)
			}
		}
	}

	for pkg, v := range typeGroup {
		val, err := BuildTypes(v)
		if err != nil {
			return err
		}
		typeFilename, err := format.FileNamingFormat(cfg.NamingFormat, typesFile)
		if err != nil {
			return err
		}

		typeFilename = typeFilename + ".go"
		filename := path.Join(dir, typesDir, typeFilename)
		subDir := typesDir
		if pkg != empty {
			filename = path.Join(dir, typesDir, pkg, typeFilename)
			subDir = path.Join(typesDir, pkg)
		}
		os.Remove(filename)
		err = genFile(fileGenConfig{
			dir:             dir,
			subdir:          subDir,
			filename:        typeFilename,
			templateName:    "typesTemplate",
			category:        category,
			templateFile:    typesTemplateFile,
			builtinTemplate: typesTemplate,
			data: map[string]any{
				"types":        val,
				"containsTime": false,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func writeType(writer io.Writer, tp spec.Type) error {
	structType, ok := tp.(spec.DefineStruct)
	if !ok {
		return fmt.Errorf("unspport struct type: %s", tp.Name())
	}

	fmt.Fprintf(writer, "type %s struct {\n", util.Title(tp.Name()))
	for _, member := range structType.Members {
		if member.IsInline {
			if _, err := fmt.Fprintf(writer, "%s\n", strings.Title(member.Type.Name())); err != nil {
				return err
			}

			continue
		}

		if err := writeProperty(writer, member.Name, member.Tag, member.GetComment(), member.Type, 1); err != nil {
			return err
		}
	}
	fmt.Fprintf(writer, "}")
	return nil
}
