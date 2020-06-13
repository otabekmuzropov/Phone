package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	pb "bitbucket.org/alien_soft/phone/genproto/phone"
	//"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

type Phone struct {
	ID uint64 `json:"id"`
	Name string `json:"name"`
	Ram string `json:"ram"`
	ScreenDiagnol uint64 `json:"screen_diagnol"`
	Battery string `json:"battery"`
	Memory string `json:"memory"`
	Status string `json:"status"`
}

func main() {
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())

	if err != nil {
		log.Println("error while dialing", err)
		return
	}
	defer conn.Close()

	p := pb.NewPhoneServiceClient(conn)

	r := gin.Default()

	r.POST("/phone/", func(context *gin.Context) {
		var phone Phone

		err := context.ShouldBindJSON(&phone)

		if err != nil {
			log.Println("error while binding json", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		resp, err := p.Create(context, &pb.Phone{Name:phone.Name, Ram:phone.Ram, ScreenDiagnol:phone.ScreenDiagnol, Battery:phone.Battery, Memory:phone.Memory,Status:phone.Status})

		if err != nil {
			log.Println("error while creating phone", err)
			context.JSON(http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusCreated, resp)

	})

	r.PUT("/phone/:id", func(context *gin.Context) {
		var phone Phone

		err := context.ShouldBindJSON(&phone)

		if err != nil {
			log.Println("error while binding json", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		a := context.Param("id")
		id, err := strconv.Atoi(a)

		if err != nil {
			log.Println("param is not uint", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		resp, err := p.Update(context, &pb.Phone{Id:uint64(id), Name:phone.Name, Ram:phone.Ram, ScreenDiagnol:phone.ScreenDiagnol, Battery:phone.Battery, Memory:phone.Memory,Status:phone.Status})

		if err != nil {
			log.Println("error while updating phone", err)
			context.JSON(http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, resp)
	})


	/*r.GET("/", func(context *gin.Context) {
		a := context.Param("offset")
		offset, err := strconv.Atoi(a)

		if err != nil {
			log.Println("error give a offset", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		resp, err := p.List2(context, &pb.List2Request{Offset:uint64(offset)})

		if err != nil {
			log.Println("error getting phones", err)
			context.JSON(http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, resp)
	})*/

		r.GET("/", func(context *gin.Context) {
		letter := context.Query("letter")
		log.Println(letter)
		//letter := string(a)

		resp, err := p.Search(context, &pb.SearchRequest{Letter:letter})

		if err != nil {
			log.Println("error while searching....", err)
			context.JSON(http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, resp)

	})


	r.DELETE("/phone/:id", func(context *gin.Context) {
		var phone Phone

		err := context.ShouldBindJSON(&phone)

		if err != nil {
			log.Println("error while binding json", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		a := context.Param("id")
		id, err := strconv.Atoi(a)

		if err != nil {
			log.Println("param is not uint64", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		resp, err := p.Delete(context, &pb.DeleteRequest{Id:uint64(id)})

		if err != nil {
			log.Println("error while deleting phone", err)
			context.JSON(http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, resp)

	})

	r.GET("/phone/", func(context *gin.Context) {
		pageVal := context.DefaultQuery("page", "1")

		page, err := strconv.Atoi(pageVal)

		if err != nil {
			log.Println("error while getting page", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		resp, err := p.List2(context, &pb.List2Request{Offset:uint64(page)})

		if err != nil {
			log.Println("error while getting all phones", err)
			context.JSON(http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, resp)

	})

	r.GET("/phone/:id", func(context *gin.Context) {

		a := context.Param("id")
		id, err := strconv.Atoi(a)

		if err != nil {
			log.Println("param is not uint64", err)
			context.JSON(http.StatusBadRequest, err)
			return
		}

		resp, err := p.GetOne(context, &pb.GetOneRequest{Id:uint64(id)})

		if err != nil {
			log.Println("error while getting phone", err)
			context.JSON(http.StatusInternalServerError, err)
			return
		}

		context.JSON(http.StatusOK, resp)

	})

	r.Run(":5051")
}



