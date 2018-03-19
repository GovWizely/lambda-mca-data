package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mmcdole/gofeed"
)

type Response struct {
	Message string
}

type McaDataItem struct {
	Title       string
	Link        string
	Description string
	PubDate     string
	GUID        string
	Categories  string
	CountryCode string
	CountryName string
}

var re = regexp.MustCompile("country/(..) - (.*)")

const FILENAME = "mca-data.json"
const BUCKET = "mca-trade-leads-data"
const REGION = "us-east-1"
const URL = "https://www.dgmarket.com/tenders/RssFeedAction.do?locationISO=&keywords=Millennium+Challenge+Account&" +
	"sub=&noticeType=gpn%2cpp%2cspn%2crfc&language"

func getCountryInfo(categories []string) (countryIndex int, countryCode string, countryName string) {
	for idx, category := range categories {
		if match := re.FindStringSubmatch(category); len(match) > 0 {
			countryIndex, countryCode, countryName = idx, strings.ToUpper(match[1]), match[2]
			break
		}
	}
	return
}

func remove(stringArray []string, i int) []string {
	stringArray[i] = stringArray[len(stringArray)-1]
	return stringArray[:len(stringArray)-1]
}

func getMcaDataItem(item *gofeed.Item) McaDataItem {
	categories := item.Categories
	countryIndex, countryCode, countryName := getCountryInfo(categories)
	remainingCategories := remove(categories, countryIndex)
	sort.Strings(remainingCategories)
	remainingCategoriesCsv := strings.Join(remainingCategories, ",")
	mcaDataItem := McaDataItem{Title: item.Title, Link: item.Link, Description: item.Description,
		PubDate: item.Published, GUID: item.GUID, CountryCode: countryCode, CountryName: countryName,
		Categories: remainingCategoriesCsv}
	return mcaDataItem
}

func processFeed() string {
	fmt.Println("Setting up S3 session from env vars...")
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{Config: aws.Config{Region: aws.String(REGION)}}))
	uploader := s3manager.NewUploader(sess)

	fmt.Println("Reading MCA Data feed...")
	fp := gofeed.NewParser()
	feed, feedErr := fp.ParseURL(URL)
	if feedErr != nil {
		// Still report "success" when there's an empty feed or the feed is down. Not our error to fix.
		return fmt.Sprintf("Error parsing MCA Data feed: %s", feedErr)
	}
	fmt.Println("Number of trade leads:", len(feed.Items))

	mcaDataItems := make([]McaDataItem, len(feed.Items))
	for i, item := range feed.Items {
		mcaDataItems[i] = getMcaDataItem(item)
	}
	mcaDataItemsJsonBytes, _ := json.Marshal(mcaDataItems)
	ioReader := bytes.NewReader(mcaDataItemsJsonBytes)

	fmt.Println("Uploading file to S3...")
	upParams := s3manager.UploadInput{
		Bucket:      aws.String(BUCKET),
		Key:         aws.String(filepath.Base(FILENAME)),
		Body:        ioReader,
		ContentType: aws.String("application/json"),
	}
	result, err := uploader.Upload(&upParams)
	if err != nil {
		fmt.Println("Error uploading MCA JSON to S3: ", err)
		os.Exit(1)
	}

	return fmt.Sprintf("Successfully uploaded %s to %s", FILENAME, result.Location)
}

func Handler() (Response, error) {
	responseStr := processFeed()
	return Response{Message: responseStr}, nil
}

func main() {
	lambda.Start(Handler)
}
