package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

func SystemEngine() {
	client, _ := CreateSession()
	allDrugs := fetchdrugs(client)

	file, err := os.OpenFile("./cache/endpoint-cache.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error Opening [cache/endpoint-cache.txt]: ", err)
		return
	}
	defer file.Close()
	index := 0
	for _, match := range allDrugs {
		if strings.Contains(match, "/drugs") && !strings.Contains(match, "SubNav") && !strings.Contains(match, "otherSiteTab") {
			// if !ContainsString(cache, strings.Split(match, `"`)[1]) {
			index++
			_, err = file.WriteString(fmt.Sprintf("%s\n", strings.Split(match, `"`)[1]))

			if err != nil {
				fmt.Println("Error Writing To [cache/endpoint-cache.txt]: ", err)
				return
			}
			// }
		}
	}
	for _, match := range allDrugs {
		fetchInteractionURL(strings.Replace(base_url+"/"+fmt.Sprintf("%s\n", strings.Split(match, `"`)[1]), "/drugs/", "drugs/", -1), client)

	}
	fmt.Printf("%d New Drug Endpoints Found & Cached\n", index)
	strMap := make(map[string]bool)
	uniqueStrArray := []string{}
	interactionsCache := CheckFile("./cache/interaction-endpoint-cache.txt")
	interactionsFile, err := os.OpenFile("./cache/interaction-endpoint-cache.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error Opening [interaction-endpoint-cache.txt]: ", err)
		return
	}
	defer interactionsFile.Close()
	for _, str := range interactions_list {
		if _, ok := strMap[str]; !ok {
			strMap[str] = true
			uniqueStrArray = append(uniqueStrArray, str)
		}
	}
	for _, interactionURL := range uniqueStrArray {
		if !ContainsString(interactionsCache, interactionURL) {

			_, err = interactionsFile.WriteString(interactionURL + "\n")
			if err != nil {
				fmt.Println("Error Writng To [interaction-endpoint-cache.txt]: ", err)
				return
			}
		}
		FetchInteractions(interactionURL, client)

	}

}

func fetchdrugs(client *http.Client) []string {
	// cache := CheckFile("./cache/endpoint-cache.txt")

	request, err := http.NewRequest("GET", test_base_url+"/"+"drugs", nil)
	if err != nil {
		fmt.Printf("Error Initiating Request")
	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Sec-Fetch-Dest", "document")
	request.Header.Set("Accept-Language", "en-GB,en;q=0.9")
	request.Header.Set("Sec-Fetch-Mode", "navigate")
	request.Header.Set("Host", "bnf.nice.org.uk")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1.2 Safari/605.1.15")
	request.Header.Set("Referer", "https://bnf.nice.org.uk/sw.js")
	request.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var re = regexp.MustCompile(`(?m)<a\s+(?:[^>]*?\s+)?href=["']?([^"']+)["']?[^>]*>`)
		return re.FindAllString(string(bodyText), -1)
	}
	return []string{}
}

func fetchInteractionURL(drugUrl string, client *http.Client) {

	request, err := http.NewRequest("GET", StripChars(drugUrl, "\n"), nil)
	if err != nil {
		fmt.Printf("Error Initiating Request")
	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Sec-Fetch-Dest", "document")
	request.Header.Set("Accept-Language", "en-GB,en;q=0.9")
	request.Header.Set("Sec-Fetch-Mode", "navigate")
	request.Header.Set("Host", "bnf.nice.org.uk")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1.2 Safari/605.1.15")
	request.Header.Set("Referer", "https://bnf.nice.org.uk/sw.js")
	request.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		re := regexp.MustCompile(`<a\s+href="([^"]*)"[^>]*>.*?<\/a>`)

		hrefs := re.FindAllStringSubmatch(string(bodyText), -1)
		var intercationURL string
		for _, href := range hrefs {
			if strings.Contains(href[1], "/interactions/") {
				if intercationURL != href[1] && href[1] != "/interactions/" {
					intercationURL = href[1]
					interactions_list = append(interactions_list, intercationURL)

					fmt.Println(intercationURL)
				}

			}
			// liPattern := regexp.MustCompile(`(?s)(<li class="{BnfInteractant-slug}-module--interactionsListItem--38172">.*?</li>)`)
			// liMatches := liPattern.FindAllString(string(bodyText), -1)
			// for _, match := range liMatches {
			// 	fmt.Println(match)
			// }
		}

	}
}

func FetchInteractions(interactionURL string, client *http.Client) {
	cacheData := make(map[string]interface{})
	request, err := http.NewRequest("GET", "https://bnf.nice.org.uk"+interactionURL, nil)
	if err != nil {
		fmt.Printf("Error Initiating Request")
	}
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Sec-Fetch-Dest", "document")
	request.Header.Set("Accept-Language", "en-GB,en;q=0.9")
	request.Header.Set("Sec-Fetch-Mode", "navigate")
	request.Header.Set("Host", "bnf.nice.org.uk")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1.2 Safari/605.1.15")
	request.Header.Set("Referer", "https://bnf.nice.org.uk/sw.js")
	request.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		doc, _ := html.Parse(strings.NewReader(string(bodyText)))
		fmt.Println(request.URL.String())
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "li" {
				for _, a := range n.Attr {
					if a.Key == "class" && strings.Contains(a.Val, "{BnfInteractant-slug}-module--interactionsListItem--38172") {
						var b bytes.Buffer
						html.Render(&b, n)
						reTag := regexp.MustCompile(`<a\s+href="/drugs/([^"]+)/">([^<]+)</a>|<h3[^>]*>([^<]+)</h3>`)

						matches := reTag.FindStringSubmatch(b.String())

						var interactionDrug string
						for _, rmatch := range matches {
							if len(rmatch) > 0 && unicode.IsUpper(rune(rmatch[0])) {
								interactionDrug = rmatch
							}
						}
						re := regexp.MustCompile(`<p class="interaction-message">(.*?)<\/p>`)
						imatches := re.FindAllStringSubmatch(b.String(), -1)
						var interactionMessage string = ""
						var interactionData []string
						for _, imatch := range imatches {
							if len(imatch[1]) > 0 && unicode.IsUpper(rune(imatch[1][0])) {

								if strings.Contains(imatch[1], "<a href") {
									removeHref := strings.Split(strings.Split(imatch[1], `<a href=`)[1], `"`)[1]
									removeTitle := strings.Split(strings.Split(imatch[1], `title=`)[1], `"`)[1]
									interactionMessage += strings.Replace(imatch[1], `<a href="`+removeHref+`" title="`+removeTitle+`">Guidance on Prescribing</a>`, "", -1)

								} else {
									interactionMessage += imatch[1]
								}

								if countFs(interactionMessage) > 1 {
									splitInteraction := strings.Split(interactionMessage, ".")

									// Remove empty strings (due to periods at the end)
									for _, s := range splitInteraction {
										if s != "" {
											interactionData = append(interactionData, strings.TrimSpace(s))
										}
									}
								} else {
									interactionData = append(interactionData, interactionMessage)

								}

							}
						}
						var Severity string = "Not Available"
						var Evidence string = "Not Available"
						if strings.Contains(b.String(), "Severity") {
							severityEx := regexp.MustCompile(`<dd>(.*?)<\/dd>`)
							severityMatches := severityEx.FindAllStringSubmatch(b.String(), -1)
							Severity = strings.Split(severityMatches[0][0], "<dd>")[1][:len(strings.Split(severityMatches[0][0], "<dd>")[1])-5]
						}
						if strings.Contains(b.String(), "Evidence") {
							evidenceEx := regexp.MustCompile(`<dd>(.*?)<\/dd>`)
							evidenceMatches := evidenceEx.FindAllStringSubmatch(b.String(), -1)
							Evidence = strings.Split(evidenceMatches[1][0], "<dd>")[1][:len(strings.Split(evidenceMatches[1][0], "<dd>")[1])-5]
						}
						cacheData[interactionDrug] = map[string]interface{}{
							"Drug":                 interactionDrug,
							"Severity":             Severity,
							"Evidence":             Evidence,
							"Interaction-Messages": interactionData,
						}

					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
		jsonData, err := json.MarshalIndent(cacheData, "", "  ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		file, err := os.Create(fmt.Sprintf("./cache/interactions/%s.json", strings.ReplaceAll(interactionURL[:len(interactionURL)-1], "/", "-")))
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

	}
}
