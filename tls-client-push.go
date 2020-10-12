package main

import (
    "crypto/tls"
    "crypto/x509"
    "io"
    "log"
    "os"
    "strconv"
)

func main(){
    cert, err := tls.LoadX509KeyPair("./dev1/dev1.crt", "./dev1/dev1.key")
    if err != nil {
        log.Fatalf("server: loadkeys: %s", err)
    }
    config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
    conn, err := tls.Dial("tcp4", "129.254.170.216:22222", &config)
    if err != nil {
        log.Fatalf("client: dial: %s", err)
    }
    defer conn.Close()
    log.Println("client: connected to: ", conn.RemoteAddr())

    state := conn.ConnectionState()
    for _, v := range state.PeerCertificates {
        x509.MarshalPKIXPublicKey(v.PublicKey)
    }
    log.Println("client: handshake: ", state.HandshakeComplete)
    log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
    
    buf := make([]byte, 20)

    var flag string = "1" 
    _, err = io.WriteString(conn,flag)
     if err != nil {
	    log.Fatalf("client: write: %s", err)
    }
    log.Printf("client: conn: write: %s", flag)

    _, err = conn.Read(buf)
    if err != nil {
            log.Printf("server: conn: read: %s", err)
            _, err = io.WriteString(conn,"server read error")
       	    os.Exit(1)
    }
    log.Printf("client: conn: read: %s", buf)

    var name string = os.Args[1]
    _, err = io.WriteString(conn,name)
     if err != nil {
	    log.Fatalf("client: write: %s", err)
    }
    log.Printf("client: conn: write: %s", name)

    _, err = conn.Read(buf)
    if err != nil {
            log.Printf("server: conn: read: %s", err)
            _, err = io.WriteString(conn,"server read error")
       	    os.Exit(1)
    }
    log.Printf("client: conn: read: %s", buf)

    message, num := filer(name+"-sign.gob")
    size := strconv.FormatInt(num,10) 
    
    _, err = io.WriteString(conn,size)
    if err != nil {
	    log.Fatalf("client: write: %s", err)
    }
    log.Printf("client: conn: write: %s", size)
    
    _, err = conn.Read(buf)
    if err != nil {
            log.Printf("client: conn: read: %s", err)
            _, err = io.WriteString(conn,"server read error")
       	    os.Exit(1)
    }
    log.Printf("client: conn: read: %s", buf)

    _, err = conn.Write(message)
    if err != nil {
            _, _ = io.WriteString(conn, "Send sign data fail")
       	    os.Exit(1)
    }
    log.Printf("client: conn: send file: %s", name+"-sign.gob")

    result := make([]byte, 50)

    n, err := conn.Read(result)
    var str string = string(result[:n])
    log.Printf("client: conn: read: %s",str)

    if str == "verify fail" {
       log.Printf("%s",str)
       os.Exit(1)
    }else if str == "vulnerability detection alot. Upload fail" {
       log.Printf("%s",str)
       os.Exit(1)
    }else if str == "resign fail" {
       log.Printf("%s",str)
       os.Exit(1)
    }else if str == "OK"{
	log.Println("client: conn: closed")
    	os.Exit(0)
    }
}

func filer(path string) ([]byte, int64) {
    fd, _ := os.Open(path)
    fi, _  := fd.Stat()
    defer fd.Close()
    var num = fi.Size()
    var data = make([]byte, num)
    _,_ = fd.Read(data)
    return data, num
}

