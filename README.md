# Form3 Go client library

A simple go client library for the [Form3 REST APIs](https://api-docs.form3.tech/api.html#organisation-accounts)

## Example usage
```go
package main

import (
	"context"
	"log"

	"github.com/althink/form3"
	"github.com/althink/form3/accounts"
)

func main() {
	ctx := context.Background()
	orgID := "0de1f73f-8af2-4316-86f9-325ce9755cb6"

	f3, err := form3.NewClient()
	if err != nil {
		log.Fatal("Failed to create client", err)
	}

	_, err := f3.Accounts.Create(ctx, accounts.NewWithGenID(orgID, &accounts.Attributes{
		Country: "PL",
		Name:    []string{"John Smith"},
	}))
	if err != nil {
		log.Fatal("Failed to create account", err)
	}

	_, err = f3.Accounts.Fetch(ctx, acc.Data.ID)
	if err != nil {
		log.Fatal("Failed to fetch account", err)
	}

	err = f3.Accounts.Delete(ctx, acc.Data.ID, *acc.Data.Version)
	if err != nil {
		log.Fatal("Failed to delete account", err)
	}
}
```

## Testing
To run unit tests `go test ./...`

To run integration tests `docker-compose up`

## About author
Krzysztof Szczesniak
E-mail: sl0w0rm@gmail.com

I'm mostly Java developer but sometimes I write mocks and tools in go. Although I don't write production code in go because it's not allowed in my current company (only JVM languages are allowed). 

## Technical decisions
I've decided to keep the structure simple as possible. There's no pkg folder because this library is so far very small. I've decided to make a dedicated package for every resource type, so it's more extendable and maintenable (your API has a lot of resource types). When second resource type will be added some common code regarding http calls should be extracted to seperate type and used in every resource client.

I've decided to make it a simple service api because there are not many operations. I've forced users to pass `context.Context` becuse it may be useful for adding custom headers to HTTP requests, for tracing purposed for example. Users can use custom `http.Client` and set Transport with custom delegating RoundTripper that adds custom headers. Other option would be to make it more object oriented while every operation creates a customizable request object that have an operation that allowes to execute given call. Something like `f3.Accounts.Creation().Do()`. Those request objects could be altered with `.WithContext(ctx)` call like in `http.Request.WithContext(ctx)`.

I've decided to create custom error types to make it more obvious what errors can be returned in given call (so users don't have to base they logic on http status codes). Although I'm not sure if that's the most idiomatic go code.

I haven't implemented rate limiter although it can be easly done by users with custom RoundTripper or by checking status code from HttpStatusError.

There's over 80% of code coverage. All happy paths and all unhappy paths that have custom errors are coverd. Not covered part is mostly error messages printing and some rare cases like errors related to parsing json etc.