package graph

import (
	"bytes"
	"context"
	"log"

	"github.com/boatnoah/spidernet/internal/store"
	"github.com/goccy/go-graphviz"
	"github.com/google/uuid"
)

type GraphService struct {
	store store.Storage
}

func New(store store.Storage) GraphService {
	return GraphService{store}
}

func (gs *GraphService) CreateGraph(ctx context.Context, jobID uuid.UUID) (*bytes.Buffer, error) {
	g, err := graphviz.New(ctx)

	if err != nil {
		return nil, err
	}

	defer g.Close()

	graph, err := g.Graph()
	if err != nil {
		return nil, err
	}

	links, err := gs.store.Links.GetAllLinksByJobID(ctx, jobID)

	for _, link := range *links {
		node1, err := graph.CreateNodeByName(link.FromURL)
		if err != nil {
			return nil, err
		}

		node2, err := graph.CreateNodeByName(link.ToURL)
		if err != nil {
			return nil, err
		}

		_, err = graph.CreateEdgeByName(link.FromURL, node1, node2)
		if err != nil {
			return nil, err
		}

	}

	var buf bytes.Buffer

	if err := g.Render(ctx, graph, "dot", &buf); err != nil {
		log.Fatal(err)
	}

	return &buf, nil
}

// create convert functions for bytes.buffer -> png, svg, jpeg

func (gs *GraphService) toPNG(b bytes.Buffer) {
}

func (gs *GraphService) toJPEG(b bytes.Buffer) {
}

func (gs *GraphService) toSVG(b bytes.Buffer) {
}
