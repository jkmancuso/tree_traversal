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

func NewService(downstream *svc) svc {
	return svc{
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
	var svcSlice []svc
	var id int

	for i := 0; i < num; i++ {

		if i == 0 {
			downstream = &svc{id: 0}
		} else {
			downstream = &svcSlice[rand.Intn(len(svcSlice))]
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

		service.id = id
		fmt.Println("Running: ", query, id)

		svcSlice = append(svcSlice, service)

	}

	traverse([]int{}, svcSlice, svcSlice[len(svcSlice)-1])

}

func traverse(visited []int, services []svc, currentSvc svc) {

	if len(services) == 0 {
		fmt.Println("Empty, returning")
		return
	}

	if slices.Contains(visited, currentSvc.id) {
		fmt.Printf("Found %d in %v\n", currentSvc.id, visited)
		return
	}

	fmt.Println("At svc node: ", currentSvc.id)

	visited = append(visited, currentSvc.id)

	if currentSvc.downstream != nil {
		fmt.Println("Next visit is downstream id: ", currentSvc.downstream.id)
		traverse(visited, services[0:len(services)-1], *currentSvc.downstream)
	} else {
		traverse(visited, services[0:len(services)-1], services[len(services)-1])
	}

}
