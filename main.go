package main

import (
	"WebCrawler/boot"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	boot.ViperSetup()
	boot.LoggerSetup()

	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{
		DomainRegexp: ``,
		Delay:        3 * time.Second,
	})
	var name string
	c.OnHTML(".h-5", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		name = e.Text
		fmt.Println(name)
		if len(href) > 1 {
			//fmt.Println(href[10 : len(href)-1])
			GetContent(href[10:len(href)-1], name)
		}

	})
	c.OnError(func(response *colly.Response, err error) {
		panic(err)
	})
	c.Visit("https://leetcode.cn/problemset/algorithms/")
}
func GetContent(href, name string) {
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{
		DomainRegexp: ``,
		Delay:        3 * time.Second,
	})
	type variables struct {
		TitleSlug string `json:"titleSlug"`
	}

	type data struct {
		OperationName string    `json:"operationName"`
		Variables     variables `json:"variables"`
		Query         string    `json:"query"`
	}

	var da = &data{
		OperationName: "questionData",
		Variables:     variables{TitleSlug: href},
		Query:         "query questionData($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n    translatedContent\n        }\n}",
	}
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(r.Request.URL)
		content := string(r.Body)
		//fmt.Println(content)

		// 去尖括号
		for i := 0; i < len(content); i++ {
			if content[i] == 60 {
				for j := i; j < len(content); j++ {
					if content[j] == 62 {
						content = content[:i] + content[j+1:]
						i = 0
						break
					} else {
						continue
					}
				}
			}
		}

		//去\n
		for i := 0; i < len(content); i++ {
			if content[i] == 92 && content[i+1] == 116 {
				content = content[:i] + content[i+2:]
				i = 0
			} else {
				continue
			}
		}

		//去&和分号的
		for i := 0; i < len(content); i++ {
			if content[i] == 38 {
				content = content[:i] + content[i+5:]
				i = 0
			} else if content[i] == 59 {
				content = content[:i] + content[i+1:]
				i = 0
			} else {
				continue
			}
		}

		if index := strings.Index(content, "\\u63d0\\u793a"); index != -1 {
			content = content[:index]
		}

		//去前面的
		content = content[42:]
		//fmt.Println(content)

		var context string
		for i := 0; i < len(content); i++ {
			if content[i] == 92 {
				if content[i+1] == 34 {
					context += "\""
					i++
					continue
				}
				if content[i+1] == 110 {
					context += "\n"
					i++
					continue
				}
				temp, err := strconv.ParseInt(content[i+2:i+6], 16, 32)
				if err != nil {
					panic(err)
				}
				context += fmt.Sprintf("%c", temp)
				i += 5
			} else {
				context += string(content[i])
			}
		}
		fmt.Println(context)
		//save([]byte(name+"\n"), []byte(context+"\n"))
		fmt.Printf("已爬取题目:%s\n", name)
	})
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Content-Length", "<calculated when request is sent>")
		request.Headers.Set("Host", "<calculated when request is sent>")
		//request.Headers.Set("Cookie", "csrftoken=QlGaAzPkx1B7pcvpv2Jdr8QWwvlRxrjZYChCJLjy3KPEBoLB3tLWVm2mkwdc7wWz")
		request.Headers.Set("content-type", "application/json")
	})
	c.OnError(func(response *colly.Response, err error) {
		panic(err.Error())
	})
	d, err := json.Marshal(da)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(d))
	err = c.PostRaw("https://leetcode.cn/graphql/", d)
	if err != nil {
		panic(err)
	}
}

func save(name, c1 []byte) {
	file, _ := os.OpenFile("problem.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	_, err := file.Write(name)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(c1)
	if err != nil {
		panic(err)
	}
	file.Close()

}
