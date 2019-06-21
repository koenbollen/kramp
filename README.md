# Kramp

[![Go Report Card](https://goreportcard.com/badge/github.com/koenbollen/kramp)](https://goreportcard.com/report/github.com/koenbollen/kramp)

## About

This project was created as an assignment for KrampHub.nl

### Assignment

> Using your favorite GO framework / libraries build a service, that
> will accept a request with text parameter on input.
> It will return maximum of 5 books and maximum of 5 albums that are
> related to the input term. The response elements will only contain
> title, authors(/artists) and information whether it's a book or an
> album. Sort the result by title alphabetically.
> 
> - For albums please use the [iTunes API](https://affiliate.itunes.apple.com/resources/documentation/itunes-store-web-service-search-api/#searching)
> - For books please use [Google Books API](https://developers.google.com/books/docs/v1/reference/volumes/list)
> 
> Make sure the software is production-ready from resilience, stability
> nd performance point of view.
> The stability of the downstream service may not be affected by the
> stability of the upstream services.
> Results originating from one upstream service (and its stability /
> performance) may not affect the results originating from the other
> upstream service.

## Local Development

This project is created using [Go](https://golang.org/). This repository followw the 
[scripts-to-rule-them-all](https://github.com/github/scripts-to-rule-them-all) method.

To get all dependencies up and running just run:

```bash
$ script/bootstrap
```

Then start the server and optionally the client a separate terminal:

```bash
$ script/server
```

```bash
$ script/client
Server running at http://localhost:1234
# or:
$ curl -s "http://localhost:8080/?q=Golang" | jq .
```

## Choices

- Since there is only one simple endpoint the default http library will suffice.
- Using zap.Logger since it's the easiest, quickest structured logging library.
- Used the [warmed](https://github.com/koenbollen/warmed) experiment since this 
  service will hardly be called in production. In a normal production environment 
  the normal http.Client would be fine (and the better choice).
- Created the sources pattern since it's really flexible and it neatly fit this 
  usecase.
- Did not implement caching since this services would be behind a cache 
  normally (which reduces complexity in this project).
