package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/D1abloRUS/checker-server/config"
	"github.com/D1abloRUS/checker-server/models"

	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
)

type postgres struct {
	User     string
	Password string
	Host     string
	DBname   string
	Port     int `default:"5432"`
}

//CheckErr func
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var p postgres

	err := envconfig.Process("postgres", &p)
	CheckErr(err)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.DBname)

	db, err := config.NewDB(psqlInfo)
	CheckErr(err)
	env := &config.Env{DB: db}

	err = models.CreateTables(env.DB)
	CheckErr(err)
	log.Print("checker-server: started normaly")

	r := httprouter.New()
	r.POST("/api/v1/activate", models.Activate(env))
	r.GET("/api/v1/gettask/:id", models.GetTask(env))
	r.POST("/api/v1/statusupdate", models.StatusUpdate(env))

	http.ListenAndServe(":3000", r)

}
