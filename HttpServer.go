package main

/*
#include "MyAsymmetricEncryption.h"
*/
import "C"
import (
 "fmt"
 "io"
 "os"
 "net/http"
)


const SO_FILE_PATH = "/home/lee/sofile"


var IMEI_INFO map[string]string 
var IMEI_STATE map[string]bool 


func handleLogin (w http.ResponseWriter, request *http.Request) {
    if request.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        _, _ = w.Write([]byte("Only GET Method is allowed!"))
        return
    }
    
    imei := request.FormValue("imei") 
    if imei == "" {
        w.WriteHeader(http.StatusBadRequest)
        _, _ = io.WriteString(w, "IMEI Empty!")
        return
    }
    
    //w.WriteHeader(http.StatusBadRequest)
    
    fmt.Println("token: " + imei)
    token, ok := IMEI_INFO [ imei ]
    if (ok) {
        _, _ = io.WriteString(w, token)
        IMEI_STATE [ token ] = true
    } else {
        w.WriteHeader(http.StatusBadRequest)
        _, _ = io.WriteString(w, "Login Error!")
    }
    
    fmt.Println("LOGIN SUCCESSED!") 
    return
   
}


func handleDownload (w http.ResponseWriter, request *http.Request) {
    if request.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        _, _ = w.Write([]byte("Only GET Method is allowed!"))
        return
    }
    
    token := request.FormValue("token")
    if token == "" {
        w.WriteHeader(http.StatusBadRequest)
        _, _ = io.WriteString(w, "TOKEN Empty!")
        return
    }
    
    if false == IMEI_STATE [ token ] {
        w.WriteHeader(http.StatusBadRequest)
        _, _ = io.WriteString(w, "Token Time Out!")
        return
    }
    
    filename := request.FormValue("filename")
    if filename == "" {
        w.WriteHeader(http.StatusBadRequest)
        _, _ = io.WriteString(w, "FILENAME Empty!")
        return
    }
    
    fmt.Println("filename: " + SO_FILE_PATH + "/" + filename)
    file, err := os.Open(SO_FILE_PATH + "/" + filename)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, _ = io.WriteString(w, "FILENAME  Incorrect!")
        return
    }
    
    defer file.Close()      
    defer fmt.Println("DOWNLOAD SUCCESSED!")                                                 

    w.Header().Add("Content-type", "application/octet-stream")
    w.Header().Add("content-disposition", "attachment; filename=\""+filename+"\"")
    
    _, err = io.Copy(w, file)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, _ = io.WriteString(w, "Download Error!")
        return
    }
    
    IMEI_STATE [ token ] = false 
}


func main() {
    fmt.Println(C.count)
    C.show()
    fmt.Println("SERVER START")
    
    IMEI_INFO = make(map[string]string)
    IMEI_INFO [ "R52R30L7RLZ" ] = "F00F3CF4D06E9D4326591EB921ECDB173478B12687D8D3EF9BA38296CF2A7E6F"
    IMEI_INFO [ "351922102365095" ] = "862D158B38A0D98DD155C4391F7271312BF8F879CF333E2040E2230C4F33BA4C"
    IMEI_INFO [ "353240111116864" ] = "5BBCB7083B38743FBE625E0A19EB29400F6428B9FC8EC433B1595FA8F2028F2D"
    IMEI_INFO [ "4bbca50019cbb1b8" ] = "D48F375BEB17F9CADA05EB5B8C5FF724D2C087E61E079885992FF8D0E1A8D13F"
    IMEI_INFO [ "ee7c62a1ca8f9fd8" ] = "5AEF65CA9E81E40ECE9E735FD70CE74565B08FE8D37D224150A5070639B3F281"
    
    IMEI_STATE = make(map[string]bool) 
    IMEI_STATE [ "F00F3CF4D06E9D4326591EB921ECDB173478B12687D8D3EF9BA38296CF2A7E6F" ] = false
    IMEI_STATE [ "862D158B38A0D98DD155C4391F7271312BF8F879CF333E2040E2230C4F33BA4C" ] = false
    IMEI_STATE [ "5BBCB7083B38743FBE625E0A19EB29400F6428B9FC8EC433B1595FA8F2028F2D" ] = false
    IMEI_STATE [ "D48F375BEB17F9CADA05EB5B8C5FF724D2C087E61E079885992FF8D0E1A8D13F" ] = false
    IMEI_STATE [ "5AEF65CA9E81E40ECE9E735FD70CE74565B08FE8D37D224150A5070639B3F281" ] = false
    
    http.HandleFunc("/login", handleLogin)
    http.HandleFunc("/download", handleDownload)

    fmt.Println("LISTEN PORT = 20213")
    err := http.ListenAndServe(":20213", nil)
    if err != nil  {
        fmt.Println("SERVER RUN ERROR!")
    }
}


