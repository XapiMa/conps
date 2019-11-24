package main

import (
  "os"
  "fmt"
)

func main(){
  for i,arg := range os.Args{
    fmt.Printf("arg %d : %v\n", i, arg)
  }
}
