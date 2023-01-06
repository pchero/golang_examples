package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Printf("hello world\n")

	variables := map[string]string{
		"JITTERBUFFER(disabled)":                "",
		"PJSIP_HEADER(add,P-Asserted-Identity)": "<sip:Trunker_Test_Call_21018@peer.Test.io>",
		"PJSIP_HEADER(add,Privacy)":             "id",
		"PJSIP_HEADER(add,X-Test-Call-ID)":      "28517198-e4f2-4c26-be7f-92e3cc0c6589",
		"PJSIP_HEADER(add,X-Test-Call-ID-B)":    "c156d356-4834-11ec-8e12-42dec4b7a538",
		"PJSIP_HEADER(add,X-Test-Dialed)":       "+14999998889 - Trunker Test Call",
		"PJSIP_HEADER(add,X-Test-Forward-URI)":  "sip:+14999998889@trunker-test-call.mon.svc.cluster.local:5062;transport=udp",
		"PJSIP_HEADER(add,X-Porta-Customer-ID)": "21018",
		"_Test_CALLED_USER":                     "NO_Test_USER_ID",
		"_Test_LINKED_CALL":                     "28517198-e4f2-4c26-be7f-92e3cc0c6589",
		"_Test_ORGANIZATION":                    "d328b0ba-c084-11ea-bfdb-3623424f375e",
		"_Test_USER":                            "",
	}
	fmt.Printf("Original variables. variables:\n%v", variables)

	res := ""
	for k, v := range variables {
		fmt.Printf("k: %s, v: %s\n", k, v)

		tmp := strings.ToLower(k)
		if !strings.HasPrefix(tmp, "_Test_") {
			continue
		}

		tmp = strings.TrimPrefix(tmp, "_")
		fmt.Printf("Trimed key: %s\n", tmp)
		res = fmt.Sprintf("%s,%s=%s", res, tmp, v)
	}

	fmt.Printf("Res: %s", res)
}
