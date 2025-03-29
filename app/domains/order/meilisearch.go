package order

import (
	"context"
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"github.com/p4xx07/order-service/configuration"
	"go.uber.org/zap"
	"log"
	"strconv"
)

type IMeilisearchService interface {
	List(ctx context.Context, request ListRequest) (interface{}, error)
	Update(orders Order) error
	Delete(orderIDs ...uint) error
}

type meilisearchService struct {
	configuration     *configuration.Configuration
	logger            *zap.SugaredLogger
	meilisearchClient meilisearch.ServiceManager
	store             IStore
}

func NewMeilisearchService(meilisearchClient meilisearch.ServiceManager, configuration *configuration.Configuration, logger *zap.SugaredLogger, store IStore) IMeilisearchService {
	s := &meilisearchService{meilisearchClient: meilisearchClient, configuration: configuration, logger: logger, store: store}
	go s.syncOrdersToMeili()
	return s
}

func (s *meilisearchService) List(ctx context.Context, request ListRequest) (interface{}, error) {
	index, err := s.meilisearchClient.GetIndexWithContext(ctx, "orders")
	if err != nil {
		s.logger.Errorw("error getting index", "error", err)
		return nil, err
	}

	attributes := s.getAttributes()
	_, err = index.UpdateFilterableAttributes(&attributes)
	if err != nil {
		s.logger.Errorw("error while updating meilisearch", "error", err)
		return nil, err
	}

	var filter string
	if request.StartDate != nil && request.EndDate != nil {
		startTimestamp := request.StartDate.UnixMilli()
		endTimestamp := request.EndDate.UnixMilli()
		filter = fmt.Sprintf("CreatedAtTimestamp >= %d AND CreatedAtTimestamp <= %d", startTimestamp, endTimestamp)
	} else if request.StartDate != nil {
		startTimestamp := request.StartDate.UnixMilli()
		filter = fmt.Sprintf("CreatedAtTimestamp >= %d", startTimestamp)
	} else if request.EndDate != nil {
		endTimestamp := request.EndDate.UnixMilli()
		filter = fmt.Sprintf("CreatedAtTimestamp <= %d", endTimestamp)
	}

	query := meilisearch.SearchRequest{
		Filter: filter,
		Limit:  request.Limit,
		Offset: request.Offset,
	}

	res, err := index.SearchWithContext(ctx, request.Input, &query)
	if err != nil {
		s.logger.Errorw("error while searching meilisearch", "error", err)
		return nil, err
	}

	return res.Hits, nil
}

func (s *meilisearchService) Delete(orderIDs ...uint) error {
	if len(orderIDs) == 0 {
		return nil
	}

	index := s.meilisearchClient.Index("orders")
	_, err := index.Delete(strconv.Itoa(int(orderIDs[0])))
	if err != nil {
		s.logger.Errorw("error while updating meilisearch", "error", err)
		return err
	}

	return nil
}

func (s *meilisearchService) Update(order Order) error {
	index := s.meilisearchClient.Index("orders")
	attributes := s.getAttributes()
	_, err := index.UpdateFilterableAttributes(&attributes)
	if err != nil {
		s.logger.Errorw("error while updating meilisearch", "error", err)
		return err
	}

	_, err = index.AddDocuments(order.toDocument(), "ID")
	if err != nil {
		s.logger.Errorw("error while updating meilisearch", "error", err)
		return err
	}

	return nil
}

func (s *meilisearchService) syncOrdersToMeili() {
	index := s.meilisearchClient.Index("orders")
	attributes := s.getAttributes()
	_, err := index.UpdateFilterableAttributes(&attributes)
	if err != nil {
		s.logger.Errorw("error while updating meilisearch", "error", err)
		return
	}

	stats, err := index.GetStats()
	if err != nil || stats.NumberOfDocuments == 0 {
		log.Println("Meilisearch empty, starting sync...")

		batchSize := 1000
		offset := 0
		for {
			orders, err := s.store.Fetch(batchSize, offset)
			if err != nil {
				s.logger.Errorw("error while fetching meilisearch orders", "error", err)
				return
			}

			if len(orders) == 0 {
				break
			}

			documents := make([]OrderMeilisearch, len(orders))
			for i, order := range orders {
				documents[i] = order.toDocument()
			}

			_, err = index.AddDocuments(documents, "ID")
			if err != nil {
				log.Println("Error syncing to Meilisearch:", err)
				return
			}

			offset += batchSize
		}
	}
}

func (s *meilisearchService) getAttributes() []string {
	return []string{"CreatedAtTimestamp", "Items.Product.Name", "Items.Product.Description"}
}
