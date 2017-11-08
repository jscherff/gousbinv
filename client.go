// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	`bytes`
	`encoding/json`
	`fmt`
	`io/ioutil`
	`net/http`
	`github.com/jscherff/cmdb/ci/peripheral/usb`
)

type httpStatus int

// Accepted returns true for a successful http response status.
func (this httpStatus) Accepted() (ok bool) {

	ok = true

	switch int(this) {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusAccepted:
	case http.StatusNoContent:
	case http.StatusNotModified:
	default:
		ok = false
	}

	return ok
}

// String implements the Stringer interface for httpStatus.
func (this httpStatus) String() (s string) {

	switch this {

	case http.StatusOK:
		s = `request processed, no errors`
	case http.StatusCreated:
		s = `request processed, object created`
	case http.StatusAccepted:
		s = `request processed, data accepted`
	case http.StatusNoContent:
		s = `request processed, no action taken`
	case http.StatusNotModified:
		s = `request processed, no changes found`
	case http.StatusBadRequest:
		s = `unsupported or malformed request`
	case http.StatusNotAcceptable:
		s = `insufficient or incorrect data`
	case http.StatusUnprocessableEntity:
		s = `unable to decode request`
	case http.StatusFailedDependency:
		s = `unsatisfied prerequisite`
	case http.StatusInternalServerError:
		s = `unable to process request`
	case http.StatusNotFound:
		s = `api endpoint not found`
	default:
		s = this.StatusText()
	}

	return s
}

// StatusText returns the HTTP status text associated with the status code.
func (this httpStatus) StatusText() (s string) {
	return http.StatusText(int(this))
}

// auth authenticates with the server using basic authentication and, if
// successful, obtains JWT for API authentication in a cookie.
func auth() error {

	url := fmt.Sprintf(`%s/%s`,
		conf.API.Server,
		conf.API.Endpoints[`cmdb_auth`],
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(conf.API.Auth.Username, conf.API.Auth.Password)

	if _, hs, err := httpRequest(req); err != nil {
		return err
	} else if !hs.Accepted() {
		return fmt.Errorf(`authentication failure - %s`, hs)
	} else {
		sl.Printf(`authentication success - %s`, hs)
	}

	return nil
}

// newSn obtains a serial number from the cmdbd server.
func newSn(dev usb.Serializer) (string, error) {

	var s string

	url := fmt.Sprintf(`%s/%s/%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoints[`usb_ci_newsn`],
		dev.Host(), dev.VID(), dev.PID(),
	)

	if j, err := dev.JSON(); err != nil {
		return ``, err
	} else if j, hs, err := httpPost(url, j); err != nil {
		return ``, err
	} else if !hs.Accepted() {
		return ``, fmt.Errorf(`serial number not generated - %s`, hs)
	} else if err = json.Unmarshal(j, &s); err != nil {
		return ``, err
	} else {
		sl.Printf(`serial number %q generated - %s`, s, hs)
		return s, nil
	}
}

// checkin checks a device in with the cmdbd server.
func checkin(dev usb.Reporter) (error) {

	url := fmt.Sprintf(`%s/%s/%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoints[`usb_ci_checkin`],
		dev.Host(), dev.VID(), dev.PID(),
	)

	if j, err := dev.JSON(); err != nil {
		return err
	} else if _, hs, err := httpPost(url, j); err != nil {
		return err
	} else if !hs.Accepted() {
		return fmt.Errorf(`checkin not accepted - %s`, hs)
	} else {
		sl.Printf(`checkin accepted - %s`, hs)
		return nil
	}
}

// checkout obtains the JSON representation of a serialized device object
// from the server using the unique key combination VID+PID+SN.
func checkout(dev usb.Auditer) ([]byte, error) {

	if dev.SN() == `` {
		sl.Printf(`device %s-%s skipping fetch, no SN`, dev.VID(), dev.PID())
		return nil, nil
	}

	url := fmt.Sprintf(`%s/%s/%s/%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoints[`usb_ci_checkout`],
		dev.Host(), dev.VID(), dev.PID(), dev.SN(),
	)

	if j, hs, err := httpGet(url); err != nil {
		return nil, err
	} else if !hs.Accepted() {
		return nil, fmt.Errorf(`device not retreived - %s`, hs)
	} else {
		sl.Printf(`device retrieved - %s`, hs)
		return j, nil
	}
}

// sendAudit submits changes from audit to the server in JSON format.
func sendAudit(dev usb.Auditer) (error) {

	url := fmt.Sprintf(`%s/%s/%s/%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoints[`usb_ci_audit`],
		dev.Host(), dev.VID(), dev.PID(), dev.SN(),
	)

	if j, err := json.Marshal(dev.GetChanges()); err != nil {
		return err
	} else if _, hs, err := httpPost(url, j); err != nil {
		return err
	} else if !hs.Accepted() {
		return fmt.Errorf(`audit not accepted - %s`, hs)
	} else {
		sl.Printf(`audit accepted - %s`, hs)
		return nil
	}
}

// vendor retrieves the vendor name given the vid.
func vendor(dev usb.Updater) (s string, err error) {

	url := fmt.Sprintf(`%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoints[`usb_meta_vendor`],
		dev.VID(),
	)

	if j, hs, err := httpGet(url); err != nil {
		return ``, err
	} else if !hs.Accepted() {
		return ``, fmt.Errorf(`lookup failed - %s`, hs)
	} else if err = json.Unmarshal(j, &s); err != nil {
		return ``, err
	} else {
		sl.Printf(`lookup succeeded - %s`, hs)
		return s, nil
	}
}

// product retrieves the product name given the vid and pid.
func product(dev usb.Updater) (s string, err error) {

	url := fmt.Sprintf(`%s/%s/%s/%s`,
		conf.API.Server,
		conf.API.Endpoints[`usb_meta_product`],
		dev.VID(), dev.PID(),
	)

	if j, hs, err := httpGet(url); err != nil {
		return ``, err
	} else if !hs.Accepted() {
		return ``, fmt.Errorf(`lookup failed - %s`, hs)
	} else if err = json.Unmarshal(j, &s); err != nil {
		return ``, err
	} else {
		sl.Printf(`lookup succeeded - %s`, hs)
		return s, nil
	}
}

// httpPost sends http POST requests to cmdbd server endpoints for other functions.
func httpPost(url string, j []byte ) (b []byte, hs httpStatus, err error) {

	if req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(j)); err != nil {
		return b, hs, err
	} else {
		req.Header.Add(`Content-Type`, `application/json; charset=UTF8`)
		return httpRequest(req)
	}
}

// httpGet sends http GET requests to cmdbd server endpoints for other functions.
func httpGet(url string) (b []byte, hs httpStatus, err error) {

	if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
		return b, hs, err
	} else {
		return httpRequest(req)
	}
}

// httpRequest sends http requests to cmdbd server endpoints for other functions.
func httpRequest(req *http.Request) (b []byte, hs httpStatus, err error) {

	req.Header.Add(`Accept`, `application/json; charset=UTF8`)
	req.Header.Add(`X-Custom-Header`, `cmdbc`)

	sl.Printf(`API call %s %s`, req.Method, req.URL)

	resp, err := httpClient.Do(req)

	if err != nil {
		return b, hs, err
	}

	defer resp.Body.Close()

	hs = httpStatus(resp.StatusCode)
	b, err = ioutil.ReadAll(resp.Body)

	return b, hs, err
}
