package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const dashWidth = 26
const tolerancePercentage float64 = 10
const dashes = "-"
const URL = "http://api.github.com/search/repositories"

type Repository struct {
	fullName string
	htmlURL  string
	language string
	stars    float64
}

func main() {

	introduceGame()
	input := getUserInput()
	selectedLanguage := getLanguages()[input]
	repositoryList := getResponse(selectedLanguage)
	guessStarsForRepositories(repositoryList)
}

func printDashes() {
	fmt.Println(strings.Repeat(dashes, dashWidth))
}

func introduceGame() {

	printDashes()
	fmt.Println("Welcome to GUESS THE STARS")
	printDashes()

	fmt.Println("\n" + "You have following langauages to choose from : ")
	for key, value := range getLanguages() {
		fmt.Printf("\n%d : %s", key, value)
	}
	fmt.Println("\nChoose the index for the language you want to guess: ")
}

func getUserInput() int {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	input = strings.TrimRight(input, "\r\n")
	handleErr(err)
	out, _ := strconv.Atoi(input)
	return out

}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getLanguages() []string {

	return []string{
		"Javascript",
		"Go",
		"Java",
		"Python",
		"Scala",
		"Rust",
		"C",
		"C#",
		"C++",
		"Perl",
		"Ruby",
		"PHP",
	}
}

func getResponse(selectedLanguage string) []Repository {
	pageId := 1
	for {
		req := createRequest(string(strconv.Itoa(pageId)))
		byteData := fireRequest(req)
		repositoryList := parseJsonData(byteData, selectedLanguage)
		if len(repositoryList) < 5 {
			pageId++
		} else {
			return repositoryList
		}
	}
}

func createRequest(pageId string) *http.Request {

	currentDate := time.Now().AddDate(0, 0, -7)
	pastDate := currentDate.Format("2006-01-02")
	createdDate := fmt.Sprintf("created:>%s", pastDate)

	req, err := http.NewRequest("GET", URL, nil)
	handleErr(err)
	params := req.URL.Query()
	params.Add("q", createdDate)
	params.Add("sort", "stars")
	params.Add("order", "desc")
	params.Add("page", pageId)
	params.Add("per_page", "100")
	req.URL.RawQuery = params.Encode()
	log.Println(req)
	return req
}

func fireRequest(req *http.Request) []byte {
	client := http.Client{}
	resp, err := client.Do(req)
	handleErr(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)
	return body
}

func shuffleSlice(arr []Repository) []Repository {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func parseJsonData(byteData []byte, selectedLanguage string) []Repository {

	var responseData map[string]interface{}
	json.Unmarshal(byteData, &responseData)

	items := responseData["items"].([]interface{})
	var repositoryList []Repository
	repositoryCount := 0

	for _, item := range items {
		if repositoryCount == 5 {
			break
		}
		repo := item.(map[string]interface{})
		lang, ok := repo["language"]
		if lang == nil {
			continue
		}
		if !ok {
			log.Fatal("Incorrect data received")
		}
		language := lang.(string)
		if strings.EqualFold(language, selectedLanguage) {
			repository := Repository{repo["full_name"].(string), repo["html_url"].(string), language, repo["stargazers_count"].(float64)}
			repositoryCount++
			repositoryList = append(repositoryList, repository)
		}
	}
	return repositoryList
}

func guessStarsForRepositories(repositoryList []Repository) {
	correctAnswers := 0
	repositoryList = shuffleSlice(repositoryList)
	for _, repository := range repositoryList {
		fmt.Println("\nGuess the stars for the repository : " + repository.fullName)
		input := getUserInput()
		tolerance := (tolerancePercentage / 100) * repository.stars
		if math.Abs(float64(input)-repository.stars) <= tolerance {
			fmt.Println("You've guessed it correctly!!!")
			correctAnswers++
		} else {
			fmt.Println("Your answer is incorrect!!!")
		}
		fmt.Printf("%s has %d stars\n\n", repository.fullName, int(repository.stars))
	}
	printDashes()
	fmt.Println("Final Result")
	printDashes()
	if correctAnswers >= 4 {
		fmt.Println("YOU WIN!!!")
	} else {
		fmt.Println("Sorry, You Lost!!!")
	}
	fmt.Println()
}
