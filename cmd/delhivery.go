package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// genrated using https://mholt.github.io/json-to-go/
type Delhivery struct {
	Meta struct {
		RequestID string `json:"requestId"`
	} `json:"meta"`
	Data []struct {
		Status struct {
			Status         string `json:"status"`
			StatusDateTime string `json:"statusDateTime"`
			StatusType     string `json:"statusType"`
			Instructions   string `json:"instructions"`
		} `json:"status"`
		Slot struct {
			StUtc  int    `json:"stUtc"`
			Src    string `json:"src"`
			From   string `json:"from"`
			To     string `json:"to"`
			Dsrc   string `json:"dsrc"`
			Date   string `json:"date"`
			SrcApp string `json:"srcApp"`
		} `json:"slot"`
		EstimatedDate string `json:"estimatedDate"`
		Flow          string `json:"flow"`
		FlowDirection string `json:"flowDirection"`
		ReferenceNo   string `json:"referenceNo"`
		PackageType   string `json:"packageType"`
		Awb           string `json:"awb"`
		Destination   string `json:"destination"`
		Scans         []struct {
			ScanDateTime         string `json:"scanDateTime"`
			ScanNslCode          string `json:"scanNslCode"`
			CityLocation         string `json:"cityLocation"`
			ScanType             string `json:"scanType"`
			Scan                 string `json:"scan"`
			ScannedLocation      string `json:"scannedLocation"`
			Instructions         string `json:"instructions"`
			Status               string `json:"status"`
			AdditionalScanRemark string `json:"additionalScanRemark,omitempty"`
			Destination          string `json:"destination,omitempty"`
		} `json:"scans"`
		AwbHash           string `json:"awbHash"`
		HqStatus          string `json:"hqStatus"`
		DispatchCenterID  string `json:"dispatchCenterId"`
		ProductType       string `json:"productType"`
		CovidZone         string `json:"covidZone"`
		Essential         bool   `json:"essential"`
		ContainmentArea   bool   `json:"containmentArea"`
		ExpectedDate      string `json:"expectedDate"`
		ProductName       string `json:"productName"`
		ClientName        string `json:"clientName"`
		ConsigneeAddress  string `json:"consigneeAddress"`
		NextTrialDate     string `json:"nextTrialDate"`
		IsAddressSpecific bool   `json:"isAddressSpecific"`
		IsPersonSpecific  bool   `json:"isPersonSpecific"`
		IsInternational   bool   `json:"isInternational"`
	} `json:"data"`
}

// delhiveryCmd represents the delhivery command
var delhiveryCmd = &cobra.Command{
	Use:   "delhivery",
	Short: "track delhivery status of your consignment",
	Long:  "track delhivery status of your consignment",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide the tracking number")
			return
		}
		detailed, _ := cmd.Flags().GetBool("detailed")
		fmt.Println("Tracking AWB:"+args[0], "on delhivery ")
		trackDelhivery(args, detailed)

	},
}

func init() {
	rootCmd.AddCommand(delhiveryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	delhiveryCmd.PersistentFlags().BoolP("detailed", "d", false, "detailed tracking information")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// delhiveryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func trackDelhivery(args []string, detailed bool) {
	var currentStatus Delhivery
	resp, err := http.Get("https://uxxbqylwa3.execute-api.ap-southeast-1.amazonaws.com/prod/track?waybillId=" + args[0])
	if err != nil {
		log.Fatal(err.Error())
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	json.Unmarshal(bodyBytes, &currentStatus)
	scans := currentStatus.Data[0].Scans
	layout := "2006-01-02T15:04:05"
	layout1 := "Mon, 01/02/06, 03:04PM"
	reportLines := 0
	if !detailed {
		reportLines = len(scans) - 4
	}
	for i := len(scans) - 1; i >= reportLines; i-- {
		time, err := time.Parse(layout, scans[i].ScanDateTime)
		if err != nil {
			log.Fatal(err.Error())
		}
		scanstatus := strings.ToLower(scans[i].Scan)
		var status string
		switch {
		case strings.Contains(scanstatus, "pending"):
			status = "⌛...Pending"
		case strings.Contains(scanstatus, "transit"):
			status = "🚚"
		case strings.Contains(scanstatus, "manifest"):
			status = "📒...Booked"
		}
		result := fmt.Sprint(time.Format(layout1), " @ ", scans[i].ScannedLocation, scans[i].Instructions, ".........................................")
		fmt.Println(result[:100], status)

	}

}
