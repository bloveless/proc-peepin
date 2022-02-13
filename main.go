package main

import (
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/v3/process"
)

func main() {
	// Ignore any errors when loading the .env file just in case it doesn't exist
	_ = godotenv.Load()

	client := influxdb2.NewClient(os.Getenv("INFLUX_API"), os.Getenv("INFLUX_TOKEN"))
	defer client.Close()

	processes, err := process.Processes()
	if err != nil {
		panic(err)
	}

	log.Printf("Logging RSS and CPU from %d processes", len(processes))

	// get non-blocking write client
	writeAPI := client.WriteAPI(os.Getenv("INFLUX_EMAIL"), os.Getenv("INFLUX_BUCKET"))

	for _, process := range processes {
		name, err := process.Name()
		if err != nil {
			log.Printf("ERROR: Unable to get process name for PID: %d: %v", process.Pid, err)
			continue
		}

		memory, err := process.MemoryInfo()
		if err != nil {
			log.Printf("ERROR: Unable to get memory info for process %s: %v", name, err)
			continue
		}

		cpuTime, err := process.CPUPercent()
		if err != nil {
			log.Printf("ERROR: Unable to get cpu time for process %s: %v", name, err)
			continue
		}

		p := influxdb2.NewPointWithMeasurement("host_process").
			AddTag("host", os.Getenv("HOST")).
			AddTag("process", name).
			AddField("rss_mb", float64(memory.RSS)/1024/1024).
			AddField("cpu_time_percent", cpuTime).
			SetTime(time.Now())

		// write point asynchronously
		writeAPI.WritePoint(p)
	}

	// Flush writes
	writeAPI.Flush()

	log.Printf("Successfully pushed metrics to influxdb for %d processes", len(processes))
}
