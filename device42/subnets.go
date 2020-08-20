package device42

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Subnet ...
type Subnet struct {
	Subnets []struct {
		Name           string      `json:"name"`
		Network        string      `json:"network"`
		MaskBits       int         `json:"mask_bits"`
		RangeBegin     string      `json:"range_begin"`
		RangeEnd       string      `json:"range_end"`
		VrfGroupName   string      `json:"vrf_group_name"`
		Assigned       string      `json:"assigned"`
		Allocated      string      `json:"allocated"`
		ParentSubnetID interface{} `json:"parent_subnet_id,omitempty"`
		SubnetID       int         `json:"subnet_id"`
		Gateway        interface{} `json:"gateway"`
	} `json:"subnets"`
}

//Child ...
type Child struct {
	SubnetID int    `json:"subnet_id"`
	MaskBits string `json:"mask_bits"`
	Network  string `json:"network"`
}

//
func (c *Client) GetSubnet(Name string) (*Subnet, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/1.0/subnets/?name=%s", c.HostURL, Name), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	subnet := Subnet{}
	err = json.Unmarshal(body, &subnet)
	if err != nil {
		return nil, err
	}

	return &subnet, nil
}

//
func (c *Client) GetSubnetID(ID string) (*Subnet, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/1.0/subnets/?subnet_id=%s", c.HostURL, ID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	subnet := Subnet{}
	err = json.Unmarshal(body, &subnet)
	if err != nil {
		return nil, err
	}

	return &subnet, nil
}

//
func (c *Client) SetSubnet(Network string, Mask int, Name string) error {
	sMask := fmt.Sprint(Mask)
	spayload := fmt.Sprintf("network=%s&mask_bits=%s&allocated=yes&name=%s", Network, sMask, Name)
	payload := strings.NewReader(spayload)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/1.0/subnets/", c.HostURL), payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

//
func (c *Client) SetChild(ParentID int, Mask int) (*Child, error) {
	sParentID := fmt.Sprint(ParentID)
	sMask := fmt.Sprint(Mask)
	spayload := fmt.Sprintf("mask_bits=%s&parent_subnet_id=%s", sMask, sParentID)

	payload := strings.NewReader(spayload)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/1.0/subnets/create_child/", c.HostURL), payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	child := Child{}
	err = json.Unmarshal(body, &child)
	if err != nil {
		return nil, err

	}

	return &child, nil
}

//
func (c *Client) SetGateway(Name string) error {
	subnet, err := c.GetSubnet(Name)
	if err != nil {
		return err
	}

	network := subnet.Subnets[0].Network
	maskbits := fmt.Sprint(subnet.Subnets[0].MaskBits)
	gateway := subnet.Subnets[0].RangeEnd

	rangeend := subnet.Subnets[0].RangeEnd
	split := strings.Split(rangeend, ".")
	x, _ := strconv.Atoi(split[3])
	x--
	rangeend = fmt.Sprintf("%s.%s.%s.%s", split[0], split[1], split[2], fmt.Sprint(x))

	spayload := fmt.Sprintf("network=%s&mask_bits=%s&gateway=%s&range_end=%s", network, maskbits, gateway, rangeend)
	payload := strings.NewReader(spayload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/1.0/subnets/", c.HostURL), payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

//
func (c *Client) DeleteSubnet(subnetID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/1.0/subnets/%s/", c.HostURL, subnetID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
