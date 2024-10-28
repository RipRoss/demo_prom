package seed

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AVIRecord struct {
	Id string `bson:"_id"`
	TrackingId string `bson:"trackingId"`
	ApplicationId string `bson:"applicationId"`
	ApplicationName string `bson:"applicationName"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	SiemRef []string `bson:"siemRef"`
} 

type SIEMRecord struct {
	Id string `bson:"_id"`
	TrackingId string `bson:"trackingId"`
	ApplicationId string `bson:"applicationId"`
	ApplicationName string `bson:"applicationName"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	CreateSaimCase string `bson:"createSaimCase"`
	AviRef []string `bson:"aviRef"`	
	SaimRef []string `bson:"saimRef"`
}

type SAIMRecord struct {
	Id string `bson:"_id"`
	TrackingId string `bson:"trackingId"`
	ApplicationId string `bson:"applicationId"`
	ApplicationName string `bson:"applicationName"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	SiemRef []string `bson:"siemRef"`
}

type UsecaseExecution struct {
	Id string `bson:"_id"`
	TrackingId string `bson:"trackingId"`
	ApplicationId string `bson:"applicationId"`
	ApplicationName string `bson:"applicationName"`
	CreatedAt time.Time `bson:"createdAt"`
	Duration int `bson:"duration"`
	ReferenceId string `bson:"referenceId"`
}

func GenerateDatabaseData() {
	/*
	We need to iterate over the number of x passed in as a parameter
	For each one, we need to generate an AVI record
	For each AVI record, we need to generate a SIEM record
	For eachn SIEM record, we need to generate a SAIM record
	We need to link the records together
	We need to insert the records into the database
	The created_at time needs to be a random number of minutes different, to simulate it taking time for records to go from AVI record, to SIEM record and to SAIM record.
	*/
	days := 1
	recordsPerDay := 400000
	endDate := time.Date(2024, 10, 25, 12, 0, 0, 0, time.UTC)
	startDate := endDate.AddDate(0, 0, -days)
	useCaseExecDate := startDate.AddDate(0, 0, -days)

	rand.New(rand.NewSource(time.Now().UnixNano())) // Seed the random number generator

	aviRecords := make([]AVIRecord, 0)
	siemRecords := make([]SIEMRecord, 0)
	saimRecords := make([]SAIMRecord, 0)
	useCaseExecs := make([]UsecaseExecution, 0)
	
	random := rand.New(rand.NewSource(time.Now().UnixNano())) // Seed the random number generator

	for startDate.Before(endDate) {
		endOfDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 23, 59, 59, 0, time.UTC)

		for i := 0; i < days; i++ {
			for i := 0; i < recordsPerDay; i++ {
				createSaimCase := "YES"

				if i % 3 == 0 {
					createSaimCase = "NO"
				}

				for x := 0; x < 2; x++ {
					newUUID := uuid.New()

					randomSeconds := 1 + random.Intn(2) // This gives 1 or 2 seconds
					randomMilliseconds := random.Intn(1000) // Random milliseconds between 0 and 999
					randomDuration := time.Duration(randomSeconds)*time.Second + time.Duration(randomMilliseconds)*time.Millisecond

					useCaseRandomSeconds := 1 + random.Intn(2) // This gives 1 or 2 seconds
					useCaseRandomMilliseconds := random.Intn(1000) // Random milliseconds between 0 and 999
					useCaseRandomDuration := time.Duration(useCaseRandomSeconds)*time.Second + time.Duration(useCaseRandomMilliseconds)*time.Millisecond

					// to fix the issue that we have, we need to change this to incorporate the fact that we also need to flick the day over to the next day when it hits 24 hours
					startDate = startDate.Add(randomDuration)
					useCaseExecDate = useCaseExecDate.Add(useCaseRandomDuration)
					
					// Generate AVI record
					aviRecord := AVIRecord{
						Id: fmt.Sprintf("avi_%s", newUUID.String()), // i needs to be more random, we are stuck at 5000 because of recordsPerDay
						TrackingId: fmt.Sprintf("avi_tracking_%s", newUUID.String()),
						ApplicationId: "avi_application",
						ApplicationName: "avi_application_name",
						CreatedAt:       startDate, // Use the random duration
						UpdatedAt: startDate,
						SiemRef: []string{},
					}

					siemRecord := SIEMRecord{
						Id: fmt.Sprintf("siem%s", newUUID.String()),
						TrackingId: fmt.Sprintf("siem_tracking_%s", newUUID.String()),
						ApplicationId: "siem_application",
						ApplicationName: "siem_application_name",
						CreatedAt:       startDate, // Use the random duration
						UpdatedAt: startDate,
						AviRef: []string{aviRecord.Id},
						CreateSaimCase: createSaimCase,
						SaimRef: []string{},
					}

					aviUseCase := UsecaseExecution{
						Id: fmt.Sprintf("usecase_%s", newUUID.String()),
						TrackingId: aviRecord.TrackingId,
						ApplicationId: "avi_application",
						ApplicationName: "avi_application_name",
						CreatedAt: useCaseExecDate,
						ReferenceId: aviRecord.Id,
						Duration: 1 + random.Intn(60),
					}

					if siemRecord.CreateSaimCase == "YES" {
						saimRecord := SAIMRecord{
							Id: fmt.Sprintf("saim_%s", newUUID.String()),
							TrackingId: fmt.Sprintf("saim_tracking_%s", newUUID.String()),
							ApplicationId: "saim_application",
							ApplicationName: "saim_application_name",
							CreatedAt:       startDate, // Use the random duration
							UpdatedAt: startDate,
							SiemRef: []string{siemRecord.Id},
						}

						siemRecord.SaimRef = append(siemRecord.SaimRef, saimRecord.Id)
						saimRecords = append(saimRecords, saimRecord)
					}

					// Link the records together
					aviRecord.SiemRef = append(aviRecord.SiemRef, siemRecord.Id)

					aviRecords = append(aviRecords, aviRecord)
					siemRecords = append(siemRecords, siemRecord)
					useCaseExecs = append(useCaseExecs, aviUseCase)
				}
				

				if startDate.After(endOfDay) {
					// Move to midnight of the next day
					startDate = endOfDay.Add(time.Second)
					break // Exit loop for the current day and move to the next day
				}
			}
		}
	}

	client := getMongoConnection()
	writeRecordsToDatabase(client, Records{aviRecords, siemRecords, saimRecords, useCaseExecs})
}

type Records struct {
	aviRecords  []AVIRecord
	siemRecords []SIEMRecord
	saimRecords []SAIMRecord
	useCaseExecs []UsecaseExecution
}

func writeRecordsToDatabase(mConn *MongoConn, records Records) {
	// Get the database and collection
	defer func() {
		if err := mConn.client.Disconnect(context.TODO()); err != nil {
				log.Fatal(err)
		}
	}()

	aviCollection := mConn.database.Collection("avi")
	siemCollection := mConn.database.Collection("siem")
	saimCollection := mConn.database.Collection("saim")
	useCases := mConn.database.Collection("use_cases")

	aviCollection.DeleteMany(context.TODO(), struct{}{})
	siemCollection.DeleteMany(context.TODO(), struct{}{})
	saimCollection.DeleteMany(context.TODO(), struct{}{})
	useCases.DeleteMany(context.TODO(), struct{}{})

	// Insert the records
	aviInterfaces := make([]interface{}, len(records.aviRecords))
	for i, v := range records.aviRecords {
		aviInterfaces[i] = v
	}

	siemInterfaces := make([]interface{}, len(records.siemRecords))
	for i, v := range records.siemRecords {
		siemInterfaces[i] = v
	}

	saimInterfaces := make([]interface{}, len(records.saimRecords))
	for i, v := range records.saimRecords {
		saimInterfaces[i] = v
	}	

	useCaseInterfaces := make([]interface{}, len(records.useCaseExecs))
	for i, v := range records.useCaseExecs {
		useCaseInterfaces[i] = v
	}

	aviCollection.InsertMany(context.TODO(), aviInterfaces)
	siemCollection.InsertMany(context.TODO(), siemInterfaces)
	saimCollection.InsertMany(context.TODO(), saimInterfaces)
	useCases.InsertMany(context.TODO(), useCaseInterfaces)
}

type MongoConn struct {
	client *mongo.Client
	database *mongo.Database
}

func getMongoConnection() *MongoConn {
	uri := "mongodb://admin:admin@localhost:27017"

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	database := client.Database("monitoring_demo")
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return &MongoConn{client: client, database: database}
}