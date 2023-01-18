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

func main() {

	cred := Credentials{
		AccessToken:       AccessToken,
		AccessTokenSecret: AccessTokenSecret,
		ConsumerKey:       APIKey,
		ConsumerSecret:    APIKeySecret,
	}
	fmt.Println("Loading Quran verses...")
	verses, err := loadVerses()
	rand.Seed(time.Now().Unix())

	if err != nil {
		fmt.Println("Error loading quran verses", err.Error())
		return
	}

	for {
		client, err := getClient(&cred)
		if err != nil {
			log.Println("Error getting Twitter Client")
			log.Println(err)
		}

		randomVerse := verseToTweetBody(&verses[rand.Intn(len(verses))])

		for len(randomVerse) > 280 {
			randomVerse = verseToTweetBody(&verses[rand.Intn(len(verses))])
		}

		fmt.Println(randomVerse)

		tweet, resp, err := client.Statuses.Update(randomVerse, nil)
		if err != nil {
			fmt.Println("ERROR:")
			log.Println(err)
		}

		fmt.Printf("Tweeted at %v:\n %v \n", time.Now(), tweet.FullText)
		fmt.Printf("Response: %v", resp.Status)
		resp.Body.Close()

		time.Sleep(2 * time.Hour)
	}

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
	const fileName string = "quran-dataset.csv"
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
