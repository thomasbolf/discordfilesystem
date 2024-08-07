package main

import (
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func upload(filepath string) {
	//load environment variables
	godotenv.Load()
	channel := os.Getenv("CHANNEL")

	//create the discord bot and open the file we want to download
	session, _ := discordgo.New("Bot " + os.Getenv("TOKEN"))
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()
	defer session.Close()

	buffer_start := 0
	//arbitrary, around 32KB
	buffer_size := 32000
	file_info, _ := file.Stat()
	size_remaining := file_info.Size()

	//for record keeping
	fmt.Println("File size: ", size_remaining)
	start := time.Now()
	timeSendingMessages := 0
	//we are iterating until we are out of bytes to read
	for {
		//if we are almost out of bytes, take the size remaining, else take the normal buffer size
		buffer_size = int(math.Min(float64(buffer_size), float64(size_remaining)))
		fmt.Println("Buffer start: ", buffer_start)
		fmt.Println("Size remaining: ", size_remaining)
		file.Seek(int64(buffer_start), 0)
		buffer := make([]byte, buffer_size)
		n, err := file.Read(buffer)
		fmt.Println("Bytes read: ", n)
		if size_remaining <= 0 {
			break
		}

		if err != nil {
			fmt.Println("Error reading file: ", err)
			return
		}
		beforeCall := time.Now()
		message := hex.EncodeToString(buffer)
		os.WriteFile("message.txt", []byte(message), 0644)
		readr, _ := os.Open("message.txt")
		session.ChannelFileSend(channel, "message.txt", readr)

		fmt.Println("Time taken to send message: ", time.Since(beforeCall))
		timeSendingMessages += int(time.Since(beforeCall).Milliseconds())
		buffer_start += buffer_size
		size_remaining -= int64(n)
	}
	fmt.Println("Time taken to send messages (ms): ", timeSendingMessages)
	fmt.Println("Time taken: ", time.Since(start))

}
