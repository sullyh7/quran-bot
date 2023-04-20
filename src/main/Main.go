package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
)

var APIKey string = getConfig("APIKey")
var APIKeySecret string = getConfig("APIKeySecret")

var AccessToken string = getConfig("AccessToken")
var AccessTokenSecret string = getConfig("AccessTokenSecret")

const fileName string = "quran-dataset.csv"

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var purple = "\033[35m"
var cyan = "\033[36m"
var gray = "\033[37m"
var white = "\033[97m"

func main() {

	cred := Credentials{
		AccessToken:       AccessToken,
		AccessTokenSecret: AccessTokenSecret,
		ConsumerKey:       APIKey,
		ConsumerSecret:    APIKeySecret,
	}
	normalLog("Welcome to quran bot")
	normalLog("Loading Quran verses from " + fileName + "...")
	verses, err := loadVerses()
	rand.Seed(time.Now().Unix())

	if err != nil {
		errorLog("Error loading quran verses from " + fileName + " Error: " + err.Error())
		return
	}
	normalLog("Authenticating twitter client")
	for {
		client, err := getClient(&cred)
		if err != nil {
			errorLog("Error getting Twitter Client")
			log.Println(err)
			log.Fatal()
		}
		normalLog("========================================================")
		normalLog("Getting random verse...")
		randomVerse := verseToTweetBody(&verses[rand.Intn(len(verses))])

		for len(randomVerse) > 280 {
			randomVerse = verseToTweetBody(&verses[rand.Intn(len(verses))])
			showLog("Verse too long, getting another verse...", purple)
		}

		normalLog("Chosen verse:")
		showLog(randomVerse, cyan)

		normalLog("Tweeting verse...")
		tweet, resp, err := client.Statuses.Update(randomVerse, nil)
		if err != nil {
			errorLog("Error: " + err.Error())
			log.Println(err)
		}

		fmt.Printf("Successfully Tweeted at %v:\n ID: %v\n", time.Now(), tweet.ID)
		fmt.Printf("Response: %v\n", resp.Status)
		resp.Body.Close()

		normalLog("========================================================")

		time.Sleep(5 * time.Hour)
	}

}

func normalLog(msg string) {
	showLog(msg, blue)
}
func errorLog(msg string) {
	showLog(msg, red)
}
func successLog(msg string) {
	showLog(msg, green)
}

func showLog(message, colour string) {
	fmt.Println(colour, message)
}

func getConfig(name string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading environment variables")
	}
	return os.Getenv(name)
}

func getClient(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func verseToTweetBody(v *Verse) (body string) {
	body += v.Text + "\n"
	body += fmt.Sprintf("%d:%d\n", v.SurahNumber, v.Ayah)
	body += fmt.Sprintf("Tafsir (Ibn Kathir): https://quran.com/%v:%v/tafsirs/en-tafisr-ibn-kathir", v.SurahNumber, v.Ayah)
	return
}

func loadVerses() (verses []Verse, e error) {
	in, err := os.Open(fileName)
	if err != nil {
		return verses, fmt.Errorf("error opening file: %s", err.Error())
	}
	defer in.Close()

	if err := gocsv.UnmarshalFile(in, &verses); err != nil {
		return verses, fmt.Errorf("error unmarshalling csv file: %s", err.Error())
	}

	return
}

type Verse struct {
	SurahName   string `csv:"surah_name_en"`
	SurahNumber int    `csv:"surah_no"`
	Ayah        int    `csv:"ayah_no_surah"`
	Text        string `csv:"ayah_en"`
}

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}
