package main

import (
	"fmt"
	"strconv"
	"testing"
)

func main() {
	test := map[string]map[string]string{
		"test1": {
			"Event":             "RTCPReceived",
			"Privilege":         "reporting,all",
			"Channel":           "PJSIP/test-voip-0012833d",
			"ChannelState":      "0",
			"ChannelStateDesc":  "Down",
			"CallerIDNum":       "testttt",
			"CallerIDName":      "<unknown>",
			"ConnectedLineNum":  "testttt",
			"ConnectedLineName": "<unknown>",
			"Language":          "en",
			"AccountCode":       "",
			"Context":           "test-voip",
			"Exten":             "s",
			"Priority":          "1",
			"Uniqueid":          "49b66532-0070-4194-943f-0f772fa0e93e",
			"Linkedid":          "49b66532-0070-4194-943f-0f772fa0e93e",
			"To":                "127.0.0.1:12427",
			"From":              "127.0.0.1:51071",
			"RTT":               "0.0000",
			"SSRC":              "0x00000008",
			"PT":                "200(SR)",
			"ReportCount":       "0",
			"SentNTP":           "2085978504.009999",
			"SentRTP":           "63940",
			"SentPackets":       "120",
			"SentOctets":        "12084",
		},
		"test2": {
			"Event": "RTCPReceived",

			"Privilege":                   "reporting,all",
			"Channel":                     "PJSIP/test-voip-001282e4",
			"ChannelState":                "6",
			"ChannelStateDesc":            "Up",
			"CallerIDNum":                 "testttt",
			"CallerIDName":                "<unknown>",
			"ConnectedLineNum":            "testttt",
			"ConnectedLineName":           "<unknown>",
			"Language":                    "en",
			"AccountCode":                 "",
			"Context":                     "test-voip",
			"Exten":                       "s",
			"Priority":                    "1",
			"Uniqueid":                    "0b7fe87e-e77b-4f67-b936-eabb3c42b2c9",
			"Linkedid":                    "0b7fe87e-e77b-4f67-b936-eabb3c42b2c9",
			"To":                          "127.0.0.1:10401",
			"From":                        "127.0.0.1:21023",
			"RTT":                         "0.3235",
			"SSRC":                        "0x00000011",
			"PT":                          "200(SR)",
			"ReportCount":                 "1",
			"SentNTP":                     "2085978513.769989",
			"SentRTP":                     "1902891976",
			"SentPackets":                 "889",
			"SentOctets":                  "142240",
			"Report0SourceSSRC":           "0x0180c008",
			"Report0FractionLost":         "0",
			"Report0CumulativeLost":       "0",
			"Report0HighestSequence":      "2421",
			"Report0SequenceNumberCycles": "0",
			"Report0IAJitter":             "0",
			"Report0LSR":                  "1061144719",
			"Report0DLSR":                 "2.8500",
		},
	}

	validateMessage(test["test1"])

	validateMessage(test["test2"])
}

func TestValidate(t *testing.T) {

}

func validateMessage(msg map[string]string) bool {
	if msg["PT"] != "200(SR)" {
		return false
	}

	if v, err := strconv.ParseFloat(msg["RTT"], 64); err == nil {
		fmt.Println(v)
	}

	if v, err := strconv.ParseFloat(msg["Report0IAJitter"], 64); err == nil {
		fmt.Println(v)
	}

	// get report count
	count, err := strconv.ParseInt(msg["ReportCount"], 10, 8)
	if err != nil {
		fmt.Println("Could not get report count.")
		return true
	}

	for i := 0; i < int(count); i++ {
		var key string

		key = fmt.Sprintf("Report%dIAJitter", i)
		if v, err := strconv.ParseFloat(msg[key], 64); err == nil {
			fmt.Println(v)
		}

		key = fmt.Sprintf("Report%dFractionLost", i)
		if v, err := strconv.ParseFloat(msg[key], 64); err == nil {
			fmt.Println(v)
		}

		key = fmt.Sprintf("Report%dDLSR", i)
		if v, err := strconv.ParseFloat(msg[key], 64); err == nil {
			fmt.Println(v)
		}
	}

	return true
}
