module siteintel

go 1.22

require (
	github.com/PuerkitoBio/goquery v1.9.2
	github.com/spf13/cobra v1.8.1
)

replace github.com/spf13/cobra => ./third_party/cobra

replace github.com/PuerkitoBio/goquery => ./third_party/goquery
