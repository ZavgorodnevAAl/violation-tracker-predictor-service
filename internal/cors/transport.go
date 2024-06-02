package cors

import "net/http"

type transport struct {
	ModifyResponseFunc func(*http.Response) error
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method == "OPTIONS" {
		resp := &http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       http.NoBody,
		}

		resp.Header.Set("Access-Control-Allow-Origin", "*")
		resp.Header.Set("Access-Control-Allow-Headers", "*")
		resp.Header.Set("Access-Control-Allow-Methods", "*")

		return resp, nil
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if t.ModifyResponseFunc != nil {
		err = t.ModifyResponseFunc(resp)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
