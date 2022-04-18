package massdriver_test

// type SNSTestClient struct {
// 	Input *sns.PublishInput
// 	Data  *string
// }

// func (c *SNSTestClient) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
// 	c.Input = params
// 	c.Data = params.Message
// 	return &sns.PublishOutput{}, nil
// }

// func TestPublishEventToSNS(t *testing.T) {

// 	type testData struct {
// 		name  string
// 		want  string
// 	}
// 	tests := []testData{
// 		{
// 			name: "Test Decommission Failed",
// 			want: `{"metadata":{"timestamp":"2021-01-01 12:00:00.4321","provisioner":"testaform","event_type":"create_pending"},"payload":{"deployment_id":"depId"}}`,
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			environmentDepId := "envDepId"
// 			t.Setenv("MASSDRIVER_DEPLOYMENT_ID", environmentDepId)
// 			testClient, err := massdriver.InitializeMassdriverClient()
// 			if err != nil {
// 				t.Fatalf("%d, unexpected error", err)
// 			}
// 			testSNSClient := SNSTestClient{}
// 			testClient.SNSClient = &testSNSClient
// 			err = testClient.PublishEventToSNS(tc.input)
// 			if err != nil {
// 				t.Fatalf("%d, unexpected error", err)
// 			}

// 			got := testSNSClient.Input
// 			if *got.Message != tc.want {
// 				t.Fatalf("want: %v, got: %v", tc.want, *got.Message)
// 			}
// 			if *got.MessageGroupId != environmentDepId {
// 				t.Fatalf("want: %v, got: %v", environmentDepId, *got.MessageGroupId)
// 			}
// 		})
// 	}
// }
