package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"

	_ "github.com/go-sql-driver/mysql"

	_usersHandler "github.com/novriyantoAli/go-insinyur-radius-v1/users/delivery/http"
	_usersRepository "github.com/novriyantoAli/go-insinyur-radius-v1/users/repository/mysql"
	_usersUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/users/usecase"

	_resellerHandler "github.com/novriyantoAli/go-insinyur-radius-v1/reseller/delivery/http"
	_resellerRepository "github.com/novriyantoAli/go-insinyur-radius-v1/reseller/repository/mysql"
	_resellerUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/reseller/usecase"

	_transactionHandler "github.com/novriyantoAli/go-insinyur-radius-v1/transaction/delivery/http"
	_transactionRepository "github.com/novriyantoAli/go-insinyur-radius-v1/transaction/repository/mysql"
	_transactionUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/transaction/usecase"

	_radcheckRepository "github.com/novriyantoAli/go-insinyur-radius-v1/radcheck/repository/mysql"

	_packageHandler "github.com/novriyantoAli/go-insinyur-radius-v1/package/delivery/http"
	_packageRepository "github.com/novriyantoAli/go-insinyur-radius-v1/package/repository/mysql"
	_packageUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/package/usecase"

	_schedulerHandler "github.com/novriyantoAli/go-insinyur-radius-v1/scheduler/delivery/udp"
	_schedulerUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/scheduler/usecase"

	_radacctRepository "github.com/novriyantoAli/go-insinyur-radius-v1/radacct/repository/mysql"
)

type responseError struct {
	Message string `json:"error"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logrus.SetReportCaller(true)
	viper.SetConfigName("config")    // name of config file (without extension)
	viper.SetConfigType("yaml")      // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/ir/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.ir") // call multiple times to add many search paths
	viper.AddConfigPath(".")         // optionally look for config in the working directory
	err := viper.ReadInConfig()      // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if viper.GetBool(`debug`) {
		logrus.Infoln("SERVICE RUN IN DEBUG MODE")
	}

	// batas
	dbConn := createDB()
	triggerName := viper.GetString(`administrator.triggerName`)
	// create all event after insert
	_, err = dbConn.Exec(`
	CREATE TRIGGER ` + triggerName + ` AFTER INSERT ON radacct FOR EACH ROW 
		
	BEGIN
		
	SET @expiration = (SELECT COUNT(*) FROM radcheck WHERE username = New.username AND attribute = 'Expiration'); 
		
	IF (@expiration = 0) THEN
		SET @validity_value = (SELECT package.validity_value FROM radpackage INNER JOIN package ON package.id = radpackage.id_package WHERE radpackage.username = New.username ORDER BY radpackage.id DESC LIMIT 1);
		SET @validity_unit = (SELECT package.validity_unit FROM radpackage INNER JOIN package ON package.id = radpackage.id_package WHERE radpackage.username = New.username ORDER BY radpackage.id DESC LIMIT 1);

		IF (@validity_unit = 'HOUR') THEN
			INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_value HOUR), "%d %b %Y %H:%I:%S"));

		ELSEIF (@validity_unit = 'DAY') THEN
			INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_value DAY), "%d %b %Y %H:%I:%S"));

		ELSEIF (@validity_unit = 'MONTH') THEN
			INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_value MONTH), "%d %b %Y %H:%I:%S"));

		ELSEIF (@validity_unit = 'YEAR') THEN
			INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_unit YEAR), "%d %b %Y %H:%I:%S"));

		END IF;

	END IF;
	END;`)

	if err != nil {
		logrus.Error(err)
	}

	dbConn.Close()
}

func createDB() *sql.DB {
	// set radacct to check if user logged in
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add(`parseTime`, "1")
	val.Add(`loc`, "Asia/Makassar")

	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return dbConn
}

func main() {
	// initialize
	f, err := os.OpenFile(`go-insinyur-radius-v1.log`, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer f.Close()

	wrt := io.MultiWriter(os.Stdout, f)

	logrus.SetOutput(wrt)

	// database initialize
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add(`parseTime`, "1")
	val.Add(`loc`, "Asia/Makassar")

	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil {
		logrus.Fatalln(err)
		// log.Fatal(err)
	}

	err = dbConn.Ping()
	if err != nil {
		logrus.Fatalln(err)
		// log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			logrus.Fatalln(err)
			// log.Fatal(err)
		}
	}()

	e := echo.New()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "X-CSRF-Token", "application/json"},
		Debug:          true,
	})
	e.Use(echo.WrapMiddleware(corsMiddleware.Handler))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Validator = &CustomValidator{validator: validator.New()}

	timeout := time.Duration(viper.GetInt("context.timeout")) * time.Second

	/**
	 * Defined Application Repository
	 */
	usersRepository := _usersRepository.NewMysqlRepository(dbConn)
	resellerRepository := _resellerRepository.NewMysqlRepository(dbConn)
	transactionRepository := _transactionRepository.NewMysqlRepository(dbConn)
	radcheckRepository := _radcheckRepository.NewMysqlRepository(dbConn)
	packageRepository := _packageRepository.NewMysqlRepository(dbConn)
	radacctRepository := _radacctRepository.NewMysqlRepository(dbConn)

	/**
	 * Defined Application Usecase
	 */
	usersUsecase := _usersUsecase.NewUsecase(timeout, usersRepository)
	resellerUsecase := _resellerUsecase.NewUsecase(timeout, resellerRepository, packageRepository, radcheckRepository, transactionRepository, radacctRepository)
	packageUsecase := _packageUsecase.NewUsecase(timeout, packageRepository)
	transactionUsecase := _transactionUsecase.NewUsecase(timeout, transactionRepository)
	schedulerUsecase := _schedulerUsecase.NewUsecase(timeout, radcheckRepository, radacctRepository)
	/**
	 * Call all Handler here
	 */
	_usersHandler.NewHandler(e, usersUsecase)
	_resellerHandler.NewHandler(e, resellerUsecase)
	_packageHandler.NewHandler(e, packageUsecase)
	_transactionHandler.NewHandler(e, transactionUsecase)
	_schedulerHandler.NewHandler(schedulerUsecase)
	/**
	 * Call Echo Framework function for run this app
	 */

	logrus.Fatal(e.Start(viper.GetString("server.address")))
}
