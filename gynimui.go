package main

import (
    "fmt" 
    "sync"
    "strconv"
)

const sendersCount int = 2
const printersCount int = 2
const maxNumbers int = 20
const EndProcess int = -1 //enumas, darbo pabaigai nurodyti
const Continue int = 0

func main() {
    var sendersWaitGroup sync.WaitGroup
    var receiverWaitGroup sync.WaitGroup
    var printersWaitGroup sync.WaitGroup

    var senderChennels [sendersCount]chan int
    senderChennels[0] = make(chan int)
    senderChennels[1] = make(chan int)

    var printerChannels [printersCount]chan int
    printerChannels[0] = make(chan int)
    printerChannels[1] = make(chan int)

    var statusChannels [sendersCount]chan int
    statusChannels[0] = make(chan int)
    statusChannels[1] = make(chan int)

    sendersWaitGroup.Add(2)
    go sender(0, senderChennels[0], statusChannels[0], &sendersWaitGroup)
    go sender(11, senderChennels[1], statusChannels[1], &sendersWaitGroup)

    receiverWaitGroup.Add(1)
    go receiver(senderChennels, printerChannels, statusChannels, &receiverWaitGroup)

    printersWaitGroup.Add(2)
    go printer(printerChannels[0], &printersWaitGroup)
    go printer(printerChannels[1], &printersWaitGroup)
    
    receiverWaitGroup.Wait()
    close(statusChannels[0])
    close(statusChannels[1])
    fmt.Println("Receiver finished \n")

    sendersWaitGroup.Wait()
    fmt.Println("Senders finished \n")
    close(senderChennels[0])
    close(senderChennels[1])
    close(printerChannels[0])
    close(printerChannels[1])
    
    printersWaitGroup.Wait()
    fmt.Println("Printers finished \n")

    fmt.Println("end")
}

func sender(start int, senderChennel chan int, statusChannel chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    var done bool
    for !done {
        select {
            case status:= <-statusChannel:
                if status == EndProcess {
                    done = true
                    break
                }
                senderChennel<- start
                start++
            default: 
        }
    }
}

func receiver(senderChennels [sendersCount]chan int, printerChannels [printersCount]chan int, statusChannel [sendersCount]chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    statusChannel[0]<- Continue
    statusChannel[1]<- Continue
    var numbersPassedCount int
    var done bool
    for !done {
        select {
            case number, ok:= <-senderChennels[0]:
                if !ok {
                    senderChennels[0] = nil
                    break
                }
                numbersPassedCount++
                if numbersPassedCount > 20 {
                    statusChannel[0]<- EndProcess
                    senderChennels[0] = nil
                    break
                } else {
                    statusChannel[0]<- Continue
                }

                if number % 2 == 0 {
                    printerChannels[0]<- number
                } else {
                    printerChannels[1]<- number
                }
            case number, ok:= <-senderChennels[1]:
                if !ok {
                    senderChennels[1] = nil
                    break
                }

                numbersPassedCount++
                if numbersPassedCount > 20 {
                    statusChannel[1]<- EndProcess
                    senderChennels[1] = nil
                    done = true
                    break
                } else {
                    statusChannel[1]<- Continue
                }

                if number % 2 == 0 {
                    printerChannels[0]<- number
                } else {
                    printerChannels[1]<- number
                }
            default:
                if senderChennels[0] == nil && senderChennels[1] == nil {
                    done = true
                    break
                }
        }
    }
}

func printer(printerChannel chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    var numbers [maxNumbers]int
    var count int

    var done bool

    for !done {
        select {
            case number, ok:= <-printerChannel:
                if !ok {
                    printerChannel = nil
                    break
                }
                numbers[count] = number
                count++
            default: 
                if printerChannel == nil {
                    done = true
                    printNumbers(numbers, count)
                    break
                }
        }
    }
}

func printNumbers (numbers [maxNumbers]int, count int) {
    for i := 0; i < count; i++ {
        fmt.Println(strconv.Itoa(numbers[i]) + "\n")
    }
}
