package auth

import (
	"reflect"
	"strings"
	"log"

	"github.com/credli/com"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	
	"github.com/credli/hcsg/base"
	"github.com/credli/hcsg/models"
)

func SignedInID(ctx *macaron.Context, sess session.Store) string {
	if !models.Connected {
		return ""
	}
	
	uid := sess.Get("uid")
	if uid == nil {
		return ""
	}
	if id, ok := uid.(string); ok {
		if _, err := models.GetUser(id); err != nil {
			if !models.IsErrUserNotExist(err) {
				log.Printf("GetUser: %v", err)
			}
			return ""
		}
		return id
	}
	return ""
}

func SignedInUser(ctx *macaron.Context, sess session.Store) (*models.User, bool) {
	if !models.Connected {
		return nil, false
	}
	
	uid := SignedInID(ctx, sess)
	if uid != "" {
		baHead := ctx.Req.Header.Get("Authorization")
		if len(baHead) > 0 {
			auths := strings.Fields(baHead)
			if len(auths) == 2 && auths[0] == "Basic" {
				uname, passwd, err := base.BasicAuthDecode(auths[1])
				
				u, err := models.UserSignIn(uname, passwd)
				if err != nil {
					if !models.IsErrUserNotExist(err) {
						log.Printf("UserSignIn: %v", err)
					}
					return nil, false
				}
				return u, true
			}
		}
		return nil, false
	}
	
	u, err := models.GetUser(uid)
	if err != nil {
		log.Printf("GetUser: %v", err)
		return nil, false
	}
	return u, false
}

func getRuleBody(field reflect.StructField, prefix string) string {
	for _, rule := range strings.Split(field.Tag.Get("binding"), ";") {
		if strings.HasPrefix(rule, prefix) {
			return rule[len(prefix) : len(rule)-1]
		}
	}
	return ""
}

func GetSize(field reflect.StructField) string {
	return getRuleBody(field, "Size(")
}

func GetMinSize(field reflect.StructField) string {
	return getRuleBody(field, "MinSize(")
}

func GetMaxSize(field reflect.StructField) string {
	return getRuleBody(field, "MaxSize(")
}

func GetInclude(field reflect.StructField) string {
	return getRuleBody(field, "Include(")
}

// FIXME: struct contains a struct
func validateStruct(obj interface{}) binding.Errors {

	return nil
}

type Form interface {
	binding.Validator
}

func init() {
	binding.SetNameMapper(com.ToSnakeCase)
}

// AssignForm assign form values back to the template data.
func AssignForm(form interface{}, data map[string]interface{}) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		} else if len(fieldName) == 0 {
			fieldName = com.ToSnakeCase(field.Name)
		}

		data[fieldName] = val.Field(i).Interface()
	}
}

func validate(errs binding.Errors, data map[string]interface{}, f Form, l macaron.Locale) binding.Errors {
	if errs.Len() == 0 {
		return errs
	}

	data["HasError"] = true
	AssignForm(f, data)

	typ := reflect.TypeOf(f)
	val := reflect.ValueOf(f)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		}

		if errs[0].FieldNames[0] == field.Name {
			data["Err_"+field.Name] = true

			trName := field.Tag.Get("locale")
			if len(trName) == 0 {
				trName = l.Tr("form." + field.Name)
			} else {
				trName = l.Tr(trName)
			}

			switch errs[0].Classification {
			case binding.ERR_REQUIRED:
				data["ErrorMsg"] = trName + l.Tr("form.require_error")
			case binding.ERR_ALPHA_DASH:
				data["ErrorMsg"] = trName + l.Tr("form.alpha_dash_error")
			case binding.ERR_ALPHA_DASH_DOT:
				data["ErrorMsg"] = trName + l.Tr("form.alpha_dash_dot_error")
			case binding.ERR_SIZE:
				data["ErrorMsg"] = trName + l.Tr("form.size_error", GetSize(field))
			case binding.ERR_MIN_SIZE:
				data["ErrorMsg"] = trName + l.Tr("form.min_size_error", GetMinSize(field))
			case binding.ERR_MAX_SIZE:
				data["ErrorMsg"] = trName + l.Tr("form.max_size_error", GetMaxSize(field))
			case binding.ERR_EMAIL:
				data["ErrorMsg"] = trName + l.Tr("form.email_error")
			case binding.ERR_URL:
				data["ErrorMsg"] = trName + l.Tr("form.url_error")
			case binding.ERR_INCLUDE:
				data["ErrorMsg"] = trName + l.Tr("form.include_error", GetInclude(field))
			default:
				data["ErrorMsg"] = l.Tr("form.unknown_error") + " " + errs[0].Classification
			}
			return errs
		}
	}
	return errs
}

type LoginForm struct {
	UserName string `binding:"Required,AlphaDashDot;MaxSize(254)"`
	Password string `binding:"Required,MaxSize(254)"`
	Remember bool
}

func (f *LoginForm) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}
