package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"fmt"
	"os"
)

type Reservation struct {
	HotelId      string `bson:"hotelId"`
	CustomerName string `bson:"customerName"`
	InDate       string `bson:"inDate"`
	OutDate      string `bson:"outDate"`
	Number       int    `bson:"number"`
}

type Number struct {
	HotelId      string `bson:"hotelId"`
	Number       int    `bson:"numberOfRoom"`
}

func initializeDatabase(url string) *mgo.Session {
	fmt.Printf("reservation db ip addr = %s\n", url)
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	swarmTaskSlot := os.Getenv("SWARM_TASK_SLOT")
	if swarmTaskSlot != "" && swarmTaskSlot != "1" {
		return session
	}

	c := session.DB("reservation-db").C("reservation")
	count, err := c.Find(&bson.M{"hotelId": "4"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0{
		err = c.Insert(&Reservation{"4", "Alice", "2015-04-09", "2015-04-10", 1})
		if err != nil {
			log.Fatal(err)
		}
	}

	c = session.DB("reservation-db").C("number")
	count, err = c.Find(&bson.M{"hotelId": "1"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0{
		err = c.Insert(&Number{"1", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelId": "2"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0{
		err = c.Insert(&Number{"2", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelId": "3"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0{
		err = c.Insert(&Number{"3", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelId": "4"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0{
		err = c.Insert(&Number{"4", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelId": "5"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0{
		err = c.Insert(&Number{"5", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	count, err = c.Find(&bson.M{"hotelId": "6"}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0{
		err = c.Insert(&Number{"6", 200})
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 7; i <= 1000; i++ {
		hotel_id := strconv.Itoa(i)
		count, err = c.Find(&bson.M{"hotelId": hotel_id}).Count()
		if err != nil {
			log.Fatal(err)
		}
		room_num := 200
		if i % 3 == 1 {
			room_num = 300
		} else if i % 3 == 2 {
			room_num = 250
		}
		if count == 0{
			err = c.Insert(&Number{hotel_id, room_num})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = c.EnsureIndexKey("hotelId")
	if err != nil {
		log.Fatal(err)
	}


	return session
}