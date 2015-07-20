package main

import "os"
import "fmt"
import "mapreduce"
import "container/list"
import "unicode"
import "strconv"
import "strings"

// our simplified version of MapReduce does not supply a
// key to the Map function, as in the paper; only a value,
// which is a part of the input file contents
func Map(value string) *list.List {
  items := list.New()
  splitFunc := func(c rune) bool{
    return !unicode.IsLetter(c) && !unicode.IsNumber(c)
  }
  array := strings.FieldsFunc(value, splitFunc)
  for _, x:= range array {
    items.PushFront(mapreduce.KeyValue{x,"1"});
  }
  return items
}

// iterate over list and add values
func Reduce(key string, values *list.List) string {
  var sum int =0
  for e:= values.Front(); e!=nil; e=e.Next(){
    tmp,_ := strconv.Atoi(e.Value.(string))
    sum += tmp;
  }
  return strconv.Itoa(sum)
}

// Can be run in 3 ways:
// 1) Sequential (e.g., go run wc.go master x.txt sequential)
// 2) Master (e.g., go run wc.go master x.txt localhost:7777)
// 3) Worker (e.g., go run wc.go worker localhost:7777 localhost:7778 &)
func main() {
  if len(os.Args) != 4 {
    fmt.Printf("%s: see usage comments in file\n", os.Args[0])
  } else if os.Args[1] == "master" {
    if os.Args[3] == "sequential" {
      mapreduce.RunSingle(5, 3, os.Args[2], Map, Reduce)
    } else {
      mr := mapreduce.MakeMapReduce(5, 3, os.Args[2], os.Args[3])    
      // Wait until MR is done
      <- mr.DoneChannel
    }
  } else {
    mapreduce.RunWorker(os.Args[2], os.Args[3], Map, Reduce, 100)
  }
}
