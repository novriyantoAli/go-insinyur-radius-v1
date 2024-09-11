package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	// I use TUN interface, so only plain IP packet, no ethernet header + mtu is set to 1300
	BUFFERSIZE = 1500
	MTU        = "1400"
	PORT       = "8000"
	TUN_ADDR   = "10.0.0.1"
	REMOTE     = "10.0.0.2"
)

func runIP(args ...string) (err error) {
	cmd := exec.Command("/sbin/ip", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if nil != err {
		log.Println("Error running /sbin/ip:", err)
	}
	return
}

func runIFConfig(args ...string) (err error) {
	cmd := exec.Command("/sbin/ifconfig", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if nil != err {
		log.Println("Error running /sbin/ifconfig: ", err)
	}
	return
}

func runServer(args ...string) (err error) {
	cmd := exec.Command("./server", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		log.Println("Error running ./server", err)
		return
	}
	// just run schedule
	for {
		// get data from secret

		// if expired
		// cmd.Process.Kill()
		time.Sleep(2 * time.Minute)
	}
}

func RunTunnel(secret string) {
	var tunFace string
	for {
		i := 0
		tunFace = fmt.Sprintf("tun%d", i)
		err := runIP("tuntap", "add", "dev", tunFace, "mode", "tun")
		if err == nil {
			runIFConfig(tunFace, TUN_ADDR, "dstaddr", REMOTE, "up")
			break
		}
		i++
	}
	runServer(tunFace, PORT, secret, "-m", MTU, "-a", REMOTE, "32", "-d", "8.8.8.8", "-r", "0.0.0.0", "0")
}
