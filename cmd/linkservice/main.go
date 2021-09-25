package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/pavelzagorodnyuk/linkservice/internal/api"
	service "github.com/pavelzagorodnyuk/linkservice/internal/linkservice"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var (
	port = ":50051"

	DBConnParams = os.ExpandEnv("user=$POSTGRES_USER password=$POSTGRES_PASSWORD host=$DB_HOST port=$DB_PORT dbname=$POSTGRES_DB sslmode=disable")
)

func main() {
	// задаем начальное значение для генератора псевдослучайных чисел
	rand.Seed(time.Now().UnixNano())

	// устанавливаем подключение к базе данных
	log.Println("Connecting to database...")

	db, err := sql.Open("postgres", DBConnParams)
	if err != nil {
		log.Fatalf("failed to connect to database: %v\n", err)
	}

	for i := 5; i > 0 && db.Ping() != nil; i-- {
		if i > 1 {
			log.Println("failed to connect to database. The next attempt is in 5 seconds...")
			time.Sleep(5 * time.Second)
		} else {
			log.Fatalln("failed to connect to database. Exit...")
		}
	}

	defer db.Close()

	// запускаем gRPC сервер
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer l.Close()

	srv := grpc.NewServer()

	api.RegisterLinkServiceServer(srv, &service.GRPCServer{Database: db})

	log.Println("Starting gRPC server...")

	if err := srv.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
