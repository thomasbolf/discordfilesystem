package main

import (
	"fmt"
	"math"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func upload(filepath string) {
	//load environment variables
	godotenv.Load()
	channel := os.Getenv("CHANNEL")
	//create a session to upload the file
	session, _ := discordgo.New("Bot " + os.Getenv("TOKEN"))
	//convert file to a byte stream
	//bot says hi
	session.ChannelMessageSend(channel, "Hi")
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()
	//we will make a buffer that will take in 2mb at a time until the file is uploaded
	//two megabytes at a time: we iterate through the file, load data into a buffer, and send bytes to discord as a raw message of bytes
	buffer_start := 0
	buffer_size := 1024 * 2
	//get size of file
	file_info, _ := file.Stat()
	size_remaining := file_info.Size()
	fmt.Println("File size: ", size_remaining)

	for {
		buffer_size = int(math.Min(float64(buffer_size), float64(size_remaining)))
		fmt.Println("Buffer start: ", buffer_start)
		fmt.Println("Size remaining: ", size_remaining)
		//seek to the buffer_start
		file.Seek(int64(buffer_start), 0)
		//create a buffer
		buffer := make([]byte, buffer_size)
		//read the file into the buffer
		//read the file into the buffer, if the buffer_size is greater than the file size, we will read the remaining bytes

		n, err := file.Read(buffer)
		fmt.Println("Bytes read: ", n)
		//if we reach the end of the file, break

		if size_remaining <= 0 {
			break
		}

		if err != nil {
			fmt.Println("Error reading file: ", err)
			return
		}
		//send the bytes as a message containing only text (the bytes)
		fmt.Println(string(buffer[:n]))
		session.ChannelMessageSend(channel, string(buffer[:n]))

		//increment the buffer_start
		buffer_start += buffer_size
		size_remaining -= int64(n)
	}
	//close the file
	file.Close()
	session.Close()

}
