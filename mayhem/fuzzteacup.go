package fuzzteacup

import (
	"bytes"
	"image"
	"strconv"

	fuzz "github.com/AdaLogics/go-fuzz-headers"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mistakenelf/teacup/code"
	"github.com/mistakenelf/teacup/filetree"
	teacupImage "github.com/mistakenelf/teacup/image"
	"github.com/mistakenelf/teacup/markdown"
)

func mayhemit(data []byte) int {
	if len(data) <= 2 {
		return 0
	}
	num, _ := strconv.Atoi(string(data[0]))
	data = data[1:]
	fuzzConsumer := fuzz.NewConsumer(data)

	switch num {
	case 1:
		testActive, _ := fuzzConsumer.GetBool()
		code.New(testActive)
		return 0
	case 2:
		var testModel code.Model
		fuzzConsumer.GenerateStruct(&testModel)
		testName, _ := fuzzConsumer.GetString()
		testModel.SetFileName(testName)
		return 0
	case 3:
		var testModel code.Model
		fuzzConsumer.GenerateStruct(&testModel)
		testW, _ := fuzzConsumer.GetInt()
		testH, _ := fuzzConsumer.GetInt()
		testModel.SetSize(testW, testH)
		return 0
	case 4:
		var testMsg tea.Msg
		var testModel code.Model
		fuzzConsumer.GenerateStruct(&testModel)
		testModel.Update(testMsg)
		return 0
	case 5:
		testInt, _ := fuzzConsumer.GetInt()
		filetree.ConvertBytesToSizeString(int64(testInt))
		return 0
	case 6:
		img, _, _ := image.Decode(bytes.NewReader(data))
		testWidth, _ := fuzzConsumer.GetInt()
		testBorderless, _ := fuzzConsumer.GetBool()
		var testColor lipgloss.AdaptiveColor
		fuzzConsumer.GenerateStruct(&testColor)
		m := teacupImage.New(true, testBorderless, testColor)
		m.SetSize(testWidth, 10)
		_ = m
		teacupImage.ToString(testWidth, img)
		return 0
	default:
		testWidth, _ := fuzzConsumer.GetInt()
		testContent, _ := fuzzConsumer.GetString()
		m := markdown.New(true)
		m.SetSize(testWidth, 10)
		_ = m
		_, _ = markdown.RenderMarkdown(testWidth, testContent)
		return 0
	}
}

func Fuzz(data []byte) int {
	_ = mayhemit(data)
	return 0
}
