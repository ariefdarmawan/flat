package test

import (
    "testing"
    "os"
)

var fn string = "/users/ariefdarmawan/temp/dewa.txt"
var f *os.File

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
    var e error
    f, e = os.Open(fn)
    check(t, e, "connect")
}

func TestClose(t *testing.T){
    skipIfNil(t)
    f.Close()
}



