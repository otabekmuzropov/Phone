package main

import (
	pb "bitbucket.org/alien_soft/phone/genproto/phone"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
)

type server struct {
	db *sqlx.DB
}

func (s *server) Create(ctx context.Context, req *pb.Phone) (*pb.Phone, error)  {
	var phone pb.Phone
	log.Println("request server", req.Name)

	id := rand.Uint32()

	create := `insert into phone(id, name, ram, screen_diagnol, battery, memory, status) values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.db.Exec(create, id, req.Name, req.Ram,req.ScreenDiagnol, req.Battery, req.Memory, req.Status)

	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow("select id, name, ram, screen_diagnol, battery, memory, status from phone where id = $1", id)

	err = row.Scan(&phone.Id, &phone.Name, &phone.Ram, &phone.ScreenDiagnol, &phone.Battery, &phone.Memory, &phone.Status)

	if err != nil {
		return nil, err
	}

	return &phone, nil
}

func (s *server) Update(ctx context.Context, req *pb.Phone) (*pb.Phone, error) {
	var p pb.Phone
	log.Println("Received request")

	update := `update phone set name = $2, ram = $3, screen_diagnol = $4, battery = $5, memory = $6, status = $7 where id = $1`

	_, err := s.db.Exec(update, req.Id, req.Name, req.Ram,req.ScreenDiagnol, req.Battery, req.Memory, req.Status)

	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow("select id, name, ram, screen_diagnol, battery, memory, status from phone where id = $1", req.Id)

	err = row.Scan(&p.Id, &p.Name, &p.Ram, &p.ScreenDiagnol, &p.Battery, &p.Memory, &p.Status)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*empty.Empty, error) {
	log.Println("Received request")

	delete := `delete from phone where id = $1`

	_, err := s.db.Exec(delete, req.Id)

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) List(ctx context.Context, req *empty.Empty) (*pb.ListResponse, error) {
	var phones []*pb.Phone

	log.Println("Received request")

	rows, err := s.db.Queryx("select id, name, ram, screen_diagnol, battery, memory, status from phone")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var phone pb.Phone

		err = rows.Scan(&phone.Id, &phone.Name, &phone.Ram, &phone.ScreenDiagnol, &phone.Battery, &phone.Memory, &phone.Status)

		if err != nil {
			return nil, err
		}

		phones = append(phones, &phone)
	}

	return &pb.ListResponse{Phones:phones}, nil
}

func (s *server) GetOne(ctx context.Context, req *pb.GetOneRequest) (*pb.Phone, error)  {
	var phone pb.Phone
	log.Println("Received request")

	row := s.db.QueryRow("select id, name, ram, screen_diagnol, battery, memory, status from phone where id = $1", req.Id)

	err := row.Scan(&phone.Id, &phone.Name, &phone.Ram, &phone.ScreenDiagnol, &phone.Battery, &phone.Memory, &phone.Status)

	if err != nil {
		return nil, err
	}

	return &phone, nil
}

func (s *server) List2(ctx context.Context, req *pb.List2Request) (*pb.List2Response, error)  {
	var phones []*pb.Phone
	log.Println("Received request")

	offset := (req.Offset - 1) * 5

	rows, err := s.db.Queryx("select id, name, ram, screen_diagnol, battery, memory, status from phone limit 5 offset $1", offset)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var phone pb.Phone

		err = rows.Scan(&phone.Id, &phone.Name, &phone.Ram, &phone.ScreenDiagnol, &phone.Battery, &phone.Memory, &phone.Status)

		if err != nil {
			return nil, err
		}

		phones = append(phones, &phone)
	}

	return &pb.List2Response{Phones:phones}, nil
}

func (s *server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	var phones []*pb.Phone

	log.Println("Received request")

	letter := req.Letter

	rows, err := s.db.Queryx("select id, name, ram, screen_diagnol, battery, memory, status from phone where name ilike $1", "%" + letter + "%")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var phone pb.Phone

		err = rows.Scan(&phone.Id, &phone.Name, &phone.Ram, &phone.ScreenDiagnol, &phone.Battery, &phone.Memory, &phone.Status)

		if err != nil {
			return nil, err
		}

		phones = append(phones, &phone)
	}

	return &pb.SearchResponse{Phones:phones}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":50053")

	if err != nil {
		log.Println("error while listening port", err)
		return
	}

	db, err := sqlx.Connect("postgres", "user=postgres dbname=test password=123 sslmode=disable")

	if err != nil {
		log.Println("error while connecting db", err)
		return
	}

	s := grpc.NewServer()

	pb.RegisterPhoneServiceServer(s, &server{db:db})

	log.Println("listening %d port", 50053)

	if err := s.Serve(listen); err != nil {
		log.Println("failed serve", err)
		return
	}
}
