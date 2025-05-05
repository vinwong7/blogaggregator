module github.com/vinwong7/blogaggregator

go 1.23.4

require github.com/vinwong7/blogaggregator/internal/config v0.0.0

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
)

replace github.com/vinwong7/blogaggregator/internal/config => ./internal/config
