package fuzzteacup

import (
    "strconv"
    "image"
    "bytes"
    fuzz "github.com/AdaLogics/go-fuzz-headers"

    "github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"

    "github.com/knipferrc/teacup/code"
    "github.com/knipferrc/teacup/filetree"
    teacupImage "github.com/knipferrc/teacup/image"
    "github.com/knipferrc/teacup/markdown"
)

func mayhemit(data []byte) int {

    var num int
    if len(data) > 2 {
        num, _ = strconv.Atoi(string(data[0]))
        data = data[1:]
        fuzzConsumer := fuzz.NewConsumer(data)
        
        switch num {
            case 0:
                testContent, _ := fuzzConsumer.GetString()
                testExtension, _ := fuzzConsumer.GetString()
                testSyntaxTheme, _ := fuzzConsumer.GetString()

                code.Highlight(testContent, testExtension, testSyntaxTheme)
                return 0

            case 1:
                testActive, _ := fuzzConsumer.GetBool()
                testBorderless, _ := fuzzConsumer.GetBool()

                var testColor lipgloss.AdaptiveColor
                fuzzConsumer.GenerateStruct(&testColor)

                code.New(testActive, testBorderless, testColor)
                return 0

            case 2:
                var testBubble code.Bubble
                fuzzConsumer.GenerateStruct(&testBubble)
                testName, _ := fuzzConsumer.GetString()

                testBubble.SetFileName(testName)
                return 0

            case 3:
                var testBubble code.Bubble
                fuzzConsumer.GenerateStruct(&testBubble)
                testW, _ := fuzzConsumer.GetInt()
                testH, _ := fuzzConsumer.GetInt()

                testBubble.SetSize(testW, testH)
                return 0

        
            case 4:
                var testMsg tea.Msg
                var testBubble code.Bubble
                fuzzConsumer.GenerateStruct(&testBubble)

                testBubble.Update(testMsg)
                return 0

            case 5:
                testInt, _ := fuzzConsumer.GetInt()
                testInt64 := int64(testInt)

                filetree.ConvertBytesToSizeString(testInt64)
                return 0

            case 6:
                img, _, _ := image.Decode(bytes.NewReader(data))
                testWidth, _ := fuzzConsumer.GetInt()

                teacupImage.ToString(testWidth, img)
                return 0

            default:
                testWidth, _ := fuzzConsumer.GetInt()
                testContent, _ := fuzzConsumer.GetString()

                markdown.RenderMarkdown(testWidth, testContent)
                return 0



        }
    }
    return 0
}

func Fuzz(data []byte) int {
    _ = mayhemit(data)
    return 0
}