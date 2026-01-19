package main

import (
	"context"

	"github.com/tmc/langchaingo/textsplitter"

	"os"
	"path/filepath"

	"github.com/tmc/langchaingo/vectorstores/qdrant"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/schema"
	"github.com/urfave/cli/v3"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"strings"
)

func index(ctx context.Context, cmd *cli.Command) error {
	fileName := cmd.String("file")
	collectionName := cmd.String("collection")

	return indexDo(ctx, collectionName, fileName)
}

func indexDo(ctx context.Context, collectionName string, fileName string) error {
	log.Info().Msgf("Indexing file: %s", fileName)

	info, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	reqCharacterSplitter := textsplitter.NewRecursiveCharacter()
	reqCharacterSplitter.ChunkSize = 1000
	reqCharacterSplitter.ChunkOverlap = 200
	reqCharacterSplitter.LenFunc = func(s string) int { return len(s) }

	store, err := qdrant.New(
		qdrant.WithURL(*urlAPI),
		qdrant.WithCollectionName(collectionName),
		qdrant.WithEmbedder(e),
	)
	logging.Error(err)

	if info.IsDir() {
		err := traverse(ctx, reqCharacterSplitter, store, fileName)
		if err != nil {
			return err
		}

		return nil
	}

	indexDocument(ctx, reqCharacterSplitter, store, fileName)
	return nil
}

func traverse(ctx context.Context, reqCharacterSplitter textsplitter.RecursiveCharacter, store qdrant.Store, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if strings.Contains(path, ".git") {
			log.Debug().Msgf("ignoring git: %s", path)
			return nil
		}

		indexDocument(ctx, reqCharacterSplitter, store, path)
		return nil
	})
}

func indexDocument(ctx context.Context, reqCharacterSplitter textsplitter.RecursiveCharacter, store qdrant.Store, fileName string) {
	log.Info().Msgf("indexing: %s", fileName)

	var pagesList []schema.Document
	var err error
	if strings.HasSuffix(fileName, ".pdf") {
		pagesList, err = indexPDF(fileName)
	} else {
		pagesList, err = indexText(fileName)
	}

	logging.Error(err)


	chunksDocList, _ := textsplitter.SplitDocuments(reqCharacterSplitter, pagesList)

	_, err = store.AddDocuments(ctx, chunksDocList)
	logging.Warn(err, "indexDocument")
}
