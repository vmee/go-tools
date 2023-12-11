package captchax

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/mojocn/base64Captcha"
)

//DriverMath captcha config for captcha math
type DriverMath struct {
	driver *base64Captcha.DriverMath
}

//NewDriverMath creates a driver of math
func NewDriverMath(height int, width int, noiseCount int, showLineOptions int, bgColor *color.RGBA, fontsStorage base64Captcha.FontsStorage, fonts []string) *DriverMath {
	d := base64Captcha.NewDriverMath(height, width, noiseCount, showLineOptions, bgColor, fontsStorage, fonts)
	return &DriverMath{
		driver: d,
	}
}

//GenerateIdQuestionAnswer creates id,captcha content and answer
func (d *DriverMath) GenerateIdQuestionAnswer() (id, question, answer string) {
	id = base64Captcha.RandomId()
	operators := []string{"+", "-", "x"}
	var mathResult int32
	switch operators[rand.Int31n(3)] {
	case "+":
		a := rand.Int31n(20)
		b := rand.Int31n(20)
		question = fmt.Sprintf("%d+%d=?", a, b)
		mathResult = a + b
	case "x":
		a := rand.Int31n(10)
		b := rand.Int31n(10)
		question = fmt.Sprintf("%dx%d=?", a, b)
		mathResult = a * b
	default:
		a := rand.Int31n(10) + rand.Int31n(10)
		b := rand.Int31n(a)

		question = fmt.Sprintf("%d-%d=?", a, b)
		mathResult = a - b

	}
	answer = fmt.Sprintf("%d", mathResult)
	return
}

//DrawCaptcha creates math captcha item
func (d *DriverMath) DrawCaptcha(question string) (item base64Captcha.Item, err error) {
	return d.driver.DrawCaptcha(question)
}

//ConvertFonts loads fonts from names
func (d *DriverMath) ConvertFonts() *DriverMath {
	d.driver = d.driver.ConvertFonts()
	return d
}
