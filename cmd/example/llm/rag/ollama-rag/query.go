package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

var (
	docsFlag      = flag.Int("docs", 5, "# of docs to include (RAG)")
	thresholdFlag = flag.Float64("threshold", 0.4, "threshold match, >= 0.4")
)

func query() {
	docsRetrieved := useRetriever()

	if len(docsRetrieved) == 0 {
		log.Warn().Msgf("No documents retrieved, err")
		return
	}

	history := memory.NewChatMessageHistory()
	for i := range docsRetrieved {
		err := history.AddAIMessage(application.Context, docsRetrieved[i].PageContent)
		logging.Warn(err, "history.AddAIMessage(ctx, doc.PageContent)")

		log.Info().Msgf("using doc: %d, %s [%f]", i, docsRetrieved[i].Metadata["file_path"], docsRetrieved[i].Score)
	}


	conversation := memory.NewConversationBuffer(memory.WithChatHistory(history))
	conversationalAgent := agents.NewConversationalAgent(
		llm, nil,
		agents.WithMemory(conversation),
	)

	executor := agents.NewExecutor(conversationalAgent)
	options := []chains.ChainCallOption{
		chains.WithTemperature(0.8),
	}

	res, err := chains.Run(application.Context, executor, *promptFlag, options...)
	logging.Panic(err)

	log.Info().Msg("result - after")
	fmt.Println(res)
}

func useRetriever() []schema.Document {
	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(float32(*thresholdFlag)),
	}

	retriever := vectorstores.ToRetriever(store, *docsFlag, optionsVector...)
	docsRetrieved, err := retriever.GetRelevantDocuments(context.Background(), *promptFlag)
	log.Info().Msgf("%v docs retrieved", len(docsRetrieved))
	logging.Panic(err)

	return docsRetrieved
}
