package main

import (
	"bytes"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/freshman-tech/news-demo/news"
	"github.com/joho/godotenv"

	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
)

var tpl = template.Must(template.ParseFiles("index.html"))

type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
}

func (s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}

func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}

	return s.NextPage - 1
}

func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ldKey := os.Getenv("LD_SDK_KEY")
		if ldKey == "" {
			log.Fatal("Env: ldKey must be set")
		}

		ldClient, _ := ld.MakeClient(ldKey, 5*time.Second)
		if ldClient.Initialized() {
			log.Println("Launch Darkly SDK successfully initialized!")
		} else {
			log.Fatal("Launch Darkly SDK failed to initialize")
			os.Exit(1)
		}

		flagValue, err := ldClient.BoolVariation("maxPageItems10", getUser(r), false)
		if err != nil {
			log.Fatal("error: " + err.Error())
		}
		log.Println("Current Page size is: ", newsapi.PageSize)

		if flagValue {
			newsapi.PageSize = 10
			log.Println("Feature flag is turned ON")
		} else {
			log.Println("Feature flag is OFF")
			newsapi.PageSize = default_item_count
		}
		log.Println("Page size after flag evaluation is: ", newsapi.PageSize)

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results:    results,
		}

		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
	}
}

const default_item_count = 50

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, apiKey, default_item_count)

	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/search", searchHandler(newsapi))
	mux.HandleFunc("/", indexHandler)
	log.Println("Open a browser and navigate to loacalhost:", port)
	http.ListenAndServe(":"+port, mux)
}

func getUserAgent(r *http.Request) string {
	ua := r.UserAgent()
	return ua
}

func getUser(r *http.Request) lduser.User {

	var userKey = "non-chrome-user"
	ua := getUserAgent(r)
	log.Println("user agent is: ", ua)
	if strings.Contains(ua, "Chrome") {
		userKey = "chrome-users"
	}
	log.Println("user-key: ", userKey)
	user := lduser.NewUserBuilder(userKey).
		Build()

	return user
}
