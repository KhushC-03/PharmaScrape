package scraper

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var (
	interactions_list        = []string{}
	base_url          string = "https://bnf.nice.org.uk"
	// intercations_url  string = "https://bnf.nice.org.uk/interactions"
	test_base_url string = "http://192.168.0.158:8080"
	// test_intercations_url string = "https://bnf.nice.org.uk/interactions"

	// COUNTER                  = []string{}
)

func StripChars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

func rts(sentence string) string {
	sentences := strings.Split(sentence, ".")
	cleanedSentences := make([]string, len(sentences))
	for i, s := range sentences {
		trimmedSentence := strings.TrimSpace(s)
		cleanedSentences[i] = trimmedSentence
	}
	return strings.Join(cleanedSentences, "")
}

func ContainsString(lines []string, x string) bool {
	for _, line := range lines {
		if line == x {
			return true
		}
	}
	return false
}

func CheckFile(fileName string) []string {

	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return []string{}
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return []string{}
	}
	lines := strings.Split(string(content), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	if len(lines) == 0 {
		return []string{}
	}
	return lines

}

func countFs(sentence string) int {
	count := 0
	for _, char := range sentence {
		if char == '.' {
			count++
		}
	}
	return count
}

func Home() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link href='https://fonts.googleapis.com/css?family=Inter' rel='stylesheet'>
    <title>Search Interactions</title>
    <style>
		html{
			font-family:'Inter';
			scroll-behavior: smooth;
		}
        #search-results {
            display: none;
            position: absolute;
            background-color: white;
            border: 1px solid #ccc;
            max-height: 200px;
            overflow-y: auto;
            width: calc(100% - 2px);
        }
        #search-results ul {
            list-style-type: none;
            padding: 0;
            margin: 0;
        }
        #search-results li {
            padding: 5px 10px;
            cursor: pointer;
        }
        #search-results li:hover {
            background-color: #f0f0f0;
        }
		.severity dd, .severity dt {
			display: inline;
			margin: 0 .5rem 0 0;
			font-weight: bold;
		}
		.interactionlist{
			border-bottom: 3px solid #000000;
		}
    </style>
</head>
<body>
<h1> <a style="color: inherit;" href="/">Search Interactions</a> </h1>
<input type="text" id="search" placeholder="Search...">
<div id="search-results"></div>

<script>
    const searchInput = document.getElementById('search');
    const searchResults = document.getElementById('search-results');

    searchInput.addEventListener('input', function() {
        const query = this.value.trim();

        if (query.length === 0) {
            searchResults.style.display = 'none';
            return;
        }

        fetch("http://192.168.0.158:8080/search?s=" + query)
            .then(response => response.json())
            .then(data => {
                displayResults(data);
            })
            .catch(error => console.error('Error fetching data:', error));
    });

    function displayResults(results) {
        searchResults.innerHTML = '';
        if (results.length === 0) {
            searchResults.style.display = 'none';
            return;
        }
        const ul = document.createElement('ul');
        results.forEach(result => {
            const li = document.createElement('li');
            const link = document.createElement('a');
			link.style.color = 'inherit';
            link.textContent = result;
            link.href = "http://192.168.0.158:8080/interaction?i=" + result;
            li.appendChild(link);
            ul.appendChild(li);
        });
        searchResults.appendChild(ul);
        searchResults.style.display = 'block';
    }
</script>

	</body>	
</html>`

}

func Interactions(jsonData map[string]interface{}) string {
	sortedData := make([]string, 0, len(jsonData))
	for key := range jsonData {
		sortedData = append(sortedData, key)
	}

	sort.Strings(sortedData)

	htmlResponse := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link href='https://fonts.googleapis.com/css?family=Inter' rel='stylesheet'>
		<title>Search Interactions</title>
		<style>
		html{
			font-family:'Inter';
		}
        #search-results {
            display: none;
            position: absolute;
            background-color: white;
            border: 1px solid #ccc;
            max-height: 200px;
            overflow-y: auto;
            width: calc(100% - 2px);
        }
        #search-results ul {
            list-style-type: none;
            padding: 0;
            margin: 0;
        }
        #search-results li {
            padding: 5px 10px;
            cursor: pointer;
        }
        #search-results li:hover {
            background-color: #f0f0f0;
        }
		.severity dd, .severity dt {
			display: inline;
			margin: 0 .5rem 0 0;
		}
		.interactionlist{
			border-bottom: 3px solid #000000;
		}
		dt{
			font-weight: bold;
		}
		</style>
	</head>
	<body>
	<h1><a style="color: inherit;" href="/">Search Interactions</a></h1>
	<input type="text" id="search" placeholder="Search...">
	<div id="search-results"></div>
	<ol class="allInteractions" style="list-style-type: none;">
	`
	for _, i := range sortedData {
		var ims string = ""
		var drugIn []string
		interactionMessage := jsonData[i].(map[string]interface{})["Interaction-Messages"]
		Severity := `<dt>Severity:</dt><dd>` + jsonData[i].(map[string]interface{})["Severity"].(string) + `</dd>`
		Evidence := `<dt>Evidence:</dt><dd>` + jsonData[i].(map[string]interface{})["Evidence"].(string) + `</dd>`
		if Severity == "<dt>Severity:</dt><dd></dd>" {
			Severity = ""
		}
		if Evidence == "<dt>Evidence:</dt><dd></dd>" {
			Evidence = ""
		}
		if interactionMessage == nil {
			continue
		}

		if len(interactionMessage.([]interface{})) == 1 {
			ims = `<p class="interaction-message"> ` + rts(interactionMessage.([]interface{})[0].(string)) + "</p>"
		} else {
			for im := range interactionMessage.([]interface{}) {
				if !ContainsString(drugIn, rts(interactionMessage.([]interface{})[im].(string))) {
					ims += `<p class="interaction-message"> ` + rts(interactionMessage.([]interface{})[im].(string)) + "</p>"
					drugIn = append(drugIn, rts(interactionMessage.([]interface{})[im].(string)))

				}

			}
		}

		htmlResponse += `<li class="interactionlist">
		<h3 class="interacton-title">` + jsonData[i].(map[string]interface{})["Drug"].(string) + `</h3>
		<ul class="interaction-information">
		<li class="Interaction-module--message">
		<dl class="severity">` + ims + Severity + Evidence + `
				</dl>
		</li>
		</ul>
	</li>`

	}
	htmlResponse += `
	</ol>
	<script>
    const searchInput = document.getElementById('search');
    const searchResults = document.getElementById('search-results');

    searchInput.addEventListener('input', function() {
        const query = this.value.trim();

        if (query.length === 0) {
            searchResults.style.display = 'none';
            return;
        }

        fetch("http://192.168.0.158:8080/search?s=" + query)
            .then(response => response.json())
            .then(data => {
                displayResults(data);
            })
            .catch(error => console.error('Error fetching data:', error));
    });

    function displayResults(results) {
		console.log(results);
        searchResults.innerHTML = '';
        if (results.length === 0) {
            searchResults.style.display = 'none';
            return;
        }
        const ul = document.createElement('ul');
        results.forEach(result => {
            const li = document.createElement('li');
            const link = document.createElement('a');
			link.style.color = 'inherit';
            link.textContent = result;
            link.href = "http://192.168.0.158:8080/interaction?i=" + result;
            li.appendChild(link);
            ul.appendChild(li);
        });
        searchResults.appendChild(ul);
        searchResults.style.display = 'block';
    }
</script>

	</body>	
</html>
	`
	return htmlResponse
}
