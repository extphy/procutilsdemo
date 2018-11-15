package main

import (
   "flag"
   "time"
   "log"
   "os"

   "github.com/extphy/procutils"
)

func main() {

   var numproc = flag.Int("numproc", 1, "num of processes")
   var timeout = flag.Int("timeout", 5, "seconds before stop processes")
   var isWorker = flag.Bool("worker", false, "run worker process")

   flag.Parse()

   if *isWorker {
      log.Println("running worker")
   } else {
      log.Printf("running with params numproc=%d, timeout=%d", *numproc, *timeout)
   }

   if *isWorker {
      for {
         time.Sleep(3 * time.Second)
         log.Println("child is alive")
         // TODO use args to manage stderr scenarios
      }
      return
   }

   var procMan = procutils.NewProcessManager()
   for tag := 0; tag < *numproc; tag++ {
      procMan.RunProcess(tag, "procutilsdemo", os.Environ())
   }

mainloop:
   for {
      select {
      case msg := <-procMan.MsgQueue:
         switch msg.MsgType {
         case procutils.MsgRunning:
            log.Printf("running proc %d", msg.Tag)
         }
      case <-time.After(time.Duration(*timeout) * time.Second):
         log.Println("stopped by timeout")
         for tag := 0; tag < *numproc; tag++ {
            procMan.StopProcess(tag)
         }
         break mainloop
      }
   }

}

