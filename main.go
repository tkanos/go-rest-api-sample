package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/tkanos/go-rest-api-sample/account"
	"github.com/tkanos/go-rest-api-sample/config"
	"github.com/tkanos/go-rest-api-sample/mongoDb"
	"gopkg.in/mgo.v2"
)

var (
	appConfig   *config.Config
	infoLogger  = log.NewJSONLogger(os.Stdout)
	errorLogger = log.NewJSONLogger(os.Stderr)
)

// exit error codes
var (
	dbError = -2
)

func init() {
	appConfig = config.GetConfig()
}

func main() {

	//Db Connection
	session, err := mgo.Dial(appConfig.MongoConnectionString)
	if err != nil {
		errorLogger.Log("mongo_session_error", err)
		os.Exit(dbError)
	}
	defer session.Close()

	// Endpoints
	accountEndpoints := getAccountEndpoints(session)

	// Errors channel
	errc := make(chan error)

	// Interrupt handler.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP Transport
	go func() {
		httpAddr := ":" + strconv.Itoa(appConfig.Port)
		mux := http.NewServeMux()

		mux.Handle("/accounts/", account.MakeHTTPHandler(errorLogger, accountEndpoints))

		mux.HandleFunc("/healthz", healthzHandler)

		http.Handle("/", mux)
		infoLogger.Log("service", "go-rest-api-sample", "transport", "http", "address", httpAddr, "msg", "listening")
		errc <- http.ListenAndServe(httpAddr, nil)
	}()

	infoLogger.Log("exit", <-errc)
}

func getAccountEndpoints(mongoSession *mgo.Session) account.Endpoints {

	accountRepository, err := mongoDb.NewAccountRepository(mongoSession)
	if err != nil {
		errorLogger.Log("mongo_account_session_error", err)
		os.Exit(dbError)
	}
	accountService := account.NewService(accountRepository)

	getByIDEndpoint := account.MakeGetAccountEndpoint(accountService)

	getListEndpoint := account.MakeGetAccountsEndpoint(accountService)

	updateEndpoint := account.MakeUpdateAccountEndpoint(accountService)

	createEndpoint := account.MakeCreateAccountEndpoint(accountService)

	deleteEndpoint := account.MakeDeleteAccountEndpoint(accountService)

	return account.Endpoints{
		GetByID: getByIDEndpoint,
		GetList: getListEndpoint,
		Update:  updateEndpoint,
		Create:  createEndpoint,
		Delete:  deleteEndpoint,
	}
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
