package shallalist

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"github.com/go-sql-driver/mysql"
)

type category struct {
	name  string
	bytes *string
}

func loadFile(category string) (*string, error) {
	b, err := ioutil.ReadFile("data/BL/" + category + "/domains")
	if err != nil {
		return nil, err
	}
	b2, err := ioutil.ReadFile("data/BL/" + category + "/urls")
	result := string(b) + string(b2)
	return &result, err
}

func Run() {
	categories := []string{"aggressive", "alcohol", "anonvpn", "drugs", "dating", "downloads", "dynamic", "education/schools",
		"finance/banking", "finance/insurance", "finance/moneylending", "finance/other", "finance/realestate", "finance/trading",
		"fortunetelling", "forum", "gamble", "government", "hacking", "hobby/cooking", "hobby/games-misc", "hobby/games-online",
		"hobby/gardening", "hobby/pets", "homestyle", "hospitals", "imagehosting","isp","jobsearch","library", "military", "models", "movies",
		"music", "news", "podcasts", "politics", "porn", "religion", "science/astronomy", "science/chemistry",  "sex/education", "sex/lingerie", 
		"searchengines", "violence", "spyware", "weapons", "webmail",
	}
	cFiles := make([]category, len(categories))
	for i, cat := range categories {
		cFiles[i].name = cat
		var err error
		cFiles[i].bytes, err = loadFile(cat)
		if err != nil {
			log.Fatal(err)
		}
	}

	//myHosts := []string{"znamost.cz", "zkhnord.de"}
	i := 800000
	db := dbConn()
	defer db.Close()
	db.SetMaxIdleConns(100)
	for i < 4000000 {
		myHosts, err := selectHosts(db, i, 10000)
		if err != nil {
			log.Fatal(err)
		}
		for _, cat := range cFiles {
			for _, host := range *myHosts {
				if strings.Contains(*cat.bytes, host) {
					fmt.Printf("Match found for category %s and url %s\n", cat.name, host)
					err := insertCategory(db, host, cat.name, 3000)
					if err != nil {
						me, ok := err.(*mysql.MySQLError)
						if !ok {
							log.Fatal(err)
						}
						if me.Number == 1062 {
							fmt.Printf("Entry already in the DB!\n")
						}
					}
					continue
				}
			}
		}
		i += 10000
	}

}
