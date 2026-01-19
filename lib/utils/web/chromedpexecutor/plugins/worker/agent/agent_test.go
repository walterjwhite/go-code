package agent

import (
	"context"
	"testing"
)

func TestName(t *testing.T) {
	c := &Conf{}
	expected := "agent"
	if c.Name() != expected {
		t.Errorf("expected %s, got %s", expected, c.Name())
	}
}

func TestGetQuestion(t *testing.T) {
	questions := []string{"How are you?", "What is Go?", "Is it raining?"}
	c := &Conf{
		questions: questions,
	}

	for i := 0; i < 10; i++ {
		q := c.getQuestion()
		found := false
		for _, v := range questions {
			if v == q {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("getQuestion() returned a value not in the original slice: %s", q)
		}
	}
}

func TestIterationIncrement(t *testing.T) {

	c := &Conf{
		iteration: 0,
		questions: []string{"test"},
	}

	ctx := context.Background()

	question := c.getQuestion()

	_ = c.ask(ctx, question)

	if c.iteration != 1 {
		t.Errorf("expected iteration to be 1, got %d", c.iteration)
	}
}
