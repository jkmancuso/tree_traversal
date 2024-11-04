package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"slices"

	_ "github.com/lib/pq"
)

const (
	num     = 6
	connStr = "user=postgres dbname=test sslmode=disable"
	query   = "INSERT INTO services(latency, cpu, err_rate, downstream)"
)

func NewService(downstream *svc) *svc {
	return &svc{
		latency:    rand.Intn(100),
		cpu:        rand.Intn(100),
		errRate:    rand.Intn(100),
		downstream: downstream,
	}
}

type svc struct {
	id         int
	latency    int
	cpu        int
	errRate    int
	downstream *svc
}

func main() {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	var downstream *svc
	var id int
	var idSlice []int

	svcMap := make(map[int]*svc)

	for i := 0; i < num; i++ {

		if i == 0 {
			downstream = &svc{id: 0}
		} else {
			randSvcId := idSlice[rand.Intn(len(idSlice))]
			fmt.Println(randSvcId)
			downstream = svcMap[randSvcId]
		}

		service := NewService(downstream)

		query := fmt.Sprintf("%v VALUES(%d, %d, %d, %d) RETURNING id",
			query,
			service.latency,
			service.cpu,
			service.errRate,
			downstream.id)

		fmt.Println("Running query")

		err = db.QueryRow(query).Scan(&id)

		if err != nil {
			log.Fatal(err)
		}

		idSlice = append(idSlice, id)

		service.id = id
		fmt.Println("Running: ", query, id)

		svcMap[id] = service

	}

	traverse([]int{}, svcMap, id)

}

func traverse(visited []int, svcMap map[int]*svc, currentId int) {

	if !slices.Contains(visited, currentId) {
		fmt.Printf("ID %d-> downstream %d \n", currentId, svcMap[currentId].downstream.id)

		if svcMap[currentId].downstream.id == 0 {
			fmt.Println("No more downstream found, exiting")
			return
		}

		visited = append(visited, currentId)

		traverse(visited, svcMap, svcMap[currentId].downstream.id)
	}

}
