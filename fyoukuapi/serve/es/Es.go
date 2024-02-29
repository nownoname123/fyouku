package els

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyoukuapi/config"
	"io/ioutil"

	"io"
	"log"

	"net/http"
)

var esUrl string

func init() {
	esUrl = "http://" + config.Url + ":9200/"

}

// ReqSearchData 解析获取到的值
type ReqSearchData struct {
	Hits HitsData `json:"hits"`
}
type HitsData struct {
	Total TotalData     `json:"total"`
	Hits  []HitsTwoData `json:"hits"`
}
type TotalData struct {
	Value    int
	Relation string
}
type HitsTwoData struct {
	Source json.RawMessage `json:"_source"`
}

// EsSearch 搜索功能
func EsSearch(indexName string, query map[string]interface{}, from int, size int,
	sort []map[string]string) HitsData {
	searchQuery := map[string]interface{}{
		"query": query,
		"from":  from,
		"size":  size,
		"sort":  sort,
	}

	// 将查询结构体转换为 JSON 格式的字节切片
	requestBody, err := json.Marshal(searchQuery)
	if err != nil {
		fmt.Println(err)
		return HitsData{} // 返回空的 HitsData 结构体
	}

	// 发送 HTTP 请求
	req, err := http.NewRequest("POST", esUrl+indexName+"/_search", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
		return HitsData{} // 返回空的 HitsData 结构体
	}
	req.Header.Set("Content-Type", "application/json")

	// 执行 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return HitsData{} // 返回空的 HitsData 结构体
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 读取响应体
	var stb ReqSearchData
	err = json.NewDecoder(resp.Body).Decode(&stb)
	if err != nil {
		fmt.Println(err)
		return HitsData{} // 返回空的 HitsData 结构体
	}

	return stb.Hits

}
func EsAdd(indexName string, id string, body map[string]interface{}) bool {
	requestBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", esUrl+indexName+"/_doc/"+id, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("error in es add:", err)
		return false
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error in sending request:", err)
		return false
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		fmt.Println("unexpected status code:", res.StatusCode)
		return false
	}

	// 读取响应体
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error in reading response body:", err)
		return false
	}

	fmt.Println("添加成功:", string(requestBody))
	fmt.Println("返回数据:", string(responseData))
	return true
}
func EsEdit(indexName string, id string, body map[string]interface{}) bool {
	bodyData := map[string]interface{}{
		"doc": body,
	}

	requestBody, _ := json.Marshal(bodyData)
	req, err := http.NewRequest("POST", esUrl+indexName+"/_doc/"+id+"/_update", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		fmt.Println(err)
		return false // 返回空的 HitsData 结构体
	}
	fmt.Println("修改成功")
	return true
}
func EsDelete(indexName string, id string) bool {

	req, err := http.NewRequest("POST", esUrl+indexName+"/_doc/"+id, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		fmt.Println(err)
		return false // 返回空的 HitsData 结构体
	}
	fmt.Println("删除成功")
	return true
}
