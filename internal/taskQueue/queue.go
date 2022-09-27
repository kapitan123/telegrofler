package taskQueue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"encoding/json"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/protobuf/types/known/timestamppb"

	log "github.com/sirupsen/logrus"
)

var defaultMessageLifeTime = 10

type TaskQueue struct {
	Name           string
	Path           string
	meta           meta
	client         *cloudtasks.Client
	ctx            context.Context
	initClientOnce sync.Once
	selfUrl        string
}

type meta interface {
	GetProjectId() string
	GetRegion() string
	GetEmail() string
}

func New(ctx context.Context, name string, meta meta, selfUrl string) *TaskQueue {
	return &TaskQueue{
		Name:           name,
		Path:           fmt.Sprintf("projects/%s/locations/%s/queues/%s", meta.GetProjectId(), meta.GetRegion(), name),
		meta:           meta,
		ctx:            ctx,
		initClientOnce: sync.Once{},
		selfUrl:        selfUrl,
	}
}

func (q *TaskQueue) EnqueueDeleteMessage(chatId int64, messageId int) error {
	var err error

	q.initClientOnce.Do(func() {
		err = q.initClient()
	})

	if err != nil {
		return err
	}

	delJson, err := json.Marshal(
		struct {
			ChatId    int64 `json:"chatId"`
			MessageId int   `json:"messageId"`
		}{
			chatId,
			messageId,
		},
	)

	if err != nil {
		return err
	}

	req := q.createPostRequest(q.selfUrl, delJson)

	createdTask, err := q.client.CreateTask(q.ctx, req)

	log.Info("Message deletion task was created ", createdTask.Name)

	if err != nil {
		return fmt.Errorf("cloudtasks.CreateTask has failed: %v", err)
	}

	return nil
}

func (q *TaskQueue) Close() error {
	if q.client == nil {
		return nil
	}

	return q.client.Close()
}

func (q *TaskQueue) initClient() error {
	var err error

	init := func() error {
		client, err := cloudtasks.NewClient(q.ctx)
		if err != nil {
			return fmt.Errorf("cloudtasks client was not created: %v", err)
		}

		q.client = client

		return nil
	}

	q.initClientOnce.Do(func() {
		err = init()
	})

	return err
}

func (q *TaskQueue) createPostRequest(url string, payload []byte) *taskspb.CreateTaskRequest {
	req := &taskspb.CreateTaskRequest{
		Parent: q.Path,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: q.meta.GetEmail(),
						},
					},
				},
			},
			ScheduleTime: getMinutesOffset(defaultMessageLifeTime),
		},
	}

	req.Task.GetHttpRequest().Body = payload

	return req
}

func getMinutesOffset(minutes int) *timestamppb.Timestamp {
	d := time.Minute * time.Duration(minutes)

	ts := &timestamppb.Timestamp{
		Seconds: time.Now().Add(d).Unix(),
	}

	return ts
}