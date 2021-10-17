package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

var (
	comicWebUrl string
    comicImgUrl string

    lowBound *int
    highBound *int
    	
	generateCmd = &cobra.Command{
        Use: "gen",
        Short: "Generate a random comic or image link",
        Args: cobra.ExactValidArgs(1),
        ValidArgs: []string{"comic", "image"},
        Run: func(cmd *cobra.Command, args []string) {

            highestFlag := cmd.Flags().Lookup("highest")

            if(!highestFlag.Changed) {
                highestFlag.Value.Set(strconv.Itoa(getNewestComicNum()))
            }
            
            comicWebUrl = getComicUrl(*lowBound, *highBound)
            comicImgUrl = getComicImgUrl(comicWebUrl)
            
            switch comicUrlType := args[0]; comicUrlType{
            case "comic":
                fmt.Println(comicWebUrl)
            
            case "image":
                fmt.Println(comicImgUrl)
            }

        },
	}

)

func init() {

	// Lowest possible number for a comic
    lowBound = generateCmd.Flags().Int("lowest", 1, "low bound for comic number")

    // The default highest number will actually be pulled during Run field
    highBound = generateCmd.Flags().Int("highest", 0, "high bound for comic number")
    
	rootCmd.AddCommand(generateCmd)
}

func getNewestComicNum() int {
    // xkcd conveniently provides a JSON link with the latest comic number,
    // leaving me the task of parsing and returning it

    type JsonStruct struct {
        Num int
    }

    var jsonObj JsonStruct
    url := "https://xkcd.com/info.0.json"

    resp, err := http.Get(url)
    if(err != nil) { panic(err) }
    defer resp.Body.Close()

    jsonByte, _ := ioutil.ReadAll(resp.Body)
    json.Unmarshal(jsonByte, &jsonObj)

    return jsonObj.Num

}

func getComicUrl(min int, max int) string {
    
    // Random seed based on the current time, I guess
    rand.Seed(time.Now().UnixNano())

    // Some math for getting a random number between the low and high bounds
    comicNumStr := strconv.Itoa(rand.Intn(max - min + 1) + min)
    comicWebUrl = fmt.Sprintf("https://xkcd.com/%s", comicNumStr)
    
    return comicWebUrl

}

func getComicImgUrl(imgUrl string) string {

    resp, err := http.Get(comicWebUrl)
    if(err != nil) { panic(err) }
    defer resp.Body.Close()

    htmlDoc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil { panic(err) }
    
    htmlDoc.Find("#comic > img").Each(func(i int, s *goquery.Selection) {
        imgUrl, _ = s.Attr("src")
    })
    
    return imgUrl

}