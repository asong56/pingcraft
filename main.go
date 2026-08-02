package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type speedtestResult struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Ping      struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
	} `json:"ping"`
	Download struct {
		Bandwidth int64 `json:"bandwidth"`
		Bytes     int64 `json:"bytes"`
	} `json:"download"`
	Upload struct {
		Bandwidth int64 `json:"bandwidth"`
		Bytes     int64 `json:"bytes"`
	} `json:"upload"`
	PacketLoss float64 `json:"packetLoss"`
	ISP        string  `json:"isp"`
	Server     struct {
		Name    string `json:"name"`
		Location string `json:"location"`
		Country string `json:"country"`
	} `json:"server"`
	Result struct {
		URL string `json:"url"`
	} `json:"result"`
}

func bandwidthToMbps(bandwidthBytesPerSec int64) float64 {
	return float64(bandwidthBytesPerSec) * 8 / 1_000_000
}

func findSpeedtestExe() (string, error) {
	exePath, err := os.Executable()
	if err == nil {
		localPath := filepath.Join(filepath.Dir(exePath), "speedtest.exe")
		if _, statErr := os.Stat(localPath); statErr == nil {
			return localPath, nil
		}
	}

	if pathExe, err := exec.LookPath("speedtest.exe"); err == nil {
		return pathExe, nil
	}
	if pathExe, err := exec.LookPath("speedtest"); err == nil {
		return pathExe, nil
	}

	return "", fmt.Errorf("speedtest.exe not found，download official CLI from https://www.speedtest.net/apps/cli")
}

func main() {
	fmt.Println("=== Simple Speedtest Tool ===")

	exePath, err := findSpeedtestExe()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	fmt.Println("(Usually need 10-30 sec)...")
	fmt.Println()

	cmd := exec.Command(exePath, "--accept-license", "--accept-gdpr", "--format=json")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed:", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Println("Details:", string(exitErr.Stderr))
		}
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	var result speedtestResult
	if err := json.Unmarshal(output, &result); err != nil {
		fmt.Println("Parsing failed:", err)
		fmt.Println("Raw output:", string(output))
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	printResult(result)

	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}

func printResult(r speedtestResult) {
	downloadMbps := bandwidthToMbps(r.Download.Bandwidth)
	uploadMbps := bandwidthToMbps(r.Upload.Bandwidth)

	packetLossStr := "Not available"
	if r.PacketLoss >= 0 {
		packetLossStr = fmt.Sprintf("%.1f%%", r.PacketLoss)
	}

	serverInfo := r.Server.Name
	if r.Server.Location != "" {
		serverInfo = fmt.Sprintf("%s (%s", r.Server.Name, r.Server.Location)
		if r.Server.Country != "" {
			serverInfo += ", " + r.Server.Country
		}
		serverInfo += ")"
	}

	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Server:          %s\n", serverInfo)
	fmt.Printf("Operators:       %s\n", r.ISP)
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("Latency:         %.2f ms  (抖动 %.2f ms)\n", r.Ping.Latency, r.Ping.Jitter)
	fmt.Printf("Pkg Loss Rate:   %s\n", packetLossStr)
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("Download Speed:  %.2f Mbps\n", downloadMbps)
	fmt.Printf("Upload Speed:    %.2f Mbps\n", uploadMbps)
	fmt.Println(strings.Repeat("=", 40))

	if r.Result.URL != "" {
		fmt.Printf("Detailed results: %s\n", r.Result.URL)
	}
}
