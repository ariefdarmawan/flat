package test

import (
    "github.com/ariefdarmawan/flat"
    "testing"
    "io"
    "github.com/eaciit/toolkit"
)

var fn string = "/Users/ariefdarmawan/Dropbox/biz/eaciit/melon/05. From Clients/DATA_POC_EACIIT/mu_user_stream.txt"
var f *flat.Flat

func check(t *testing.T, e error, pre string){
    if e!=nil {
        if pre=="" {
            t.Fatal(e.Error)
        } else {
            t.Fatalf("%s %s", pre, e.Error())
        }
    }
}

func skipIfNil(t *testing.T){
    if f==nil {
        t.Skip()
    }
}

func TestConnect(t *testing.T){
    f = flat.New(fn,true,true)
    f.Delimeter = '|'
    check(t, f.Open(), "connect")
}

func TestRead(t *testing.T){
    skipIfNil(t)
    check(t,f.Reset(),"Read.Reset")
    iseof := false
    i := 0
    for ;!iseof; {
        //data := toolkit.M{}
        i++
        data, eof := f.ReadM()
        if eof==nil{
            toolkit.Printf("Data: %d => %v \n", i, data)
        } else if eof==io.EOF{
            iseof=true
        } else {
            check(t,eof,"Read.Iterate")
        }
    }
}

func TestClose(t *testing.T){
    skipIfNil(t)
    f.Close()
}



