package generator

import (
	_ "embed"
	"fmt"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"path/filepath"
)

//go:embed entity.tpl
var entityTemplate string

func (g *Generator) GenEntity(ctx DirContext, serviceName, packageName, originPackageName string) error {
	dir := ctx.GetEntity()
	childPkg, err := dir.GetChildPackage(packageName)
	serviceDir := filepath.Base(childPkg)
	fileName := filepath.Join(dir.Filename, serviceDir, fmt.Sprintf("%v", "entity.go"))
	text, err := pathx.LoadTemplate(category, entityTemplateFile, entityTemplate)
	if err != nil {
		return err
	}
	if err = util.With(entity).GoFmt(true).Parse(text).SaveTo(map[string]any{
		"packageName":       packageName,
		"originPackageName": originPackageName,
		"imports":           "",
	}, fileName, false); err != nil {
		return err
	}
	return nil
}
