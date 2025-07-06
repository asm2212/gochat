package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type ChatService struct {
	rdb *redis.Client
}

func NewChatService(rdb *redis.Client) *ChatService {
	return &ChatService{rdb: rdb}
}

// --- Direct Messages ---

// SendDirectMessage stores a direct message between two users
func (s *ChatService) SendDirectMessage(ctx context.Context, from, to, content string) error {
	msg := Message{
		ID:        uuid.NewString(),
		From:      from,
		To:        to,
		Content:   content,
		Type:      "direct",
		Timestamp: time.Now(),
	}
	key := directKey(from, to)
	return s.rdb.RPush(ctx, key, serializeMsg(msg)).Err()
}

// GetDirectMessages fetches messages between two users
func (s *ChatService) GetDirectMessages(ctx context.Context, user1, user2 string) ([]Message, error) {
	key := directKey(user1, user2)
	msgs, err := s.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return deserializeMsgs(msgs), nil
}

func directKey(u1, u2 string) string {
	if u1 < u2 {
		return fmt.Sprintf("dm:%s:%s", u1, u2)
	}
	return fmt.Sprintf("dm:%s:%s", u2, u1)
}

// --- Group ---

// CreateGroup creates a group
func (s *ChatService) CreateGroup(ctx context.Context, group string) error {
	return s.rdb.SAdd(ctx, "groups", group).Err()
}

// SendGroupMessage stores a message to a group
func (s *ChatService) SendGroupMessage(ctx context.Context, from, group, content string) error {
	exists, _ := s.rdb.SIsMember(ctx, "groups", group).Result()
	if !exists {
		return errors.New("group does not exist")
	}
	msg := Message{
		ID:        uuid.NewString(),
		From:      from,
		Group:     group,
		Content:   content,
		Type:      "group",
		Timestamp: time.Now(),
	}
	key := fmt.Sprintf("group:%s", group)
	return s.rdb.RPush(ctx, key, serializeMsg(msg)).Err()
}

// GetGroupMessages fetches group chat history
func (s *ChatService) GetGroupMessages(ctx context.Context, group string) ([]Message, error) {
	key := fmt.Sprintf("group:%s", group)
	msgs, err := s.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return deserializeMsgs(msgs), nil
}

// --- Broadcast ---

// Broadcast stores a global message
func (s *ChatService) Broadcast(ctx context.Context, from, content string) error {
	msg := Message{
		ID:        uuid.NewString(),
		From:      from,
		Content:   content,
		Type:      "broadcast",
		Timestamp: time.Now(),
	}
	return s.rdb.RPush(ctx, "broadcast", serializeMsg(msg)).Err()
}

// GetBroadcasts fetches all broadcast messages
func (s *ChatService) GetBroadcasts(ctx context.Context) ([]Message, error) {
	msgs, err := s.rdb.LRange(ctx, "broadcast", 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return deserializeMsgs(msgs), nil
}

// --- Serialization helpers ---

func serializeMsg(msg Message) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%d", msg.ID, msg.From, msg.To, msg.Group, msg.Content, msg.Type, msg.Timestamp.Unix())
}

func deserializeMsgs(ss []string) []Message {
	var out []Message
	for _, s := range ss {
		var msg Message
		var ts int64
		fmt.Sscanf(s, "%[^|]|%[^|]|%[^|]|%[^|]|%[^|]|%[^|]|%d",
			&msg.ID, &msg.From, &msg.To, &msg.Group, &msg.Content, &msg.Type, &ts)
		msg.Timestamp = time.Unix(ts, 0)
		out = append(out, msg)
	}
	return out
}
