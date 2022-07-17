# News Application Demo with Launch Darkly

This is a fork of the News Application that has been updated with feature management functionality provided by Launch Darkly.


Here's what the [completed application](https://freshman-news.herokuapp.com/)
looks like:

![demo](https://ik.imagekit.io/freshman/news-demo_MrYio9GKlzSi.png)

To use this web application you simply have to enter your search term and hit enter. Results of new itesm related to your search term are displayed in pages. By default 50 line items are displayed on a single page. Number of pages are calculated based on total results returned and number of news items allowed on a single page.

We are demonstrating a new feature here for Chrome users to have only 10 results per page. For all other users, default behavior of 50 pages stays intact. This implementation uses Feature Flags with Launch Darkly to make a decision how many items to be displayed on a single page. It also uses User Segmentation functionality of Launch Darkly to segment the users into Chrome and Non-Chrome users for implementation of targeting (Bonus point!)



## Prerequisites

- You need to have [Go](https://golang.org/dl/) installed on your computer. The
version used to test the code in this repository is **1.18.4**.

- Sign up for a [News API account](https://newsapi.org/register) and get your
free API key. You will need this key to retrieve the news items.

- You need to have a [Launch Darkly API key](https://app.launchdarkly.com/) and have to implement user segmentation based on the user keys. User keys used are `chrome-users` and `non-chrome-user`. I can demo how this is implemented in my account.

## Get started

- Clone this repository to your filesystem.

```bash
$ git clone https://github.com/Freshman-tech/news-demo
```

- Rename the `.env.example` file to `.env` and enter your News API Key.
- Also enter the Launch Darkly SDK Key in the `.env` file. Your completed file will look like this:<p>
```
PORT=3000
NEWS_API_KEY= your news api key here
LD_SDK_KEY= your launch darkly api key here
```
- `cd` into it and run the following command: `go build && ./news-demo` to start the server on port 3000.
- Visit http://localhost:3000 in your browser.
