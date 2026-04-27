package graph

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"net/url"
	"strings"

	"github.com/boatnoah/spidernet/internal/store"
	"github.com/goccy/go-graphviz"
	"github.com/google/uuid"
)

type GraphService struct {
	store store.Storage
}

type NormalizedNode struct {
	Key   string
	Label string
}

func New(store store.Storage) GraphService {
	return GraphService{store}
}

func (gs *GraphService) CreateGraph(ctx context.Context, jobID uuid.UUID, fileType string) ([]byte, error) {

	g, err := graphviz.New(ctx)

	if err != nil {
		return nil, err
	}

	defer g.Close()

	g.SetLayout(graphviz.TWOPI)

	graph, err := g.Graph()
	if err != nil {
		return nil, err
	}

	graph.SafeSet("splines", "line", "")
	graph.SafeSet("ranksep", "3.0", "")
	graph.SafeSet("overlap", "scale", "")

	links, err := gs.store.Links.GetAllLinksByJobID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	for _, link := range links {

		fromNode, err := normalizeURL(link.FromURL)
		if err != nil {
			return nil, err
		}

		toNode, err := normalizeURL(link.ToURL)
		if err != nil {
			return nil, err
		}

		node1, err := graph.CreateNodeByName(nodeID(fromNode.Key))
		if err != nil {
			return nil, err
		}
		node1.SetLabel(fromNode.Label)

		node2, err := graph.CreateNodeByName(nodeID(toNode.Key))
		if err != nil {
			return nil, err
		}
		node2.SetLabel(toNode.Label)

		_, err = graph.CreateEdgeByName(edgeID(fromNode.Key, toNode.Key), node1, node2)
		if err != nil {
			return nil, err
		}

	}

	var buf bytes.Buffer

	if err := g.Render(ctx, graph, graphviz.Format(fileType), &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func normalizeURL(raw string) (NormalizedNode, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return NormalizedNode{}, err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return NormalizedNode{}, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	host := strings.ToLower(u.Hostname())

	return NormalizedNode{
		Key:   host,
		Label: host,
	}, nil
}

func nodeID(key string) string {
	sum := sha1.Sum([]byte(key))
	return fmt.Sprintf("n_%x", sum[:8])
}

func edgeID(fromKey, toKey string) string {
	sum := sha1.Sum([]byte(fromKey + "->" + toKey))
	return fmt.Sprintf("e_%x", sum[:8])
}
