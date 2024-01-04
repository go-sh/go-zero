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

func genTypes(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	typeGroup := map[string][]spec.Type{}
	for _, obj := range api.Service.Groups {
		group := empty
		if v, ok := obj.Annotation.Properties["group"]; ok {
			group = v
		}
		var types []spec.Type
		typesM := map[string]struct{}{}
		for _, route := range obj.Routes {
			if _, ok := typesM[route.ResponseType.Name()]; !ok {
				types = append(types, route.ResponseType)
			}
			if _, ok := typesM[route.RequestType.Name()]; !ok {
				types = append(types, route.RequestType)
			}
			typesM[route.ResponseType.Name()] = struct{}{}
			typesM[route.RequestType.Name()] = struct{}{}
		}
		typeGroup[group] = types
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
