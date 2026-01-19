package main

import (
	"context"
	"fmt"

	"strings"

	"github.com/tmc/langchaingo/llms"

	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/qdrant"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func query(ctx context.Context, cmd *cli.Command) error {
	queryString := cmd.String("query")
	collectionName := cmd.String("collection")

	nDocs := cmd.Int("docs")
	threshold := cmd.Float64("threshold")

	log.Info().Msgf("Querying: %s [%d]@%f", queryString, nDocs, threshold)

	question := strings.ReplaceAll(queryString, "\"", "")
	question = strings.ToLower(question)


	store, err := qdrant.New(
		qdrant.WithURL(*urlAPI),
		qdrant.WithCollectionName(collectionName),
		qdrant.WithEmbedder(e),
	)
	logging.Error(err)

	docs, err := store.SimilaritySearch(ctx,
		question, nDocs,
		vectorstores.WithScoreThreshold(float32(threshold)))
	logging.Error(err)

	stringContext := ""
	for i := range len(docs) {
		stringContext += docs[i].PageContent
		log.Info().Msgf("using doc: %d -> %s [%f] - %s", i, docs[i].PageContent, docs[i].Score, docs[i].Metadata["file_path"])
	}

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a helpful assistant."),
		llms.TextParts(llms.ChatMessageTypeHuman, `Use the following pieces of context to answer the question at the end. If you don't know the answer, just say that you don't know, don't try to make up an answer.

			`+stringContext+`

			Question:`+question+`
			Helpful Answer:",
			"context","question")))
		`),
	}

	output, err := llm.GenerateContent(ctx, content,
		llms.WithMaxTokens(1024),
		llms.WithTemperature(0),
	)
	logging.Error(err)

	fmt.Println(output.Choices[0].Content)

	return nil
}
