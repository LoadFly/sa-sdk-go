/*
 * Created by dengshiwei on 2020/01/06.
 * Copyright 2015－2020 Sensors Data Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func DoRequest(url, args string, to time.Duration) error {
	var resp *http.Response

	data := bytes.NewBufferString(args)

	req, _ := http.NewRequest("POST", url, data)

	client := &http.Client{Timeout: to}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		data, _ := json.Marshal(struct {
			StatusCode int
			Body       string
		}{resp.StatusCode, string(body)})
		return errors.New(string(data))
	}
	return nil
}

var (
	spWriter sync.Pool
	spBuffer sync.Pool
)

func init() {
	// 公共对象池,更极致的优化可以建多个池
	spWriter = sync.Pool{New: func() interface{} {
		buf := new(bytes.Buffer)
		return gzip.NewWriter(buf)
	}}
	spBuffer = sync.Pool{New: func() interface{} {
		return new(bytes.Buffer)
	}}
}

func gzipData(data string) ([]byte, error) {
	buf := spBuffer.Get().(*bytes.Buffer)
	w := spWriter.Get().(*gzip.Writer)
	w.Reset(buf)
	defer func() {
		// 归还buff
		buf.Reset()
		spBuffer.Put(buf)
		// 归还Writer
		spWriter.Put(w)
	}()
	_, err := w.Write([]byte(data))
	if err != nil {
		w.Close()
		return nil, err
	}
	w.Close()

	return buf.Bytes(), nil
}

func encodeData(data string) (string, error) {
	gdata, err := gzipData(data)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(gdata)
	return encoded, nil
}

func GeneratePostDataList(data string) (string, error) {
	edata, err := encodeData(data)
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Add("data_list", edata)
	v.Add("gzip", "1")

	uedata := v.Encode()

	return uedata, nil
}

func GeneratePostData(data string) (string, error) {
	edata, err := encodeData(data)
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Add("data", edata)
	v.Add("gzip", "1")

	uedata := v.Encode()

	return uedata, nil
}

func NowMs() int64 {
	return time.Now().UnixNano() / 1000000
}

// 合并公共属性
func MergeSuperProperty(superProperty map[string]interface{}, properties map[string]interface{}) map[string]interface{} {
	if superProperty == nil {
		return properties
	}

	for key, value := range superProperty {
		_, ok := properties[key]
		if !ok {
			properties[key] = value
		}
	}
	return properties
}

func DeepCopy(value map[string]interface{}) map[string]interface{} {
	ncopy := deepCopy(value)
	if nmap, ok := ncopy.(map[string]interface{}); ok {
		return nmap
	}

	return nil
}

func deepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = deepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = deepCopy(v)
		}

		return newSlice
	}

	return value
}
