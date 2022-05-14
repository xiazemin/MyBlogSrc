---
title: How to Manage Database Timeouts and Cancellations in Go
layout: post
category: golang
author: 夏泽民
---
https://www.alexedwards.net/blog/how-to-manage-database-timeouts-and-cancellations-in-go
One of the great features of Go is that it's possible to cancel database queries while they are still running via a context.Context instance (so long as cancellation is supported by your database driver).

On the face of it, using this functionality is quite straightforward (here's a basic example). But once you start digging into the details there's a lot a nuance and quite a few gotchas... especially if you are using this functionality in the context of a web application or API.

So in this post I want to explain how to cancel database queries in a web application, what behavioral quirks and edge cases it is important to be aware of, and try to provide answers to the questions that you might have when working through all this.
<!-- more -->
But first off, why would you want to cancel a database query? Two scenarios spring to mind:

When a query is taking a lot longer to complete than expected. If this happens, it suggests a problem — either with that particular query or your database or application more generally. In this scenario, you would probably want to cancel the query after a set period of time (so that resources are freed-up and the database connection is returned to the sql.DB connection pool for reuse), log an error for further investigation, and return a 500 Internal Server Error response to the client.

When a client goes away unexpectedly before the query completes. This could happen for a number of reasons, such as a user closing a browser tab or terminating a process. In this scenario, nothing has really gone 'wrong', but there is no client left to return a response to so you may as well cancel the query and free-up the resources.

Mimicking a long-running query
Let's start with the first scenario. To demonstrate this, I'll make a very basic web application with a handler that executes a SELECT pg_sleep(10) SQL query against a PostgreSQL database using the pq driver. The pg_sleep(10) function will make the query sleep for 10 seconds before returning, essentially mimicking a slow-running query.

package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    _ "github.com/lib/pq"
)

var db *sql.DB

func slowQuery() error {
    _, err := db.Exec("SELECT pg_sleep(10)")
    return err
}

func main() {
    var err error

    db, err = sql.Open("postgres", "postgres://user:pa$$word@localhost/example_db")
    if err != nil {
        log.Fatal(err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/", exampleHandler)

    log.Println("Listening...")
    err = http.ListenAndServe(":5000", mux)
    if err != nil {
        log.Fatal(err)
    }
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
    err := slowQuery()
    if err != nil {
        serverError(w, err)
        return
    }

    fmt.Fprintln(w, "OK")
}

func serverError(w http.ResponseWriter, err error) {
    log.Printf("ERROR: %s", err.Error())
    http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
}
If you were to run this code, then make a GET / request to the application you should find that the request hangs for 10 seconds before you finally get an "OK" response. Like so:

$ curl -i localhost:5000/
HTTP/1.1 200 OK
Date: Fri, 17 Apr 2020 07:46:40 GMT
Content-Length: 3
Content-Type: text/plain; charset=utf-8

OK
Note: The structure of the application code above is deliberately over-simplified. In a real project I would recommend using dependency injection to make the sql.DB connection pool and logger available to your handlers, instead of using global variables.
Adding a context timeout
OK, now that we've got some code that mimics a long-running query, let's enforce a timeout on the query so it is automatically canceled if it doesn't complete within 5 seconds.

To do this we need to:
Use the context.WithTimeout() function to create a context.Context instance with a 5-second timeout duration.
Execute the SQL query using the ExecContext() method, passing the context.Context instance as a parameter.
I'll demonstrate:

package main

import (
    "context" // New import
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "time" // New import

    _ "github.com/lib/pq"
)

var db *sql.DB

func slowQuery(ctx context.Context) error {
    // Create a new child context with a 5-second timeout, using the
    // provided ctx parameter as the parent.
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // Pass the child context (the one with the timeout) as the first
    // parameter to ExecContext().
    _, err := db.ExecContext(ctx, "SELECT pg_sleep(10)")
    return err
}

...

func exampleHandler(w http.ResponseWriter, r *http.Request) {
    // Pass the request context to slowQuery(), so it can be used as the 
    // parent context.
    err := slowQuery(r.Context())
    if err != nil {
        serverError(w, err)
        return
    }

    fmt.Fprintln(w, "OK")
}

...
There are a few things about this that I'd like to emphasize and explain:

Note that we pass r.Context() (the request context) to slowQuery() to use as the parent context. As we'll see in the next section, this is important because it means that any cancellation signal on the request context will be able to 'bubble down' to the context that we use in ExecContext().

The defer cancel() line is important because it ensures that the resources associated with our child context (the one with the timeout) will be released before the slowQuery() function returns. If we don't call cancel() it may cause a memory leak: the resources won't be released until either the parent r.Context() is canceled or the 5-second timeout is hit (whichever happens first).

The timeout countdown begins from the moment that the child context is created using context.WithTimeout(). If you want more control over this you could use the alternative context.WithDeadline() function, which allows you to set an explicit time.Time value for when the context should timeout instead.

OK, let's try this out. If you run the application again and make a GET / request, after a 5-second delay you should get a response like this:

$ curl -i localhost:5000/
HTTP/1.1 500 Internal Server Error
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 17 Apr 2020 08:21:14 GMT
Content-Length: 28

Sorry, something went wrong
And if you go back to the terminal window running the application you should see a log message similar to this:

$ go run .
2020/04/17 10:21:07 Listening...
2020/04/17 10:21:14 ERROR: pq: canceling statement due to user request
That log message might seem a bit odd... until you realize that the error message is actually coming from PostgreSQL. In that light it makes sense: our web application is the user and we're canceling the query after 5 seconds.

So this is actually really good; things are working as we want.

Specifically, after 5 seconds the context timeout is reached and the pq driver sends a cancellation signal to PostgreSQL†. PostgreSQL then terminates the running query (thereby freeing-up resources). The client is sent a 500 Internal Server Error response, and the error message is logged so we know that something has gone wrong.

† More precisely, our child context (the one with the 5-second timeout) has a Done channel, and when the timeout is reached it will close the Done channel. While the SQL query is running, our database driver pq is also running a background goroutine which listens on this Done channel. If the channel is closed, then it sends a cancellation signal to PostgreSQL. PostgreSQL terminates the query, and then sends the error message that we see above as a response to the original pq goroutine. That error message is then returned to our slowQuery() function.

Dealing with closed connections
OK, let's try one more thing. Let's use curl to make a GET / request and then very quickly (within 5 seconds) press Ctrl+C to cancel the request.

If you look at the logs for the application again, you should see another log line with exactly the same error message that we saw before.

$ go run .
2020/04/17 10:21:07 Listening...
2020/04/17 10:21:14 ERROR: pq: canceling statement due to user request
2020/04/17 10:41:18 ERROR: pq: canceling statement due to user request
So what's happening here?

In this case, the request context (which we use as the parent in our code above) is canceled because the client closed the connection. From the net/http docs:

For incoming server requests, the [request] context is canceled when the client's connection closes, the request is canceled (with HTTP/2), or when the ServeHTTP method returns.

This cancellation signal bubbles down to our child context, it's Done channel is closed, and the pq driver terminates the running query in exactly the same way as before.

With that in mind, it's not surprising that we see the same error message... From a PostgreSQL point of view exactly the same thing is happening as when the timeout was reached.

But from the perspective of our web application the scenario is very different. A client connection being closed can happen for many different, innocuous, reasons. It's not really an error from our application's point of view, although it is probably sensible to log it as a warning (if we start to see elevated rates, it could be a sign that something is wrong).

Fortunately, it's possible to tell these two scenarios apart by calling the ctx.Err() method on our child context. If the context was canceled (due to a client closing the connection), then ctx.Err() will return context.Canceled. If the timeout was reached, then it will return context.DeadlineExceeded. If both the deadline is reached and the context is canceled, then ctx.Err() will surface whichever happened first.

There's another important thing to point out here: it's possible that a timeout/cancellation will happen before the PostgreSQL query even starts. For example you might have set MaxOpenConns() on your sql.DB connection pool, and if that open connection limit is reached and all connections are in-use, then the query will be 'queued' by sql.DB until a connection becomes available. In this scenario — or any other which causes a delay — it's quite possible that the timeout/cancellation will occur before a free database connection even becomes available. In this case ExecContext() will directly return the ctx.Err() value as the error response (instead of the "pq: canceling statement due to user request" error that we see above).

If you're using the QueryContext() method then it's also possible that the timeout/cancellation will occur when processing the data with Scan(). If this happens, then Scan() will directly return the ctx.Err() value as an error. As far as I can see this behavior isn't mentioned in the database/sql docs, but I can confirm that this is the case with Go 1.14 and the comments on issue #28842 suggest that it is intentional.

Putting all that together, a sensible approach is to check for the error "pq: canceling statement due to user request" and then wrap this with the error from ctx.Err() before returning from our slowQuery() function.

Then in our handler, we can use the errors.Is() function to check if the error from slowQuery() is equal to (or wraps) context.Canceled and manage it accordingly. Like so:

package main

import (
	"context"
	"database/sql"
	"errors" // New import
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func slowQuery(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, "SELECT pg_sleep(10)")
    // If we get a "pq: canceling statement..." error wrap it with the 
    // context error before returning.
    if err != nil && err.Error() == "pq: canceling statement due to user request" {
		return fmt.Errorf("%w: %v", ctx.Err(), err)
	}

	return err
}

...

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	err := slowQuery(r.Context())
	if err != nil {
        // Check if the returned error equals or wraps context.Canceled and 
        // record a warning if it does.
		switch {
		case errors.Is(err, context.Canceled):
			serverWarning(err)
		default:
			serverError(w, err)
		}
		return
	}

	fmt.Fprintln(w, "OK")
}

func serverWarning(err error) {
	log.Printf("WARNING: %s", err.Error())
}

...
If you were to run this application again now and make two different GET / requests — one that times out, and the other that you cancel — you should see clearly different messages in the application log, like so:

$ go run .
2020/04/17 13:09:25 Listening...
2020/04/17 13:09:45 ERROR: context deadline exceeded: pq: canceling statement due to user request
2020/04/17 13:09:47 WARNING: context canceled: pq: canceling statement due to user request
Other context-aware methods
The database/sql package provides context-aware variants for most actions on sql.DB, including PingContext(), QueryContext(), and QueryRowContext(). We can (and should!) update the main() function in the code above to use PingContext() instead of Ping().

In this case there is no request context to use as the parent, so we need to create an empty parent context with context.Background() instead. Like so:

...

func main() {
	var err error

	db, err = sql.Open("postgres", "postgres://user:pa$$word@localhost/example_db")
	if err != nil {
		log.Fatal(err)
	}

    // Create a context with a 10-second timeout, using the empty 
    // context.Background() as the parent.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

    // Use this when testing the connection pool.
	if err = db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", exampleHandler)

	log.Println("Listening...")
	err = http.ListenAndServe(":5000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

...
Can I set a global timeout for all requests?
Sure, you could create and use some middleware on your routes which adds a timeout to the current request context, similar to this:

func setTimeout(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
        defer cancel()

        // This gives you a copy of the request with a the request context 
        // changed to the new context with the 5-second timeout created 
        // above.
        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    })
}
If you take this approach there are a few of things to be aware of:

The timeout starts from the moment the context is created, so any code running in your handlers before the database query counts towards the timeout.
If you have multiple queries being executed in a handler, then they all have to complete within that one time.
The timeout will continue to apply even if you derive a child context with a different timeout duration. So while you can enforce an earlier timeout in a child context, you can't make it longer.
What about http.TimeoutHandler?
Go provides a http.TimeoutHandler() middleware function which you can use to wrap your handlers or router/servemux. This works similar to the middleware above in the sense that it sets a timeout on the request context... so the warnings above also apply when using this.

However, http.TimeoutHandler() also sends the client a 503 Service Unavailable response and a HTML error message. So, if you're using this in your application, you shouldn't (or at least, you don't need to) send the client an error response yourself when encountering a context.DeadlineExceeded error.

How about transactions? How does context work in those?
The database/sql package provides a BeginTx() method which you can use to initiate a context-aware transaction. A code example can be seen here.

It's important to understand that the context you provide to BeginTx() applies to the whole transaction. In the event of a timeout/cancellation on the context, then the queries in the transaction will automatically be rolled-back.

It's perfectly fine to pass the same context as a parameter for all the queries in the transaction, in which case it ensures that they all (as a whole) complete before any timeout/cancellation . Alternatively, if you want per-query timeouts you can create different child contexts with different timeouts for each in the queries in the transaction. But you must derive these child contexts from the context you passed to BeginTX(). Otherwise there is a risk that the BeginTX() context timeout/cancellation occurs and the automatic rollback happens, but your code still may try to execute the query with a still-live context. If that happened you would receive the error "sql: transaction has already been committed or rolled back".

What about background processing?
When doing background-processing in a different goroutine, bear in mind that if a parent context is canceled, the cancellation signal 'bubbles down' to its children. And also bear in mind what I quoted earlier about request context cancellation:

For incoming server requests, the [request] context is canceled ... when the ServeHTTP method returns.

Combine those two things, and it means that if you use a context which is a child of the request context in the background-process, the background-process will get a cancellation signal when the HTTP response is sent for the initial request. If you don't want that to be the case (and you probably don't), then you should create a brand-new context for the background-process using context.Background() and copy over any values that you need... or just pass them as regular parameters instead.

If a context is canceled, can I be confident that it's due to a closed connection?
Yes — so long as it's within the main goroutine for the request, it's a child of the request context, and you haven't manually canceled it yourself yet using defer cancel(). Otherwise, no.

Is the behavior the same with other databases and drivers?
I'm not sure. I've only used these features extensively with PostgreSQL and the pq driver. I imagine that things will be roughly the same with other databases and drivers, but you'll need to check.

Anything else I should know?
Yep. This is a strange one and it's not officially documented yet, but if a client makes a request with a non-empty request body then closes the connection, the context won't be canceled until after you have read the request body.

This doesn’t apply to requests without a request body, where the cancellation signal will be received immediately.

You should also be aware of the WriteTimeout setting on your http.Server (if you have set one). Your context timeouts should always be shorter than your WriteTimeout value, otherwise the WriteTimeout will be hit first, the connection will be closed, and the client won’t get any response.
