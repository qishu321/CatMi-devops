package deploy

import (
	"CatMi-devops/utils/tools"
	"fmt"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	Sshd         = &SshdController{}
	CommandLog   = &LogController{}
	Template     = &TemplateController{}
	Template_Log = &Template_Log_Controller{}
	Task         = &TaskController{}
	TaskEnv      = &TaskEnvController{}
	validate     = validator.New()
	trans        ut.Translator
)

func Run(c *gin.Context, req interface{}, fn func() (interface{}, interface{})) {
	var err error
	// bind struct
	err = c.Bind(req)
	if err != nil {
		tools.Err(c, tools.NewValidatorError(err), nil)
		return
	}
	// 校验
	err = validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			tools.Err(c, tools.NewValidatorError(fmt.Errorf(err.Translate(trans))), nil)
			return
		}
	}
	data, err1 := fn()
	if err1 != nil {
		tools.Err(c, tools.ReloadErr(err1), data)
		return
	}
	tools.Success(c, data)
}
