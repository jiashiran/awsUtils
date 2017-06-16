package dynomadb

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"strings"
	"strconv"
)

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

type DynoUtil struct {
	cfg Config
	svc *dynamodb.DynamoDB
}

func NewdynoUtil(region string, accessKey string, secretAccessKey string) *DynoUtil {
	cfg := Config{}
	cfg.Region = region
	cfg.Limit = 5
	svc := dynamodb.New(session.Must(session.NewSession(&aws.Config{
		Region: &region,
		Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretAccessKey,
		}),
	})))
	return &DynoUtil{cfg: cfg, svc: svc}
}

func (d *DynoUtil) Get(tableName string, param map[string]string) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		//IndexName:aws.String("cdr_enterprise_id"),
	}
	expressionAttributeNames := make(map[string]*dynamodb.AttributeValue)
	expressionAttributeNames[":cdrenterpriseId"] = &dynamodb.AttributeValue{N: aws.String(param["enterpriseId"])}
	expressionAttributeNames[":cdrTimeUniqueId"] = &dynamodb.AttributeValue{S: aws.String(param["cdrTimeUniqueId"])}
	//expressionAttributeNames[":cdrTimeUniqueIdStart"] = &dynamodb.AttributeValue{S: aws.String("1479818209.9")}

	params.ExpressionAttributeValues = expressionAttributeNames
	params.FilterExpression = aws.String("cdr_enterprise_id = :cdrenterpriseId and cdr_time_unique_id = :cdrTimeUniqueId")

	//params.FilterExpression = aws.String("cdr_enterprise_id = :cdrenterpriseId")
	if d.cfg.Limit > 0 {
		params.Limit = aws.Int64(d.cfg.Limit)
	}

	// Make the DynamoDB Query API call
	result, err := d.svc.Scan(params)
	if err != nil {
		exitWithError(fmt.Errorf("failed to make Query API call, %v", err))
	}

	// Unmarshal the Items field in the result value to the Item Go type.
	maps := result.Items
	for _, v := range maps {
		for key, value := range v {
			/*if key == "cdr_enterprise_id"{
				fmt.Print(key, " = ", value)
			}*/
			fmt.Print(key, " = ", value)
		}

	}
	//err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		exitWithError(fmt.Errorf("failed to unmarshal Query result items, %v", err))
	}

}

func (d *DynoUtil) Query(tableName string, param map[string]string) map[string]*string {
	attributeValueList := make(map[string]*dynamodb.AttributeValue)

	//keyConditions
	/*keyConditionsMap := make(map[string]*dynamodb.Condition)
	keyAttributeValueList := make([]*dynamodb.AttributeValue,1)
	keyConditionsMap[":a"] = &dynamodb.Condition{
		ComparisonOperator: aws.String("EQ"),
		AttributeValueList:append(keyAttributeValueList,&dynamodb.AttributeValue{N: aws.String("7000040"),
		})}
	keyAttributeValueList = make([]*dynamodb.AttributeValue,1)
	keyConditionsMap[":b"] = &dynamodb.Condition{
		ComparisonOperator: aws.String("EQ"),
		AttributeValueList:append(keyAttributeValueList,&dynamodb.AttributeValue{S: aws.String("1497597329.162682-sip-1"),
		})}
	queryInput.KeyConditions = keyConditionsMap
	*/

	queryInput := &dynamodb.QueryInput{TableName: aws.String(tableName)}
	queryInput.Limit = aws.Int64(d.cfg.Limit)
	keyConditionExpression := make([]string, 0)
	if param["enterpriseId"] != "" {
		attributeValueList[":cdrenterpriseId"] = &dynamodb.AttributeValue{N: aws.String(param["enterpriseId"])}
		keyConditionExpression = append(keyConditionExpression, "cdr_enterprise_id = :cdrenterpriseId ")
	}
	if param["cdrTimeUniqueId"] != "" {
		attributeValueList[":cdrTimeUniqueId"] = &dynamodb.AttributeValue{S: aws.String(param["cdrTimeUniqueId"])}
		keyConditionExpression = append(keyConditionExpression, " cdr_time_unique_id = :cdrTimeUniqueId ")
	}
	queryInput.KeyConditionExpression = aws.String(strings.Join(keyConditionExpression, "and"))

	if len(attributeValueList) > 0 {
		queryInput.ExpressionAttributeValues = attributeValueList
	}
	if param["scanIndexForward"] != ""{
		if b,err := strconv.ParseBool(param["scanIndexForward"]);err == nil{
			queryInput.ScanIndexForward = aws.Bool(b)//升序:true , 降序:false
		}
	}
	result, err := d.svc.Query(queryInput)
	if err != nil {
		exitWithError(fmt.Errorf("failed to make Query API call, %v", err))
	}
        resultMap := make(map[string]*string)
	maps := result.Items
	for _, v := range maps {
		for key, value := range v {
			resultMap[key] = value.S
		}

	}
	return resultMap
}

type Config struct {
	Table  string // required
	Region string // optional
	Limit  int64  // optional

}
