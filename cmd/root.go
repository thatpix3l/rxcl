package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
    lowBound int
    highBound int

    rootCmd = &cobra.Command{
        Use: "xclg",
        Short: "Random xkcd comic link generator",
        Long: `
This application creates a valid, random xkcd comic link by
generating a random number in a range between 1 and the latest comic number,
or a custom range if provided.`,
        Run: func(cmd *cobra.Command, args []string) {
            
            highest := cmd.Flags().Lookup("highest")
            
            if(!highest.Changed) {
                highest.Value.Set(strconv.Itoa(getNewestNum()))

            }

            generateComicLink()

        },
    }
)

func Execute() {

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

}

func init() {
    
    // Lowest possible number for a comic
    rootCmd.Flags().IntVar(&lowBound, "lowest", 1, "low bound for comic number")

    // The default highest number will actually be pulled during the rootCmd's Run field
    rootCmd.Flags().IntVar(&highBound, "highest", 0, "high bound for comic number")
    
}

func getNewestNum() int {

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

func generateComicLink() {
    
    // Random seed based on the current time, I guess
    rand.Seed(time.Now().UnixNano())
    min := lowBound
    max := highBound

    // Some math for getting a random number between the low and high bounds
    comicNumStr := strconv.Itoa(rand.Intn(max - min + 1) + min)

    fmt.Printf("https://xkcd.com/%s\n", comicNumStr)

}
