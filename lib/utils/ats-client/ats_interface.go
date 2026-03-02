package atsclient

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/time/delay"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"
)

type Executor struct {
	ctx            context.Context
	cancel         context.CancelFunc
	visibleTimeout time.Duration
	locateDelay    delay.Delayer
}

func NewExecutor(headless bool) (*Executor, error) {
	delayAmount := 100 * time.Millisecond

	conf := &provider.Conf{
		Headless: headless,
		HeadlessViewport: &provider.HeadlessViewport{
			Width:  1920,
			Height: 1080,
		},
		Delay:     delayAmount,
		DelayType: delay.Fixed,
	}

	ctx, cancel := provider.New(conf, application.Context)

	return &Executor{
		ctx:            ctx,
		cancel:         cancel,
		visibleTimeout: 30 * time.Second,
		locateDelay:    delay.New(delayAmount),
	}, nil
}

func NewExecutorWithConfig(conf *provider.Conf) (*Executor, error) {
	ctx, cancel := provider.New(conf, application.Context)

	var delayer delay.Delayer
	if conf.DelayType == delay.Random {
		delayer = delay.NewRandom(conf.Delay)
	} else {
		delayer = delay.New(conf.Delay)
	}

	return &Executor{
		ctx:            ctx,
		cancel:         cancel,
		visibleTimeout: 30 * time.Second,
		locateDelay:    delayer,
	}, nil
}

func (e *Executor) SetVisibleTimeout(timeout time.Duration) {
	e.visibleTimeout = timeout
}

func (e *Executor) WaitAndClick(selector string) error {
	return action.Click(e.ctx, e.visibleTimeout, e.locateDelay, selector, chromedp.ByQuery)
}

func (e *Executor) SetValue(selector, value string) error {
	return action.SetValue(e.ctx, e.visibleTimeout, e.locateDelay, selector, value, chromedp.ByQuery)
}

func (e *Executor) WaitForElement(selector string) error {
	return action.Locate(e.ctx, e.visibleTimeout, e.locateDelay, selector, chromedp.ByQuery)
}

func (e *Executor) Click(selector string) error {
	return action.Click(e.ctx, e.visibleTimeout, e.locateDelay, selector, chromedp.ByQuery)
}

func (e *Executor) Navigate(url string) error {
	return action.Execute(e.ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	)
}

func (e *Executor) RunTasks(actions ...chromedp.Action) error {
	return action.Execute(e.ctx, actions...)
}

func (e *Executor) Sleep(duration time.Duration) error {
	return action.Execute(e.ctx, chromedp.Sleep(duration))
}

func (e *Executor) GetText(selector string) (string, error) {
	return action.Get(e.ctx, selector)
}

func (e *Executor) GetAttribute(selector, attribute string) (string, error) {
	var value string
	err := action.Execute(e.ctx,
		chromedp.AttributeValue(selector, attribute, &value, nil, chromedp.ByQuery),
	)
	return value, err
}

func (e *Executor) Screenshot(quality int) ([]byte, error) {
	return action.TakeScreenshot(e.ctx)
}

func (e *Executor) Evaluate(script string, res any) error {
	return action.Execute(e.ctx,
		chromedp.Evaluate(script, res),
	)
}

func (e *Executor) SendKeys(selector, keys string) error {
	return action.SendKeys(e.ctx, e.visibleTimeout, e.locateDelay, selector, keys, chromedp.ByQuery)
}

func (e *Executor) SelectOption(selector, value string) error {
	return action.SetValue(e.ctx, e.visibleTimeout, e.locateDelay, selector, value, chromedp.ByQuery)
}

func (e *Executor) IsVisible(selector string) bool {
	return action.IsVisible(e.ctx, selector, "", chromedp.ByQuery)
}

func (e *Executor) WaitNotVisible(selector string) error {
	ctx, cancel := context.WithTimeout(e.ctx, e.visibleTimeout)
	defer cancel()

	return action.Execute(ctx,
		chromedp.WaitNotVisible(selector, chromedp.ByQuery),
	)
}

func (e *Executor) GetContext() context.Context {
	return e.ctx
}

func (e *Executor) Close() {
	if e.cancel != nil {
		e.cancel()
	}
}

type ATS interface {
	RegisterAccount(executor *Executor, account *Account) error
	LoginAccount(executor *Executor, email, password string) error
	ApplyForJob(executor *Executor, resumePath string, qaMap map[string]string, aiEnabled bool) error
	GetName() string
}

type Account struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
	Country   string
	Address   string
	City      string
	State     string
	ZipCode   string
}
