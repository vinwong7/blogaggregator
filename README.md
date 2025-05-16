# blogaggregator

6th project in Boot.dev - Blog Aggregator CLI
---------------------------------------------

Requires Postgres and Go in order to run the program.

To install the CLI, use "go install"

A config file has to be set up so that the program can access the Postgresql database. Create a .gatorconfig json with this below, replace db_url string with your own database connection string:

>{"db_url":"postgres://postgres:postgres@localhost:5432/gator?>sslmode=disable","current_user_name":"kahya"}

Some functions for gator:
- register "name": add new user to gator
- login "name": log into gator as named user
- addfeed "url": subscribe to feed, adding feed into datbase
- follow "url": follow feed already in database
- agg: run feed aggregation, adding posts of followed feeds to database
- browse "number": look at "number" of post from followed feeds, default 2

