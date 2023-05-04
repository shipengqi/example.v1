package validator

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/stretchr/testify/assert"
)

const (
	minDescriptionLength = 5
	minNameLength        = 2
)

// 约束
//
// 范围约束
// 对于数值，则约束其值
// 对于字符串，则约束其长度
// 对于切片、数组和 map，则约束其长度
//
// 范围约束 tags
// len: 等于参数值，例如 len=10；
// max: 小于等于参数值，例如 max=10；
// min: 大于等于参数值，例如 min=10；
// eq: 等于参数值，注意与 len 不一样。对于字符串，eq 约束字符串自己的值，而 len 约束字符串长度。例如 eq=10
// ne: 不等于参数值，例如 ne=10；
// gt: 大于参数值，例如 gt=10；
// gte: 大于等于参数值，例如 gte=10；
// lt: 小于参数值，例如 lt=10；
// lte: 小于等于参数值，例如 lte=10；
// oneof: 只能是列举出的值其中一个，这些值必须是数值或字符串，以空格分隔，若是字符串中有空格，将字符串用单引号包围，例如 oneof=red green。
//
// 跨字段约束
// validator 可以定义跨字段的约束，即该字段与其余字段之间的关系。这种约束实际上分为两种，一种是参数字段就是同一个结构中的平级字段，另外一种是参数字段为其他结构体中的字段。
// 若是是约束同一个结构中的字段，则在后面添加一个 field，使用 eqfield 定义字段间的相等约束。若是是更深层次的字段，在 field 以前还须要加上cs（cross-struct），
// eq 就变为 eqcsfield
//
// 字符串
//
// 字符串约束 tags
//
// contains=: 包含参数子串，例如 contains=email
// containsany: 包含参数中任意的 UNICODE 字符，例如 containsany=abcd
// containsrune: 包含参数表示的 rune 字符，例如 containsrune=☻
// excludes: 不包含参数子串，例如 excludes=email
// excludesall: 不包含参数中任意的 UNICODE 字符，例如 excludesall=abcd
// excludesrune: 不包含参数表示的 rune 字符，excludesrune=☻
// startswith: 以参数子串为前缀，例如 startswith=hello
// endswith: 以参数子串为后缀，例如 endswith=bye
//
// 惟一性
// 使用 unqiue 来指定惟一性约束，对不一样类型的处理以下:
//
// 对于数组和切片，unique 约束没有重复的元素；
// 对于 map，unique 约束没有重复的值；
// 对于元素类型为结构体的切片，unique 约束结构体对象的某个字段不重复，使用 unqiue=field 指定这个字段名。
//
// 特殊的约束
//
// -: 跳过该字段，不检验
// |: 使用多个约束，只须要知足其中一个，例如 rgb|rgba
// required: 字段必须设置，不能为默认值
// omitempty: 若是字段未设置，则忽略它

type User struct {
	Name string `validate:"min=6,max=10"`
	Age  int    `validate:"min=1,max=100"`
}

type UserCustom struct {
	Name        string `validate:"name"`
	Description string `validate:"description"`
}

type UserP struct {
	Password string `validate:"required"`
}

type UserE struct {
	Email string `validate:"required,email,min=1,max=100"`
}

type UserC struct {
	Name    string    `validate:"ne=admin"`          // 字符串不能是 admin
	Age     int       `validate:"gte=18"`            // 必须大于等于 18
	Sex     string    `validate:"oneof=male female"` // 性别必须是 male 和 female 其中一个
	RegTime time.Time `validate:"lte"`               // 注册时间必须小于当前的 UTC 时间，注意若是字段类型是 time.Time，使用 gt/gte/lt/lte 等约束时不用指定参数值，默认与当前的 UTC 时间比较。
}

type UserCS struct {
	Name      string `validate:"min=2"`
	Age       int    `validate:"min=18"`
	Password  string `validate:"min=10"`
	Password2 string `validate:"eqfield=Password"` // 约束两次输入的密码必须相等
}

type UserStr struct {
	Name string `validate:"containsrune=☻"`
	Age  int    `validate:"min=18"`
}

type UserUni struct {
	Name    string    `validate:"min=2"`
	Age     int       `validate:"min=18"`
	Hobbies []string  `validate:"unique"`      // Hobbies 中不能有重复元素
	Friends []UserUni `validate:"unique=Name"` // Friends 的各个元素不能有一样的 Name
}

func TestValidateStruct(t *testing.T) {
	validate := validator.New()

	t.Run("min and max", func(t *testing.T) {
		tests := []struct {
			u        User
			expected string
		}{
			{User{Name: "lidajun", Age: 18}, ""},
			{User{Name: "dj", Age: 101}, `Key: 'User.Name' Error:Field validation for 'Name' failed on the 'min' tag
Key: 'User.Age' Error:Field validation for 'Age' failed on the 'max' tag`},
		}

		for _, v := range tests {
			err := validate.Struct(v.u)
			if v.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.expected)
			}
		}

	})

	t.Run("required", func(t *testing.T) {
		tests := []struct {
			u        UserP
			expected string
		}{
			{UserP{Password: "lidajun"}, ""},
			{UserP{Password: ""}, "Key: 'UserP.Password' Error:Field validation for 'Password' failed on the 'required' tag"},
			{UserP{}, "Key: 'UserP.Password' Error:Field validation for 'Password' failed on the 'required' tag"},
		}

		for _, v := range tests {
			err := validate.Struct(v.u)
			if v.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.expected)
			}
		}
	})

	t.Run("email", func(t *testing.T) {
		tests := []struct {
			u        UserE
			expected string
		}{
			{UserE{}, "Key: 'UserE.Email' Error:Field validation for 'Email' failed on the 'required' tag"},
			{UserE{Email: "lidajun"}, "Key: 'UserE.Email' Error:Field validation for 'Email' failed on the 'email' tag"},
			{UserE{Email: "pooky.shi@gmail"}, "Key: 'UserE.Email' Error:Field validation for 'Email' failed on the 'email' tag"},
			{UserE{Email: "pooky.shi@gmail.com"}, ""},
		}

		for _, v := range tests {
			err := validate.Struct(v.u)
			if v.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.expected)
			}
		}
	})

	t.Run("range tags", func(t *testing.T) {
		tests := []struct {
			u        UserC
			expected string
		}{
			{UserC{Name: "dj", Age: 18, Sex: "male", RegTime: time.Now().UTC()}, ""},
			{UserC{Name: "admin", Age: 15, Sex: "none", RegTime: time.Now().UTC().Add(1 * time.Hour)}, `Key: 'UserC.Name' Error:Field validation for 'Name' failed on the 'ne' tag
Key: 'UserC.Age' Error:Field validation for 'Age' failed on the 'gte' tag
Key: 'UserC.Sex' Error:Field validation for 'Sex' failed on the 'oneof' tag
Key: 'UserC.RegTime' Error:Field validation for 'RegTime' failed on the 'lte' tag`},
		}

		for _, v := range tests {
			err := validate.Struct(v.u)
			if v.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.expected)
			}
		}
	})

	t.Run("cross-struct field", func(t *testing.T) {
		tests := []struct {
			u        UserCS
			expected string
		}{
			{UserCS{
				Name:      "dj",
				Age:       18,
				Password:  "1234567890",
				Password2: "1234567890",
			}, ""},
			{UserCS{
				Name:      "dj",
				Age:       18,
				Password:  "1234567890",
				Password2: "123",
			}, "Key: 'UserCS.Password2' Error:Field validation for 'Password2' failed on the 'eqfield' tag"},
		}

		for _, v := range tests {
			err := validate.Struct(v.u)
			if v.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.expected)
			}
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			u        UserStr
			expected string
		}{
			{UserStr{"d☻j", 18}, ""},
			{UserStr{"dj", 18}, "Key: 'UserStr.Name' Error:Field validation for 'Name' failed on the 'containsrune' tag"},
		}

		for _, v := range tests {
			err := validate.Struct(v.u)
			if v.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.expected)
			}
		}
	})

	t.Run("unique", func(t *testing.T) {
		u1 := UserUni{
			Name: "dj2",
			Age:  18,
		}
		u2 := UserUni{
			Name: "dj3",
			Age:  18,
		}
		u3 := UserUni{
			Name: "dj3",
			Age:  18,
		}
		tests := []struct {
			u        UserUni
			expected string
		}{
			{UserUni{
				Name:    "dj",
				Age:     18,
				Hobbies: []string{"pingpong", "chess", "programming"},
				Friends: []UserUni{u1, u2},
			}, ""},
			{UserUni{
				Name:    "dj",
				Age:     18,
				Hobbies: []string{"programming", "programming"},
				Friends: []UserUni{u1, u2, u3},
			}, `Key: 'UserUni.Hobbies' Error:Field validation for 'Hobbies' failed on the 'unique' tag
Key: 'UserUni.Friends' Error:Field validation for 'Friends' failed on the 'unique' tag`},
		}

		for _, v := range tests {
			err := validate.Struct(v.u)
			if v.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.expected)
			}
		}
	})
}

func TestValidateLocale(t *testing.T) {
	// locale 通常取决于 http 请求头的 'Accept-Language'
	zhT := zh.New() // chinese
	enT := en.New() // english
	uni := ut.New(enT, zhT, enT)

	validate := validator.New()

	t.Run("English", func(t *testing.T) {
		lang := "en"
		trans, _ := uni.GetTranslator(lang)
		_ = en_translations.RegisterDefaultTranslations(validate, trans)
		tests := []struct {
			u        User
			expected string
		}{
			{User{Name: "lidajun", Age: 18}, ""},
			{User{Name: "dj", Age: 101}, `Key: 'User.Name' Error:Field validation for 'Name' failed on the 'min' tag
Key: 'User.Age' Error:Field validation for 'Age' failed on the 'max' tag`},
		}

		for _, v := range tests {
			err := validate.Struct(v.u)
			if v.expected == "" {
				assert.NoError(t, err)
			} else {
				if tes, ok := err.(validator.ValidationErrors); ok {
					for _, val := range tes {
						t.Log(val.Translate(trans))
						// Name must be at least 6 characters in length
						// Age must be 100 or less
					}
				}
			}
		}
	})

	t.Run("Chinese", func(t *testing.T) {
		lang := "zh"
		trans, _ := uni.GetTranslator(lang)
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
		tests := []struct {
			u        User
			expected string
		}{
			{User{Name: "lidajun", Age: 18}, ""},
			{User{Name: "dj", Age: 101}, `Key: 'User.Name' Error:Field validation for 'Name' failed on the 'min' tag
Key: 'User.Age' Error:Field validation for 'Age' failed on the 'max' tag`},
		}

		for _, v := range tests {
			err := validate.Struct(v.u)
			if v.expected == "" {
				assert.NoError(t, err)
			} else {
				if tes, ok := err.(validator.ValidationErrors); ok {
					for _, val := range tes {
						t.Log(val.Translate(trans))
						// Name长度必须至少为6个字符
						// Age必须小于或等于100
					}
				}
			}
		}
	})
}

// validator 校验返回的结果只有 3 种状况:
//
// nil: 没有错误
// InvalidValidationError: 输入参数错误
// ValidationErrors: 字段违反约束
func TestValidateError(t *testing.T) {
	validate := validator.New()

	err := validate.Struct(1)
	handleErr(err, t) // param error: validator: (nil int)

	err = validate.VarWithValue(1, 2, "eqfield")
	handleErr(err, t) // Key: '' Error:Field validation for '' failed on the 'eqfield' tag
}

// 一些很简单的状况下，仅仅想对两个变量进行比较，可以使用 VarWithValue，只须要传入要验证的两个变量和约束
func TestWithValue(t *testing.T) {
	name1 := "dj"
	name2 := "dj2"

	validate := validator.New()
	err := validate.VarWithValue(name1, name2, "eqfield")
	assert.EqualError(t, err, "Key: '' Error:Field validation for '' failed on the 'eqfield' tag")

	err = validate.VarWithValue(name1, name2, "nefield")
	assert.NoError(t, err)
}

func TestCustomValidator(t *testing.T) {
	zhT := zh.New() // chinese
	enT := en.New() // english
	uni := ut.New(enT, zhT, enT)

	validate := validator.New()

	lang := "en"
	trans, _ := uni.GetTranslator(lang)
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	// RegisterValidation 方法将该约束注册到指定的 tag 上
	_ = validate.RegisterValidation("description", validateDescription) // nolint: errcheck // no need
	_ = validate.RegisterValidation("name", validateName)               // nolint: errcheck // no need

	// additional translations
	translations := []struct {
		tag         string
		translation string
	}{
		{
			tag:         "description",
			translation: fmt.Sprintf("{0} must be at least %d characters in length", minDescriptionLength),
		},
		{
			tag:         "name",
			translation: "{0}: {1} is not an invalid name",
		},
	}
	for _, tl := range translations {
		err := validate.RegisterTranslation(tl.tag, trans, registrationFunc(tl.tag, tl.translation), translateFunc)
		if err != nil {
			assert.NoError(t, err)
		}
	}
	tests := []struct {
		u        UserCustom
		expected string
	}{
		{UserCustom{Name: "lidajun", Description: "df"}, "sdf"},
		{UserCustom{Name: "d", Description: "dfdddd"}, "sdf"},
	}
	for _, v := range tests {
		err := validate.Struct(v.u)
		if v.expected == "" {
			assert.NoError(t, err)
		} else {
			if tes, ok := err.(validator.ValidationErrors); ok {
				for _, val := range tes {
					t.Log(val.Translate(trans))
					// Description must be at least 5 characters in length
					// Name is not an invalid name
				}
			}
		}
	}
}

func TestIsValidPassword(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		u        UserP
		expected string
	}{
		{UserP{Password: "lidajun"}, "uppercase letter missing, at least one numeric character required, special character missing, password length must be between 8 to 16 characters long"},
		{UserP{Password: "XXX"}, "lowercase letter missing, at least one numeric character required, special character missing, password length must be between 8 to 16 characters long"},
		{UserP{Password: "XXXlll"}, "at least one numeric character required, special character missing, password length must be between 8 to 16 characters long"},
		{UserP{Password: "XXXlll@"}, "at least one numeric character required, password length must be between 8 to 16 characters long"},
		{UserP{Password: "XXXlll@111"}, ""},
	}

	for _, v := range tests {
		err := validate.Struct(v.u)
		assert.NoError(t, err)

		err = IsValidPassword(v.u.Password)
		if v.expected == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, v.expected)
		}
	}
}

func registrationFunc(tag string, translation string) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		// adds a normal translation for a particular language/locale
		// {#} is the only replacement type accepted and are ad infinitum
		// eg. one: '{0} day left' other: '{0} days left'
		if err = ut.Add(tag, translation, true); err != nil {
			return
		}
		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field(), reflect.ValueOf(fe.Value()).String())
	if err != nil {
		return fe.(error).Error()
	}

	return t
}

// validateDescription checks if a given description is illegal.
func validateDescription(fl validator.FieldLevel) bool {
	// 通过 FieldLevel 取出要检查的字段的信息
	description := fl.Field().String()

	return len(description) >= minDescriptionLength
}

// validateName checks if a given name is illegal.
func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	return len(name) >= minNameLength
}

// handleErr print error details
func handleErr(err error, t *testing.T) {
	if err == nil {
		return
	}

	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		t.Logf("param error: %s", invalid)
		return
	}

	validationErrs := err.(validator.ValidationErrors)
	for _, validationErr := range validationErrs {
		t.Log(validationErr)
	}
}
