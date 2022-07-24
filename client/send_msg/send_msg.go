package send_msg

import (
	"fmt"
	"messanger/internal"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

func SendMessage(ch chan internal.Message, serverAddr string) {
	var client fasthttp.Client

	for {
		msg := <-ch
		data, err := jsoniter.Marshal(&msg)
		if err != nil {
			fmt.Errorf("can't marshal msg, error %v", err)
		}

		req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.Header.SetMethod(fasthttp.MethodPost)
		req.SetRequestURI("http://" + serverAddr + "/send")
		req.SetBody(data)

		if err := client.Do(req, resp); err != nil {
			fmt.Errorf("message is't delivered, error %v", err)
			return
		}
		if resp.StatusCode() == fasthttp.StatusOK {
			fmt.Println("message successfully delivered")
		} else {
			fmt.Printf("message is't delivered, response code %v", resp.StatusCode())
		}
	}
}

func Authorize(message *internal.AuthorizeMessage, serverAddr string) (int, error) {
	var client fasthttp.Client

	for {
		data, err := jsoniter.Marshal(&message)
		if err != nil {
			fmt.Errorf("can't marshal msg, error %v", err)
		}

		req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.Header.SetMethod(fasthttp.MethodPost)
		req.SetRequestURI("http://" + serverAddr + "/authorize")
		req.SetBody(data)

		if err := client.Do(req, resp); err != nil {
			fmt.Errorf("authorization failed, error %v", err)
			return resp.StatusCode(), err
		}
	}
}
