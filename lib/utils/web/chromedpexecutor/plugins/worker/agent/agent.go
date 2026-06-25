package agent

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"time"
)

func (c *Conf) Name() string {
	return "agent"
}

func (c *Conf) Work(ctx context.Context, headless bool) {
	err := c.ask(ctx, c.getQuestion())
	if err != nil {
		log.Warn().Msgf("agent.Work.error: %v", err)
		return
	}

	c.waitForAnswer()
}

func (c *Conf) getQuestion() string {
	num, err := rand.Int(rand.Reader, big.NewInt(int64(len(c.questions))))
	if err != nil {
		log.Error().Msgf("failed to generate random index: %v", err)
		return c.questions[0]
	}
	randomIndex := int(num.Int64())
	log.Info().Msgf("selected question: [%d]", randomIndex)

	return c.questions[randomIndex]
}

func (c *Conf) ask(ctx context.Context, question string) error {
	log.Info().Msgf("%v.ask[%d] - %s", c, c.iteration, question)
	c.iteration++

	return action.Execute(ctx,
		chromedp.KeyEvent(question),
		chromedp.Sleep(200*time.Millisecond),
		chromedp.KeyEvent(kb.Enter))
}

func (c *Conf) waitForAnswer() {


	time.Sleep(30 * time.Second)
}

func (c *Conf) Cleanup() {
}
