package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "os/exec"
    "strings"
)

func main() {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
    	panic(err)
    }
    for _, addr := range addrs {
    	ipNet, err := addr.(*net.IPNet)
    	if !err {
    		panic(err)
    	}
    	if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
    		fmt.Println("listening on", ipNet.IP)
    	}
    }
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/output/") {
	    program := r.URL.Path[8:]
    	    if !checkIfAllowed(program) {
		fmt.Fprint(w, "not allowed!", "\n")
    		return
	    }
	    fmt.Fprint(w, program, "\n")
	    fmt.Fprint(w, runProgram(program, w, true))
	    return
	} else if r.URL.Path != "/" {
	    program := r.URL.Path[1:]
	    if !checkIfAllowed(program) {
		fmt.Fprint(w, "not allowed!", "\n")
		return
	    }
	    fmt.Fprint(w, runProgram(program, w, false))
	    return
	} else {
	   http.Error(w, "bozo", 401)
	}
    })
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func runProgram(program string, client http.ResponseWriter, output bool) string {
    fmt.Println(program)
    if output {
	out, err := exec.Command(program).Output()
	if err != nil {
	    fmt.Println(err)
	    }
	fmt.Fprintln(client, string(out))
	return "finnished"
    } else {
	cmd := exec.Command(program, "--detached")
	err := cmd.Start()
	if err != nil {
	    panic(err)
	}
	fmt.Fprint(client, program, "\n")
	cmd.Process.Release()
	return "started"
    }
}
func checkIfAllowed(program string) bool {
    file, err := os.Open("allowlist")
    if err != nil {
	panic(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
    	if program == scanner.Text() {
	    return true
	}
    }
    return false
}
