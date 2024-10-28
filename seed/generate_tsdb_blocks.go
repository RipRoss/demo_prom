package seed

import (
	"context"
	"fmt"
	"log"
	"os"
)

type Metric struct {
	Name string
	Value float64
	Timestamp int64
	Labels map[string]string
}

func GenerateTimeSeriesData() {
	// openmetrics format
	mConn := getMongoConnection()

	defer func() {
		if err := mConn.client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	aviRecords := getAVIRecords(mConn)
	siemRecords := getSIEMRecords(mConn)
	saimRecords := getSAIMRecords(mConn)
	useCaseExecs := getUseCaseRecords(mConn)

	aviMetrics := getAVIMetrics(aviRecords)
	siemMetrics := getSIEMMetrics(siemRecords)
	saimMetrics := getSAIMMetrics(saimRecords)
	useCaseMetrics := getUsecaseMetrics(useCaseExecs)

	writeMetricsToFile(aviMetrics, siemMetrics, saimMetrics, useCaseMetrics)

	// GenerateTSDBBlocks()
}

func labelsToString(labels map[string]string) string {
	labelStr := ""
	for key, value := range labels {
		labelStr += fmt.Sprintf("%s=\"%s\",", key, value)
	}
	// Remove the trailing comma if there are labels
	if len(labelStr) > 0 {
		labelStr = labelStr[:len(labelStr)-1]
	}
	return labelStr
}

func writeMetricsToFile(aviMetrics []Metric, siemMetrics []Metric, saimMetrics []Metric, usecaseMetrics []Metric) error {
	file, err := os.Create("metrics.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header comments
	_, err = file.WriteString("# HELP my_random_metric This is a random metric\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString("# TYPE my_random_metric gauge\n")
	if err != nil {
		return err
	}

	// Write each metric
	for _, metric := range aviMetrics {
		_, err = file.WriteString(fmt.Sprintf("%s{%s} %f %d\n",
			metric.Name,
			labelsToString(metric.Labels),
			metric.Value,
			metric.Timestamp))
		if err != nil {
			return err
		}
	}

	for _, metric := range siemMetrics {
		_, err = file.WriteString(fmt.Sprintf("%s{%s} %f %d\n",
			metric.Name,
			labelsToString(metric.Labels),
			metric.Value,
			metric.Timestamp))
		if err != nil {
			return err
		}
	}

	for _, metric := range saimMetrics {
		_, err = file.WriteString(fmt.Sprintf("%s{%s} %f %d\n",
			metric.Name,
			labelsToString(metric.Labels),
			metric.Value,
			metric.Timestamp))
		if err != nil {
			return err
		}
	}

	for _, metric := range usecaseMetrics {
		// this is going to be different to the others as we are doing this as a histogram

		_, err = file.WriteString(fmt.Sprintf("%s_bucket{%s} %f %d\n",
			metric.Name,
			labelsToString(metric.Labels),
			metric.Value,
			metric.Timestamp,
		))

		if err != nil {
			return err
		}	
	}

	// Add EOF at the end of the file
	_, err = file.WriteString("# EOF\n")
	if err != nil {
		return err
	}

	return nil
}

func getAVIMetrics(records []AVIRecord) []Metric {
	counter := 0

	metrics := []Metric{}

	for _, record := range records {
		counter++
		
		metrics = append(metrics, Metric{
			Name: "record_ingest",
			Value: float64(counter),
			Timestamp: record.CreatedAt.Unix(),
			Labels: map[string]string{
				"appName": record.ApplicationName,
			},
		})
	}

	return metrics
}

func getSIEMMetrics(records []SIEMRecord) []Metric {
	counter := 0

	metrics := []Metric{}

	for _, record := range records {
		counter++
		
		metrics = append(metrics, Metric{
			Name: "record_ingest",
			Value: float64(counter),
			Timestamp: record.CreatedAt.Unix(),
			Labels: map[string]string{
				"appName": record.ApplicationName,
			},
		})
	}

	return metrics
}

func getSAIMMetrics(records []SAIMRecord) []Metric {
	counter := 0

	metrics := []Metric{}

	for _, record := range records {
		counter++
		
		metrics = append(metrics, Metric{
			Name: "record_ingest",
			Value: float64(counter),
			Timestamp: record.CreatedAt.Unix(),
			Labels: map[string]string{
				"appName": record.ApplicationName,
			},
		})
	}

	return metrics
}

func getUsecaseMetrics(records []UsecaseExecution) []Metric {
	// this is going to be refactored because it needs to accommodate the fact that it is a histogram.
	// buckets := []float64{1, 5, 10, 20, 30, 40, 50, 60}

	bucketMap := map[float64]float64{
		1: 0,
		5: 0,
		10: 0,
		20: 0,
		30: 0,
		40: 0,
		50: 0,
		60: 0,
	}

	metrics := []Metric{}

	for key := range bucketMap {
		for _, record := range records {
			if float64(record.Duration) <= key {
				bucketMap[key]++

				metrics = append(metrics, Metric{
					Name: "use_case_execs",
					Value: bucketMap[key],
					Timestamp: record.CreatedAt.Unix(),
					Labels: map[string]string{
						"le": fmt.Sprintf("%f", key),  
						"appName": record.ApplicationName,
					},
				})
			}
		}
	}

	return metrics
}

func getAVIRecords(mConn *MongoConn) []AVIRecord {
	aviCollection := mConn.database.Collection("avi")
	cursor, err := aviCollection.Find(context.TODO(), struct{}{})
	if err != nil {
		log.Fatal(err)
	}

	var records []AVIRecord
	if err = cursor.All(context.TODO(), &records); err != nil {
		log.Fatal(err)
	}

	return records
}

func getSIEMRecords(mConn *MongoConn) []SIEMRecord {
	siemCollection := mConn.database.Collection("siem")
	cursor, err := siemCollection.Find(context.TODO(), struct{}{})
	if err != nil {
		log.Fatal(err)
	}

	var records []SIEMRecord
	if err = cursor.All(context.TODO(), &records); err != nil {
		log.Fatal(err)
	}

	return records
}

func getSAIMRecords(mConn *MongoConn) []SAIMRecord {
	saimCollection := mConn.database.Collection("saim")
	cursor, err := saimCollection.Find(context.TODO(), struct{}{})
	if err != nil {
		log.Fatal(err)
	}

	var records []SAIMRecord
	if err = cursor.All(context.TODO(), &records); err != nil {
		log.Fatal(err)
	}

	return records
}

func getUseCaseRecords(mConn *MongoConn) []UsecaseExecution {
	useCasesCollection := mConn.database.Collection("use_cases")
	cursor, err := useCasesCollection.Find(context.TODO(), struct{}{})
	if err != nil {
		log.Fatal(err)
	}

	var records []UsecaseExecution
	if err = cursor.All(context.TODO(), &records); err != nil {
		log.Fatal(err)
	}

	return records
}