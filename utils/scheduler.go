package utils

import (
	"encoding/xml"
	"net"
	"strconv"

	"github.com/novriyantoAli/go-phc/helper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Response struct {
	XMLName xml.Name `xml:"GetAttLogResponse"`
    Row   []Row   `xml:"Row"`
}

type Row struct {
	XMLName xml.Name `xml:"Row"`
	PIN int `xml:"PIN"`
	DateTime helper.CustomTime `xml:"DateTime"`
	Status int `xml:"Status"`
}

func getData(IPAddress string, key string, ch chan Response) {
	conn, err := net.Dial("tcp", IPAddress + ":80")
	if err != nil {
		logrus.Error(err)
		return
	}
	defer conn.Close()

	soapRequest := "<GetAttLog><ArgComKey xsi:type=\"xsd:integer\">"+key+"</ArgComKey><Arg><PIN xsi:type=\"xsd:integer\">All</PIN></Arg></GetAttLog>"

	newLine := "\r\n"

	fmt.Fprintf(conn, "POST /iWsService HTTP/1.0" + newLine)
	fmt.Fprintf(conn, "Content-Type: text/xml" + newLine)
	fmt.Fprintf(conn, "Content-Length: " + strconv.Itoa(len(soapRequest)) +newLine+newLine)
	fmt.Fprintf(conn, soapRequest + newLine)

    var buf bytes.Buffer
    io.Copy(&buf, conn)
}

// RunScheduler is a function to schedule finger print machine
func RunScheduler(){
	channel := make(chan Response)
	for {
		go getData(viper.GetString("finger.ip1"), viper.GetString("finger.key"), channel)
		
	}
}