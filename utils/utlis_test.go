package utils

import "testing"

func TestValidateString(t *testing.T) {
	str := `user_01JM6AAMNHBJX5TPPHEYCS8PKW%3A%3AeyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhdXRoMHx1c2VyXzAxSk02QUFNTkhCSlg1VFBQSEVZQ1M4UEtXIiwidGltZSI6IjE3Mzk2NzM5NjYiLCJyYW5kb21uZXNzIjoiYzdlNTQ4MDctNWVkZi00MmIyIiwiZXhwIjo0MzMxNjczOTY2LCJpc3MiOiJodHRwczovL2F1dGhlbnRpY2F0aW9uLmN1cnNvci5zaCIsInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwgb2ZmbGluZV9hY2Nlc3MiLCJhdWQiOiJodHRwczovL2N1cnNvci5jb20ifQ.mpV46iuHQrAezFS4LS6faBQAe7iYr7H2glT0YdWba9w`
	if _, _, ok := ValidateToken(str); !ok {
		t.Fatal("测试未通过")
	}
	t.Log("测试通过")
}
