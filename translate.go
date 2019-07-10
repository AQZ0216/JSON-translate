// Sample translate-quickstart translates "Hello, world!" into Russian.
package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"

	"encoding/json"
	"io/ioutil"
	"strings"
)

// reference to https://github.com/GoogleCloudPlatform/golang-samples/blob/master/translate/snippets/snippet.go
func translateText(targetLanguage, text string) (string, error) {
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", err
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", err
	}
	return resp[0].Text, nil
}

func translateInterface(source interface{}, target interface{}, targetLanguage string) (interface{}, error) {
	sourceMap := source.(map[string]interface{})
	targetMap := make(map[string]interface{})
	switch v := target.(type) {
	case map[string]interface{}:
		targetMap = v
	default:
	}

	for k, i := range sourceMap {
		switch v := i.(type) {
		case map[string]interface{}:
			_, ok := targetMap[k]
			if ok == false {
				targetMap[k] = make(map[string]interface{})
			}
			var err error
			targetMap[k], err = translateInterface(v, targetMap[k], targetLanguage)
			if err != nil {
				return target, err
			}
		case interface{}:
			_, ok := targetMap[k]
			if ok == false {
				srcStr := v.(string)
				targetStr, err := translateText(targetLanguage, srcStr)
				if err != nil {
					return target, err
				}

				// process "{...}""
				targetLeftIdx := 0
				srcLeftIdx := 0
				targetRightIdx := 0
				srcRightIdx := 0
				for strings.Contains(targetStr[targetRightIdx:], "{") {
					targetLeftIdx = strings.Index(targetStr[targetRightIdx:], "{") + targetRightIdx
					srcLeftIdx = strings.Index(srcStr[srcRightIdx:], "{") + srcRightIdx
					targetRightIdx = strings.Index(targetStr[targetLeftIdx:], "}") + targetLeftIdx
					srcRightIdx = strings.Index(srcStr[srcLeftIdx:], "}") + srcLeftIdx
					if targetStr[targetLeftIdx+1] < 'A' || targetStr[targetLeftIdx+1] > 'z' {
						targetStr = targetStr[:targetLeftIdx] + srcStr[srcLeftIdx:srcRightIdx] + targetStr[targetRightIdx:]
					}
				}
				targetMap[k] = targetStr
				fmt.Println(v, ", ", targetMap[k])
			}
		default:
		}
	}

	return targetMap, nil
}

func translateJSON(source []byte, target []byte, targetLanguage string) ([]byte, error) {
	//step 1: decode JSON
	var sourceInterface interface{}
	err := json.Unmarshal(source, &sourceInterface)
	if err != nil {
		return target, err
	}

	var targetInterface interface{}
	err = json.Unmarshal(target, &targetInterface)
	if err != nil {
		return target, err
	}

	//step 2 : translate
	targetInterface, err = translateInterface(sourceInterface, targetInterface, targetLanguage)
	if err != nil {
		return target, err
	}

	//step 3: encode JSON
	b, err := json.Marshal(targetInterface)
	if err != nil {
		return target, err
	}

	return b, err
}

func main() {
	//part 1 : read file
	enJSON, err := ioutil.ReadFile("./en.json")
	if err != nil {
		log.Fatal(err)
	}

	jaJSON, err := ioutil.ReadFile("./ja.json")
	if err != nil {
		log.Fatal(err)
	}

	koJSON, err := ioutil.ReadFile("./ko.json")
	if err != nil {
		log.Fatal(err)
	}

	zhCNJSON, err := ioutil.ReadFile("./zh-cn.json")
	if err != nil {
		log.Fatal(err)
	}

	zhTWJSON, err := ioutil.ReadFile("./zh-tw.json")
	if err != nil {
		log.Fatal(err)
	}

	//part 2 : en -> zh-tw
	zhTWJSON, err = translateJSON(enJSON, zhTWJSON, "zh-TW")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./zh-tw.json", zhTWJSON, 0777)
	if err != nil {
		log.Fatal(err)
	}

	//part 3 : zh-tw -> zh-cn
	zhCNJSON, err = translateJSON(zhTWJSON, zhCNJSON, "zh-CN")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./zh-cn.json", zhCNJSON, 0777)
	if err != nil {
		log.Fatal(err)
	}

	//part 4 : zh-tw -> ja
	jaJSON, err = translateJSON(zhTWJSON, jaJSON, "ja")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./ja.json", jaJSON, 0777)
	if err != nil {
		log.Fatal(err)
	}

	//part 5 : zh-tw -> ko
	koJSON, err = translateJSON(zhTWJSON, koJSON, "ko")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./ko.json", koJSON, 0777)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Finish")
}
