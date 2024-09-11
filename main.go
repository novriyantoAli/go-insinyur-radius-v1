package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	// "gopkg.in/go-playground/validator.v9"
	"github.com/go-playground/validator/v10"

	"github.com/novriyantoAli/go-insinyur-radius-v1/initializers"

	_usersHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_users/delivery/http"
	_usersRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_users/repository/mysql"
	_usersUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_users/usecase"

	_radcheckRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_radcheck/repository/mysql"

	_radreplyRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_radreply/repository/mysql"

	_packageHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_package/delivery/http"
	_packageRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_package/repository/mysql"
	_packageUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_package/usecase"

	_radacctRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_radacct/repository/mysql"

	_profilesHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_profiles/delivery/http"
	_profilesRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_profiles/repository/mysql"
	_profilesUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_profiles/usecase"

	_vouchersHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_vouchers/delivery/http"
	_vouchersRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_vouchers/repository/mysql"
	_vouchersUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_vouchers/usecase"

	_clientsHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_clients/delivery/http"
	_clientsRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_clients/repository/mysql"
	_clientsUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_clients/usecase"

	_orderHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_order/delivery/http"
	_orderRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_order/repository/mysql"
	_orderUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_order/usecase"

	_customerHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_customer/delivery/http"
	_customerRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_customer/repository/mysql"
	_customerUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_customer/usecase"

	_vcrRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_vouchers/repository/mysql"

	_radgroupHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_radgroup/delivery/http"
	_radgroupRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_radgroup/repository/mysql"
	_radgroupUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_radgroup/usecase"

	_firewallRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_firewall/repository/ros"
	_ipbindingRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_ipbinding/repository/ros"
	_simpleQueueRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_simplequeue/repository/ros"

	_schedulerHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_scheduler/delivery/udp"
	_schedulerUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_scheduler/usecase"

	_paymentHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_payment/delivery/http"
	_paymentRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_payment/repository/mysql"
	_paymentUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_payment/usecase"

	_dashboardHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_dashboard/delivery/http"
	_dashboardUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_dashboard/usecase"

	_timelineRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_timeline/repository/mysql"

	_messageHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_message/delivery/mbroker"
	_messageUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_message/usecase"

	_nasHandler "github.com/novriyantoAli/go-insinyur-radius-v1/m_nas/delivery/http"
	_nasRepository "github.com/novriyantoAli/go-insinyur-radius-v1/m_nas/repository/mysql"
	_nasUsecase "github.com/novriyantoAli/go-insinyur-radius-v1/m_nas/usecase"
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

	// meload config dengan viper
	initializers.LoadConfig()
	// mengkoneksikan ke database
	initializers.ConnectDB()
	// load library casbin
	// initializers.LoadAAA(initializers.GORM)
	// migrasi database
	initializers.Migrations()
	// koneksi ke mikrotik
	initializers.ConnectROS()
	// koneksi ke redis
	initializers.ConnectRedis()

	triggerName := viper.GetString(`administrator.triggerName`)
	// create all event after insert
	_ = initializers.GORM.Exec(`
	CREATE TRIGGER ` + triggerName + ` AFTER INSERT ON radacct FOR EACH ROW

	BEGIN

	SET @expiration = (SELECT COUNT(*) FROM radcheck WHERE username = New.username AND attribute = 'Expiration');

	IF (@expiration = 0) THEN
		SET @validity_value = (SELECT pkgs.validity_value FROM vcrs INNER JOIN pkgs ON pkgs.id = vcrs.pkg WHERE vcrs.username = New.username ORDER BY vcrs.id DESC LIMIT 1);
		SET @validity_value = (SELECT pkgs.validity_value FROM vcrs INNER JOIN pkgs ON pkgs.id = vcrs.pkg WHERE vcrs.username = New.username ORDER BY vcrs.id DESC LIMIT 1);
		SET @validity_unit = (SELECT pkgs.validity_unit FROM vcrs INNER JOIN pkgs ON pkgs.id = vcrs.pkg WHERE vcrs.username = New.username ORDER BY vcrs.id DESC LIMIT 1);

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
}

type User struct {
	Username string `json:"username"`
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

	e := echo.New()

	e.Static("/public/upload/img", "public/upload/img")

	e.Use(middleware.CORS())

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Validator = &CustomValidator{validator: validator.New()}

	timeout := time.Duration(viper.GetInt("context.timeout")) * time.Second

	/**
	 * Defined Application Repository
	 */
	usersRepository := _usersRepository.NewMysqlRepository(initializers.DB, initializers.GORM)
	radcheckRepository := _radcheckRepository.NewMysqlRepository(initializers.DB, initializers.GORM)
	radreplyRepository := _radreplyRepository.NewMysqlRepository(initializers.GORM)
	packageRepository := _packageRepository.NewMysqlRepository(initializers.DB, initializers.GORM)
	radacctRepository := _radacctRepository.NewMysqlRepository(initializers.DB, initializers.GORM)
	profilesRepository := _profilesRepository.NewMysqlRepository(initializers.GORM)
	vouchersRepository := _vouchersRepository.NewMysqlRepository(initializers.GORM)
	clientsRepository := _clientsRepository.NewMysqlRepository(initializers.GORM)
	orderRepository := _orderRepository.NewMysqlRepository(initializers.GORM)
	customerRepository := _customerRepository.NewMysqlRepository(initializers.GORM)
	vcrRepository := _vcrRepository.NewMysqlRepository(initializers.GORM)
	radgroupRepository := _radgroupRepository.NewMysqlRepository(initializers.GORM)
	paymentRepo := _paymentRepository.NewMysqlRepository(initializers.GORM)
	ipbindingRepository := _ipbindingRepository.NewROSRepository(initializers.ROS[0])
	firewallRepository := _firewallRepository.NewROSRepository(initializers.ROS[0])
	simpleQueueRepository := _simpleQueueRepository.NewROSRepository(initializers.ROS[0])
	timelineRepository := _timelineRepository.NewMysqlRepository(initializers.GORM)
	nasRepository := _nasRepository.NewMysqlRepository(initializers.GORM)
	/**
	 * Defined Application Usecase
	 */
	usersUsecase := _usersUsecase.NewUsecase(timeout, usersRepository)
	packageUsecase := _packageUsecase.NewUsecase(timeout, packageRepository)
	profilesUsecase := _profilesUsecase.NewUsecase(timeout, profilesRepository)
	vouchersUsecase := _vouchersUsecase.NewUsecase(timeout, vouchersRepository, radcheckRepository, packageRepository)
	clientsUsecase := _clientsUsecase.NewUsecase(clientsRepository, packageRepository)
	orderUsecase := _orderUsecase.NewUsecase(
		timeout,
		orderRepository,
		packageRepository,
		customerRepository,
		vcrRepository,
		radcheckRepository,
		radreplyRepository,
		clientsRepository,
		ipbindingRepository,
		firewallRepository,
		simpleQueueRepository,
		initializers.ROS_CLIENT,
	)
	customerUsecase := _customerUsecase.NewUsecase(timeout, customerRepository)
	radgroupUsecase := _radgroupUsecase.NewUsecase(radgroupRepository)
	schedulerUsecase := _schedulerUsecase.NewUsecase(
		timeout,
		radcheckRepository,
		radacctRepository,
		orderRepository,
		ipbindingRepository,
		firewallRepository,
		simpleQueueRepository,
		initializers.ROS_CLIENT,
	)
	paymentUcase := _paymentUsecase.NewUsecase(timeout, paymentRepo)
	dashboardUsecase := _dashboardUsecase.NewUsecase(timeout,
		customerRepository,
		clientsRepository,
		usersRepository,
		packageRepository,
		timelineRepository,
		radacctRepository,
	)
	messageUsecase := _messageUsecase.NewUsecase(timeout, vcrRepository)
	nasUsecase := _nasUsecase.NewUsecase(nasRepository, timeout)
	/**
	 * Call all Handler here
	 */
	_usersHandler.NewHandler(e, usersUsecase)
	_packageHandler.NewHandler(e, packageUsecase)
	_profilesHandler.NewHandler(e, profilesUsecase)
	_vouchersHandler.NewHandler(e, vouchersUsecase)
	_clientsHandler.NewHandler(e, clientsUsecase)
	_orderHandler.NewHandler(e, orderUsecase)
	_customerHandler.NewHandler(e, customerUsecase)
	_radgroupHandler.NewHandler(e, radgroupUsecase)
	_schedulerHandler.NewHandler(schedulerUsecase)
	_paymentHandler.NewHandler(e, paymentUcase)
	_dashboardHandler.NewHandler(e, dashboardUsecase)
	_messageHandler.NewHandler(initializers.REDIS, messageUsecase)
	_nasHandler.NewHandler(e, nasUsecase)
	/**
	 * Call Echo Framework function for run this app
	 */

	logrus.Fatal(e.Start(viper.GetString("server.address")))
}
