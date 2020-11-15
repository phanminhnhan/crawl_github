package main

import (
	"context"
	"fmt"
	"github-trending/database"
	"github-trending/crawler"
	"github-trending/handler"
	"github-trending/log"
	"github-trending/repository/repo_impl"
	"github-trending/router"
	"github.com/labstack/echo"
	"os"
	"time"
)


func init() {
	os.Setenv("APP_NAME", "github")
	log.InitLogger(false)
}

func main() {
	sql := &database.Sql{
		Host:     "localhost",
		Port:      5432,
		Username: "postgres",
		Password: "phanminhnhan",
		DbName:   "github",
	}
	sql.ConnectDB()
	defer sql.Close()
	var email string 
	err := sql.Db.GetContext(context.Background(), &email, "select email from users where email = $1", "a@admin.com ") 
	if err != nil {
		os.Setenv("APP_NAME", "github")
		log.InitLogger(false)
		log.Error(err)
	}
	fmt.Println("======================")


	
	//log.Error("Error happend")

	e := echo.New()
	
	userhandler := handler.UserHandler{UserRepo: repo_impl.NewRepo(sql),}

	repoHandler := handler.RepoHandler{
		GithubRepo: repo_impl.NewGithubRepo(sql),
	}
	api := router.API{
		Echo: e,
		UserHandler: userhandler,
		RepoHandler: repoHandler,
	}
	api.SetupRouter()

	go scheduleUpdateTrending(1000*time.Second, repoHandler)
	e.Logger.Fatal(e.Start(":8888"))
}




func scheduleUpdateTrending(timeSchedule time.Duration, handler handler.RepoHandler) {
	ticker := time.NewTicker(timeSchedule)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Checking from github...")
				crawler.CrawlRepo(handler.GithubRepo)
			}
		}
	}()
}
