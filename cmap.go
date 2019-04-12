package main

type ChannelMap struct {

	// Channels for Asking, Adding, Reducing, Kill
	AskSendChan chan string
	AskReceiveChan chan int
	AddChan chan string
	ReduceChan chan ReduceFunc
	ReduceRetStrChan chan string
	ReduceRetIntChan chan int
	KillChan chan int

	// Map
	Words map[string]int

}

// Decides what action to take based on what channel has data in it. 
// Each case performs one task using the into it read from the channel and possible provides a response.
// Loops indefinitely until Stop() is called.
func (c *ChannelMap) Listen() {
	for {
   		select {
   			// Do the reduce function
   			case functor := <-c.ReduceChan:
   				accum_str := "INVALID"
   				accum_int := 0
   				for k, v := range c.Words {
   					accum_str, accum_int = functor(accum_str, accum_int, k, v)
   				}
   				go func(){ c.ReduceRetIntChan <- accum_int }()
   				go func(){ c.ReduceRetStrChan <- accum_str }()
   			// Lookup the word count
   			case wordToSearch := <-c.AskSendChan:
   				go func(){c.AskReceiveChan <- c.Words[wordToSearch]}()
   			// Add word to map
   			case wordToAdd := <-c.AddChan:
   				c.Words[wordToAdd]++
	  		// Stop the loop if there's something in the kill channel
			case <-c.KillChan:
				return
		}
	}
}

// Send 1 to KillChan to stop the Listen() loop
func (c *ChannelMap) Stop() {
	c.KillChan <- 1
}

// Send the functor to the channel to be processed in Listen()
func (c *ChannelMap) Reduce(functor ReduceFunc, accum_str string, accum_int int) (string, int) {
	
	go func(){ c.ReduceChan <- functor }()

	accum_int = <-c.ReduceRetIntChan
	accum_str = <-c.ReduceRetStrChan

	return accum_str, accum_int
}

// Add word to channel
func (c *ChannelMap) AddWord(word string) {
	c.AddChan <- word
}

// Send the word to ChannelMap AddSendChan, return the number in AskReceiveChan
func (c *ChannelMap) GetCount(word string) int {
	c.AskSendChan <- word
	return <-c.AskReceiveChan
}

// Create a new ChannelMap with map and buffered channels
func NewChannelMap() *ChannelMap {
	cm := new(ChannelMap)
	cm.AskSendChan = make(chan string, ASK_BUFFER_SIZE)
	cm.AskReceiveChan = make(chan int)
	cm.AddChan = make(chan string, ADD_BUFFER_SIZE)
	cm.ReduceChan = make(chan ReduceFunc, REDUCE_BUFFER_SIZE)
	cm.ReduceRetStrChan = make(chan string, REDUCE_BUFFER_SIZE)
	cm.ReduceRetIntChan = make(chan int, REDUCE_BUFFER_SIZE)
	cm.KillChan = make(chan int)
	cm.Words = make(map[string]int)
	return cm
}
